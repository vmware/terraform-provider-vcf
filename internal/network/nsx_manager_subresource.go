// © Broadcom. All Rights Reserved.
// The term “Broadcom” refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

package network

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/vmware/terraform-provider-vcf/internal/resource_utils"
	"github.com/vmware/vcf-sdk-go/vcf"

	validationutils "github.com/vmware/terraform-provider-vcf/internal/validation"
)

// NsxManagerNodeSchema this helper function extracts the NSX Manager Node schema, which contains
// the parameters required to install and configure NSX Manager in a workload domain.
func NsxManagerNodeSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				Description:  "Name of the NSX Manager appliance, e.g., sfo-w01-nsx01 ",
				ValidateFunc: validation.NoZeroValues,
			},
			"ip_address": {
				Type:         schema.TypeString,
				Required:     true,
				Description:  "IPv4 address of the NSX Manager appliance",
				ValidateFunc: validationutils.ValidateIPv4AddressSchema,
			},
			"fqdn": {
				Type:         schema.TypeString,
				Required:     true,
				Description:  "Fully qualified domain name of the NSX Manager appliance, e.g., sfo-w01-nsx01a.sfo.rainpole.io",
				ValidateFunc: validation.NoZeroValues,
			},
			"subnet_mask": {
				Type:         schema.TypeString,
				Required:     true,
				Description:  "IPv4 subnet mask for the NSX Manager appliance",
				ValidateFunc: validationutils.ValidateIPv4AddressSchema,
			},
			"gateway": {
				Type:         schema.TypeString,
				Required:     true,
				Description:  "IPv4 gateway the NSX Manager appliance",
				ValidateFunc: validationutils.ValidateIPv4AddressSchema,
			},
		},
	}
}

func TryConvertToNsxManagerNodeSpec(object map[string]interface{}) (vcf.NsxManagerSpec, error) {
	result := vcf.NsxManagerSpec{}
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
	result.NetworkDetailsSpec = vcf.NetworkDetailsSpec{
		IpAddress: &ipAddress,
	}
	if fqdn, ok := object["fqdn"]; ok && !validationutils.IsEmpty(fqdn) {
		result.NetworkDetailsSpec.DnsName = fqdn.(string)
	}
	if subnetMask, ok := object["subnet_mask"]; ok && !validationutils.IsEmpty(subnetMask) {
		result.NetworkDetailsSpec.SubnetMask = resource_utils.ToStringPointer(subnetMask)
	}
	if gateway, ok := object["gateway"]; ok && !validationutils.IsEmpty(gateway) {
		result.NetworkDetailsSpec.Gateway = resource_utils.ToStringPointer(gateway)
	}
	return result, nil
}
