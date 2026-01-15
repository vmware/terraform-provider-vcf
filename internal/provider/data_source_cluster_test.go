// © Broadcom. All Rights Reserved.
// The term “Broadcom” refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccDataSourceVcfCluster(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: muxedFactories(),
		Steps: []resource.TestStep{
			{
				Config: testAccVcfClusterDataSourceConfig("c1fef121-592f-45e1-9153-2d326a8f4e71"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.vcf_cluster.cluster1", "name"),
					resource.TestCheckResourceAttrSet("data.vcf_cluster.cluster1", "domain_id"),
					resource.TestCheckResourceAttrSet("data.vcf_cluster.cluster1", "primary_datastore_name"),
					resource.TestCheckResourceAttrSet("data.vcf_cluster.cluster1", "primary_datastore_type"),
					resource.TestCheckResourceAttrSet("data.vcf_cluster.cluster1", "is_default"),
					resource.TestCheckResourceAttrSet("data.vcf_cluster.cluster1", "is_stretched"),
					resource.TestCheckResourceAttrSet("data.vcf_cluster.cluster1", "host.0.id"),
					resource.TestCheckResourceAttrSet("data.vcf_cluster.cluster1", "host.0.host_name"),
					resource.TestCheckResourceAttrSet("data.vcf_cluster.cluster1", "host.0.ip_address"),
					resource.TestCheckResourceAttrSet("data.vcf_cluster.cluster1", "host.0.pnic.#"),
					resource.TestCheckResourceAttrSet("data.vcf_cluster.cluster1", "host.1.id"),
					resource.TestCheckResourceAttrSet("data.vcf_cluster.cluster1", "host.1.host_name"),
					resource.TestCheckResourceAttrSet("data.vcf_cluster.cluster1", "host.1.ip_address"),
					resource.TestCheckResourceAttrSet("data.vcf_cluster.cluster1", "host.1.pnic.#"),
					resource.TestCheckResourceAttrSet("data.vcf_cluster.cluster1", "host.2.id"),
					resource.TestCheckResourceAttrSet("data.vcf_cluster.cluster1", "host.2.host_name"),
					resource.TestCheckResourceAttrSet("data.vcf_cluster.cluster1", "host.2.ip_address"),
					resource.TestCheckResourceAttrSet("data.vcf_cluster.cluster1", "host.2.pnic.#"),
					resource.TestCheckResourceAttrSet("data.vcf_cluster.cluster1", "vds.0.name"),
					resource.TestCheckResourceAttrSet("data.vcf_cluster.cluster1", "vds.0.is_used_by_nsx"),
					resource.TestCheckResourceAttrSet("data.vcf_cluster.cluster1", "vds.0.portgroup.#"),
					resource.TestCheckResourceAttrSet("data.vcf_cluster.cluster1", "vds.0.nioc_bandwidth_allocations.#"),
				),
			},
		},
	})
}

func testAccVcfClusterDataSourceConfig(domainId string) string {
	return fmt.Sprintf(`
	data "vcf_cluster" "cluster1" {
		cluster_id = %q
	}`, domainId)
}
