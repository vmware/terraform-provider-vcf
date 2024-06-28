// Copyright 2024 Broadcom. All Rights Reserved.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/vmware/terraform-provider-vcf/internal/api_client"
	"github.com/vmware/terraform-provider-vcf/internal/constants"
	"github.com/vmware/terraform-provider-vcf/internal/nsx_edge_cluster"
	validationUtils "github.com/vmware/terraform-provider-vcf/internal/validation"
	vcfClient "github.com/vmware/vcf-sdk-go/client"
	"github.com/vmware/vcf-sdk-go/client/nsxt_edge_clusters"
	"github.com/vmware/vcf-sdk-go/models"
	"math"
	"time"
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
				RequiredWith: []string{"tier0_name"},
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
				RequiredWith: []string{"tier0_name"},
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
				Type:         schema.TypeInt,
				Optional:     true,
				Description:  "ASN for the cluster",
				ValidateFunc: validation.IntBetween(1, int(math.Pow(2, 31)-1)),
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

	validationErr := validateClusterCreationSpec(client, ctx, spec)

	if validationErr != nil {
		return validationErr
	}

	createClusterParams := &nsxt_edge_clusters.CreateEdgeClusterParams{
		EdgeCreationSpec: spec,
		Context:          ctx,
	}

	_, task, err := client.NSXTEdgeClusters.CreateEdgeCluster(createClusterParams.WithTimeout(constants.DefaultVcfApiCallTimeout))

	if err != nil {
		return diag.FromErr(err)
	}

	tflog.Info(ctx, "Edge cluster creation has started.")
	err = meta.(*api_client.SddcManagerClient).WaitForTaskComplete(ctx, task.Payload.ID, false)
	if err != nil {
		return diag.FromErr(err)
	}

	getClusterParams := &nsxt_edge_clusters.GetEdgeClustersParams{}
	clusters, err := client.NSXTEdgeClusters.GetEdgeClusters(getClusterParams.WithTimeout(constants.DefaultVcfApiCallTimeout))

	if err != nil {
		return diag.FromErr(err)
	}

	for _, cluster := range clusters.Payload.Elements {
		if cluster.Name == data.Get("name") {
			data.SetId(cluster.ID)
			tflog.Info(ctx, "Edge cluster created successfully.")
			return nil
		}
	}

	return diag.Errorf("Edge cluster creation failed - cluster not found in inventory")
}

func resourceNsxEdgeClusterRead(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*api_client.SddcManagerClient).ApiClient

	_, err := getEdgeCluster(ctx, client, data.Id())

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

	edgeClusterOk, err := getEdgeCluster(ctx, client, data.Id())

	if err != nil {
		return diag.FromErr(err)
	}

	if data.HasChange("edge_node") {
		oldNodesRaw, newNodesRaw := data.GetChange("edge_node")
		oldNodes, newNodes := oldNodesRaw.([]interface{}), newNodesRaw.([]interface{})

		if len(oldNodes) == len(newNodes) {
			return diag.Errorf("Adding and removing edge nodes is not supported in a single configuration change. Apply each change separately.")
		}

		updateParams := nsxt_edge_clusters.NewUpdateEdgeClusterParamsWithContext(ctx)
		updateParams.ID = data.Id()
		updateParams.EdgeClusterUpdateSpec = &models.EdgeClusterUpdateSpec{}

		// Shrink
		if len(oldNodes) > len(newNodes) {
			operation := shrinkage
			updateParams.EdgeClusterUpdateSpec.Operation = &operation
			updateParams.EdgeClusterUpdateSpec.EdgeClusterShrinkageSpec =
				nsx_edge_cluster.GetNsxEdgeClusterShrinkageSpec(edgeClusterOk.Payload.EdgeNodes, newNodes)
			tflog.Info(ctx, "Shrinking edge cluster")
		}

		// Expand
		if len(oldNodes) < len(newNodes) {
			operation := expansion
			updateParams.EdgeClusterUpdateSpec.Operation = &operation
			spec, err := nsx_edge_cluster.GetNsxEdgeClusterExpansionSpec(edgeClusterOk.Payload.EdgeNodes, newNodes, client)

			if err != nil {
				return diag.FromErr(err)
			}
			updateParams.EdgeClusterUpdateSpec.EdgeClusterExpansionSpec = spec
			tflog.Info(ctx, "Expanding edge cluster")
		}

		_, task, err := client.NSXTEdgeClusters.UpdateEdgeCluster(updateParams.WithTimeout(constants.DefaultVcfApiCallTimeout))
		if err != nil {
			return diag.FromErr(err)
		}

		err = meta.(*api_client.SddcManagerClient).WaitForTaskComplete(ctx, task.Payload.ID, false)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	return nil
}

func getEdgeCluster(ctx context.Context, client *vcfClient.VcfClient, id string) (*nsxt_edge_clusters.GetEdgeClusterOK, error) {
	params := nsxt_edge_clusters.NewGetEdgeClusterParamsWithContext(ctx)
	params.ID = id

	return client.NSXTEdgeClusters.GetEdgeCluster(params.WithTimeout(constants.DefaultVcfApiCallTimeout))
}

func validateClusterCreationSpec(client *vcfClient.VcfClient, ctx context.Context, spec *models.EdgeClusterCreationSpec) diag.Diagnostics {
	validateClusterParams := &nsxt_edge_clusters.ValidateEdgeClusterCreationSpecParams{
		EdgeCreationSpec: spec,
	}

	_, validateResponse, err := client.NSXTEdgeClusters.ValidateEdgeClusterCreationSpec(validateClusterParams.WithTimeout(constants.DefaultVcfApiCallTimeout))

	if err != nil {
		return validationUtils.ConvertVcfErrorToDiag(err)
	}

	validationResult := validateResponse.Payload
	if validationUtils.HasValidationFailed(validationResult) {
		return validationUtils.ConvertValidationResultToDiag(validationResult)
	}

	for {
		getClusterValidationParams := nsxt_edge_clusters.NewGetEdgeClusterValidationByIDParamsWithContext(ctx).
			WithTimeout(constants.DefaultVcfApiCallTimeout)
		getClusterValidationParams.SetID(validateResponse.Payload.ID)
		getValidationResponse, err := client.NSXTEdgeClusters.GetEdgeClusterValidationByID(getClusterValidationParams)
		if err != nil {
			return validationUtils.ConvertVcfErrorToDiag(err)
		}
		validationResult = getValidationResponse.Payload
		if validationUtils.HaveValidationChecksFinished(validationResult.ValidationChecks) {
			break
		}
		time.Sleep(10 * time.Second)
	}

	if validationUtils.HasValidationFailed(validationResult) {
		return validationUtils.ConvertValidationResultToDiag(validationResult)
	}

	return nil
}
