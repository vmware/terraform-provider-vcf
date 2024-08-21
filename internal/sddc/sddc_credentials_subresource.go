// © Broadcom. All Rights Reserved.
// The term “Broadcom” refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

package sddc

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/vmware/vcf-sdk-go/models"

	utils "github.com/vmware/terraform-provider-vcf/internal/resource_utils"
	"github.com/vmware/terraform-provider-vcf/internal/validation"
)

func getCredentialsSchema() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeList,
		Optional: true,
		MaxItems: 1,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"password": {
					Type:         schema.TypeString,
					Required:     true,
					ValidateFunc: validation.ValidatePassword,
				},
				"username": {
					Type:     schema.TypeString,
					Required: true,
				},
			},
		},
	}
}

func getCredentialsFromSchema(rawData []interface{}) *models.SDDCCredentials {
	if len(rawData) <= 0 {
		return nil
	}
	data := rawData[0].(map[string]interface{})
	password := utils.ToStringPointer(data["password"])
	username := utils.ToStringPointer(data["username"])

	credentialsBinding := &models.SDDCCredentials{
		Password: password,
		Username: username,
	}
	return credentialsBinding
}
