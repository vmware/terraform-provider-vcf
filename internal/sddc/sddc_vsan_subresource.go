// © Broadcom. All Rights Reserved.
// The term “Broadcom” refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

package sddc

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/vmware/vcf-sdk-go/installer"
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

func GetVsanSpecFromSchema(rawData []interface{}) *installer.VsanSpec {
	if len(rawData) <= 0 {
		return nil
	}
	data := rawData[0].(map[string]interface{})
	datastoreName := data["datastore_name"].(string)
	vsanDedup := data["vsan_dedup"].(bool)
	esaEnabled := data["esa_enabled"].(bool)

	vsanSpecBinding := &installer.VsanSpec{
		DatastoreName: &datastoreName,
		VsanDedup:     &vsanDedup,
		EsaConfig:     &installer.VsanEsaConfig{Enabled: &esaEnabled},
	}

	return vsanSpecBinding
}
