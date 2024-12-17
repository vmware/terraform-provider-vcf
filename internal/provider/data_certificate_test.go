// © Broadcom. All Rights Reserved.
// The term “Broadcom” refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceCertificate(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccSDDCManagerOrCloudBuilderPreCheck(t) },
		ProtoV6ProviderFactories: muxedFactories(),
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceCertificateConfig(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.vcf_certificate.cert", "certificate.0.subject_cn", "sddc-manager.vrack.vsphere.local"),
					resource.TestCheckResourceAttr("data.vcf_certificate.cert", "certificate.0.subject_locality", "Palo Alto"),
					resource.TestCheckResourceAttr("data.vcf_certificate.cert", "certificate.0.subject_st", "California"),
					resource.TestCheckResourceAttr("data.vcf_certificate.cert", "certificate.0.subject_country", "US"),
					resource.TestCheckResourceAttr("data.vcf_certificate.cert", "certificate.0.subject_org", "VMware"),
					resource.TestCheckResourceAttr("data.vcf_certificate.cert", "certificate.0.subject_ou", "VMware Engineering"),
					resource.TestCheckResourceAttr("data.vcf_certificate.cert", "certificate.0.key_size", "2048"),
					resource.TestCheckResourceAttr("data.vcf_certificate.cert", "certificate.0.signature_algorithm", "SHA256WITHRSA"),
					resource.TestCheckResourceAttr("data.vcf_certificate.cert", "certificate.0.expiration_status", "ACTIVE"),
					resource.TestCheckResourceAttr("data.vcf_certificate.cert", "certificate.0.issued_to", "sddc-manager.vrack.vsphere.local"),
					resource.TestCheckResourceAttr("data.vcf_certificate.cert", "certificate.0.issued_by", "OU=VMware Engineering, O=vcenter-1.vrack.vsphere.local, ST=California, C=US, DC=local, DC=vsphere, CN=CA"),
					resource.TestCheckResourceAttr("data.vcf_certificate.cert", "certificate.0.version", "V3"),
				),
			},
		},
	})
}

func testAccDataSourceCertificateConfig() string {
	return `
	data "vcf_domain" "w01" {
	  name = "sddcId-1001"
	}
	data "vcf_certificate" "cert" {
	  domain_id     = data.vcf_domain.w01.id
	  resource_fqdn = "sddc-manager.vrack.vsphere.local"
	}
`
}
