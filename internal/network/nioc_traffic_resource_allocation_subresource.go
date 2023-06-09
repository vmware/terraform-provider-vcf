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

// NiocTrafficResourceAllocationSchema this helper function extracts the NiocTrafficResourceAllocation
// Schema, so that it's made available for both Domain and Cluster creation.
func NiocTrafficResourceAllocationSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"limit": {
				Type:         schema.TypeInt,
				Required:     true,
				Description:  "Specify network traffic limit allocation",
				ValidateFunc: validation.NoZeroValues,
			},
			"reservation": {
				Type:         schema.TypeInt,
				Required:     true,
				Description:  "Specify network traffic reservation allocation",
				ValidateFunc: validation.NoZeroValues,
			},
			"shares": {
				Type:         schema.TypeInt,
				Optional:     true,
				Description:  "The number of shares allocated",
				ValidateFunc: validation.NoZeroValues,
			},
			"shares_level": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The allocation level. One among: low, normal, high, custom",
				ValidateFunc: validation.StringInSlice([]string{
					"low", "normal", "high", "custom",
				}, true),
				DiffSuppressFunc: func(k, oldValue, newValue string, d *schema.ResourceData) bool {
					return oldValue == strings.ToUpper(newValue) || strings.ToUpper(oldValue) == newValue
				},
			},
		},
	}
}
