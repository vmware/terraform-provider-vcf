/*
 *  Copyright 2023 VMware, Inc.
 *    SPDX-License-Identifier: MPL-2.0
 */

package provider

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/vmware/terraform-provider-vcf/internal/api_client"
	"log"
	"testing"
)

const (
	testUserName1 = "testuser1@vrack.vsphere.local"
	testUserName2 = "serviceuser1"
)

func TestAccResourceVcfUser(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		CheckDestroy:      testCheckVcfUserDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccVcfUserConfig(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("vcf_user.testuser1", "id"),
					resource.TestCheckResourceAttrSet("vcf_user.testuser1", "creation_timestamp"),
					resource.TestCheckResourceAttrSet("vcf_user.serviceuser1", "id"),
					resource.TestCheckResourceAttrSet("vcf_user.serviceuser1", "api_key"),
					resource.TestCheckResourceAttrSet("vcf_user.serviceuser1", "creation_timestamp"),
				),
			},
		},
	})
}

func testAccVcfUserConfig() string {
	return fmt.Sprintf(`
	resource "vcf_user" "testuser1" {
		name      = %q
		domain    = "vrack.vsphere.local"
		type      = "USER"
		role_name = "VIEWER"
	}

	resource "vcf_user" "serviceuser1" {
		name      = %q
		domain    = "vrack.vsphere.local"
		type      = "SERVICE"
		role_name = "VIEWER"
	}
`, testUserName1, testUserName2)
}

func testCheckVcfUserDestroy(_ *terraform.State) error {
	vcfClient := testAccProvider.Meta().(*api_client.SddcManagerClient)
	apiClient := vcfClient.ApiClient

	ok, err := apiClient.Users.GetUsers(nil)
	if err != nil {
		log.Println("error = ", err)
		return err
	}

	// Check if the users with the known usernames exist
	for _, user := range ok.Payload.Elements {
		if *user.Name == testUserName1 || *user.Name == testUserName2 {
			return fmt.Errorf("found user with username %q", *user.Name)
		}
	}

	// Didn't find the test users
	return nil
}
