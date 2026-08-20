// Package tasks provides Task types and related utilities.
package tasks

import (
	"context"
	"fmt"

	ext "github.com/alanktwong/dum/internal/external"
	l "github.com/alanktwong/dum/internal/logging"
	ty "github.com/alanktwong/dum/internal/types"
)

// DirTask creates a directory.
type DirTask struct {
	ty.Attributes
	Utils ext.Ext
	Log   l.Logger
}

// NewDirTask returns a new DirTask for creating a directory.
func NewDirTask(attributes *ty.Attributes) (*DirTask, error) {
	if attributes == nil {
		return nil, fmt.Errorf("attributes cannot be nil")
	}
	if attributes.ID == "" {
		return nil, fmt.Errorf("task ID cannot be empty")
	}
	return &DirTask{
		Attributes: *attributes,
		Utils:      ext.NewExt(),
		Log:        l.Log(),
	}, nil
}

// GetAttributes returns the Attributes.
func (t *DirTask) GetAttributes() ty.Attributes {
	return t.Attributes
}

// GetID returns the ID.
func (t *DirTask) GetID() string {
	return t.ID
}

// IsEnabled returns whether the task is enabled.
func (t *DirTask) IsEnabled() bool {
	return t.Enabled
}

// Install installs the task.
func (t *DirTask) Install(ctx context.Context, input *ty.TaskInput) (*ty.TaskResult, error) {
	t.Log.Debugf("%s %s START ... play: %s taskID: %s", TaskEllipsis, TaskTypeName(t), input.Play, t.ID)
	if !t.Enabled {
		if input.DryRun {
			err := t.Log.Printlnf("%v %s: mkdir -p %s", TaskEllipsis, TaskTypeName(t), t.ID)
			if err != nil {
				return nil, fmt.Errorf("failed to log disabled dir: %w", err)
			}
		}
		result, err := t.CreateTaskResult(input, false)
		if err != nil {
			return nil, fmt.Errorf("failed to create task result: %w", err)
		}
		return result, nil
	}
	path, err := t.Utils.ExpandUser(t.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to expand user for %s: %w", t.ID, err)
	}
	if t.Utils.IsDir(path) {
		t.Log.Infof("%s %s: %s already exists", TaskEllipsis, input.Play, t.ID)
		result, err := t.CreateTaskResult(input, false)
		if err != nil {
			return nil, fmt.Errorf("failed to create task result: %w", err)
		}
		return result, nil
	}
	t.Log.Infof("%s %s: mkdir -p %s", TaskEllipsis, input.Play, t.ID)
	if !input.DryRun {
		err = t.Utils.CreateDirectory(ctx, path, t.Sudo)
		if err != nil {
			return nil, fmt.Errorf("failed to mkdir -p %v: %w", t.ID, err)
		}
	}
	result, err := t.CreateTaskResult(input, true)
	if err != nil {
		return nil, fmt.Errorf("failed to create task result: %w", err)
	}
	return result, nil
}
