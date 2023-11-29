// Copyright 2023 Broadcom. All Rights Reserved.
// SPDX-License-Identifier: MPL-2.0

package sddc

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	utils "github.com/vmware/terraform-provider-vcf/internal/resource_utils"
	"github.com/vmware/vcf-sdk-go/models"
)

var teamingPolicies = []string{"loadbalance_loadbased", "loadbalance_ip", "loadbalance_srcmac", "loadbalance_srcid", "failover_explicit"}

func GetNetworkSpecsSchema() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeList,
		Required: true,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"network_type": {
					Type:        schema.TypeString,
					Description: "Network Type. One among: VSAN, VMOTION, MANAGEMENT, VM_MANAGEMENT or any custom network type",
					Required:    true,
				},
				"vlan_id": {
					Type:         schema.TypeString,
					Description:  "VLAN Id",
					Required:     true,
					ValidateFunc: validation.StringLenBetween(1, 4),
				},
				"active_up_links": {
					Type:        schema.TypeList,
					Description: "Active Uplinks for teaming policy, specify uplink1 for failover_explicit VSAN Teaming Policy",
					Optional:    true,
					Elem: &schema.Schema{
						Type: schema.TypeString,
					},
				},
				"standby_uplinks": {
					Type:        schema.TypeList,
					Description: "Standby Uplinks for teaming policy, specify uplink2 for failover_explicit VSAN Teaming Policy",
					Optional:    true,
					Elem: &schema.Schema{
						Type: schema.TypeString,
					},
				},
				"exclude_ip_address_ranges": {
					Type:        schema.TypeList,
					Description: "IP Address ranges to be excluded",
					Optional:    true,
					Elem: &schema.Schema{
						Type: schema.TypeString,
					},
				},
				"exclude_ip_addresses": {
					Type:        schema.TypeList,
					Description: "IP Addresses to be excluded",
					Optional:    true,
					Elem: &schema.Schema{
						Type:         schema.TypeString,
						ValidateFunc: validation.IsIPAddress,
					},
				},
				"gateway": {
					Type:         schema.TypeString,
					Optional:     true,
					ValidateFunc: validation.IsIPAddress,
				},
				"include_ip_address": {
					Type:     schema.TypeList,
					Optional: true,
					Elem: &schema.Schema{
						Type:         schema.TypeString,
						ValidateFunc: validation.IsIPAddress,
					},
				},
				"include_ip_address_ranges": getIncludeIPAddressRangesSchema(),
				"mtu": {
					Type:         schema.TypeString,
					Description:  "MTU size",
					Required:     true,
					ValidateFunc: validation.StringLenBetween(4, 4),
				},

				"port_group_key": {
					Type:        schema.TypeString,
					Description: "Portgroup key name. When adding a cluster with a new DVS, this value must be provided. When adding a cluster to an existing DVS, this value must not be provided.",
					Optional:    true,
				},
				"subnet": {
					Type:         schema.TypeString,
					Optional:     true,
					ValidateFunc: validation.StringLenBetween(7, 18),
				},
				"subnet_mask": {
					Type:         schema.TypeString,
					Optional:     true,
					ValidateFunc: validation.StringLenBetween(7, 15),
				},
				"teaming_policy": {
					Type:         schema.TypeString,
					Optional:     true,
					Description:  "Teaming Policy for VSAN and VMOTION network types, Default is loadbalance_loadbased. One among: loadbalance_ip, loadbalance_srcmac, loadbalance_srcid, failover_explicit, loadbalance_loadbased",
					Default:      teamingPolicies[0],
					ValidateFunc: validation.StringInSlice(teamingPolicies, false),
				},
			},
		},
	}
}

func getIncludeIPAddressRangesSchema() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeList,
		Optional: true,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"end_ip_address": {
					Type:         schema.TypeString,
					Description:  "End IPv4 Address",
					Required:     true,
					ValidateFunc: validation.IsIPAddress,
				},
				"start_ip_address": {
					Type:         schema.TypeString,
					Description:  "Start IPv4 Address",
					Required:     true,
					ValidateFunc: validation.IsIPAddress,
				},
			},
		},
	}
}

func GetNetworkSpecsBindingFromSchema(rawData []interface{}) []*models.SDDCNetworkSpec {
	var networkSpecsBindingsList []*models.SDDCNetworkSpec
	for _, networkSpec := range rawData {
		data := networkSpec.(map[string]interface{})
		subnet := data["subnet"].(string)
		vlanID := data["vlan_id"].(string)
		mtu := data["mtu"].(string)
		portGroupKey := data["port_group_key"].(string)
		networkType := data["network_type"].(string)
		gateway := data["gateway"].(string)
		subnetMask := data["subnet_mask"].(string)
		teamingPolicy := data["teaming_policy"].(string)

		networkSpecsBinding := &models.SDDCNetworkSpec{
			Gateway:       gateway,
			Mtu:           mtu,
			NetworkType:   utils.ToStringPointer(networkType),
			PortGroupKey:  portGroupKey,
			Subnet:        subnet,
			SubnetMask:    subnetMask,
			TeamingPolicy: teamingPolicy,
			VlanID:        utils.ToStringPointer(vlanID),
		}
		if activeUpLinksData, ok := data["active_up_links"].([]interface{}); ok {
			networkSpecsBinding.ActiveUplinks = utils.ToStringSlice(activeUpLinksData)
		}
		if excludeIPAddressRangesData, ok := data["exclude_ip_address_ranges"].([]interface{}); ok {
			networkSpecsBinding.ExcludeIPAddressRanges = utils.ToStringSlice(excludeIPAddressRangesData)
		}
		if excludeIPAddressesData, ok := data["exclude_ip_addresses"].([]interface{}); ok {
			networkSpecsBinding.ExcludeIpaddresses = utils.ToStringSlice(excludeIPAddressesData)
		}
		if includeIPAddressData, ok := data["include_ip_address"].([]interface{}); ok {
			networkSpecsBinding.IncludeIPAddress = utils.ToStringSlice(includeIPAddressData)
		}
		if includeIPAddressRangesData := getIncludeIPAddressRangesBindingFromSchema(data["include_ip_address_ranges"].([]interface{})); len(includeIPAddressRangesData) > 0 {
			networkSpecsBinding.IncludeIPAddressRanges = includeIPAddressRangesData
		}
		if standbyUplinksData, ok := data["standby_uplinks"].([]interface{}); ok {
			networkSpecsBinding.StandbyUplinks = utils.ToStringSlice(standbyUplinksData)
		}
		networkSpecsBindingsList = append(networkSpecsBindingsList, networkSpecsBinding)
	}
	return networkSpecsBindingsList
}

func getIncludeIPAddressRangesBindingFromSchema(rawData []interface{}) []*models.IPRange {
	var ipAddressRangesBindindsList []*models.IPRange
	for _, ipAddressRange := range rawData {
		data := ipAddressRange.(map[string]interface{})
		startIPAddress := data["start_ip_address"].(string)
		endIPAddress := data["end_ip_address"].(string)

		ipAddressRangesBinding := &models.IPRange{
			StartIPAddress: utils.ToStringPointer(startIPAddress),
			EndIPAddress:   utils.ToStringPointer(endIPAddress),
		}
		ipAddressRangesBindindsList = append(ipAddressRangesBindindsList, ipAddressRangesBinding)
	}
	return ipAddressRangesBindindsList
}
