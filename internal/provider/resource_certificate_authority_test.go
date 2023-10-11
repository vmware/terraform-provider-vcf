/*
 *  Copyright 2023 VMware, Inc.
 *    SPDX-License-Identifier: MPL-2.0
 */

package provider

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/vmware/terraform-provider-vcf/internal/api_client"
	"github.com/vmware/vcf-sdk-go/client/certificates"
	"strings"
	"testing"
)

func TestAccResourceVcfCertificateAuthority(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		CheckDestroy:      testVerifyVcfCertificateAuthorityDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccVcfCertificateAuthorityMsft(),
				Check:  testVerifyVcfCertificateAuthorityCreate,
			},
			{
				Config: testAccVcfCertificateAuthorityOpenSsl(),
				Check:  testVerifyVcfCertificateAuthorityUpdate,
			},
		},
	})
}

func testAccVcfCertificateAuthorityMsft() string {
	return `
	resource "vcf_certificate_authority" "ca" {
  		microsoft {
			username = "Admin"
			secret = "VMwareInfra@1"
			server_url = "https://AD-vcf.eng.vmware.com/certsrv"
			template_name = "Vcms"
		}
	}`
}

func testAccVcfCertificateAuthorityOpenSsl() string {
	return `
	resource "vcf_certificate_authority" "ca" {
  		open_ssl {
			common_name = "test.openssl.eng.vmware.com"
			country = "BG"
			state = "Sofia-grad"
			locality = "Sofia"
			organization = "VMware"
			ogranization_unit = "CIBG"
		}
	}`
}

func testVerifyVcfCertificateAuthority(caType string) error {
	vcfClient := testAccProvider.Meta().(*api_client.SddcManagerClient)
	apiClient := vcfClient.ApiClient

	getCertificateAuthorityParams := &certificates.GetCertificateAuthorityByIDParams{
		ID: caType,
	}
	getCertificateAuthorityResponse, err := apiClient.Certificates.GetCertificateAuthorityByID(getCertificateAuthorityParams)
	if err != nil {
		return err
	}
	if *getCertificateAuthorityResponse.Payload.ID == caType {
		return nil
	} else {
		return fmt.Errorf("CA not the expected type: %q", caType)
	}

}

func testVerifyVcfCertificateAuthorityCreate(_ *terraform.State) error {
	return testVerifyVcfCertificateAuthority("Microsoft")
}

func testVerifyVcfCertificateAuthorityUpdate(_ *terraform.State) error {
	return testVerifyVcfCertificateAuthority("OpenSSL")
}

func testVerifyVcfCertificateAuthorityDestroy(_ *terraform.State) error {
	err := testVerifyVcfCertificateAuthority("OpenSSL")
	if !strings.Contains(err.Error(), "404") {
		return fmt.Errorf("expected CA to not be found after delete")
	}
	return nil
}
