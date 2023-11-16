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
	"github.com/vmware/terraform-provider-vcf/internal/certificates"
	"github.com/vmware/terraform-provider-vcf/internal/constants"
	"github.com/vmware/terraform-provider-vcf/internal/resource_utils"
	certificatesSdk "github.com/vmware/vcf-sdk-go/client/certificates"
	"github.com/vmware/vcf-sdk-go/models"
	"strconv"
	"time"
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
				Type:     schema.TypeString,
				Required: true,
				// TODO when migrating to 5.x.x support check if these are still accurate
				Description:  "Resources for which the CSRs are to be generated. One among: SDDC_MANAGER, VCENTER, NSX_MANAGER, NSXT_MANAGER, VROPS, VRSLCM, VXRAIL_MANAGER",
				ValidateFunc: validation.StringInSlice([]string{"SDDC_MANAGER", "VCENTER", "NSX_MANAGER", "NSXT_MANAGER", "VROPS", "VRSLCM", "VXRAIL_MANAGER"}, false),
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

	resourceFqdn, err := certificates.GetFqdnOfResourceTypeInDomain(ctx, domainId, resourceType, apiClient)
	if err != nil {
		return diag.FromErr(err)
	}
	if resourceFqdn == nil {
		return diag.FromErr(fmt.Errorf("could not determine FQDN for resourceType %s in domain %s", resourceType, domainId))
	}

	country := data.Get("country").(string)
	email := data.Get("email").(string)
	keySize := strconv.Itoa(data.Get("key_size").(int))
	locality := data.Get("locality").(string)
	organization := data.Get("organization").(string)
	organizationUnit := data.Get("organization_unit").(string)
	state := data.Get("state").(string)

	csrGenerationSpec := &models.CSRGenerationSpec{
		Country:          &country,
		Email:            email,
		KeyAlgorithm:     resource_utils.ToStringPointer("RSA"),
		KeySize:          &keySize,
		Locality:         &locality,
		Organization:     &organization,
		OrganizationUnit: &organizationUnit,
		State:            &state,
	}

	csrsGenerationSpec := &models.CSRSGenerationSpec{
		CSRGenerationSpec: csrGenerationSpec,
		Resources: []*models.Resource{
			{
				Fqdn: *resourceFqdn,
				Type: &resourceType,
			},
		},
	}

	generateCsrParams := certificatesSdk.NewGeneratesCSRsParamsWithContext(ctx).
		WithTimeout(constants.DefaultVcfApiCallTimeout).
		WithDomainName(domainId).
		WithCSRSGenerationSpec(csrsGenerationSpec)

	var taskId string
	_, task, err := apiClient.Certificates.GeneratesCSRs(generateCsrParams)
	if err != nil {
		return diag.FromErr(err)
	}
	if task != nil {
		taskId = task.Payload.ID
	}
	err = vcfClient.WaitForTaskComplete(ctx, taskId, true)
	if err != nil {
		return diag.FromErr(err)
	}
	data.SetId("csr:" + domainId + ":" + resourceType + ":" + taskId)

	return resourceCsrRead(ctx, data, meta)
}

func resourceCsrRead(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	apiClient := meta.(*api_client.SddcManagerClient).ApiClient

	domainId := data.Get("domain_id").(string)
	getCsrsParams := certificatesSdk.NewGetCSRsParamsWithContext(ctx).
		WithTimeout(constants.DefaultVcfApiCallTimeout).
		WithDomainName(domainId)

	getCsrResponse, err := apiClient.Certificates.GetCSRs(getCsrsParams)
	if err != nil {
		return diag.FromErr(err)
	}
	resourceType := data.Get("resource").(string)
	resourceFqdn, err := certificates.GetFqdnOfResourceTypeInDomain(ctx, domainId, resourceType, apiClient)
	if err != nil {
		return diag.FromErr(err)
	}

	csr := getCsrByResourceFqdn(*resourceFqdn, getCsrResponse.Payload.Elements)
	flattenedCsr := certificates.FlattenCsr(csr)
	_ = data.Set("csr", []interface{}{flattenedCsr})

	return nil
}

func resourceCsrUpdate(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return resourceCsrCreate(ctx, data, meta)
}

func resourceCsrDelete(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

// getCsrByResourceFqdn SDDC Manager API doesn't return CSR resource type, just FQDN.
func getCsrByResourceFqdn(resourceFqdn string, csrs []*models.CSR) *models.CSR {
	if len(resourceFqdn) < 1 || len(csrs) < 1 {
		return nil
	}
	for _, csr := range csrs {
		if len(csr.Resource.Fqdn) > 0 && resourceFqdn == csr.Resource.Fqdn {
			return csr
		}
	}
	return nil
}
