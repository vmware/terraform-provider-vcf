/* Copyright 2023 VMware, Inc.
   SPDX-License-Identifier: MPL-2.0 */

package provider

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	vcf "github.com/vmware/terraform-provider-vcf/internal/vcf"
)

// Provider returns the resource configuration of the VCF provider.
func Provider() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"sddc_manager_username": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "SDDC Manager username.",
			},
			"sddc_manager_password": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "SDDC Manager password.",
			},
			"sddc_manager_host": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "SDDC Manager host.",
			},
		},

		DataSourcesMap: map[string]*schema.Resource{},

		ResourcesMap: map[string]*schema.Resource{
			"vcf_user":         vcf.ResourceUser(),
			"vcf_network_pool": vcf.ResourceNetworkPool(),
			"vcf_ceip":         vcf.ResourceCeip(),
			"vcf_host":         vcf.ResourceHost(),
		},

		ConfigureContextFunc: providerConfigure,
	}
}

func providerConfigure(_ context.Context, data *schema.ResourceData) (interface{}, diag.Diagnostics) {
	username, isSetUsername := data.GetOk("sddc_manager_username")
	password, isSetPassword := data.GetOk("sddc_manager_password")
	hostName, isSetHost := data.GetOk("sddc_manager_host")

	if !isSetUsername || !isSetPassword || !isSetHost {
		return nil, diag.Errorf("SDDC Manager username, password, and host must be provided")
	}

	var newClient = vcf.NewSddcManagerClient(username.(string), password.(string), hostName.(string))
	newClient.Connect()
	return newClient, nil
}
