/* Copyright 2023 VMware, Inc.
   SPDX-License-Identifier: MPL-2.0 */

package provider

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/vmware/terraform-provider-vcf/internal/constants"
	"log"
	"os"
	"testing"
)

func TestAccResourceVcfHost(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		CheckDestroy:      testCheckVcfHostDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccVcfHostConfig(
					os.Getenv(constants.VcfTestHost1Fqdn),
					os.Getenv(constants.VcfTestHost1Pass)),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("vcf_host.host1", "host_id"),
				),
			},
		},
	})
}

func testAccVcfHostConfig(hostFqdn, hostSshPassword string) string {
	return fmt.Sprintf(`
	resource "vcf_network_pool" "eng_pool" {
		name    = "engineering-pool"
		network {
			gateway   = "192.168.8.1"
			mask      = "255.255.255.0"
			mtu       = 9000
			subnet    = "192.168.8.0"
			type      = "VSAN"
			vlan_id   = 100
			ip_pools {
				start = "192.168.8.5"
				end   = "192.168.8.50"
			}
		}
		network {
			gateway   = "192.168.9.1"
			mask      = "255.255.255.0"
			mtu       = 9000
			subnet    = "192.168.9.0"
			type      = "vMotion"
			vlan_id   = 100
			ip_pools {
			  start = "192.168.9.5"
			  end   = "192.168.9.50"
			}
		  }
	}

	resource "vcf_host" "host1" {
		fqdn      = %q
		username  = "root"
		password  = %q
		network_pool_id = vcf_network_pool.eng_pool.id
		storage_type = "VSAN"
	}`, hostFqdn, hostSshPassword)
}

func testCheckVcfHostDestroy(_ *terraform.State) error {
	vcfClient := testAccProvider.Meta().(*SddcManagerClient)
	apiClient := vcfClient.ApiClient

	hosts, err := apiClient.Hosts.GetHosts(nil)
	if err != nil {
		log.Println("error = ", err)
		return err
	}

	for _, host := range hosts.Payload.Elements {
		if host.Fqdn == os.Getenv(constants.VcfTestHost1Fqdn) {
			return fmt.Errorf("found host %q", host.ID)
		}
	}

	// Found the host
	return nil
}
