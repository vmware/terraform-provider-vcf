/*
 *  Copyright 2023 VMware, Inc.
 *    SPDX-License-Identifier: MPL-2.0
 */

package provider

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/vmware/terraform-provider-vcf/internal/constants"
	"os"
	"testing"
)

func TestAccResourceVcfResourceExternalCertificate(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccVcfResourceExternalCertificatePreCheck(t)
		},
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccVcfResourceExternalCertificate(
					os.Getenv(constants.VcfTestDomainDataSourceId),
					os.Getenv(constants.VcfTestResourceCertificate),
					os.Getenv(constants.VcfTestResourceCaCertificate),
					os.Getenv(constants.VcfTestResourceCertificateChain)),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("vcf_external_certificate.vcenter_cert", "certificate.0.issued_by"),
					resource.TestCheckResourceAttrSet("vcf_external_certificate.vcenter_cert", "certificate.0.issued_to"),
					resource.TestCheckResourceAttrSet("vcf_external_certificate.vcenter_cert", "certificate.0.expiration_status"),
					resource.TestCheckResourceAttrSet("vcf_external_certificate.vcenter_cert", "certificate.0.certificate_error"),
					resource.TestCheckResourceAttrSet("vcf_external_certificate.vcenter_cert", "certificate.0.is_installed"),
					resource.TestCheckResourceAttrSet("vcf_external_certificate.vcenter_cert", "certificate.0.issued_by"),
					resource.TestCheckResourceAttrSet("vcf_external_certificate.vcenter_cert", "certificate.0.issued_to"),
					resource.TestCheckResourceAttrSet("vcf_external_certificate.vcenter_cert", "certificate.0.key_size"),
					resource.TestCheckResourceAttrSet("vcf_external_certificate.vcenter_cert", "certificate.0.not_after"),
					resource.TestCheckResourceAttrSet("vcf_external_certificate.vcenter_cert", "certificate.0.not_before"),
					resource.TestCheckResourceAttrSet("vcf_external_certificate.vcenter_cert", "certificate.0.pem_encoded"),
					resource.TestCheckResourceAttrSet("vcf_external_certificate.vcenter_cert", "certificate.0.public_key"),
					resource.TestCheckResourceAttrSet("vcf_external_certificate.vcenter_cert", "certificate.0.public_key_algorithm"),
					resource.TestCheckResourceAttrSet("vcf_external_certificate.vcenter_cert", "certificate.0.serial_number"),
					resource.TestCheckResourceAttrSet("vcf_external_certificate.vcenter_cert", "certificate.0.signature_algorithm"),
					resource.TestCheckResourceAttrSet("vcf_external_certificate.vcenter_cert", "certificate.0.subject"),
					resource.TestCheckResourceAttrSet("vcf_external_certificate.vcenter_cert", "certificate.0.subject_alternative_name"),
					resource.TestCheckResourceAttrSet("vcf_external_certificate.vcenter_cert", "certificate.0.thumbprint"),
					resource.TestCheckResourceAttrSet("vcf_external_certificate.vcenter_cert", "certificate.0.thumbprint_algorithm"),
					resource.TestCheckResourceAttrSet("vcf_external_certificate.vcenter_cert", "certificate.0.version"),
					resource.TestCheckResourceAttrSet("vcf_external_certificate.vcenter_cert", "certificate.0.number_of_days_to_expire")),
			},
		},
	})
}

func testAccVcfResourceExternalCertificatePreCheck(t *testing.T) {
	if v := os.Getenv(constants.VcfTestResourceCertificate); v == "" {
		t.Fatal(constants.VcfTestResourceCertificate + " must be set for acceptance tests")
	}
	if v := os.Getenv(constants.VcfTestResourceCaCertificate); v == "" {
		t.Fatal(constants.VcfTestResourceCaCertificate + " must be set for acceptance tests")
	}
	if v := os.Getenv(constants.VcfTestResourceCertificateChain); v == "" {
		t.Fatal(constants.VcfTestResourceCertificateChain + " must be set for acceptance tests")
	}
}

func testAccVcfResourceExternalCertificate(domainID, resourceCert, caCert, certChain string) string {
	return fmt.Sprintf(`
	resource "vcf_external_certificate" "vcenter_cert" {
		domain_id = %q
		resource = "VCENTER"
		resource_certificate = %q
		ca_certificate = %q
		certificate_chain = %q
	}
	`,
		domainID,
		resourceCert,
		caCert,
		certChain,
	)
}
