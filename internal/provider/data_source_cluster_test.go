/*
 *  Copyright 2023 VMware, Inc.
 *    SPDX-License-Identifier: MPL-2.0
 */

package provider

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/vmware/terraform-provider-vcf/internal/constants"
	"os"
	"testing"
)

func TestAccDataSourceVcfCluster(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccVcfClusterDataSourceConfig(
					os.Getenv(constants.VcfTestClusterDataSourceId)),
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
					resource.TestCheckResourceAttrSet("data.vcf_cluster.cluster1", "host.1.id"),
					resource.TestCheckResourceAttrSet("data.vcf_cluster.cluster1", "host.1.host_name"),
					resource.TestCheckResourceAttrSet("data.vcf_cluster.cluster1", "host.1.ip_address"),
					resource.TestCheckResourceAttrSet("data.vcf_cluster.cluster1", "host.2.id"),
					resource.TestCheckResourceAttrSet("data.vcf_cluster.cluster1", "host.2.host_name"),
					resource.TestCheckResourceAttrSet("data.vcf_cluster.cluster1", "host.2.ip_address"),
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
