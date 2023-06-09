/* Copyright 2023 VMware, Inc.
   SPDX-License-Identifier: MPL-2.0 */

package datastores

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

// VsanDatastoreSchema this helper function extracts the VSAN Datastore schema, so that
// it's made available for both Domain and Cluster creation.
func VsanDatastoreSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"datastore_name": {
				Type:         schema.TypeString,
				Required:     true,
				Description:  "Datastore name used for cluster creation",
				ValidateFunc: validation.NoZeroValues,
			},
			"dedup_and_compression_enabled": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Enable vSAN deduplication and compression",
			},
			"failures_to_tolerate": {
				Type:         schema.TypeInt,
				Optional:     true,
				Description:  "Number of vSphere host failures to tolerate in the vSAN cluster (can be 0, 1, or 2)",
				ValidateFunc: validation.IntBetween(0, 2),
			},
			"license_key": {
				Type:         schema.TypeString,
				Optional:     true,
				Description:  "License key for the vSAN data store to be applied in vCenter",
				ValidateFunc: validation.NoZeroValues,
			},
		},
	}
}
