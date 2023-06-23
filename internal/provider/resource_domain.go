/* Copyright 2023 VMware, Inc.
   SPDX-License-Identifier: MPL-2.0 */

package provider

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/vmware/terraform-provider-vcf/internal/network"
	validation_utils "github.com/vmware/terraform-provider-vcf/internal/validation"
	"github.com/vmware/vcf-sdk-go/client"
	"github.com/vmware/vcf-sdk-go/client/clusters"
	"github.com/vmware/vcf-sdk-go/client/domains"
	"github.com/vmware/vcf-sdk-go/models"
	"strconv"
	"strings"
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
			"vcenter_name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.NoZeroValues,
				Description:  "Name of the vCenter virtual machine to be created with the domain",
			},
			"vcenter_datacenter_name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.NoZeroValues,
				Description:  "vCenter datacenter name",
			},
			"vcenter_root_password": {
				Type:         schema.TypeString,
				Required:     true,
				Sensitive:    true,
				ForceNew:     true,
				Description:  "Password for the vCenter root shell user (8-20 characters)",
				ValidateFunc: validation_utils.ValidatePassword,
			},
			"vcenter_vm_size": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "VCenter VM size. One among: xlarge, large, medium, small, tiny",
				ValidateFunc: validation.StringInSlice([]string{
					"xlarge", "large", "medium", "small", "tiny",
				}, true),
				DiffSuppressFunc: func(k, oldValue, newValue string, d *schema.ResourceData) bool {
					return oldValue == strings.ToUpper(newValue) || strings.ToUpper(oldValue) == newValue
				},
			},
			"vcenter_storage_size": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "VCenter storage size. One among: lstorage, xlstorage",
				ValidateFunc: validation.StringInSlice([]string{
					"lstorage", "xlstorage",
				}, true),
				DiffSuppressFunc: func(k, oldValue, newValue string, d *schema.ResourceData) bool {
					return oldValue == strings.ToUpper(newValue) || strings.ToUpper(oldValue) == newValue
				},
			},
			"vcenter_ip_address": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				Description:  "IPv4 address of the vCenter virtual machine",
				ValidateFunc: validation_utils.ValidateIPv4AddressSchema,
			},
			"vcenter_subnet_mask": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				Description:  "Subnet mask",
				ValidateFunc: validation_utils.ValidateIPv4AddressSchema,
			},
			"vcenter_gateway": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				Description:  "IPv4 gateway the vCenter VM can use to connect to the outside world",
				ValidateFunc: validation_utils.ValidateIPv4AddressSchema,
			},
			"vcenter_dns_name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				Description:  "DNS name of the virtual machine, e.g., vc-1.domain1.rainpole.io",
				ValidateFunc: validation.NoZeroValues,
			},
			"nsxt_configuration": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Specification details for NSX-T configuration",
				MaxItems:    1,
				Elem:        network.NsxTSchema(),
			},
			"cluster": {
				Type:        schema.TypeList,
				Required:    true,
				Description: "Specification representing the clusters to be added to the workload domain",
				MinItems:    1,
				Elem:        clusterSubresourceSchema(),
			},
			"vcenter_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "ID of the vCenter",
			},
			"vcenter_fqdn": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "FQDN of the vCenter",
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
	validateDomainSpec := domains.NewValidateDomainsOperationsParamsWithContext(ctx)
	validateDomainSpec.DomainCreationSpec = domainCreationSpec

	validateResponse, err := apiClient.Domains.ValidateDomainsOperations(validateDomainSpec)
	if err != nil {
		return validation_utils.ConvertVcfErrorToDiag(err)
	}
	if validation_utils.HasValidationFailed(validateResponse.Payload) {
		return validation_utils.ConvertValidationResultToDiag(validateResponse.Payload)
	}

	domainCreationParams := domains.NewCreateDomainParamsWithContext(ctx)
	domainCreationParams.DomainCreationSpec = domainCreationSpec

	_, accepted, err := apiClient.Domains.CreateDomain(domainCreationParams)
	if err != nil {
		return validation_utils.ConvertVcfErrorToDiag(err)
	}
	taskId := accepted.Payload.ID
	err = vcfClient.WaitForTaskComplete(taskId)
	if err != nil {
		return diag.FromErr(err)
	}
	domainId, err := vcfClient.GetResourceIdAssociatedWithTask(taskId)
	if err != nil {
		return diag.FromErr(err)
	}

	data.SetId(domainId)

	return resourceDomainRead(ctx, data, meta)
}

func resourceDomainRead(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	vcfClient := meta.(*SddcManagerClient)
	apiClient := vcfClient.ApiClient

	getDomainParams := domains.NewGetDomainParamsWithContext(ctx)
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
	_ = data.Set("vcenter_id", domain.VCENTERS[0].ID)
	_ = data.Set("vcenter_fqdn", domain.VCENTERS[0].Fqdn)

	err = readAndSetClustersData(domain.Clusters, data, apiClient)
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

	domainUpdateSpec, err := createDomainUpdateSpec(data, false)
	if err != nil {
		tflog.Error(ctx, err.Error())
		return diag.FromErr(err)
	}
	domainUpdateParams := domains.NewUpdateDomainParamsWithContext(ctx)
	domainUpdateParams.DomainUpdateSpec = domainUpdateSpec
	domainUpdateParams.ID = data.Id()

	_, accepted, err := apiClient.Domains.UpdateDomain(domainUpdateParams)
	if err != nil {
		return diag.FromErr(err)
	}
	taskId := accepted.Payload.ID
	err = vcfClient.WaitForTaskComplete(taskId)
	if err != nil {
		return diag.FromErr(err)
	}

	return resourceDomainRead(ctx, data, meta)
}

func resourceDomainDelete(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	vcfClient := meta.(*SddcManagerClient)
	apiClient := vcfClient.ApiClient

	markForDeleteUpdateSpec, _ := createDomainUpdateSpec(data, true)
	domainUpdateParams := domains.NewUpdateDomainParamsWithContext(ctx)
	domainUpdateParams.DomainUpdateSpec = markForDeleteUpdateSpec
	domainUpdateParams.ID = data.Id()

	_, acceptedUpdateTask, err := apiClient.Domains.UpdateDomain(domainUpdateParams)
	if err != nil {
		return diag.FromErr(err)
	}
	taskId := acceptedUpdateTask.Payload.ID
	err = vcfClient.WaitForTaskComplete(taskId)
	if err != nil {
		return diag.FromErr(err)
	}

	domainDeleteParams := domains.NewDeleteDomainParamsWithContext(ctx)
	domainDeleteParams.ID = data.Id()

	_, acceptedDeleteTask, err := apiClient.Domains.DeleteDomain(domainDeleteParams)
	if err != nil {
		return diag.FromErr(err)
	}
	taskId = acceptedDeleteTask.Payload.ID
	err = vcfClient.WaitForTaskComplete(taskId)
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

	// Create Domain with empty, but non-null ComputeSpec, to pass backend validation.
	// Clusters can be attached to a domain with a separate Cluster resource, resulting
	// in a cleaner TF config files and easier management of Domain resource diffs.
	result.ComputeSpec = new(models.ComputeSpec)

	vcenterName := data.Get("vcenter_name").(string)
	vcenterDatacenterName := data.Get("vcenter_datacenter_name").(string)
	vcenterRootPassword := data.Get("vcenter_root_password").(string)
	vcenterStorageSize := data.Get("vcenter_storage_size").(string)
	vcenterVmSize := data.Get("vcenter_vm_size").(string)

	vcenterIp := data.Get("vcenter_ip_address").(string)
	vcenterSubnetMask := data.Get("vcenter_subnet_mask").(string)
	vcenterGateway := data.Get("vcenter_gateway").(string)
	vcenterDnsName := data.Get("vcenter_dns_name").(string)
	networkDetailsSpec := new(models.NetworkDetailsSpec)
	networkDetailsSpec.IPAddress = &vcenterIp
	networkDetailsSpec.SubnetMask = vcenterSubnetMask
	networkDetailsSpec.DNSName = vcenterDnsName
	networkDetailsSpec.Gateway = vcenterGateway

	vcenterSpec := models.VcenterSpec{
		DatacenterName:     &vcenterDatacenterName,
		Name:               &vcenterName,
		RootPassword:       &vcenterRootPassword,
		StorageSize:        vcenterStorageSize,
		VMSize:             vcenterVmSize,
		NetworkDetailsSpec: networkDetailsSpec,
	}
	result.VcenterSpec = &vcenterSpec

	nsxtSpec, err := generateNsxtSpecFromResourceData(data)
	if err == nil {
		result.NsxTSpec = nsxtSpec
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

func readAndSetClustersData(domainClusterRefs []*models.ClusterReference, data *schema.ResourceData, apiClient *client.VcfClient) error {
	// TODO get all cluster infos and match them via name, then set the clusters data
	clusterIds := make(map[string]interface{}, len(domainClusterRefs))
	for _, clusterReference := range domainClusterRefs {
		clusterIds[*clusterReference.ID] = true
	}

	getClustersParams := clusters.GetClustersParams{
		IsStretched: toBoolPointer(false),
	}

	// TODO: consider parallel GetCluster(clusterId) calls
	clustersResult, err := apiClient.Clusters.GetClusters(&getClustersParams)
	if err != nil {
		return err
	}
	domainClusterData := data.Get("cluster")
	domainClusterDataList := domainClusterData.([]interface{})
	allClusters := clustersResult.Payload.Elements
	for i, domainClusterRaw := range domainClusterDataList {
		domainCluster := domainClusterRaw.(map[string]interface{})
		for _, cluster := range allClusters {
			if domainCluster["name"] == cluster.Name {
				clusterIndex := fmt.Sprintf("cluster.%d", i)
				_ = data.Set(clusterIndex+".id", cluster.ID)
				_ = data.Set(clusterIndex+".primary_datastore_name", cluster.PrimaryDatastoreName)
				_ = data.Set(clusterIndex+".primary_datastore_type", cluster.PrimaryDatastoreType)
				_ = data.Set(clusterIndex+".is_default", cluster.IsDefault)
				_ = data.Set(clusterIndex+".is_streched", cluster.IsStretched)
			}
		}
	}

	//_ = data.Set("cluster_ids", cluster_ids)

	return nil
}

func createDomainUpdateSpec(data *schema.ResourceData, markForDeletion bool) (*models.DomainUpdateSpec, error) {
	result := new(models.DomainUpdateSpec)
	if markForDeletion {
		result.MarkForDeletion = true
		return result, nil
	}
	result.Name = data.Get("name").(string)
	nsxtSpec, err := generateNsxtSpecFromResourceData(data)
	if err == nil {
		result.NsxTSpec = nsxtSpec
	} else {
		return nil, err
	}

	return result, nil
}

func generateNsxtSpecFromResourceData(data *schema.ResourceData) (*models.NsxTSpec, error) {
	if nsxtConfigRaw, ok := data.GetOk("nsxt_configuration"); ok && len(nsxtConfigRaw.([]interface{})) > 0 {
		nsxtConfigList := nsxtConfigRaw.([]interface{})
		nsxtConfigListEntry := nsxtConfigList[0].(map[string]interface{})
		nsxtSpec, err := network.TryConvertToNsxtSpec(nsxtConfigListEntry)
		return nsxtSpec, err
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
