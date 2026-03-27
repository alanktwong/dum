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

// VsCodePluginTask installs a VS Code extension.
type VsCodePluginTask struct {
	ty.Attributes
	Code ext.Code
	Log  l.Logger
}

// NewVsCodePluginTask returns a new VsCodePluginTask for installing a VS Code extension.
func NewVsCodePluginTask(attributes *ty.Attributes) (*VsCodePluginTask, error) {
	if attributes == nil {
		return nil, fmt.Errorf("attributes cannot be nil")
	}
	if attributes.ID == "" {
		return nil, fmt.Errorf("task ID cannot be empty")
	}
	return &VsCodePluginTask{
		Attributes: *attributes,
		Code:       ext.NewCode(),
		Log:        l.Log(),
	}, nil
}

// GetAttributes returns the Attributes.
func (t *VsCodePluginTask) GetAttributes() ty.Attributes {
	return t.Attributes
}

// GetID returns the ID.
func (t *VsCodePluginTask) GetID() string {
	return t.ID
}

// IsEnabled returns whether the task is enabled.
func (t *VsCodePluginTask) IsEnabled() bool {
	return t.Enabled
}

// Install installs the task.
func (t *VsCodePluginTask) Install(ctx context.Context, input *ty.TaskInput) (*ty.TaskResult, error) {
	t.Log.Debugf("%s START ... play: %s taskID: %s", TaskEllipsis, input.Play, t.ID)
	if !t.Enabled {
		result, err := t.CreateTaskResult(input, false)
		if err != nil {
			return nil, fmt.Errorf("failed to create task result: %w", err)
		}
		return result, nil
	}
	extList, err := t.Code.ListExtensions(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to list VS Code extensions: %w", err)
	}
	if strings.Contains(extList, t.ID) {
		t.Log.Infof("%s %s: %s is already installed", TaskEllipsis, input.Play, t.ID)
		result, err := t.CreateTaskResult(input, true)
		if err != nil {
			return nil, fmt.Errorf("failed to create task result: %w", err)
		}
		return result, nil
	}
	t.Log.Infof("%s %s: code --install-extension %s", TaskEllipsis, input.Play, t.ID)
	if !input.DryRun {
		if err = t.Code.InstallExtension(ctx, t.ID); err != nil {
			return nil, fmt.Errorf("%s %s: code --install-extension %s failed: %w", TaskEllipsis, input.Play, t.ID, err)
		}
	}
	result, err := t.CreateTaskResult(input, true)
	if err != nil {
		return nil, fmt.Errorf("failed to create task result: %w", err)
	}
	return result, nil
}

// List lists the task.
func (t *VsCodePluginTask) List(_ context.Context, input *ty.TaskInput) (*ty.TaskResult, error) {
	err := t.Log.Printlnf("%v code --install-extensions %s", TaskEllipsis, t.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to list code: %w", err)
	}
	result, err := t.CreateTaskResult(input, t.Enabled)
	if err != nil {
		return nil, fmt.Errorf("failed to create task result: %w", err)
	}
	return result, nil
}
