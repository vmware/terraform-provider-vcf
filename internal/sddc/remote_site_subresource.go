/*
 *  Copyright 2023 VMware, Inc.
 *    SPDX-License-Identifier: MPL-2.0
 */

package sddc

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/vmware/vcf-sdk-go/models"
)

func GetRemoteSiteSchema() *schema.Schema {
	return &schema.Schema{
		Type:        schema.TypeList,
		Description: "Parameters for Remote site products",
		MaxItems:    1,
		Optional:    true,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"psc_address": {
					Type:        schema.TypeString,
					Description: "Remote region vCenter address",
					Optional:    true,
				},
				"ssl_thumbprint": {
					Type:        schema.TypeString,
					Description: "Remote region vCenter SSL thumbprint (SHA256)",
					Optional:    true,
				},
				"vc_credentials": getCredentialsSchema(),
			},
		},
	}
}

func GetRemoteSiteSpecFromSchema(rawData []interface{}) *models.RemoteSiteSpec {
	if len(rawData) <= 0 {
		return nil
	}
	data := rawData[0].(map[string]interface{})
	pscAddress := data["psc_address"].(string)
	sslThumbprint := data["ssl_thumbprint"].(string)

	remoteSiteSpecBinding := &models.RemoteSiteSpec{
		PscAddress:    pscAddress,
		SSLThumbprint: sslThumbprint,
	}
	if vcCredentials := getCredentialsFromSchema(data["vc_credentials"].([]interface{})); vcCredentials != nil {
		remoteSiteSpecBinding.VcCredentials = vcCredentials
	}

	return remoteSiteSpecBinding
}
