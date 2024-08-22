// © Broadcom. All Rights Reserved.
// The term “Broadcom” refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceCredentialsAll(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { testAccSDDCManagerOrCloudBuilderPreCheck(t) },
		ProtoV6ProviderFactories: muxedFactories(),
		Steps: []resource.TestStep{{
			Config: testAccDataSourceCredentialsAll(),
			Check:  resource.TestCheckResourceAttrSet("data.vcf_credentials.creds", "credentials.#"),
		}},
	})
}

func TestAccDataSourceCredentials_VC(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { testAccSDDCManagerOrCloudBuilderPreCheck(t) },
		ProtoV6ProviderFactories: muxedFactories(),
		Steps: []resource.TestStep{{
			Config: testAccDataSourceCredentialsVc(),
			Check: resource.ComposeTestCheckFunc(
				resource.TestCheckResourceAttr("data.vcf_credentials.creds", "credentials.0.credential_type", "SSH"),
				resource.TestCheckResourceAttr("data.vcf_credentials.creds", "credentials.1.credential_type", "SSO"),
			),
		}},
	})
}

func testAccDataSourceCredentialsAll() string {
	return `
	data "vcf_credentials" "creds" {

	}
`
}

func testAccDataSourceCredentialsVc() string {
	return `
	data "vcf_credentials" "creds" {
		resource_type = "VCENTER"
	}
`
}
