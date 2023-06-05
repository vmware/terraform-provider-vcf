/* Copyright 2023 VMware, Inc.
   SPDX-License-Identifier: MPL-2.0 */

package provider

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/vmware/terraform-provider-vcf/internal/constants"
	"os"
	"testing"
)

var testAccProvider *schema.Provider

func init() {
	testAccProvider = Provider()
}

func TestProvider(t *testing.T) {
	if err := Provider().InternalValidate(); err != nil {
		t.Fatalf("err: %s", err)
	}
}

var providerFactories = map[string]func() (*schema.Provider, error){
	"vcf": func() (*schema.Provider, error) {
		return testAccProvider, nil
	},
}

// testAccPreCheck validates all required environment variables for running acceptance
// tests are set.
func testAccPreCheck(t *testing.T) {
	if v := os.Getenv(constants.VCF_TEST_URL); v == "" {
		t.Fatal(constants.VCF_TEST_URL + " must be set for acceptance tests")
	}
	if v := os.Getenv(constants.VCF_TEST_USERNAME); v == "" {
		t.Fatal(constants.VCF_TEST_USERNAME + " must be set for acceptance tests")
	}
	if v := os.Getenv(constants.VCF_TEST_PASSWORD); v == "" {
		t.Fatal(constants.VCF_TEST_PASSWORD + " must be set for acceptance tests")
	}
	if v := os.Getenv(constants.VCF_TEST_COMMISSIONED_HOST_FQDN); v == "" {
		t.Fatal(constants.VCF_TEST_COMMISSIONED_HOST_FQDN + " must be set for acceptance tests")
	}
	if v := os.Getenv(constants.VCF_TEST_COMMISSIONED_HOST_PASS); v == "" {
		t.Fatal(constants.VCF_TEST_COMMISSIONED_HOST_PASS + " must be set for acceptance tests")
	}
}
