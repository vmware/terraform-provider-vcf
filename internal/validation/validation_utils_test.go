// Copyright 2023 Broadcom. All Rights Reserved.
// SPDX-License-Identifier: MPL-2.0

package validation

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
			_, err := ValidatePassword(passTest.password, "")
			if len(err) == 0 {
				t.Errorf("failed. expected one error for password %s, but got zero", passTest.password)
				break
			}
			if !strings.Contains(err[0].Error(), passTest.expectedErr) {
				t.Errorf("failed. Unexpected error for password %s : %s, expected %s", passTest.password, err[0].Error(), passTest.expectedErr)
			}
		}
	})

	t.Run("Nil password validation", func(t *testing.T) {
		var expectedError = "expected not nil and type of \"\" to be string"

		_, err := ValidatePassword(nil, "")
		if len(err) == 0 {
			t.Fatalf("failed. expected one error for nil password, but got zero")
		}
		if !strings.Contains(err[0].Error(), expectedError) {
			t.Errorf("Failed. Unexpected error for nil password: %s, expected %s", err[0].Error(), expectedError)
		}
	})
}

func TestValidateSddcId(t *testing.T) {
	t.Run("Validate sddc Id", func(t *testing.T) {
		var sddcIdTests = []struct {
			sddcId      string
			expectedErr string
		}{
			{"Test1!", "can contain only letters, numbers and the following symbol: '-'"},
			{"te", "sddcId can have length of 3-20 characters"},
			{"Test%", "can contain only letters, numbers and the following symbol: '-'"},
		}

		for _, sddcIdTest := range sddcIdTests {
			_, err := ValidateSddcId(sddcIdTest.sddcId, "")
			if len(err) == 0 {
				t.Errorf("failed. expected one error for sddcId %s, but got zero", sddcIdTest.sddcId)
				break
			}
			if !strings.Contains(err[0].Error(), sddcIdTest.expectedErr) {
				t.Errorf("failed. Unexpected error for sddcId %s : %s, expected %s", sddcIdTest.sddcId, err[0].Error(), sddcIdTest.expectedErr)
			}
		}
	})

	t.Run("Nil sddcId validation", func(t *testing.T) {
		var expectedError = "expected not nil and type of \"\" to be string"

		_, err := ValidateSddcId(nil, "")
		if len(err) == 0 {
			t.Fatalf("failed. expected one error for nil sddcId, but got zero")
		}
		if !strings.Contains(err[0].Error(), expectedError) {
			t.Errorf("Failed. Unexpected error for nil sddcId: %s, expected %s", err[0].Error(), expectedError)
		}
	})
}

func TestValidateParsingFloatToInt(t *testing.T) {
	var testFloatNotInt = 3.14
	var testFloatInt float64 = 3
	var expectedErr = "expected an integer, got a float"

	if _, err := ValidateParsingFloatToInt(testFloatNotInt, ""); len(err) == 0 {
		t.Errorf("Failed. Expected error: \"%s\", for float64 %f", expectedErr, testFloatNotInt)
	}

	if _, err := ValidateParsingFloatToInt(testFloatInt, ""); len(err) != 0 {
		t.Errorf("Failed. Expected no errors for float64 %f, got: \"%s\"", testFloatInt, err[0].Error())
	}
}

func TestConvertToStringSlice(t *testing.T) {
	var expectedStringSlice = []string{"test1", "test2"}
	testInterface := make([]interface{}, len(expectedStringSlice))
	for i, testString := range expectedStringSlice {
		testInterface[i] = testString
	}

	stringSlice := ConvertToStringSlice(testInterface)
	if !reflect.DeepEqual(stringSlice, expectedStringSlice) {
		t.Errorf("Failed. Expected string slice %v, got %v", expectedStringSlice, stringSlice)
	}
}

func TestValidateIpv4Address(t *testing.T) {
	t.Run("Validate ipv4 address", func(t *testing.T) {
		var ipTests = []struct {
			ip          string
			expectError bool
		}{
			{"192.168.0.1", false},
			{"255.255.255.0", false},
			{"random text", true},
			{"420.168.0.1", true},
			{"120.168.01", true},
			{"01.168.0.1", true},
		}

		for _, ipTest := range ipTests {
			err := validateIPv4Address(ipTest.ip)
			if ipTest.expectError && err == nil {
				t.Errorf("failed. Unexpected error occurred.")
			}
			if !ipTest.expectError && err != nil {
				t.Errorf("failed. Expected error.")
			}
		}
	})
}

func TestIsEmpty(t *testing.T) {
	t.Run("is object empty", func(t *testing.T) {
		var nonEmptyMap = make(map[string]interface{})
		nonEmptyMap["name"] = "SlavaZSU"
		var isEmptyTests = []struct {
			object         interface{}
			expectedResult bool
		}{
			{"some string", false},
			{nil, true},
			{true, false},
			{new([]interface{}), true},
			{make(map[string]interface{}), true},
			{append(make([]interface{}, 0), "first", "second"), false},
			{nonEmptyMap, false},
		}

		for _, emptyTest := range isEmptyTests {
			result := IsEmpty(emptyTest.object)
			if emptyTest.expectedResult != result {
				t.Errorf("%s test failed", emptyTest.object)
			}
		}
	})
}
