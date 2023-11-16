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
	certificatesSdk "github.com/vmware/vcf-sdk-go/client/certificates"
	"github.com/vmware/vcf-sdk-go/models"
	"time"
)

func ResourceExternalCertificate() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceResourceExternalCertificateCreate,
		ReadContext:   resourceResourceExternalCertificateRead,
		UpdateContext: resourceResourceExternalCertificateUpdate,
		DeleteContext: resourceResourceExternalCertificateDelete,
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(20 * time.Minute),
			Read:   schema.DefaultTimeout(20 * time.Minute),
			Update: schema.DefaultTimeout(20 * time.Minute),
			Delete: schema.DefaultTimeout(20 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"domain_id": {
				Type:         schema.TypeString,
				Description:  "Domain Id or Name for which the certificate should be rotated",
				Required:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},
			"resource": {
				Type:         schema.TypeString,
				Required:     true,
				Description:  "Resources for which the certificate should be rotated. One among: SDDC_MANAGER, VCENTER, NSX_MANAGER, NSXT_MANAGER, VROPS, VRSLCM, VXRAIL_MANAGER",
				ValidateFunc: validation.StringInSlice([]string{"SDDC_MANAGER", "VCENTER", "NSX_MANAGER", "NSXT_MANAGER", "VROPS", "VRSLCM", "VXRAIL_MANAGER"}, false),
			},
			"ca_certificate": {
				Type:         schema.TypeString,
				Required:     true,
				Description:  "Certificate of the CA issuing the replacement certificate",
				ValidateFunc: validation.StringIsNotEmpty,
			},
			"certificate_chain": {
				Type:         schema.TypeString,
				Required:     true,
				Description:  "Certificate Chain",
				ValidateFunc: validation.StringIsNotEmpty,
			},
			"resource_certificate": {
				Type:         schema.TypeString,
				Required:     true,
				Description:  "Resource Certificate",
				ValidateFunc: validation.StringIsNotEmpty,
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

	domainID := data.Get("domain_id").(string)
	resourceType := data.Get("resource").(string)

	resourceFqdn, err := certificates.GetFqdnOfResourceTypeInDomain(ctx, domainID, resourceType, apiClient)
	if err != nil {
		return diag.FromErr(err)
	}
	if resourceFqdn == nil {
		return diag.FromErr(fmt.Errorf("could not determine FQDN for resourceType %s in domain %s", resourceType, domainID))
	}

	caCertificate := data.Get("ca_certificate").(string)
	certificateChain := data.Get("certificate_chain").(string)
	resourceCertificate := data.Get("resource_certificate").(string)

	resourceCertificateSpec := &models.ResourceCertificateSpec{
		ResourceFqdn:        *resourceFqdn,
		CaCertificate:       caCertificate,
		CertificateChain:    certificateChain,
		ResourceCertificate: resourceCertificate,
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
	task, err := apiClient.Certificates.ReplaceResourceCertificates(replaceResourceCertificatesParams)
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
	data.SetId("ext_cert:" + domainID + ":" + resourceType + ":" + taskId)

	return resourceResourceExternalCertificateRead(ctx, data, meta)
}

func resourceResourceExternalCertificateRead(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	vcfClient := meta.(*api_client.SddcManagerClient)
	apiClient := vcfClient.ApiClient

	domainID := data.Get("domain_id").(string)
	resourceType := data.Get("resource").(string)

	cert, err := certificates.GetCertificateForResourceInDomain(ctx, apiClient, domainID, resourceType)
	if err != nil {
		return diag.FromErr(err)
	}

	flattenedCert := certificates.FlattenCertificate(cert)
	_ = data.Set("certificate", flattenedCert)

	return nil
}

func resourceResourceExternalCertificateUpdate(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return resourceResourceExternalCertificateRead(ctx, data, meta)
}

func resourceResourceExternalCertificateDelete(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}
