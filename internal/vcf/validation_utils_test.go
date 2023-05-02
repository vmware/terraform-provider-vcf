/* Copyright 2023 VMware, Inc.
   SPDX-License-Identifier: MPL-2.0 */

package vcf

import (
	"reflect"
	"strings"
	"testing"
)

func TestValidatePassword(t *testing.T) {
	t.Run("Validate password", func(t *testing.T) {
		var passwordTests = []struct {
			password    string
			expectedErr string
		}{
			{"Test1!", "the password must be at least 8 characters long"},
			{"testpassword1!", "the password must contain at least one upper case letter"},
			{"TESTPASSWORD1!", "the password must contain at least one lower case letter"},
			{"Testpassword!", "the password must contain at least one digit"},
			{"Testpassword1", "the password must contain at least one special symbol"},
		}

		for _, passTest := range passwordTests {
			err := validatePassword(passTest.password, "")
			if len(err) == 0 {
				t.Errorf("Failed. Expected one error for password %s, but got zero", passTest.password)
				break
			}
			if !strings.Contains(err[0].Error(), passTest.expectedErr) {
				t.Errorf("Failed. Unexpected error for password %s : %s, expected %s", passTest.password, err[0].Error(), passTest.expectedErr)
			}
		}
	})

	t.Run("Nil password validation", func(t *testing.T) {
		var expectedError = "expected not nil and type of \"\" to be string"

		err := validatePassword(nil, "")
		if len(err) == 0 {
			t.Fatalf("Failed. Expected one error for nil password, but got zero")
		}
		if !strings.Contains(err[0].Error(), expectedError) {
			t.Errorf("Failed. Unexpected error for nil password: %s, expected %s", err[0].Error(), expectedError)
		}
	})
}

func TestValidateParsingFloatToInt(t *testing.T) {
	var testFloatNotInt = 3.14
	var testFloatInt float64 = 3
	var expectedErr = "expected an integer, got a float"

	if err := validateParsingFloatToInt(testFloatNotInt); len(err) == 0 {
		t.Errorf("Failed. Expected error: \"%s\", for float64 %f", expectedErr, testFloatNotInt)
	}

	if err := validateParsingFloatToInt(testFloatInt); len(err) != 0 {
		t.Errorf("Failed. Expected no errors for float64 %f, got: \"%s\"", testFloatInt, err[0].Error())
	}
}

func TestConvertToStringSlice(t *testing.T) {
	var expectedStringSlice = []string{"test1", "test2"}
	testInterface := make([]interface{}, len(expectedStringSlice))
	for i, testString := range expectedStringSlice {
		testInterface[i] = testString
	}

	stringSlice := convertToStringSlice(testInterface)
	if !reflect.DeepEqual(stringSlice, expectedStringSlice) {
		t.Errorf("Failed. Expected string slice %v, got %v", expectedStringSlice, stringSlice)
	}
}
