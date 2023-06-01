/* Copyright 2023 VMware, Inc.
   SPDX-License-Identifier: MPL-2.0 */

package provider

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"log"
	"testing"
)

func TestAccResourceVcfCeip(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		CheckDestroy:      testCheckVcfCeipDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccVcfCeip(),
				Check:  testCheckVcfCeipCreate,
			},
		},
	})
}

func testAccVcfCeip() string {
	return `
	resource "vcf_ceip" "ceip" {
  		status    = "ENABLED"
	}`
}

func testVerifyVcfCeip(enabledState string) error {
	vcfClient := testAccProvider.Meta().(*SddcManagerClient)
	apiClient := vcfClient.ApiClient

	ceipResult, err := apiClient.CEIP.GetCEIPStatus(nil)
	if err != nil {
		log.Println("error = ", err)
		return err
	}
	if *ceipResult.Payload.Status == enabledState {
		return nil
	} else {
		return fmt.Errorf("CEIP not in status %q", enabledState)
	}

}

func testCheckVcfCeipCreate(_ *terraform.State) error {
	return testVerifyVcfCeip(ENABLED_STATE)
}

func testCheckVcfCeipDestroy(_ *terraform.State) error {
	return testVerifyVcfCeip(DISABLED_STATE)
}
