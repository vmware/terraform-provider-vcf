/*
 *  Copyright 2023 VMware, Inc.
 *    SPDX-License-Identifier: MPL-2.0
 */

package sddc

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	utils "github.com/vmware/terraform-provider-vcf/internal/resource_utils"
	"github.com/vmware/vcf-sdk-go/models"
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
					Description: "Tenant domain. Example: rainpole.io",
					Required:    true,
				},
				"name_server": {
					Type:         schema.TypeString,
					Description:  "Primary nameserver to be configured for vCenter/PSC/ESXi's/NSX. Example: 172.0.0.4",
					Optional:     true,
					ValidateFunc: validation.IsIPAddress,
				},
				"secondary_name_server": {
					Type:         schema.TypeString,
					Description:  "Secondary nameserver to be configured for vCenter/PSC/ESXi's/NSX. Example: 172.0.0.5",
					Optional:     true,
					ValidateFunc: validation.IsIPAddress,
				},
				"subdomain": {
					Type:        schema.TypeString,
					Description: "Tenant Sub domain. Example: sfo.rainpole.io",
					Required:    true,
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
	subdomain := utils.ToStringPointer(data["subdomain"])

	dnsSpecBinding := &models.DNSSpec{
		Domain:              domain,
		Nameserver:          nameServer,
		SecondaryNameserver: secondaryNameserver,
		Subdomain:           subdomain,
	}
	return dnsSpecBinding
}
