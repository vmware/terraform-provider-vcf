// © Broadcom. All Rights Reserved.
// The term “Broadcom” refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

package certificates

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/vmware/vcf-sdk-go/vcf"
)

func CsrSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"csr_pem": {
				Type:        schema.TypeString,
				Description: "The CSR encoded content",
				Computed:    true,
			},
			"csr_string": {
				Type:        schema.TypeString,
				Description: "The CSR decoded content",
				Computed:    true,
			},
			"resource": {
				Type:        schema.TypeList,
				Description: "Resource associated with CSR",
				Computed:    true,
				Elem:        ResourceSchema(),
			},
		},
	}
}

func ResourceSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"resource_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Resource ID",
			},
			"type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Resource type",
			},
			"fqdn": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Resource FQDN",
			},
			"name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Name of the resource",
			},
		},
	}
}

func FlattenCsr(csr *vcf.Csr) map[string]interface{} {
	result := make(map[string]interface{})
	result["csr_pem"] = *csr.CsrEncodedContent
	result["csr_string"] = *csr.CsrDecodedContent
	flattenedResource := FlattenResource(csr.Resource)
	result["resource"] = []interface{}{flattenedResource}
	return result
}

func FlattenResource(resource *vcf.Resource) map[string]interface{} {
	result := make(map[string]interface{})

	result["resource_id"] = resource.ResourceId
	result["type"] = resource.Type
	result["fqdn"] = resource.Fqdn
	result["name"] = resource.Name

	return result
}
