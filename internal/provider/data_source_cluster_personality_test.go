// © Broadcom. All Rights Reserved.
// The term “Broadcom” refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccDataSourceClusterPersonality_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: muxedFactories(),
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceClusterPersonalityConfigBasic(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.vcf_cluster_personality.personality1", "id"),
				),
			},
		},
	})
}

func testAccDataSourceClusterPersonalityConfigBasic() string {
	return `
	data "vcf_cluster_personality" "personality1" {
        # The default personality name, no need for a variable
		name = "Management-Domain-ESXi-Personality"
	}`
}
