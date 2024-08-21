// © Broadcom. All Rights Reserved.
// The term “Broadcom” refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

package sddc

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/vmware/vcf-sdk-go/models"

	utils "github.com/vmware/terraform-provider-vcf/internal/resource_utils"
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
					Description: "A path (URL or local path) to an HCL file that will be uploaded to vCenter prior to configuring vSAN",
					Optional:    true,
				},
				"license": {
					Type:        schema.TypeString,
					Description: "VSAN License",
					Optional:    true,
				},
				"vsan_dedup": {
					Type:        schema.TypeBool,
					Description: "VSAN feature Deduplication and Compression flag, one flag for both features",
					Optional:    true,
				},
				"esa_enabled": {
					Type:        schema.TypeBool,
					Optional:    true,
					Description: "Enable vSAN ESA",
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
	license := data["license"].(string)
	vsanDedup := data["vsan_dedup"].(bool)
	esaEnabled := data["esa_enabled"].(bool)

	vsanSpecBinding := &models.VSANSpec{
		DatastoreName: utils.ToStringPointer(datastoreName),
		HclFile:       hclFile,
		LicenseFile:   license,
		VSANDedup:     vsanDedup,
		EsaConfig:     &models.VSANEsaConfig{Enabled: &esaEnabled},
	}

	return vsanSpecBinding
}
