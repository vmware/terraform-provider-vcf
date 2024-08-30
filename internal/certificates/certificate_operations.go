// © Broadcom. All Rights Reserved.
// The term “Broadcom” refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

package certificates

import (
	"context"
	md52 "crypto/md5"
	"encoding/hex"
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	vcfclient "github.com/vmware/vcf-sdk-go/client"
	"github.com/vmware/vcf-sdk-go/client/certificates"
	"github.com/vmware/vcf-sdk-go/models"

	"github.com/vmware/terraform-provider-vcf/internal/api_client"
	"github.com/vmware/terraform-provider-vcf/internal/constants"
	validationutils "github.com/vmware/terraform-provider-vcf/internal/validation"
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

func ReadCertificate(ctx context.Context, client *vcfclient.VcfClient,
	domainId, resourceFqdn string) (*models.Certificate, error) {
	viewCertificatesParams := certificates.NewGetCertificatesByDomainParamsWithContext(ctx).
		WithTimeout(constants.DefaultVcfApiCallTimeout)
	viewCertificatesParams.ID = domainId

	certificatesResponse, _, err := client.Certificates.GetCertificatesByDomain(viewCertificatesParams)
	if err != nil {
		return nil, fmt.Errorf("failed to get certificate by domain: %w", err)
	}

	// Check if any certificates are found
	if certificatesResponse.Payload == nil || len(certificatesResponse.Payload.Elements) == 0 {
		return nil, fmt.Errorf("no certificates found for domain ID %s", domainId)
	}

	allCertsForDomain := certificatesResponse.Payload.Elements
	for _, cert := range allCertsForDomain {
		if cert.IssuedTo != nil && *cert.IssuedTo == resourceFqdn {
			return cert, nil
		}
	}
	return nil, fmt.Errorf("no certificate found for resource FQDN %s in domain ID %s", resourceFqdn, domainId)
}

func FlattenCertificateWithSubject(cert *models.Certificate) map[string]interface{} {
	result := make(map[string]interface{})
	if cert.Domain == nil {
		result["domain"] = nil
	} else {
		result["domain"] = *cert.Domain
	}
	if cert.GetCertificateError == nil {
		result["certificate_error"] = nil
	} else {
		result["certificate_error"] = *cert.GetCertificateError
	}

	result["expiration_status"] = *cert.ExpirationStatus
	result["issued_by"] = *cert.IssuedBy
	result["issued_to"] = *cert.IssuedTo
	result["key_size"] = *cert.KeySize
	result["not_after"] = *cert.NotAfter
	result["not_before"] = *cert.NotBefore
	result["number_of_days_to_expire"] = *cert.NumberOfDaysToExpire
	result["pem_encoded"] = *cert.PemEncoded
	result["public_key"] = *cert.PublicKey
	result["public_key_algorithm"] = *cert.PublicKeyAlgorithm
	result["serial_number"] = *cert.SerialNumber
	result["signature_algorithm"] = *cert.SignatureAlgorithm
	result["subject"] = *cert.Subject
	result["subject_alternative_name"] = cert.SubjectAlternativeName
	result["thumbprint"] = *cert.Thumbprint
	result["thumbprint_algorithm"] = *cert.ThumbprintAlgorithm
	result["version"] = *cert.Version

	// Parse the subject string to extract CN, OU, O, L, ST, C
	subjectDetails := parseSubject(*cert.Subject)

	// Add parsed subject components to the result map
	result["subject_cn"] = subjectDetails["CN"]
	result["subject_ou"] = subjectDetails["OU"]
	result["subject_org"] = subjectDetails["O"]
	result["subject_locality"] = subjectDetails["L"]
	result["subject_st"] = subjectDetails["ST"]
	result["subject_country"] = subjectDetails["C"]

	return result
}

func parseSubject(subject string) map[string]string {
	parsedSubject := make(map[string]string)

	// Split the subject string by commas to separate key-value pairs
	subjectParts := strings.Split(subject, ",")

	for _, part := range subjectParts {
		// Split each part by the equals sign to separate the key and value
		keyValue := strings.SplitN(strings.TrimSpace(part), "=", 2)
		if len(keyValue) == 2 {
			// Store the value in the map with the key as the subject component
			parsedSubject[keyValue[0]] = keyValue[1]
		}
	}

	return parsedSubject
}

func HashFields(fields []string) (string, error) {
	md5 := md52.New()
	_, err := io.WriteString(md5, strings.Join(fields, ""))

	if err != nil {
		return "", err
	}

	return hex.EncodeToString(md5.Sum(nil)), nil
}
