/* Copyright 2023 VMware, Inc.
   SPDX-License-Identifier: MPL-2.0 */

package datastores

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"strings"
)

// VvolDatastoreSchema this helper function extracts the VVOL Datastore schema, so that
// it's made available for both Domain and Cluster creation.
func VvolDatastoreSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"datastore_name": {
				Type:         schema.TypeString,
				Required:     true,
				Description:  "Datastore name used for cluster creation",
				ValidateFunc: validation.NoZeroValues,
			},
			"storage_container_id": {
				Type:         schema.TypeString,
				Required:     true,
				Description:  "UUID of the VASA storage container",
				ValidateFunc: validation.IsUUID,
			},
			"storage_protocol_type": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Type of the VASA storage protocol. One among: ISCSI, NFS, FC",
				ValidateFunc: validation.StringInSlice(
					[]string{"ISCSI", "NFS", "FC"}, true),
				DiffSuppressFunc: func(k, oldValue, newValue string, d *schema.ResourceData) bool {
					return oldValue == strings.ToUpper(newValue) || strings.ToUpper(oldValue) == newValue
				},
			},
			"user_id": {
				Type:         schema.TypeString,
				Required:     true,
				Description:  "UUID of the VASA storage user",
				ValidateFunc: validation.IsUUID,
			},
			"vasa_provider_id": {
				Type:         schema.TypeString,
				Required:     true,
				Description:  "UUID of the VASA storage provider",
				ValidateFunc: validation.IsUUID,
			},
		},
	}
}
