// © Broadcom. All Rights Reserved.
// The term “Broadcom” refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/vmware/terraform-provider-vcf/internal/constants"
)

func TestAccDataSourceVcfHost(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: muxedFactories(),
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceVcfHostConfig(os.Getenv(constants.VcfTestHost1Fqdn), os.Getenv(constants.VcfTestHost1Pass)),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.vcf_host.test_host", "id"),
					resource.TestCheckResourceAttr("data.vcf_host.test_host", "fqdn", os.Getenv(constants.VcfTestHost1Fqdn)),
				),
			},
		},
	})
}

func testAccDataSourceVcfHostConfig(hostFqdn, hostPass string) string {
	return fmt.Sprintf(`
	resource "vcf_network_pool" "test_pool" {
		name    = "test_network_pool"
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

	resource "vcf_host" "test_host" {
		fqdn = %q
		username = "root"
		password = %q
		network_pool_id = vcf_network_pool.test_pool.id
		storage_type = "VSAN"
	}

	data "vcf_host" "test_host" {
		fqdn = vcf_host.test_host.fqdn
	}
    `, hostFqdn, hostPass)
}
