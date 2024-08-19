// Copyright 2023 Broadcom. All Rights Reserved.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/vmware/terraform-provider-vcf/internal/constants"
)

/*
 * To execute this test one has to log in to a VCF environment, generate a CSR for the
 * VCENTER resource, download it, paste its contents into an external CA, generate a
 * certificate, based on that CSR, copy the PEM format of the Certificate along with
 * Certificate Chain and Certificate CA and assign them to the appropriate env variables.
 */
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
					os.Getenv(constants.VcfTestResourceCaCertificate)),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("vcf_external_certificate.vcenter_cert", "certificate.0.issued_by"),
					resource.TestCheckResourceAttrSet("vcf_external_certificate.vcenter_cert", "certificate.0.issued_to"),
					resource.TestCheckResourceAttrSet("vcf_external_certificate.vcenter_cert", "certificate.0.expiration_status"),
					resource.TestCheckResourceAttrSet("vcf_external_certificate.vcenter_cert", "certificate.0.certificate_error"),
					resource.TestCheckResourceAttrSet("vcf_external_certificate.vcenter_cert", "certificate.0.key_size"),
					resource.TestCheckResourceAttrSet("vcf_external_certificate.vcenter_cert", "certificate.0.not_after"),
					resource.TestCheckResourceAttrSet("vcf_external_certificate.vcenter_cert", "certificate.0.not_before"),
					resource.TestCheckResourceAttrSet("vcf_external_certificate.vcenter_cert", "certificate.0.pem_encoded"),
					resource.TestCheckResourceAttrSet("vcf_external_certificate.vcenter_cert", "certificate.0.public_key"),
					resource.TestCheckResourceAttrSet("vcf_external_certificate.vcenter_cert", "certificate.0.public_key_algorithm"),
					resource.TestCheckResourceAttrSet("vcf_external_certificate.vcenter_cert", "certificate.0.serial_number"),
					resource.TestCheckResourceAttrSet("vcf_external_certificate.vcenter_cert", "certificate.0.signature_algorithm"),
					resource.TestCheckResourceAttrSet("vcf_external_certificate.vcenter_cert", "certificate.0.subject"),
					resource.TestCheckResourceAttrSet("vcf_external_certificate.vcenter_cert", "certificate.0.subject_alternative_name.#"),
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
}

func testAccVcfResourceExternalCertificate(domainID, resourceCert, caCert string) string {
	return fmt.Sprintf(`

	resource "vcf_external_certificate" "vcenter_cert" {
		csr_id = "csr:%s:VCENTER:some-task-id"
		resource_certificate = %q
		ca_certificate = %q
	}
	`,
		domainID,
		resourceCert,
		caCert,
	)
}
