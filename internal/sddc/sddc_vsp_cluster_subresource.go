// © Broadcom. All Rights Reserved.
// The term “Broadcom” refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

package sddc

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	utils "github.com/vmware/terraform-provider-vcf/internal/resource_utils"
	"github.com/vmware/vcf-sdk-go/installer"
)

func GetVspClusterSchema() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeList,
		Optional: true,
		MaxItems: 1,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"platform_fqdn": {
					Type:        schema.TypeString,
					Description: "FQDN of the VCF management services platform cluster",
					Required:    true,
				},
				"instance_fqdn": {
					Type:        schema.TypeString,
					Description: "FQDN of the VCF instance",
					Required:    true,
				},
				"fleet_fqdn": {
					Type:        schema.TypeString,
					Description: "FQDN of the fleet. Provided for VVF and for the primary VCF instance only",
					Optional:    true,
				},
				"ipv4_pool": getVspIpv4PoolSchema(),
				"system_user_password": {
					Type:        schema.TypeString,
					Description: "SSH password for vmware-system-user and admin@vsp.local on the cluster nodes. If blank the password will be auto-generated",
					Optional:    true,
					Sensitive:   true,
				},
				"size": {
					Type:         schema.TypeString,
					Description:  "Size of the cluster. One among: small, small_ha, medium, large",
					Optional:     true,
					ValidateFunc: validation.StringInSlice([]string{"small", "small_ha", "medium", "large"}, false),
				},
				"internal_cluster_cidr_ipv4": {
					Type:         schema.TypeString,
					Description:  "Internal cluster CIDR for IPv4. One among: 198.18.0.0/15, 240.0.0.0/15, 250.0.0.0/15",
					Optional:     true,
					ValidateFunc: validation.StringLenBetween(9, 18),
				},
				"version": {
					Type:        schema.TypeString,
					Description: "Version of the cluster",
					Optional:    true,
				},
				"use_existing_deployment": {
					Type:        schema.TypeBool,
					Description: "Import an existing deployment instead of deploying one",
					Optional:    true,
				},
				"ssl_thumbprint": {
					Type:        schema.TypeString,
					Description: "SSL thumbprint (SHA256) of the product certificate. Required when importing an existing deployment",
					Optional:    true,
				},
			},
		},
	}
}

func getVspIpv4PoolSchema() *schema.Schema {
	return &schema.Schema{
		Type:        schema.TypeList,
		Description: "IPv4 addresses available to the cluster. One of cidr, ip_range or addresses is required",
		Required:    true,
		MaxItems:    1,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"cidr": {
					Type:          schema.TypeString,
					Description:   "Network CIDR",
					Optional:      true,
					ValidateFunc:  validation.StringLenBetween(9, 18),
					ConflictsWith: []string{"ip_range", "addresses"},
				},
				"ip_range": {
					Type:          schema.TypeList,
					Description:   "Range of IP addresses",
					Optional:      true,
					MaxItems:      1,
					ConflictsWith: []string{"cidr", "addresses"},
					Elem: &schema.Resource{
						Schema: map[string]*schema.Schema{
							"start_ip_address": {
								Type:         schema.TypeString,
								Required:     true,
								ValidateFunc: validation.IsIPAddress,
							},
							"end_ip_address": {
								Type:         schema.TypeString,
								Required:     true,
								ValidateFunc: validation.IsIPAddress,
							},
						},
					},
				},
				"addresses": {
					Type:          schema.TypeList,
					Description:   "List of IP addresses",
					Optional:      true,
					ConflictsWith: []string{"cidr", "ip_range"},
					Elem:          &schema.Schema{Type: schema.TypeString},
				},
				"excluded_addresses": {
					Type:        schema.TypeList,
					Description: "List of IP addresses to exclude. Applies to cidr and ip_range",
					Optional:    true,
					Elem:        &schema.Schema{Type: schema.TypeString},
				},
			},
		},
	}
}

func GetVspClusterSpecFromSchema(rawData []interface{}) *installer.SddcVspClusterSpec {
	if len(rawData) <= 0 {
		return nil
	}
	data := rawData[0].(map[string]interface{})

	spec := &installer.SddcVspClusterSpec{
		PlatformFqdn: data["platform_fqdn"].(string),
		InstanceFqdn: data["instance_fqdn"].(string),
		Ipv4Pool:     getVspIpv4PoolFromSchema(data["ipv4_pool"].([]interface{})),
	}

	if fleetFqdn := data["fleet_fqdn"].(string); fleetFqdn != "" {
		spec.FleetFqdn = &fleetFqdn
	}
	if systemUserPassword := data["system_user_password"].(string); systemUserPassword != "" {
		spec.SystemUserPassword = &systemUserPassword
	}
	if size := data["size"].(string); size != "" {
		spec.Size = &size
	}
	if internalClusterCidr := data["internal_cluster_cidr_ipv4"].(string); internalClusterCidr != "" {
		spec.InternalClusterCidrIpv4 = &internalClusterCidr
	}
	if version := data["version"].(string); version != "" {
		spec.Version = &version
	}
	if useExistingDeployment, ok := data["use_existing_deployment"].(bool); ok && useExistingDeployment {
		spec.UseExistingDeployment = &useExistingDeployment
	}
	if sslThumbprint := data["ssl_thumbprint"].(string); sslThumbprint != "" {
		spec.SslThumbprint = &sslThumbprint
	}

	return spec
}

func getVspIpv4PoolFromSchema(rawData []interface{}) installer.IPv4Pool {
	pool := installer.IPv4Pool{}
	if len(rawData) <= 0 || rawData[0] == nil {
		return pool
	}
	data := rawData[0].(map[string]interface{})

	if cidr := data["cidr"].(string); cidr != "" {
		pool.Cidr = &cidr
	}
	if ipRanges := getIncludeIPAddressRangesBindingFromSchema(data["ip_range"].([]interface{})); len(ipRanges) > 0 {
		pool.IpRange = &ipRanges[0]
	}
	if addresses := utils.ToStringSlice(data["addresses"].([]interface{})); len(addresses) > 0 {
		pool.Addresses = &addresses
	}
	if excluded := utils.ToStringSlice(data["excluded_addresses"].([]interface{})); len(excluded) > 0 {
		pool.ExcludedAddresses = &excluded
	}

	return pool
}
