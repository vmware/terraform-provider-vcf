// Copyright 2023 Broadcom. All Rights Reserved.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	certificatesSdk "github.com/vmware/vcf-sdk-go/client/certificates"
	"github.com/vmware/vcf-sdk-go/models"

	"github.com/vmware/terraform-provider-vcf/internal/api_client"
	"github.com/vmware/terraform-provider-vcf/internal/certificates"
	"github.com/vmware/terraform-provider-vcf/internal/constants"
	validation_utils "github.com/vmware/terraform-provider-vcf/internal/validation"
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

	var resourceCertificateSpec *models.ResourceCertificateSpec

	if !validation_utils.IsEmpty(resourceCertificate) && !validation_utils.IsEmpty(caCertificate) {
		resourceCertificateSpec = &models.ResourceCertificateSpec{
			ResourceFqdn:        resourceFqdn,
			CaCertificate:       caCertificate,
			ResourceCertificate: resourceCertificate,
		}
	} else if !validation_utils.IsEmpty(certificateChain) {
		resourceCertificateSpec = &models.ResourceCertificateSpec{
			ResourceFqdn:     resourceFqdn,
			CertificateChain: certificateChain,
		}
	} else {
		return diag.FromErr(fmt.Errorf("no certificate_chain or (ca_certificate, resource_certificate) defined"))
	}

	resourceCertificateSpecs := []*models.ResourceCertificateSpec{resourceCertificateSpec}

	diags := certificates.ValidateResourceCertificates(ctx, apiClient, domainID, resourceCertificateSpecs)
	if diags != nil {
		return diags
	}

	replaceResourceCertificatesParams := certificatesSdk.NewReplaceResourceCertificatesParamsWithContext(ctx).
		WithTimeout(constants.DefaultVcfApiCallTimeout).
		WithID(domainID)
	replaceResourceCertificatesParams.SetResourceCertificateSpecs(resourceCertificateSpecs)

	var taskId string
	_, responseAcc, err := apiClient.Certificates.ReplaceResourceCertificates(replaceResourceCertificatesParams)
	if err != nil {
		return diag.FromErr(err)
	}
	if responseAcc != nil {
		taskId = responseAcc.Payload.ID
	}
	err = vcfClient.WaitForTaskComplete(ctx, taskId, true)
	if err != nil {
		return diag.FromErr(err)
	}
	data.SetId("ext_cert:" + domainID + ":" + resourceType + ":" + taskId)

	return resourceResourceExternalCertificateRead(ctx, data, meta)
}

func resourceResourceExternalCertificateRead(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	vcfClient := meta.(*api_client.SddcManagerClient)
	apiClient := vcfClient.ApiClient

	csrID := data.Get("csr_id").(string)
	csrIdComponents := strings.Split(csrID, ":")
	if len(csrIdComponents) != 5 {
		return diag.FromErr(fmt.Errorf("CSR ID invalid"))
	}

	domainID := csrIdComponents[1]
	resourceType := csrIdComponents[2]

	cert, err := certificates.GetCertificateForResourceInDomain(ctx, apiClient, domainID, resourceType)
	if err != nil {
		return diag.FromErr(err)
	}

	flattenedCert := certificates.FlattenCertificate(cert)
	_ = data.Set("certificate", []interface{}{flattenedCert})

	return nil
}

func resourceResourceExternalCertificateUpdate(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return resourceResourceExternalCertificateCreate(ctx, data, meta)
}

func resourceResourceExternalCertificateDelete(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}
