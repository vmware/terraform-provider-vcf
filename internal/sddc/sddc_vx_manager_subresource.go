// © Broadcom. All Rights Reserved.
// The term “Broadcom” refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

package sddc

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/vmware/vcf-sdk-go/models"

	utils "github.com/vmware/terraform-provider-vcf/internal/resource_utils"
)

func GetVxManagerSchema() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeList,
		Optional: true,
		MaxItems: 1,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"default_admin_user_credentials": getCredentialsSchema(),
				"default_root_user_credentials":  getCredentialsSchema(),
				"ssh_thumbprint": {
					Type:        schema.TypeString,
					Description: "VxRail Manager SSH thumbprint (RSA SHA256)",
					Optional:    true,
				},
				"ssl_thumbprint": {
					Type:        schema.TypeString,
					Description: "VxRail Manager SSL thumbprint (SHA256)",
					Optional:    true,
				},
				"vx_manager_hostname": {
					Type:        schema.TypeString,
					Description: "VxManager host name",
					Required:    true,
				},
			},
		},
	}
}

func GetVxManagerSpecFromSchema(rawData []interface{}) *models.VxManagerSpec {
	if len(rawData) <= 0 {
		return nil
	}
	data := rawData[0].(map[string]interface{})
	sshThumbprint := data["ssh_thumbprint"].(string)
	sslThumbprint := data["ssl_thumbprint"].(string)
	vxManagerHostName := data["vx_manager_hostname"].(string)

	vxManagerSpecBinding := &models.VxManagerSpec{
		SSHThumbprint:     sshThumbprint,
		SSLThumbprint:     sslThumbprint,
		VxManagerHostName: utils.ToStringPointer(vxManagerHostName),
	}
	if defaultAdminUserCredentials := getCredentialsFromSchema(data["default_admin_user_credentials"].([]interface{})); defaultAdminUserCredentials != nil {
		vxManagerSpecBinding.DefaultAdminUserCredentials = defaultAdminUserCredentials
	}
	if defaultRootUserCredentials := getCredentialsFromSchema(data["default_root_user_credentials"].([]interface{})); defaultRootUserCredentials != nil {
		vxManagerSpecBinding.DefaultRootUserCredentials = defaultRootUserCredentials
	}

	return vxManagerSpecBinding
}
