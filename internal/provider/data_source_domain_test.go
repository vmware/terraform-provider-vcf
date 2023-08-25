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

func TestAccDataSourceVcfDomain(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccVcfDomainDataSourceConfig(
					os.Getenv(constants.VcfTestDomainDataSourceId)),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.vcf_domain.domain1", "vcenter.0.id"),
					resource.TestCheckResourceAttrSet("data.vcf_domain.domain1", "vcenter.0.fqdn"),
					resource.TestCheckResourceAttrSet("data.vcf_domain.domain1", "status"),
					resource.TestCheckResourceAttrSet("data.vcf_domain.domain1", "type"),
					resource.TestCheckResourceAttrSet("data.vcf_domain.domain1", "sso_id"),
					resource.TestCheckResourceAttrSet("data.vcf_domain.domain1", "sso_name"),
					resource.TestCheckResourceAttrSet("data.vcf_domain.domain1", "cluster.0.id"),
					resource.TestCheckResourceAttrSet("data.vcf_domain.domain1", "cluster.0.name"),
					resource.TestCheckResourceAttrSet("data.vcf_domain.domain1", "cluster.0.primary_datastore_name"),
					resource.TestCheckResourceAttrSet("data.vcf_domain.domain1", "cluster.0.primary_datastore_type"),
					resource.TestCheckResourceAttrSet("data.vcf_domain.domain1", "cluster.0.is_default"),
					resource.TestCheckResourceAttrSet("data.vcf_domain.domain1", "cluster.0.is_stretched"),
					resource.TestCheckResourceAttrSet("data.vcf_domain.domain1", "cluster.0.host.0.id"),
					resource.TestCheckResourceAttrSet("data.vcf_domain.domain1", "cluster.0.host.1.id"),
					resource.TestCheckResourceAttrSet("data.vcf_domain.domain1", "cluster.0.host.2.id"),
					resource.TestCheckResourceAttrSet("data.vcf_domain.domain1", "cluster.0.host.3.id"),
					resource.TestCheckResourceAttrSet("data.vcf_domain.domain1", "nsx_configuration.0.id"),
					resource.TestCheckResourceAttrSet("data.vcf_domain.domain1", "nsx_configuration.0.vip"),
					resource.TestCheckResourceAttrSet("data.vcf_domain.domain1", "nsx_configuration.0.vip_fqdn"),
					resource.TestCheckResourceAttrSet("data.vcf_domain.domain1", "nsx_configuration.0.nsx_manager_node.0.name"),
					resource.TestCheckResourceAttrSet("data.vcf_domain.domain1", "nsx_configuration.0.nsx_manager_node.0.ip_address"),
					resource.TestCheckResourceAttrSet("data.vcf_domain.domain1", "nsx_configuration.0.nsx_manager_node.0.dns_name"),
				),
			},
		},
	})
}

func testAccVcfDomainDataSourceConfig(domainId string) string {
	return fmt.Sprintf(`
	data "vcf_domain" "domain1" {
		domain_id = %q
	}`, domainId)
}
