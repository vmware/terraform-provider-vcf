/* Copyright 2023 VMware, Inc.
   SPDX-License-Identifier: MPL-2.0 */

package datastores

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// VmfsDatastoreSchema this helper function extracts the VMFS Datastore schema, so that
// it's made available for both Domain and Cluster creation.
func VmfsDatastoreSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"datastore_names": {
				Type:        schema.TypeList,
				Required:    true,
				Description: "Datastore names used for VMFS on FC for cluster creation",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}
