// © Broadcom. All Rights Reserved.
// The term “Broadcom” refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/vmware/vcf-sdk-go/vcf"

	"github.com/vmware/terraform-provider-vcf/internal/api_client"
	"github.com/vmware/terraform-provider-vcf/internal/certificates"
	"github.com/vmware/terraform-provider-vcf/internal/constants"
)

func ResourceCsr() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceCsrCreate,
		ReadContext:   resourceCsrRead,
		UpdateContext: resourceCsrUpdate,
		DeleteContext: resourceCsrDelete,
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(20 * time.Minute),
			Read:   schema.DefaultTimeout(20 * time.Minute),
			Update: schema.DefaultTimeout(20 * time.Minute),
			Delete: schema.DefaultTimeout(20 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"domain_id": {
				Type:         schema.TypeString,
				Description:  "Domain Id or Name for which the CSRs should be generated",
				Required:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},
			"country": {
				Type:         schema.TypeString,
				Description:  "ISO 3166 country code where company is legally registered",
				Required:     true,
				ValidateFunc: validation.StringInSlice(constants.GetIso3166CountryCodes(), false),
			},
			"email": {
				Type:         schema.TypeString,
				Description:  "Contact email address",
				Required:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},
			"key_size": {
				Type:         schema.TypeInt,
				Description:  "Certificate public key size. One among: 2048, 3072, 4096",
				Required:     true,
				ValidateFunc: validation.IntInSlice([]int{2048, 3072, 4096}),
			},
			"locality": {
				Type:         schema.TypeString,
				Description:  "The city or locality where company is legally registered",
				Required:     true,
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
			"resource": {
				Type:         schema.TypeString,
				Required:     true,
				Description:  "Resources for which the CSRs are to be generated. One among: SDDC_MANAGER, PSC, VCENTER, NSX_MANAGER, NSXT_MANAGER, VROPS, VRSLCM, VXRAIL_MANAGER",
				ValidateFunc: validation.StringInSlice([]string{"SDDC_MANAGER", "PSC", "VCENTER", "NSX_MANAGER", "NSXT_MANAGER", "VROPS", "VRSLCM", "VXRAIL_MANAGER"}, false),
			},
			"fqdn": {
				Type:         schema.TypeString,
				Required:     true,
				Description:  "FQDN of the resource",
				ValidateFunc: validation.NoZeroValues,
			},
			"csr": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Resulting CSR",
				Elem:        certificates.CsrSchema(),
			},
		},
	}
}

func resourceCsrCreate(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	vcfClient := meta.(*api_client.SddcManagerClient)
	apiClient := vcfClient.ApiClient

	domainId := data.Get("domain_id").(string)
	resourceType := data.Get("resource").(string)
	resourceFqdn := data.Get("fqdn").(string)

	country := data.Get("country").(string)
	email := data.Get("email").(string)
	keySize := strconv.Itoa(data.Get("key_size").(int))
	locality := data.Get("locality").(string)
	organization := data.Get("organization").(string)
	organizationUnit := data.Get("organization_unit").(string)
	state := data.Get("state").(string)

	csrGenerationSpec := vcf.CsrGenerationSpec{
		Country:          country,
		Email:            &email,
		KeyAlgorithm:     "RSA",
		KeySize:          keySize,
		Locality:         locality,
		Organization:     organization,
		OrganizationUnit: organizationUnit,
		State:            state,
	}

	csrsGenerationSpec := vcf.CsrsGenerationSpec{
		CsrGenerationSpec: csrGenerationSpec,
		Resources: &[]vcf.Resource{
			{
				Fqdn: &resourceFqdn,
				Type: resourceType,
			},
		},
	}

	res, err := apiClient.GeneratesCSRsWithResponse(ctx, domainId, csrsGenerationSpec)
	if err != nil {
		return diag.FromErr(err)
	}
	task, vcfErr := api_client.GetResponseAs[vcf.Task](res.Body)
	if vcfErr != nil {
		api_client.LogError(vcfErr)
		return diag.FromErr(errors.New(*vcfErr.Message))
	}
	err = vcfClient.WaitForTaskComplete(ctx, *task.Id, true)
	if err != nil {
		return diag.FromErr(err)
	}
	data.SetId(fmt.Sprintf("csr:%s:%s:%s:%s", domainId, resourceType, resourceFqdn, *task.Id))

	getCsrResponse, err := apiClient.GetCSRsWithResponse(ctx, domainId)
	if err != nil {
		return diag.FromErr(err)
	}
	page, vcfErr := api_client.GetResponseAs[vcf.PageOfCsr](getCsrResponse.Body)
	if vcfErr != nil {
		api_client.LogError(vcfErr)
		return diag.FromErr(errors.New(*vcfErr.Message))
	}

	csr := getCsrByResourceFqdn(resourceFqdn, page.Elements)
	flattenedCsr := certificates.FlattenCsr(csr)
	_ = data.Set("csr", []interface{}{flattenedCsr})

	return resourceCsrRead(ctx, data, meta)
}

func resourceCsrRead(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return nil
}

func resourceCsrUpdate(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return resourceCsrCreate(ctx, data, meta)
}

func resourceCsrDelete(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

// getCsrByResourceFqdn SDDC Manager API doesn't return CSR resource type, just FQDN.
func getCsrByResourceFqdn(resourceFqdn string, csrs *[]vcf.Csr) *vcf.Csr {
	if len(resourceFqdn) < 1 || csrs != nil && len(*csrs) < 1 {
		return nil
	}
	for _, csr := range *csrs {
		if csr.Resource != nil && csr.Resource.Fqdn != nil && resourceFqdn == *csr.Resource.Fqdn {
			return &csr
		}
	}
	return nil
}
