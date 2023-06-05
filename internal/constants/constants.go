/* Copyright 2023 VMware, Inc.
   SPDX-License-Identifier: MPL-2.0 */

package constants

const (

	// VCF_TEST_URL URL of a VCF instance, used for Acceptance tests.
	VCF_TEST_URL = "VCF_TEST_URL"
	// VCF_TEST_USERNAME username of SSO user, used for Acceptance tests.
	VCF_TEST_USERNAME = "VCF_TEST_USERNAME"
	// VCF_TEST_PASSWORD passowrd of SSO user, used for Acceptance tests.
	VCF_TEST_PASSWORD = "VCF_TEST_PASSWORD"

	// VCF_TEST_COMMISSIONED_HOST_FQDN the FQDN of an ESXi host, that has not been commissioned
	// with the SDDC Manager.
	VCF_TEST_COMMISSIONED_HOST_FQDN = "VCF_TEST_COMMISSIONED_HOST_FQDN"

	// VCF_TEST_COMMISSIONED_HOST_PASS the SSH pass of an ESXi host, that has not been commissioned
	// with the SDDC Manager.
	VCF_TEST_COMMISSIONED_HOST_PASS = "VCF_TEST_COMMISSIONED_HOST_PASS"

	// VCF_TEST_NETWORK_POOL_NAME used in vcf_network_pool Acceptance tests.
	VCF_TEST_NETWORK_POOL_NAME = "terraform-test-pool"
)
