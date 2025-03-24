// © Broadcom. All Rights Reserved.
// The term “Broadcom” refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

package sddc

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	validation_utils "github.com/vmware/terraform-provider-vcf/internal/validation"
	"github.com/vmware/vcf-sdk-go/installer"
)

var vmSizeValues = []string{"tiny", "small", "medium", "large", "xlarge"}
var storageSizes = []string{"lstorage", "xlstorage"}

func GetVcenterSchema() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeList,
		Required: true,
		MaxItems: 1,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"root_vcenter_password": {
					Type:         schema.TypeString,
					Description:  "vCenter root password. The password must be between 8 characters and 20 characters long. It must also contain at least one uppercase and lowercase letter, one number, and one character from '! \" # $ % & ' ( ) * + , - . / : ; < = > ? @ [ \\ ] ^ _ ` { &Iota; } ~' and all characters must be ASCII. Space is not allowed in password.",
					Required:     true,
					Sensitive:    true,
					ValidateFunc: validation_utils.ValidatePassword,
				},
				"ssl_thumbprint": {
					Type:        schema.TypeString,
					Description: "vCenter Server SSL thumbprint (SHA256)",
					Optional:    true,
				},
				"storage_size": {
					Type:         schema.TypeString,
					Description:  "vCenter VM storage size. One among:lstorage, xlstorage",
					Optional:     true,
					ValidateFunc: validation.StringInSlice(storageSizes, false),
				},
				"vcenter_hostname": {
					Type:        schema.TypeString,
					Description: "vCenter Server hostname address. If just the short hostname is provided, then FQDN will be generated using the \"domain\" from dns configuration",
					Required:    true,
				},
				"vm_size": {
					Type:         schema.TypeString,
					Description:  "vCenter Server Appliance  size. One among: tiny, small, medium, large, xlarge",
					Optional:     true,
					ValidateFunc: validation.StringInSlice(vmSizeValues, false),
				},
			},
		},
	}
}

func GetVcenterSpecFromSchema(rawData []interface{}) *installer.SddcVcenterSpec {
	if len(rawData) <= 0 {
		return nil
	}
	data := rawData[0].(map[string]interface{})
	rootVcenterPassword := data["root_vcenter_password"].(string)
	sslThumbprint := data["ssl_thumbprint"].(string)
	storageSize := data["storage_size"].(string)
	vcenterHostname := data["vcenter_hostname"].(string)
	vmSize := data["vm_size"].(string)

	vcenterSpecBinding := &installer.SddcVcenterSpec{
		RootVcenterPassword: rootVcenterPassword,
		SslThumbprint:       &sslThumbprint,
		StorageSize:         &storageSize,
		VcenterHostname:     vcenterHostname,
		VmSize:              &vmSize,
	}
	return vcenterSpecBinding
}
