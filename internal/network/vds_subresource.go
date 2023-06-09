/*
 *  Copyright 2023 VMware, Inc.
 *    SPDX-License-Identifier: MPL-2.0
 */

package network

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

// VdsSchema this helper function extracts the VDS Schema, so that
// it's made available for both Domain and Cluster creation.
// This specification contains vSphere distributed switch configurations.
func VdsSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				Description:  "vSphere Distributed Switch name",
				ValidateFunc: validation.NoZeroValues,
			},
			"is_used_by_nsxt": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Identifies if the vSphere distributed switch is used by NSX-T",
			},
			"portgroups": {
				Type:        schema.TypeList,
				Description: "List of portgroups to be associated with the vSphere Distributed Switch",
				Elem:        PortgroupSchema(),
			},
			"nioc_bandwidth_allocations": {
				Type:        schema.TypeList,
				Description: "List of Network I/O Control Bandwidth Allocations for System Traffic",
				Elem:        NiocBandwidthAllocationSchema(),
			},
		},
	}
}
