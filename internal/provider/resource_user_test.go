// © Broadcom. All Rights Reserved.
// The term “Broadcom” refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"fmt"
	"log"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"

	"github.com/vmware/terraform-provider-vcf/internal/api_client"
)

const (
	testUserName1 = "testuser1@vrack.vsphere.local"
	testUserName2 = "serviceuser1"
)

func TestAccResourceVcfUser(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: muxedFactories(),
		CheckDestroy:             testCheckVcfUserDestroy,
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
	apiClient := testAccProvider.Meta().(*api_client.SddcManagerClient).ApiClient

	ok, err := apiClient.GetUsersWithResponse(context.TODO())
	if err != nil {
		log.Println("error = ", err)
		return err
	}

	// Check if the users with the known usernames exist
	for _, user := range *ok.JSON200.Elements {
		if user.Name == testUserName1 || user.Name == testUserName2 {
			return fmt.Errorf("found user with username %q", user.Name)
		}
	}

	// Didn't find the test users
	return nil
}
