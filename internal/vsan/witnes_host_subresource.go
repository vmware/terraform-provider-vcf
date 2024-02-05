// Copyright 2024 Broadcom. All Rights Reserved.
// SPDX-License-Identifier: MPL-2.0

package vsan

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	validationUtils "github.com/vmware/terraform-provider-vcf/internal/validation"
)

func WitnessHostSubresource() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"vsan_ip": {
				Type:         schema.TypeString,
				Required:     true,
				Description:  "TODO",
				ValidateFunc: validation.IsIPv4Address,
			},
			"vsan_cidr": {
				Type:         schema.TypeString,
				Required:     true,
				Description:  "TODO",
				ValidateFunc: validationUtils.ValidateCidrIPv4AddressSchema,
			},
			"fqdn": {
				Type:         schema.TypeString,
				Required:     true,
				Description:  "TODO",
				ValidateFunc: validation.NoZeroValues,
			},
		},
	}
}
