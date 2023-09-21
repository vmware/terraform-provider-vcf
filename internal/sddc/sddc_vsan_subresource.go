/*
 *  Copyright 2023 VMware, Inc.
 *    SPDX-License-Identifier: MPL-2.0
 */

package sddc

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	utils "github.com/vmware/terraform-provider-vcf/internal/resource_utils"
	"github.com/vmware/vcf-sdk-go/models"
)

func GetVsanSchema() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeList,
		Optional: true,
		MaxItems: 1,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"datastore_name": {
					Type:        schema.TypeString,
					Description: "Datastore Name",
					Required:    true,
				},
				"hcl_file": {
					Type:        schema.TypeString,
					Description: "HCL File",
					Optional:    true,
				},
				"license_file": {
					Type:        schema.TypeString,
					Description: "License File",
					Optional:    true,
				},
				"vsan_dedup": {
					Type:        schema.TypeBool,
					Description: "VSAN feature Deduplication and Compression flag, one flag for both features",
					Optional:    true,
				},
			},
		},
	}
}

func GetVsanSpecFromSchema(rawData []interface{}) *models.VSANSpec {
	if len(rawData) <= 0 {
		return nil
	}
	data := rawData[0].(map[string]interface{})
	datastoreName := data["datastore_name"].(string)
	hclFile := data["hcl_file"].(string)
	licenseFile := data["license_file"].(string)
	vsanDedup := data["vsan_dedup"].(bool)

	vsanSpecBinding := &models.VSANSpec{
		DatastoreName: utils.ToStringPointer(datastoreName),
		HclFile:       hclFile,
		LicenseFile:   licenseFile,
		VSANDedup:     vsanDedup,
	}
	return vsanSpecBinding
}
