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

func GetVcfOperationsCollectorSchema() *schema.Schema {
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
				"appliance_size": {
					Type:         schema.TypeString,
					Description:  "Appliance size",
					Optional:     true,
					ValidateFunc: validation.StringInSlice([]string{"xsmall", "small", "medium", "large", "xlarge"}, true),
				},
			},
		},
	}
}

func GetVcfOperationsCollectorSpecFromSchema(rawData []interface{}) *installer.VcfOperationsCollectorSpec {
	if len(rawData) <= 0 {
		return nil
	}
	data := rawData[0].(map[string]interface{})

	var rootPassword *string
	if data["root_user_password"].(string) != "" {
		rootPassword = utils.ToPointer[string](data["root_user_password"])
	}

	var applianceSize *string
	if data["appliance_size"].(string) != "" {
		applianceSize = utils.ToPointer[string](data["appliance_size"])
	}

	spec := &installer.VcfOperationsCollectorSpec{
		ApplianceSize:    applianceSize,
		RootUserPassword: rootPassword,
		Hostname:         data["hostname"].(string),
	}
	return spec
}
