/*
 *  Copyright 2023 VMware, Inc.
 *    SPDX-License-Identifier: MPL-2.0
 */

package provider

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/vmware/terraform-provider-vcf/internal/constants"
	"github.com/vmware/vcf-sdk-go/client/domains"
	"log"
	"os"
	"testing"
)

func TestAccResourceVcfDomain(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		CheckDestroy:      testCheckVcfDomainDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccVcfDomainConfig(
					os.Getenv(constants.VcfTestHost2Fqdn),
					os.Getenv(constants.VcfTestHost2Pass),
					os.Getenv(constants.VcfTestHost3Fqdn),
					os.Getenv(constants.VcfTestHost3Pass),
					os.Getenv(constants.VcfTestHost4Fqdn),
					os.Getenv(constants.VcfTestHost4Pass),
					os.Getenv(constants.VcfTestNsxtLicenseKey),
					os.Getenv(constants.VcfTestEsxiLicenseKey),
					os.Getenv(constants.VcfTestVsanLicenseKey)),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("vcf_domain.domain1", "id"),
					resource.TestCheckResourceAttrSet("vcf_domain.domain1", "vcenter_id"),
					resource.TestCheckResourceAttrSet("vcf_domain.domain1", "vcenter_fqdn"),
					resource.TestCheckResourceAttrSet("vcf_domain.domain1", "status"),
					resource.TestCheckResourceAttrSet("vcf_domain.domain1", "type"),
					resource.TestCheckResourceAttrSet("vcf_domain.domain1", "sso_id"),
					resource.TestCheckResourceAttrSet("vcf_domain.domain1", "sso_name"),
				),
			},
		},
	})
}

func testAccVcfDomainConfig(host1Fqdn, host1SshPassword, host2Fqdn, host2SshPassword,
	host3Fqdn, host3SshPassword, nsxtLicenseKey, esxiLicenseKey, vsanLicenseKey string) string {
	return fmt.Sprintf(`
	resource "vcf_network_pool" "domain_pool" {
		name    = "engineering-pool"
		network {
			gateway   = "192.168.10.1"
			mask      = "255.255.255.0"
			mtu       = 9000
			subnet    = "192.168.10.0"
			type      = "VSAN"
			vlan_id   = 100
			ip_pools {
				start = "192.168.10.5"
				end   = "192.168.10.50"
			}
		}
		network {
			gateway   = "192.168.11.1"
			mask      = "255.255.255.0"
			mtu       = 9000
			subnet    = "192.168.11.0"
			type      = "vMotion"
			vlan_id   = 100
			ip_pools {
			  start = "192.168.11.5"
			  end   = "192.168.11.50"
			}
		  }
	}

	resource "vcf_host" "host1" {
		fqdn      = %q
		username  = "root"
		password  = %q
		network_pool_id = vcf_network_pool.domain_pool.id
		storage_type = "VSAN"
	}
	resource "vcf_host" "host2" {
		fqdn      = %q
		username  = "root"
		password  = %q
		network_pool_id = vcf_network_pool.domain_pool.id
		storage_type = "VSAN"
	}
	resource "vcf_host" "host3" {
		fqdn      = %q
		username  = "root"
		password  = %q
		network_pool_id = vcf_network_pool.domain_pool.id
		storage_type = "VSAN"
	}
	resource "vcf_domain" "domain1" {
		name                    = "test-domain"
		vcenter_name            = "test-vcenter"
        vcenter_datacenter_name = "test-datacenter"
		vcenter_root_password   = "S@mpleP@ss123!"
		vcenter_vm_size         = "tiny"
        vcenter_storage_size    = "lstorage"
		vcenter_ip_address      = "10.0.1.6"
		vcenter_subnet_mask     = "255.255.255.0"
		vcenter_gateway         = "10.0.0.250"
		vcenter_dns_name        = "test-vcenter.rainpole.io"
		nsxt_configuration {
			vip        					= "10.0.1.30"
			vip_fqdn   					= "nsx-mgmt1.rainpole.io"
			nsx_manager_admin_password	= "Nqkva_parola1"
			license_key                 = %q
			nsx_manager {
				name        = "nsxt-manager-test1"
				ip_address  = "10.0.1.40"
				dns_name    = "nsxt-manager-test1.rainpole.io"
				subnet_mask = "255.255.255.0"
				gateway     = "10.0.0.250"
			}
			nsx_manager {
				name        = "nsxt-manager-test2"
				ip_address  = "10.0.1.41"
				dns_name    = "nsxt-manager-test2.rainpole.io"
				subnet_mask = "255.255.255.0"
				gateway     = "10.0.0.250"
			}
			nsx_manager {
				name        = "nsxt-manager-test3"
				ip_address  = "10.0.1.42"
				dns_name    = "nsxt-manager-test3.rainpole.io"
				subnet_mask = "255.255.255.0"
				gateway     = "10.0.0.250"
			}
        }
		cluster {
			name = "test-cluster"
			host {
				id = vcf_host.host1.host_id
				license_key = %q
				vmnic {
					id = "vmnic0"
					vds_name = "sfo-w01-cl01-vds01"
				}
				vmnic {
					id = "vmnic1"
					vds_name = "sfo-w01-cl01-vds01"
				}
			}
			host {
				id = vcf_host.host2.host_id
				license_key = %q
				vmnic {
					id = "vmnic0"
					vds_name = "sfo-w01-cl01-vds01"
				}
				vmnic {
					id = "vmnic1"
					vds_name = "sfo-w01-cl01-vds01"
				}
			}
			host {
				id = vcf_host.host3.host_id
				license_key = %q
				vmnic {
					id = "vmnic0"
					vds_name = "sfo-w01-cl01-vds01"
				}
				vmnic {
					id = "vmnic1"
					vds_name = "sfo-w01-cl01-vds01"
				}
			}
			vds {
				name = "sfo-w01-cl01-vds01"
				portgroup {
					name = "sfo-w01-cl01-vds01-pg-mgmt"
					transport_type = "MANAGEMENT"
				}
				portgroup {
					name = "sfo-w01-cl01-vds01-pg-vsan"
					transport_type = "VSAN"
				}
				portgroup {
					name = "sfo-w01-cl01-vds01-pg-vmotion"
					transport_type = "VMOTION"
				}
			}
			vsan_datastore {
				datastore_name = "sfo-w01-cl01-ds-vsan01"
				failures_to_tolerate = 1
				license_key = %q
			}
			geneve_vlan_id = 2
		}
	}`, host1Fqdn, host1SshPassword, host2Fqdn, host2SshPassword,
		host3Fqdn, host3SshPassword, nsxtLicenseKey, esxiLicenseKey,
		esxiLicenseKey, esxiLicenseKey, vsanLicenseKey)
}

func testCheckVcfDomainDestroy(state *terraform.State) error {
	vcfClient := testAccProvider.Meta().(*SddcManagerClient)
	apiClient := vcfClient.ApiClient

	for _, rs := range state.RootModule().Resources {
		if rs.Type != "vcf_domain" {
			continue
		}

		domainId := rs.Primary.Attributes["id"]
		getDomainParams := domains.GetDomainParams{
			ID: domainId,
		}

		domainResult, err := apiClient.Domains.GetDomain(&getDomainParams)
		if err != nil {
			log.Println("error = ", err)
			return err
		}
		if domainResult.Payload != nil {
			return fmt.Errorf("domain with id %q not destroyed", domainId)
		}

	}

	// Did not find the domain
	return nil
}
