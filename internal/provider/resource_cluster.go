/* Copyright 2023 VMware, Inc.
   SPDX-License-Identifier: MPL-2.0 */

package provider

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/vmware/terraform-provider-vcf/internal/cluster"
	"github.com/vmware/terraform-provider-vcf/internal/datastores"
	"github.com/vmware/terraform-provider-vcf/internal/network"
	"strings"
)

// clusterSchema this helper function extracts the Cluster schema, so that
// it's made available for merging in the Domain resource schema.
func clusterSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"name": {
			Type:         schema.TypeString,
			Required:     true,
			Description:  "Name of the new cluster that will be added to the specified workload domain",
			ValidateFunc: validation.NoZeroValues,
		},
		"hosts": {
			Type:        schema.TypeList,
			Required:    true,
			Description: "List of vSphere host information from the free pool to consume in the workload domain",
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
			Type:        schema.TypeInt,
			Optional:    true,
			Description: "Vlan id of Geneve for NSX-T based workload domains",
		},
		"vds": {
			Type:        schema.TypeList,
			Required:    true,
			Description: "Distributed switches to add to the cluster",
			Elem:        network.VdsSchema(),
		},
	}
}
