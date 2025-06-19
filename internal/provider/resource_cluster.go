// © Broadcom. All Rights Reserved.
// The term “Broadcom” refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"errors"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	utils "github.com/vmware/terraform-provider-vcf/internal/resource_utils"
	"github.com/vmware/vcf-sdk-go/vcf"

	"github.com/vmware/terraform-provider-vcf/internal/api_client"
	"github.com/vmware/terraform-provider-vcf/internal/cluster"
	"github.com/vmware/terraform-provider-vcf/internal/datastores"
	"github.com/vmware/terraform-provider-vcf/internal/network"
	validationUtils "github.com/vmware/terraform-provider-vcf/internal/validation"
	"github.com/vmware/terraform-provider-vcf/internal/vsan"
)

func ResourceCluster() *schema.Resource {
	clusterResourceSchema := clusterSubresourceSchema().Schema
	clusterResourceSchema["domain_id"] = &schema.Schema{
		Type:         schema.TypeString,
		Optional:     true,
		Computed:     true,
		Description:  "The ID of a workload domain that the cluster belongs to",
		ValidateFunc: validation.NoZeroValues,
	}

	clusterResourceSchema["domain_name"] = &schema.Schema{
		Type:         schema.TypeString,
		Optional:     true,
		Computed:     true,
		Description:  "The name of a workload domain that the cluster belongs to",
		ValidateFunc: validation.NoZeroValues,
	}

	return &schema.Resource{
		CreateContext: resourceClusterCreate,
		ReadContext:   resourceClusterRead,
		UpdateContext: resourceClusterUpdate,
		DeleteContext: resourceClusterDelete,
		Importer: &schema.ResourceImporter{
			StateContext: func(ctx context.Context, data *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
				apiClient := meta.(*api_client.SddcManagerClient).ApiClient
				clusterId := data.Id()
				return cluster.ImportCluster(ctx, data, apiClient, clusterId)
			},
		},
		Schema: clusterResourceSchema,
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(2 * time.Hour),
			Read:   schema.DefaultTimeout(10 * time.Minute),
			Update: schema.DefaultTimeout(2 * time.Hour),
			Delete: schema.DefaultTimeout(1 * time.Hour),
		},
	}
}

// clusterSubresourceSchema this helper function extracts the Cluster schema, so that
// it's made available for merging in the Domain resource schema.
func clusterSubresourceSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "ID of the cluster",
			},
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				Description:  "Name of the cluster to add to the workload domain",
				ValidateFunc: validation.NoZeroValues,
			},
			"host": {
				Type:        schema.TypeList,
				Required:    true,
				Description: "List of ESXi host information from the free pool to consume in a workload domain",
				MinItems:    2,
				Elem:        cluster.HostSpecSchema(),
			},
			"cluster_image_id": {
				Type:         schema.TypeString,
				Optional:     true,
				Description:  "ID of the cluster image to be used with the cluster",
				ValidateFunc: validation.NoZeroValues,
			},
			"evc_mode": {
				Type:     schema.TypeString,
				Optional: true,
				Description: "EVC mode for new cluster, if needed. One among: INTEL_MEROM, " +
					"INTEL_PENRYN, INTEL_NEALEM, INTEL_WESTMERE, INTEL_SANDYBRIDGE, " +
					"INTEL_IVYBRIDGE, INTEL_HASWELL, INTEL_BROADWELL, INTEL_SKYLAKE, " +
					"INTEL_CASCADELAKE, AMD_REV_E, AMD_REV_F, AMD_GREYHOUND_NO3DNOW, " +
					"AMD_GREYHOUND, AMD_BULLDOZER, AMD_PILEDRIVER, AMD_STREAMROLLER, AMD_ZEN",
				ValidateFunc: validation.StringInSlice([]string{
					"INTEL_MEROM",
					"INTEL_PENRYN",
					"INTEL_NEALEM",
					"INTEL_WESTMERE",
					"INTEL_SANDYBRIDGE",
					"INTEL_IVYBRIDGE",
					"INTEL_HASWELL",
					"INTEL_BROADWELL",
					"INTEL_SKYLAKE",
					"INTEL_CASCADELAKE",
					"AMD_REV_E",
					"AMD_REV_F",
					"AMD_GREYHOUND_NO3DNOW",
					"AMD_GREYHOUND",
					"AMD_BULLDOZER",
					"AMD_PILEDRIVER",
					"AMD_STREAMROLLER",
					"AMD_ZEN"}, true),
				DiffSuppressFunc: func(k, oldValue, newValue string, d *schema.ResourceData) bool {
					return oldValue == strings.ToUpper(newValue) || strings.ToUpper(oldValue) == newValue
				},
			},
			"high_availability_enabled": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "vSphere High Availability settings for the cluster",
			},
			"vsan_datastore": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Cluster storage configuration for vSAN",
				MaxItems:    1,
				Elem:        datastores.VsanDatastoreSchema(),
			},
			"vmfs_datastore": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Cluster storage configuration for VMFS",
				MaxItems:    1,
				Elem:        datastores.VmfsDatastoreSchema(),
			},
			"vsan_remote_datastore_cluster": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Cluster storage configuration for vSAN Remote Datastore",
				MaxItems:    1,
				Elem:        datastores.VsanRemoteDatastoreClusterSchema(),
			},
			"nfs_datastores": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Cluster storage configuration for NFS",
				Elem:        datastores.NfsDatastoreSchema(),
			},
			"vvol_datastores": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Cluster storage configuration for VVOL",
				Elem:        datastores.VvolDatastoreSchema(),
			},
			"geneve_vlan_id": {
				Type:         schema.TypeInt,
				Optional:     true,
				Description:  "VLAN ID use for NSX Geneve in the workload domain",
				ValidateFunc: validation.IntBetween(0, 4095),
			},
			"ip_address_pool": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Description: "Contains the parameters required to create or reuse an IP address pool. Omit for DHCP, " +
					"provide name only to reuse existing IP Pool, if subnets are provided a new IP Pool will be created",
				Elem: network.IpAddressPoolSchema(),
			},
			"vds": {
				Type:        schema.TypeList,
				Required:    true,
				MinItems:    1,
				Description: "vSphere Distributed Switches to add to the cluster",
				Elem:        network.VdsSchema(),
			},
			"vsan_stretch_configuration": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Settings for stretched vSAN clusters",
				MaxItems:    1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"witness_host": {
							Type:        schema.TypeList,
							Optional:    true,
							MaxItems:    1,
							Description: "Configuration for the witness host",
							Elem:        vsan.WitnessHostSubresource(),
						},
						"secondary_fd_host": {
							Type:        schema.TypeList,
							Optional:    true,
							MinItems:    1,
							Description: "The list of hosts that will go into the secondary fault domain",
							Elem:        cluster.HostSpecSchema(),
						},
					},
				},
			},
			"primary_datastore_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Name of the primary datastore",
			},
			"primary_datastore_type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Storage type of the primary datastore",
			},
			"is_default": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Status of the cluster if default or not",
			},
			"is_stretched": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Status of the cluster if stretched or not",
			},
		},
	}
}

func resourceClusterCreate(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	vcfClient := meta.(*api_client.SddcManagerClient)

	clusterSpec, err := cluster.TryConvertResourceDataToClusterSpec(data)
	if err != nil {
		return diag.FromErr(err)
	}

	domainId, err := getDomainId(data, vcfClient.ApiClient)
	if err != nil {
		return diag.FromErr(err)
	}

	clusterId, diagnostics := createCluster(ctx, domainId, *clusterSpec, vcfClient)
	if diagnostics != nil {
		return diagnostics
	}

	data.SetId(clusterId)

	return resourceClusterRead(ctx, data, meta)
}

func resourceClusterRead(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	apiClient := meta.(*api_client.SddcManagerClient).ApiClient

	clusterResult, err := apiClient.GetClusterWithResponse(ctx, data.Id())
	if err != nil {
		return diag.FromErr(err)
	}
	clusterObj, vcfErr := api_client.GetResponseAs[vcf.Cluster](clusterResult)
	if vcfErr != nil {
		api_client.LogError(vcfErr, ctx)
		return diag.FromErr(errors.New(*vcfErr.Message))
	}

	_ = data.Set("primary_datastore_name", clusterObj.PrimaryDatastoreName)
	_ = data.Set("primary_datastore_type", clusterObj.PrimaryDatastoreType)
	_ = data.Set("is_default", clusterObj.IsDefault)
	_ = data.Set("is_stretched", clusterObj.IsStretched)

	return nil
}

func resourceClusterUpdate(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	vcfClient := meta.(*api_client.SddcManagerClient)

	clusterUpdateSpec, err := cluster.CreateClusterUpdateSpec(data, false)
	if err != nil {
		return diag.FromErr(err)
	}

	diagnostics := updateCluster(ctx, data.Id(), *clusterUpdateSpec, vcfClient)
	if diagnostics != nil {
		return diagnostics
	}

	return resourceClusterRead(ctx, data, meta)
}

func resourceClusterDelete(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	vcfClient := meta.(*api_client.SddcManagerClient)

	diagnostics := deleteCluster(ctx, data.Id(), vcfClient)
	if diagnostics != nil {
		return diagnostics
	}

	return nil
}

func createCluster(ctx context.Context, domainId string, clusterSpec vcf.ClusterSpec,
	vcfClient *api_client.SddcManagerClient) (string, diag.Diagnostics) {
	apiClient := vcfClient.ApiClient
	clusterCreationSpec := vcf.ClusterCreationSpec{
		ComputeSpec: vcf.ComputeSpec{
			ClusterSpecs: []vcf.ClusterSpec{clusterSpec},
		},
		DomainId:                 domainId,
		DeployWithoutLicenseKeys: utils.ToPointer[bool](true),
	}

	validateResponse, err := apiClient.ValidateClusterCreationSpecWithResponse(ctx, nil, clusterCreationSpec)
	if err != nil {
		return "", validationUtils.ConvertVcfErrorToDiag(err)
	}
	validationResult, vcfErr := api_client.GetResponseAs[vcf.Validation](validateResponse)
	if vcfErr != nil {
		api_client.LogError(vcfErr, ctx)
		return "", diag.FromErr(errors.New(*vcfErr.Message))
	}
	if validationUtils.HasValidationFailed(validationResult) {
		return "", validationUtils.ConvertValidationResultToDiag(validationResult)
	}

	accepted, err := apiClient.CreateClusterWithResponse(ctx, clusterCreationSpec)
	if err != nil {
		return "", validationUtils.ConvertVcfErrorToDiag(err)
	}
	task, vcfErr := api_client.GetResponseAs[vcf.Task](accepted)
	if vcfErr != nil {
		api_client.LogError(vcfErr, ctx)
		return "", diag.FromErr(errors.New(*vcfErr.Message))
	}
	if err = api_client.NewTaskTracker(ctx, apiClient, *task.Id).WaitForTask(); err != nil {
		return "", diag.FromErr(err)
	}
	clusterId, err := vcfClient.GetResourceIdAssociatedWithTask(ctx, *task.Id, "Cluster")
	if err != nil {
		return "", diag.FromErr(err)
	}
	return clusterId, nil
}

func updateCluster(ctx context.Context, clusterId string, clusterUpdateSpec vcf.ClusterUpdateSpec,
	vcfClient *api_client.SddcManagerClient) diag.Diagnostics {
	apiClient := vcfClient.ApiClient
	validationDiagnostics := cluster.ValidateClusterUpdateOperation(ctx, clusterId, clusterUpdateSpec, apiClient)
	if validationDiagnostics != nil {
		return validationDiagnostics
	}

	acceptedUpdateTask, err := apiClient.UpdateClusterWithResponse(ctx, clusterId, clusterUpdateSpec)
	if err != nil {
		return diag.FromErr(err)
	}

	task, vcfErr := api_client.GetResponseAs[vcf.Task](acceptedUpdateTask)
	if vcfErr != nil {
		api_client.LogError(vcfErr, ctx)
		return diag.FromErr(errors.New(*vcfErr.Message))
	}

	if err = api_client.NewTaskTracker(ctx, apiClient, *task.Id).WaitForTask(); err != nil {
		return diag.FromErr(err)
	}
	return nil
}

func deleteCluster(ctx context.Context, clusterId string, vcfClient *api_client.SddcManagerClient) diag.Diagnostics {
	clusterUpdateSpec, err := cluster.CreateClusterUpdateSpec(nil, true)
	if err != nil {
		return diag.FromErr(err)
	}

	apiClient := vcfClient.ApiClient
	log.Printf("Marking Cluster %s for deletion", clusterId)
	acceptedUpdateRes, err := apiClient.UpdateClusterWithResponse(ctx, clusterId, *clusterUpdateSpec)
	if err != nil {
		return diag.FromErr(err)
	}
	if acceptedUpdateRes.StatusCode() != 200 {
		return diag.FromErr(fmt.Errorf("failed to mark cluster for deletion"))
	}

	log.Printf("Deleting Cluster %s", clusterId)
	acceptedDeleteTask, err := apiClient.DeleteClusterWithResponse(ctx, clusterId, nil)
	if err != nil {
		return diag.FromErr(err)
	}
	task, vcfErr := api_client.GetResponseAs[vcf.Task](acceptedDeleteTask)
	if vcfErr != nil {
		api_client.LogError(vcfErr, ctx)
		return diag.FromErr(errors.New(*vcfErr.Message))
	}

	if err = api_client.NewTaskTracker(ctx, apiClient, *task.Id).WaitForTask(); err != nil {
		return diag.FromErr(err)
	}
	return nil
}

func getDomainId(data *schema.ResourceData, client *vcf.ClientWithResponses) (string, error) {
	domainId := data.Get("domain_id").(string)
	domainName := data.Get("domain_name").(string)

	if domainName != "" {
		if domainId != "" {
			return "", errors.New("you cannot set domain_id and domain_name at the same time")
		}

		domain, err := getDomain(domainName, client)

		if err != nil {
			return "", err
		}

		domainId = *domain.Id
	}

	return domainId, nil
}

func getDomain(name string, client *vcf.ClientWithResponses) (*vcf.Domain, error) {
	ok, err := client.GetDomainsWithResponse(context.Background(), nil)

	if err != nil {
		return nil, err
	}
	domainsList, vcfErr := api_client.GetResponseAs[vcf.PageOfDomain](ok)
	if vcfErr != nil {
		return nil, errors.New(*vcfErr.Message)
	}

	if domainsList != nil && len(*domainsList.Elements) > 0 {
		for _, domain := range *domainsList.Elements {
			if *domain.Name == name {
				return &domain, nil
			}
		}
	}

	return nil, fmt.Errorf("domain %s not found", name)
}
