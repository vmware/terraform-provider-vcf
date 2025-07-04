// © Broadcom. All Rights Reserved.
// The term “Broadcom” refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"

	"github.com/vmware/terraform-provider-vcf/internal/api_client"
	"github.com/vmware/terraform-provider-vcf/internal/constants"
	validationUtils "github.com/vmware/terraform-provider-vcf/internal/validation"
)

func TestAccResourceVcfDomainCreate(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: muxedFactories(),
		CheckDestroy:             testCheckVcfDomainDestroy,
		Steps: []resource.TestStep{
			{
				// Initial config: 1 network pool, 3 commissioned hosts, Domain with cluster with those 3 hosts
				Config: testAccVcfDomainConfig(
					testGenerateCommissionHostConfigs(
						3,
						os.Getenv(constants.VcfTestHost2Fqdn),
						os.Getenv(constants.VcfTestHost2Pass),
						os.Getenv(constants.VcfTestHost3Fqdn),
						os.Getenv(constants.VcfTestHost3Pass),
						os.Getenv(constants.VcfTestHost4Fqdn),
						os.Getenv(constants.VcfTestHost4Pass)),
					testAccVcfClusterInDomainConfig(
						"sfo-w01-cl01",
						testGenerateHostsInClusterInDomainConfig(
							"sfo-w01-cl01",
							"host1", "host2", "host3")),
					""),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("vcf_domain.domain1", "id"),
					resource.TestCheckResourceAttrSet("vcf_domain.domain1", "vcenter_configuration.0.id"),
					resource.TestCheckResourceAttrSet("vcf_domain.domain1", "vcenter_configuration.0.fqdn"),
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
				ResourceName:     "vcf_domain.domain1",
				ImportState:      true,
				ImportStateCheck: domainImportStateCheck,
			},
		},
	})
}

func TestAccResourceVcfDomainFull(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: muxedFactories(),
		CheckDestroy:             testCheckVcfDomainDestroy,
		Steps: []resource.TestStep{
			{
				// Initial config: 1 network pool, 3 commissioned hosts, Domain with cluster with those 3 hosts
				Config: testAccVcfDomainConfig(
					testGenerateCommissionHostConfigs(
						3,
						os.Getenv(constants.VcfTestHost2Fqdn),
						os.Getenv(constants.VcfTestHost2Pass),
						os.Getenv(constants.VcfTestHost3Fqdn),
						os.Getenv(constants.VcfTestHost3Pass),
						os.Getenv(constants.VcfTestHost4Fqdn),
						os.Getenv(constants.VcfTestHost4Pass)),
					testAccVcfClusterInDomainConfig(
						"sfo-w01-cl01",
						testGenerateHostsInClusterInDomainConfig(
							"sfo-w01-cl01",
							"host1", "host2", "host3")),
					""),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("vcf_domain.domain1", "id"),
					resource.TestCheckResourceAttrSet("vcf_domain.domain1", "vcenter_configuration.0.id"),
					resource.TestCheckResourceAttrSet("vcf_domain.domain1", "vcenter_configuration.0.fqdn"),
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
					testAccVcfClusterInDomainConfig(
						"sfo-w01-cl01",
						testGenerateHostsInClusterInDomainConfig(
							"sfo-w01-cl01",
							"host1", "host2", "host3")),
					testAccVcfClusterInDomainConfig(
						"sfo-w01-cl02",
						testGenerateHostsInClusterInDomainConfig(
							"sfo-w01-cl02",
							"host4", "host5", "host6"))),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("vcf_domain.domain1", "id"),
					resource.TestCheckResourceAttrSet("vcf_domain.domain1", "vcenter_configuration.0.id"),
					resource.TestCheckResourceAttrSet("vcf_domain.domain1", "vcenter_configuration.0.fqdn"),
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
					testAccVcfClusterInDomainConfig(
						"sfo-w01-cl01",
						testGenerateHostsInClusterInDomainConfig(
							"sfo-w01-cl01",
							"host1", "host2", "host3")),
					testAccVcfClusterInDomainConfig(
						"sfo-w01-cl02",
						testGenerateHostsInClusterInDomainConfig(
							"sfo-w01-cl02",
							"host4", "host5", "host6", "host7"))),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("vcf_domain.domain1", "id"),
					resource.TestCheckResourceAttrSet("vcf_domain.domain1", "vcenter_configuration.0.id"),
					resource.TestCheckResourceAttrSet("vcf_domain.domain1", "vcenter_configuration.0.fqdn"),
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
			{
				//remove  additional host in the second cluster in the domain
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
					testAccVcfClusterInDomainConfig(
						"sfo-w01-cl01",
						testGenerateHostsInClusterInDomainConfig(
							"sfo-w01-cl01",
							"host1", "host2", "host3")),
					testAccVcfClusterInDomainConfig(
						"sfo-w01-cl02",
						testGenerateHostsInClusterInDomainConfig(
							"sfo-w01-cl02",
							"host4", "host5", "host6"))),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("vcf_domain.domain1", "id"),
					resource.TestCheckResourceAttrSet("vcf_domain.domain1", "vcenter_configuration.0.id"),
					resource.TestCheckResourceAttrSet("vcf_domain.domain1", "vcenter_configuration.0.fqdn"),
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
					resource.TestCheckNoResourceAttr("vcf_domain.domain1", "cluster.1.host.3.id"),
				),
			},
			{
				// Remove additional cluster
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
					testAccVcfClusterInDomainConfig(
						"sfo-w01-cl01",
						testGenerateHostsInClusterInDomainConfig(
							"sfo-w01-cl01",
							"host1", "host2", "host3")),
					""),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("vcf_domain.domain1", "id"),
					resource.TestCheckResourceAttrSet("vcf_domain.domain1", "vcenter_configuration.0.id"),
					resource.TestCheckResourceAttrSet("vcf_domain.domain1", "vcenter_configuration.0.fqdn"),
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
					resource.TestCheckNoResourceAttr("vcf_domain.domain1", "cluster.1.id"),
				),
			},
		},
	})
}

func testAccVcfDomainConfig(commissionHostConfig,
	clusterConfig, additionalClusterConfig string) string {
	return fmt.Sprintf(`
	resource "vcf_network_pool" "domain_pool" {
		name    = "engineering-pool"
		network {
			gateway   = "192.168.10.1"
			mask      = "255.255.255.0"
			mtu       = 8940
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
			mtu       = 8940
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
		sso {
			domain_name = "acc-test.vrack.vsphere.local"
			domain_password = "S@mpleL0ngP@ss123!"
		}
		vcenter_configuration {
			name            = "test-vcenter"
			datacenter_name = "test-datacenter"
			root_password   = "S@mpleL0ngP@ss123!"
			vm_size         = "small"
			storage_size    = "lstorage"
			ip_address      = "10.0.0.143"
			subnet_mask     = "255.255.255.0"
			gateway         = "10.0.0.250"
			fqdn            = "sfo-w01-vc01.vrack.vsphere.local"
		}
		nsx_configuration {
			vip        					= "10.0.0.166"
			vip_fqdn   					= "sfo-w01-nsx01.vrack.vsphere.local"
			nsx_manager_admin_password	= "S@mpleL0ngP@ss123!"
			form_factor                 = "small"
			nsx_manager_node {
				name        = "sfo-w01-nsx01a"
				ip_address  = "10.0.0.162"
				fqdn    = "sfo-w01-nsx01a.vrack.vsphere.local"
				subnet_mask = "255.255.255.0"
				gateway     = "10.0.0.250"
			}
			nsx_manager_node {
				name        = "sfo-w01-nsx01b"
				ip_address  = "10.0.0.163"
				fqdn    = "sfo-w01-nsx01b.vrack.vsphere.local"
				subnet_mask = "255.255.255.0"
				gateway     = "10.0.0.250"
			}
			nsx_manager_node {
				name        = "sfo-w01-nsx01c"
				ip_address  = "10.0.0.164"
				fqdn    = "sfo-w01-nsx01c.vrack.vsphere.local"
				subnet_mask = "255.255.255.0"
				gateway     = "10.0.0.250"
			}
        }
		// cluster 1 config
		%s
		// cluster 2 config
		%s
	}`, commissionHostConfig, clusterConfig, additionalClusterConfig)
}

func testAccVcfClusterInDomainConfig(clusterName, hostConfig string) string {
	return fmt.Sprintf(`
		cluster {
			name = %q
			high_availability_enabled = true
			cluster_image_id = %q
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
			}
			geneve_vlan_id = 3
		}`, clusterName, os.Getenv(constants.VcfTestClusterImageId), hostConfig, clusterName, clusterName, clusterName,
		clusterName, clusterName)
}

func testGenerateHostsInClusterInDomainConfig(clusterName string, hostsRefs ...string) string {
	var result string
	for _, hostRef := range hostsRefs {
		result += "\t" + testAccVcfHostInClusterConfig(hostRef, clusterName)
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
	apiClient := testAccProvider.Meta().(*api_client.SddcManagerClient).ApiClient

	for _, rs := range state.RootModule().Resources {
		if rs.Type != "vcf_domain" {
			continue
		}

		domainId := rs.Primary.Attributes["id"]

		domainResult, err := apiClient.GetDomainWithResponse(context.TODO(), domainId)
		if err != nil {
			log.Println("error = ", err)
			return nil
		}
		if domainResult != nil && domainResult.JSON200 != nil {
			return fmt.Errorf("domain with id %q not destroyed", domainId)
		}

	}

	// Did not find the domain
	return nil
}

func domainImportStateCheck(states []*terraform.InstanceState) error {
	for _, state := range states {
		if state.Ephemeral.Type != "vcf_domain" {
			continue
		}
		if validationUtils.IsEmpty(state.Attributes["id"]) {
			return fmt.Errorf("domain has no id attribute set")
		}
		if state.Attributes["name"] != "sfo-w01-vc01" {
			return fmt.Errorf("domain has wrong name attribute set")
		}
		if state.Attributes["vcenter_configuration.0.fqdn"] != "sfo-w01-vc01.vrack.vsphere.local" {
			return fmt.Errorf("domain has wrong name attribute set")
		}
		if state.Attributes["cluster.0.name"] != "sfo-w01-cl01" {
			return fmt.Errorf("domain has wrong cluster.0.name attribute set")
		}
		if validationUtils.IsEmpty(state.Attributes["vcenter_configuration.0.id"]) {
			return fmt.Errorf("domain has no vcenter_configuration.0.id attribute set")
		}
		if validationUtils.IsEmpty(state.Attributes["cluster.0.host.0.id"]) {
			return fmt.Errorf("domain has no cluster.0.host.0.id attribute set")
		}
		if validationUtils.IsEmpty(state.Attributes["cluster.0.host.1.id"]) {
			return fmt.Errorf("domain has no cluster.0.host.1.id attribute set")
		}
		if validationUtils.IsEmpty(state.Attributes["cluster.0.host.2.id"]) {
			return fmt.Errorf("domain has no cluster.0.host.2.id attribute set")
		}
		if validationUtils.IsEmpty(state.Attributes["cluster.0.host.0.ip_address"]) {
			return fmt.Errorf("domain has no cluster.0.host.0.ip_address attribute set")
		}
		if validationUtils.IsEmpty(state.Attributes["cluster.0.host.1.ip_address"]) {
			return fmt.Errorf("domain has no cluster.0.host.1.ip_address attribute set")
		}
		if validationUtils.IsEmpty(state.Attributes["cluster.0.host.2.ip_address"]) {
			return fmt.Errorf("domain has no cluster.0.host.2.ip_address attribute set")
		}
		if validationUtils.IsEmpty(state.Attributes["cluster.0.host.0.host_name"]) {
			return fmt.Errorf("domain has no cluster.0.host.0.host_name attribute set")
		}
		if validationUtils.IsEmpty(state.Attributes["cluster.0.host.1.host_name"]) {
			return fmt.Errorf("domain has no cluster.0.host.1.host_name attribute set")
		}
		if validationUtils.IsEmpty(state.Attributes["cluster.0.host.2.host_name"]) {
			return fmt.Errorf("domain has no cluster.0.host.2.host_name attribute set")
		}
		if validationUtils.IsEmpty(state.Attributes["nsx_configuration.0.id"]) {
			return fmt.Errorf("domain has no nsx_configuration.0.id attribute set")
		}
		if validationUtils.IsEmpty(state.Attributes["nsx_configuration.0.vip"]) {
			return fmt.Errorf("domain has no nsx_configuration.0.vip attribute set")
		}
		if validationUtils.IsEmpty(state.Attributes["nsx_configuration.0.vip_fqdn"]) {
			return fmt.Errorf("domain has no nsx_configuration.0.vip_fqdn attribute set")
		}
		if validationUtils.IsEmpty(state.Attributes["nsx_configuration.0.nsx_manager_node.0.name"]) {
			return fmt.Errorf("domain has no nsx_configuration.0.nsx_manager_node.0.name attribute set")
		}
		if validationUtils.IsEmpty(state.Attributes["nsx_configuration.0.nsx_manager_node.0.ip_address"]) {
			return fmt.Errorf("domain has no nsx_configuration.0.nsx_manager_node.0.ip_address attribute set")
		}
		if validationUtils.IsEmpty(state.Attributes["nsx_configuration.0.nsx_manager_node.0.fqdn"]) {
			return fmt.Errorf("domain has no nsx_configuration.0.nsx_manager_node.0.fqdn attribute set")
		}
		return nil
	}
	return fmt.Errorf("domain InstanceState not found! Import failed")
}
