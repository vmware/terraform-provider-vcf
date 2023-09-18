/* Copyright 2023 VMware, Inc.
   SPDX-License-Identifier: MPL-2.0 */

package provider

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/vmware/terraform-provider-vcf/internal/api_client"
	"github.com/vmware/terraform-provider-vcf/internal/constants"
)

// Provider returns the resource configuration of the VCF provider.
func Provider() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"sddc_manager_username": {
				Type:          schema.TypeString,
				Optional:      true,
				Description:   "Username to authenticate to SDDC Manager",
				ConflictsWith: []string{"cloud_builder_username", "cloud_builder_password", "cloud_builder_host"},
				RequiredWith:  []string{"sddc_manager_password", "sddc_manager_host"},
				DefaultFunc:   schema.EnvDefaultFunc(constants.VcfTestUsername, nil),
			},
			"sddc_manager_password": {
				Type:          schema.TypeString,
				Optional:      true,
				Description:   "Password to authenticate to SDDC Manager",
				ConflictsWith: []string{"cloud_builder_username", "cloud_builder_password", "cloud_builder_host"},
				RequiredWith:  []string{"sddc_manager_username", "sddc_manager_host"},
				DefaultFunc:   schema.EnvDefaultFunc(constants.VcfTestPassword, nil),
			},
			"sddc_manager_host": {
				Type:          schema.TypeString,
				Optional:      true,
				Description:   "Fully qualified domain name or IP address of the SDDC Manager",
				ConflictsWith: []string{"cloud_builder_username", "cloud_builder_password", "cloud_builder_host"},
				RequiredWith:  []string{"sddc_manager_username", "sddc_manager_password"},
				DefaultFunc:   schema.EnvDefaultFunc(constants.VcfTestUrl, nil),
			},
			"cloud_builder_username": {
				Type:          schema.TypeString,
				Optional:      true,
				Description:   "Username to authenticate to CloudBuilder",
				ConflictsWith: []string{"sddc_manager_username", "sddc_manager_password", "sddc_manager_host"},
				RequiredWith:  []string{"cloud_builder_password", "cloud_builder_host"},
				DefaultFunc:   schema.EnvDefaultFunc(constants.CloudBuilderTestUsername, nil),
			},
			"cloud_builder_password": {
				Type:          schema.TypeString,
				Optional:      true,
				Description:   "Password to authenticate to CloudBuilder",
				ConflictsWith: []string{"sddc_manager_username", "sddc_manager_password", "sddc_manager_host"},
				RequiredWith:  []string{"cloud_builder_username", "cloud_builder_host"},
				DefaultFunc:   schema.EnvDefaultFunc(constants.CloudBuilderTestPassword, nil),
			},
			"cloud_builder_host": {
				Type:          schema.TypeString,
				Optional:      true,
				Description:   "Fully qualified domain name or IP address of the CloudBuilder",
				ConflictsWith: []string{"sddc_manager_username", "sddc_manager_password", "sddc_manager_host"},
				RequiredWith:  []string{"cloud_builder_username", "cloud_builder_password"},
				DefaultFunc:   schema.EnvDefaultFunc(constants.CloudBuilderTestUrl, nil),
			},
			"allow_unverified_tls": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "If set, VMware VCF client will permit unverifiable TLS certificates.",
				DefaultFunc: schema.EnvDefaultFunc(constants.VcfTestAllowUnverifiedTls, false),
			},
		},

		DataSourcesMap: map[string]*schema.Resource{
			"vcf_domain":  DataSourceDomain(),
			"vcf_cluster": DataSourceCluster(),
		},

		ResourcesMap: map[string]*schema.Resource{
			"vcf_instance":     ResourceVcfInstance(),
			"vcf_user":         ResourceUser(),
			"vcf_network_pool": ResourceNetworkPool(),
			"vcf_ceip":         ResourceCeip(),
			"vcf_host":         ResourceHost(),
			"vcf_domain":       ResourceDomain(),
			"vcf_cluster":      ResourceCluster(),
		},

		ConfigureContextFunc: providerConfigure,
	}
}

func providerConfigure(_ context.Context, data *schema.ResourceData) (interface{}, diag.Diagnostics) {
	vcfUsername, isVcfUsernameSet := data.GetOk("sddc_manager_username")
	allowUnverifiedTls := data.Get("allow_unverified_tls")
	if isVcfUsernameSet {
		password, isSetPassword := data.GetOk("sddc_manager_password")
		hostName, isSetHost := data.GetOk("sddc_manager_host")
		if !isVcfUsernameSet || !isSetPassword || !isSetHost {
			return nil, diag.Errorf("SDDC Manager vcfUsername, password, and URL must be provided")
		}
		var sddcManagerClient = api_client.NewSddcManagerClient(vcfUsername.(string), password.(string),
			hostName.(string), allowUnverifiedTls.(bool))
		err := sddcManagerClient.Connect()
		if err != nil {
			return nil, diag.FromErr(err)
		}
		return sddcManagerClient, nil
	} else {
		cbUsername, isCbUsernameSet := data.GetOk("cloud_builder_username")
		password, isSetPassword := data.GetOk("cloud_builder_password")
		hostName, isSetHost := data.GetOk("cloud_builder_host")
		if !isCbUsernameSet || !isSetPassword || !isSetHost {
			return nil, diag.Errorf("CloudBuilder username, password, and URL must be provided")
		}
		var cloudBuilderClient = api_client.NewCloudBuilderClient(cbUsername.(string), password.(string),
			hostName.(string), allowUnverifiedTls.(bool))
		return cloudBuilderClient, nil
	}
}
