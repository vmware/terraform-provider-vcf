// © Broadcom. All Rights Reserved.
// The term “Broadcom” refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/vmware/terraform-provider-vcf/internal/version"

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
				ConflictsWith: []string{"installer_username", "installer_password", "installer_host"},
				RequiredWith:  []string{"sddc_manager_password", "sddc_manager_host"},
				DefaultFunc:   schema.EnvDefaultFunc(constants.VcfTestUsername, nil),
			},
			"sddc_manager_password": {
				Type:          schema.TypeString,
				Optional:      true,
				Description:   "The password to authenticate to the SDDC Manager instance.",
				ConflictsWith: []string{"installer_username", "installer_password", "installer_host"},
				RequiredWith:  []string{"sddc_manager_username", "sddc_manager_host"},
				DefaultFunc:   schema.EnvDefaultFunc(constants.VcfTestPassword, nil),
			},
			"sddc_manager_host": {
				Type:          schema.TypeString,
				Optional:      true,
				Description:   "The fully qualified domain name or IP address of the SDDC Manager instance.",
				ConflictsWith: []string{"installer_username", "installer_password", "installer_host"},
				RequiredWith:  []string{"sddc_manager_username", "sddc_manager_password"},
				DefaultFunc:   schema.EnvDefaultFunc(constants.VcfTestUrl, nil),
			},
			"installer_username": {
				Type:          schema.TypeString,
				Optional:      true,
				Description:   "The username to authenticate to the installer.",
				ConflictsWith: []string{"sddc_manager_username", "sddc_manager_password", "sddc_manager_host"},
				RequiredWith:  []string{"installer_password", "installer_host"},
				DefaultFunc:   schema.EnvDefaultFunc(constants.InstallerTestUsername, nil),
			},
			"installer_password": {
				Type:          schema.TypeString,
				Optional:      true,
				Description:   "The password to authenticate to the installer.",
				ConflictsWith: []string{"sddc_manager_username", "sddc_manager_password", "sddc_manager_host"},
				RequiredWith:  []string{"installer_username", "installer_host"},
				DefaultFunc:   schema.EnvDefaultFunc(constants.InstallerTestPassword, nil),
			},
			"installer_host": {
				Type:          schema.TypeString,
				Optional:      true,
				Description:   "The fully qualified domain name or IP address of the installer.",
				ConflictsWith: []string{"sddc_manager_username", "sddc_manager_password", "sddc_manager_host"},
				RequiredWith:  []string{"installer_username", "installer_password"},
				DefaultFunc:   schema.EnvDefaultFunc(constants.InstallerTestUrl, nil),
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
			"vcf_credentials":  DataSourceCredentials(),
			"vcf_domain":       DataSourceDomain(),
			"vcf_host":         DataSourceHost(),
			"vcf_network_pool": DataSourceNetworkPool(),
			"vcf_certificate":  DataSourceCertificate(),
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
	installerUsername, isInstallerUsernameSet := data.GetOk("installer_username")
	allowUnverifiedTLS := data.Get("allow_unverified_tls")

	if !isVcfUsernameSet && !isInstallerUsernameSet {
		return nil, diag.Errorf("Either SDDC Manager or Installer configuration must be provided.")
	}

	if isVcfUsernameSet {
		password, isSetPassword := data.GetOk("sddc_manager_password")
		hostName, isSetHost := data.GetOk("sddc_manager_host")
		if !isSetPassword || !isSetHost {
			return nil, diag.Errorf("SDDC Manager username, password, and host must be provided.")
		}
		sddcManagerClient := api_client.NewSddcManagerClient(
			sddcManagerUsername.(string),
			password.(string),
			hostName.(string),
			version.ProviderVersion,
			allowUnverifiedTLS.(bool))
		err := sddcManagerClient.Connect()
		if err != nil {
			return nil, diag.FromErr(err)
		}
		return sddcManagerClient, nil
	}

	if isInstallerUsernameSet {
		password, isSetPassword := data.GetOk("installer_password")
		hostName, isSetHost := data.GetOk("installer_host")
		if !isSetPassword || !isSetHost {
			return nil, diag.Errorf("Installer username, password, and host must be provided.")
		}
		installerClient := api_client.NewInstallerClient(installerUsername.(string), password.(string),
			hostName.(string), allowUnverifiedTLS.(bool))
		err := installerClient.Connect()
		if err != nil {
			return nil, diag.FromErr(err)
		}
		return installerClient, nil
	}

	return nil, diag.Errorf("Failed to configure the provider. Please check the provider configuration settings.")
}
