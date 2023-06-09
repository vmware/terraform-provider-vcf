/* Copyright 2023 VMware, Inc.
   SPDX-License-Identifier: MPL-2.0 */

package datastores

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

// NfsDatastoreSchema this helper function extracts the NFS Datastore schema, so that
// it's made available for both Domain and Cluster creation.
func NfsDatastoreSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"datastore_name": {
				Type:         schema.TypeString,
				Required:     true,
				Description:  "Datastore name used for cluster creation",
				ValidateFunc: validation.NoZeroValues,
			},
			"path": {
				Type:         schema.TypeString,
				Required:     true,
				Description:  "Shared directory path used for NFS based cluster creation",
				ValidateFunc: validation.NoZeroValues,
			},
			"read_only": {
				Type:        schema.TypeBool,
				Required:    true,
				Description: "Readonly is used to identify whether to mount the directory as readOnly or not",
			},
			"server_name": {
				Type:         schema.TypeString,
				Required:     true,
				Description:  "NFS Server name used for cluster creation",
				ValidateFunc: validation.NoZeroValues,
			},
			"user_tag": {
				Type:         schema.TypeString,
				Optional:     true,
				Description:  "User tag used to annotate NFS share",
				ValidateFunc: validation.NoZeroValues,
			},
		},
	}
}
