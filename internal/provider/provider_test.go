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
	if v := os.Getenv(constants.VcfTestUrl); v == "" {
		t.Fatal(constants.VcfTestUrl + " must be set for acceptance tests")
	}
	if v := os.Getenv(constants.VcfTestUsername); v == "" {
		t.Fatal(constants.VcfTestUsername + " must be set for acceptance tests")
	}
	if v := os.Getenv(constants.VcfTestPassword); v == "" {
		t.Fatal(constants.VcfTestPassword + " must be set for acceptance tests")
	}
	if v := os.Getenv(constants.VcfTestHost1Fqdn); v == "" {
		t.Fatal(constants.VcfTestHost1Fqdn + " must be set for acceptance tests")
	}
	if v := os.Getenv(constants.VcfTestHost1Pass); v == "" {
		t.Fatal(constants.VcfTestHost1Pass + " must be set for acceptance tests")
	}
	if v := os.Getenv(constants.VcfTestHost2Fqdn); v == "" {
		t.Fatal(constants.VcfTestHost2Fqdn + " must be set for acceptance tests")
	}
	if v := os.Getenv(constants.VcfTestHost2Pass); v == "" {
		t.Fatal(constants.VcfTestHost2Pass + " must be set for acceptance tests")
	}
	if v := os.Getenv(constants.VcfTestHost3Fqdn); v == "" {
		t.Fatal(constants.VcfTestHost3Fqdn + " must be set for acceptance tests")
	}
	if v := os.Getenv(constants.VcfTestHost3Pass); v == "" {
		t.Fatal(constants.VcfTestHost3Pass + " must be set for acceptance tests")
	}
	if v := os.Getenv(constants.VcfTestHost4Fqdn); v == "" {
		t.Fatal(constants.VcfTestHost2Fqdn + " must be set for acceptance tests")
	}
	if v := os.Getenv(constants.VcfTestHost4Pass); v == "" {
		t.Fatal(constants.VcfTestHost2Pass + " must be set for acceptance tests")
	}
	if v := os.Getenv(constants.VcfTestNsxLicenseKey); v == "" {
		t.Fatal(constants.VcfTestNsxLicenseKey + " must be set for acceptance tests")
	}
	if v := os.Getenv(constants.VcfTestEsxiLicenseKey); v == "" {
		t.Fatal(constants.VcfTestEsxiLicenseKey + " must be set for acceptance tests")
	}
	if v := os.Getenv(constants.VcfTestVsanLicenseKey); v == "" {
		t.Fatal(constants.VcfTestVsanLicenseKey + " must be set for acceptance tests")
	}
}
