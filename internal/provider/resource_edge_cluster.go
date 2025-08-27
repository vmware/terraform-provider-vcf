// © Broadcom. All Rights Reserved.
// The term “Broadcom” refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/vmware/vcf-sdk-go/vcf"

	"github.com/vmware/terraform-provider-vcf/internal/api_client"
	"github.com/vmware/terraform-provider-vcf/internal/nsx_edge_cluster"
	validationUtils "github.com/vmware/terraform-provider-vcf/internal/validation"
)

const (
	shrinkage = "SHRINKAGE"
	expansion = "EXPANSION"
)

func ResourceEdgeCluster() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceNsxEdgeClusterCreate,
		ReadContext:   resourceNsxEdgeClusterRead,
		UpdateContext: resourceNsxEdgeClusterUpdate,
		DeleteContext: resourceNsxEdgeClusterDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(180 * time.Minute),
			Update: schema.DefaultTimeout(180 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				Description:  "The name of the edge cluster",
				ValidateFunc: validation.NoZeroValues,
			},
			"root_password": {
				Type:         schema.TypeString,
				Required:     true,
				Description:  "Root user password for the NSX manager",
				ValidateFunc: validationUtils.ValidateNsxEdgePassword,
			},
			"admin_password": {
				Type:         schema.TypeString,
				Required:     true,
				Description:  "Administrator password for the NSX manager",
				ValidateFunc: validationUtils.ValidateNsxEdgePassword,
			},
			"audit_password": {
				Type:         schema.TypeString,
				Required:     true,
				Description:  "Audit user password for the NSX manager",
				ValidateFunc: validationUtils.ValidateNsxEdgePassword,
			},
			"tier0_name": {
				Type:         schema.TypeString,
				Optional:     true,
				Description:  "Name for the Tier-0 gateway",
				ValidateFunc: validation.NoZeroValues,
			},
			"tier1_name": {
				Type:         schema.TypeString,
				Optional:     true,
				Description:  "Name for the Tier-1 gateway",
				ValidateFunc: validation.NoZeroValues,
			},
			"profile_type": {
				Type:         schema.TypeString,
				Required:     true,
				Description:  "One among: DEFAULT, CUSTOM. If set to CUSTOM a 'profile' must be provided",
				ValidateFunc: validation.StringInSlice([]string{"DEFAULT", "CUSTOM"}, false),
			},
			"profile": {
				Type:        schema.TypeList,
				Optional:    true,
				MaxItems:    1,
				Description: "The specification for the edge cluster profile",
				Elem:        nsx_edge_cluster.ClusterProfileSchema(),
			},
			"routing_type": {
				Type:         schema.TypeString,
				Optional:     true,
				Description:  "One among: EBGP, STATIC",
				ValidateFunc: validation.StringInSlice([]string{"EBGP", "STATIC"}, false),
			},
			"form_factor": {
				Type:         schema.TypeString,
				Required:     true,
				Description:  "One among: XLARGE, LARGE, MEDIUM, SMALL",
				ValidateFunc: validation.StringInSlice([]string{"XLARGE", "LARGE", "MEDIUM", "SMALL"}, false),
			},
			"high_availability": {
				Type:         schema.TypeString,
				Optional:     true,
				Description:  "One among: ACTIVE_ACTIVE, ACTIVE_STANDBY",
				ValidateFunc: validation.StringInSlice([]string{"ACTIVE_ACTIVE", "ACTIVE_STANDBY"}, false),
			},
			"mtu": {
				Type:         schema.TypeInt,
				Required:     true,
				Description:  "Maximum transmission unit size for the cluster",
				ValidateFunc: validation.IntBetween(1600, 9000),
			},
			"asn": {
				Type:         schema.TypeString,
				Optional:     true,
				Description:  "ASN for the cluster",
				ValidateFunc: validationUtils.ValidASN,
			},
			"skip_tep_routability_check": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Set to true to bypass normal ICMP-based check of Edge TEP / host TEP routability (default is false, meaning do check)",
				Default:     false,
			},
			"tier1_unhosted": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Select whether Tier-1 being created per this spec is hosted on the new Edge cluster or not (default value is false, meaning hosted)",
				Default:     false,
			},
			"internal_transit_subnets": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Subnet addresses in CIDR notation that are used to assign addresses to logical links connecting service routers and distributed routers",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"transit_subnets": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Transit subnet addresses in CIDR notation that are used to assign addresses to logical links connecting Tier-0 and Tier-1s",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"edge_node": {
				Type:        schema.TypeList,
				Required:    true,
				Description: "The nodes in the edge cluster",
				Elem:        nsx_edge_cluster.EdgeNodeSchema(),
			},
		},
	}
}

func resourceNsxEdgeClusterCreate(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*api_client.SddcManagerClient).ApiClient

	spec, err := nsx_edge_cluster.GetNsxEdgeClusterCreationSpec(data, client)
	if err != nil {
		return diag.FromErr(err)
	}

	validationErr := validateClusterCreationSpec(client, ctx, *spec)
	if validationErr != nil {
		return validationErr
	}

	res, err := client.CreateEdgeClusterWithResponse(ctx, *spec)
	if err != nil {
		return diag.FromErr(err)
	}
	task, vcfErr := api_client.GetResponseAs[vcf.Task](res)
	if vcfErr != nil {
		api_client.LogError(vcfErr, ctx)
		return diag.FromErr(errors.New(*vcfErr.Message))
	}

	tflog.Info(ctx, "Edge cluster creation has started.")
	if err = api_client.NewTaskTracker(ctx, client, *task.Id).WaitForTask(); err != nil {
		return diag.FromErr(err)
	}

	const maxRetries = 30
	const retryDelay = 10 * time.Second

	for attempt := 0; attempt < maxRetries; attempt++ {
		tflog.Info(ctx, fmt.Sprintf("Attempt %d/%d to find edge cluster", attempt+1, maxRetries))

		clusters, err := client.GetEdgeClustersWithResponse(ctx, nil)
		if err != nil {
			return diag.FromErr(err)
		}
		page, vcfErr := api_client.GetResponseAs[vcf.PageOfEdgeCluster](clusters)
		if vcfErr != nil {
			api_client.LogError(vcfErr, ctx)
			return diag.FromErr(errors.New(*vcfErr.Message))
		}

		tflog.Info(ctx, fmt.Sprintf("Looking for cluster name: '%s'", data.Get("name").(string)))
		tflog.Info(ctx, fmt.Sprintf("Total clusters found in API response: %d", len(*page.Elements)))

		for _, cluster := range *page.Elements {
			if cluster.Name != nil && *cluster.Name == data.Get("name") {
				data.SetId(*cluster.Id)
				tflog.Info(ctx, "Edge cluster created successfully.")
				return nil
			}
		}

		// Sleep after all attempts other than the last one.
		if attempt < maxRetries-1 {
			time.Sleep(retryDelay)
		}
	}

	return diag.Errorf("Edge cluster creation failed - cluster not found in inventory")
}

func resourceNsxEdgeClusterRead(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*api_client.SddcManagerClient).ApiClient

	_, err := client.GetEdgeClusterWithResponse(ctx, data.Id())

	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func resourceNsxEdgeClusterDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	tflog.Info(ctx, "Edge cluster deletion is not implemented. See KB article 78635 for more information.")
	return nil
}

func resourceNsxEdgeClusterUpdate(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*api_client.SddcManagerClient).ApiClient

	edgeClusterOk, err := client.GetEdgeClusterWithResponse(ctx, data.Id())

	if err != nil {
		return diag.FromErr(err)
	}

	resp, vcfErr := api_client.GetResponseAs[vcf.EdgeCluster](edgeClusterOk)
	if vcfErr != nil {
		api_client.LogError(vcfErr, ctx)
		return diag.FromErr(errors.New(*vcfErr.Message))
	}

	if data.HasChange("edge_node") {
		oldNodesRaw, newNodesRaw := data.GetChange("edge_node")
		oldNodes, newNodes := oldNodesRaw.([]interface{}), newNodesRaw.([]interface{})

		if len(oldNodes) == len(newNodes) {
			return diag.Errorf("Adding and removing edge nodes is not supported in a single configuration change. Apply each change separately.")
		}

		updateSpec := vcf.EdgeClusterUpdateSpec{}

		// Shrink
		if len(oldNodes) > len(newNodes) {
			operation := shrinkage
			updateSpec.Operation = operation
			updateSpec.EdgeClusterShrinkageSpec =
				nsx_edge_cluster.GetNsxEdgeClusterShrinkageSpec(*resp.EdgeNodes, newNodes)
			tflog.Info(ctx, "Shrinking edge cluster")
		}

		// Expand
		if len(oldNodes) < len(newNodes) {
			operation := expansion
			updateSpec.Operation = operation
			spec, err := nsx_edge_cluster.GetNsxEdgeClusterExpansionSpec(*resp.EdgeNodes, newNodes, client)

			if err != nil {
				return diag.FromErr(err)
			}
			updateSpec.EdgeClusterExpansionSpec = spec
			tflog.Info(ctx, "Expanding edge cluster")
		}

		taskRes, err := client.UpdateEdgeClusterWithResponse(ctx, data.Id(), updateSpec)
		if err != nil {
			return diag.FromErr(err)
		}
		task, vcfErr := api_client.GetResponseAs[vcf.Task](taskRes)
		if vcfErr != nil {
			api_client.LogError(vcfErr, ctx)
			return diag.FromErr(errors.New(*vcfErr.Message))
		}

		if err = api_client.NewTaskTracker(ctx, client, *task.Id).WaitForTask(); err != nil {
			return diag.FromErr(err)
		}
	}

	return nil
}

func validateClusterCreationSpec(client *vcf.ClientWithResponses, ctx context.Context, spec vcf.EdgeClusterCreationSpec) diag.Diagnostics {
	validateResponse, err := client.ValidateEdgeClusterCreationSpecWithResponse(ctx, spec)

	if err != nil {
		return validationUtils.ConvertVcfErrorToDiag(err)
	}
	validationResult, vcfErr := api_client.GetResponseAs[vcf.Validation](validateResponse)
	if vcfErr != nil {
		api_client.LogError(vcfErr, ctx)
		return diag.FromErr(errors.New(*vcfErr.Message))
	}

	if validationUtils.HasValidationFailed(validationResult) {
		return validationUtils.ConvertValidationResultToDiag(validationResult)
	}

	for {
		getValidationResponse, err := client.GetEdgeClusterValidationByIDWithResponse(ctx, *validationResult.Id)
		if err != nil {
			return validationUtils.ConvertVcfErrorToDiag(err)
		}
		validationStatus, vcfErr := api_client.GetResponseAs[vcf.Validation](getValidationResponse)
		if vcfErr != nil {
			api_client.LogError(vcfErr, ctx)
			return diag.FromErr(errors.New(*vcfErr.Message))
		}

		if validationUtils.HaveValidationChecksFinished(*validationStatus.ValidationChecks) {
			break
		}
		// TODO: reimplement this block without timeouts
		time.Sleep(10 * time.Second)
	}

	if validationUtils.HasValidationFailed(validationResult) {
		return validationUtils.ConvertValidationResultToDiag(validationResult)
	}

	return nil
}
