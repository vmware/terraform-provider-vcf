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

const (
	edgeClusterName = "testCluster1"
	edgeNode1Name   = "nsxt-edge-node-3.vrack.vsphere.local"
	edgeNode2Name   = "nsxt-edge-node-4.vrack.vsphere.local"
	edgeNode3Name   = "nsxt-edge-node-5.vrack.vsphere.local"
)

// same as the "full" test but will most optional inputs omitted.
func TestAccResourceEdgeCluster_basic(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			// Create
			{
				Config: getEdgeClusterConfigBasicInitial(),
				Check: resource.ComposeTestCheckFunc(
					getEdgeClusterChecks(2)...,
				),
			},
			// Update
			// Expand the cluster with an additional node
			{
				Config: getEdgeClusterConfigBasicExpansion(),
				Check: resource.ComposeTestCheckFunc(
					getEdgeClusterChecks(3)...,
				),
			},
			// Update
			// Shrink the cluster to its original set of nodes
			{
				Config: getEdgeClusterConfigBasicInitial(),
				Check: resource.ComposeTestCheckFunc(
					getEdgeClusterChecks(2)...,
				),
			},
		},
	})
}

func TestAccResourceEdgeCluster_full(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			// Create
			{
				Config: getEdgeClusterConfigFullInitial(),
				Check: resource.ComposeTestCheckFunc(
					getEdgeClusterChecks(2)...,
				),
			},
			// Update
			// Expand the cluster with an additional node
			{
				Config: getEdgeClusterConfigFullExpansion(),
				Check: resource.ComposeTestCheckFunc(
					getEdgeClusterChecks(3)...,
				),
			},
			// Update
			// Shrink the cluster to its original set of nodes
			{
				Config: getEdgeClusterConfigFullInitial(),
				Check: resource.ComposeTestCheckFunc(
					getEdgeClusterChecks(2)...,
				),
			},
		},
	})
}

func getEdgeClusterConfigFullInitial() string {
	edgeNode1 := getEdgeNodeConfigFull(
		edgeNode1Name,
		"10.0.0.52/24",
		"192.168.52.12/24",
		"192.168.52.13/24",
		"192.168.18.2/24",
		"192.168.19.2/24")
	edgeNode2 := getEdgeNodeConfigFull(
		edgeNode2Name,
		"10.0.0.53/24",
		"192.168.52.14/24",
		"192.168.52.15/24",
		"192.168.18.3/24",
		"192.168.19.3/24")

	return fmt.Sprintf(`
		resource "vcf_edge_cluster" "testCluster1" {
			name      = %q
			root_password = %q
			admin_password = %q
			audit_password = %q
			tier0_name = "T0_testCluster1"
			tier1_name = "T1_testCluster1"
			form_factor = "MEDIUM"
			profile_type = "DEFAULT"
			routing_type = "EBGP"
			high_availability = "ACTIVE_ACTIVE"
			mtu = 8940
			asn = 65004
			%s
			%s
		}
		`,
		edgeClusterName,
		os.Getenv(constants.VcfTestEdgeClusterRootPass),
		os.Getenv(constants.VcfTestEdgeClusterAdminPass),
		os.Getenv(constants.VcfTestEdgeClusterAuditPass),
		edgeNode1,
		edgeNode2)
}

func getEdgeClusterConfigFullExpansion() string {
	edgeNode1 := getEdgeNodeConfigFull(
		edgeNode1Name,
		"10.0.0.52/24",
		"192.168.52.12/24",
		"192.168.52.13/24",
		"192.168.18.2/24",
		"192.168.19.2/24")
	edgeNode2 := getEdgeNodeConfigFull(
		edgeNode2Name,
		"10.0.0.53/24",
		"192.168.52.14/24",
		"192.168.52.15/24",
		"192.168.18.3/24",
		"192.168.19.3/24")
	edgeNode3 := getEdgeNodeConfigFull(
		edgeNode3Name,
		"10.0.0.54/24",
		"192.168.52.16/24",
		"192.168.52.17/24",
		"192.168.18.6/24",
		"192.168.19.6/24")

	return fmt.Sprintf(`
		resource "vcf_edge_cluster" "testCluster1" {
			name      = %q
			root_password = %q
			admin_password = %q
			audit_password = %q
			tier0_name = "T0_testCluster1"
			tier1_name = "T1_testCluster1"
			form_factor = "MEDIUM"
			profile_type = "DEFAULT"
			routing_type = "EBGP"
			high_availability = "ACTIVE_ACTIVE"
			mtu = 8940
			asn = 65004
			%s
			%s
			%s
		}
		`,
		edgeClusterName,
		os.Getenv(constants.VcfTestEdgeClusterRootPass),
		os.Getenv(constants.VcfTestEdgeClusterAdminPass),
		os.Getenv(constants.VcfTestEdgeClusterAuditPass),
		edgeNode1,
		edgeNode2,
		edgeNode3)
}

func getEdgeClusterConfigBasicInitial() string {
	edgeNode1 := getEdgeNodeConfigBasic(
		edgeNode1Name,
		"10.0.0.52/24",
		"192.168.52.12/24",
		"192.168.52.13/24",
		"192.168.18.2/24",
		"192.168.19.2/24")
	edgeNode2 := getEdgeNodeConfigBasic(
		edgeNode2Name,
		"10.0.0.53/24",
		"192.168.52.14/24",
		"192.168.52.15/24",
		"192.168.18.3/24",
		"192.168.19.3/24")

	return fmt.Sprintf(`
		resource "vcf_edge_cluster" "testCluster1" {
			name      = %q
			root_password = %q
			admin_password = %q
			audit_password = %q
			form_factor = "MEDIUM"
			profile_type = "DEFAULT"
			mtu = 8940
			%s
			%s
		}
		`,
		edgeClusterName,
		os.Getenv(constants.VcfTestEdgeClusterRootPass),
		os.Getenv(constants.VcfTestEdgeClusterAdminPass),
		os.Getenv(constants.VcfTestEdgeClusterAuditPass),
		edgeNode1,
		edgeNode2)
}

func getEdgeClusterConfigBasicExpansion() string {
	edgeNode1 := getEdgeNodeConfigBasic(
		edgeNode1Name,
		"10.0.0.52/24",
		"192.168.52.12/24",
		"192.168.52.13/24",
		"192.168.18.2/24",
		"192.168.19.2/24")
	edgeNode2 := getEdgeNodeConfigBasic(
		edgeNode2Name,
		"10.0.0.53/24",
		"192.168.52.14/24",
		"192.168.52.15/24",
		"192.168.18.3/24",
		"192.168.19.3/24")
	edgeNode3 := getEdgeNodeConfigBasic(
		edgeNode3Name,
		"10.0.0.54/24",
		"192.168.52.16/24",
		"192.168.52.17/24",
		"192.168.18.6/24",
		"192.168.19.6/24")

	return fmt.Sprintf(`
		resource "vcf_edge_cluster" "testCluster1" {
			name      = %q
			root_password = %q
			admin_password = %q
			audit_password = %q
			form_factor = "MEDIUM"
			mtu = 8940
			asn = 65004
			%s
			%s
			%s
		}
		`,
		edgeClusterName,
		os.Getenv(constants.VcfTestEdgeClusterRootPass),
		os.Getenv(constants.VcfTestEdgeClusterAdminPass),
		os.Getenv(constants.VcfTestEdgeClusterAuditPass),
		edgeNode1,
		edgeNode2,
		edgeNode3)
}

func getEdgeNodeConfigFull(name, ip, tep1, tep2, uplink1, uplink2 string) string {
	return fmt.Sprintf(`
		edge_node {
			name = %q
			compute_cluster_id = %q
			root_password = %q
			admin_password = %q
			audit_password = %q
			management_ip = %q
			management_gateway = "10.0.0.250"
			tep_gateway = "192.168.52.1"
			tep1_ip = %q
			tep2_ip = %q
			tep_vlan = 1252
			inter_rack_cluster = false

			uplink {
				vlan = 2083
				interface_ip = %q
				bgp_peer {
					ip = "192.168.18.10/24"
					password = "VMware1!"
					asn = 65001
				}
			}

			uplink {
				vlan = 2084
				interface_ip = %q
				bgp_peer {
					ip = "192.168.19.10/24"
					password = "VMware1!"
					asn = 65001
				}
			}
		}
		`,
		name,
		os.Getenv(constants.VcfTestComputeClusterId),
		os.Getenv(constants.VcfTestEdgeNodeRootPass),
		os.Getenv(constants.VcfTestEdgeNodeAdminPass),
		os.Getenv(constants.VcfTestEdgeNodeAuditPass),
		ip,
		tep1,
		tep2,
		uplink1,
		uplink2)
}

func getEdgeNodeConfigBasic(name, ip, tep1, tep2, uplink1, uplink2 string) string {
	return fmt.Sprintf(`
		edge_node {
			name = %q
			compute_cluster_name = %q
			root_password = %q
			admin_password = %q
			audit_password = %q
			management_ip = %q
			management_gateway = "10.0.0.250"
			tep_gateway = "192.168.52.1"
			tep1_ip = %q
			tep2_ip = %q
			tep_vlan = 1252
			inter_rack_cluster = false

			uplink {
				vlan = 2083
				interface_ip = %q
			}

			uplink {
				vlan = 2084
				interface_ip = %q
			}
		}
		`,
		name,
		os.Getenv(constants.VcfTestComputeClusterName),
		os.Getenv(constants.VcfTestEdgeNodeRootPass),
		os.Getenv(constants.VcfTestEdgeNodeAdminPass),
		os.Getenv(constants.VcfTestEdgeNodeAuditPass),
		ip,
		tep1,
		tep2,
		uplink1,
		uplink2)
}

func getEdgeClusterChecks(numNodes int) []resource.TestCheckFunc {
	checks := []resource.TestCheckFunc{
		resource.TestCheckResourceAttrSet("vcf_edge_cluster.testCluster1", "id"),
		resource.TestCheckResourceAttrSet("vcf_edge_cluster.testCluster1", "name"),
		resource.TestCheckResourceAttrSet("vcf_edge_cluster.testCluster1", "root_password"),
		resource.TestCheckResourceAttrSet("vcf_edge_cluster.testCluster1", "admin_password"),
		resource.TestCheckResourceAttrSet("vcf_edge_cluster.testCluster1", "audit_password"),
		resource.TestCheckResourceAttrSet("vcf_edge_cluster.testCluster1", "tier0_name"),
		resource.TestCheckResourceAttrSet("vcf_edge_cluster.testCluster1", "tier1_name"),
		resource.TestCheckResourceAttrSet("vcf_edge_cluster.testCluster1", "form_factor"),
		resource.TestCheckResourceAttrSet("vcf_edge_cluster.testCluster1", "profile_type"),
		resource.TestCheckResourceAttrSet("vcf_edge_cluster.testCluster1", "routing_type"),
		resource.TestCheckResourceAttrSet("vcf_edge_cluster.testCluster1", "high_availability"),
		resource.TestCheckResourceAttrSet("vcf_edge_cluster.testCluster1", "mtu"),
	}

	for i := 0; i < numNodes; i++ {
		checks = append(checks, getEdgeNodeChecks(i)...)
	}

	return checks
}

func getEdgeNodeChecks(i int) []resource.TestCheckFunc {
	return []resource.TestCheckFunc{
		resource.TestCheckResourceAttrSet("vcf_edge_cluster.testCluster1", fmt.Sprintf("edge_node.%d.name", i)),
		resource.TestCheckResourceAttrSet("vcf_edge_cluster.testCluster1", fmt.Sprintf("edge_node.%d.compute_cluster_id", i)),
		resource.TestCheckResourceAttrSet("vcf_edge_cluster.testCluster1", fmt.Sprintf("edge_node.%d.root_password", i)),
		resource.TestCheckResourceAttrSet("vcf_edge_cluster.testCluster1", fmt.Sprintf("edge_node.%d.admin_password", i)),
		resource.TestCheckResourceAttrSet("vcf_edge_cluster.testCluster1", fmt.Sprintf("edge_node.%d.audit_password", i)),
		resource.TestCheckResourceAttrSet("vcf_edge_cluster.testCluster1", fmt.Sprintf("edge_node.%d.management_ip", i)),
		resource.TestCheckResourceAttrSet("vcf_edge_cluster.testCluster1", fmt.Sprintf("edge_node.%d.management_gateway", i)),
		resource.TestCheckResourceAttrSet("vcf_edge_cluster.testCluster1", fmt.Sprintf("edge_node.%d.tep_gateway", i)),
		resource.TestCheckResourceAttrSet("vcf_edge_cluster.testCluster1", fmt.Sprintf("edge_node.%d.tep1_ip", i)),
		resource.TestCheckResourceAttrSet("vcf_edge_cluster.testCluster1", fmt.Sprintf("edge_node.%d.tep2_ip", i)),
		resource.TestCheckResourceAttrSet("vcf_edge_cluster.testCluster1", fmt.Sprintf("edge_node.%d.tep_vlan", i)),
		resource.TestCheckResourceAttrSet("vcf_edge_cluster.testCluster1", fmt.Sprintf("edge_node.%d.inter_rack_cluster", i)),

		resource.TestCheckResourceAttrSet("vcf_edge_cluster.testCluster1", fmt.Sprintf("edge_node.%d.uplink.0.interface_ip", i)),
		resource.TestCheckResourceAttrSet("vcf_edge_cluster.testCluster1", fmt.Sprintf("edge_node.%d.uplink.0.vlan", i)),
		resource.TestCheckResourceAttrSet("vcf_edge_cluster.testCluster1", fmt.Sprintf("edge_node.%d.uplink.0.bgp_peer.0.ip", i)),
		resource.TestCheckResourceAttrSet("vcf_edge_cluster.testCluster1", fmt.Sprintf("edge_node.%d.uplink.0.bgp_peer.0.password", i)),
		resource.TestCheckResourceAttrSet("vcf_edge_cluster.testCluster1", fmt.Sprintf("edge_node.%d.uplink.0.bgp_peer.0.asn", i)),

		resource.TestCheckResourceAttrSet("vcf_edge_cluster.testCluster1", fmt.Sprintf("edge_node.%d.uplink.1.interface_ip", i)),
		resource.TestCheckResourceAttrSet("vcf_edge_cluster.testCluster1", fmt.Sprintf("edge_node.%d.uplink.1.vlan", i)),
		resource.TestCheckResourceAttrSet("vcf_edge_cluster.testCluster1", fmt.Sprintf("edge_node.%d.uplink.1.bgp_peer.0.ip", i)),
		resource.TestCheckResourceAttrSet("vcf_edge_cluster.testCluster1", fmt.Sprintf("edge_node.%d.uplink.1.bgp_peer.0.password", i)),
		resource.TestCheckResourceAttrSet("vcf_edge_cluster.testCluster1", fmt.Sprintf("edge_node.%d.uplink.1.bgp_peer.0.asn", i)),
	}
}
