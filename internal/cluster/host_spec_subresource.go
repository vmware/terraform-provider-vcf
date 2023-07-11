/*
 *  Copyright 2023 VMware, Inc.
 *    SPDX-License-Identifier: MPL-2.0
 */

package cluster

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/vmware/terraform-provider-vcf/internal/network"
	validation_utils "github.com/vmware/terraform-provider-vcf/internal/validation"
	"github.com/vmware/vcf-sdk-go/models"
)

// HostSpecSchema this helper function extracts the Host
// schema, so that it's made available for both Domain and Cluster creation.
func HostSpecSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "ID of a vSphere host in the free pool",
			},
			"host_name": {
				Type:         schema.TypeString,
				Optional:     true,
				Description:  "Host name of the vSphere host",
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
				Description:  "IP address of the vSphere host",
				ValidateFunc: validation_utils.ValidateIPv4AddressSchema,
			},
			"license_key": {
				Type:     schema.TypeString,
				Optional: true,
				Description: "License key of a vSphere host in the free pool. This is required except in cases where the " +
					"ESXi host has already been licensed outside of the VMware Cloud Foundation system",
				ValidateFunc: validation.NoZeroValues,
			},
			"username": {
				Type:         schema.TypeString,
				Optional:     true,
				Description:  "Username of the vSphere host",
				ValidateFunc: validation.NoZeroValues,
			},
			"password": {
				Type:         schema.TypeString,
				Optional:     true,
				Sensitive:    true,
				Description:  "SSH password of the vSphere host",
				ValidateFunc: validation.NoZeroValues,
			},
			"serial_number": {
				Type:         schema.TypeString,
				Optional:     true,
				Description:  "Serial Number of the vSphere hosts",
				ValidateFunc: validation.NoZeroValues,
			},
			"ssh_thumbprint": {
				Type:         schema.TypeString,
				Optional:     true,
				Sensitive:    true,
				Description:  "SSH thumbprint (fingerprint) of the vSphere host. Note: This field will be mandatory in future releases.",
				ValidateFunc: validation.NoZeroValues,
			},
			"vmnic": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Contains vmnic configurations for vSphere host",
				Elem:        network.VMNicSchema(),
			},
		},
	}
}

func TryConvertToHostSpec(object map[string]interface{}) (*models.HostSpec, error) {
	result := &models.HostSpec{}
	if object == nil {
		return nil, fmt.Errorf("cannot conver to HostSpec, object is nil")
	}
	id := object["id"].(string)
	if len(id) == 0 {
		return nil, fmt.Errorf("cannot conver to HostSpec, id is required")
	}
	result.ID = &id
	if hostName, ok := object["host_name"]; ok && !validation_utils.IsEmpty(hostName) {
		result.HostName = hostName.(string)
	}
	if availabilityZoneName, ok := object["availability_zone_name"]; ok && !validation_utils.IsEmpty(availabilityZoneName) {
		result.AzName = availabilityZoneName.(string)
	}
	if ipAddress, ok := object["ip_address"]; ok && !validation_utils.IsEmpty(ipAddress) {
		result.IPAddress = ipAddress.(string)
	}
	if licenseKey, ok := object["license_key"]; ok && !validation_utils.IsEmpty(licenseKey) {
		result.LicenseKey = licenseKey.(string)
	}
	if userName, ok := object["username"]; ok && !validation_utils.IsEmpty(userName) {
		result.Username = userName.(string)
	}
	if password, ok := object["password"]; ok && !validation_utils.IsEmpty(password) {
		result.Password = password.(string)
	}
	if serialNumber, ok := object["serial_number"]; ok && !validation_utils.IsEmpty(serialNumber) {
		result.SerialNumber = serialNumber.(string)
	}
	if sshThumbprint, ok := object["ssh_thumbprint"]; ok && !validation_utils.IsEmpty(sshThumbprint) {
		result.SSHThumbprint = sshThumbprint.(string)
	}
	if vmNicsRaw, ok := object["vmnic"]; ok && !validation_utils.IsEmpty(vmNicsRaw) {
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
			return nil, fmt.Errorf("cannot convert to ClusterSpec, hosts list is empty")
		}
	}

	return result, nil
}
