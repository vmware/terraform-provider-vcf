/*
 *  Copyright 2023 VMware, Inc.
 *    SPDX-License-Identifier: MPL-2.0
 */

package sddc

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	utils "github.com/vmware/terraform-provider-vcf/internal/resource_utils"
	validation_utils "github.com/vmware/terraform-provider-vcf/internal/validation"
	"github.com/vmware/vcf-sdk-go/models"
)

func GetSddcManagerSchema() *schema.Schema {
	sddcManagerSchema := &schema.Schema{
		Type:     schema.TypeList,
		MaxItems: 1,
		Optional: true,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"hostname": {
					Type:         schema.TypeString,
					Description:  "SDDC Manager Hostname. If just the short hostname is provided, then FQDN will be generated using the \"domain\" from dns configuration, length 3-63",
					Optional:     true,
					ValidateFunc: validation.StringLenBetween(3, 63),
				},
				"ip_address": {
					Type:         schema.TypeString,
					Description:  "SDDC Manager IPv4 address",
					Optional:     true,
					ValidateFunc: validation.IsIPAddress,
				},
				"local_user_password": {
					Type:         schema.TypeString,
					Description:  "The local account is a built-in admin account (password for the break glass user admin@local) in VCF that can be used in emergency scenarios. The password of this account must be at least 12 characters long. It also must contain at-least 1 uppercase, 1 lowercase, 1 special character specified in braces [!%@$^#?] and 1 digit. In addition, a character cannot be repeated more than 3 times consecutively.",
					Optional:     true,
					ValidateFunc: validation_utils.ValidatePassword,
				},
				"root_user_credentials":   getCredentialsSchema(),
				"second_user_credentials": getCredentialsSchema(),
			},
		},
	}
	sddcManagerSchemaElements := sddcManagerSchema.Elem.(*schema.Resource)
	sddcManagerSchemaElements.Schema["root_user_credentials"].Description = "Root user credentials for the SDDC Manager VM, UserName must be root. Password needs to be a strong password with at least one alphabet and one special character and at least 8 characters in length."
	sddcManagerSchemaElements.Schema["second_user_credentials"].Description = "Second user credentials for the SDDC Manager VM, UserName must be vcf.  Password needs to be a strong password with at least one alphabet and one special character and at least 8 characters in length."

	return sddcManagerSchema
}

func GetSddcManagerSpecFromSchema(rawData []interface{}) *models.SDDCManagerSpec {
	if len(rawData) <= 0 {
		return nil
	}
	data := rawData[0].(map[string]interface{})
	hostname := data["hostname"].(string)
	ipAddress := data["ip_address"].(string)
	localUserPassword := data["local_user_password"].(string)

	sddcManagerSpec := &models.SDDCManagerSpec{
		Hostname:          utils.ToStringPointer(hostname),
		IPAddress:         utils.ToStringPointer(ipAddress),
		LocalUserPassword: localUserPassword,
	}
	if rootUserCredentialsData := getCredentialsFromSchema(data["root_user_credentials"].([]interface{})); rootUserCredentialsData != nil {
		sddcManagerSpec.RootUserCredentials = rootUserCredentialsData
	}
	if secondUserCredentialsData := getCredentialsFromSchema(data["second_user_credentials"].([]interface{})); secondUserCredentialsData != nil {
		sddcManagerSpec.SecondUserCredentials = secondUserCredentialsData
	}
	return sddcManagerSpec
}
