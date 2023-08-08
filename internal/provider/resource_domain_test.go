/*
 *  Copyright 2023 VMware, Inc.
 *    SPDX-License-Identifier: MPL-2.0
 */

package provider

import (
	"context"
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
					testGenerateCommissionHostConfigs(
						3,
						os.Getenv(constants.VcfTestHost2Fqdn),
						os.Getenv(constants.VcfTestHost2Pass),
						os.Getenv(constants.VcfTestHost3Fqdn),
						os.Getenv(constants.VcfTestHost3Pass),
						os.Getenv(constants.VcfTestHost4Fqdn),
						os.Getenv(constants.VcfTestHost4Pass)),
					os.Getenv(constants.VcfTestNsxLicenseKey),
					testAccVcfClusterInDomainConfig(
						"sfo-w01-cl01",
						testGenerateHostsInClusterInDomainConfig(
							os.Getenv(constants.VcfTestEsxiLicenseKey),
							"sfo-w01-cl01",
							"host1", "host2", "host3"),
						os.Getenv(constants.VcfTestVsanLicenseKey)),
					""),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("vcf_domain.domain1", "id"),
					resource.TestCheckResourceAttrSet("vcf_domain.domain1", "vcenter.0.id"),
					resource.TestCheckResourceAttrSet("vcf_domain.domain1", "vcenter.0.fqdn"),
					resource.TestCheckResourceAttrSet("vcf_domain.domain1", "status"),
					resource.TestCheckResourceAttrSet("vcf_domain.domain1", "type"),
					resource.TestCheckResourceAttrSet("vcf_domain.domain1", "sso_id"),
					resource.TestCheckResourceAttrSet("vcf_domain.domain1", "sso_name"),
					resource.TestCheckResourceAttrSet("vcf_domain.domain1", "cluster.0.id"),
					resource.TestCheckResourceAttrSet("vcf_domain.domain1", "cluster.0.name"),
					resource.TestCheckResourceAttrSet("vcf_domain.domain1", "cluster.0.primary_datastore_name"),
					resource.TestCheckResourceAttrSet("vcf_domain.domain1", "cluster.0.primary_datastore_type"),
					resource.TestCheckResourceAttrSet("vcf_domain.domain1", "cluster.0.is_default"),
					resource.TestCheckResourceAttrSet("vcf_domain.domain1", "cluster.0.is_stretched"),
					resource.TestCheckResourceAttrSet("vcf_domain.domain1", "cluster.0.host.0.id"),
					resource.TestCheckResourceAttrSet("vcf_domain.domain1", "cluster.0.host.1.id"),
					resource.TestCheckResourceAttrSet("vcf_domain.domain1", "cluster.0.host.2.id"),
				),
			},
			{
				// add second cluster inside the domain
				Config: testAccVcfDomainConfig(
					testGenerateCommissionHostConfigs(
						6,
						os.Getenv(constants.VcfTestHost2Fqdn),
						os.Getenv(constants.VcfTestHost2Pass),
						os.Getenv(constants.VcfTestHost3Fqdn),
						os.Getenv(constants.VcfTestHost3Pass),
						os.Getenv(constants.VcfTestHost4Fqdn),
						os.Getenv(constants.VcfTestHost4Pass),
						os.Getenv(constants.VcfTestHost5Fqdn),
						os.Getenv(constants.VcfTestHost5Pass),
						os.Getenv(constants.VcfTestHost6Fqdn),
						os.Getenv(constants.VcfTestHost6Pass),
						os.Getenv(constants.VcfTestHost7Fqdn),
						os.Getenv(constants.VcfTestHost7Pass)),
					os.Getenv(constants.VcfTestNsxLicenseKey),
					testAccVcfClusterInDomainConfig(
						"sfo-w01-cl01",
						testGenerateHostsInClusterInDomainConfig(
							os.Getenv(constants.VcfTestEsxiLicenseKey),
							"sfo-w01-cl01",
							"host1", "host2", "host3"),
						os.Getenv(constants.VcfTestVsanLicenseKey)),
					testAccVcfClusterInDomainConfig(
						"sfo-w01-cl02",
						testGenerateHostsInClusterInDomainConfig(
							os.Getenv(constants.VcfTestEsxiLicenseKey),
							"sfo-w01-cl02",
							"host4", "host5", "host6"),
						os.Getenv(constants.VcfTestVsanLicenseKey))),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("vcf_domain.domain1", "id"),
					resource.TestCheckResourceAttrSet("vcf_domain.domain1", "vcenter.0.id"),
					resource.TestCheckResourceAttrSet("vcf_domain.domain1", "vcenter.0.fqdn"),
					resource.TestCheckResourceAttrSet("vcf_domain.domain1", "status"),
					resource.TestCheckResourceAttrSet("vcf_domain.domain1", "type"),
					resource.TestCheckResourceAttrSet("vcf_domain.domain1", "sso_id"),
					resource.TestCheckResourceAttrSet("vcf_domain.domain1", "sso_name"),
					resource.TestCheckResourceAttrSet("vcf_domain.domain1", "cluster.0.id"),
					resource.TestCheckResourceAttrSet("vcf_domain.domain1", "cluster.0.name"),
					resource.TestCheckResourceAttrSet("vcf_domain.domain1", "cluster.0.primary_datastore_name"),
					resource.TestCheckResourceAttrSet("vcf_domain.domain1", "cluster.0.primary_datastore_type"),
					resource.TestCheckResourceAttrSet("vcf_domain.domain1", "cluster.0.is_default"),
					resource.TestCheckResourceAttrSet("vcf_domain.domain1", "cluster.0.is_stretched"),
					resource.TestCheckResourceAttrSet("vcf_domain.domain1", "cluster.0.host.0.id"),
					resource.TestCheckResourceAttrSet("vcf_domain.domain1", "cluster.0.host.1.id"),
					resource.TestCheckResourceAttrSet("vcf_domain.domain1", "cluster.0.host.2.id"),
					resource.TestCheckResourceAttrSet("vcf_domain.domain1", "cluster.1.id"),
					resource.TestCheckResourceAttrSet("vcf_domain.domain1", "cluster.1.name"),
					resource.TestCheckResourceAttrSet("vcf_domain.domain1", "cluster.1.primary_datastore_name"),
					resource.TestCheckResourceAttrSet("vcf_domain.domain1", "cluster.1.primary_datastore_type"),
					resource.TestCheckResourceAttrSet("vcf_domain.domain1", "cluster.1.is_default"),
					resource.TestCheckResourceAttrSet("vcf_domain.domain1", "cluster.1.is_stretched"),
					resource.TestCheckResourceAttrSet("vcf_domain.domain1", "cluster.1.host.0.id"),
					resource.TestCheckResourceAttrSet("vcf_domain.domain1", "cluster.1.host.1.id"),
					resource.TestCheckResourceAttrSet("vcf_domain.domain1", "cluster.1.host.2.id"),
				),
			},
			{
				// add additional host in the second cluster in the domain
				Config: testAccVcfDomainConfig(
					testGenerateCommissionHostConfigs(
						7,
						os.Getenv(constants.VcfTestHost2Fqdn),
						os.Getenv(constants.VcfTestHost2Pass),
						os.Getenv(constants.VcfTestHost3Fqdn),
						os.Getenv(constants.VcfTestHost3Pass),
						os.Getenv(constants.VcfTestHost4Fqdn),
						os.Getenv(constants.VcfTestHost4Pass),
						os.Getenv(constants.VcfTestHost5Fqdn),
						os.Getenv(constants.VcfTestHost5Pass),
						os.Getenv(constants.VcfTestHost6Fqdn),
						os.Getenv(constants.VcfTestHost6Pass),
						os.Getenv(constants.VcfTestHost7Fqdn),
						os.Getenv(constants.VcfTestHost7Pass),
						os.Getenv(constants.VcfTestHost8Fqdn),
						os.Getenv(constants.VcfTestHost8Pass)),
					os.Getenv(constants.VcfTestNsxLicenseKey),
					testAccVcfClusterInDomainConfig(
						"sfo-w01-cl01",
						testGenerateHostsInClusterInDomainConfig(
							os.Getenv(constants.VcfTestEsxiLicenseKey),
							"sfo-w01-cl01",
							"host1", "host2", "host3"),
						os.Getenv(constants.VcfTestVsanLicenseKey)),
					testAccVcfClusterInDomainConfig(
						"sfo-w01-cl02",
						testGenerateHostsInClusterInDomainConfig(
							os.Getenv(constants.VcfTestEsxiLicenseKey),
							"sfo-w01-cl02",
							"host4", "host5", "host6", "host7"),
						os.Getenv(constants.VcfTestVsanLicenseKey))),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("vcf_domain.domain1", "id"),
					resource.TestCheckResourceAttrSet("vcf_domain.domain1", "vcenter.0.id"),
					resource.TestCheckResourceAttrSet("vcf_domain.domain1", "vcenter.0.fqdn"),
					resource.TestCheckResourceAttrSet("vcf_domain.domain1", "status"),
					resource.TestCheckResourceAttrSet("vcf_domain.domain1", "type"),
					resource.TestCheckResourceAttrSet("vcf_domain.domain1", "sso_id"),
					resource.TestCheckResourceAttrSet("vcf_domain.domain1", "sso_name"),
					resource.TestCheckResourceAttrSet("vcf_domain.domain1", "cluster.0.id"),
					resource.TestCheckResourceAttrSet("vcf_domain.domain1", "cluster.0.name"),
					resource.TestCheckResourceAttrSet("vcf_domain.domain1", "cluster.0.primary_datastore_name"),
					resource.TestCheckResourceAttrSet("vcf_domain.domain1", "cluster.0.primary_datastore_type"),
					resource.TestCheckResourceAttrSet("vcf_domain.domain1", "cluster.0.is_default"),
					resource.TestCheckResourceAttrSet("vcf_domain.domain1", "cluster.0.is_stretched"),
					resource.TestCheckResourceAttrSet("vcf_domain.domain1", "cluster.0.host.0.id"),
					resource.TestCheckResourceAttrSet("vcf_domain.domain1", "cluster.0.host.1.id"),
					resource.TestCheckResourceAttrSet("vcf_domain.domain1", "cluster.0.host.2.id"),
					resource.TestCheckResourceAttrSet("vcf_domain.domain1", "cluster.1.id"),
					resource.TestCheckResourceAttrSet("vcf_domain.domain1", "cluster.1.name"),
					resource.TestCheckResourceAttrSet("vcf_domain.domain1", "cluster.1.primary_datastore_name"),
					resource.TestCheckResourceAttrSet("vcf_domain.domain1", "cluster.1.primary_datastore_type"),
					resource.TestCheckResourceAttrSet("vcf_domain.domain1", "cluster.1.is_default"),
					resource.TestCheckResourceAttrSet("vcf_domain.domain1", "cluster.1.is_stretched"),
					resource.TestCheckResourceAttrSet("vcf_domain.domain1", "cluster.1.host.0.id"),
					resource.TestCheckResourceAttrSet("vcf_domain.domain1", "cluster.1.host.1.id"),
					resource.TestCheckResourceAttrSet("vcf_domain.domain1", "cluster.1.host.2.id"),
					resource.TestCheckResourceAttrSet("vcf_domain.domain1", "cluster.1.host.3.id"),
				),
			},
		},
	})
}

func testAccVcfDomainConfig(commissionHostConfig, nsxLicenseKey,
	clusterConfig, additionalClusterConfig string) string {
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

	// Host commission configs
	%s

	resource "vcf_domain" "domain1" {
		name                    = "sfo-w01-vc01"
		vcenter {
			name            = "test-vcenter"
			datacenter_name = "test-datacenter"
			root_password   = "S@mpleP@ss123!"
			vm_size         = "medium"
			storage_size    = "lstorage"
			ip_address      = "10.0.0.43"
			subnet_mask     = "255.255.255.0"
			gateway         = "10.0.0.250"
			dns_name        = "sfo-w01-vc01.sfo.rainpole.io"
		}
		nsx_configuration {
			vip        					= "10.0.0.66"
			vip_fqdn   					= "sfo-w01-nsx01.sfo.rainpole.io"
			nsx_manager_admin_password	= "Nqkva_parola1"
			form_factor                 = "small"
			license_key                 = %q
			nsx_manager_node {
				name        = "sfo-w01-nsx01a"
				ip_address  = "10.0.0.62"
				dns_name    = "sfo-w01-nsx01a.sfo.rainpole.io"
				subnet_mask = "255.255.255.0"
				gateway     = "10.0.0.250"
			}
			nsx_manager_node {
				name        = "sfo-w01-nsx01b"
				ip_address  = "10.0.0.63"
				dns_name    = "sfo-w01-nsx01b.sfo.rainpole.io"
				subnet_mask = "255.255.255.0"
				gateway     = "10.0.0.250"
			}
			nsx_manager_node {
				name        = "sfo-w01-nsx01c"
				ip_address  = "10.0.0.64"
				dns_name    = "sfo-w01-nsx01c.sfo.rainpole.io"
				subnet_mask = "255.255.255.0"
				gateway     = "10.0.0.250"
			}
        }
		// cluster 1 config
		%s
		// cluster 2 config
		%s
	}`, commissionHostConfig, nsxLicenseKey, clusterConfig, additionalClusterConfig)
}

func testAccVcfClusterInDomainConfig(clusterName, hostConfig, vsanLicenseKey string) string {
	return fmt.Sprintf(`
		cluster {
			name = %q
			// hosts config
			%s
			vds {
				name = "%s-vds01"
				portgroup {
					name = "%s-vds01-pg-mgmt"
					transport_type = "MANAGEMENT"
				}
				portgroup {
					name = "%s-vds01-pg-vsan"
					transport_type = "VSAN"
				}
				portgroup {
					name = "%s-vds01-pg-vmotion"
					transport_type = "VMOTION"
				}
			}
			vsan_datastore {
				datastore_name = "%s-ds-vsan01"
				failures_to_tolerate = 1
				license_key = %q
			}
			geneve_vlan_id = 3
		}`, clusterName, hostConfig, clusterName, clusterName, clusterName,
		clusterName, clusterName, vsanLicenseKey)
}

func testGenerateHostsInClusterInDomainConfig(esxLicenseKey, clusterName string, hostsRefs ...string) string {
	var result string
	for _, hostRef := range hostsRefs {
		result += "\t" + testAccVcfHostInClusterConfig(hostRef, esxLicenseKey, clusterName)
		result += "\n"
	}
	return result
}

func testGenerateCommissionHostConfigs(numberOfCommissionedHosts int, commissionHostsCredentials ...string) string {
	var result string
	for i := 0; i < numberOfCommissionedHosts; i++ {
		result += fmt.Sprintf(
			`resource "vcf_host" "host%d" {
				fqdn      = %q
				username  = "root"
				password  = %q
				network_pool_id = vcf_network_pool.domain_pool.id
				storage_type = "VSAN"
		}
		`, i+1, commissionHostsCredentials[i*2], commissionHostsCredentials[i*2+1])
	}
	return result
}

func testCheckVcfDomainDestroy(state *terraform.State) error {
	vcfClient := testAccProvider.Meta().(*SddcManagerClient)
	apiClient := vcfClient.ApiClient

	for _, rs := range state.RootModule().Resources {
		if rs.Type != "vcf_domain" {
			continue
		}

		domainId := rs.Primary.Attributes["id"]
		getDomainParams := domains.NewGetDomainParams().
			WithTimeout(constants.DefaultVcfApiCallTimeout).
			WithContext(context.TODO())
		getDomainParams.ID = domainId

		domainResult, err := apiClient.Domains.GetDomain(getDomainParams)
		if err != nil {
			log.Println("error = ", err)
			return nil
		}
		if domainResult != nil && domainResult.Payload != nil {
			return fmt.Errorf("domain with id %q not destroyed", domainId)
		}

	}

	// Did not find the domain
	return nil
}
