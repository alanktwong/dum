// Package tasks provides Task types and related utilities.
package tasks

import (
	"context"
	"fmt"

	ext "github.com/alanktwong/dum/internal/external"
	l "github.com/alanktwong/dum/internal/logging"
	ty "github.com/alanktwong/dum/internal/types"
)

// BashTask executes bash commands or scripts.
type BashTask struct {
	ty.Attributes
	Command string
	Script  string
	Utils   ext.Ext
	Log     l.Logger
}

// NewBashTask returns a new BashTask for executing bash commands or scripts.
func NewBashTask(attributes *ty.Attributes, command, script string) (*BashTask, error) {
	if attributes == nil {
		return nil, fmt.Errorf("attributes cannot be nil")
	}
	if attributes.ID == "" {
		return nil, fmt.Errorf("task ID cannot be empty")
	}
	if command == "" && script == "" {
		return nil, fmt.Errorf("either command or script must be specified")
	}
	return &BashTask{
		Attributes: *attributes,
		Command:    command,
		Script:     script,
		Utils:      ext.NewExt(),
		Log:        l.Log(),
	}, nil
}

// GetAttributes returns the Attributes.
func (t *BashTask) GetAttributes() ty.Attributes {
	return t.Attributes
}

// GetID returns the ID.
func (t *BashTask) GetID() string {
	return t.ID
}

// IsEnabled returns whether the task is enabled.
func (t *BashTask) IsEnabled() bool {
	return t.Enabled
}

// Install installs the task.
func (t *BashTask) Install(ctx context.Context, input *ty.TaskInput) (*ty.TaskResult, error) {
	t.Log.Debugf("%s %s START ... play: %s taskID: %s", TaskEllipsis, TaskTypeName(t), input.Play, t.ID)
	if !t.Enabled {
		if input.DryRun {
			if err := t.logDisabledDryRun(); err != nil {
				return nil, err
			}
		}
		result, err := t.CreateTaskResult(input, false)
		if err != nil {
			return nil, fmt.Errorf("failed to create task result: %w", err)
		}
		return result, nil
	}

	cmd := t.Command
	if cmd == "" {
		cmd = t.Script
	}

	t.Log.Infof("%s %s: %s", TaskEllipsis, input.Play, cmd)
	if !input.DryRun {
		err := t.Utils.RunCommand(ctx, cmd, t.Sudo)
		if err != nil {
			return nil, fmt.Errorf("failed to run command %s: %w", t.ID, err)
		}
	}
	result, err := t.CreateTaskResult(input, true)
	if err != nil {
		return nil, fmt.Errorf("failed to create task result: %w", err)
	}
	return result, nil
}

// logDisabledDryRun logs the command a disabled task would run during a dry run.
func (t *BashTask) logDisabledDryRun() error {
	cmd := t.Command
	if cmd == "" {
		cmd = t.Script
	}
	err := t.Log.Printlnf("%s %s: %s", TaskEllipsis, TaskTypeName(t), cmd)
	if err != nil {
		return fmt.Errorf("failed to log disabled command: %w", err)
	}
	return nil
}
