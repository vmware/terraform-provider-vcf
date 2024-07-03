// Copyright 2023-2024 Broadcom. All Rights Reserved.
// SPDX-License-Identifier: MPL-2.0

package certificates

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/vmware/terraform-provider-vcf/internal/api_client"
	"github.com/vmware/terraform-provider-vcf/internal/constants"
	validationutils "github.com/vmware/terraform-provider-vcf/internal/validation"
	vcfclient "github.com/vmware/vcf-sdk-go/client"
	"github.com/vmware/vcf-sdk-go/client/certificates"
	"github.com/vmware/vcf-sdk-go/models"
	"time"
)

func ValidateResourceCertificates(ctx context.Context, client *vcfclient.VcfClient,
	domainId string, resourceCertificateSpecs []*models.ResourceCertificateSpec) diag.Diagnostics {
	validateResourceCertificatesParams := certificates.NewValidateResourceCertificatesParams().
		WithContext(ctx).WithTimeout(constants.DefaultVcfApiCallTimeout).
		WithID(domainId)
	validateResourceCertificatesParams.SetResourceCertificateSpecs(resourceCertificateSpecs)

	var validationResponse *models.CertificateValidationTask
	okResponse, acceptedResponse, err := client.Certificates.ValidateResourceCertificates(validateResourceCertificatesParams)
	if okResponse != nil {
		validationResponse = okResponse.Payload
	}
	if acceptedResponse != nil {
		validationResponse = acceptedResponse.Payload
	}
	if err != nil {
		return validationutils.ConvertVcfErrorToDiag(err)
	}
	if validationutils.HaveCertificateValidationsFailed(validationResponse) {
		return validationutils.ConvertCertificateValidationsResultToDiag(validationResponse)
	}
	validationId := validationResponse.ValidationID
	// Wait for certificate validation to finish
	if !validationutils.HasCertificateValidationFinished(validationResponse) {
		for {
			getResourceCertificatesValidationResultParams := certificates.NewGetResourceCertificatesValidationByIDParams().
				WithContext(ctx).
				WithTimeout(constants.DefaultVcfApiCallTimeout).
				WithID(*validationId)
			getValidationResponse, err := client.Certificates.GetResourceCertificatesValidationByID(getResourceCertificatesValidationResultParams)
			if err != nil {
				return validationutils.ConvertVcfErrorToDiag(err)
			}
			validationResponse = getValidationResponse.Payload
			if validationutils.HasCertificateValidationFinished(validationResponse) {
				break
			}
			time.Sleep(10 * time.Second)
		}
	}
	if err != nil {
		return validationutils.ConvertVcfErrorToDiag(err)
	}
	if validationutils.HaveCertificateValidationsFailed(validationResponse) {
		return validationutils.ConvertCertificateValidationsResultToDiag(validationResponse)
	}

	return nil
}

func GetCertificateForResourceInDomain(ctx context.Context, client *vcfclient.VcfClient,
	domainId, resourceFqdn string) (*models.Certificate, error) {
	viewCertificatesParams := certificates.NewGetCertificatesByDomainParamsWithContext(ctx).
		WithTimeout(constants.DefaultVcfApiCallTimeout)
	viewCertificatesParams.ID = domainId

	certificatesResponse, _, err := client.Certificates.GetCertificatesByDomain(viewCertificatesParams)
	if err != nil {
		return nil, err
	}

	allCertsForDomain := certificatesResponse.Payload.Elements
	for _, cert := range allCertsForDomain {
		if cert.IssuedTo != nil && *cert.IssuedTo == resourceFqdn {
			return cert, nil
		}
	}
	return nil, nil
}

func GenerateCertificateForResource(ctx context.Context, client *api_client.SddcManagerClient,
	domainId, resourceType, resourceFqdn, caType *string) error {

	certificateGenerationSpec := &models.CertificatesGenerationSpec{
		CaType: caType,
		Resources: []*models.Resource{{
			Fqdn: *resourceFqdn,
			Type: resourceType,
		}},
	}
	generateCertificatesParam := certificates.NewGenerateCertificatesParamsWithContext(ctx).
		WithTimeout(constants.DefaultVcfApiCallTimeout).
		WithID(*domainId)
	generateCertificatesParam.SetCertificateGenerationSpec(certificateGenerationSpec)

	var taskId string
	responseOk, responseAccepted, err := client.ApiClient.Certificates.GenerateCertificates(generateCertificatesParam)
	if err != nil {
		return err
	}
	if responseOk != nil {
		taskId = responseOk.Payload.ID
	}
	if responseAccepted != nil {
		taskId = responseAccepted.Payload.ID
	}
	err = client.WaitForTaskComplete(ctx, taskId, true)
	if err != nil {
		return err
	}
	return nil
}
