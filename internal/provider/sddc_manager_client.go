/* Copyright 2023 VMware, Inc.
   SPDX-License-Identifier: MPL-2.0 */

package provider

import (
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"github.com/go-openapi/runtime"
	"github.com/vmware/terraform-provider-vcf/internal/constants"
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

func (sddcManagerClient *SddcManagerClient) Connect() error {
	// Disable cert checks
	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}

	cfg := vcfclient.DefaultTransportConfig()
	openApiClient := openapiclient.New(sddcManagerClient.SddcManagerHost, cfg.BasePath, cfg.Schemes)
	// openApiClient.SetDebug(true)

	openApiClient.Transport = newTransport()

	// create the API client, with the transport
	vclient := vcfclient.New(openApiClient, strfmt.Default)
	// save the client for later use
	sddcManagerClient.ApiClient = vclient
	// Get access token
	tokenSpec := &models.TokenCreationSpec{
		Username: sddcManagerClient.SddcManagerUsername,
		Password: sddcManagerClient.SddcManagerPassword,
	}
	params := tokens.NewCreateTokenParams().
		WithTokenCreationSpec(tokenSpec).WithTimeout(constants.DefaultVcfApiCallTimeout)

	ok, _, err := vclient.Tokens.CreateToken(params)
	if err != nil {
		return err
	}

	accessToken = &ok.Payload.AccessToken
	// save the access token for later use
	sddcManagerClient.AccessToken = &ok.Payload.AccessToken
	return nil
}

// WaitForTask Wait for a task to complete (waits for up to a minute).
func (sddcManagerClient *SddcManagerClient) WaitForTask(ctx context.Context, taskId string) error {
	// Fetch task status 10 times with a delay of 20 seconds each time
	taskStatusRetry := 10

	for taskStatusRetry > 0 {
		task, err := sddcManagerClient.getTask(ctx, taskId)
		if err != nil {
			log.Println("error = ", err)
			return err
		}

		if task.Status == "In Progress" || task.Status == "Pending" {
			time.Sleep(20 * time.Second)
			taskStatusRetry--
			continue
		}

		if task.Status == "Failed" || task.Status == "Cancelled" {
			errorMsg := fmt.Sprintf("Task with ID = %s is in state %s", taskId, task.Status)
			log.Println(errorMsg)
			return errors.New(errorMsg)
		}

		log.Printf("Task with ID = %s is in state %s, completed at %s", taskId, task.Status, task.CompletionTimestamp)
		return nil
	}

	return fmt.Errorf("timedout waiting for task %s", taskId)
}

// WaitForTaskComplete Wait for task till it completes (either succeeds or fails).
func (sddcManagerClient *SddcManagerClient) WaitForTaskComplete(ctx context.Context, taskId string) error {
	log.Printf("Getting status of task %s", taskId)
	for {
		task, err := sddcManagerClient.getTask(ctx, taskId)
		if err != nil {
			return err
		}

		if task.Status == "In Progress" || task.Status == "Pending" {
			time.Sleep(20 * time.Second)
			continue
		}

		if task.Status == "Failed" || task.Status == "Cancelled" {
			errorMsg := fmt.Sprintf("Task with ID = %s is in state %s", taskId, task.Status)
			log.Println(errorMsg)
			return errors.New(errorMsg)
		}

		log.Printf("Task with ID = %s is in state %s, completed at %s", taskId, task.Status, task.CompletionTimestamp)
		return nil
	}
}

func (sddcManagerClient *SddcManagerClient) GetResourceIdAssociatedWithTask(ctx context.Context, taskId string) (string, error) {
	task, err := sddcManagerClient.getTask(ctx, taskId)
	if err != nil {
		return "", err
	}
	if len(task.Resources) == 0 {
		return "", fmt.Errorf("no resources associated with Task with ID %q", taskId)
	}
	return *task.Resources[0].ResourceID, nil
}

func (sddcManagerClient *SddcManagerClient) getTask(ctx context.Context, taskId string) (*models.Task, error) {
	apiClient := sddcManagerClient.ApiClient
	getTaskParams := tasks.NewGetTaskParamsWithTimeout(constants.DefaultVcfApiCallTimeout).
		WithContext(ctx)
	getTaskParams.ID = taskId

	getTaskResult, err := apiClient.Tasks.GetTask(getTaskParams)
	if err != nil {
		apiError, isApiError := err.(*runtime.APIError)
		// if the error is 4xx, then relogin and retry getting the task
		if isApiError && apiError.IsClientError() {
			err := sddcManagerClient.Connect()
			if err != nil {
				return nil, err
			}
			return sddcManagerClient.getTask(ctx, taskId)
		}

		log.Println("error = ", err)
		return nil, err
	}
	return getTaskResult.Payload, nil
}
