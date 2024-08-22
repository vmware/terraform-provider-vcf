// © Broadcom. All Rights Reserved.
// The term “Broadcom” refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/vmware/vcf-sdk-go/client/certificates"

	"github.com/vmware/terraform-provider-vcf/internal/api_client"
	"github.com/vmware/terraform-provider-vcf/internal/constants"
	validationUtils "github.com/vmware/terraform-provider-vcf/internal/validation"
)

func TestAccResourceVcfCertificateAuthority(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccVcfCertificateAuthorityPreCheck(t)
		},
		ProtoV6ProviderFactories: muxedFactories(),
		CheckDestroy:      testVerifyVcfCertificateAuthorityDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccVcfCertificateAuthorityMsft(),
				Check:  testVerifyVcfCertificateAuthorityCreate,
			},
			{
				ResourceName:     "vcf_certificate_authority.ca",
				ImportState:      true,
				ImportStateCheck: caImportStateCheck,
			},
			{
				Config: testAccVcfCertificateAuthorityOpenSsl(),
				Check:  testVerifyVcfCertificateAuthorityUpdate,
			},
		},
	})
}

func testAccVcfCertificateAuthorityPreCheck(t *testing.T) {
	if v := os.Getenv(constants.VcfTestMsftCaServerUrl); v == "" {
		t.Fatal(constants.VcfTestMsftCaServerUrl + " must be set for acceptance tests")
	}
	if v := os.Getenv(constants.VcfTestMsftCaUser); v == "" {
		t.Fatal(constants.VcfTestMsftCaUser + " must be set for acceptance tests")
	}
	if v := os.Getenv(constants.VcfTestMsftCaSecret); v == "" {
		t.Fatal(constants.VcfTestMsftCaSecret + " must be set for acceptance tests")
	}
}

func testAccVcfCertificateAuthorityMsft() string {
	return fmt.Sprintf(`
	resource "vcf_certificate_authority" "ca" {
  		microsoft {
			username = %q
			secret = %q
			server_url = %q
			template_name = "vcf"
		}
	}`,
		os.Getenv(constants.VcfTestMsftCaUser),
		os.Getenv(constants.VcfTestMsftCaSecret),
		os.Getenv(constants.VcfTestMsftCaServerUrl))
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
			organization_unit = "CIBG"
		}
	}`
}

func testVerifyVcfCertificateAuthority(caType string) error {
	apiClient := testAccProvider.Meta().(*api_client.SddcManagerClient).ApiClient

	getCertificateAuthorityParams := &certificates.GetCertificateAuthorityByIDParams{
		ID:      caType,
		Context: context.Background(),
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

func caImportStateCheck(states []*terraform.InstanceState) error {
	for _, state := range states {
		if state.Ephemeral.Type != "vcf_certificate_authority" {
			continue
		}
		if validationUtils.IsEmpty(state.Attributes["id"]) {
			return fmt.Errorf("CA has no id attribute set")
		}
		if state.Attributes["type"] != "Microsoft" {
			return fmt.Errorf("CA has wrong type attribute set")
		}
		if state.Attributes["microsoft.0.server_url"] != os.Getenv(constants.VcfTestMsftCaServerUrl) {
			return fmt.Errorf("CA has wrong server_url attribute set")
		}
		if state.Attributes["microsoft.0.template_name"] != "vcf" {
			return fmt.Errorf("CA has wrong template_name attribute set")
		}
		if state.Attributes["microsoft.0.username"] != os.Getenv(constants.VcfTestMsftCaUser) {
			return fmt.Errorf("CA has wrong username attribute set")
		}
		return nil
	}
	return fmt.Errorf("CA InstanceState not found! Import failed")
}
