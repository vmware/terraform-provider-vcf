// © Broadcom. All Rights Reserved.
// The term “Broadcom” refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

package api_client

import (
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/vmware/vcf-sdk-go/installer"
)

// InstallerClient model that represents properties to authenticate against VCF installer.
type InstallerClient struct {
	username           string
	password           string
	vcfInstallerUrl    string
	accessToken        *string
	ApiClient          *installer.ClientWithResponses
	allowUnverifiedTls bool
	lastRefreshTime    time.Time
	isRefreshing       bool
}

// NewInstallerClient constructs new Client instance with vcf credentials.
func NewInstallerClient(username, password, url string, allowUnverifiedTls bool) *InstallerClient {
	return &InstallerClient{
		username:           username,
		password:           password,
		vcfInstallerUrl:    url,
		allowUnverifiedTls: allowUnverifiedTls,
		lastRefreshTime:    time.Now(),
		isRefreshing:       false,
	}
}

func (installerClient *InstallerClient) authEditor(ctx context.Context, req *http.Request) error {
	// Refresh the access token every 20 minutes so that SDK operations won't start to
	// fail with 401, 403 because of token expiration, during long-running tasks
	if time.Since(installerClient.lastRefreshTime) > 20*time.Minute &&
		!installerClient.isRefreshing {
		err := installerClient.Connect()
		if err != nil {
			return err
		}
	}

	if installerClient.accessToken != nil {
		req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", *installerClient.accessToken))
	}

	req.Header.Add("Content-Type", "application/json")

	return nil
}

func (installerClient *InstallerClient) Connect() error {
	installerClient.isRefreshing = true

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: installerClient.allowUnverifiedTls},
	}
	httpClient := &http.Client{Transport: tr}
	client, err := installer.NewClientWithResponses(fmt.Sprintf("https://%s", installerClient.vcfInstallerUrl),
		installer.WithRequestEditorFn(installerClient.authEditor), installer.WithHTTPClient(httpClient))
	if err != nil {
		return err
	}

	installerClient.ApiClient = client

	tokenCreationSpec := installer.TokenCreationSpec{
		Username: &installerClient.username,
		Password: &installerClient.password,
	}

	res, err := client.CreateTokenWithResponse(context.Background(), tokenCreationSpec)
	if err != nil {
		return err
	}

	tokenPair, vcfErr := GetResponseAs[installer.TokenPair](res)
	if vcfErr != nil {
		return errors.New(*vcfErr.Message)
	}
	installerClient.accessToken = tokenPair.AccessToken
	installerClient.lastRefreshTime = time.Now()
	installerClient.isRefreshing = false

	return nil
}

func (installerClient *InstallerClient) GetResourceIdAssociatedWithTask(ctx context.Context, taskId, resourceType string) (string, error) {
	task, err := installerClient.getTask(ctx, taskId)
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

func (installerClient *InstallerClient) getTask(ctx context.Context, taskId string) (*installer.Task, error) {
	apiClient := installerClient.ApiClient
	res, err := apiClient.GetTaskWithResponse(ctx, taskId)
	task, vcfErr := GetResponseAs[installer.Task](res)
	if err != nil || vcfErr != nil {
		log.Println("error = ", err)
		return nil, err
	}

	return task, nil
}
