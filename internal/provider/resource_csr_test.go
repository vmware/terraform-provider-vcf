// Copyright 2023-2024 Broadcom. All Rights Reserved.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/vmware/terraform-provider-vcf/internal/constants"
	"os"
	"testing"
)

func TestAccResourceVcfCsr_vCenter(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccVcfCsrConfig(os.Getenv(constants.VcfTestDomainDataSourceId), "VCENTER", os.Getenv(constants.VcfTestVcenterFqdn)),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("vcf_csr.csr1", "csr.0.csr_pem"),
					resource.TestCheckResourceAttrSet("vcf_csr.csr1", "csr.0.csr_string"),
					resource.TestCheckResourceAttrSet("vcf_csr.csr1", "csr.0.resource.0.fqdn")),
			},
		},
	})
}

func TestAccResourceVcfCsr_sddcManager(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccVcfCsrConfig(os.Getenv(constants.VcfTestDomainDataSourceId), "SDDC_MANAGER", os.Getenv(constants.VcfTestSddcManagerFqdn)),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("vcf_csr.csr1", "csr.0.csr_pem"),
					resource.TestCheckResourceAttrSet("vcf_csr.csr1", "csr.0.csr_string"),
					resource.TestCheckResourceAttrSet("vcf_csr.csr1", "csr.0.resource.0.fqdn")),
			},
		},
	})
}

func TestAccResourceVcfCsr_nsxManager(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccVcfCsrConfig(os.Getenv(constants.VcfTestDomainDataSourceId), "NSXT_MANAGER", os.Getenv(constants.VcfTestNsxManagerFqdn)),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("vcf_csr.csr1", "csr.0.csr_pem"),
					resource.TestCheckResourceAttrSet("vcf_csr.csr1", "csr.0.csr_string"),
					resource.TestCheckResourceAttrSet("vcf_csr.csr1", "csr.0.resource.0.fqdn")),
			},
		},
	})
}

func testAccVcfCsrConfig(domainID, resource, fqdn string) string {
	return fmt.Sprintf(`
	resource "vcf_csr" "csr1" {
  		domain_id = %q
		country = "BG"
		email = "admin@vmware.com"
		key_size = "3072"
		locality = "Sofia"
		state = "Sofia-grad"
		organization = "VMware Inc."
		organization_unit = "VCF"
		resource = %q
		fqdn = %q
	}`,
		domainID, resource, fqdn,
	)
}
