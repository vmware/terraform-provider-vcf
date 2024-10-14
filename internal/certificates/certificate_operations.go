// © Broadcom. All Rights Reserved.
// The term “Broadcom” refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

package certificates

import (
	"context"
	md52 "crypto/md5"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/vmware/vcf-sdk-go/vcf"

	"github.com/vmware/terraform-provider-vcf/internal/api_client"
	validationutils "github.com/vmware/terraform-provider-vcf/internal/validation"
)

func ValidateResourceCertificates(ctx context.Context, client *vcf.ClientWithResponses,
	domainId string, resourceCertificateSpecs []vcf.ResourceCertificateSpec) diag.Diagnostics {
	okResponse, err := client.ValidateResourceCertificatesWithResponse(ctx, domainId, resourceCertificateSpecs)
	if err != nil {
		return diag.FromErr(err)
	}
	if okResponse.StatusCode() != 201 {
		vcfError := api_client.GetError(okResponse.Body)
		api_client.LogError(vcfError)
		return diag.FromErr(errors.New(*vcfError.Message))
	}
	if validationutils.HaveCertificateValidationsFailed(okResponse.JSON201) {
		return validationutils.ConvertCertificateValidationsResultToDiag(okResponse.JSON201)
	}
	// Wait for certificate validation to finish
	if !validationutils.HasCertificateValidationFinished(okResponse.JSON201) {
		for {
			getValidationResponse, err := client.GetResourceCertificatesValidationByIDWithResponse(ctx, domainId, okResponse.JSON201.ValidationId)
			if err != nil {
				return validationutils.ConvertVcfErrorToDiag(err)
			}
			if getValidationResponse.StatusCode() != 201 {
				vcfError := api_client.GetError(getValidationResponse.Body)
				api_client.LogError(vcfError)
				return diag.FromErr(errors.New(*vcfError.Message))
			}
			if validationutils.HasCertificateValidationFinished(okResponse.JSON201) {
				break
			}
			time.Sleep(10 * time.Second)
		}
	}
	if err != nil {
		return validationutils.ConvertVcfErrorToDiag(err)
	}
	if validationutils.HaveCertificateValidationsFailed(okResponse.JSON201) {
		return validationutils.ConvertCertificateValidationsResultToDiag(okResponse.JSON201)
	}

	return nil
}

func GetCertificateForResourceInDomain(ctx context.Context, client *vcf.ClientWithResponses,
	domainId, resourceFqdn string) (*vcf.Certificate, error) {
	certificatesResponse, err := client.GetCertificatesByDomainWithResponse(ctx, domainId)
	if err != nil {
		return nil, err
	}
	if certificatesResponse.StatusCode() != 200 {
		vcfError := api_client.GetError(certificatesResponse.Body)
		api_client.LogError(vcfError)
		return nil, errors.New(*vcfError.Message)
	}

	allCertsForDomain := certificatesResponse.JSON200.Elements
	for _, cert := range *allCertsForDomain {
		if cert.IssuedTo != nil && *cert.IssuedTo == resourceFqdn {
			return &cert, nil
		}
	}
	return nil, nil
}

func GenerateCertificateForResource(ctx context.Context, client *api_client.SddcManagerClient,
	domainId, resourceType, resourceFqdn, caType *string) error {

	certificateGenerationSpec := vcf.CertificatesGenerationSpec{
		CaType: *caType,
		Resources: &[]vcf.Resource{{
			Fqdn: resourceFqdn,
			Type: *resourceType,
		}},
	}

	var taskId string
	res, err := client.ApiClient.GenerateCertificatesWithResponse(ctx, *domainId, certificateGenerationSpec)
	if err != nil {
		return err
	}
	if res != nil && res.JSON202 != nil {
		taskId = *res.JSON202.Id
	}
	err = client.WaitForTaskComplete(ctx, taskId, true)
	if err != nil {
		return err
	}
	return nil
}

func ReadCertificate(ctx context.Context, client *vcf.ClientWithResponses,
	domainId, resourceFqdn string) (*vcf.Certificate, error) {

	certificatesResponse, err := client.GetCertificatesByDomainWithResponse(ctx, domainId)
	if err != nil {
		return nil, fmt.Errorf("failed to get certificate by domain: %w", err)
	}

	// Check if any certificates are found
	if certificatesResponse.JSON200 == nil || len(*certificatesResponse.JSON200.Elements) == 0 {
		return nil, fmt.Errorf("no certificates found for domain ID %s", domainId)
	}

	allCertsForDomain := certificatesResponse.JSON200.Elements
	for _, cert := range *allCertsForDomain {
		if cert.IssuedTo != nil && *cert.IssuedTo == resourceFqdn {
			return &cert, nil
		}
	}
	return nil, fmt.Errorf("no certificate found for resource FQDN %s in domain ID %s", resourceFqdn, domainId)
}

func FlattenCertificateWithSubject(cert *vcf.Certificate) map[string]interface{} {
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
