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

func GetVcfOperationsSchema() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeList,
		Optional: true,
		MaxItems: 1,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"admin_user_password": {
					Type:        schema.TypeString,
					Description: "Administrator password",
					Optional:    true,
					Sensitive:   true,
				},
				"appliance_size": {
					Type:         schema.TypeString,
					Description:  "Appliance size",
					Optional:     true,
					ValidateFunc: validation.StringInSlice([]string{"xsmall", "small", "medium", "large", "xlarge"}, true),
				},
				"load_balancer_fqdn": {
					Type:        schema.TypeString,
					Description: "FQDN of the load balancer",
					Optional:    true,
				},
				"node": getNodesSchema(),
			},
		},
	}
}

func getNodesSchema() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeList,
		Required: true,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"hostname": {
					Type:        schema.TypeString,
					Description: "Host name for the node",
					Required:    true,
				},
				"root_user_password": {
					Type:        schema.TypeString,
					Description: "root password",
					Optional:    true,
					Sensitive:   true,
				},
				"type": {
					Type:         schema.TypeString,
					Description:  "Type of the node",
					Required:     true,
					ValidateFunc: validation.StringInSlice([]string{"master", "replica", "data"}, true),
				},
			},
		},
	}
}

func GetVcfOperationsSpecFromSchema(rawData []interface{}) *installer.VcfOperationsSpec {
	if len(rawData) <= 0 {
		return nil
	}
	data := rawData[0].(map[string]interface{})

	var adminPassword *string
	if data["admin_user_password"].(string) != "" {
		adminPassword = utils.ToPointer[string](data["admin_user_password"])
	}

	var applianceSize *string
	if data["appliance_size"].(string) != "" {
		applianceSize = utils.ToPointer[string](data["appliance_size"])
	}

	var loadBalancerFqdn *string
	if data["load_balancer_fqdn"].(string) != "" {
		loadBalancerFqdn = utils.ToPointer[string](data["load_balancer_fqdn"])
	}

	spec := &installer.VcfOperationsSpec{
		AdminUserPassword: adminPassword,
		ApplianceSize:     applianceSize,
		LoadBalancerFqdn:  loadBalancerFqdn,
		Nodes:             getNodesFromSchema(data["node"].([]interface{})),
	}
	return spec
}

func getNodesFromSchema(data []interface{}) []installer.VcfOperationsNode {
	var nodes []installer.VcfOperationsNode
	if data != nil {
		nodes = make([]installer.VcfOperationsNode, len(data))
		for i, d := range data {
			nodeData := d.(map[string]interface{})
			node := installer.VcfOperationsNode{
				Hostname:         nodeData["hostname"].(string),
				RootUserPassword: utils.ToPointer[string](nodeData["root_user_password"].(string)),
				Type:             utils.ToPointer[string](nodeData["type"].(string)),
			}

			nodes[i] = node
		}
	}

	return nodes
}
