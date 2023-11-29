// Copyright 2023 Broadcom. All Rights Reserved.
// SPDX-License-Identifier: MPL-2.0

package validation

import (
	"errors"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/vmware/vcf-sdk-go/client/clusters"
	"github.com/vmware/vcf-sdk-go/client/domains"
	"github.com/vmware/vcf-sdk-go/models"
	"net/netip"
	"strings"
	"unicode"
)

func ValidatePassword(v interface{}, k string) (warnings []string, errors []error) {
	password, ok := v.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected not nil and type of %q to be string", k))
		return
	}
	var containsUpperCase, containsLowerCase, containsDigit, containsSymbol bool
	var specialSymbols = []rune{'\'', '!', '"', '#', '$', '%', '&', '(', ')', '*', '+', '-', '.', '/', ':', ';', '<', '=', '>', '?', '@', '[', '\\', ']', '^', '_', '`', '{', 'Ι', '}', '~'}
	for _, char := range password {
		if !containsLowerCase && unicode.IsLower(char) {
			containsLowerCase = true
		} else if !containsUpperCase && unicode.IsUpper(char) {
			containsUpperCase = true
		} else if !containsDigit && unicode.IsDigit(char) {
			containsDigit = true
		}
	}
	for _, symbol := range specialSymbols {
		if strings.ContainsRune(password, symbol) {
			containsSymbol = true
			break
		}
	}
	if len(password) < 8 {
		errors = append(errors, fmt.Errorf("the password must be at least 8 characters long"))
	}
	if !containsLowerCase {
		errors = append(errors, fmt.Errorf("the password must contain at least one lower case letter"))
	}
	if !containsUpperCase {
		errors = append(errors, fmt.Errorf("the password must contain at least one upper case letter"))
	}
	if !containsDigit {
		errors = append(errors, fmt.Errorf("the password must contain at least one digit"))
	}
	if !containsSymbol {
		errors = append(errors, fmt.Errorf("the password must contain at least one special symbol"))
	}
	return
}

func ValidateSddcId(v interface{}, k string) (warnings []string, errors []error) {
	sddcId, ok := v.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected not nil and type of %q to be string", k))
		return
	}
	if len(sddcId) < 3 || len(sddcId) > 20 {
		errors = append(errors, fmt.Errorf("sddcId can have length of 3-20 characters"))
		return
	}
	for _, char := range sddcId {
		if !unicode.IsLetter(char) && !unicode.IsDigit(char) && char != '-' {
			errors = append(errors, fmt.Errorf("can contain only letters, numbers and the following symbol: '-'"))
			return
		}
	}
	return
}

func ValidateParsingFloatToInt(v interface{}, k string) (warnings []string, errors []error) {
	floatNum := v.(float64)
	var intNum = int(floatNum)
	if floatNum != float64(intNum) {
		errors = append(errors, fmt.Errorf("expected an integer, got a float"))
	}
	return
}

func ConvertToStringSlice(params []interface{}) []string {
	var paramSlice []string
	for _, p := range params {
		if param, ok := p.(string); ok {
			paramSlice = append(paramSlice, param)
		}
	}
	return paramSlice
}

func validateIPv4Address(value string) error {
	addr, err := netip.ParseAddr(value)
	if err != nil {
		return err
	}

	if !addr.Is4() {
		return errors.New("invalid IPv4 address")
	}
	return nil
}

func ValidateIPv4AddressSchema(i interface{}, k string) (_ []string, errors []error) {
	ipAddress, ok := i.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected type of %s to be string", k))
		return nil, errors
	}
	ipValidationError := validateIPv4Address(ipAddress)
	if ipValidationError != nil {
		return nil, []error{}
	} else {
		return nil, nil
	}
}

func ConvertVcfErrorToDiag(err interface{}) diag.Diagnostics {
	if err == nil {
		return nil
	}
	domainsBadRequest, ok := err.(*domains.ValidateDomainsOperationsBadRequest)
	if ok {
		return convertVcfErrorsToDiagErrors(domainsBadRequest.Payload)
	}
	clustersBadRequest, ok := err.(*clusters.ValidateClusterOperationsBadRequest)
	if ok {
		return convertVcfErrorsToDiagErrors(clustersBadRequest.Payload)
	}
	createDomainBadRequest, ok := err.(*domains.CreateDomainBadRequest)
	if ok {
		return convertVcfErrorsToDiagErrors(createDomainBadRequest.Payload)
	}

	return diag.FromErr(err.(error))
}

func convertVcfErrorsToDiagErrors(err *models.Error) []diag.Diagnostic {
	var result []diag.Diagnostic

	var errorDetail string
	if IsEmpty(err.ReferenceToken) {
		errorDetail = err.RemediationMessage
	} else {
		errorDetail = fmt.Sprintf("look for reference token %q in service logs", err.ReferenceToken)
	}

	result = append(result, diag.Diagnostic{
		Severity: diag.Error,
		Summary:  err.Message,
		Detail:   errorDetail,
	})

	for _, nestedErr := range err.NestedErrors {
		result = append(result, convertVcfErrorsToDiagErrors(nestedErr)...)
	}
	return result
}

func HasValidationFailed(validationResult *models.Validation) bool {
	if validationResult == nil {
		return false
	}
	return validationResult.ResultStatus == "FAILED"
}

func HaveCertificateValidationsFailed(validationTask *models.CertificateValidationTask) bool {
	if validationTask == nil {
		return true
	}
	validationResult := validationTask.Validations
	for _, certValidation := range validationResult {
		if validationResult == nil && certValidation.ValidationStatus != nil {
			continue
		}
		if *certValidation.ValidationStatus == "FAILED" {
			return true
		}
	}
	return false
}

func ConvertValidationResultToDiag(validationResult *models.Validation) diag.Diagnostics {
	return convertValidationChecksToDiagErrors(validationResult.ValidationChecks)
}

func convertValidationChecksToDiagErrors(validationChecks []*models.ValidationCheck) []diag.Diagnostic {
	var result []diag.Diagnostic
	for _, validationCheck := range validationChecks {
		if validationCheck.Severity == "ERROR" || validationCheck.ResultStatus != "SUCCEEDED" {
			var validationErrorDetail string
			if len(validationCheck.ErrorResponse.NestedErrors) > 0 {
				for _, nestedError := range validationCheck.ErrorResponse.NestedErrors {
					validationErrorDetail += nestedError.Message + "\n"
				}
			} else {
				validationErrorDetail = validationCheck.ErrorResponse.Message
			}
			result = append(result, diag.Diagnostic{
				Severity: diag.Error,
				Summary:  validationCheck.Description,
				Detail:   validationErrorDetail,
			})
		}
		if len(validationCheck.NestedValidationChecks) > 0 {
			result = append(result, convertValidationChecksToDiagErrors(validationCheck.NestedValidationChecks)...)
		}
	}
	return result
}

func ConvertCertificateValidationsResultToDiag(validationTask *models.CertificateValidationTask) diag.Diagnostics {
	if validationTask == nil || validationTask.Validations == nil {
		return diag.FromErr(fmt.Errorf("provided certificate validation task is nil"))
	}
	return convertCertificateValidationChecksToDiagErrors(validationTask.Validations)
}

func convertCertificateValidationChecksToDiagErrors(validationChecks []*models.CertificateValidation) []diag.Diagnostic {
	var result []diag.Diagnostic
	for _, validationCheck := range validationChecks {
		if *validationCheck.ValidationStatus != "SUCCEEDED" {
			validationMessage := validationCheck.ValidationMessage
			result = append(result, diag.Diagnostic{
				Severity: diag.Error,
				Summary:  validationMessage,
			})
		}
	}
	return result
}

func HaveValidationChecksFinished(validationChecks []*models.ValidationCheck) bool {
	for _, validationCheck := range validationChecks {
		if validationCheck.ResultStatus == "IN_PROGRESS" {
			return false
		}
		if !HaveValidationChecksFinished(validationCheck.NestedValidationChecks) {
			return false
		}
	}
	return true
}

func HasCertificateValidationFinished(validationTask *models.CertificateValidationTask) bool {
	if validationTask == nil {
		return false
	}
	return *validationTask.Completed
}

func IsEmpty(object interface{}) bool {
	if object == nil {
		return true
	}
	_, ok := object.(bool)
	if ok {
		return false
	}
	objectStr, ok := object.(string)
	if ok {
		if len(objectStr) > 0 {
			return false
		}
	}
	objectAnySlice, ok := object.([]interface{})
	if ok {
		if len(objectAnySlice) > 0 {
			return false
		}
	}
	objectAnyMap, ok := object.(map[string]interface{})
	if ok {
		if len(objectAnyMap) > 0 {
			return false
		}
	}
	_, ok = object.(int)

	return !ok
}
