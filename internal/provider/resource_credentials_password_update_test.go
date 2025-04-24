// © Broadcom. All Rights Reserved.
// The term “Broadcom” refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/vmware/terraform-provider-vcf/internal/constants"
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
			resource_name = %[1]q
			resource_type = "ESXI"
			credentials {
				credential_type = "SSH"
				user_name = "root"
				password = %[2]q
			}
		}

		data "vcf_credentials" "esx_creds" {
			resource_type = "ESXI"
			account_type = "USER"
			resource_name = %[1]q

			depends_on = [
				vcf_credentials_update.vc_0_update
			]
		}

`, os.Getenv(constants.VcfTestHost1Fqdn), newPassword)
}
