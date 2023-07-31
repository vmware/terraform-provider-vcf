/* Copyright 2023 VMware, Inc.
   SPDX-License-Identifier: MPL-2.0 */

package provider

import (
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"github.com/hashicorp/terraform-plugin-log/tflog"
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
	allowUnverifiedTls  bool
	lastRefreshTime     time.Time
	isRefreshing        bool
	getTaskRetries      int
}

// NewSddcManagerClient constructs new Client instance with vcf credentials.
func NewSddcManagerClient(username, password, host string, allowUnverifiedTls bool) *SddcManagerClient {
	return &SddcManagerClient{
		SddcManagerUsername: username,
		SddcManagerPassword: password,
		SddcManagerHost:     host,
		allowUnverifiedTls:  allowUnverifiedTls,
		lastRefreshTime:     time.Now(),
		isRefreshing:        false,
		getTaskRetries:      0,
	}
}

var accessToken *string

const maxGetTaskRetries int = 10
const maxTaskRetries int = 3

func newTransport(sddcManagerClient *SddcManagerClient) *customTransport {
	return &customTransport{
		originalTransport: http.DefaultTransport,
		sddcManagerClient: sddcManagerClient,
	}
}

type customTransport struct {
	originalTransport http.RoundTripper
	sddcManagerClient *SddcManagerClient
}

func (c *customTransport) RoundTrip(r *http.Request) (*http.Response, error) {
	// Refresh the access token every 20 minutes so that SDK operations won't start to
	// fail with 401, 403 because of token expiration, during long-running tasks
	if time.Since(c.sddcManagerClient.lastRefreshTime) > 20*time.Minute &&
		!c.sddcManagerClient.isRefreshing {
		err := c.sddcManagerClient.Connect()
		if err != nil {
			return nil, err
		}
	}

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
	sddcManagerClient.isRefreshing = true
	// Disable cert checks
	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{
		InsecureSkipVerify: sddcManagerClient.allowUnverifiedTls}

	cfg := vcfclient.DefaultTransportConfig()
	openApiClient := openapiclient.New(sddcManagerClient.SddcManagerHost, cfg.BasePath, cfg.Schemes)
	// openApiClient.SetDebug(true)

	openApiClient.Transport = newTransport(sddcManagerClient)

	// create the API client, with the transport
	vcfClient := vcfclient.New(openApiClient, strfmt.Default)
	// save the client for later use
	sddcManagerClient.ApiClient = vcfClient
	// Get access token
	tokenSpec := &models.TokenCreationSpec{
		Username: sddcManagerClient.SddcManagerUsername,
		Password: sddcManagerClient.SddcManagerPassword,
	}
	params := tokens.NewCreateTokenParams().
		WithTokenCreationSpec(tokenSpec).WithTimeout(constants.DefaultVcfApiCallTimeout)

	ok, _, err := vcfClient.Tokens.CreateToken(params)
	if err != nil {
		return err
	}

	accessToken = &ok.Payload.AccessToken
	// save the access token for later use
	sddcManagerClient.lastRefreshTime = time.Now()
	sddcManagerClient.AccessToken = &ok.Payload.AccessToken
	sddcManagerClient.isRefreshing = false
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
func (sddcManagerClient *SddcManagerClient) WaitForTaskComplete(ctx context.Context, taskId string, retry bool) error {
	log.Printf("Getting status of task %s", taskId)
	currentTaskRetries := 0
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
			errorMsg := fmt.Sprintf("Task with ID = %s , Name: %q Type: %q is in state %s", taskId, task.Name, task.Type, task.Status)
			tflog.Error(ctx, errorMsg)

			if retry && currentTaskRetries < maxTaskRetries {
				currentTaskRetries++
				err := sddcManagerClient.retryTask(ctx, taskId)
				if err != nil {
					tflog.Error(ctx, fmt.Sprintf("Task %q %q failed after %d retries",
						taskId, task.Type, currentTaskRetries))
					return err
				}
			} else {
				return errors.New(errorMsg)
			}
		}

		log.Printf("Task with ID = %s is in state %s, completed at %s", taskId, task.Status, task.CompletionTimestamp)
		return nil
	}
}

func (sddcManagerClient *SddcManagerClient) GetResourceIdAssociatedWithTask(ctx context.Context, taskId, resourceType string) (string, error) {
	task, err := sddcManagerClient.getTask(ctx, taskId)
	if err != nil {
		return "", err
	}
	if len(task.Resources) == 0 {
		return "", fmt.Errorf("no resources associated with Task with ID %q", taskId)
	}
	for _, resource := range task.Resources {
		if *resource.Type == resourceType {
			return *resource.ResourceID, nil
		}
	}
	return "", fmt.Errorf("task %q did not contain resources of type %q", taskId, resourceType)
}

func (sddcManagerClient *SddcManagerClient) getTask(ctx context.Context, taskId string) (*models.Task, error) {
	apiClient := sddcManagerClient.ApiClient
	getTaskParams := tasks.NewGetTaskParamsWithTimeout(constants.DefaultVcfApiCallTimeout).
		WithContext(ctx)
	getTaskParams.ID = taskId

	getTaskResult, err := apiClient.Tasks.GetTask(getTaskParams)
	if err != nil {
		// retry the task up to maxGetTaskRetries
		if sddcManagerClient.getTaskRetries < maxGetTaskRetries {
			sddcManagerClient.getTaskRetries++
			return sddcManagerClient.getTask(ctx, taskId)
		}
		log.Println("error = ", err)
		return nil, err
	}
	// reset the counter
	sddcManagerClient.getTaskRetries = 0
	return getTaskResult.Payload, nil
}

func (sddcManagerClient *SddcManagerClient) retryTask(ctx context.Context, taskId string) error {
	apiClient := sddcManagerClient.ApiClient
	retryTaskParams := tasks.NewRetryTaskParamsWithTimeout(constants.DefaultVcfApiCallTimeout).
		WithContext(ctx)
	retryTaskParams.ID = taskId
	_, err := apiClient.Tasks.RetryTask(retryTaskParams)
	if err != nil {
		return err
	}
	return nil
}
