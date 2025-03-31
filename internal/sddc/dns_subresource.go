// © Broadcom. All Rights Reserved.
// The term “Broadcom” refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

package sddc

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	utils "github.com/vmware/terraform-provider-vcf/internal/resource_utils"
	"github.com/vmware/vcf-sdk-go/installer"
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

func GetDnsSpecFromSchema(rawData []interface{}) *installer.DnsSpec {
	if len(rawData) <= 0 {
		return nil
	}
	data := rawData[0].(map[string]interface{})
	domain := utils.ToStringPointer(data["domain"])
	nameServer := data["name_server"].(string)
	secondaryNameserver := data["secondary_name_server"].(string)

	nameservers := make([]string, 0)
	if len(nameServer) > 0 {
		nameservers = append(nameservers, nameServer)
	}

	if len(secondaryNameserver) > 0 {
		nameservers = append(nameservers, secondaryNameserver)
	}

	dnsSpecBinding := &installer.DnsSpec{
		Nameservers: &nameservers,
		Subdomain:   *domain,
	}
	return dnsSpecBinding
}
