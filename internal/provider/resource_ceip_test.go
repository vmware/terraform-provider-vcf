// Copyright 2023 Broadcom. All Rights Reserved.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"fmt"
	"log"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/vmware/terraform-provider-vcf/internal/api_client"
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
	vcfClient := testAccProvider.Meta().(*api_client.SddcManagerClient)
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
	return testVerifyVcfCeip(EnabledState)
}

func testCheckVcfCeipDestroy(_ *terraform.State) error {
	return testVerifyVcfCeip(DisabledState)
}
