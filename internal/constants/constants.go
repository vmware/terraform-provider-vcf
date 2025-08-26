// © Broadcom. All Rights Reserved.
// The term “Broadcom” refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

package constants

const (
	ProviderName = "terraform-provider-vcf"

	// VcfTestUrl URL of a VCF instance, used for acceptance tests.
	VcfTestUrl = "VCF_TEST_URL"
	// VcfTestUsername username of SSO user, used for acceptance tests.
	VcfTestUsername = "VCF_TEST_USERNAME"
	// VcfTestPassword an SSO user with the ADMIN role or admin@local API user, used for acceptance tests.
	VcfTestPassword = "VCF_TEST_PASSWORD"

	// InstallerTestUrl URL of a CloudBuilder instance, used for acceptance tests.
	InstallerTestUrl = "INSTALLER_TEST_URL"
	// InstallerTestUsername username of CloudBuilder user, used for acceptance tests.
	InstallerTestUsername = "INSTALLER_TEST_USERNAME"
	// InstallerTestPassword an CloudBuilder user, used for acceptance tests.
	InstallerTestPassword = "INSTALLER_TEST_PASSWORD"

	// VcfTestAllowUnverifiedTls allows VCF environments with self-signed certificates
	// to be used in acceptance tests.
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

	// VcfTestWitnessHostIp the IP address of the witness host to be used to stretch a cluster.
	VcfTestWitnessHostIp = "VCF_TEST_WITNESS_HOST_IP"

	// VcfTestWitnessHostCidr the CIDR notatin address of the witness host to be used to stretch a cluster.
	VcfTestWitnessHostCidr = "VCF_TEST_WITNESS_HOST_CIDR"

	// VcfTestWitnessHostFqdn the FQDN of the witness host to be used to stretch a cluster.
	VcfTestWitnessHostFqdn = "VCF_TEST_WITNESS_HOST_FQDN"

	// VcfTestClusterId the identifier of the cluster within its vCenter server.
	VcfTestClusterId = "VCF_TEST_CLUSTER_ID"

	// VcfTestClusterImageId the identifier of the lifecycle image for the cluster.
	VcfTestClusterImageId = "VCF_TEST_CLUSTER_IMAGE_ID"

	// VcfTestDomainDataSourceId id of a workload domain used in workload domain data source acceptance test.
	// Typically, the id of management domain is used as it is already created during bringup.
	VcfTestDomainDataSourceId = "VCF_DOMAIN_DATA_SOURCE_ID"

	// VcfTestDomainDataSourceName name of a workload domain used in workload domain data source acceptance test.
	// Typically, the name of management domain is used as it is already created during bringup.
	VcfTestDomainDataSourceName = "VCF_DOMAIN_DATA_SOURCE_NAME"

	// VcfTestClusterDataSourceId id of cluster used in cluster data source acceptance test.
	// Typically, the id of the default cluster in the management domain is used as it is
	// already created during bringup.
	VcfTestClusterDataSourceId = "VCF_CLUSTER_DATA_SOURCE_ID"

	// VcfTestDomainName display name of the workload domain used in the acceptance tests.
	VcfTestDomainName = "VCF_DOMAIN_NAME"

	// VcfTestNetworkPoolName used in vcf_network_pool acceptance tests.
	VcfTestNetworkPoolName = "networkpool-1"

	// VcfTestMsftCaServerUrl used in vcf_certificate_authority tests.
	VcfTestMsftCaServerUrl = "VCF_TEST_MSFT_CA_SERVER_URL"

	// VcfTestMsftCaUser used in vcf_certificate_authority tests.
	VcfTestMsftCaUser = "VCF_TEST_MSFT_CA_USER"

	// VcfTestMsftCaSecret used in vcf_certificate_authority tests.
	VcfTestMsftCaSecret = "VCF_TEST_MSFT_CA_SECRET"

	// VcfTestResourceCertificate used in vcf_external_certificate tests.
	VcfTestResourceCertificate = "VCF_TEST_RESOURCE_CERTIFICATE"

	// VcfTestResourceCaCertificate used in vcf_external_certificate tests.
	VcfTestResourceCaCertificate = "VCF_TEST_RESOURCE_CA_CERTIFICATE"

	// VcfTestEdgeClusterRootPass the root user password for the NSX manager.
	VcfTestEdgeClusterRootPass = "VCF_TEST_EDGE_CLUSTER_ROOT_PASS"

	// VcfTestEdgeClusterAdminPass the admin user password for the NSX manager.
	VcfTestEdgeClusterAdminPass = "VCF_TEST_EDGE_CLUSTER_ADMIN_PASS"

	// VcfTestEdgeClusterAuditPass the audit user password for the NSX manager.
	VcfTestEdgeClusterAuditPass = "VCF_TEST_EDGE_CLUSTER_AUDIT_PASS"

	// VcfTestEdgeNodeRootPass the root user password for the edge nodes.
	VcfTestEdgeNodeRootPass = "VCF_TEST_EDGE_NODE_ROOT_PASS"

	// VcfTestEdgeNodeAdminPass the admin user password for the edge nodes.
	VcfTestEdgeNodeAdminPass = "VCF_TEST_EDGE_NODE_ADMIN_PASS"

	// VcfTestEdgeNodeAuditPass the audit user password for the edge nodes.
	VcfTestEdgeNodeAuditPass = "VCF_TEST_EDGE_NODE_AUDIT_PASS"

	// VcfTestComputeClusterId the identifier of the compute cluster that will contain the edge nodes.
	VcfTestComputeClusterId = "VCF_TEST_COMPUTE_CLUSTER_ID"

	// VcfTestVcenterFqdn the FQDN of the vcenter server.
	VcfTestVcenterFqdn = "VCF_TEST_VCENTER_FQDN"

	// VcfTestSddcManagerFqdn the FQDN of the SDDC manager.
	VcfTestSddcManagerFqdn = "VCF_TEST_SDDC_MANAGER_FQDN"

	// VcfTestNsxManagerFqdn the FQDN of the NSX manager.
	VcfTestNsxManagerFqdn = "VCF_TEST_NSX_MANAGER_FQDN"
)

func GetIso3166CountryCodes() []string {
	return []string{"US", "CA", "AX", "AD", "AE", "AF", "AG", "AI", "AL", "AM", "AN", "AO", "AQ", "AR", "AS", "AT", "AU",
		"AW", "AZ", "BA", "BB", "BD", "BE", "BF", "BG", "BH", "BI", "BJ", "BM", "BN", "BO", "BR", "BS", "BT", "BV", "BW", "BZ", "CA", "CC", "CF", "CH", "CI", "CK",
		"CL", "CM", "CN", "CO", "CR", "CS", "CV", "CX", "CY", "CZ", "DE", "DJ", "DK", "DM", "DO", "DZ", "EC", "EE", "EG", "EH", "ER", "ES", "ET", "FI", "FJ", "FK",
		"FM", "FO", "FR", "FX", "GA", "GB", "GD", "GE", "GF", "GG", "GH", "GI", "GL", "GM", "GN", "GP", "GQ", "GR", "GS", "GT", "GU", "GW", "GY", "HK", "HM", "HN",
		"HR", "HT", "HU", "ID", "IE", "IL", "IM", "IN", "IO", "IS", "IT", "JE", "JM", "JO", "JP", "KE", "KG", "KH", "KI", "KM", "KN", "KR", "KW", "KY", "KZ", "LA",
		"LC", "LI", "LK", "LS", "LT", "LU", "LV", "LY", "MA", "MC", "MD", "ME", "MG", "MH", "MK", "ML", "MM", "MN", "MO", "MP", "MQ", "MR", "MS", "MT", "MU", "MV",
		"MW", "MX", "MY", "MZ", "NA", "NC", "NE", "NF", "NG", "NI", "NL", "NO", "NP", "NR", "NT", "NU", "NZ", "OM", "PA", "PE", "PF", "PG", "PH", "PK", "PL", "PM",
		"PN", "PR", "PS", "PT", "PW", "PY", "QA", "RE", "RO", "RS", "RU", "RW", "SA", "SB", "SC", "SE", "SG", "SH", "SI", "SJ", "SK", "SL", "SM", "SN", "SR", "ST",
		"SU", "SV", "SZ", "TC", "TD", "TF", "TG", "TH", "TJ", "TK", "TM", "TN", "TO", "TP", "TR", "TT", "TV", "TW", "TZ", "UA", "UG", "UM", "US", "UY", "UZ", "VA",
		"VC", "VE", "VG", "VI", "VN", "VU", "WF", "WS", "YE", "YT", "ZA", "ZM", "COM", "EDU", "GOV", "INT", "MIL", "NET", "ORG", "ARPA"}
}
