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

var esxiCertsModes = []string{"Custom", "VMCA"}

func GetSecuritySchema() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeList,
		MaxItems: 1,
		Optional: true,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"esxi_certs_mode": {
					Type:         schema.TypeString,
					Description:  "ESXi certificates mode. One among: Custom, VMCA",
					Optional:     true,
					ValidateFunc: validation.StringInSlice(esxiCertsModes, false),
				},
				"root_ca_certs": getRootCaCertsSchema(),
			},
		},
	}
}

func getRootCaCertsSchema() *schema.Schema {
	return &schema.Schema{
		Type:        schema.TypeList,
		Description: "Root Certificate Authority certificate list",
		Optional:    true,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"alias": {
					Type:        schema.TypeString,
					Description: "Certificate alias",
					Optional:    true,
				},
				"cert_chain": {
					Type:        schema.TypeList,
					Description: "List of Base64 encoded certificates",
					Optional:    true,
					Elem: &schema.Schema{
						Type: schema.TypeString,
					},
				},
			},
		},
	}
}

func GetSecuritySpecSchema(rawData []interface{}) *installer.SecuritySpec {
	if len(rawData) <= 0 {
		return nil
	}
	data := rawData[0].(map[string]interface{})
	esxiCertsMode := data["esxi_certs_mode"].(string)

	securitySpecBinding := &installer.SecuritySpec{
		EsxiCertsMode: &esxiCertsMode,
	}
	if rootCaCerts := getRootCaCertsBindingFromSchema(data["root_ca_certs"].([]interface{})); len(rootCaCerts) > 0 {
		securitySpecBinding.RootCaCerts = &rootCaCerts
	}

	return securitySpecBinding
}

func getRootCaCertsBindingFromSchema(rawData []interface{}) []installer.RootCaCerts {
	var rootCaCertsBindingsList []installer.RootCaCerts
	for _, rootCaCerts := range rawData {
		data := rootCaCerts.(map[string]interface{})
		alias := data["alias"].(string)

		rootCaCertsBinding := installer.RootCaCerts{
			Alias: &alias,
		}
		if certChain, ok := data["cert_chain"].([]interface{}); ok {
			certChainSlice := utils.ToStringSlice(certChain)
			rootCaCertsBinding.CertChain = &certChainSlice
		}

		rootCaCertsBindingsList = append(rootCaCertsBindingsList, rootCaCertsBinding)
	}
	return rootCaCertsBindingsList
}
