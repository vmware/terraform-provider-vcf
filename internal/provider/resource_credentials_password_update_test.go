// © Broadcom. All Rights Reserved.
// The term “Broadcom” refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccCredentialsResourcePasswordUpdate(t *testing.T) {
	newPassword := fmt.Sprintf("%s$1A", acctest.RandString(7))

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccSDDCManagerOrCloudBuilderPreCheck(t) },
		ProtoV6ProviderFactories: muxedFactories(),
		Steps: []resource.TestStep{{
			Config: testAccResourceCredentialsPasswordUpdateConfig(newPassword),
			Check:  resource.TestCheckResourceAttr("data.vcf_credentials.esx_creds", "credentials.0.password", newPassword),
		}},
	})
}

func testAccResourceCredentialsPasswordUpdateConfig(newPassword string) string {

	return fmt.Sprintf(`
		resource "vcf_credentials_update" "vc_0_update" {
			resource_name = "esxi-4.vrack.vsphere.local"
			resource_type = "ESXI"
			credentials {
				credential_type = "SSH"
				user_name = "root"
				password = %[1]q
			}
		}

		data "vcf_credentials" "esx_creds" {
			resource_type = "ESXI"
			account_type = "USER"
			resource_name = "esxi-4.vrack.vsphere.local"

			depends_on = [
				vcf_credentials_update.vc_0_update
			]
		}

`, newPassword)
}
