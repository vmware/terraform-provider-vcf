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
				"is_used_by_nsxt": {
					Type:        schema.TypeBool,
					Description: "Flag indicating whether the DVS is used by NSX",
					Optional:    true,
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
					Description: "Types of networks in this portgroup. Possible values: VSAN, VMOTION, PUBLIC, MANAGEMENT, NSX_VTEP, HOST_MANAGEMENT,\n\tCLOUD_VENDOR_API, REPLICATION, DATACENTER_NETWORK, NSX_VXLAN, NON_ROUTABLE, CLOUD_VENDOR_API,\n\tOOB, CROSS_VPC, UPLINK01, UPLINK02, STORAGE, UDLR, DLR, X_REGION, REGION_SPECIFIC,\n\tREMOTE_REGION_SPECIFIC, COMPUTE, MANAGEMENT_VM, NSXT_EDGE_TEP, NSXT_HOST_OVERLAY",
					Required:    true,
					Elem: &schema.Schema{
						Type:         schema.TypeString,
						ValidateFunc: validation.StringInSlice(NetworkSpecNetworkType, false),
					},
				},
				"nioc": getNiocSchema(),
				"vmnics": {
					Type:        schema.TypeList,
					Description: "Vmnics to be attached to the DVS",
					Required:    true,
					Elem:        &schema.Schema{Type: schema.TypeString},
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
					Description: "NIOC Value",
					Required:    true,
				},
			},
		},
	}
}

func GetDvsSpecsFromSchema(rawData []interface{}) []*models.DvsSpec {
	var dvsSpecs []*models.DvsSpec
	for _, dvsSpecListEntry := range rawData {
		dvsSpecRaw := dvsSpecListEntry.(map[string]interface{})
		dvsName := utils.ToStringPointer(dvsSpecRaw["dvs_name"])
		isUsedByNsxt := dvsSpecRaw["is_used_by_nsxt"].(bool)
		mtu := dvsSpecRaw["mtu"].(int32)

		dvsSpec := &models.DvsSpec{
			DvsName:      dvsName,
			IsUsedByNSXT: isUsedByNsxt,
			Mtu:          mtu,
		}
		if networksData, ok := dvsSpecRaw["networks"].([]interface{}); ok {
			dvsSpec.Networks = utils.ToStringSlice(networksData)
		}
		if niocSpecsData := getNiocSpecsFromSchema(dvsSpecRaw["nioc"].([]interface{})); len(niocSpecsData) > 0 {
			dvsSpec.NiocSpecs = niocSpecsData
		}
		if vmnicsData, ok := dvsSpecRaw["vmnics"].([]interface{}); ok {
			dvsSpec.Vmnics = utils.ToStringSlice(vmnicsData)
		}
		dvsSpecs = append(dvsSpecs, dvsSpec)
	}
	return dvsSpecs
}

func getNiocSpecsFromSchema(rawData []interface{}) []*models.NiocSpec {
	var niocSpecBindingsList []*models.NiocSpec
	for _, niocSpecListEntry := range rawData {
		niocSpecRaw := niocSpecListEntry.(map[string]interface{})
		trafficType := utils.ToStringPointer(niocSpecRaw["traffic_type"])
		value := utils.ToStringPointer(niocSpecRaw["value"])

		niocSpecsBinding := &models.NiocSpec{
			TrafficType: trafficType,
			Value:       value,
		}
		niocSpecBindingsList = append(niocSpecBindingsList, niocSpecsBinding)
	}
	return niocSpecBindingsList
}
