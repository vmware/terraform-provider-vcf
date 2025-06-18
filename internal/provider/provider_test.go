// © Broadcom. All Rights Reserved.
// The term “Broadcom” refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"log"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-mux/tf5to6server"
	"github.com/hashicorp/terraform-plugin-mux/tf6muxserver"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/vmware/terraform-provider-vcf/internal/constants"
	validationUtils "github.com/vmware/terraform-provider-vcf/internal/validation"
)

var testAccProvider *schema.Provider
var testAccFrameworkProvider provider.Provider

func TestProvider(t *testing.T) {
	if err := Provider().InternalValidate(); err != nil {
		t.Fatalf("err: %s", err)
	}
}

func muxedFactories() map[string]func() (tfprotov6.ProviderServer, error) {
	testAccProvider = Provider()
	testAccFrameworkProvider = New()

	ctx := context.Background()
	upgradedSdkServer, err := tf5to6server.UpgradeServer(
		ctx,
		testAccProvider.GRPCProvider,
	)

	if err != nil {
		log.Fatal(err)
	}

	providers := []func() tfprotov6.ProviderServer{
		providerserver.NewProtocol6(testAccFrameworkProvider),
		func() tfprotov6.ProviderServer {
			return upgradedSdkServer
		},
	}

	muxServer, err := tf6muxserver.NewMuxServer(ctx, providers...)

	if err != nil {
		log.Fatal(err)
	}

	return map[string]func() (tfprotov6.ProviderServer, error){
		"vcf": func() (tfprotov6.ProviderServer, error) { return muxServer, err },
	}
}

// testAccPreCheck validates all required environment variables for running acceptance
// tests are set.
func testAccPreCheck(t *testing.T) {
	testAccSDDCManagerOrCloudBuilderPreCheck(t)
	testAccHostsPreCheck(t, 8)
	if v := os.Getenv(constants.VcfTestDomainDataSourceId); v == "" {
		t.Fatal(constants.VcfTestDomainDataSourceId + " must be set for acceptance tests")
	}
	if v := os.Getenv(constants.VcfTestClusterDataSourceId); v == "" {
		t.Fatal(constants.VcfTestClusterDataSourceId + " must be set for acceptance tests")
	}
}

func testAccSDDCManagerOrCloudBuilderPreCheck(t *testing.T) {
	testSddcManagerUrl := os.Getenv(constants.VcfTestUrl)
	testCloudBuilderUrl := os.Getenv(constants.InstallerTestUrl)
	if validationUtils.IsEmpty(testSddcManagerUrl) && validationUtils.IsEmpty(testCloudBuilderUrl) {
		t.Fatal(constants.VcfTestUrl + " or " + constants.InstallerTestUrl +
			" must be set for acceptance tests")
	}
	testSddcManagerUsername := os.Getenv(constants.VcfTestUsername)
	testCloudBuilderUsername := os.Getenv(constants.InstallerTestUsername)
	if validationUtils.IsEmpty(testSddcManagerUsername) && validationUtils.IsEmpty(testCloudBuilderUsername) {
		t.Fatal(constants.VcfTestUsername + " or " + constants.InstallerTestUsername +
			" must be set for acceptance tests")
	}
	testSddcManagerPassword := os.Getenv(constants.VcfTestPassword)
	testCloudBuilderPassword := os.Getenv(constants.InstallerTestPassword)
	if validationUtils.IsEmpty(testSddcManagerPassword) && validationUtils.IsEmpty(testCloudBuilderPassword) {
		t.Fatal(constants.VcfTestPassword + " or " + constants.InstallerTestPassword +
			" must be set for acceptance tests")
	}
}

func testAccHostsPreCheck(t *testing.T, numberOfHosts int) {
	hostList := []string{
		constants.VcfTestHost1Fqdn,
		constants.VcfTestHost2Fqdn,
		constants.VcfTestHost3Fqdn,
		constants.VcfTestHost4Fqdn,
		constants.VcfTestHost5Fqdn,
		constants.VcfTestHost6Fqdn,
		constants.VcfTestHost7Fqdn,
		constants.VcfTestHost8Fqdn,
	}

	passwordList := []string{
		constants.VcfTestHost1Pass,
		constants.VcfTestHost2Pass,
		constants.VcfTestHost3Pass,
		constants.VcfTestHost4Pass,
		constants.VcfTestHost5Pass,
		constants.VcfTestHost6Pass,
		constants.VcfTestHost7Pass,
		constants.VcfTestHost8Pass,
	}

	if numberOfHosts < len(hostList) {
		t.Fatal("Too many hosts required")
		return
	}

	for i := numberOfHosts - 1; i >= 0; i-- {
		hostNameEnvVar := hostList[i]
		passwordEnvVar := passwordList[i]
		if v := os.Getenv(hostNameEnvVar); v == "" {
			t.Fatal(hostNameEnvVar + " must be set for acceptance tests")
		}

		if v := os.Getenv(passwordEnvVar); v == "" {
			t.Fatal(passwordEnvVar + " must be set for acceptance tests")
		}
	}
}
