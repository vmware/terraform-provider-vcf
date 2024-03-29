// Copyright 2024 Broadcom. All Rights Reserved.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/vmware/terraform-provider-vcf/internal/constants"
	"os"
	"testing"
)

func TestAccClusterPersonality_basic(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPersonalityPreCheck(t) },
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: getClusterPersonalityConfig(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("vcf_cluster_personality.personality", "id"),
				),
			},
		},
	})
}

func getClusterPersonalityConfig() string {
	return fmt.Sprintf(`
		resource "vcf_cluster_personality" "personality" {
			name      = "personality1"
			cluster_id = %q
			domain_id = %q
		}
		`,
		os.Getenv(constants.VcfTestClusterId),
		os.Getenv(constants.VcfTestDomainDataSourceId))
}

// testAccPreCheck validates all required environment variables for running these acceptance
// tests are set.
func testAccPersonalityPreCheck(t *testing.T) {
	if v := os.Getenv(constants.VcfTestClusterId); v == "" {
		t.Fatalf("%s must be set for acceptance tests", constants.VcfTestClusterId)
	}
	if v := os.Getenv(constants.VcfTestDomainDataSourceId); v == "" {
		t.Fatalf("%s must be set for acceptance tests", constants.VcfTestDomainDataSourceId)
	}
}
