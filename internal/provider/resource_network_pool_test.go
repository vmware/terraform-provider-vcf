// Copyright 2023 Broadcom. All Rights Reserved.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/vmware/terraform-provider-vcf/internal/api_client"
	"github.com/vmware/terraform-provider-vcf/internal/constants"
	"log"
	"testing"
)

func TestAccResourceVcfNetworkPool(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		CheckDestroy:      testCheckVcfNetworkPoolDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccVcfNetworkPoolConfig(constants.VcfTestNetworkPoolName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("vcf_network_pool.test_pool", "id"),
				),
			},
		},
	})
}

func testAccVcfNetworkPoolConfig(networkPoolName string) string {
	return fmt.Sprintf(`
	resource "vcf_network_pool" "test_pool" {
		name    = %q
		network {
			gateway   = "192.168.4.1"
			mask      = "255.255.255.0"
			mtu       = 9000
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
			mtu       = 9000
			subnet    = "192.168.5.0"
			type      = "vMotion"
			vlan_id   = 100
			ip_pools {
			  start = "192.168.5.5"
			  end   = "192.168.5.50"
			}
		  }
	}`, networkPoolName)
}

func testCheckVcfNetworkPoolDestroy(_ *terraform.State) error {
	vcfClient := testAccProvider.Meta().(*api_client.SddcManagerClient)
	apiClient := vcfClient.ApiClient

	hosts, err := apiClient.NetworkPools.GetNetworkPool(nil)
	if err != nil {
		log.Println("error = ", err)
		return err
	}

	for _, networkPool := range hosts.Payload.Elements {
		if networkPool.Name == constants.VcfTestNetworkPoolName {
			return fmt.Errorf("found networkPool %q", networkPool.ID)
		}
	}

	// Didn't find the networkPool
	return nil
}
