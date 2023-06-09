/*
 *  Copyright 2023 VMware, Inc.
 *    SPDX-License-Identifier: MPL-2.0
 */

package validation

import (
	"errors"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/vmware/vcf-sdk-go/models"
	"net/netip"
	"strings"
	"unicode"
)

func ValidatePassword(v interface{}, k string) (errors []error) {
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

func ValidateParsingFloatToInt(v interface{}) (errors []error) {
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
	return nil, []error{validateIPv4Address(ipAddress)}
}

func HasValidationFailed(validationResult *models.Validation) bool {
	if validationResult == nil {
		return false
	}
	return validationResult.ExecutionStatus == "FAILED"
}

func ConvertValidationResultToDiag(validationResult *models.Validation) diag.Diagnostics {
	return diag.Diagnostics{}
}
