package action

import (
	"errors"

	bosherr "github.com/cloudfoundry/bosh-utils/errors"

	boshtask "github.com/cloudfoundry/bosh-agent/v2/agent/task"
)

type GetTaskAction struct {
	taskService boshtask.Service
}

func NewGetTask(taskService boshtask.Service) (getTask GetTaskAction) {
	getTask.taskService = taskService
	return
}

func (a GetTaskAction) IsAsynchronous(_ ProtocolVersion) bool {
	return false
}

func (a GetTaskAction) IsPersistent() bool {
	return false
}

func (a GetTaskAction) IsLoggable() bool {
	return true
}

func (a GetTaskAction) Run(taskID string) (interface{}, error) {
	task, found := a.taskService.FindTaskWithID(taskID)
	if !found {
		return nil, bosherr.Errorf("Task with id %s could not be found", taskID)
	}

	if task.State == boshtask.StateRunning {
		return boshtask.StateValue{
			AgentTaskID: task.ID,
			State:       task.State,
		}, nil
	}

	if task.Error != nil {
		return task.Value, bosherr.WrapErrorf(task.Error, "Task %s result", taskID)
	}

	return task.Value, nil
}

func (a GetTaskAction) Resume() (interface{}, error) {
	return nil, errors.New("not supported")
}

func (a GetTaskAction) Cancel() error {
	return errors.New("not supported")
}
