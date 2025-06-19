// © Broadcom. All Rights Reserved.
// The term “Broadcom” refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	"github.com/vmware/terraform-provider-vcf/internal/constants"
)

func TestAccDataSourceVcfHost(t *testing.T) {
	hosts := []string{
		os.Getenv(constants.VcfTestHost1Fqdn),
		os.Getenv(constants.VcfTestHost2Fqdn),
		os.Getenv(constants.VcfTestHost3Fqdn),
		os.Getenv(constants.VcfTestHost4Fqdn),
	}

	var steps []resource.TestStep
	for _, fqdn := range hosts {
		steps = append(steps, resource.TestStep{
			Config: testAccDataSourceVcfHostConfig(fqdn),
			Check: resource.ComposeTestCheckFunc(
				resource.TestCheckResourceAttrSet("data.vcf_host.test_host", "id"),
				resource.TestCheckResourceAttr("data.vcf_host.test_host", "fqdn", fqdn),
			),
		})
	}

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: muxedFactories(),
		Steps:                    steps,
	})
}

func testAccDataSourceVcfHostConfig(hostFqdn string) string {
	return fmt.Sprintf(`
    data "vcf_host" "test_host" {
        fqdn = %q
    }
    `, hostFqdn)
}
