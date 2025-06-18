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

var trafficTypeValues = []string{"VSAN", "VMOTION", "VIRTUALMACHINE", "MANAGEMENT", "NFS", "VDP", "HBR", "FAULTTOLERANCE", "ISCSI"}

func GetDvsSchema() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeList,
		Required: true,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"dvs_name": {
					Type:        schema.TypeString,
					Description: "DVS Name",
					Required:    true,
				},
				"mtu": {
					Type:         schema.TypeInt,
					Description:  "DVS MTU (default value is 9000). In between 1500 and 9000",
					Optional:     true,
					Default:      9000,
					ValidateFunc: validation.IntBetween(1500, 9000),
				},
				"networks": {
					Type:        schema.TypeList,
					Description: "Types of networks in this portgroup. Possible values: VSAN, VMOTION, MANAGEMENT, VM_MANAGEMENT",
					Required:    true,
					Elem: &schema.Schema{
						Type: schema.TypeString,
					},
				},
				"nioc": getNiocSchema(),
				"vmnic_mapping": {
					Type:        schema.TypeList,
					Description: "Vmnic to uplink mappings",
					Required:    true,
					Elem: &schema.Resource{
						Schema: map[string]*schema.Schema{
							"vmnic": {
								Description: "Vmnic identifier",
								Type:        schema.TypeString,
								Required:    true,
							},
							"uplink": {
								Description: "Uplink identifier",
								Type:        schema.TypeString,
								Required:    true,
							},
						},
					},
				},
			},
		},
	}

}

func getNiocSchema() *schema.Schema {
	return &schema.Schema{
		Type:        schema.TypeList,
		Description: "List of NIOC specs for networks",
		Optional:    true,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"traffic_type": {
					Type:         schema.TypeString,
					Description:  "Traffic Type One among:VSAN, VMOTION, VIRTUALMACHINE, MANAGEMENT, NFS, VDP, HBR, FAULTTOLERANCE, ISCSI",
					Required:     true,
					ValidateFunc: validation.StringInSlice(trafficTypeValues, false),
				},
				"value": {
					Type:        schema.TypeString,
					Description: "NIOC Value. Example: LOW, NORMAL, HIGH",
					Required:    true,
				},
			},
		},
	}
}

func GetDvsSpecsFromSchema(rawData []interface{}) *[]installer.DvsSpec {
	var dvsSpecs []installer.DvsSpec
	for _, dvsSpecListEntry := range rawData {
		dvsSpecRaw := dvsSpecListEntry.(map[string]interface{})
		dvsName := utils.ToStringPointer(dvsSpecRaw["dvs_name"])
		mtu := int32(dvsSpecRaw["mtu"].(int))

		dvsSpec := installer.DvsSpec{
			DvsName: dvsName,
			Mtu:     &mtu,
		}
		if networksData, ok := dvsSpecRaw["networks"].([]interface{}); ok {
			networks := utils.ToStringSlice(networksData)
			dvsSpec.Networks = &networks
		}

		if vmnicMappings, ok := dvsSpecRaw["vmnic_mapping"].([]interface{}); ok {
			dvsSpec.VmnicsToUplinks = make([]installer.VmnicToUplink, len(vmnicMappings))
			for i, vmnicMapping := range vmnicMappings {
				dvsSpec.VmnicsToUplinks[i] = getVmnicToUplink(vmnicMapping.(map[string]interface{}))
			}
		}

		dvsSpecs = append(dvsSpecs, dvsSpec)
	}
	return &dvsSpecs
}

func getVmnicToUplink(rawData map[string]interface{}) installer.VmnicToUplink {
	return installer.VmnicToUplink{
		Uplink: rawData["uplink"].(string),
		Id:     rawData["vmnic"].(string),
	}
}
