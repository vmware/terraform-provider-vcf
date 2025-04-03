// © Broadcom. All Rights Reserved.
// The term “Broadcom” refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

package sddc

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	utils "github.com/vmware/terraform-provider-vcf/internal/resource_utils"
	"github.com/vmware/vcf-sdk-go/installer"
)

func GetVcfAutomationSchema() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeList,
		Optional: true,
		MaxItems: 1,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"hostname": {
					Type:        schema.TypeString,
					Description: "Host name for the automation appliance",
					Required:    true,
				},
				"admin_user_password": {
					Type:        schema.TypeString,
					Description: "Administrator password",
					Optional:    true,
					Sensitive:   true,
				},
				"internal_cluster_cidr": {
					Type:        schema.TypeString,
					Description: "Internal Cluster CIDR. One among: 198.18.0.0/15, 240.0.0.0/15, 250.0.0.0/15",
					Required:    true,
				},
				"node_prefix": {
					Type:        schema.TypeString,
					Description: "Node Prefix. It cannot be blank and must begin and end with an alphanumeric character, and can only contain lowercase alphanumeric characters or hyphens.",
					Optional:    true,
				},
				"ip_pool": {
					Type:        schema.TypeList,
					Description: "List of IP addresses.  For Standard deployment model two IP addresses need to be specified and for High Availability four IP addresses need to be specified",
					Required:    true,
					Elem: &schema.Schema{
						Type: schema.TypeString,
					},
				},
			},
		},
	}
}

func GetVcfAutomationSpecFromSchema(rawData []interface{}) *installer.VcfAutomationSpec {
	if len(rawData) <= 0 {
		return nil
	}
	data := rawData[0].(map[string]interface{})

	var adminPassword *string
	if data["admin_user_password"].(string) != "" {
		adminPassword = utils.ToPointer[string](data["admin_user_password"])
	}

	var nodePrefix *string
	if data["node_prefix"].(string) != "" {
		nodePrefix = utils.ToPointer[string](data["node_prefix"])
	}

	var internalClusterCidr *string
	if data["internal_cluster_cidr"].(string) != "" {
		internalClusterCidr = utils.ToPointer[string](data["internal_cluster_cidr"])
	}

	var ipPools *[]string
	if data["ip_pool"] != nil && len(data["ip_pool"].([]interface{})) > 0 {
		ipPools = utils.ToPointer[[]string](utils.ToStringSlice(data["ip_pool"].([]interface{})))
	}

	spec := &installer.VcfAutomationSpec{
		AdminUserPassword:   adminPassword,
		Hostname:            data["hostname"].(string),
		InternalClusterCidr: internalClusterCidr,
		IpPool:              ipPools,
		NodePrefix:          nodePrefix,
	}
	return spec
}
