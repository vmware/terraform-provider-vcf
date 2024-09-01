// © Broadcom. All Rights Reserved.
// The term “Broadcom” refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceCertificate(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { testAccSDDCManagerOrCloudBuilderPreCheck(t) },
		ProtoV6ProviderFactories: muxedFactories(),
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceCertificateConfig(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.vcf_certificate.cert", "certificate.0.subject_cn", "sfo-w01-vc01.sfo.rainpole.io"),
					resource.TestCheckResourceAttr("data.vcf_certificate.cert", "certificate.0.subject_locality", "Palo Alto"),
					resource.TestCheckResourceAttr("data.vcf_certificate.cert", "certificate.0.subject_st", "California"),
					resource.TestCheckResourceAttr("data.vcf_certificate.cert", "certificate.0.subject_country", "US"),
					resource.TestCheckResourceAttr("data.vcf_certificate.cert", "certificate.0.subject_org", "VMware Inc."),
					resource.TestCheckResourceAttr("data.vcf_certificate.cert", "certificate.0.subject_ou", "VCF"),
					resource.TestCheckResourceAttr("data.vcf_certificate.cert", "certificate.0.public_key_algorithm", "RSA"),
					resource.TestCheckResourceAttr("data.vcf_certificate.cert", "certificate.0.key_size", "3072"),
					resource.TestCheckResourceAttr("data.vcf_certificate.cert", "certificate.0.signature_algorithm", "SHA256withRSA"),
					resource.TestCheckResourceAttr("data.vcf_certificate.cert", "certificate.0.expiration_status", "ACTIVE"),
					resource.TestCheckResourceAttr("data.vcf_certificate.cert", "certificate.0.issued_to", "sfo-w01-vc01.sfo.rainpole.io"),
					resource.TestCheckResourceAttr("data.vcf_certificate.cert", "certificate.0.issued_by", "CN=rainpole-RPL-AD01-CA, DC=rainpole, DC=io"),
					resource.TestCheckResourceAttr("data.vcf_certificate.cert", "certificate.0.version", "V3"),
				),
			},
		},
	})
}

func testAccDataSourceCertificateConfig() string {
	return `
data "vcf_domain" "w01" {
  name = "sfo-w01"
}
data "vcf_certificate" "cert" {
  domain_id     = data.vcf_domain.w01.id
  resource_fqdn = "sfo-w01-vc01.sfo.rainpole.io"
}
`
}
