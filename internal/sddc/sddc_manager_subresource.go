// © Broadcom. All Rights Reserved.
// The term “Broadcom” refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

package sddc

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/vmware/vcf-sdk-go/installer"

	validation_utils "github.com/vmware/terraform-provider-vcf/internal/validation"
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
				"local_user_password": {
					Type:         schema.TypeString,
					Description:  "The local account is a built-in admin account (password for the break glass user admin@local) in VCF that can be used in emergency scenarios. The password of this account must be at least 12 characters long. It also must contain at-least 1 uppercase, 1 lowercase, 1 special character specified in braces [!%@$^#?] and 1 digit. In addition, a character cannot be repeated more than 3 times consecutively.",
					Optional:     true,
					ValidateFunc: validation_utils.ValidatePassword,
				},
				"root_user_password": {
					Type:         schema.TypeString,
					Description:  "The password for the root user",
					Required:     true,
					ValidateFunc: validation_utils.ValidatePassword,
				},
				"ssh_password": {
					Type:         schema.TypeString,
					Description:  "The password for the vcf user (ssh connections only)",
					Required:     true,
					ValidateFunc: validation_utils.ValidatePassword,
				},
			},
		},
	}

	return sddcManagerSchema
}

func GetSddcManagerSpecFromSchema(rawData []interface{}) *installer.SddcManagerSpec {
	if len(rawData) <= 0 {
		return nil
	}
	data := rawData[0].(map[string]interface{})
	hostname := data["hostname"].(string)
	localUserPassword := data["local_user_password"].(string)
	rootUserPassword := data["root_user_password"].(string)
	sshPassword := data["ssh_password"].(string)

	sddcManagerSpec := &installer.SddcManagerSpec{
		Hostname: hostname,
	}

	if localUserPassword != "" {
		sddcManagerSpec.LocalUserPassword = &localUserPassword
	}

	if rootUserPassword != "" {
		sddcManagerSpec.RootPassword = &rootUserPassword
	}

	if sshPassword != "" {
		sddcManagerSpec.SshPassword = &sshPassword
	}

	return sddcManagerSpec
}
