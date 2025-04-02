// © Broadcom. All Rights Reserved.
// The term “Broadcom” refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

package api_client

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/vmware/vcf-sdk-go/vcf"
)

// SddcManagerClient model that represents properties to authenticate against VCF instance.
type SddcManagerClient struct {
	username           string
	password           string
	sddcManagerUrl     string
	accessToken        *string
	ApiClient          *vcf.ClientWithResponses
	allowUnverifiedTls bool
	lastRefreshTime    time.Time
	isRefreshing       bool
}

// NewSddcManagerClient constructs new Client instance with vcf credentials.
func NewSddcManagerClient(username, password, url string, allowUnverifiedTls bool) *SddcManagerClient {
	return &SddcManagerClient{
		username:           username,
		password:           password,
		sddcManagerUrl:     url,
		allowUnverifiedTls: allowUnverifiedTls,
		lastRefreshTime:    time.Now(),
		isRefreshing:       false,
	}
}

func (sddcManagerClient *SddcManagerClient) authEditor(ctx context.Context, req *http.Request) error {
	// Refresh the access token every 20 minutes so that SDK operations won't start to
	// fail with 401, 403 because of token expiration, during long-running tasks
	if time.Since(sddcManagerClient.lastRefreshTime) > 20*time.Minute &&
		!sddcManagerClient.isRefreshing {
		err := sddcManagerClient.Connect()
		if err != nil {
			return err
		}
	}

	if sddcManagerClient.accessToken != nil {
		req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", *sddcManagerClient.accessToken))
	}

	req.Header.Add("Content-Type", "application/json")

	return nil
}

func (sddcManagerClient *SddcManagerClient) Connect() error {
	sddcManagerClient.isRefreshing = true

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: sddcManagerClient.allowUnverifiedTls},
	}
	httpClient := &http.Client{Transport: tr}
	client, err := vcf.NewClientWithResponses(fmt.Sprintf("https://%s", sddcManagerClient.sddcManagerUrl),
		vcf.WithRequestEditorFn(sddcManagerClient.authEditor), vcf.WithHTTPClient(httpClient))
	if err != nil {
		return err
	}

	sddcManagerClient.ApiClient = client

	tokenCreationSpec := vcf.TokenCreationSpec{
		Username: &sddcManagerClient.username,
		Password: &sddcManagerClient.password,
	}

	res, err := client.CreateTokenWithResponse(context.Background(), tokenCreationSpec)
	if err != nil {
		return err
	}

	tokenPair, vcfErr := GetResponseAs[vcf.TokenPair](res)
	if vcfErr != nil {
		LogError(vcfErr)
		return errors.New(*vcfErr.Message)
	}
	sddcManagerClient.accessToken = tokenPair.AccessToken
	sddcManagerClient.lastRefreshTime = time.Now()
	sddcManagerClient.isRefreshing = false

	return nil
}

func (sddcManagerClient *SddcManagerClient) GetResourceIdAssociatedWithTask(ctx context.Context, taskId, resourceType string) (string, error) {
	task, err := sddcManagerClient.getTask(ctx, taskId)
	if err != nil {
		return "", err
	}
	if len(*task.Resources) == 0 {
		return "", fmt.Errorf("no resources associated with Task with ID %q", taskId)
	}
	for _, resource := range *task.Resources {
		if resource.Type == resourceType {
			return resource.ResourceId, nil
		}
	}
	return "", fmt.Errorf("task %q did not contain resources of type %q", taskId, resourceType)
}

func (sddcManagerClient *SddcManagerClient) getTask(ctx context.Context, taskId string) (*vcf.Task, error) {
	apiClient := sddcManagerClient.ApiClient
	res, err := apiClient.GetTaskWithResponse(ctx, taskId)
	task, vcfErr := GetResponseAs[vcf.Task](res)
	if err != nil || vcfErr != nil {
		log.Println("error = ", err)
		return nil, err
	}

	return task, nil
}

type Response interface {
	GetBody() []byte
	StatusCode() int
}

// GetResponseAs attempts to parse the response body into the provided type
// If it fails it attempts to parse it as a vcf.Error.
func GetResponseAs[T interface{}](res Response) (*T, *vcf.Error) {
	if res.StatusCode() < 200 || res.StatusCode() >= 300 {
		return nil, GetError(res.GetBody())
	}

	var resp T
	if json.Unmarshal(res.GetBody(), &resp) != nil {
		return nil, nil
	}

	return &resp, nil
}

// GetError when the API responds with an error code the response is unmarshalled into the appropriate field for that code
// all error code fields are of type *vcf.Error and only one can be != nil at any time
// if the status code is an error code the body is always *vcf.Error
//
// use this method if you are not interested in the error code but only in the error itself.
func GetError(body []byte) *vcf.Error {
	var dest vcf.Error
	if json.Unmarshal(body, &dest) != nil {
		return nil
	}

	return &dest
}

// LogError traverses a vcf.Error structure and logs its error message as well as
// the messages of any nested errors.
func LogError(err *vcf.Error) {
	if err != nil {
		if err.Message != nil {
			tflog.Error(context.Background(), *err.Message)
		}
		if err.NestedErrors != nil {
			for _, nestedErr := range *err.NestedErrors {
				LogError(&nestedErr)
			}
		}
	}
}
