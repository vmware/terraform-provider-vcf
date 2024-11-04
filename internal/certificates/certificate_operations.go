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
	task, vcfErr := api_client.GetResponseAs[vcf.CertificateValidationTask](okResponse.Body)
	if vcfErr != nil {
		api_client.LogError(vcfErr)
		return diag.FromErr(errors.New(*vcfErr.Message))
	}
	if validationutils.HaveCertificateValidationsFailed(task) {
		return validationutils.ConvertCertificateValidationsResultToDiag(task)
	}
	// Wait for certificate validation to finish
	if !validationutils.HasCertificateValidationFinished(task) {
		for {
			getValidationResponse, err := client.GetResourceCertificatesValidationByIDWithResponse(ctx, domainId, task.ValidationId)
			if err != nil {
				return validationutils.ConvertVcfErrorToDiag(err)
			}
			task, vcfErr = api_client.GetResponseAs[vcf.CertificateValidationTask](getValidationResponse.Body)
			if vcfErr != nil {
				api_client.LogError(vcfErr)
				return diag.FromErr(errors.New(*vcfErr.Message))
			}
			if validationutils.HasCertificateValidationFinished(task) {
				break
			}
			time.Sleep(10 * time.Second)
		}
	}
	if err != nil {
		return validationutils.ConvertVcfErrorToDiag(err)
	}
	if validationutils.HaveCertificateValidationsFailed(task) {
		return validationutils.ConvertCertificateValidationsResultToDiag(task)
	}

	return nil
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

	res, err := client.ApiClient.GenerateCertificatesWithResponse(ctx, *domainId, certificateGenerationSpec)
	if err != nil {
		return err
	}
	task, vcfErr := api_client.GetResponseAs[vcf.Task](res.Body)
	if vcfErr != nil {
		api_client.LogError(vcfErr)
		return errors.New(*vcfErr.Message)
	}
	err = client.WaitForTaskComplete(ctx, *task.Id, true)
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

	page, vcfErr := api_client.GetResponseAs[vcf.PageOfCertificate](certificatesResponse.Body)
	if vcfErr != nil {
		api_client.LogError(vcfErr)
		return nil, errors.New(*vcfErr.Message)
	}
	// Check if any certificates are found
	if page == nil || len(*page.Elements) == 0 {
		return nil, fmt.Errorf("no certificates found for domain ID %s", domainId)
	}

	allCertsForDomain := page.Elements
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
