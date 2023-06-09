/*
 *  Copyright 2023 VMware, Inc.
 *    SPDX-License-Identifier: MPL-2.0
 */

package network

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

// VMNicSchema this helper function extracts the VMNic Schema, so that
// it's made available for both Domain and Cluster creation.
func VMNicSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"id": {
				Type:         schema.TypeString,
				Required:     true,
				Description:  "VmNic ID of vSphere host to be associated with VDS, once added to cluster",
				ValidateFunc: validation.NoZeroValues,
			},
			"move_to_nvds": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "This flag determines if the vmnic must be on N-VDS",
			},
			"uplink": {
				Type:         schema.TypeString,
				Optional:     true,
				Description:  "Uplink to be associated with vmnic",
				ValidateFunc: validation.NoZeroValues,
			},
			"vds_name": {
				Type:         schema.TypeString,
				Optional:     true,
				Description:  "VDS name to associate with vSphere host",
				ValidateFunc: validation.NoZeroValues,
			},
		},
	}
}
