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
				"nsx_teaming": {
					Type:        schema.TypeList,
					Description: "NSX teaming policies for uplink profiles",
					Optional:    true,
					Elem: &schema.Resource{
						Schema: map[string]*schema.Schema{
							"policy": {
								Type:        schema.TypeString,
								Description: "Teaming policy (e.g., FAILOVER_ORDER, LOADBALANCE_SRCID)",
								Required:    true,
							},
							"active_uplinks": {
								Type:        schema.TypeList,
								Description: "List of active uplinks",
								Required:    true,
								Elem:        &schema.Schema{Type: schema.TypeString},
							},
							"standby_uplinks": {
								Type:        schema.TypeList,
								Description: "List of standby uplinks",
								Optional:    true,
								Elem:        &schema.Schema{Type: schema.TypeString},
							},
						},
					},
				},
				"nsxt_switch_config": {
					Type:        schema.TypeList,
					Description: "NSX-T switch configuration",
					Optional:    true,
					MaxItems:    1,
					Elem: &schema.Resource{
						Schema: map[string]*schema.Schema{
							"host_switch_operational_mode": {
								Type:        schema.TypeString,
								Description: "Host switch operational mode (e.g., STANDARD, ENS)",
								Optional:    true,
							},
							"ip_assignment_type": {
								Type:        schema.TypeString,
								Description: "IP assignment type for host switch",
								Optional:    true,
							},
							"transport_zones": {
								Type:        schema.TypeList,
								Description: "Transport zones for NSX switch",
								Required:    true,
								Elem: &schema.Resource{
									Schema: map[string]*schema.Schema{
										"name": {
											Type:        schema.TypeString,
											Description: "Transport zone name",
											Optional:    true,
										},
										"transport_type": {
											Type:        schema.TypeString,
											Description: "Transport type (e.g., OVERLAY, VLAN)",
											Required:    true,
										},
									},
								},
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

		if nsxTeamings, ok := dvsSpecRaw["nsx_teaming"].([]interface{}); ok {
			dvsSpec.NsxTeamings = convertNsxTeamings(nsxTeamings)
		}

		if nsxtSwitchConfig, ok := dvsSpecRaw["nsxt_switch_config"].([]interface{}); ok && len(nsxtSwitchConfig) > 0 {
			dvsSpec.NsxtSwitchConfig = convertNsxtSwitchConfig(nsxtSwitchConfig[0].(map[string]interface{}))
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

func convertNsxTeamings(rawData []interface{}) *[]installer.TeamingSpec {
	if len(rawData) == 0 {
		return nil
	}

	teamings := make([]installer.TeamingSpec, len(rawData))
	for i, teamingRaw := range rawData {
		teaming := teamingRaw.(map[string]interface{})
		teamingSpec := installer.TeamingSpec{
			Policy:        teaming["policy"].(string),
			ActiveUplinks: utils.ToStringSlice(teaming["active_uplinks"].([]interface{})),
		}

		if standbyUplinks, ok := teaming["standby_uplinks"].([]interface{}); ok {
			standbySlice := utils.ToStringSlice(standbyUplinks)
			teamingSpec.StandByUplinks = &standbySlice
		}

		teamings[i] = teamingSpec
	}
	return &teamings
}

func convertNsxtSwitchConfig(rawData map[string]interface{}) *installer.NsxtSwitchConfig {
	config := &installer.NsxtSwitchConfig{}
	hasConfig := false

	if mode, ok := rawData["host_switch_operational_mode"].(string); ok && mode != "" {
		config.HostSwitchOperationalMode = utils.ToStringPointer(mode)
		hasConfig = true
	}

	if ipType, ok := rawData["ip_assignment_type"].(string); ok && ipType != "" {
		config.IpAssignmentType = utils.ToStringPointer(ipType)
		hasConfig = true
	}

	if tzRaw, ok := rawData["transport_zones"].([]interface{}); ok {
		transportZones := make([]installer.TransportZone, len(tzRaw))
		for i, tzData := range tzRaw {
			tz := tzData.(map[string]interface{})
			transportZones[i] = installer.TransportZone{
				TransportType: tz["transport_type"].(string),
			}
			if name, ok := tz["name"].(string); ok && name != "" {
				transportZones[i].Name = utils.ToStringPointer(name)
			}
		}
		config.TransportZones = transportZones
		hasConfig = true
	}

	if !hasConfig {
		return nil
	}

	return config
}
