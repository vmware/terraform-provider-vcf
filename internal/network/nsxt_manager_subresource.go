/*
 *  Copyright 2023 VMware, Inc.
 *    SPDX-License-Identifier: MPL-2.0
 */

package network

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	validation_utils "github.com/vmware/terraform-provider-vcf/internal/validation"
)

// NsxtManagerSchema this helper function extracts the Nsxt Manager schema, which contains
// the parameters required to install and configure NSX Manager in a workload domain.
func NsxtManagerSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				Description:  "Name of the NSX Manager virtual machine",
				ValidateFunc: validation.NoZeroValues,
			},
			"ip_address": {
				Type:         schema.TypeString,
				Required:     true,
				Description:  "IPv4 address of the virtual machine",
				ValidateFunc: validation_utils.ValidateIPv4AddressSchema,
			},
			"dns_name": {
				Type:         schema.TypeString,
				Optional:     true,
				Description:  "DNS name of the virtual machine, e.g., vc-1.domain1.rainpole.io",
				ValidateFunc: validation.NoZeroValues,
			},
			"subnet_mask": {
				Type:         schema.TypeString,
				Optional:     true,
				Description:  "Subnet mask",
				ValidateFunc: validation_utils.ValidateIPv4AddressSchema,
			},
			"gateway": {
				Type:         schema.TypeString,
				Optional:     true,
				Description:  "IPv4 gateway the VM can use to connect to the outside world",
				ValidateFunc: validation_utils.ValidateIPv4AddressSchema,
			},
		},
	}
}
