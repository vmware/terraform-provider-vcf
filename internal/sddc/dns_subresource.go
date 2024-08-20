// Copyright 2023 Broadcom. All Rights Reserved.
// SPDX-License-Identifier: MPL-2.0

package sddc

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/vmware/vcf-sdk-go/models"

	utils "github.com/vmware/terraform-provider-vcf/internal/resource_utils"
)

func GetDnsSchema() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeList,
		Required: true,
		MaxItems: 1,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"domain": {
					Type:        schema.TypeString,
					Description: "Tenant domain. Parent tenant domain including TLD suffix Example: vmware.com",
					Required:    true,
				},
				"name_server": {
					Type:         schema.TypeString,
					Description:  "Primary nameserver IPv4 address. Example: 172.0.0.4",
					Optional:     true,
					ValidateFunc: validation.IsIPAddress,
				},
				"secondary_name_server": {
					Type:         schema.TypeString,
					Description:  "Secondary nameserver IPv4 address. Example: 172.0.0.5",
					Optional:     true,
					ValidateFunc: validation.IsIPAddress,
				},
			},
		},
	}
}

func GetDnsSpecFromSchema(rawData []interface{}) *models.DNSSpec {
	if len(rawData) <= 0 {
		return nil
	}
	data := rawData[0].(map[string]interface{})
	domain := utils.ToStringPointer(data["domain"])
	nameServer := data["name_server"].(string)
	secondaryNameserver := data["secondary_name_server"].(string)

	dnsSpecBinding := &models.DNSSpec{
		Nameserver:          nameServer,
		SecondaryNameserver: secondaryNameserver,
		Domain:              domain,
		Subdomain:           domain,
	}
	return dnsSpecBinding
}
