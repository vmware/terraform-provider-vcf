// Copyright 2023-2024 Broadcom. All Rights Reserved.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	"github.com/vmware/terraform-provider-vcf/internal/constants"
)

func TestAccDataSourceVcfNetworkPool(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: muxedFactories(),
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceVcfNetworkPoolConfig(constants.VcfTestNetworkPoolName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.vcf_network_pool.test_pool", "id"),
					resource.TestCheckResourceAttr("data.vcf_network_pool.test_pool", "name", constants.VcfTestNetworkPoolName),
				),
			},
		},
	})
}

func testAccDataSourceVcfNetworkPoolConfig(networkPoolName string) string {
	return fmt.Sprintf(`
    data "vcf_network_pool" "test_pool" {
        name = %q
    }
    `, networkPoolName)
}
