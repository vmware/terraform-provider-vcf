// Copyright 2023-2024 Broadcom. All Rights Reserved.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/vmware/terraform-provider-vcf/internal/constants"
)

func TestAccDataSourceVcfNetworkPool(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: protoV6ProviderFactories(),
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
    resource "vcf_network_pool" "test_pool" {
        name    = %q
        network {
            gateway   = "192.168.4.1"
            mask      = "255.255.255.0"
            mtu       = 8940
            subnet    = "192.168.4.0"
            type      = "VSAN"
            vlan_id   = 100
            ip_pools {
                start = "192.168.4.5"
                end   = "192.168.4.50"
            }
        }
        network {
            gateway   = "192.168.5.1"
            mask      = "255.255.255.0"
            mtu       = 8940
            subnet    = "192.168.5.0"
            type      = "vMotion"
            vlan_id   = 100
			ip_pools {
				start = "192.168.5.5"
				end   = "192.168.5.50"
			}
		}
    }

    data "vcf_network_pool" "test_pool" {
        name = vcf_network_pool.test_pool.name
    }
    `, networkPoolName)
}
