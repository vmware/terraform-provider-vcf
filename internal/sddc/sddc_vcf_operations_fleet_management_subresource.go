// © Broadcom. All Rights Reserved.
// The term “Broadcom” refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

package sddc

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	utils "github.com/vmware/terraform-provider-vcf/internal/resource_utils"
	"github.com/vmware/vcf-sdk-go/installer"
)

func GetVcfOperationsFleetManagementSchema() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeList,
		Optional: true,
		MaxItems: 1,
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
				"admin_user_password": {
					Type:        schema.TypeString,
					Description: "root password",
					Optional:    true,
					Sensitive:   true,
				},
			},
		},
	}
}

func GetVcfOperationsFleetManagementSpecFromSchema(rawData []interface{}) *installer.VcfOperationsFleetManagementSpec {
	if len(rawData) <= 0 {
		return nil
	}
	data := rawData[0].(map[string]interface{})

	var rootPassword *string
	if data["root_user_password"].(string) != "" {
		rootPassword = utils.ToPointer[string](data["root_user_password"])
	}

	var adminPassword *string
	if data["admin_user_password"].(string) != "" {
		adminPassword = utils.ToPointer[string](data["admin_user_password"])
	}

	spec := &installer.VcfOperationsFleetManagementSpec{
		AdminUserPassword: adminPassword,
		RootUserPassword:  rootPassword,
		Hostname:          data["hostname"].(string),
	}
	return spec
}
