// © Broadcom. All Rights Reserved.
// The term “Broadcom” refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

package network

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	utils "github.com/vmware/terraform-provider-vcf/internal/resource_utils"
	"github.com/vmware/vcf-sdk-go/installer"
	"github.com/vmware/vcf-sdk-go/vcf"

	validationutils "github.com/vmware/terraform-provider-vcf/internal/validation"
)

// IpAddressPoolSchema this helper function extracts the IpAddressPoolSpec schema, which
// contains the parameters required to create or reuse an IP address pool.
func IpAddressPoolSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"name": {
				Type: schema.TypeString,
				Description: "Providing only name of existing IP Address Pool reuses it, " +
					"while providing a new name with subnets creates a new one",
				Required:     true,
				ValidateFunc: validation.NoZeroValues,
			},
			"description": {
				Type:        schema.TypeString,
				Description: "Description of the IP address pool",
				Optional:    true,
			},
			"ignore_unavailable_nsx_cluster": {
				Type:        schema.TypeBool,
				Description: "Ignore unavailable NSX cluster(s) during IP pool spec validation",
				Optional:    true,
			},
			"subnet": {
				Type:        schema.TypeList,
				Description: "List of IP address pool subnet specifications",
				Optional:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"cidr": {
							Type:         schema.TypeString,
							Description:  "The subnet representation, contains the network address and the prefix length",
							Required:     true,
							ValidateFunc: validation.IsCIDR,
						},
						"gateway": {
							Type:         schema.TypeString,
							Description:  "The default gateway address of the network",
							Required:     true,
							ValidateFunc: validationutils.ValidateIPv4AddressSchema,
						},
						"ip_address_pool_range": {
							Type:        schema.TypeList,
							Description: "List of the IP allocation ranges. At least 1 IP address range has to be specified",
							Optional:    true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"start": {
										Type:         schema.TypeString,
										Description:  "The first IP Address of the IP Address Range",
										Required:     true,
										ValidateFunc: validationutils.ValidateIPv4AddressSchema,
									},
									"end": {
										Type:         schema.TypeString,
										Description:  "The last IP Address of the IP Address Range",
										Required:     true,
										ValidateFunc: validationutils.ValidateIPv4AddressSchema,
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

func GetIpAddressPoolSpecFromSchema(object map[string]interface{}) (*vcf.IpAddressPoolSpec, error) {
	result := &vcf.IpAddressPoolSpec{}
	if object == nil {
		return nil, fmt.Errorf("cannot convert to IpAddressPoolSpec, object is nil")
	}
	name := object["name"].(string)
	if len(name) == 0 {
		return nil, fmt.Errorf("cannot convert to IpAddressPoolSpec, name is required")
	}
	result.Name = name
	if description, ok := object["description"]; ok && !validationutils.IsEmpty(description) {
		result.Description = utils.ToStringPointer(description)
	}
	if ignoreUnavailableNsxCluster, ok := object["ignore_unavailable_nsx_cluster"]; ok && !validationutils.IsEmpty(ignoreUnavailableNsxCluster) {
		result.IgnoreUnavailableNsxtCluster = utils.ToBoolPointer(ignoreUnavailableNsxCluster)
	}
	if subnetsRaw, ok := object["subnet"]; ok {
		subnetsList := subnetsRaw.([]interface{})
		if len(subnetsList) > 0 {
			subnets := []vcf.IpAddressPoolSubnetSpec{}
			for _, subnetsListEntry := range subnetsList {
				ipAddressPoolSubnetSpec, err := getIpAddressPoolSubnetSpecFromSchema(subnetsListEntry.(map[string]interface{}))
				if err != nil {
					return nil, err
				}
				subnets = append(subnets, *ipAddressPoolSubnetSpec)
			}
			result.Subnets = &subnets
		}
	}

	return result, nil
}

func GetInstallerIpAddressPoolSpecFromSchema(object map[string]interface{}) (*installer.IpAddressPoolSpec, error) {
	result := &installer.IpAddressPoolSpec{}
	if object == nil {
		return nil, fmt.Errorf("cannot convert to IpAddressPoolSpec, object is nil")
	}
	name := object["name"].(string)
	if len(name) == 0 {
		return nil, fmt.Errorf("cannot convert to IpAddressPoolSpec, name is required")
	}
	result.Name = name
	if description, ok := object["description"]; ok && !validationutils.IsEmpty(description) {
		result.Description = utils.ToStringPointer(description)
	}
	if ignoreUnavailableNsxCluster, ok := object["ignore_unavailable_nsx_cluster"]; ok && !validationutils.IsEmpty(ignoreUnavailableNsxCluster) {
		result.IgnoreUnavailableNsxtCluster = utils.ToBoolPointer(ignoreUnavailableNsxCluster)
	}
	if subnetsRaw, ok := object["subnet"]; ok {
		subnetsList := subnetsRaw.([]interface{})
		if len(subnetsList) > 0 {
			subnets := []installer.IpAddressPoolSubnetSpec{}
			for _, subnetsListEntry := range subnetsList {
				ipAddressPoolSubnetSpec, err := getInstallerIpAddressPoolSubnetSpecFromSchema(subnetsListEntry.(map[string]interface{}))
				if err != nil {
					return nil, err
				}
				subnets = append(subnets, *ipAddressPoolSubnetSpec)
			}
			result.Subnets = &subnets
		}
	}

	return result, nil
}

func getIpAddressPoolSubnetSpecFromSchema(object map[string]interface{}) (*vcf.IpAddressPoolSubnetSpec, error) {
	result := &vcf.IpAddressPoolSubnetSpec{}
	if object == nil {
		return nil, fmt.Errorf("cannot convert to IPAddressPoolSubnetSpec, object is nil")
	}
	cidr := object["cidr"].(string)
	if len(cidr) == 0 {
		return nil, fmt.Errorf("cannot convert to IPAddressPoolSubnetSpec, cidr is required")
	}
	gateway := object["gateway"].(string)
	if len(gateway) == 0 {
		return nil, fmt.Errorf("cannot convert to IPAddressPoolSubnetSpec, gateway is required")
	}
	result.Cidr = cidr
	result.Gateway = gateway
	if ipAddressPoolRangeRaw, ok := object["ip_address_pool_range"]; ok {
		ipAddressPoolRangeList := ipAddressPoolRangeRaw.([]interface{})
		if len(ipAddressPoolRangeList) > 0 {
			result.IpAddressPoolRanges = []vcf.IpAddressPoolRangeSpec{}
			for _, ipAddressPoolRangeEntry := range ipAddressPoolRangeList {
				ipAddressPoolSubnetSpec := vcf.IpAddressPoolRangeSpec{}
				ipAddressPoolRangeMap := ipAddressPoolRangeEntry.(map[string]interface{})
				start := ipAddressPoolRangeMap["start"].(string)
				end := ipAddressPoolRangeMap["end"].(string)
				ipAddressPoolSubnetSpec.Start = start
				ipAddressPoolSubnetSpec.End = end
				result.IpAddressPoolRanges = append(result.IpAddressPoolRanges, ipAddressPoolSubnetSpec)
			}
		}
	}

	return result, nil
}

func getInstallerIpAddressPoolSubnetSpecFromSchema(object map[string]interface{}) (*installer.IpAddressPoolSubnetSpec, error) {
	result := &installer.IpAddressPoolSubnetSpec{}
	if object == nil {
		return nil, fmt.Errorf("cannot convert to IPAddressPoolSubnetSpec, object is nil")
	}
	cidr := object["cidr"].(string)
	if len(cidr) == 0 {
		return nil, fmt.Errorf("cannot convert to IPAddressPoolSubnetSpec, cidr is required")
	}
	gateway := object["gateway"].(string)
	if len(gateway) == 0 {
		return nil, fmt.Errorf("cannot convert to IPAddressPoolSubnetSpec, gateway is required")
	}
	result.Cidr = cidr
	result.Gateway = gateway
	if ipAddressPoolRangeRaw, ok := object["ip_address_pool_range"]; ok {
		ipAddressPoolRangeList := ipAddressPoolRangeRaw.([]interface{})
		if len(ipAddressPoolRangeList) > 0 {
			result.IpAddressPoolRanges = []installer.IpAddressPoolRangeSpec{}
			for _, ipAddressPoolRangeEntry := range ipAddressPoolRangeList {
				ipAddressPoolSubnetSpec := installer.IpAddressPoolRangeSpec{}
				ipAddressPoolRangeMap := ipAddressPoolRangeEntry.(map[string]interface{})
				start := ipAddressPoolRangeMap["start"].(string)
				end := ipAddressPoolRangeMap["end"].(string)
				ipAddressPoolSubnetSpec.Start = start
				ipAddressPoolSubnetSpec.End = end
				result.IpAddressPoolRanges = append(result.IpAddressPoolRanges, ipAddressPoolSubnetSpec)
			}
		}
	}

	return result, nil
}
