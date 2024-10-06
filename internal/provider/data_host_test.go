// © Broadcom. All Rights Reserved.
// The term “Broadcom” refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/vmware/terraform-provider-vcf/internal/constants"
)

func TestAccDataSourceVcfHost(t *testing.T) {
	hosts := []string{
		constants.VcfTestHost1Fqdn,
		constants.VcfTestHost2Fqdn,
		constants.VcfTestHost3Fqdn,
		constants.VcfTestHost4Fqdn,
	}

	var steps []resource.TestStep
	for _, fqdn := range hosts {
		steps = append(steps, resource.TestStep{
			Config: testAccDataSourceVcfHostConfig(fqdn),
			Check: resource.ComposeTestCheckFunc(
				resource.TestCheckResourceAttrSet("data.vcf_host.test_host", "id"),
				resource.TestCheckResourceAttr("data.vcf_host.test_host", "fqdn", fqdn),
			),
		})
	}

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: muxedFactories(),
		Steps:                    steps,
	})
}

func testAccDataSourceVcfHostConfig(hostFqdn string) string {
	return fmt.Sprintf(`
    resource "vcf_host" "test_host" {
        fqdn = %q
		username = "root"
		password = "password"
        network_pool_id = "test_network_pool_id"
		storage_type = "VSAN"
    }

    data "vcf_host" "test_host" {
        fqdn = vcf_host.test_host.fqdn
    }
    `, hostFqdn)
}
