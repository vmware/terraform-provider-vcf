/*
 *  Copyright 2023 VMware, Inc.
 *    SPDX-License-Identifier: MPL-2.0
 */

package network

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	validation_utils "github.com/vmware/terraform-provider-vcf/internal/validation"
	"github.com/vmware/vcf-sdk-go/models"
)

// NsxManagerNodeSchema this helper function extracts the NSX Manager Node schema, which contains
// the parameters required to install and configure NSX Manager in a workload domain.
func NsxManagerNodeSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				Description:  "Name of the NSX Manager virtual machine",
				ValidateFunc: validation.NoZeroValues,
			},
			"ip_address": {
				Type:         schema.TypeString,
				Required:     true,
				Description:  "IPv4 address of the virtual machine",
				ValidateFunc: validation_utils.ValidateIPv4AddressSchema,
			},
			"dns_name": {
				Type:         schema.TypeString,
				Required:     true,
				Description:  "DNS name of the virtual machine, e.g., nsx-1.domain1.rainpole.io",
				ValidateFunc: validation.NoZeroValues,
			},
			"subnet_mask": {
				Type:         schema.TypeString,
				Required:     true,
				Description:  "Subnet mask",
				ValidateFunc: validation_utils.ValidateIPv4AddressSchema,
			},
			"gateway": {
				Type:         schema.TypeString,
				Required:     true,
				Description:  "IPv4 gateway the VM can use to connect to the outside world",
				ValidateFunc: validation_utils.ValidateIPv4AddressSchema,
			},
		},
	}
}

func TryConvertToNsxManagerNodeSpec(object map[string]interface{}) (models.NsxManagerSpec, error) {
	result := models.NsxManagerSpec{}
	if object == nil {
		return result, fmt.Errorf("cannot convert to NsxManagerSpec, object is nil")
	}
	name := object["name"].(string)
	if len(name) == 0 {
		return result, fmt.Errorf("cannot convert to NsxManagerSpec, name is required")
	}
	ipAddress := object["ip_address"].(string)
	if len(ipAddress) == 0 {
		return result, fmt.Errorf("cannot convert to NsxManagerSpec, ip_address is required")
	}
	result.Name = &name
	result.NetworkDetailsSpec = &models.NetworkDetailsSpec{
		IPAddress: &ipAddress,
	}
	if dnsName, ok := object["dns_name"]; ok && !validation_utils.IsEmpty(dnsName) {
		result.NetworkDetailsSpec.DNSName = dnsName.(string)
	}
	if subnetMask, ok := object["subnet_mask"]; ok && !validation_utils.IsEmpty(subnetMask) {
		result.NetworkDetailsSpec.SubnetMask = subnetMask.(string)
	}
	if gateway, ok := object["gateway"]; ok && !validation_utils.IsEmpty(gateway) {
		result.NetworkDetailsSpec.Gateway = gateway.(string)
	}
	return result, nil
}
