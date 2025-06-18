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

func GetSddcHostSchema() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeList,
		Required: true,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"credentials": getCredentialsSchema(),
				"hostname": {
					Type:         schema.TypeString,
					Description:  "ESXi hostname. If just the short hostname is provided, then FQDN will be generated using the \"domain\" from dns configuration. Must also adhere to RFC 1123 naming conventions. Example: \"esx-1\" length from 3 to 63",
					Required:     true,
					ValidateFunc: validation.StringLenBetween(3, 63),
				},
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
			},
		},
	}
}

func GetSddcHostSpecsFromSchema(rawData []interface{}) *[]installer.SddcHostSpec {
	var hostSpecs []installer.SddcHostSpec
	for _, rawListEntity := range rawData {
		hostSpecRaw := rawListEntity.(map[string]interface{})
		hostname := hostSpecRaw["hostname"].(string)
		sshThumbprint := utils.ToStringPointer(hostSpecRaw["ssh_thumbprint"])
		sslThumbprint := utils.ToStringPointer(hostSpecRaw["ssl_thumbprint"])

		hostSpec := installer.SddcHostSpec{
			Hostname:      hostname,
			SshThumbprint: sshThumbprint,
			SslThumbprint: sslThumbprint,
		}
		if credentialsData := getCredentialsFromSchema(hostSpecRaw["credentials"].([]interface{})); credentialsData != nil {
			hostSpec.Credentials = credentialsData
		}
		hostSpecs = append(hostSpecs, hostSpec)
	}
	return &hostSpecs
}
