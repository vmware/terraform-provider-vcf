// Copyright 2023-2024 Broadcom. All Rights Reserved.
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
	"github.com/vmware/terraform-provider-vcf/internal/resource_utils"
)

func ResourceCertificate() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceResourceCertificateCreate,
		ReadContext:   resourceResourceCertificateRead,
		UpdateContext: resourceResourceCertificateUpdate,
		DeleteContext: resourceResourceCertificateDelete,
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
				Description:  "The ID of the CSR generated for a resource",
				ValidateFunc: validation.StringIsNotEmpty,
			},
			"ca_id": {
				Type:         schema.TypeString,
				Required:     true,
				Description:  "Certificate of the CA issuing the replacement certificate",
				ValidateFunc: validation.StringIsNotEmpty,
			},
			"certificate": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "The resulting Certificate details",
				Elem:        certificates.CertificateSchema(),
			},
		},
	}
}

func resourceResourceCertificateCreate(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
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
	caType := data.Get("ca_id").(string)

	err := certificates.GenerateCertificateForResource(ctx, vcfClient, &domainID, &resourceType, &resourceFqdn, &caType)
	if err != nil {
		return diag.FromErr(err)
	}

	certificateOperationSpec := &models.CertificateOperationSpec{
		OperationType: resource_utils.ToStringPointer("INSTALL"),
		Resources: []*models.Resource{{
			Fqdn: resourceFqdn,
			Type: &resourceType,
		}},
	}
	replaceCertificatesParams := certificatesSdk.NewReplaceCertificatesParamsWithContext(ctx).
		WithTimeout(constants.DefaultVcfApiCallTimeout).
		WithID(domainID)
	replaceCertificatesParams.SetCertificateOperationSpec(certificateOperationSpec)

	var taskId string
	responseOk, responseAccepted, err := apiClient.Certificates.ReplaceCertificates(replaceCertificatesParams)
	if err != nil {
		return diag.FromErr(err)
	}
	if responseOk != nil {
		taskId = responseOk.Payload.ID
	}
	if responseAccepted != nil {
		taskId = responseAccepted.Payload.ID
	}
	err = vcfClient.WaitForTaskComplete(ctx, taskId, true)
	if err != nil {
		return diag.FromErr(err)
	}
	data.SetId("cert:" + domainID + ":" + resourceType + ":" + taskId)

	return resourceResourceCertificateRead(ctx, data, meta)
}

func resourceResourceCertificateRead(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	vcfClient := meta.(*api_client.SddcManagerClient)
	apiClient := vcfClient.ApiClient

	csrID := data.Get("csr_id").(string)
	csrIdComponents := strings.Split(csrID, ":")
	if len(csrIdComponents) != 5 {
		return diag.FromErr(fmt.Errorf("CSR ID invalid"))
	}

	domainID := csrIdComponents[1]
	resourceFqdn := csrIdComponents[3]

	cert, err := certificates.GetCertificateForResourceInDomain(ctx, apiClient, domainID, resourceFqdn)
	if err != nil {
		return diag.FromErr(err)
	}

	flattenedCert := certificates.FlattenCertificate(cert)
	_ = data.Set("certificate", []interface{}{flattenedCert})

	return nil
}

func resourceResourceCertificateUpdate(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return resourceResourceCertificateCreate(ctx, data, meta)
}

func resourceResourceCertificateDelete(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}
