// © Broadcom. All Rights Reserved.
// The term “Broadcom” refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"errors"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccCredentialsResourcePasswordRotate(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccSDDCManagerOrCloudBuilderPreCheck(t) },
		ProtoV6ProviderFactories: muxedFactories(),
		Steps: []resource.TestStep{{
			Config: testAccResourceCredentialsPasswordRotateConfig(),
			Check:  testAccResourceCredentialsPasswordRotateCheck,
		}},
	})
}

func testAccResourceCredentialsPasswordRotateConfig() string {

	return `
		data "vcf_credentials" "sddc_creds" {
			resource_type = "VCENTER"
			account_type = "USER"
		}

		resource "vcf_credentials_rotate" "vc_0_rotate" {
			resource_name = data.vcf_credentials.sddc_creds.credentials[0].resource[0].name
			resource_type = data.vcf_credentials.sddc_creds.credentials[0].resource[0].type
			credentials {
				credential_type = data.vcf_credentials.sddc_creds.credentials[0].credential_type
				user_name = data.vcf_credentials.sddc_creds.credentials[0].user_name
			}
		}
`

}

func testAccResourceCredentialsPasswordRotateCheck(state *terraform.State) error {
	resources := state.RootModule().Resources
	dataCredentials := resources["data.vcf_credentials.sddc_creds"]
	resourceRotate := resources["vcf_credentials_rotate.vc_0_rotate"]
	passFromCredentials := dataCredentials.Primary.Attributes["credentials.0.password"]
	passFromRotation := resourceRotate.Primary.Attributes["credentials.0.password"]

	if passFromRotation == passFromCredentials {
		return errors.New("rotation is not passed, passwords are the same")
	}

	return nil
}
