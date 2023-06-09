/*
 *  Copyright 2023 VMware, Inc.
 *    SPDX-License-Identifier: MPL-2.0
 */

package datastores

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// VsanRemoteDatastoreClusterSchema this helper function extracts the VSAN Datastore Cluster
// schema, so that it's made available for both Domain and Cluster creation.
func VsanRemoteDatastoreClusterSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"datastore_uuids": {
				Type:        schema.TypeList,
				Required:    true,
				Description: "vSAN Remote Datastore UUID",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}
