// © Broadcom. All Rights Reserved.
// The term “Broadcom” refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

package credentials

import (
	"context"
	md52 "crypto/md5"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"log"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	utils "github.com/vmware/terraform-provider-vcf/internal/resource_utils"
	"github.com/vmware/vcf-sdk-go/vcf"

	"github.com/vmware/terraform-provider-vcf/internal/api_client"
)

func ReadCredentials(ctx context.Context, data *schema.ResourceData, apiClient *vcf.ClientWithResponses) ([]vcf.Credential, error) {
	getCredentialsParam := &vcf.GetCredentialsParams{}
	resourceName, nameOk := data.Get("resource_name").(string)
	if nameOk && len(resourceName) > 0 {
		getCredentialsParam.ResourceName = &resourceName
	}

	ip, ipOk := data.Get("resource_ip").(string)
	if ipOk && len(ip) > 0 {
		getCredentialsParam.ResourceIp = &ip
	}

	resType, resTypeOK := data.Get("resource_type").(string)
	if resTypeOK && len(resType) > 0 {
		getCredentialsParam.ResourceType = &resType
	}

	domainName, domainNameOk := data.Get("domain_name").(string)
	if domainNameOk && len(domainName) > 0 {
		getCredentialsParam.DomainName = &domainName
	}

	accountType, accountTypeOk := data.Get("account_type").(string)
	if accountTypeOk && len(accountType) > 0 {
		getCredentialsParam.AccountType = &accountType
	}

	page, pageOk := data.Get("page").(int)
	if pageOk && page > 0 {
		pageNum := strconv.Itoa(page)
		getCredentialsParam.PageNumber = &pageNum
	}

	pageSize, pageSizeOk := data.Get("page_size").(int)
	if pageSizeOk && pageSize > 0 {
		pageSizeNum := strconv.Itoa(pageSize)
		getCredentialsParam.PageSize = &pageSizeNum
	}

	res, err := apiClient.GetCredentialsWithResponse(ctx, getCredentialsParam)
	if err != nil {
		return nil, err
	}
	pageOfCredential, vcfErr := api_client.GetResponseAs[vcf.PageOfCredential](res)
	if vcfErr != nil {
		api_client.LogError(vcfErr, ctx)
		return nil, errors.New(*vcfErr.Message)
	}

	return *pageOfCredential.Elements, nil
}

func FlattenCredentials(creds []vcf.Credential) []map[string]interface{} {
	if creds == nil {
		return []map[string]interface{}{}
	}

	credsArray := make([]map[string]interface{}, 0)

	for _, entry := range creds {
		entryMap := map[string]interface{}{
			"id":                entry.Id,
			"account_type":      entry.AccountType,
			"creation_time":     entry.CreationTimestamp,
			"credential_type":   entry.CredentialType,
			"modification_time": entry.ModificationTimestamp,
			"user_name":         entry.Username,
			"password":          entry.Password,
			"resource": []map[string]string{{
				"id":     entry.Resource.ResourceId,
				"domain": *entry.Resource.DomainName,
				"ip":     entry.Resource.ResourceIp,
				"name":   entry.Resource.ResourceName,
				"type":   entry.Resource.ResourceType,
			}},
		}

		if entry.AutoRotatePolicy != nil {
			entryMap["auto_rotate_frequency_days"] = entry.AutoRotatePolicy.FrequencyInDays
			entryMap["auto_rotate_next_schedule"] = entry.AutoRotatePolicy.NextSchedule
		}

		credsArray = append(credsArray, entryMap)
	}

	return credsArray
}

func CreateAutoRotatePolicy(ctx context.Context, data *schema.ResourceData, meta interface{}) error {
	log.Print("[DEBUG] About to create password rotation schedule")

	autoRotationEnabled := data.Get("enable_auto_rotation").(bool)
	autoRotationDays := data.Get("auto_rotate_days").(int)
	if !autoRotationEnabled {
		autoRotationDays = 0
	}
	resourceType := data.Get("resource_type").(string)
	resourceId := data.Get("resource_id").(string)
	resourceName := data.Get("resource_name").(string)
	userName := data.Get("user_name").(string)

	credentialsUpdateSpec, err := makeAutoRotatePolicySpec(autoRotationEnabled, int32(autoRotationDays), resourceName, resourceId, resourceType, userName)
	if err != nil {
		return err
	}

	sddcClient := meta.(*api_client.SddcManagerClient)
	return executeCredentialsUpdate(ctx, credentialsUpdateSpec, sddcClient)

}

func RotatePasswords(ctx context.Context, data *schema.ResourceData, meta interface{}) error {
	return mutatePassword(ctx, data, meta, Rotate)
}

func UpdatePasswords(ctx context.Context, data *schema.ResourceData, meta interface{}) error {
	return mutatePassword(ctx, data, meta, Update)
}

func mutatePassword(ctx context.Context, data *schema.ResourceData, meta interface{}, operationName string) error {
	resourceName := data.Get("resource_name").(string)
	resourceType := data.Get("resource_type").(string)

	creds := data.Get("credentials").([]interface{})

	credentialsUpdateSpec := makeCredentialsChangeSpec(resourceType, resourceName, creds, operationName)

	sddcClient := meta.(*api_client.SddcManagerClient)

	return executeCredentialsUpdate(ctx, credentialsUpdateSpec, sddcClient)
}

func RemoveAutoRotatePolicy(ctx context.Context, data *schema.ResourceData, meta interface{}) error {
	resourceType := data.Get("resource_type").(string)
	resourceId := data.Get("resource_id").(string)
	resourceName := data.Get("resource_name").(string)
	userName := data.Get("user_name").(string)

	credentialsUpdateSpec, err := makeAutoRotatePolicySpec(false, 0, resourceName, resourceId, resourceType, userName)
	if err != nil {
		return err
	}

	sddcClient := meta.(*api_client.SddcManagerClient)
	return executeCredentialsUpdate(ctx, credentialsUpdateSpec, sddcClient)
}

func makeAutoRotatePolicySpec(autoRotateEnabled bool, autoRotateDays int32, resourceName string, resourceId string, resourceType string, userName string) (*vcf.CredentialsUpdateSpec, error) {
	if len(resourceId) == 0 && len(resourceName) == 0 {
		log.Print("[ERROR] resource_id or resource_name attributes must be set")
		return nil, errors.New("resource_id or resource_name must be set")
	}

	operation := ConfigAutoRotate
	return &vcf.CredentialsUpdateSpec{
		AutoRotatePolicy: &vcf.AutoRotateCredentialPolicyInputSpec{
			EnableAutoRotatePolicy: autoRotateEnabled,
			FrequencyInDays:        &autoRotateDays,
		},
		Elements: []vcf.ResourceCredentials{
			{
				ResourceId:   &resourceId,
				ResourceName: &resourceName,
				ResourceType: resourceType,
				Credentials: []vcf.BaseCredential{{
					Username: userName,
				}},
			},
		},
		OperationType: operation,
	}, nil
}

func executeCredentialsUpdate(ctx context.Context, updateSpec *vcf.CredentialsUpdateSpec, sddcClient *api_client.SddcManagerClient) error {
	apiClient := sddcClient.ApiClient
	res, err := apiClient.UpdateOrRotatePasswordsWithResponse(ctx, *updateSpec)
	if err != nil {
		return err
	}
	task, vcfErr := api_client.GetResponseAs[vcf.Task](res)
	if vcfErr != nil {
		api_client.LogError(vcfErr, ctx)
		return errors.New(*vcfErr.Message)
	}
	return api_client.NewTaskTracker(ctx, apiClient, *task.Id).WaitForTask()
}

func CreatePasswordChangeID(data *schema.ResourceData, operation string) (string, error) {
	params := []string{
		operation,
		data.Get("resource_name").(string),
		data.Get("resource_type").(string),
	}

	authDetails := data.Get("credentials").([]interface{})
	for _, authDetail := range authDetails {
		entry := authDetail.(map[string]interface{})
		credentialType := entry["credential_type"].(string)
		username := entry["user_name"].(string)
		params = append(params, fmt.Sprintf("credential:%s|%s", credentialType, username))
	}

	return HashFields(params)
}

func HashFields(fields []string) (string, error) {
	md5 := md52.New()
	_, err := io.WriteString(md5, strings.Join(fields, ""))

	if err != nil {
		return "", err
	}

	return hex.EncodeToString(md5.Sum(nil)), nil
}

func makeCredentialsChangeSpec(resourceType string, resourceName string, credentialsList []interface{}, operation string) *vcf.CredentialsUpdateSpec {
	baseCredentials := make([]vcf.BaseCredential, 0)
	for _, listEntry := range credentialsList {
		entry := listEntry.(map[string]interface{})
		userName := entry["user_name"].(string)
		password, passwordOk := entry["password"].(string)
		credential := vcf.BaseCredential{
			Username:       userName,
			CredentialType: utils.ToStringPointer(entry["credential_type"]),
		}

		if passwordOk && len(password) > 0 {
			credential.Password = &password
		}

		baseCredentials = append(baseCredentials, credential)
	}

	return &vcf.CredentialsUpdateSpec{
		Elements: []vcf.ResourceCredentials{
			{
				ResourceName: &resourceName,
				ResourceType: resourceType,
				Credentials:  baseCredentials,
			},
		},
		OperationType: operation,
	}

}
