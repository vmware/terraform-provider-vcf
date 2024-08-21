// © Broadcom. All Rights Reserved.
// The term “Broadcom” refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

package cluster

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	utils "github.com/vmware/terraform-provider-vcf/internal/resource_utils"
	"github.com/vmware/vcf-sdk-go/vcf"

	"github.com/vmware/terraform-provider-vcf/internal/network"
	validationutils "github.com/vmware/terraform-provider-vcf/internal/validation"
)

// HostSpecSchema this helper function extracts the Host
// schema, so that it's made available for both workload domain and cluster creation.
func HostSpecSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "ID of the ESXi host in the free pool",
			},
			"host_name": {
				Type:         schema.TypeString,
				Optional:     true,
				Description:  "Host name of the ESXi host",
				ValidateFunc: validation.NoZeroValues,
			},
			"availability_zone_name": {
				Type:         schema.TypeString,
				Optional:     true,
				Description:  "Availability Zone Name. This is required while performing a stretched cluster expand operation",
				ValidateFunc: validation.NoZeroValues,
			},
			"ip_address": {
				Type:         schema.TypeString,
				Optional:     true,
				Description:  "IPv4 address of the ESXi host",
				ValidateFunc: validationutils.ValidateIPv4AddressSchema,
			},
			"license_key": {
				Type:      schema.TypeString,
				Optional:  true,
				Sensitive: true,
				Description: "License key for an ESXi host in the free pool. This is required except in cases where the " +
					"ESXi host has already been licensed outside of the VMware Cloud Foundation system",
				ValidateFunc: validation.NoZeroValues,
			},
			"username": {
				Type:         schema.TypeString,
				Optional:     true,
				Description:  "Username to authenticate to the ESXi host",
				ValidateFunc: validation.NoZeroValues,
			},
			"password": {
				Type:         schema.TypeString,
				Optional:     true,
				Sensitive:    true,
				Description:  "Password to authenticate to the ESXi host",
				ValidateFunc: validation.NoZeroValues,
			},
			"serial_number": {
				Type:         schema.TypeString,
				Optional:     true,
				Description:  "Serial number of the ESXi host",
				ValidateFunc: validation.NoZeroValues,
			},
			"ssh_thumbprint": {
				Type:         schema.TypeString,
				Optional:     true,
				Sensitive:    true,
				Description:  "SSH thumbprint of the ESXi host",
				ValidateFunc: validation.NoZeroValues,
			},
			"vmnic": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Physical NIC configuration for the ESXi host",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Name of the physical NIC",
						},
						"mac_address": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "MAC address of the physical NIC",
						},
					},
				},
			},
		},
	}
}

func FlattenHost(host vcf.Host) *map[string]interface{} {
	result := make(map[string]interface{})
	result["id"] = host.Id
	result["host_name"] = host.Fqdn
	ipAddresses := *host.IpAddresses
	if len(ipAddresses) > 0 {
		result["ip_address"] = ipAddresses[0].IpAddress
	}
	if len(host.PhysicalNics) > 0 {
		var physicalNics []map[string]interface{}
		for _, nic := range host.PhysicalNics {
			nicMap := make(map[string]interface{})
			nicMap["name"] = nic.DeviceName
			nicMap["mac_address"] = nic.MacAddress
			physicalNics = append(physicalNics, nicMap)
		}
		result["vmnic"] = physicalNics
	}
	if len(host.PhysicalNics) > 0 {
		var physicalNics []map[string]interface{}
		for _, nic := range host.PhysicalNics {
			nicMap := make(map[string]interface{})
			nicMap["name"] = nic.DeviceName
			nicMap["mac_address"] = nic.MacAddress
			physicalNics = append(physicalNics, nicMap)
		}
		result["vmnic"] = physicalNics
	}

	return &result
}

func TryConvertToHostSpec(object map[string]interface{}) (*vcf.HostSpec, error) {
	result := &vcf.HostSpec{}
	if object == nil {
		return nil, fmt.Errorf("cannot convert to HostSpec, object is nil")
	}
	id := object["id"].(string)
	if len(id) == 0 {
		return nil, fmt.Errorf("cannot convert to HostSpec, id is required")
	}
	result.Id = id
	if hostName, ok := object["host_name"]; ok && !validationutils.IsEmpty(hostName) {
		result.HostName = utils.ToStringPointer(hostName)
	}
	if availabilityZoneName, ok := object["availability_zone_name"]; ok && !validationutils.IsEmpty(availabilityZoneName) {
		result.AzName = utils.ToStringPointer(availabilityZoneName)
	}
	if ipAddress, ok := object["ip_address"]; ok && !validationutils.IsEmpty(ipAddress) {
		result.IpAddress = utils.ToStringPointer(ipAddress)
	}
	if licenseKey, ok := object["license_key"]; ok && !validationutils.IsEmpty(licenseKey) {
		result.LicenseKey = utils.ToStringPointer(licenseKey)
	}
	if userName, ok := object["username"]; ok && !validationutils.IsEmpty(userName) {
		result.Username = utils.ToStringPointer(userName)
	}
	if password, ok := object["password"]; ok && !validationutils.IsEmpty(password) {
		result.Password = utils.ToStringPointer(password)
	}
	if serialNumber, ok := object["serial_number"]; ok && !validationutils.IsEmpty(serialNumber) {
		result.SerialNumber = utils.ToStringPointer(serialNumber)
	}
	if sshThumbprint, ok := object["ssh_thumbprint"]; ok && !validationutils.IsEmpty(sshThumbprint) {
		result.SshThumbprint = utils.ToStringPointer(sshThumbprint)
	}
	if vmNicsRaw, ok := object["vmnic"]; ok && !validationutils.IsEmpty(vmNicsRaw) {
		vmNicsList := vmNicsRaw.([]interface{})
		if len(vmNicsList) > 0 {
			result.HostNetworkSpec = &vcf.HostNetworkSpec{}
			vmNics := []vcf.VmNic{}
			for _, vmNicListEntry := range vmNicsList {
				vmNic, err := network.TryConvertToVmNic(vmNicListEntry.(map[string]interface{}))
				if err != nil {
					return nil, err
				}
				vmNics = append(vmNics, *vmNic)
			}
			result.HostNetworkSpec.VmNics = &vmNics
		} else {
			return nil, fmt.Errorf("cannot convert to HostSpec, vmnic list is empty")
		}
	}

	return result, nil
}
