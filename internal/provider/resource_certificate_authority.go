// © Broadcom. All Rights Reserved.
// The term “Broadcom” refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/vmware/vcf-sdk-go/vcf"

	"github.com/vmware/terraform-provider-vcf/internal/api_client"
	"github.com/vmware/terraform-provider-vcf/internal/constants"
	validationUtils "github.com/vmware/terraform-provider-vcf/internal/validation"
)

func ResourceCertificateAuthority() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceCertificateAuthorityCreate,
		ReadContext:   resourceCertificateAuthorityRead,
		DeleteContext: resourceCertificateAuthorityDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(20 * time.Minute),
			Read:   schema.DefaultTimeout(20 * time.Minute),
			Update: schema.DefaultTimeout(20 * time.Minute),
			Delete: schema.DefaultTimeout(20 * time.Minute),
		},
		CustomizeDiff: validateRequiredAttributesForCertificateAuthority,
		Schema: map[string]*schema.Schema{
			"type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Certificate authority type. \"id\" has the same value. Microsoft or OpenSSL",
			},
			"microsoft": {
				Type:          schema.TypeList,
				MaxItems:      1,
				Optional:      true,
				ForceNew:      true,
				Description:   "Configuration describing Microsoft CA server",
				ConflictsWith: []string{"open_ssl"},
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"server_url": {
							Type:         schema.TypeString,
							Required:     true,
							Description:  "Microsoft CA server URL",
							ValidateFunc: validation.StringIsNotEmpty,
						},
						"template_name": {
							Type:         schema.TypeString,
							Required:     true,
							Description:  "Microsoft CA server template name",
							ValidateFunc: validation.StringIsNotEmpty,
						},
						"username": {
							Type:         schema.TypeString,
							Required:     true,
							Description:  "Microsoft CA server username",
							ValidateFunc: validation.StringIsNotEmpty,
						},
						"secret": {
							Type:         schema.TypeString,
							Required:     true,
							Sensitive:    true,
							Description:  "Microsoft CA server password",
							ValidateFunc: validation.StringIsNotEmpty,
						},
					},
				},
			},
			"open_ssl": {
				Type:          schema.TypeList,
				MaxItems:      1,
				Optional:      true,
				ForceNew:      true,
				Description:   "Configuration describing OpenSSL CA server",
				ConflictsWith: []string{"microsoft"},
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"common_name": {
							Type:         schema.TypeString,
							Required:     true,
							Description:  "OpenSSL CA domain name",
							ValidateFunc: validation.StringIsNotEmpty,
						},
						"country": {
							Type:         schema.TypeString,
							Required:     true,
							Description:  "ISO 3166 country code where company is legally registered",
							ValidateFunc: validation.StringInSlice(constants.GetIso3166CountryCodes(), false),
						},
						"locality": {
							Type:         schema.TypeString,
							Required:     true,
							Description:  "The city or locality where company is legally registered",
							ValidateFunc: validation.StringIsNotEmpty,
						},
						"organization": {
							Type:         schema.TypeString,
							Required:     true,
							Description:  "The name under which your company is known. The listed organization must be the legal registrant of the domain name in the certificate request.",
							ValidateFunc: validation.StringIsNotEmpty,
						},
						"organization_unit": {
							Type:         schema.TypeString,
							Required:     true,
							Description:  "Organization with which the certificate is associated",
							ValidateFunc: validation.StringIsNotEmpty,
						},
						"state": {
							Type:         schema.TypeString,
							Required:     true,
							Description:  "Full name (do not abbreviate) of the state, province, region, or territory where your company is legally registered.",
							ValidateFunc: validation.StringIsNotEmpty,
						},
					}},
			},
		},
	}
}

func validateRequiredAttributesForCertificateAuthority(_ context.Context, diff *schema.ResourceDiff, _ interface{}) error {
	microsoftConfig := diff.Get("microsoft")
	openSslConfig := diff.Get("open_ssl")

	if validationUtils.IsEmpty(microsoftConfig) && validationUtils.IsEmpty(openSslConfig) {
		return fmt.Errorf("one of \"microsoft\" or \"open_ssl\" configuration has to be provided")
	}

	return nil
}

func resourceCertificateAuthorityCreate(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	apiClient := meta.(*api_client.SddcManagerClient).ApiClient

	certificateAuthorityCreationSpec := getCertificateAuthorityCreationSpec(data)
	if certificateAuthorityCreationSpec == nil {
		return diag.FromErr(fmt.Errorf("certificateAuthorityCreationSpec is empty, there was an error converting schema attributes to SDK spec"))
	}

	_, err := apiClient.CreateCertificateAuthorityWithResponse(ctx, *certificateAuthorityCreationSpec)
	if err != nil {
		return diag.FromErr(err)
	}
	data.SetId(*getCaType(data))

	return resourceCertificateAuthorityRead(ctx, data, meta)
}

func resourceCertificateAuthorityRead(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	apiClient := meta.(*api_client.SddcManagerClient).ApiClient

	authorityId := data.Id()

	authorityResponse, err := apiClient.GetCertificateAuthorityByIdWithResponse(ctx, authorityId)
	if err != nil {
		return diag.FromErr(err)
	}

	if authorityResponse.StatusCode() != 200 {
		vcfError := api_client.GetError(authorityResponse.Body)
		api_client.LogError(vcfError)
		return diag.FromErr(errors.New(*vcfError.Message))
	}

	// The ID doubles as type as per API
	_ = data.Set("type", authorityId)
	certificateAuthority := authorityResponse.JSON200
	if authorityId == "Microsoft" {
		microsoftConfigAttribute, microsoftConfigExists := data.GetOk("microsoft")
		var microsoftConfigRaw []interface{}
		if microsoftConfigExists {
			microsoftConfigRaw = microsoftConfigAttribute.([]interface{})
		} else {
			microsoftConfigRaw = *new([]interface{})
			microsoftConfigRaw = append(microsoftConfigRaw, make(map[string]interface{}))
		}
		microsoftConfig := microsoftConfigRaw[0].(map[string]interface{})
		microsoftConfig["server_url"] = certificateAuthority.ServerUrl
		microsoftConfig["template_name"] = certificateAuthority.TemplateName
		microsoftConfig["username"] = certificateAuthority.Username
		_ = data.Set("microsoft", microsoftConfigRaw)
	}
	if authorityId == "OpenSSL" {
		openSslConfigAttribute, openSslConfigExists := data.GetOk("open_ssl")
		var openSslConfigRaw []interface{}
		if openSslConfigExists {
			openSslConfigRaw = openSslConfigAttribute.([]interface{})
		} else {
			openSslConfigRaw = *new([]interface{})
			openSslConfigRaw = append(openSslConfigRaw, make(map[string]interface{}))
		}
		openSslConfig := openSslConfigRaw[0].(map[string]interface{})
		openSslConfig["common_name"] = certificateAuthority.CommonName
		openSslConfig["country"] = certificateAuthority.Country
		openSslConfig["locality"] = certificateAuthority.Locality
		openSslConfig["organization"] = certificateAuthority.Organization
		openSslConfig["organization_unit"] = certificateAuthority.OrganizationUnit
		openSslConfig["state"] = certificateAuthority.State
		_ = data.Set("open_ssl", openSslConfigRaw)
	}

	return nil
}

func resourceCertificateAuthorityDelete(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	apiClient := meta.(*api_client.SddcManagerClient).ApiClient

	caType := getCaType(data)
	if caType == nil {
		return diag.FromErr(fmt.Errorf("error deleting Certificate Authority: could not determine CA type"))
	}

	_, err := apiClient.RemoveCertificateAuthorityWithResponse(ctx, *caType)
	if err != nil {
		return diag.FromErr(err)
	}

	data.SetId("")
	return nil
}

func getCertificateAuthorityCreationSpec(data *schema.ResourceData) *vcf.CertificateAuthorityCreationSpec {
	certificateAuthorityCreationSpec := &vcf.CertificateAuthorityCreationSpec{}
	microsoftConfig := data.Get("microsoft").([]interface{})
	openSslConfig := data.Get("open_ssl").([]interface{})

	caType := getCaType(data)
	if caType == nil {
		return nil
	}

	if *caType == "Microsoft" {
		microsoftConfigMap := microsoftConfig[0].(map[string]interface{})
		serverUrl := microsoftConfigMap["server_url"].(string)
		templateName := microsoftConfigMap["template_name"].(string)
		username := microsoftConfigMap["username"].(string)
		secret := microsoftConfigMap["secret"].(string)
		certificateAuthorityCreationSpec.MicrosoftCertificateAuthoritySpec = &vcf.MicrosoftCertificateAuthoritySpec{
			ServerUrl:    serverUrl,
			TemplateName: templateName,
			Username:     username,
			Secret:       secret,
		}
	}
	if *caType == "OpenSSL" {
		openSslConfigMap := openSslConfig[0].(map[string]interface{})
		commonName := openSslConfigMap["common_name"].(string)
		country := openSslConfigMap["country"].(string)
		locality := openSslConfigMap["locality"].(string)
		organization := openSslConfigMap["organization"].(string)
		organizationUnit := openSslConfigMap["organization_unit"].(string)
		state := openSslConfigMap["state"].(string)
		certificateAuthorityCreationSpec.OpenSSLCertificateAuthoritySpec = &vcf.OpenSSLCertificateAuthoritySpec{
			CommonName:       commonName,
			Country:          country,
			Locality:         locality,
			Organization:     organization,
			OrganizationUnit: organizationUnit,
			State:            state,
		}
	}
	return certificateAuthorityCreationSpec
}

func getCaType(data *schema.ResourceData) *string {
	var caType string
	if !validationUtils.IsEmpty(data.Get("microsoft")) {
		caType = "Microsoft"
		return &caType
	}
	if !validationUtils.IsEmpty(data.Get("open_ssl")) {
		caType = "OpenSSL"
		return &caType
	} else {
		return nil
	}
}
