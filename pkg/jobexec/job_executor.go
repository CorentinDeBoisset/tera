package jobexec

import (
	"context"
	"sync"

	"github.com/corentindeboisset/tera/pkg/cfg"
	"github.com/corentindeboisset/tera/pkg/cmdrunr"
)

type TaskState int

const (
	STATE_NOT_STARTED TaskState = iota
	STATE_RUNNING
	STATE_SUCCESSFUL
	STATE_FAILED
)

type Stater struct {
	mtx   sync.Mutex
	state TaskState
}

func (s *Stater) SetState(newState TaskState) {
	s.mtx.Lock()
	defer s.mtx.Unlock()
	s.state = newState
}

func (s *Stater) State() TaskState {
	s.mtx.Lock()
	defer s.mtx.Unlock()
	return s.state
}

type TaskStatus struct {
	Stater

	Output cmdrunr.SafeBuffer
}

type StepStatus struct {
	Stater

	BeforeHooks *TaskStatus
	Tasks       []TaskStatus
	AfterHooks  *TaskStatus
}

func executeJob(ctx context.Context, basePath string, config *cfg.JobConfig, stepStatuses []StepStatus, readyToDisplay, done chan struct{}) {
	// First, initialize the status structs
	for stepIdx, step := range config.Steps {
		if len(step.RunBefore) > 0 {
			stepStatuses[stepIdx].BeforeHooks = &TaskStatus{}
		}
		if len(step.RunAfter) > 0 {
			stepStatuses[stepIdx].AfterHooks = &TaskStatus{}
		}
		stepStatuses[stepIdx].Tasks = make([]TaskStatus, len(step.Tasks))
	}

	// Notify that the statuses are ready to be displayed
	close(readyToDisplay)
	defer close(done)

	for stepIdx, step := range config.Steps {
		stepStatuses[stepIdx].SetState(STATE_RUNNING)

		if stepStatuses[stepIdx].BeforeHooks != nil {
			stepStatuses[stepIdx].BeforeHooks.SetState(STATE_RUNNING)
		}

		for _, hook := range step.RunBefore {
			if !cmdrunr.RunCommand(ctx, basePath, hook.Path, hook.Cmd, &stepStatuses[stepIdx].BeforeHooks.Output, 120, 24) {
				stepStatuses[stepIdx].BeforeHooks.SetState(STATE_FAILED)
				stepStatuses[stepIdx].SetState(STATE_FAILED)
				return
			}
		}

		if stepStatuses[stepIdx].BeforeHooks != nil {
			stepStatuses[stepIdx].BeforeHooks.SetState(STATE_SUCCESSFUL)
		}

		// Within each step, the tasks are run asynchronously
		var taskWg sync.WaitGroup
		for taskIdx, task := range step.Tasks {
			taskWg.Go(func() {
				runTask(ctx, basePath, task, &stepStatuses[stepIdx].Tasks[taskIdx])
			})
		}
		taskWg.Wait()

		// Read the status from the tasks
		tasksOk := true
		for taskIdx := range stepStatuses[stepIdx].Tasks {
			js := &stepStatuses[stepIdx].Tasks[taskIdx]
			if js.State() == STATE_FAILED {
				tasksOk = false
			}
		}
		if !tasksOk {
			stepStatuses[stepIdx].SetState(STATE_FAILED)
			return
		}
		if stepStatuses[stepIdx].AfterHooks != nil {
			stepStatuses[stepIdx].AfterHooks.SetState(STATE_RUNNING)
		}

		for _, hook := range step.RunAfter {
			if !cmdrunr.RunCommand(ctx, basePath, hook.Path, hook.Cmd, &stepStatuses[stepIdx].AfterHooks.Output, 120, 24) {
				stepStatuses[stepIdx].AfterHooks.SetState(STATE_FAILED)
				stepStatuses[stepIdx].SetState(STATE_FAILED)
				return
			}
		}

		if stepStatuses[stepIdx].AfterHooks != nil {
			stepStatuses[stepIdx].AfterHooks.SetState(STATE_SUCCESSFUL)
		}

		stepStatuses[stepIdx].SetState(STATE_SUCCESSFUL)
	}
}

func runTask(ctx context.Context, basePath string, config cfg.TaskConfig, taskStatus *TaskStatus) bool {
	taskStatus.SetState(STATE_RUNNING)

	for _, hook := range config.RunBefore {
		if !cmdrunr.RunCommand(ctx, basePath, hook.Path, hook.Cmd, &taskStatus.Output, 120, 24) {
			taskStatus.SetState(STATE_FAILED)
			return false
		}
	}

	// run config.Cmd
	if !cmdrunr.RunCommand(ctx, basePath, config.Path, config.Cmd, &taskStatus.Output, 120, 24) {
		taskStatus.SetState(STATE_FAILED)
		return false
	}

	for _, hook := range config.RunAfter {
		if !cmdrunr.RunCommand(ctx, basePath, hook.Path, hook.Cmd, &taskStatus.Output, 120, 24) {
			taskStatus.SetState(STATE_FAILED)
			return false
		}
	}

	taskStatus.SetState(STATE_SUCCESSFUL)

	return true
}
