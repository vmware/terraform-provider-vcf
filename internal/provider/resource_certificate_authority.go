/*
 *  Copyright 2023 VMware, Inc.
 *    SPDX-License-Identifier: MPL-2.0
 */

package provider

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/vmware/terraform-provider-vcf/internal/api_client"
	"github.com/vmware/terraform-provider-vcf/internal/constants"
	validationUtils "github.com/vmware/terraform-provider-vcf/internal/validation"
	"github.com/vmware/vcf-sdk-go/client/certificates"
	"github.com/vmware/vcf-sdk-go/models"
	"time"
)

func ResourceCertificateAuthority() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceCertificateAuthorityCreate,
		ReadContext:   resourceCertificateAuthorityRead,
		UpdateContext: resourceCertificateAuthorityUpdate,
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
				Description:   "",
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
				Description:   "",
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
							Type:        schema.TypeString,
							Required:    true,
							Description: "ISO 3166 country code where company is legally registered",
							ValidateFunc: validation.StringInSlice([]string{"US", "CA", "AX", "AD", "AE", "AF", "AG", "AI", "AL", "AM", "AN", "AO", "AQ", "AR", "AS", "AT", "AU",
								"AW", "AZ", "BA", "BB", "BD", "BE", "BF", "BG", "BH", "BI", "BJ", "BM", "BN", "BO", "BR", "BS", "BT", "BV", "BW", "BZ", "CA", "CC", "CF", "CH", "CI", "CK",
								"CL", "CM", "CN", "CO", "CR", "CS", "CV", "CX", "CY", "CZ", "DE", "DJ", "DK", "DM", "DO", "DZ", "EC", "EE", "EG", "EH", "ER", "ES", "ET", "FI", "FJ", "FK",
								"FM", "FO", "FR", "FX", "GA", "GB", "GD", "GE", "GF", "GG", "GH", "GI", "GL", "GM", "GN", "GP", "GQ", "GR", "GS", "GT", "GU", "GW", "GY", "HK", "HM", "HN",
								"HR", "HT", "HU", "ID", "IE", "IL", "IM", "IN", "IO", "IS", "IT", "JE", "JM", "JO", "JP", "KE", "KG", "KH", "KI", "KM", "KN", "KR", "KW", "KY", "KZ", "LA",
								"LC", "LI", "LK", "LS", "LT", "LU", "LV", "LY", "MA", "MC", "MD", "ME", "MG", "MH", "MK", "ML", "MM", "MN", "MO", "MP", "MQ", "MR", "MS", "MT", "MU", "MV",
								"MW", "MX", "MY", "MZ", "NA", "NC", "NE", "NF", "NG", "NI", "NL", "NO", "NP", "NR", "NT", "NU", "NZ", "OM", "PA", "PE", "PF", "PG", "PH", "PK", "PL", "PM",
								"PN", "PR", "PS", "PT", "PW", "PY", "QA", "RE", "RO", "RS", "RU", "RW", "SA", "SB", "SC", "SE", "SG", "SH", "SI", "SJ", "SK", "SL", "SM", "SN", "SR", "ST",
								"SU", "SV", "SZ", "TC", "TD", "TF", "TG", "TH", "TJ", "TK", "TM", "TN", "TO", "TP", "TR", "TT", "TV", "TW", "TZ", "UA", "UG", "UM", "US", "UY", "UZ", "VA",
								"VC", "VE", "VG", "VI", "VN", "VU", "WF", "WS", "YE", "YT", "ZA", "ZM", "COM", "EDU", "GOV", "INT", "MIL", "NET", "ORG", "ARPA"}, false),
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
	microsoftConfig := diff.Get("microsoft").(string)
	openSslConfig := diff.Get("open_ssl").(string)

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

	createCertificateAuthorityParams := certificates.NewCreateCertificateAuthorityParamsWithContext(ctx).
		WithTimeout(constants.DefaultVcfApiCallTimeout).
		WithCertificateAuthorityCreationSpec(certificateAuthorityCreationSpec)

	_, err := apiClient.Certificates.CreateCertificateAuthority(createCertificateAuthorityParams)
	if err != nil {
		return diag.FromErr(err)
	}
	data.SetId(*getCaType(data))

	return resourceCertificateAuthorityRead(ctx, data, meta)
}

func resourceCertificateAuthorityRead(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	apiClient := meta.(*api_client.SddcManagerClient).ApiClient

	authorityId := data.Id()
	getAuthorityParams := certificates.NewGetCertificateAuthorityByIDParamsWithContext(ctx).
		WithID(authorityId).WithTimeout(constants.DefaultVcfApiCallTimeout)

	authorityResponse, err := apiClient.Certificates.GetCertificateAuthorityByID(getAuthorityParams)
	if err != nil {
		return diag.FromErr(err)
	}

	// The ID doubles as type as per API
	_ = data.Set("type", authorityId)
	certificateAuthority := authorityResponse.Payload
	if authorityId == "Microsoft" {
		microsoftConfigRaw := *new([]interface{})
		microsoftConfigRaw = append(microsoftConfigRaw, make(map[string]interface{}))
		microsoftConfig := microsoftConfigRaw[0].(map[string]interface{})
		microsoftConfig["server_url"] = certificateAuthority.ServerURL
		microsoftConfig["template_name"] = certificateAuthority.TemplateName
		microsoftConfig["username"] = certificateAuthority.Username
		_ = data.Set("microsoft", microsoftConfigRaw)
	}
	if authorityId == "OpenSSL" {
		openSslConfigRaw := *new([]interface{})
		openSslConfigRaw = append(openSslConfigRaw, make(map[string]interface{}))
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

func resourceCertificateAuthorityUpdate(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	apiClient := meta.(*api_client.SddcManagerClient).ApiClient

	certificateAuthorityCreationSpec := getCertificateAuthorityCreationSpec(data)

	configureCertificateAuthorityParams := certificates.NewConfigureCertificateAuthorityParamsWithContext(ctx).
		WithTimeout(constants.DefaultVcfApiCallTimeout).
		WithCertificateAuthoritySpec(certificateAuthorityCreationSpec)

	_, err := apiClient.Certificates.ConfigureCertificateAuthority(configureCertificateAuthorityParams)
	if err != nil {
		return diag.FromErr(err)
	}

	return resourceCertificateAuthorityRead(ctx, data, meta)
}

func resourceCertificateAuthorityDelete(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	apiClient := meta.(*api_client.SddcManagerClient).ApiClient

	caType := getCaType(data)
	if caType == nil {
		return diag.FromErr(fmt.Errorf("error deleting Certificate Authority: could not determine CA type"))
	}
	deleteCaConfigurationParams := certificates.NewDeleteCaConfigurationParamsWithContext(ctx).
		WithTimeout(constants.DefaultVcfApiCallTimeout).WithCaType(*caType)

	_, _, err := apiClient.Certificates.DeleteCaConfiguration(deleteCaConfigurationParams)
	if err != nil {
		return diag.FromErr(err)
	}

	data.SetId("")
	return nil
}

func getCertificateAuthorityCreationSpec(data *schema.ResourceData) *models.CertificateAuthorityCreationSpec {
	certificateAuthorityCreationSpec := &models.CertificateAuthorityCreationSpec{}
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
		certificateAuthorityCreationSpec.MicrosoftCertificateAuthoritySpec = &models.MicrosoftCertificateAuthoritySpec{
			ServerURL:    &serverUrl,
			TemplateName: &templateName,
			Username:     &username,
			Secret:       &secret,
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
		certificateAuthorityCreationSpec.OpenSSLCertificateAuthoritySpec = &models.OpenSSLCertificateAuthoritySpec{
			CommonName:       &commonName,
			Country:          &country,
			Locality:         &locality,
			Organization:     &organization,
			OrganizationUnit: &organizationUnit,
			State:            &state,
		}
	}
	return certificateAuthorityCreationSpec
}

func getCaType(data *schema.ResourceData) *string {
	var caType string
	if !validationUtils.IsEmpty(data.Get("microsoft")) {
		caType = "Microsoft"
	}
	if !validationUtils.IsEmpty(data.Get("open_ssl")) {
		caType = "OpenSSL"
	} else {
		return nil
	}
	return &caType
}
