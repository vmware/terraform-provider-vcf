// © Broadcom. All Rights Reserved.
// The term “Broadcom” refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

package sddc

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/vmware/vcf-sdk-go/vcf"

	"github.com/vmware/terraform-provider-vcf/internal/network"
	validation_utils "github.com/vmware/terraform-provider-vcf/internal/validation"
)

func GetNsxSpecSchema() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeList,
		Optional: true,
		MaxItems: 1,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"vip": {
					Type:        schema.TypeString,
					Description: "Virtual IP address which would act as proxy/alias for NSX Managers",
					Required:    true,
				},
				"vip_fqdn": {
					Type:        schema.TypeString,
					Description: "FQDN for VIP so that common SSL certificates can be installed across all managers",
					Required:    true,
				},
				"root_nsx_manager_password": {
					Type:         schema.TypeString,
					Description:  "NSX Manager root password. Password should have 1) At least eight characters, 2) At least one lower-case letter, 3) At least one upper-case letter 4) At least one digit 5) At least one special character, 6) At least five different characters , 7) No dictionary words, 6) No palindromes",
					Required:     true,
					Sensitive:    true,
					ValidateFunc: validation_utils.ValidatePassword,
				},
				"ip_address_pool": {
					Type:        schema.TypeList,
					Description: "NSX IP address pool specification",
					Optional:    true,
					MaxItems:    1,
					Elem:        network.IpAddressPoolSchema(),
				},
				"nsx_admin_password": {
					Type:         schema.TypeString,
					Description:  "NSX admin password. The password must be at least 12 characters long. Must contain at-least 1 uppercase, 1 lowercase, 1 special character and 1 digit. In addition, a character cannot be repeated 3 or more times consecutively.",
					Optional:     true,
					Sensitive:    true,
					ValidateFunc: validation_utils.ValidatePassword,
				},
				"nsx_audit_password": {
					Type:         schema.TypeString,
					Description:  "NSX audit password. The password must be at least 12 characters long. Must contain at-least 1 uppercase, 1 lowercase, 1 special character and 1 digit. In addition, a character cannot be repeated 3 or more times consecutively.",
					Optional:     true,
					Sensitive:    true,
					ValidateFunc: validation_utils.ValidatePassword,
				},
				"license": {
					Type:        schema.TypeString,
					Description: "NSX Manager license",
					Optional:    true,
					Sensitive:   true,
				},
				"nsx_manager_size": {
					Type:         schema.TypeString,
					Description:  "NSX-T Manager size. One among: medium, large",
					Required:     true,
					ValidateFunc: validation.StringInSlice([]string{"medium", "large"}, true),
				},
				"nsx_manager":            getNsxManagerSpecSchema(),
				"overlay_transport_zone": getTransportZoneSchema(),
				"transport_vlan_id": {
					Type:         schema.TypeInt,
					Description:  "Transport VLAN ID",
					Required:     true,
					ValidateFunc: validation.IntBetween(0, 4095),
				},
			},
		},
	}
}

func getNsxManagerSpecSchema() *schema.Schema {
	return &schema.Schema{
		Type:        schema.TypeList,
		Description: "Parameters for NSX manager",
		Required:    true,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"hostname": {
					Type:        schema.TypeString,
					Description: "NSX Manager hostname. If just the short hostname is provided, then FQDN will be generated using the \"domain\" from dns configuration",
					Optional:    true,
				},
				"ip": {
					Type:         schema.TypeString,
					Description:  "NSX Manager IPv4 Address",
					Optional:     true,
					ValidateFunc: validation.IsIPAddress,
				},
			},
		},
	}
}

func getTransportZoneSchema() *schema.Schema {
	return &schema.Schema{
		Type:        schema.TypeList,
		Description: "NSX OverLay Transport zone",
		Optional:    true,
		MaxItems:    1,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"network_name": {
					Type:        schema.TypeString,
					Description: "Transport zone network name",
					Required:    true,
				},
				"zone_name": {
					Type:        schema.TypeString,
					Description: "Transport zone name",
					Required:    true,
				},
			},
		},
	}
}

func GetNsxSpecFromSchema(rawData []interface{}) *vcf.SddcNsxtSpec {
	if len(rawData) <= 0 {
		return nil
	}
	data := rawData[0].(map[string]interface{})
	nsxAdminPassword := data["nsx_admin_password"].(string)
	nsxAuditPassword := data["nsx_audit_password"].(string)
	nsxLicense := data["license"].(string)
	nsxManagerSize := data["nsx_manager_size"].(string)
	rootNsxManagerPassword := data["root_nsx_manager_password"].(string)
	transportVlanID := int32(data["transport_vlan_id"].(int))
	vip := data["vip"].(string)
	vipFqdn := data["vip_fqdn"].(string)

	nsxtSpecBinding := &vcf.SddcNsxtSpec{
		NsxtAdminPassword:       &nsxAdminPassword,
		NsxtAuditPassword:       &nsxAuditPassword,
		NsxtLicense:             &nsxLicense,
		NsxtManagerSize:         &nsxManagerSize,
		RootNsxtManagerPassword: &rootNsxManagerPassword,
		TransportVlanId:         &transportVlanID,
		Vip:                     &vip,
		VipFqdn:                 vipFqdn,
	}
	if nsxtManagersData := getNsxManagerSpecFromSchema(data["nsx_manager"].([]interface{})); len(nsxtManagersData) > 0 {
		nsxtSpecBinding.NsxtManagers = nsxtManagersData
	}
	if overLayTransportZoneData := getTransportZoneFromSchema(data["overlay_transport_zone"].([]interface{})); overLayTransportZoneData != nil {
		nsxtSpecBinding.OverLayTransportZone = overLayTransportZoneData
	}

	if ipAddressPoolRaw, ok := data["ip_address_pool"]; ok && !validation_utils.IsEmpty(ipAddressPoolRaw) {
		ipAddressPoolList := ipAddressPoolRaw.([]interface{})
		// Only one IP Address pool spec is allowed in the resource
		if ipAddressPoolSpec, err := network.GetIpAddressPoolSpecFromSchema(ipAddressPoolList[0].(map[string]interface{})); err == nil {
			nsxtSpecBinding.IpAddressPoolSpec = ipAddressPoolSpec
		}
	}
	return nsxtSpecBinding
}

func getNsxManagerSpecFromSchema(rawData []interface{}) []vcf.NsxtManagerSpec {
	var nsxtManagerSpecBindingsList []vcf.NsxtManagerSpec
	for _, nsxtManager := range rawData {
		data := nsxtManager.(map[string]interface{})
		hostname := data["hostname"].(string)
		ip := data["ip"].(string)

		nsxManagerSpec := vcf.NsxtManagerSpec{
			Hostname: &hostname,
			Ip:       &ip,
		}
		nsxtManagerSpecBindingsList = append(nsxtManagerSpecBindingsList, nsxManagerSpec)
	}
	return nsxtManagerSpecBindingsList
}

func getTransportZoneFromSchema(rawData []interface{}) *vcf.NsxtTransportZone {
	if len(rawData) <= 0 {
		return nil
	}
	data := rawData[0].(map[string]interface{})
	networkName := data["network_name"].(string)
	zoneName := data["zone_name"].(string)

	transportZoneBinding := &vcf.NsxtTransportZone{
		NetworkName: networkName,
		ZoneName:    zoneName,
	}
	return transportZoneBinding
}
