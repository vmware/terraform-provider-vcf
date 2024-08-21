// © Broadcom. All Rights Reserved.
// The term “Broadcom” refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

package sddc

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/vmware/vcf-sdk-go/models"

	utils "github.com/vmware/terraform-provider-vcf/internal/resource_utils"
)

func GetSddcHostSchema() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeList,
		Required: true,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"association": {
					Type:        schema.TypeString,
					Description: "Host Association: Location/Datacenter",
					Required:    true,
				},
				"credentials": getCredentialsSchema(),
				"hostname": {
					Type:         schema.TypeString,
					Description:  "ESXi hostname. If just the short hostname is provided, then FQDN will be generated using the \"domain\" from dns configuration. Must also adhere to RFC 1123 naming conventions. Example: \"esx-1\" length from 3 to 63",
					Required:     true,
					ValidateFunc: validation.StringLenBetween(3, 63),
				},
				"ip_address_private": getIPAllocationSchema(),
				"ssh_thumbprint": {
					Type:        schema.TypeString,
					Description: "Host SSH thumbprint (RSA SHA256)",
					Optional:    true,
				},
				"ssl_thumbprint": {
					Type:        schema.TypeString,
					Description: "Host SSH thumbprint (RSA SHA256)",
					Optional:    true,
				},
				"vswitch": {
					Type:        schema.TypeString,
					Description: "Host vSwitch name",
					Required:    true,
				},
			},
		},
	}
}

func getIPAllocationSchema() *schema.Schema {
	return &schema.Schema{
		Type:        schema.TypeList,
		Description: "Host Private Management IP",
		MaxItems:    1,
		Required:    true,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"cidr": {
					Type:        schema.TypeString,
					Description: "Classless Inter-Domain Routing (CIDR), Example: 172.0.0.0/24",
					Optional:    true,
				},
				"gateway": {
					Type:         schema.TypeString,
					Description:  "Gateway",
					Required:     true,
					ValidateFunc: validation.IsIPAddress,
				},
				"ip_address": {
					Type:         schema.TypeString,
					Description:  "IP Address",
					Required:     true,
					ValidateFunc: validation.IsIPAddress,
				},
				"subnet": {
					Type:         schema.TypeString,
					Description:  "Subnet",
					Optional:     true,
					ValidateFunc: validation.IsIPAddress,
				},
			},
		},
	}
}

func GetSddcHostSpecsFromSchema(rawData []interface{}) []*models.SDDCHostSpec {
	var hostSpecs []*models.SDDCHostSpec
	for _, rawListEntity := range rawData {
		hostSpecRaw := rawListEntity.(map[string]interface{})
		association := utils.ToStringPointer(hostSpecRaw["association"])
		hostname := utils.ToStringPointer(hostSpecRaw["hostname"])
		sshThumbprint := hostSpecRaw["ssh_thumbprint"].(string)
		sslThumbprint := hostSpecRaw["ssl_thumbprint"].(string)
		vswitch := utils.ToStringPointer(hostSpecRaw["vswitch"])

		hostSpec := &models.SDDCHostSpec{
			Association:   association,
			Hostname:      hostname,
			SSHThumbprint: sshThumbprint,
			SSLThumbprint: sslThumbprint,
			VSwitch:       vswitch,
		}
		if credentialsData := getCredentialsFromSchema(hostSpecRaw["credentials"].([]interface{})); credentialsData != nil {
			hostSpec.Credentials = credentialsData
		}
		if ipAllocation := getIPAllocationBindingFromSchema(hostSpecRaw["ip_address_private"].([]interface{})); ipAllocation != nil {
			hostSpec.IPAddressPrivate = ipAllocation
		}
		hostSpecs = append(hostSpecs, hostSpec)
	}
	return hostSpecs
}

func getIPAllocationBindingFromSchema(rawData []interface{}) *models.IPAllocation {
	if len(rawData) <= 0 {
		return nil
	}
	data := rawData[0].(map[string]interface{})
	cidr := data["cidr"].(string)
	gateway := data["gateway"].(string)
	ipAddress := utils.ToStringPointer(data["ip_address"])
	subnet := data["subnet"].(string)

	ipAllocationBinding := &models.IPAllocation{
		Cidr:      cidr,
		Gateway:   gateway,
		IPAddress: ipAddress,
		Subnet:    subnet,
	}
	return ipAllocationBinding
}
