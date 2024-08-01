// Copyright 2024 Broadcom. All Rights Reserved.
// SPDX-License-Identifier: MPL-2.0

package api_client

import (
	"context"
	"errors"
	"fmt"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/vmware/terraform-provider-vcf/internal/constants"
	"github.com/vmware/vcf-sdk-go/client"
	"github.com/vmware/vcf-sdk-go/client/tasks"
	"github.com/vmware/vcf-sdk-go/models"
	"time"
)

const (
	defaultPollingInterval = time.Second * 20
)

type TaskTracker struct {
	ctx             context.Context
	client          *client.VcfClient
	taskId          string
	pollingInterval time.Duration

	completedTasks map[string]bool
}

func NewTaskTracker(ctx context.Context, client *client.VcfClient, taskId string) *TaskTracker {
	return &TaskTracker{
		ctx:             ctx,
		client:          client,
		taskId:          taskId,
		pollingInterval: defaultPollingInterval,
		completedTasks:  make(map[string]bool),
	}
}

func NewTaskTrackerWithCustomPollingInterval(ctx context.Context, client *client.VcfClient, taskId string, pollingInterval time.Duration) *TaskTracker {
	tracker := NewTaskTracker(ctx, client, taskId)
	tracker.pollingInterval = pollingInterval
	return tracker
}

func (t *TaskTracker) WaitForTask() error {
	ticker := time.NewTicker(t.pollingInterval)

	for {
		select {
		case <-t.ctx.Done():
		case <-ticker.C:
			task, err := t.getTask()
			if err != nil {
				return err
			}

			t.logTask(task)

			if task.Status == "In Progress" || task.Status == "Pending" || task.Status == "IN_PROGRESS" {
				return nil
			}

			if task.Status == "Failed" || task.Status == "Cancelled" {
				errorMsg := fmt.Sprintf("Task with ID = %s , Name: %q Type: %q is in state %s",
					task.ID, task.Name, task.Type, task.Status)
				tflog.Error(t.ctx, errorMsg)

				return errors.New(errorMsg)
			}

			tflog.Info(t.ctx, fmt.Sprintf("Task with ID = %s is in state %s, completed at %s",
				task.ID, task.Status, task.CompletionTimestamp))
			return nil
		}
	}
}

func (t *TaskTracker) getTask() (*models.Task, error) {
	getTaskParams := tasks.NewGetTaskParamsWithTimeout(constants.DefaultVcfApiCallTimeout).
		WithContext(t.ctx)
	getTaskParams.ID = t.taskId

	getTaskResult, err := t.client.Tasks.GetTask(getTaskParams)

	if err != nil {
		return nil, err
	}

	return getTaskResult.Payload, nil
}

func (t *TaskTracker) logTask(task *models.Task) {
	if task.SubTasks == nil {
		messagePack := task.LocalizableDescriptionPack
		if messagePack != nil && messagePack.Message != "" && t.shouldLog(messagePack.Message) {
			t.log(messagePack.Message, task.Status)
		}
	} else {
		for _, subtask := range task.SubTasks {
			t.logSubTask(subtask)
		}
	}
}

func (t *TaskTracker) logSubTask(task *models.SubTask) {
	if task.Status != "IN_PROGRESS" && task.Status != "PENDING" && task.Status != "NOT_APPLICABLE" {
		if t.shouldLog(task.Description) {
			t.log(task.Description, task.Status)
		}
	}
}

func (t *TaskTracker) shouldLog(message string) bool {
	val, ok := t.completedTasks[message]
	return !val || !ok
}

func (t *TaskTracker) log(message, status string) {
	tflog.Info(t.ctx, fmt.Sprintf("[%s] %s", status, message))
	t.completedTasks[message] = true
}
