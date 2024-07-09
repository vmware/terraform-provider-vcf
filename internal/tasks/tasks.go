// Copyright 2024 Broadcom. All Rights Reserved.
// SPDX-License-Identifier: MPL-2.0

package tasks

import (
	"context"
	"errors"
	"fmt"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/vmware/terraform-provider-vcf/internal/api_client"
	"github.com/vmware/terraform-provider-vcf/internal/constants"
	"github.com/vmware/vcf-sdk-go/client/tasks"
	"github.com/vmware/vcf-sdk-go/models"
	"time"
)

const (
	defaultPollingInterval = 20 * time.Second
)

// resource types
const (
	ResourceTypeSddcManager   ResourceType = "SDDC_MANAGER"
	ResourceTypePsc           ResourceType = "PSC"
	ResourceTypeVcenre        ResourceType = "VCENTER"
	ResourceTypeNsxManager    ResourceType = "NSX_MANAGER"
	ResourceTypeVra           ResourceType = "VRA"
	ResourceTypeVrli          ResourceType = "VRLI"
	ResourceTypeVrops         ResourceType = "VROPS"
	ResourceTypeVrslcm        ResourceType = "VRSLCM"
	ResourceTypeVxRailManager ResourceType = "VXRAIL_MANAGER"
)

type ResourceType string

type Manager struct {
	client          *api_client.SddcManagerClient
	pollingInterval time.Duration
}

func NewManager(client *api_client.SddcManagerClient) *Manager {
	return &Manager{
		client:          client,
		pollingInterval: defaultPollingInterval,
	}
}

func NewManagerWithCustomInterval(client *api_client.SddcManagerClient, pollingInterval time.Duration) *Manager {
	return &Manager{
		client:          client,
		pollingInterval: pollingInterval,
	}
}

func (manager *Manager) WaitForCompletion(ctx context.Context, taskId string) error {
	tflog.Info(ctx, fmt.Sprintf("Getting status of task %s", taskId))
	ticker := time.NewTicker(manager.pollingInterval)
	for {
		select {
		case <-ctx.Done():
		case <-ticker.C:
			task, err := manager.getTask(ctx, taskId)
			if err != nil {
				return err
			}

			if task.Status == "In Progress" || task.Status == "Pending" || task.Status == "IN_PROGRESS" {
				logRunningTask(ctx, task)
				continue
			}

			if task.Status == "Failed" || task.Status == "Cancelled" {
				errorMsg := fmt.Sprintf("Task with ID = %s , Name: %q Type: %q is in state %s", taskId, task.Name, task.Type, task.Status)
				tflog.Error(ctx, errorMsg)

				return errors.New(errorMsg)
			}

			tflog.Info(ctx, fmt.Sprintf("Task with ID = %s, Name: %q is in state %s, completed at %s", taskId, task.Name, task.Status, task.CompletionTimestamp))
			return nil
		}
	}
}

func (manager *Manager) GetResourceIdAssociatedWithTask(ctx context.Context, taskId string, resourceType ResourceType) (string, error) {
	task, err := manager.getTask(ctx, taskId)
	if err != nil {
		return "", err
	}
	if len(task.Resources) == 0 {
		return "", fmt.Errorf("no resources associated with Task with ID %q", taskId)
	}
	for _, resource := range task.Resources {
		if *resource.Type == string(resourceType) {
			return *resource.ResourceID, nil
		}
	}
	return "", fmt.Errorf("task %q did not contain resources of type %q", taskId, resourceType)
}

func (manager *Manager) getTask(ctx context.Context, taskId string) (*models.Task, error) {
	apiClient := manager.client.ApiClient
	params := tasks.NewGetTaskParamsWithTimeout(constants.DefaultVcfApiCallTimeout).
		WithContext(ctx)
	params.ID = taskId

	result, err := apiClient.Tasks.GetTask(params)

	if err != nil {
		return nil, err
	}

	return result.Payload, nil
}

func logRunningTask(ctx context.Context, task *models.Task) {
	if task.SubTasks == nil {
		// no subtasks, log the message of the root task
		messagePack := task.LocalizableDescriptionPack
		if messagePack != nil && messagePack.Message != "" {
			tflog.Info(ctx, fmt.Sprintf("Running task: %s", task.LocalizableDescriptionPack.Message))
		}
	} else {
		// loop over the subtasks, find the one that is running and log it
		for _, subtask := range task.SubTasks {
			if subtask.Status == "IN_PROGRESS" {
				tflog.Info(ctx, fmt.Sprintf("Running task: %s", subtask.Description))
			}
		}
	}
}
