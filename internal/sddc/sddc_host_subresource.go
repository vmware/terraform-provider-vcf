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

var portGroup = []string{"VSAN", "VMOTION", "PUBLIC", "MANAGEMENT", "NSX_VTEP", "HOSTMANAGEMENT", "CLOUD_VENDOR_API", "REPLICATION"}

func GetSddcHostSchema() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeList,
		Required: true,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"association": {
					Type:        schema.TypeString,
					Description: "Host Association: Location/Datacenter",
					Required:    true,
				},
				"credentials": getCredentialsSchema(),
				"hostname": {
					Type:         schema.TypeString,
					Description:  "Host FQDN Example: esx-1, length from 3 to 63",
					Required:     true,
					ValidateFunc: validation.StringLenBetween(3, 63),
				},
				"ip_address_private": getIPAllocationSchema(),
				"key": {
					Type:        schema.TypeString,
					Description: "Host key",
					Optional:    true,
				},
				"server_id": {
					Type:        schema.TypeString,
					Description: "Host server ID",
					Optional:    true,
				},
				"ssh_thumbprint": {
					Type:        schema.TypeString,
					Description: "Host SSH thumbprint (RSA SHA256)",
					Optional:    true,
				},
				"vmknic_specs": getHostVmknicSchema(),
				"vswitch": {
					Type:     schema.TypeString,
					Required: true,
				},
			},
		},
	}
}

func getIPAllocationSchema() *schema.Schema {
	return &schema.Schema{
		Type:        schema.TypeList,
		Description: "Host Private Management IP",
		MaxItems:    1,
		Required:    true,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"cidr": {
					Type:        schema.TypeString,
					Description: "Classless Inter-Domain Routing (CIDR), Example: 172.0.0.0/24",
					Optional:    true,
				},
				"gateway": {
					Type:         schema.TypeString,
					Description:  "Gateway",
					Required:     true,
					ValidateFunc: validation.IsIPAddress,
				},
				"ip_address": {
					Type:         schema.TypeString,
					Description:  "IP Address",
					Required:     true,
					ValidateFunc: validation.IsIPAddress,
				},
				"subnet": {
					Type:         schema.TypeString,
					Description:  "Subnet",
					Optional:     true,
					ValidateFunc: validation.IsIPAddress,
				},
			},
		},
	}
}

func getHostVmknicSchema() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeList,
		Optional: true,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"ip_address": {
					Type:         schema.TypeString,
					Description:  "Vmknic IP address",
					Optional:     true,
					ValidateFunc: validation.IsIPAddress,
				},
				"mac_address": {
					Type:         schema.TypeString,
					Description:  "Vmknic mac address",
					Optional:     true,
					ValidateFunc: validation.IsMACAddress,
				},
				"port_group": {
					Type:        schema.TypeList,
					Description: "Portgroup type. One among: VSAN, VMOTION, PUBLIC, MANAGEMENT, NSX_VTEP, HOSTMANAGEMENT, CLOUD_VENDOR_API, REPLICATION",
					Required:    true,
					Elem: &schema.Schema{
						Type:         schema.TypeString,
						ValidateFunc: validation.StringInSlice(portGroup, false),
					},
				},
			},
		},
	}
}

func GetSddcHostSpecsFromSchema(rawData []interface{}) []*models.SDDCHostSpec {
	var hostSpecs []*models.SDDCHostSpec
	for _, rawListEntity := range rawData {
		hostSpecRaw := rawListEntity.(map[string]interface{})
		association := utils.ToStringPointer(hostSpecRaw["association"])
		hostname := utils.ToStringPointer(hostSpecRaw["hostname"])
		key := hostSpecRaw["key"].(string)
		serverID := hostSpecRaw["server_id"].(string)
		sshThumbprint := hostSpecRaw["ssh_thumbprint"].(string)
		vswitch := utils.ToStringPointer(hostSpecRaw["vswitch"])

		hostSpec := &models.SDDCHostSpec{
			Association:   association,
			Hostname:      hostname,
			Key:           key,
			ServerID:      serverID,
			SSHThumbprint: sshThumbprint,
			VSwitch:       vswitch,
		}
		if credentialsData := getCredentialsFromSchema(hostSpecRaw["credentials"].([]interface{})); credentialsData != nil {
			hostSpec.Credentials = credentialsData
		}
		if ipAllocation := getIPAllocationBindingFromSchema(hostSpecRaw["ip_address_private"].([]interface{})); ipAllocation != nil {
			hostSpec.IPAddressPrivate = ipAllocation
		}
		if vmknicSpecs := getHostVmknicSpecsFromSchema(hostSpecRaw["vmknic_specs"].([]interface{})); len(vmknicSpecs) != 0 {
			hostSpec.VmknicSpecs = vmknicSpecs
		}
		hostSpecs = append(hostSpecs, hostSpec)
	}
	return hostSpecs
}

func getIPAllocationBindingFromSchema(rawData []interface{}) *models.IPAllocation {
	if len(rawData) <= 0 {
		return nil
	}
	data := rawData[0].(map[string]interface{})
	cidr := data["cidr"].(string)
	gateway := data["gateway"].(string)
	ipAddress := utils.ToStringPointer(data["ip_address"])
	subnet := data["subnet"].(string)

	ipAllocationBinding := &models.IPAllocation{
		Cidr:      cidr,
		Gateway:   gateway,
		IPAddress: ipAddress,
		Subnet:    subnet,
	}
	return ipAllocationBinding
}

func getHostVmknicSpecsFromSchema(rawData []interface{}) []*models.HostVmknicSpec {
	var hostVmknicSpecs []*models.HostVmknicSpec
	for _, rawListEntity := range rawData {
		hostVmknicSpecRaw := rawListEntity.(map[string]interface{})
		ipAddress := hostVmknicSpecRaw["ip_address"].(string)
		macAddress := hostVmknicSpecRaw["mac_address"].(string)

		hostVmknicSpec := &models.HostVmknicSpec{
			IPAddress:  ipAddress,
			MacAddress: macAddress,
		}
		if portGroups, ok := hostVmknicSpecRaw["port_group"].([]interface{}); ok {
			hostVmknicSpec.Portgroup = utils.ToStringPointer(portGroups)
		}

		hostVmknicSpecs = append(hostVmknicSpecs, hostVmknicSpec)
	}
	return hostVmknicSpecs
}
