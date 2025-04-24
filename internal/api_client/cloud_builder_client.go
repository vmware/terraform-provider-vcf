// © Broadcom. All Rights Reserved.
// The term “Broadcom” refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

package api_client

import (
	"context"
	"crypto/tls"
	"fmt"
	"net/http"

	"github.com/vmware/vcf-sdk-go/vcf"
)

// CloudBuilderClient is an API client that can execute the APIs of the CloudBuilder appliance.
// Supports only HTTP Basic Authentication.
type CloudBuilderClient struct {
	username           string
	password           string
	cloudBuilderUrl    string
	providerVersion    string
	ApiClient          *vcf.ClientWithResponses
	allowUnverifiedTls bool
}

func NewCloudBuilderClient(username, password, url, providerVersion string, allowUnverifiedTls bool) *CloudBuilderClient {
	result := &CloudBuilderClient{
		username:           username,
		password:           password,
		cloudBuilderUrl:    url,
		providerVersion:    providerVersion,
		allowUnverifiedTls: allowUnverifiedTls,
	}
	result.init()
	return result
}

func (cloudBuilderClient *CloudBuilderClient) init() {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: cloudBuilderClient.allowUnverifiedTls},
	}
	httpClient := &http.Client{Transport: tr}
	client, err := vcf.NewClientWithResponses(fmt.Sprintf("https://%s", cloudBuilderClient.cloudBuilderUrl),
		vcf.WithRequestEditorFn(cloudBuilderClient.authEditor), vcf.WithHTTPClient(httpClient))
	if err == nil {
		cloudBuilderClient.ApiClient = client
	}
}

func (cloudBuilderClient *CloudBuilderClient) authEditor(ctx context.Context, req *http.Request) error {
	req.SetBasicAuth(cloudBuilderClient.username, cloudBuilderClient.password)
	req.Header.Add("Content-Type", "application/json")
	req.Header.Set("User-Agent", fmt.Sprintf("terraform-provider-vcf/%s", cloudBuilderClient.providerVersion))

	return nil
}
