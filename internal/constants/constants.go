/* Copyright 2023 VMware, Inc.
   SPDX-License-Identifier: MPL-2.0 */

package constants

import "time"

const (
	DefaultVcfApiCallTimeout = 2 * time.Minute

	// VcfTestUrl URL of a VCF instance, used for Acceptance tests.
	VcfTestUrl = "VCF_TEST_URL"
	// VcfTestUsername username of SSO user, used for Acceptance tests.
	VcfTestUsername = "VCF_TEST_USERNAME"
	// VcfTestPassword an SSO user with the ADMIN role or admin@local API user, used for Acceptance tests.
	VcfTestPassword = "VCF_TEST_PASSWORD"

	// CloudBuilderTestUrl URL of a CloudBuilder instance, used for Acceptance tests.
	CloudBuilderTestUrl = "CLOUDBUILDER_TEST_URL"
	// CloudBuilderTestUsername username of CloudBuilder user, used for Acceptance tests.
	CloudBuilderTestUsername = "CLOUDBUILDER_TEST_USERNAME"
	// CloudBuilderTestPassword an CloudBuilder user, used for Acceptance tests.
	CloudBuilderTestPassword = "CLOUDBUILDER_TEST_PASSWORD"

	// VcfTestAllowUnverifiedTls allows VCF environments with self-signed certificates
	// to be used in Acceptance tests.
	VcfTestAllowUnverifiedTls = "VCF_TEST_ALLOW_UNVERIFIED_TLS"

	// VcfTestHost1Fqdn the FQDN of the first ESXi host, that has not been commissioned
	// with the SDDC Manager.
	VcfTestHost1Fqdn = "VCF_TEST_HOST1_FQDN"

	// VcfTestHost1Pass the password of the first ESXi host, that has not been commissioned
	// with the SDDC Manager.
	VcfTestHost1Pass = "VCF_TEST_HOST1_PASS"

	// VcfTestHost2Fqdn the FQDN of the second ESXi host, that has not been commissioned
	// with the SDDC Manager.
	VcfTestHost2Fqdn = "VCF_TEST_HOST2_FQDN"

	// VcfTestHost2Pass the password of the second ESXi host, that has not been commissioned
	// with the SDDC Manager.
	VcfTestHost2Pass = "VCF_TEST_HOST2_PASS"

	// VcfTestHost3Fqdn the FQDN of the third ESXi host, that has not been commissioned
	// with the SDDC Manager.
	VcfTestHost3Fqdn = "VCF_TEST_HOST3_FQDN"

	// VcfTestHost3Pass the password of the third ESXi host, that has not been commissioned
	// with the SDDC Manager.
	VcfTestHost3Pass = "VCF_TEST_HOST3_PASS"

	// VcfTestHost4Fqdn the FQDN of the forth ESXi host, that has not been commissioned
	// with the SDDC Manager.
	VcfTestHost4Fqdn = "VCF_TEST_HOST4_FQDN"

	// VcfTestHost4Pass the password of the forth ESXi host, that has not been commissioned
	// with the SDDC Manager.
	VcfTestHost4Pass = "VCF_TEST_HOST4_PASS"

	// VcfTestHost5Fqdn the FQDN of the fifth ESXi host, that has not been commissioned
	// with the SDDC Manager.
	VcfTestHost5Fqdn = "VCF_TEST_HOST5_FQDN"

	// VcfTestHost5Pass the password of the fifth ESXi host, that has not been commissioned
	// with the SDDC Manager.
	VcfTestHost5Pass = "VCF_TEST_HOST5_PASS"

	// VcfTestHost6Fqdn the FQDN of the sixth ESXi host, that has not been commissioned
	// with the SDDC Manager.
	VcfTestHost6Fqdn = "VCF_TEST_HOST6_FQDN"

	// VcfTestHost6Pass the password of the sixth ESXi host, that has not been commissioned
	// with the SDDC Manager.
	VcfTestHost6Pass = "VCF_TEST_HOST6_PASS"

	// VcfTestHost7Fqdn the FQDN of the seventh ESXi host, that has not been commissioned
	// with the SDDC Manager.
	VcfTestHost7Fqdn = "VCF_TEST_HOST7_FQDN"

	// VcfTestHost7Pass the password of the seventh ESXi host, that has not been commissioned
	// with the SDDC Manager.
	VcfTestHost7Pass = "VCF_TEST_HOST7_PASS"

	// VcfTestHost8Fqdn the FQDN of the eight ESXi host, that has not been commissioned
	// with the SDDC Manager.
	VcfTestHost8Fqdn = "VCF_TEST_HOST8_FQDN"

	// VcfTestHost8Pass the password of the eight ESXi host, that has not been commissioned
	// with the SDDC Manager.
	VcfTestHost8Pass = "VCF_TEST_HOST8_PASS"

	// VcfTestNsxLicenseKey license key for NSX required for domain and cluster acceptance tests.
	VcfTestNsxLicenseKey = "VCF_TEST_NSX_LICENSE_KEY"

	// VcfTestEsxiLicenseKey license key for vSphere required for workload domain and cluster acceptance tests.
	VcfTestEsxiLicenseKey = "VCF_TEST_ESXI_LICENSE_KEY"

	// VcfTestVsanLicenseKey license key for vSAN required for workload domain and cluster acceptance tests.
	VcfTestVsanLicenseKey = "VCF_TEST_VSAN_LICENSE_KEY"

	// VcfTestVcenterLicenseKey license key for vCenter required for bringup acceptance tests.
	VcfTestVcenterLicenseKey = "VCF_TEST_VCENTER_LICENSE_KEY"

	// VcfTestDomainDataSourceId id of a workload domain used in workload domain data source acceptance test.
	// Typically, the id of management domain is used as it is already created during bringup.
	VcfTestDomainDataSourceId = "VCF_DOMAIN_DATA_SOURCE_ID"

	// VcfTestClusterDataSourceId id of cluster used in cluster data source acceptance test.
	// Typically, the id of the default cluster in the management domain is used as it is
	// already created during bringup.
	VcfTestClusterDataSourceId = "VCF_CLUSTER_DATA_SOURCE_ID"

	// VcfTestNetworkPoolName used in vcf_network_pool Acceptance tests.
	VcfTestNetworkPoolName = "terraform-test-pool"
)
