// © Broadcom. All Rights Reserved.
// The term “Broadcom” refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

package certificates

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/vmware/vcf-sdk-go/models"
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

func FlattenCsr(csr *models.CSR) map[string]interface{} {
	result := make(map[string]interface{})
	result["csr_pem"] = *csr.CSREncodedContent
	result["csr_string"] = *csr.CSRDecodedContent
	flattenedResource := FlattenResource(csr.Resource)
	result["resource"] = []interface{}{flattenedResource}
	return result
}

func FlattenResource(resource *models.Resource) map[string]interface{} {
	result := make(map[string]interface{})

	if resource.ResourceID != nil {
		result["resource_id"] = *resource.ResourceID
	}
	if resource.Type != nil {
		result["type"] = *resource.Type
	}
	result["fqdn"] = resource.Fqdn
	result["name"] = resource.Name

	return result
}
