/*
 *  Copyright 2023 VMware, Inc.
 *    SPDX-License-Identifier: MPL-2.0
 */

package sddc

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	utils "github.com/vmware/terraform-provider-vcf/internal/resource_utils"
	"github.com/vmware/terraform-provider-vcf/internal/validation"
	"github.com/vmware/vcf-sdk-go/models"
)

func GetPscSchema() *schema.Schema {
	return &schema.Schema{
		Type:        schema.TypeList,
		Description: "Parameters for deployment/configuration of Platform Services Controller",
		Optional:    true,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"admin_user_sso_password": {
					Type:         schema.TypeString,
					Description:  "Admin user sso password. Password needs to be a strong password with at least one Uppercase alphabet, one lowercase alphabet, one digit and one special character specified in braces [!$%^] and 8-20 characters in length,and 3 maximum identical adjacent characters!",
					Required:     true,
					ValidateFunc: validation.ValidatePassword,
				},
				"psc_sso_domain": {
					Type:        schema.TypeString,
					Description: "PSC SSO Domain. Example: vsphere.local",
					Optional:    true,
				},
			},
		},
	}
}

func GetPscSpecsFromSchema(rawData []interface{}) []*models.PscSpec {
	var pscSpecsBindingsList []*models.PscSpec
	for _, pscSpec := range rawData {
		data := pscSpec.(map[string]interface{})
		adminUserSsoPassword := data["admin_user_sso_password"].(string)
		pscSsoDomain := data["psc_sso_domain"].(string)

		pscSpecsBinding := &models.PscSpec{
			AdminUserSSOPassword: utils.ToStringPointer(adminUserSsoPassword),
			PscSSOSpec: &models.PscSSOSpec{
				SSODomain: pscSsoDomain,
			},
		}
		pscSpecsBindingsList = append(pscSpecsBindingsList, pscSpecsBinding)
	}
	return pscSpecsBindingsList
}
