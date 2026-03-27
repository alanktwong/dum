// Package tasks provides Task types and related utilities.
package tasks

import (
	ext "awong/dotfiles/internal/external"
	l "awong/dotfiles/internal/logging"
	ty "awong/dotfiles/internal/types"
	"context"
	"fmt"
	"strings"
)

// MasTask installs a Mac App Store application.
type MasTask struct {
	ty.Attributes
	Mas ext.Mas
	Log l.Logger
}

// NewMasTask returns a new MasTask for installing a Mac App Store application.
func NewMasTask(attributes *ty.Attributes) (*MasTask, error) {
	if attributes == nil {
		return nil, fmt.Errorf("attributes cannot be nil")
	}
	if attributes.ID == "" {
		return nil, fmt.Errorf("task ID cannot be empty")
	}
	return &MasTask{
		Attributes: *attributes,
		Mas:        ext.NewMas(),
		Log:        l.Log(),
	}, nil
}

// GetAttributes returns the Attributes.
func (t *MasTask) GetAttributes() ty.Attributes {
	return t.Attributes
}

// GetID returns the ID.
func (t *MasTask) GetID() string {
	return t.ID
}

// IsEnabled returns whether the task is enabled.
func (t *MasTask) IsEnabled() bool {
	return t.Enabled
}

// Install installs the task.
func (t *MasTask) Install(ctx context.Context, input *ty.TaskInput) (*ty.TaskResult, error) {
	t.Log.Debugf("%s START ... play: %s taskID: %s", TaskEllipsis, input.Play, t.ID)
	if t.Description == "" {
		result, err := t.CreateTaskResult(input, false)
		if err != nil {
			return nil, fmt.Errorf("failed to create task result: %w", err)
		}
		return result, nil
	}
	if !t.Enabled {
		result, err := t.CreateTaskResult(input, false)
		if err != nil {
			return nil, fmt.Errorf("failed to create task result: %w", err)
		}
		return result, nil
	}
	masList, err := t.Mas.List(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to list mas apps: %w", err)
	}
	count := strings.Count(masList, t.ID)
	if count > 0 {
		t.Log.Infof("%s %s: %s is already installed", TaskEllipsis, input.Play, t.ID)
		result, err := t.CreateTaskResult(input, false)
		if err != nil {
			return nil, fmt.Errorf("failed to create task result: %w", err)
		}
		return result, nil
	}
	t.Log.Infof("%s %s: mas install %s", TaskEllipsis, input.Play, t.ID)
	if !input.DryRun {
		if err := t.Mas.Install(ctx, t.ID); err != nil {
			return nil, fmt.Errorf("%s %s: mas install %s failed: %w", TaskEllipsis, input.Play, t.ID, err)
		}
	}
	result, err := t.CreateTaskResult(input, true)
	if err != nil {
		return nil, fmt.Errorf("failed to create task result: %w", err)
	}
	return result, nil
}

// List lists the task.
func (t *MasTask) List(_ context.Context, input *ty.TaskInput) (*ty.TaskResult, error) {
	err := t.Log.Printlnf("%v mas install %s ... desc: %s", TaskEllipsis, t.ID, t.Description)
	if err != nil {
		return nil, fmt.Errorf("failed to list mas: %w", err)
	}
	result, err := t.CreateTaskResult(input, true)
	if err != nil {
		return nil, fmt.Errorf("failed to create task result: %w", err)
	}
	return result, nil
}
