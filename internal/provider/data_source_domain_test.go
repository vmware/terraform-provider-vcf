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

func testAccVcfDomainDataSourceSteps(config string) []resource.TestStep {
	return []resource.TestStep{
		{
			Config: config,
			Check: resource.ComposeTestCheckFunc(
				resource.TestCheckResourceAttrSet("data.vcf_domain.domain1", "vcenter_configuration.0.id"),
				resource.TestCheckResourceAttrSet("data.vcf_domain.domain1", "vcenter_configuration.0.fqdn"),
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
				resource.TestCheckResourceAttrSet("data.vcf_domain.domain1", "nsx_configuration.0.nsx_manager_node.0.fqdn"),
			),
		},
	}
}

func testAccVcfDomainDataSourceConfigById(domainId string) string {
	return fmt.Sprintf(`  
    data "vcf_domain" "domain1" {  
        domain_id = %q  
    }`, domainId)
}

func testAccVcfDomainDataSourceConfigByName(name string) string {
	return fmt.Sprintf(`  
    data "vcf_domain" "domain1" {  
        name = %q  
    }`, name)
}

func TestAccDataSourceVcfDomainById(t *testing.T) {
	config := testAccVcfDomainDataSourceConfigById(os.Getenv(constants.VcfTestDomainDataSourceId))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		Steps:             testAccVcfDomainDataSourceSteps(config),
	})
}

func TestAccDataSourceVcfDomainByName(t *testing.T) {
	config := testAccVcfDomainDataSourceConfigByName(os.Getenv(constants.VcfTestDomainDataSourceName))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		Steps:             testAccVcfDomainDataSourceSteps(config),
	})
}
