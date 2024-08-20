// © Broadcom. All Rights Reserved.
// The term “Broadcom” refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/vmware/terraform-provider-vcf/internal/api_client"
	"github.com/vmware/terraform-provider-vcf/internal/constants"
)

// Provider returns the resource configuration of the provider.
func Provider() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"sddc_manager_username": {
				Type:          schema.TypeString,
				Optional:      true,
				Description:   "The username to authenticate to the SDDC Manager instance.",
				ConflictsWith: []string{"cloud_builder_username", "cloud_builder_password", "cloud_builder_host"},
				RequiredWith:  []string{"sddc_manager_password", "sddc_manager_host"},
				DefaultFunc:   schema.EnvDefaultFunc(constants.VcfTestUsername, nil),
			},
			"sddc_manager_password": {
				Type:          schema.TypeString,
				Optional:      true,
				Description:   "The password to authenticate to the SDDC Manager instance.",
				ConflictsWith: []string{"cloud_builder_username", "cloud_builder_password", "cloud_builder_host"},
				RequiredWith:  []string{"sddc_manager_username", "sddc_manager_host"},
				DefaultFunc:   schema.EnvDefaultFunc(constants.VcfTestPassword, nil),
			},
			"sddc_manager_host": {
				Type:          schema.TypeString,
				Optional:      true,
				Description:   "The fully qualified domain name or IP address of the SDDC Manager instance",
				ConflictsWith: []string{"cloud_builder_username", "cloud_builder_password", "cloud_builder_host"},
				RequiredWith:  []string{"sddc_manager_username", "sddc_manager_password"},
				DefaultFunc:   schema.EnvDefaultFunc(constants.VcfTestUrl, nil),
			},
			"cloud_builder_username": {
				Type:          schema.TypeString,
				Optional:      true,
				Description:   "The username to authenticate to the Cloud Builder instance.",
				ConflictsWith: []string{"sddc_manager_username", "sddc_manager_password", "sddc_manager_host"},
				RequiredWith:  []string{"cloud_builder_password", "cloud_builder_host"},
				DefaultFunc:   schema.EnvDefaultFunc(constants.CloudBuilderTestUsername, nil),
			},
			"cloud_builder_password": {
				Type:          schema.TypeString,
				Optional:      true,
				Description:   "The password to authenticate to the Cloud Builder instance.",
				ConflictsWith: []string{"sddc_manager_username", "sddc_manager_password", "sddc_manager_host"},
				RequiredWith:  []string{"cloud_builder_username", "cloud_builder_host"},
				DefaultFunc:   schema.EnvDefaultFunc(constants.CloudBuilderTestPassword, nil),
			},
			"cloud_builder_host": {
				Type:          schema.TypeString,
				Optional:      true,
				Description:   "The fully qualified domain name or IP address of the Cloud Builder instance.",
				ConflictsWith: []string{"sddc_manager_username", "sddc_manager_password", "sddc_manager_host"},
				RequiredWith:  []string{"cloud_builder_username", "cloud_builder_password"},
				DefaultFunc:   schema.EnvDefaultFunc(constants.CloudBuilderTestUrl, nil),
			},
			"allow_unverified_tls": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Allow unverified TLS certificates.",
				DefaultFunc: schema.EnvDefaultFunc(constants.VcfTestAllowUnverifiedTls, false),
			},
		},

		DataSourcesMap: map[string]*schema.Resource{
			"vcf_cluster":      DataSourceCluster(),
			"vcf_domain":       DataSourceDomain(),
			"vcf_credentials":  DataSourceCredentials(),
			"vcf_network_pool": DataSourceNetworkPool(),
		},

		ResourcesMap: map[string]*schema.Resource{
			"vcf_certificate":                    ResourceCertificate(),
			"vcf_certificate_authority":          ResourceCertificateAuthority(),
			"vcf_ceip":                           ResourceCeip(),
			"vcf_cluster":                        ResourceCluster(),
			"vcf_cluster_personality":            ResourceClusterPersonality(),
			"vcf_credentials_auto_rotate_policy": ResourceCredentialsAutoRotatePolicy(),
			"vcf_credentials_rotate":             ResourceCredentialsRotate(),
			"vcf_credentials_update":             ResourceCredentialsUpdate(),
			"vcf_csr":                            ResourceCsr(),
			"vcf_domain":                         ResourceDomain(),
			"vcf_edge_cluster":                   ResourceEdgeCluster(),
			"vcf_external_certificate":           ResourceExternalCertificate(),
			"vcf_host":                           ResourceHost(),
			"vcf_instance":                       ResourceVcfInstance(),
			"vcf_user":                           ResourceUser(),
		},

		ConfigureContextFunc: providerConfigure,
	}
}

func providerConfigure(_ context.Context, data *schema.ResourceData) (interface{}, diag.Diagnostics) {
	sddcManagerUsername, isVcfUsernameSet := data.GetOk("sddc_manager_username")
	allowUnverifiedTLS := data.Get("allow_unverified_tls")
	if isVcfUsernameSet {
		password, isSetPassword := data.GetOk("sddc_manager_password")
		hostName, isSetHost := data.GetOk("sddc_manager_host")
		if !isVcfUsernameSet || !isSetPassword || !isSetHost {
			return nil, diag.Errorf("SDDC Manager username, password, and host must be provided")
		}
		var sddcManagerClient = api_client.NewSddcManagerClient(sddcManagerUsername.(string), password.(string),
			hostName.(string), allowUnverifiedTLS.(bool))
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
			return nil, diag.Errorf("Cloud Builder username, password, and host must be provided")
		}
		var cloudBuilderClient = api_client.NewCloudBuilderClient(cbUsername.(string), password.(string),
			hostName.(string), allowUnverifiedTLS.(bool))
		return cloudBuilderClient, nil
	}
}
