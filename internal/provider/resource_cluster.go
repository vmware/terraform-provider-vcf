/* Copyright 2023 VMware, Inc.
   SPDX-License-Identifier: MPL-2.0 */

package provider

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/vmware/terraform-provider-vcf/internal/cluster"
	"github.com/vmware/terraform-provider-vcf/internal/datastores"
	"github.com/vmware/terraform-provider-vcf/internal/network"
	validation_utils "github.com/vmware/terraform-provider-vcf/internal/validation"
	"github.com/vmware/vcf-sdk-go/models"
	"strings"
)

func ResourceCluster() *schema.Resource {
	clusterResourceSchema := clusterSubresourceSchema().Schema
	clusterResourceSchema["domain_id"] = &schema.Schema{
		Type:         schema.TypeString,
		Required:     true,
		Description:  "The ID of a domain that the cluster belongs to",
		ValidateFunc: validation.NoZeroValues,
	}

	return &schema.Resource{
		CreateContext: resourceClusterCreate,
		ReadContext:   resourceClusterRead,
		UpdateContext: resourceClusterUpdate,
		DeleteContext: resourceClusterDelete,
		Schema:        clusterResourceSchema,
	}
}

func resourceClusterDelete(ctx context.Context, data *schema.ResourceData, i interface{}) diag.Diagnostics {
	return nil
}

func resourceClusterUpdate(ctx context.Context, data *schema.ResourceData, i interface{}) diag.Diagnostics {
	return nil
}

func resourceClusterRead(ctx context.Context, data *schema.ResourceData, i interface{}) diag.Diagnostics {
	return nil
}

func resourceClusterCreate(ctx context.Context, data *schema.ResourceData, i interface{}) diag.Diagnostics {
	return nil
}

// clusterSubresourceSchema this helper function extracts the Cluster schema, so that
// it's made available for merging in the Domain resource schema.
func clusterSubresourceSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Cluster ID",
			},
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				Description:  "Name of the new cluster that will be added to the specified workload domain",
				ValidateFunc: validation.NoZeroValues,
			},
			"host": {
				Type:        schema.TypeList,
				Required:    true,
				Description: "List of vSphere host information from the free pool to consume in the workload domain",
				MinItems:    1,
				Elem:        cluster.CommissionedHostSchema(),
			},
			"cluster_image_id": {
				Type:         schema.TypeString,
				Optional:     true,
				Description:  "ID of the Cluster Image to be used with the Cluster",
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
				Description: "High availability settings for the cluster",
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
				Description:  "Vlan id of Geneve for NSX-T based workload domains",
				ValidateFunc: validation.IntBetween(0, 4095),
			},
			"vds": {
				Type:        schema.TypeList,
				Required:    true,
				MinItems:    1,
				Description: "Distributed switches to add to the cluster",
				Elem:        network.VdsSchema(),
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
			"is_streched": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Status of the cluster if Stretched or not",
			},
		},
	}
}

func tryConvertToClusterSpec(object map[string]interface{}) (*models.ClusterSpec, error) {
	if object == nil {
		return nil, fmt.Errorf("cannot convert to ClusterSpec, object is nil")
	}
	name := object["name"].(string)
	if len(name) == 0 {
		return nil, fmt.Errorf("cannot convert to ClusterSpec, name is required")
	}
	result := &models.ClusterSpec{}
	result.Name = &name
	if clusterImageId, ok := object["cluster_image_id"]; ok && !validation_utils.IsEmpty(clusterImageId) {
		result.ClusterImageID = clusterImageId.(string)
	}
	if evcMode, ok := object["evc_mode"]; ok && len(evcMode.(string)) > 0 {
		if result.AdvancedOptions == nil {
			result.AdvancedOptions = &models.AdvancedOptions{}
		}
		result.AdvancedOptions.EvcMode = evcMode.(string)
	}
	if highAvailabilityEnabled, ok := object["high_availability_enabled"]; ok && !validation_utils.IsEmpty(highAvailabilityEnabled) {
		if result.AdvancedOptions == nil {
			result.AdvancedOptions = &models.AdvancedOptions{}
		}
		result.AdvancedOptions.HighAvailability = &models.HighAvailability{
			Enabled: toBoolPointer(highAvailabilityEnabled),
		}
	}

	result.NetworkSpec = &models.NetworkSpec{}
	result.NetworkSpec.NsxClusterSpec = &models.NsxClusterSpec{}
	result.NetworkSpec.NsxClusterSpec.NsxTClusterSpec = &models.NsxTClusterSpec{}

	if geneveVlanId, ok := object["geneve_vlan_id"]; ok && !validation_utils.IsEmpty(geneveVlanId) {
		result.NetworkSpec.NsxClusterSpec.NsxTClusterSpec.GeneveVlanID = int32(geneveVlanId.(int))
	}

	if hostsRaw, ok := object["host"]; ok {
		hostsList := hostsRaw.([]interface{})
		if len(hostsList) > 0 {
			result.HostSpecs = []*models.HostSpec{}
			for _, hostListEntry := range hostsList {
				hostSpec, err := cluster.TryConvertToHostSpec(hostListEntry.(map[string]interface{}))
				if err != nil {
					return nil, err
				}
				result.HostSpecs = append(result.HostSpecs, hostSpec)
			}
		} else {
			return nil, fmt.Errorf("cannot convert to ClusterSpec, hosts list is empty")
		}
	} else {
		return nil, fmt.Errorf("cannot convert to ClusterSpec, hosts list is not set")
	}

	if vdsRaw, ok := object["vds"]; ok {
		vdsList := vdsRaw.([]interface{})
		if len(vdsList) > 0 {
			result.NetworkSpec.VdsSpecs = []*models.VdsSpec{}
			for _, vdsListEntry := range vdsList {
				vdsSpec, err := network.TryConvertToVdsSpec(vdsListEntry.(map[string]interface{}))
				if err != nil {
					return nil, err
				}
				result.NetworkSpec.VdsSpecs = append(result.NetworkSpec.VdsSpecs, vdsSpec)
			}
		} else {
			return nil, fmt.Errorf("cannot convert to ClusterSpec, vds list is empty")
		}
	} else {
		return nil, fmt.Errorf("cannot convert to ClusterSpec, vds list is not set")
	}

	datastoreSpec, err := tryConvertToClusterDatastoreSpec(object, name)
	if err != nil {
		return nil, err
	} else {
		result.DatastoreSpec = datastoreSpec
	}

	return result, nil
}

func tryConvertToClusterDatastoreSpec(object map[string]interface{}, clusterName string) (*models.DatastoreSpec, error) {
	result := &models.DatastoreSpec{}
	atLeastOneTypeOfDatastoreConfigured := false
	if vsanDatastoreRaw, ok := object["vsan_datastore"]; ok && !validation_utils.IsEmpty(vsanDatastoreRaw) {
		if len(vsanDatastoreRaw.([]interface{})) > 1 {
			return nil, fmt.Errorf("more than one vsan_datastore config for cluster %q", clusterName)
		}
		vsanDatastoreListEntry := vsanDatastoreRaw.([]interface{})[0]
		vsanDatastoreSpec, err := datastores.TryConvertToVsanDatastoreSpec(vsanDatastoreListEntry.(map[string]interface{}))
		if err != nil {
			return nil, err
		}
		atLeastOneTypeOfDatastoreConfigured = true
		result.VSANDatastoreSpec = vsanDatastoreSpec
	}
	if vmfsDatastoreRaw, ok := object["vmfs_datastore"]; ok && !validation_utils.IsEmpty(vmfsDatastoreRaw) {
		if len(vmfsDatastoreRaw.([]interface{})) > 1 {
			return nil, fmt.Errorf("more than one vmfs_datastore config for cluster %q", clusterName)
		}
		vmfsDatastoreListEntry := vmfsDatastoreRaw.([]interface{})[0]
		vmfsDatastoreSpec, err := datastores.TryConvertToVmfsDatastoreSpec(vmfsDatastoreListEntry.(map[string]interface{}))
		if err != nil {
			return nil, err
		}
		atLeastOneTypeOfDatastoreConfigured = true
		result.VmfsDatastoreSpec = vmfsDatastoreSpec
	}
	if vsanRemoteDatastoreClusterRaw, ok := object["vsan_remote_datastore_cluster"]; ok && !validation_utils.IsEmpty(vsanRemoteDatastoreClusterRaw) {
		if len(vsanRemoteDatastoreClusterRaw.([]interface{})) > 1 {
			return nil, fmt.Errorf("more than one vsan_remote_datastore_cluster config for cluster %q", clusterName)
		}
		vsanRemoteDatastoreClusterListEntry := vsanRemoteDatastoreClusterRaw.([]interface{})[0]
		vsanRemoteDatastoreClusterSpec, err := datastores.TryConvertToVSANRemoteDatastoreClusterSpec(
			vsanRemoteDatastoreClusterListEntry.(map[string]interface{}))
		if err != nil {
			return nil, err
		}
		atLeastOneTypeOfDatastoreConfigured = true
		result.VSANRemoteDatastoreClusterSpec = vsanRemoteDatastoreClusterSpec
	}
	if nfsDatastoresRaw, ok := object["nfs_datastores"]; ok && !validation_utils.IsEmpty(nfsDatastoresRaw) {
		nfsDatastoresList := nfsDatastoresRaw.([]interface{})
		if len(nfsDatastoresList) > 0 {
			result.NfsDatastoreSpecs = []*models.NfsDatastoreSpec{}
			for _, nfsDatastoreListEntry := range nfsDatastoresList {
				nfsDatastoreSpec, err := datastores.TryConvertToNfsDatastoreSpec(
					nfsDatastoreListEntry.(map[string]interface{}))
				if err != nil {
					return nil, err
				}
				result.NfsDatastoreSpecs = append(result.NfsDatastoreSpecs, nfsDatastoreSpec)
			}
			atLeastOneTypeOfDatastoreConfigured = true
		}
	}
	if vvolDatastoresRaw, ok := object["vvol_datastores"]; ok && !validation_utils.IsEmpty(vvolDatastoresRaw) {
		vvolDatastoresList := vvolDatastoresRaw.([]interface{})
		if len(vvolDatastoresList) > 0 {
			result.VvolDatastoreSpecs = []*models.VvolDatastoreSpec{}
			for _, vvolDatastoreListEntry := range vvolDatastoresList {
				vvolDatastoreSpec, err := datastores.TryConvertToVvolDatastoreSpec(
					vvolDatastoreListEntry.(map[string]interface{}))
				if err != nil {
					return nil, err
				}
				result.VvolDatastoreSpecs = append(result.VvolDatastoreSpecs, vvolDatastoreSpec)
			}
			atLeastOneTypeOfDatastoreConfigured = true
		}
	}
	if !atLeastOneTypeOfDatastoreConfigured {
		return nil, fmt.Errorf("at least one type of datastore configuration required for cluster %q", clusterName)
	}

	return result, nil
}

func toBoolPointer(object interface{}) *bool {
	if object == nil {
		return nil
	}
	objectAsBool := object.(bool)
	return &objectAsBool
}
