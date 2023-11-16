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

func TestAccResourceVcfResourceCertificate(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccVcfCertificateAuthorityPreCheck(t)
		},
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccVcfResourceCertificate(
					os.Getenv(constants.VcfTestDomainDataSourceId),
					os.Getenv(constants.VcfTestMsftCaServerUrl),
					os.Getenv(constants.VcfTestMsftCaUser),
					os.Getenv(constants.VcfTestMsftCaSecret)),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("vcf_certificate.vcenter_cert", "certificate.0.issued_by"),
					resource.TestCheckResourceAttrSet("vcf_certificate.vcenter_cert", "certificate.0.issued_to"),
					resource.TestCheckResourceAttrSet("vcf_certificate.vcenter_cert", "certificate.0.expiration_status"),
					resource.TestCheckResourceAttrSet("vcf_certificate.vcenter_cert", "certificate.0.certificate_error"),
					resource.TestCheckResourceAttrSet("vcf_certificate.vcenter_cert", "certificate.0.is_installed"),
					resource.TestCheckResourceAttrSet("vcf_certificate.vcenter_cert", "certificate.0.issued_by"),
					resource.TestCheckResourceAttrSet("vcf_certificate.vcenter_cert", "certificate.0.issued_to"),
					resource.TestCheckResourceAttrSet("vcf_certificate.vcenter_cert", "certificate.0.key_size"),
					resource.TestCheckResourceAttrSet("vcf_certificate.vcenter_cert", "certificate.0.not_after"),
					resource.TestCheckResourceAttrSet("vcf_certificate.vcenter_cert", "certificate.0.not_before"),
					resource.TestCheckResourceAttrSet("vcf_certificate.vcenter_cert", "certificate.0.pem_encoded"),
					resource.TestCheckResourceAttrSet("vcf_certificate.vcenter_cert", "certificate.0.public_key"),
					resource.TestCheckResourceAttrSet("vcf_certificate.vcenter_cert", "certificate.0.public_key_algorithm"),
					resource.TestCheckResourceAttrSet("vcf_certificate.vcenter_cert", "certificate.0.serial_number"),
					resource.TestCheckResourceAttrSet("vcf_certificate.vcenter_cert", "certificate.0.signature_algorithm"),
					resource.TestCheckResourceAttrSet("vcf_certificate.vcenter_cert", "certificate.0.subject"),
					resource.TestCheckResourceAttrSet("vcf_certificate.vcenter_cert", "certificate.0.subject_alternative_name"),
					resource.TestCheckResourceAttrSet("vcf_certificate.vcenter_cert", "certificate.0.thumbprint"),
					resource.TestCheckResourceAttrSet("vcf_certificate.vcenter_cert", "certificate.0.thumbprint_algorithm"),
					resource.TestCheckResourceAttrSet("vcf_certificate.vcenter_cert", "certificate.0.version"),
					resource.TestCheckResourceAttrSet("vcf_certificate.vcenter_cert", "certificate.0.number_of_days_to_expire")),
			},
		},
	})
}

func testAccVcfResourceCertificate(domainID, msftCaServerUrl, msftCaUser, msftCaSecret string) string {
	return fmt.Sprintf(`
	resource "vcf_certificate_authority" "ca" {
  		microsoft {
			username = %q
			secret = %q
			server_url = %q
			template_name = "vcf"
		}
	}

	resource "vcf_resource_csr" "csr1" {
  		domain_id = %q
		country = "BG"
		email = "admin@vmware.com"
		key_size = "3072"
		locality = "Sofia"
		state = "Sofia-grad"
		organization = "VMware Inc."
		organization_unit = "VCF"
		resource = "VCENTER"
	}


	resource "vcf_certificate" "vcenter_cert" {
		csr_id = vcf_resource_csr.csr1.id
		ca_id = vcf_certificate_authority.ca.id
	}
	`,
		msftCaUser,
		msftCaSecret,
		msftCaServerUrl,
		domainID,
	)
}
