// © Broadcom. All Rights Reserved.
// The term “Broadcom” refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

package nsx_edge_cluster

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func ClusterProfileSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				Description:  "The name of the profile",
				ValidateFunc: validation.NoZeroValues,
			},
			"bfd_allowed_hop": {
				Type:         schema.TypeInt,
				Required:     true,
				Description:  "BFD allowed hop",
				ValidateFunc: validation.IntBetween(1, 255),
			},
			"bfd_declare_dead_multiple": {
				Type:         schema.TypeInt,
				Required:     true,
				Description:  "BFD declare dead multiple",
				ValidateFunc: validation.IntBetween(2, 16),
			},
			"bfd_probe_interval": {
				Type:         schema.TypeInt,
				Required:     true,
				Description:  "BFD probe interval",
				ValidateFunc: validation.IntBetween(500, 60000),
			},
			"standby_relocation_threshold": {
				Type:         schema.TypeInt,
				Required:     true,
				Description:  "Standby relocation threshold",
				ValidateFunc: validation.IntBetween(10, 1000),
			},
		},
	}
}
