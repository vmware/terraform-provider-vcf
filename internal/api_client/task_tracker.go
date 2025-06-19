// © Broadcom. All Rights Reserved.
// The term “Broadcom” refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

package api_client

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/vmware/vcf-sdk-go/vcf"
)

const (
	// Default polling interval for task tracking.
	defaultPollingInterval = 20 * time.Second

	// Task status constants.
	statusInProgress          = "In Progress"
	statusInProgressUppercase = "IN_PROGRESS"
	statusPending             = "Pending"
	statusFailed              = "Failed"
	statusCancelled           = "Cancelled"
	statusNotApplicable       = "NOT_APPLICABLE"
)

type TaskTracker struct {
	ctx             context.Context
	client          *vcf.ClientWithResponses
	taskId          string
	pollingInterval time.Duration
	completedTasks  map[string]bool
}

func NewTaskTracker(ctx context.Context, client *vcf.ClientWithResponses, taskId string) *TaskTracker {
	return &TaskTracker{
		ctx:             ctx,
		client:          client,
		taskId:          taskId,
		pollingInterval: defaultPollingInterval,
		completedTasks:  make(map[string]bool),
	}
}

func NewTaskTrackerWithCustomPollingInterval(ctx context.Context, client *vcf.ClientWithResponses, taskId string, pollingInterval time.Duration) *TaskTracker {
	tracker := NewTaskTracker(ctx, client, taskId)
	tracker.pollingInterval = pollingInterval
	return tracker
}

func (t *TaskTracker) WaitForTask() error {
	ticker := time.NewTicker(t.pollingInterval)
	defer ticker.Stop()

	for {
		select {
		case <-t.ctx.Done():
		case <-ticker.C:
			task, err := t.getTask()
			if err != nil {
				LogError(err, t.ctx)
				return errors.New(*err.Message)
			}

			t.logTask(*task)

			switch *task.Status {
			case statusInProgress, statusInProgressUppercase, statusPending:
				continue
			case statusFailed, statusCancelled:
				errorMsg := fmt.Sprintf("Task with ID = %s , Name: %q Type: %q is in state %s",
					*task.Id, *task.Name, *task.Type, *task.Status)
				tflog.Error(t.ctx, errorMsg)

				return errors.New(errorMsg)
			default:
				tflog.Info(t.ctx, fmt.Sprintf("Task with ID = %s , Name: %q Type: %q is in state %s",
					*task.Id, *task.Name, *task.Type, *task.Status))
				return nil
			}
		}
	}
}

func (t *TaskTracker) getTask() (*vcf.Task, *vcf.Error) {
	res, _ := t.client.GetTaskWithResponse(t.ctx, t.taskId)

	return GetResponseAs[vcf.Task](res)
}

func (t *TaskTracker) logTask(task vcf.Task) {
	if task.SubTasks == nil {
		messagePack := task.LocalizableDescriptionPack
		if messagePack != nil && messagePack.Message != nil && task.Status != nil &&
			t.shouldLog(*messagePack.Message, *task.Status) {
			t.log(*messagePack.Message, *task.Status)
		}
	} else if task.SubTasks != nil {
		for _, subtask := range *task.SubTasks {
			if t.shouldLog(*subtask.Description, *subtask.Status) {
				t.log(*subtask.Description, *subtask.Status)
				if subtask.Errors != nil {
					t.logErrors(*subtask.Errors)
				}
			}
		}
	}
}

func (t *TaskTracker) shouldLog(message, status string) bool {
	running := t.statusEqual(status, statusInProgressUppercase) ||
		t.statusEqual(status, statusInProgress) ||
		t.statusEqual(status, statusPending) ||
		t.statusEqual(status, statusNotApplicable)
	val, ok := t.completedTasks[message]
	return !running && (!val || !ok)
}

func (t *TaskTracker) log(message, status string) {
	tflog.Info(t.ctx, fmt.Sprintf("[%s] %s", status, message))
	t.completedTasks[message] = true
}

func (t *TaskTracker) logErrors(errors []vcf.Error) {
	for _, err := range errors {
		LogError(&err, t.ctx)
	}
}

func (t *TaskTracker) statusEqual(a, b string) bool {
	return strings.EqualFold(a, b)
}
