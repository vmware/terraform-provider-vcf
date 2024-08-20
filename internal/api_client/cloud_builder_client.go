// Copyright 2023 Broadcom. All Rights Reserved.
// SPDX-License-Identifier: MPL-2.0

package api_client

import (
	"crypto/tls"
	"net/http"

	openapiclient "github.com/go-openapi/runtime/client"
	"github.com/go-openapi/strfmt"
	vcfclient "github.com/vmware/vcf-sdk-go/client"
)

// CloudBuilderClient is an API client that can execute the APIs of the CloudBuilder appliance.
// Supports only HTTP Basic Authentication.
type CloudBuilderClient struct {
	username           string
	password           string
	cloudBuilderUrl    string
	ApiClient          *vcfclient.VcfClient
	allowUnverifiedTls bool
}

func NewCloudBuilderClient(username, password, url string, allowUnverifiedTls bool) *CloudBuilderClient {
	result := &CloudBuilderClient{
		username:           username,
		password:           password,
		cloudBuilderUrl:    url,
		allowUnverifiedTls: allowUnverifiedTls,
	}
	result.init()
	return result
}

func (cloudBuilderClient *CloudBuilderClient) init() {
	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{
		InsecureSkipVerify: cloudBuilderClient.allowUnverifiedTls}

	cfg := vcfclient.DefaultTransportConfig()
	openApiClient := openapiclient.New(cloudBuilderClient.cloudBuilderUrl, cfg.BasePath, cfg.Schemes)

	openApiClient.Transport = cloudBuilderClient.newTransport()

	// create the API client, with the transport
	cloudBuilderOpenApiClient := vcfclient.New(openApiClient, strfmt.Default)
	// save the client for later use
	cloudBuilderClient.ApiClient = cloudBuilderOpenApiClient
}

func (cloudBuilderClient *CloudBuilderClient) newTransport() *cloudBuilderCustomHttpTransport {
	return &cloudBuilderCustomHttpTransport{
		originalTransport:  http.DefaultTransport,
		cloudBuilderClient: cloudBuilderClient,
	}
}

type cloudBuilderCustomHttpTransport struct {
	originalTransport  http.RoundTripper
	cloudBuilderClient *CloudBuilderClient
}

func (c *cloudBuilderCustomHttpTransport) RoundTrip(r *http.Request) (*http.Response, error) {

	r.SetBasicAuth(c.cloudBuilderClient.username, c.cloudBuilderClient.password)
	r.Header.Add("Content-Type", "application/json")

	resp, err := c.originalTransport.RoundTrip(r)
	if err != nil {
		return nil, err
	}

	return resp, nil
}
