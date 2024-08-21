// © Broadcom. All Rights Reserved.
// The term “Broadcom” refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"fmt"
	"reflect"
	"time"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/vmware/vcf-sdk-go/client/domains"
	"github.com/vmware/vcf-sdk-go/models"

	"github.com/vmware/terraform-provider-vcf/internal/api_client"
	"github.com/vmware/terraform-provider-vcf/internal/cluster"
	"github.com/vmware/terraform-provider-vcf/internal/constants"
	"github.com/vmware/terraform-provider-vcf/internal/domain"
	"github.com/vmware/terraform-provider-vcf/internal/network"
	"github.com/vmware/terraform-provider-vcf/internal/resource_utils"
	validationUtils "github.com/vmware/terraform-provider-vcf/internal/validation"
	"github.com/vmware/terraform-provider-vcf/internal/vcenter"
)

func ResourceDomain() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceDomainCreate,
		ReadContext:   resourceDomainRead,
		UpdateContext: resourceDomainUpdate,
		DeleteContext: resourceDomainDelete,
		Importer: &schema.ResourceImporter{
			StateContext: func(ctx context.Context, data *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
				vcfClient := meta.(*api_client.SddcManagerClient)
				apiClient := vcfClient.ApiClient
				domainId := data.Id()
				// NOTE: Management domain cannot be imported, to not allow users to accidentally delete it,
				// but it can be used as datasource
				return domain.ImportDomain(ctx, data, apiClient, domainId, false)
			},
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(4 * time.Hour),
			Read:   schema.DefaultTimeout(20 * time.Minute),
			Update: schema.DefaultTimeout(4 * time.Hour),
			Delete: schema.DefaultTimeout(1 * time.Hour),
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringLenBetween(3, 20),
				Description:  "Name of the domain (from 3 to 20 characters)",
			},
			"org_name": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringLenBetween(3, 20),
				Description:  "Organization name of the workload domain",
			},
			"vcenter_configuration": {
				Type:        schema.TypeList,
				Required:    true,
				Description: "Specification describing vCenter Server instance settings",
				MinItems:    1,
				MaxItems:    1,
				Elem:        vcenter.VCSubresourceSchema(),
			},
			"nsx_configuration": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Specification details for NSX configuration",
				MaxItems:    1,
				Elem:        network.NsxSchema(),
			},
			"cluster": {
				Type:        schema.TypeList,
				Required:    true,
				Description: "Specification representing the clusters to be added to the workload domain",
				MinItems:    1,
				Elem:        clusterSubresourceSchema(),
			},
			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Status of the workload domain",
			},
			"type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Type of the workload domain",
			},
			"sso_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "ID of the SSO domain associated with the workload domain",
			},
			"sso_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Name of the SSO domain associated with the workload domain",
			},
			"is_management_sso_domain": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Shows whether the workload domain is joined to the management domain SSO",
			},
		},
	}
}

func resourceDomainCreate(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	vcfClient := meta.(*api_client.SddcManagerClient)
	apiClient := vcfClient.ApiClient

	domainCreationSpec, err := domain.CreateDomainCreationSpec(data)
	if err != nil {
		return diag.FromErr(err)
	}
	validateDomainSpec := domains.NewValidateDomainCreationSpecParamsWithContext(ctx).
		WithTimeout(constants.DefaultVcfApiCallTimeout)
	validateDomainSpec.DomainCreationSpec = domainCreationSpec

	validateResponse, err := apiClient.Domains.ValidateDomainCreationSpec(validateDomainSpec)
	if err != nil {
		return validationUtils.ConvertVcfErrorToDiag(err)
	}
	if validationUtils.HasValidationFailed(validateResponse.Payload) {
		return validationUtils.ConvertValidationResultToDiag(validateResponse.Payload)
	}

	domainCreationParams := domains.NewCreateDomainParamsWithContext(ctx).
		WithTimeout(constants.DefaultVcfApiCallTimeout)
	domainCreationParams.DomainCreationSpec = domainCreationSpec

	_, accepted, err := apiClient.Domains.CreateDomain(domainCreationParams)
	if err != nil {
		return validationUtils.ConvertVcfErrorToDiag(err)
	}
	taskId := accepted.Payload.ID
	err = vcfClient.WaitForTaskComplete(ctx, taskId, true)
	if err != nil {
		return diag.FromErr(err)
	}
	domainId, err := vcfClient.GetResourceIdAssociatedWithTask(ctx, taskId, "Domain")
	if err != nil {
		return diag.FromErr(err)
	}

	data.SetId(domainId)

	return resourceDomainRead(ctx, data, meta)
}

func resourceDomainRead(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	vcfClient := meta.(*api_client.SddcManagerClient)
	apiClient := vcfClient.ApiClient

	domainObj, err := domain.SetBasicDomainAttributes(ctx, data.Id(), data, apiClient)
	if err != nil {
		return diag.FromErr(err)
	}

	err = domain.ReadAndSetClustersDataToDomainResource(domainObj.Clusters, ctx, data, apiClient)
	if err != nil {
		return diag.FromErr(err)
	}
	nsxtClusterConfigRaw := data.Get("nsx_configuration").([]interface{})
	nsxtClusterConfig := nsxtClusterConfigRaw[0].(map[string]interface{})
	nsxtClusterConfig["id"] = domainObj.NSXTCluster.ID
	_ = data.Set("nsx_configuration", nsxtClusterConfigRaw)

	return nil
}

func resourceDomainUpdate(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	vcfClient := meta.(*api_client.SddcManagerClient)
	apiClient := vcfClient.ApiClient

	// Domain Update API supports only changes to domain name and Cluster Import
	if data.HasChange("name") {
		domainUpdateSpec := domain.CreateDomainUpdateSpec(data, false)
		domainUpdateParams := domains.NewUpdateDomainParamsWithContext(ctx).
			WithTimeout(constants.DefaultVcfApiCallTimeout)
		domainUpdateParams.DomainUpdateSpec = domainUpdateSpec
		domainUpdateParams.ID = data.Id()

		_, accepted, err := apiClient.Domains.UpdateDomain(domainUpdateParams)
		if err != nil {
			return diag.FromErr(err)
		}
		taskId := accepted.Payload.ID
		err = vcfClient.WaitForTaskComplete(ctx, taskId, false)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	if data.HasChange("cluster") {
		oldClustersValue, newClustersValue := data.GetChange("cluster")
		newClustersList := newClustersValue.([]interface{})
		oldClustersList := oldClustersValue.([]interface{})
		if len(oldClustersList) == len(newClustersList) {
			diags := handleClusterUpdateInDomain(ctx, newClustersList, oldClustersList, vcfClient)
			if diags != nil {
				return diags
			}
		} else {
			diags := handleClusterAddRemoveToDomain(ctx, data.Id(), newClustersList, oldClustersList, vcfClient)
			if diags != nil {
				return diags
			}
		}
	}

	return resourceDomainRead(ctx, data, meta)
}

func handleClusterAddRemoveToDomain(ctx context.Context, domainId string, newClustersList, oldClustersList []interface{},
	vcfClient *api_client.SddcManagerClient) diag.Diagnostics {
	addedClustersList, removedClustersList := resource_utils.CalculateAddedRemovedResources(newClustersList, oldClustersList)
	for _, addedCluster := range addedClustersList {
		clusterSpec, err := cluster.TryConvertToClusterSpec(addedCluster)
		if err != nil {
			return diag.FromErr(err)
		}
		// subsequent domain read will set the cluster ID, so we can discard it here
		_, diags := createCluster(ctx, domainId, clusterSpec, vcfClient)
		if diags != nil {
			return diags
		}
	}

	for _, removedCluster := range removedClustersList {
		clusterId := removedCluster["id"].(string)
		diags := deleteCluster(ctx, clusterId, vcfClient)
		if diags != nil {
			return diags
		}
	}

	return nil
}

func handleClusterUpdateInDomain(ctx context.Context, newClustersStateList, oldClustersStateList []interface{},
	vcfClient *api_client.SddcManagerClient) diag.Diagnostics {
	if len(oldClustersStateList) != len(newClustersStateList) {
		return diag.FromErr(fmt.Errorf("expecting old and new cluster list to have the same length"))
	}
	for i, newClusterState := range newClustersStateList {
		// skip the clusters that have no changes
		if reflect.DeepEqual(newClusterState, oldClustersStateList[i]) {
			continue
		}
		oldClusterStateMap := oldClustersStateList[i].(map[string]interface{})
		newClusterStateMap := newClusterState.(map[string]interface{})
		// sanity check that we're comparing the same clusters for changes to their hosts
		newClusterStateId := newClusterStateMap["id"].(string)
		oldClusterStateId := oldClusterStateMap["id"].(string)
		if newClusterStateId != oldClusterStateId {
			return diag.FromErr(fmt.Errorf("cluster order has changed, updating hosts in cluster not supported"))
		}
		oldHostsList := oldClusterStateMap["host"].([]interface{})
		newHostsList := newClusterStateMap["host"].([]interface{})
		if reflect.DeepEqual(oldHostsList, newHostsList) {
			tflog.Warn(ctx, "only expand/contract cluster update is supported")
			continue
		}

		clusterUpdateSpec := new(models.ClusterUpdateSpec)
		populatedClusterUpdateSpec, err := cluster.SetExpansionOrContractionSpec(clusterUpdateSpec, oldHostsList, newHostsList)
		if err != nil {
			return diag.FromErr(err)
		}

		diags := updateCluster(ctx, newClusterStateId, populatedClusterUpdateSpec, vcfClient)
		if diags != nil {
			return diags
		}
	}
	return nil
}

func resourceDomainDelete(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	vcfClient := meta.(*api_client.SddcManagerClient)
	apiClient := vcfClient.ApiClient

	markForDeleteUpdateSpec := domain.CreateDomainUpdateSpec(data, true)
	domainUpdateParams := domains.NewUpdateDomainParamsWithContext(ctx).
		WithTimeout(constants.DefaultVcfApiCallTimeout)
	domainUpdateParams.DomainUpdateSpec = markForDeleteUpdateSpec
	domainUpdateParams.ID = data.Id()

	acceptedUpdateTask, _, err := apiClient.Domains.UpdateDomain(domainUpdateParams)
	if err != nil {
		return diag.FromErr(err)
	}
	taskId := acceptedUpdateTask.Payload.ID
	err = vcfClient.WaitForTaskComplete(ctx, taskId, false)
	if err != nil {
		return diag.FromErr(err)
	}

	domainDeleteParams := domains.NewDeleteDomainParamsWithContext(ctx).
		WithTimeout(constants.DefaultVcfApiCallTimeout)
	domainDeleteParams.ID = data.Id()

	acceptedDeleteTask, acceptedDeleteTask2, err := apiClient.Domains.DeleteDomain(domainDeleteParams)
	if err != nil {
		return diag.FromErr(err)
	}
	if acceptedDeleteTask != nil {
		taskId = acceptedDeleteTask.Payload.ID
	}
	if acceptedDeleteTask2 != nil {
		taskId = acceptedDeleteTask2.Payload.ID
	}
	err = vcfClient.WaitForTaskComplete(ctx, taskId, true)
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}
