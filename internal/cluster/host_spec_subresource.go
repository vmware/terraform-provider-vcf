// Copyright 2023 Broadcom. All Rights Reserved.
// SPDX-License-Identifier: MPL-2.0

package cluster

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/vmware/vcf-sdk-go/models"

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
				Description: "vmnic configuration for the ESXi host",
				Elem:        network.VMNicSchema(),
			},
		},
	}
}

func FlattenHostReference(host *models.HostReference) *map[string]interface{} {
	result := make(map[string]interface{})
	if host == nil {
		return &result
	}
	result["id"] = host.ID
	result["host_name"] = host.Fqdn
	result["ip_address"] = host.IPAddress
	result["availability_zone_name"] = host.AzName

	return &result
}

func FlattenHost(host *models.Host) *map[string]interface{} {
	result := make(map[string]interface{})
	if host == nil {
		return &result
	}
	result["id"] = host.ID
	result["host_name"] = host.Fqdn
	if len(host.IPAddresses) > 0 && host.IPAddresses[0] != nil {
		result["ip_address"] = host.IPAddresses[0].IPAddress
	}

	return &result
}

func TryConvertToHostSpec(object map[string]interface{}) (*models.HostSpec, error) {
	result := &models.HostSpec{}
	if object == nil {
		return nil, fmt.Errorf("cannot convert to HostSpec, object is nil")
	}
	id := object["id"].(string)
	if len(id) == 0 {
		return nil, fmt.Errorf("cannot convert to HostSpec, id is required")
	}
	result.ID = &id
	if hostName, ok := object["host_name"]; ok && !validationutils.IsEmpty(hostName) {
		result.HostName = hostName.(string)
	}
	if availabilityZoneName, ok := object["availability_zone_name"]; ok && !validationutils.IsEmpty(availabilityZoneName) {
		result.AzName = availabilityZoneName.(string)
	}
	if ipAddress, ok := object["ip_address"]; ok && !validationutils.IsEmpty(ipAddress) {
		result.IPAddress = ipAddress.(string)
	}
	if licenseKey, ok := object["license_key"]; ok && !validationutils.IsEmpty(licenseKey) {
		result.LicenseKey = licenseKey.(string)
	}
	if userName, ok := object["username"]; ok && !validationutils.IsEmpty(userName) {
		result.Username = userName.(string)
	}
	if password, ok := object["password"]; ok && !validationutils.IsEmpty(password) {
		result.Password = password.(string)
	}
	if serialNumber, ok := object["serial_number"]; ok && !validationutils.IsEmpty(serialNumber) {
		result.SerialNumber = serialNumber.(string)
	}
	if sshThumbprint, ok := object["ssh_thumbprint"]; ok && !validationutils.IsEmpty(sshThumbprint) {
		result.SSHThumbprint = sshThumbprint.(string)
	}
	if vmNicsRaw, ok := object["vmnic"]; ok && !validationutils.IsEmpty(vmNicsRaw) {
		vmNicsList := vmNicsRaw.([]interface{})
		if len(vmNicsList) > 0 {
			result.HostNetworkSpec = &models.HostNetworkSpec{}
			result.HostNetworkSpec.VMNics = []*models.VMNic{}
			for _, vmNicListEntry := range vmNicsList {
				vmNic, err := network.TryConvertToVmNic(vmNicListEntry.(map[string]interface{}))
				if err != nil {
					return nil, err
				}
				result.HostNetworkSpec.VMNics = append(result.HostNetworkSpec.VMNics, vmNic)
			}
		} else {
			return nil, fmt.Errorf("cannot convert to HostSpec, vmnic list is empty")
		}
	}

	return result, nil
}
