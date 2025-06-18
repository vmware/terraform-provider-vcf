// © Broadcom. All Rights Reserved.
// The term “Broadcom” refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

package sddc

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/vmware/terraform-provider-vcf/internal/network"
	validation_utils "github.com/vmware/terraform-provider-vcf/internal/validation"
	"github.com/vmware/vcf-sdk-go/installer"
)

func GetNsxSpecSchema() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeList,
		Optional: true,
		MaxItems: 1,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
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
				"nsx_manager_size": {
					Type:         schema.TypeString,
					Description:  "NSX-T Manager size. One among: medium, large",
					Required:     true,
					ValidateFunc: validation.StringInSlice([]string{"medium", "large"}, true),
				},
				"nsx_manager": getNsxManagerSpecSchema(),
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
			},
		},
	}
}

func GetNsxSpecFromSchema(rawData []interface{}) *installer.SddcNsxtSpec {
	if len(rawData) <= 0 {
		return nil
	}
	data := rawData[0].(map[string]interface{})
	nsxAdminPassword := data["nsx_admin_password"].(string)
	nsxAuditPassword := data["nsx_audit_password"].(string)
	nsxManagerSize := data["nsx_manager_size"].(string)
	rootNsxManagerPassword := data["root_nsx_manager_password"].(string)
	transportVlanID := int32(data["transport_vlan_id"].(int))
	vipFqdn := data["vip_fqdn"].(string)

	nsxtSpecBinding := &installer.SddcNsxtSpec{
		NsxtAdminPassword:       &nsxAdminPassword,
		NsxtAuditPassword:       &nsxAuditPassword,
		NsxtManagerSize:         &nsxManagerSize,
		RootNsxtManagerPassword: &rootNsxManagerPassword,
		TransportVlanId:         &transportVlanID,
		VipFqdn:                 vipFqdn,
	}
	if nsxtManagersData := getNsxManagerSpecFromSchema(data["nsx_manager"].([]interface{})); len(nsxtManagersData) > 0 {
		nsxtSpecBinding.NsxtManagers = nsxtManagersData
	}

	if ipAddressPoolRaw, ok := data["ip_address_pool"]; ok && !validation_utils.IsEmpty(ipAddressPoolRaw) {
		ipAddressPoolList := ipAddressPoolRaw.([]interface{})
		// Only one IP Address pool spec is allowed in the resource
		if ipAddressPoolSpec, err := network.GetInstallerIpAddressPoolSpecFromSchema(ipAddressPoolList[0].(map[string]interface{})); err == nil {
			nsxtSpecBinding.IpAddressPoolSpec = ipAddressPoolSpec
		}
	}
	return nsxtSpecBinding
}

func getNsxManagerSpecFromSchema(rawData []interface{}) []installer.NsxtManagerSpec {
	var nsxtManagerSpecBindingsList []installer.NsxtManagerSpec
	for _, nsxtManager := range rawData {
		data := nsxtManager.(map[string]interface{})
		hostname := data["hostname"].(string)

		nsxManagerSpec := installer.NsxtManagerSpec{
			Hostname: &hostname,
		}
		nsxtManagerSpecBindingsList = append(nsxtManagerSpecBindingsList, nsxManagerSpec)
	}
	return nsxtManagerSpecBindingsList
}
