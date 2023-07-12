/* Copyright 2023 VMware, Inc.
   SPDX-License-Identifier: MPL-2.0 */

package provider

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/vmware/terraform-provider-vcf/internal/constants"
	"github.com/vmware/terraform-provider-vcf/internal/network"
	validation_utils "github.com/vmware/terraform-provider-vcf/internal/validation"
	"github.com/vmware/terraform-provider-vcf/internal/vcenter"
	"github.com/vmware/vcf-sdk-go/client"
	"github.com/vmware/vcf-sdk-go/client/clusters"
	"github.com/vmware/vcf-sdk-go/client/domains"
	"github.com/vmware/vcf-sdk-go/models"
	"strconv"
	"time"
)

func ResourceDomain() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceDomainCreate,
		ReadContext:   resourceDomainRead,
		UpdateContext: resourceDomainUpdate,
		DeleteContext: resourceDomainDelete,
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(12 * time.Hour),
			Read:   schema.DefaultTimeout(20 * time.Minute),
			Update: schema.DefaultTimeout(6 * time.Hour),
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
			"vcenter": {
				Type:        schema.TypeList,
				Required:    true,
				Description: "Specification describing vcenter settings",
				MinItems:    1,
				MaxItems:    1,
				Elem:        vcenter.VCSubresourceSchema(),
			},
			"nsx_configuration": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Specification details for NSX-T configuration",
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
				Description: "Shows whether the workload domain is joined to the Management domain SSO",
			},
			"total_cpu_capacity": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Represents a cpu total metric for the domain",
			},
			"used_cpu_capacity": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Represents a cpu used metric for the domain",
			},
			"total_memory_capacity": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Represents a memory total metric for the domain",
			},
			"used_memory_capacity": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Represents a memory used metric for the domain",
			},
			"total_storage_capacity": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Represents a storage total metric for the domain",
			},
			"used_storage_capacity": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Represents a storage used metric for the domain",
			},
		},
	}
}

func resourceDomainCreate(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	vcfClient := meta.(*SddcManagerClient)
	apiClient := vcfClient.ApiClient

	domainCreationSpec, err := createDomainCreationSpec(data)
	if err != nil {
		return diag.FromErr(err)
	}
	validateDomainSpec := domains.NewValidateDomainsOperationsParamsWithContext(ctx).
		WithTimeout(constants.DefaultVcfApiCallTimeout)
	validateDomainSpec.DomainCreationSpec = domainCreationSpec

	validateResponse, err := apiClient.Domains.ValidateDomainsOperations(validateDomainSpec)
	if err != nil {
		return validation_utils.ConvertVcfErrorToDiag(err)
	}
	if validation_utils.HasValidationFailed(validateResponse.Payload) {
		return validation_utils.ConvertValidationResultToDiag(validateResponse.Payload)
	}

	domainCreationParams := domains.NewCreateDomainParamsWithContext(ctx).
		WithTimeout(constants.DefaultVcfApiCallTimeout)
	domainCreationParams.DomainCreationSpec = domainCreationSpec

	_, accepted, err := apiClient.Domains.CreateDomain(domainCreationParams)
	if err != nil {
		return validation_utils.ConvertVcfErrorToDiag(err)
	}
	taskId := accepted.Payload.ID
	err = vcfClient.WaitForTaskComplete(ctx, taskId)
	if err != nil {
		return diag.FromErr(err)
	}
	domainId, err := vcfClient.GetResourceIdAssociatedWithTask(ctx, taskId)
	if err != nil {
		return diag.FromErr(err)
	}

	data.SetId(domainId)

	return resourceDomainRead(ctx, data, meta)
}

func resourceDomainRead(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	vcfClient := meta.(*SddcManagerClient)
	apiClient := vcfClient.ApiClient

	getDomainParams := domains.NewGetDomainParamsWithContext(ctx).
		WithTimeout(constants.DefaultVcfApiCallTimeout)
	getDomainParams.ID = data.Id()
	domainResult, err := apiClient.Domains.GetDomain(getDomainParams)
	if err != nil {
		return diag.FromErr(err)
	}
	domain := domainResult.Payload

	_ = data.Set("name", domain.Name)
	_ = data.Set("status", domain.Status)
	_ = data.Set("type", domain.Type)
	_ = data.Set("sso_id", domain.SSOID)
	_ = data.Set("sso_name", domain.SSOName)
	_ = data.Set("is_management_sso_domain", domain.IsManagementSSODomain)
	if len(domain.VCENTERS) < 1 {
		return diag.FromErr(fmt.Errorf("no vCenters found for domain %q", data.Id()))
	}
	_ = data.Set("vcenter.0.id", domain.VCENTERS[0].ID)
	_ = data.Set("vcenter.0.fqdn", domain.VCENTERS[0].Fqdn)

	err = readAndSetClustersDataToDomainResource(data, apiClient)
	if err != nil {
		return diag.FromErr(err)
	}

	totalCpuCapacity := strconv.FormatFloat(domain.Capacity.CPU.Total.Value,
		'f', 2, 64) + " " + domain.Capacity.CPU.Total.Unit
	usedCpuCapacity := strconv.FormatFloat(domain.Capacity.CPU.Used.Value,
		'f', 2, 64) + " " + domain.Capacity.CPU.Used.Unit
	_ = data.Set("total_cpu_capacity", totalCpuCapacity)
	_ = data.Set("used_cpu_capacity", usedCpuCapacity)

	totalMemoryCapacity := strconv.FormatFloat(domain.Capacity.Memory.Total.Value,
		'f', 2, 64) + " " + domain.Capacity.Memory.Total.Unit
	usedMemoryCapacity := strconv.FormatFloat(domain.Capacity.Memory.Used.Value,
		'f', 2, 64) + " " + domain.Capacity.Memory.Used.Unit
	_ = data.Set("total_memory_capacity", totalMemoryCapacity)
	_ = data.Set("used_memory_capacity", usedMemoryCapacity)

	totalStorageCapacity := strconv.FormatFloat(domain.Capacity.Storage.Total.Value,
		'f', 2, 64) + " " + domain.Capacity.Storage.Total.Unit
	usedStorageCapacity := strconv.FormatFloat(domain.Capacity.Storage.Used.Value,
		'f', 2, 64) + " " + domain.Capacity.Storage.Used.Unit
	_ = data.Set("total_storage_capacity", totalStorageCapacity)
	_ = data.Set("used_storage_capacity", usedStorageCapacity)

	return nil
}

func resourceDomainUpdate(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	vcfClient := meta.(*SddcManagerClient)
	apiClient := vcfClient.ApiClient

	// Domain Update API supports only changes to domain name and Cluster Import
	// TODO implement cluster import scenario
	if data.HasChange("name") {
		domainUpdateSpec := createDomainUpdateSpec(data, false)
		domainUpdateParams := domains.NewUpdateDomainParamsWithContext(ctx).
			WithTimeout(constants.DefaultVcfApiCallTimeout)
		domainUpdateParams.DomainUpdateSpec = domainUpdateSpec
		domainUpdateParams.ID = data.Id()

		_, accepted, err := apiClient.Domains.UpdateDomain(domainUpdateParams)
		if err != nil {
			return diag.FromErr(err)
		}
		taskId := accepted.Payload.ID
		err = vcfClient.WaitForTaskComplete(ctx, taskId)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	return resourceDomainRead(ctx, data, meta)
}

func resourceDomainDelete(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	vcfClient := meta.(*SddcManagerClient)
	apiClient := vcfClient.ApiClient

	markForDeleteUpdateSpec := createDomainUpdateSpec(data, true)
	domainUpdateParams := domains.NewUpdateDomainParamsWithContext(ctx).
		WithTimeout(constants.DefaultVcfApiCallTimeout)
	domainUpdateParams.DomainUpdateSpec = markForDeleteUpdateSpec
	domainUpdateParams.ID = data.Id()

	_, acceptedUpdateTask, err := apiClient.Domains.UpdateDomain(domainUpdateParams)
	if err != nil {
		return diag.FromErr(err)
	}
	taskId := acceptedUpdateTask.Payload.ID
	err = vcfClient.WaitForTaskComplete(ctx, taskId)
	if err != nil {
		return diag.FromErr(err)
	}

	domainDeleteParams := domains.NewDeleteDomainParamsWithContext(ctx).
		WithTimeout(constants.DefaultVcfApiCallTimeout)
	domainDeleteParams.ID = data.Id()

	_, acceptedDeleteTask, err := apiClient.Domains.DeleteDomain(domainDeleteParams)
	if err != nil {
		return diag.FromErr(err)
	}
	taskId = acceptedDeleteTask.Payload.ID
	err = vcfClient.WaitForTaskComplete(ctx, taskId)
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func createDomainCreationSpec(data *schema.ResourceData) (*models.DomainCreationSpec, error) {
	result := new(models.DomainCreationSpec)
	domainName := data.Get("name").(string)
	result.DomainName = &domainName

	if orgName, ok := data.GetOk("org_name"); ok {
		result.OrgName = orgName.(string)
	}

	vcenterSpec, err := generateVcenterSpecFromResourceData(data)
	if err == nil {
		result.VcenterSpec = vcenterSpec
	} else {
		return nil, err
	}

	nsxSpec, err := generateNsxSpecFromResourceData(data)
	if err == nil {
		result.NsxTSpec = nsxSpec
	} else {
		return nil, err
	}

	computeSpec, err := generateComputeSpecFromResourceData(data)
	if err == nil {
		result.ComputeSpec = computeSpec
	} else {
		return nil, err
	}

	return result, nil
}

func readAndSetClustersDataToDomainResource(data *schema.ResourceData, apiClient *client.VcfClient) error {
	getClustersParams := clusters.GetClustersParams{
		// TODO: handle stretched clusters
		IsStretched: toBoolPointer(false),
	}
	getClustersParams.WithTimeout(constants.DefaultVcfApiCallTimeout)

	// TODO: consider parallel GetCluster(clusterId) calls
	clustersResult, err := apiClient.Clusters.GetClusters(&getClustersParams)
	if err != nil {
		return err
	}
	domainClusterData := data.Get("cluster")
	domainClusterDataList := domainClusterData.([]interface{})
	allClusters := clustersResult.Payload.Elements
	for _, domainClusterRaw := range domainClusterDataList {
		domainCluster := domainClusterRaw.(map[string]interface{})
		for _, cluster := range allClusters {
			if domainCluster["name"] == cluster.Name {
				domainCluster["id"] = cluster.ID
				domainCluster["primary_datastore_name"] = cluster.PrimaryDatastoreName
				domainCluster["primary_datastore_type"] = cluster.PrimaryDatastoreType
				domainCluster["is_default"] = cluster.IsDefault
				domainCluster["is_stretched"] = cluster.IsStretched
			}
		}
	}
	_ = data.Set("cluster", domainClusterData)

	return nil
}

func createDomainUpdateSpec(data *schema.ResourceData, markForDeletion bool) *models.DomainUpdateSpec {
	result := new(models.DomainUpdateSpec)
	if markForDeletion {
		result.MarkForDeletion = true
		return result
	}
	if data.HasChange("name") {
		result.Name = data.Get("name").(string)
	}

	// TODO implement support for IPPoolSpecs in NsxTSpec, then implement this "Cluster_Import" scenario
	// by placing the added cluster spec in the DomainUpdateSpec
	//nsxtSpec, err := generateNsxSpecFromResourceData(data)
	//if err == nil {
	//	result.NsxTSpec = nsxtSpec
	//} else {
	//	return nil, err
	//}

	return result
}

func generateNsxSpecFromResourceData(data *schema.ResourceData) (*models.NsxTSpec, error) {
	if nsxConfigRaw, ok := data.GetOk("nsx_configuration"); ok && len(nsxConfigRaw.([]interface{})) > 0 {
		nsxConfigList := nsxConfigRaw.([]interface{})
		nsxConfigListEntry := nsxConfigList[0].(map[string]interface{})
		nsxSpec, err := network.TryConvertToNsxSpec(nsxConfigListEntry)
		return nsxSpec, err
	}
	return nil, nil
}

func generateVcenterSpecFromResourceData(data *schema.ResourceData) (*models.VcenterSpec, error) {
	if vcenterConfigRaw, ok := data.GetOk("vcenter"); ok && len(vcenterConfigRaw.([]interface{})) > 0 {
		vcenterConfigList := vcenterConfigRaw.([]interface{})
		vcenterConfigListEntry := vcenterConfigList[0].(map[string]interface{})
		vcenterSpec, err := vcenter.TryConvertToVcenterSpec(vcenterConfigListEntry)
		return vcenterSpec, err
	}
	return nil, nil
}

func generateComputeSpecFromResourceData(data *schema.ResourceData) (*models.ComputeSpec, error) {
	if clusterConfigRaw, ok := data.GetOk("cluster"); ok && !validation_utils.IsEmpty(clusterConfigRaw) {
		clusterConfigList := clusterConfigRaw.([]interface{})
		result := new(models.ComputeSpec)
		var clusterSpecs []*models.ClusterSpec
		for _, clusterConfigListEntry := range clusterConfigList {
			clusterSpec, err := tryConvertToClusterSpec(clusterConfigListEntry.(map[string]interface{}))
			if err != nil {
				return nil, err
			}
			clusterSpecs = append(clusterSpecs, clusterSpec)
		}
		result.ClusterSpecs = clusterSpecs
		return result, nil
	}
	return nil, fmt.Errorf("no cluster configuration")
}
