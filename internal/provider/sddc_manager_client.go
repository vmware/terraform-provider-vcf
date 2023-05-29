/* Copyright 2023 VMware, Inc.
   SPDX-License-Identifier: MPL-2.0 */

package provider

import (
	"crypto/tls"
	"errors"
	"fmt"
	"github.com/vmware/vcf-sdk-go/client/tasks"
	"github.com/vmware/vcf-sdk-go/client/tokens"
	"github.com/vmware/vcf-sdk-go/models"
	"log"
	"net/http"
	"time"

	openapiclient "github.com/go-openapi/runtime/client"
	"github.com/go-openapi/strfmt"
	vcfclient "github.com/vmware/vcf-sdk-go/client"
)

// SddcManagerClient model that represents properties to authenticate against VCF instance.
type SddcManagerClient struct {
	SddcManagerUsername string
	SddcManagerPassword string
	SddcManagerHost     string
	AccessToken         *string
	ApiClient           *vcfclient.VcfClient
}

// NewSddcManagerClient constructs new Client instance with vcf credentials.
func NewSddcManagerClient(username, password, host string) *SddcManagerClient {
	return &SddcManagerClient{
		SddcManagerUsername: username,
		SddcManagerPassword: password,
		SddcManagerHost:     host,
	}
}

var accessToken *string

func newTransport() *customTransport {
	return &customTransport{
		originalTransport: http.DefaultTransport,
	}
}

type customTransport struct {
	originalTransport http.RoundTripper
}

func (c *customTransport) RoundTrip(r *http.Request) (*http.Response, error) {
	if accessToken != nil {
		r.Header.Add("Authorization", fmt.Sprintf("Bearer %s", *accessToken))
	}

	r.Header.Add("Content-Type", "application/json")

	resp, err := c.originalTransport.RoundTrip(r)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (sddcManagerClient *SddcManagerClient) Connect() {
	// Disable cert checks
	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}

	cfg := vcfclient.DefaultTransportConfig()
	oclient := openapiclient.New(sddcManagerClient.SddcManagerHost, cfg.BasePath, cfg.Schemes)
	// oclient.SetDebug(true)

	oclient.Transport = newTransport()

	// create the API client, with the transport
	vclient := vcfclient.New(oclient, strfmt.Default)
	// save the client for later use
	sddcManagerClient.ApiClient = vclient
	// Get access token
	tokenSpec := &models.TokenCreationSpec{
		Username: sddcManagerClient.SddcManagerUsername,
		Password: sddcManagerClient.SddcManagerPassword,
	}
	params := tokens.NewCreateTokenParams().WithTokenCreationSpec(tokenSpec)

	ok, _, err := vclient.Tokens.CreateToken(params)
	if err != nil {
		log.Fatal(err)
	}

	accessToken = &ok.Payload.AccessToken
	// save the access token for later use
	sddcManagerClient.AccessToken = &ok.Payload.AccessToken
}

// WaitForTask Wait for a task to complete (waits for up to a minute).
func (sddcManagerClient *SddcManagerClient) WaitForTask(taskId string) error {
	apiClient := sddcManagerClient.ApiClient
	// Fetch task status 10 times with a delay of 20 seconds each time
	taskStatusRetry := 10

	for taskStatusRetry > 0 {
		log.Printf("Getting status of task %s, retry left %d", taskId, taskStatusRetry)
		getTaskParams := tasks.NewGETTaskParams()
		getTaskParams.ID = taskId

		getTaskOk, err := apiClient.Tasks.GETTask(getTaskParams)
		if err != nil {
			log.Println("error = ", err)
			return err
		}

		if getTaskOk.Payload.Status == "In Progress" || getTaskOk.Payload.Status == "Pending" {
			time.Sleep(20 * time.Second)
			taskStatusRetry--
			continue
		}

		if getTaskOk.Payload.Status == "Failed" || getTaskOk.Payload.Status == "Cancelled" {
			errorMsg := fmt.Sprintf("Task with ID = %s is in state %s", getTaskParams.ID, getTaskOk.Payload.Status)
			log.Println(errorMsg)
			return errors.New(errorMsg)
		}

		log.Printf("Task with ID = %s is in state %s, completed at %s", getTaskParams.ID, getTaskOk.Payload.Status, getTaskOk.Payload.CompletionTimestamp)
		return nil
	}

	return fmt.Errorf("timedout waiting for task %s", taskId)
}

// WaitForTaskComplete Wait for task till it completes (either succeeds or fails).
func (sddcManagerClient *SddcManagerClient) WaitForTaskComplete(taskId string) error {
	apiClient := sddcManagerClient.ApiClient
	log.Printf("Getting status of task %s", taskId)
	for {
		getTaskParams := tasks.NewGETTaskParams()
		getTaskParams.ID = taskId

		getTaskOk, err := apiClient.Tasks.GETTask(getTaskParams)
		if err != nil {
			log.Println("error = ", err)
			return err
		}

		if getTaskOk.Payload.Status == "In Progress" || getTaskOk.Payload.Status == "Pending" {
			time.Sleep(20 * time.Second)
			continue
		}

		if getTaskOk.Payload.Status == "Failed" || getTaskOk.Payload.Status == "Cancelled" {
			errorMsg := fmt.Sprintf("Task with ID = %s is in state %s", getTaskParams.ID, getTaskOk.Payload.Status)
			log.Println(errorMsg)
			return errors.New(errorMsg)
		}

		log.Printf("Task with ID = %s is in state %s, completed at %s", getTaskParams.ID, getTaskOk.Payload.Status, getTaskOk.Payload.CompletionTimestamp)
		return nil
	}
}
