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
	"github.com/vmware/vcf-sdk-go/client/certificates"
	"github.com/vmware/vcf-sdk-go/models"
	"time"
)

func getSupportedCertificateAuthorityTypes() []string {
	return []string{"Microsoft", "OpenSSL"}
}

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
			Create: schema.DefaultTimeout(12 * time.Hour),
		},
		CustomizeDiff: validateRequiredAttributesForCertificateAuthority,
		Schema: map[string]*schema.Schema{
			"type": {
				Type:         schema.TypeString,
				Required:     true,
				Description:  "Certificate authority type. Only one CA from each type allowed. One among: Microsoft, OpenSSL",
				ValidateFunc: validation.StringInSlice(getSupportedCertificateAuthorityTypes(), false),
			},
			// Microsoft CA
			"server_url": {
				Type:         schema.TypeString,
				Optional:     true,
				Description:  "Microsoft CA server URL",
				ValidateFunc: validation.StringIsNotEmpty,
			},
			"template_name": {
				Type:         schema.TypeString,
				Optional:     true,
				Description:  "Microsoft CA server template name",
				ValidateFunc: validation.StringIsNotEmpty,
			},
			"username": {
				Type:         schema.TypeString,
				Optional:     true,
				Description:  "Microsoft CA server username",
				ValidateFunc: validation.StringIsNotEmpty,
			},
			"secret": {
				Type:         schema.TypeString,
				Optional:     true,
				Sensitive:    true,
				Description:  "Microsoft CA server password",
				ValidateFunc: validation.StringIsNotEmpty,
			},
			// OpenSSL CA
			"common_name": {
				Type:         schema.TypeString,
				Optional:     true,
				Description:  "OpenSSL CA domain name",
				ValidateFunc: validation.StringIsNotEmpty,
			},
			"country": {
				Type:         schema.TypeString,
				Optional:     true,
				Description:  "ISO 3166 country code where company is legally registered",
				ValidateFunc: validation.StringIsNotEmpty,
			},
			"locality": {
				Type:         schema.TypeString,
				Optional:     true,
				Description:  "The city or locality where company is legally registered",
				ValidateFunc: validation.StringIsNotEmpty,
			},
			"organization": {
				Type:         schema.TypeString,
				Optional:     true,
				Description:  "The name under which company is legally registered",
				ValidateFunc: validation.StringIsNotEmpty,
			},
			"organization_unit": {
				Type:         schema.TypeString,
				Optional:     true,
				Description:  "Organization with which the certificate is associated",
				ValidateFunc: validation.StringIsNotEmpty,
			},
			"state": {
				Type:         schema.TypeString,
				Optional:     true,
				Description:  "The full name of the state where company is legally registered",
				ValidateFunc: validation.StringIsNotEmpty,
			},
		},
	}
}

func validateRequiredAttributesForCertificateAuthority(_ context.Context, diff *schema.ResourceDiff, _ interface{}) error {
	caType := diff.Get("type").(string)

	if caType == "Microsoft" {
		_, ok := diff.GetOk("server_url")
		if !ok {
			return fmt.Errorf("server_url required for Microsoft CA")
		}
		_, ok = diff.GetOk("template_name")
		if !ok {
			return fmt.Errorf("template_name required for Microsoft CA")
		}
		_, ok = diff.GetOk("username")
		if !ok {
			return fmt.Errorf("username required for Microsoft CA")
		}
		_, ok = diff.GetOk("secret")
		if !ok {
			return fmt.Errorf("secret required for Microsoft CA")
		}
	}

	if caType == "OpenSSL" {
		_, ok := diff.GetOk("common_name")
		if !ok {
			return fmt.Errorf("common_name required for OpenSSL CA")
		}
		_, ok = diff.GetOk("country")
		if !ok {
			return fmt.Errorf("country required for OpenSSL CA")
		}
		_, ok = diff.GetOk("locality")
		if !ok {
			return fmt.Errorf("locality required for OpenSSL CA")
		}
		_, ok = diff.GetOk("organization")
		if !ok {
			return fmt.Errorf("organization required for OpenSSL CA")
		}
		_, ok = diff.GetOk("organization_unit")
		if !ok {
			return fmt.Errorf("organization_unit required for OpenSSL CA")
		}
		_, ok = diff.GetOk("state")
		if !ok {
			return fmt.Errorf("state required for OpenSSL CA")
		}
	}

	return nil
}

func resourceCertificateAuthorityCreate(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	apiClient := meta.(*api_client.SddcManagerClient).ApiClient

	caType := data.Get("type").(string)
	certificateAuthorityCreationSpec := getCertificateAuthorityCreationSpec(data)

	createCertificateAuthorityParams := certificates.NewCreateCertificateAuthorityParamsWithContext(ctx).
		WithTimeout(constants.DefaultVcfApiCallTimeout).
		WithCertificateAuthorityCreationSpec(certificateAuthorityCreationSpec)

	_, err := apiClient.Certificates.CreateCertificateAuthority(createCertificateAuthorityParams)
	if err != nil {
		return diag.FromErr(err)
	}
	data.SetId(caType)

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
	certificateAuthority := authorityResponse.Payload
	if authorityId == "Microsoft" {
		_ = data.Set("server_url", certificateAuthority.ServerURL)
		_ = data.Set("template_name", certificateAuthority.TemplateName)
		_ = data.Set("username", certificateAuthority.Username)
	}
	if authorityId == "OpenSSL" {
		_ = data.Set("common_name", certificateAuthority.CommonName)
		_ = data.Set("country", certificateAuthority.Country)
		_ = data.Set("locality", certificateAuthority.Locality)
		_ = data.Set("organization", certificateAuthority.Organization)
		_ = data.Set("organization_unit", certificateAuthority.OrganizationUnit)
		_ = data.Set("state", certificateAuthority.State)
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

	caType := data.Get("type").(string)
	deleteCaConfigurationParams := certificates.NewDeleteCaConfigurationParamsWithContext(ctx).
		WithTimeout(constants.DefaultVcfApiCallTimeout).WithCaType(caType)

	_, _, err := apiClient.Certificates.DeleteCaConfiguration(deleteCaConfigurationParams)
	if err != nil {
		return diag.FromErr(err)
	}

	data.SetId("")
	return nil
}

func getCertificateAuthorityCreationSpec(data *schema.ResourceData) *models.CertificateAuthorityCreationSpec {
	caType := data.Get("type").(string)
	certificateAuthorityCreationSpec := &models.CertificateAuthorityCreationSpec{}
	if caType == "Microsoft" {
		serverUrl := data.Get("server_url").(string)
		templateName := data.Get("template_name").(string)
		username := data.Get("username").(string)
		secret := data.Get("secret").(string)
		certificateAuthorityCreationSpec.MicrosoftCertificateAuthoritySpec = &models.MicrosoftCertificateAuthoritySpec{
			ServerURL:    &serverUrl,
			TemplateName: &templateName,
			Username:     &username,
			Secret:       &secret,
		}
	}
	if caType == "OpenSSL" {
		commonName := data.Get("common_name").(string)
		country := data.Get("country").(string)
		locality := data.Get("locality").(string)
		organization := data.Get("organization").(string)
		organizationUnit := data.Get("organization_unit").(string)
		state := data.Get("state").(string)
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
