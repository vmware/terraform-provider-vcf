// © Broadcom. All Rights Reserved.
// The term “Broadcom” refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

package validation

import (
	"errors"
	"fmt"
	"net/netip"
	"strconv"
	"strings"
	"unicode"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/vmware/vcf-sdk-go/vcf"
)

func ValidatePassword(v interface{}, k string) (warnings []string, errors []error) {
	var specialSymbols = []rune{'\'', '!', '"', '#', '$', '%', '&', '(', ')', '*', '+', '-', '.', '/', ':', ';', '<', '=', '>', '?', '@', '[', '\\', ']', '^', '_', '`', '{', 'Ι', '}', '~'}
	return validatePasswordInternal(v, k, specialSymbols, 8)
}

func ValidateNsxEdgePassword(v interface{}, k string) (warnings []string, errors []error) {
	var specialSymbols = []rune{'!', '@', '^', '=', '*', '+'}
	return validatePasswordInternal(v, k, specialSymbols, 12)
}

func validatePasswordInternal(v interface{}, k string, specialSymbols []rune, minLength int) (warnings []string, errors []error) {
	password, ok := v.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected not nil and type of %q to be string", k))
		return
	}
	var containsUpperCase, containsLowerCase, containsDigit, containsSymbol bool
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
	if len(password) < minLength {
		errors = append(errors, fmt.Errorf("the password must be at least %d characters long", minLength))
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

func validateCidrIPv4Address(value string) error {
	prefix, err := netip.ParsePrefix(value)
	if err != nil {
		return err
	}

	if !prefix.IsValid() || !prefix.Addr().Is4() {
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
		return nil, []error{ipValidationError}
	} else {
		return nil, nil
	}
}

func ValidateCidrIPv4AddressSchema(i interface{}, k string) (_ []string, errors []error) {
	ipAddress, ok := i.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected type of %s to be string", k))
		return nil, errors
	}
	ipValidationError := validateCidrIPv4Address(ipAddress)
	if ipValidationError != nil {
		return nil, []error{ipValidationError}
	} else {
		return nil, nil
	}
}

func ConvertVcfErrorToDiag(err interface{}) diag.Diagnostics {
	if err == nil {
		return nil
	}

	if vcfError, ok := err.(*vcf.Error); ok {
		return convertVcfErrorsToDiagErrors(vcfError)
	}

	return diag.FromErr(err.(error))
}

func convertVcfErrorsToDiagErrors(err *vcf.Error) []diag.Diagnostic {
	var result []diag.Diagnostic

	var errorDetail string
	if err.RemediationMessage != nil && IsEmpty(*err.ReferenceToken) {
		errorDetail = *err.RemediationMessage
	} else {
		errorDetail = fmt.Sprintf("look for reference token %q in service logs", *err.ReferenceToken)
	}

	result = append(result, diag.Diagnostic{
		Severity: diag.Error,
		Summary:  *err.Message,
		Detail:   errorDetail,
	})

	return result
}

func HasValidationFailed(validationResult *vcf.Validation) bool {
	if validationResult == nil {
		return false
	}
	return validationResult.ResultStatus != nil && *validationResult.ResultStatus == "FAILED"
}

func HaveCertificateValidationsFailed(validationTask *vcf.CertificateValidationTask) bool {
	if validationTask == nil {
		return true
	}
	validationResult := validationTask.Validations
	for _, certValidation := range validationResult {
		if validationResult == nil && certValidation.ValidationStatus != "" {
			continue
		}
		if certValidation.ValidationStatus == "FAILED" {
			return true
		}
	}
	return false
}

func ConvertValidationResultToDiag(validationResult *vcf.Validation) diag.Diagnostics {
	return convertValidationChecksToDiagErrors(validationResult.ValidationChecks)
}

func convertValidationChecksToDiagErrors(validationChecks *[]vcf.ValidationCheck) []diag.Diagnostic {
	var result []diag.Diagnostic
	if validationChecks != nil {
		for _, validationCheck := range *validationChecks {
			severity := validationCheck.Severity
			if (severity != nil && *severity == "ERROR") || validationCheck.ResultStatus != "SUCCEEDED" {
				var validationErrorDetail string
				if validationCheck.ErrorResponse != nil &&
					validationCheck.ErrorResponse.NestedErrors != nil &&
					len(*validationCheck.ErrorResponse.NestedErrors) > 0 {
					for _, nestedError := range *validationCheck.ErrorResponse.NestedErrors {
						validationErrorDetail += *nestedError.Message + "\n"
					}
				} else {
					validationErrorDetail = *validationCheck.ErrorResponse.Message
				}
				diagnostic := diag.Diagnostic{
					Severity: diag.Error,
					Detail:   validationErrorDetail,
				}

				if validationCheck.Description != nil {
					diagnostic.Summary = *validationCheck.Description
				}

				result = append(result, diagnostic)
			}
		}
	}
	return result
}

func ConvertCertificateValidationsResultToDiag(validationTask *vcf.CertificateValidationTask) diag.Diagnostics {
	if validationTask == nil || validationTask.Validations == nil {
		return diag.FromErr(fmt.Errorf("provided certificate validation task is nil"))
	}
	return convertCertificateValidationChecksToDiagErrors(validationTask.Validations)
}

func convertCertificateValidationChecksToDiagErrors(validationChecks []vcf.CertificateValidation) []diag.Diagnostic {
	var result []diag.Diagnostic
	for _, validationCheck := range validationChecks {
		if validationCheck.ValidationStatus != "SUCCEEDED" {
			validationMessage := validationCheck.ValidationMessage
			result = append(result, diag.Diagnostic{
				Severity: diag.Error,
				Summary:  *validationMessage,
			})
		}
	}
	return result
}

func HaveValidationChecksFinished(validationChecks []vcf.ValidationCheck) bool {
	for _, validationCheck := range validationChecks {
		if validationCheck.ResultStatus == "IN_PROGRESS" || validationCheck.ResultStatus == "UNKNOWN" {
			return false
		}
	}
	return true
}

func HasCertificateValidationFinished(validationTask *vcf.CertificateValidationTask) bool {
	if validationTask == nil {
		return false
	}
	return validationTask.Completed
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

func ValidASN(v interface{}, k string) (ws []string, errors []error) {
	value := v.(string)

	asn, err := strconv.ParseInt(value, 10, 64)
	if err != nil {
		errors = append(errors, fmt.Errorf("%q (%q) must be a 64-bit integer", k, v))
		return
	}

	isLegacyAsn := func(a int64) bool {
		return a == 7224 || a == 9059 || a == 10124 || a == 17493
	}

	if !isLegacyAsn(asn) && ((asn < 64512) || (asn > 65534 && asn < 4200000000) || (asn > 4294967294)) {
		errors = append(errors, fmt.Errorf("%q (%q) must be 7224, 9059, 10124 or 17493 or in the range 64512 to 65534 or 4200000000 to 4294967294", k, v))
	}
	return
}
