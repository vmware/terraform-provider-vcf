/*
 *  Copyright 2023 VMware, Inc.
 *    SPDX-License-Identifier: MPL-2.0
 */

package network

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"strings"
)

// PortgroupSchema this helper function extracts the Portgroup Schema, so that
// it's made available for both Domain and Cluster creation.
func PortgroupSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				Description:  "Port group name",
				ValidateFunc: validation.NoZeroValues,
			},
			"transport_type": {
				Type:     schema.TypeString,
				Required: true,
				Description: "Port group transport type, One among: VSAN, VMOTION, MANAGEMENT, PUBLIC, " +
					"NFS, VREALIZE, ISCSI, EDGE_INFRA_OVERLAY_UPLINK",
				ValidateFunc: validation.StringInSlice([]string{
					"VSAN", "VMOTION", "MANAGEMENT", "PUBLIC", "NFS", "VREALIZE", "ISCSI", "EDGE_INFRA_OVERLAY_UPLINK",
				}, true),
				DiffSuppressFunc: func(k, oldValue, newValue string, d *schema.ResourceData) bool {
					return oldValue == strings.ToUpper(newValue) || strings.ToUpper(oldValue) == newValue
				},
			},
			"active_uplinks": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "List of active uplinks associated with portgroup. This is only supported for VxRail.",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}
