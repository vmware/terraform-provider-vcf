/* Copyright 2023 VMware, Inc.
   SPDX-License-Identifier: MPL-2.0 */

package vcf

import (
	"fmt"
	"strings"
	"unicode"
)

func validatePassword(v interface{}, k string) (warnings []string, errors []error) {
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

func validateParsingFloatToInt(v interface{}, k string) (warnings []string, errors []error) {
	floatNum := v.(float64)
	var intNum = int(floatNum)
	if floatNum != float64(intNum) {
		errors = append(errors, fmt.Errorf("expected an integer, got a float"))
	}
	return
}

func convertToStringSlice(params []interface{}) []string {
	var paramSlice []string
	for _, p := range params {
		if param, ok := p.(string); ok {
			paramSlice = append(paramSlice, param)
		}
	}
	return paramSlice
}
