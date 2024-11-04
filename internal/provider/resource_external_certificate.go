// © Broadcom. All Rights Reserved.
// The term “Broadcom” refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/vmware/vcf-sdk-go/vcf"

	"github.com/vmware/terraform-provider-vcf/internal/api_client"
	"github.com/vmware/terraform-provider-vcf/internal/certificates"
	validationutils "github.com/vmware/terraform-provider-vcf/internal/validation"
)

func ResourceExternalCertificate() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceResourceExternalCertificateCreate,
		ReadContext:   resourceResourceExternalCertificateRead,
		UpdateContext: resourceResourceExternalCertificateUpdate,
		DeleteContext: resourceResourceExternalCertificateDelete,
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(50 * time.Minute),
			Read:   schema.DefaultTimeout(20 * time.Minute),
			Update: schema.DefaultTimeout(50 * time.Minute),
			Delete: schema.DefaultTimeout(20 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"csr_id": {
				Type:         schema.TypeString,
				Required:     true,
				Description:  "The ID of the CSR generated for a resource. A generated CSR is required for certificate replacement.",
				ValidateFunc: validation.StringIsNotEmpty,
			},
			"resource_certificate": {
				Type:          schema.TypeString,
				Optional:      true,
				Description:   "Resource Certificate",
				RequiredWith:  []string{"ca_certificate"},
				ConflictsWith: []string{"certificate_chain"},
				ValidateFunc:  validation.StringIsNotEmpty,
			},
			"ca_certificate": {
				Type:          schema.TypeString,
				Optional:      true,
				Description:   "Certificate of the CA issuing the replacement certificate",
				RequiredWith:  []string{"resource_certificate"},
				ConflictsWith: []string{"certificate_chain"},
				ValidateFunc:  validation.StringIsNotEmpty,
			},
			"certificate_chain": {
				Type:          schema.TypeString,
				Optional:      true,
				Description:   "Certificate Chain",
				ConflictsWith: []string{"resource_certificate", "ca_certificate"},
				ValidateFunc:  validation.StringIsNotEmpty,
			},

			"certificate": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "The resulting certificate details",
				Elem:        certificates.CertificateSchema(),
			},
		},
	}
}

func resourceResourceExternalCertificateCreate(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	vcfClient := meta.(*api_client.SddcManagerClient)
	apiClient := vcfClient.ApiClient

	csrID := data.Get("csr_id").(string)
	csrIdComponents := strings.Split(csrID, ":")
	if len(csrIdComponents) != 5 {
		return diag.FromErr(fmt.Errorf("CSR ID invalid"))
	}

	domainID := csrIdComponents[1]
	resourceType := csrIdComponents[2]
	resourceFqdn := csrIdComponents[3]

	caCertificate := data.Get("ca_certificate").(string)
	certificateChain := data.Get("certificate_chain").(string)
	resourceCertificate := data.Get("resource_certificate").(string)

	var resourceCertificateSpec vcf.ResourceCertificateSpec

	if !validationutils.IsEmpty(resourceCertificate) && !validationutils.IsEmpty(caCertificate) {
		resourceCertificateSpec = vcf.ResourceCertificateSpec{
			ResourceFqdn:        &resourceFqdn,
			CaCertificate:       &caCertificate,
			ResourceCertificate: &resourceCertificate,
		}
	} else if !validationutils.IsEmpty(certificateChain) {
		resourceCertificateSpec = vcf.ResourceCertificateSpec{
			ResourceFqdn:     &resourceFqdn,
			CertificateChain: &certificateChain,
		}
	} else {
		return diag.FromErr(fmt.Errorf("no certificate_chain or (ca_certificate, resource_certificate) defined"))
	}

	resourceCertificateSpecs := []vcf.ResourceCertificateSpec{resourceCertificateSpec}

	diags := certificates.ValidateResourceCertificates(ctx, apiClient, domainID, resourceCertificateSpecs)
	if diags != nil {
		return diags
	}

	responseAcc, err := apiClient.ReplaceResourceCertificatesWithResponse(ctx, domainID, resourceCertificateSpecs)
	if err != nil {
		return diag.FromErr(err)
	}
	task, vcfErr := api_client.GetResponseAs[vcf.Task](responseAcc.Body, responseAcc.StatusCode())
	if vcfErr != nil {
		api_client.LogError(vcfErr)
		return diag.FromErr(errors.New(*vcfErr.Message))
	}
	err = vcfClient.WaitForTaskComplete(ctx, *task.Id, true)
	if err != nil {
		return diag.FromErr(err)
	}
	data.SetId("ext_cert:" + domainID + ":" + resourceType + ":" + *task.Id)

	return resourceResourceExternalCertificateRead(ctx, data, meta)
}

func resourceResourceExternalCertificateRead(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	apiClient := meta.(*api_client.SddcManagerClient).ApiClient

	csrID := data.Get("csr_id").(string)
	csrIdComponents := strings.Split(csrID, ":")
	if len(csrIdComponents) != 5 {
		return diag.FromErr(fmt.Errorf("CSR ID invalid"))
	}

	domainID := csrIdComponents[1]
	resourceType := csrIdComponents[2]

	cert, err := certificates.ReadCertificate(ctx, apiClient, domainID, resourceType)
	if err != nil {
		return diag.FromErr(err)
	}

	flattenedCert := certificates.FlattenCertificate(*cert)
	_ = data.Set("certificate", []interface{}{flattenedCert})

	return nil
}

func resourceResourceExternalCertificateUpdate(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return resourceResourceExternalCertificateCreate(ctx, data, meta)
}

func resourceResourceExternalCertificateDelete(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}
