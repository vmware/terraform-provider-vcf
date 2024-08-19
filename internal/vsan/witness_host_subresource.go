// Copyright 2024 Broadcom. All Rights Reserved.
// SPDX-License-Identifier: MPL-2.0

package vsan

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	validation_utils "github.com/vmware/terraform-provider-vcf/internal/validation"
)

func WitnessHostSubresource() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"vsan_ip": {
				Type:         schema.TypeString,
				Required:     true,
				Description:  "IP address for the witness host on the vSAN network",
				ValidateFunc: validation.IsIPv4Address,
			},
			"vsan_cidr": {
				Type:         schema.TypeString,
				Required:     true,
				Description:  "CIDR address for the witness host on the vSAN network",
				ValidateFunc: validation_utils.ValidateCidrIPv4AddressSchema,
			},
			"fqdn": {
				Type:         schema.TypeString,
				Required:     true,
				Description:  "Fully qualified domain name of the witness host. It should be routable on the vSAN network",
				ValidateFunc: validation.NoZeroValues,
			},
		},
	}
}
