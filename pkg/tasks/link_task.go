// Package tasks provides Task types and related utilities.
package tasks

import (
	ext "awong/dotfiles/pkg/external"
	l "awong/dotfiles/pkg/logging"
	ty "awong/dotfiles/pkg/types"
	"context"
	"fmt"
	"path/filepath"
)

// LinkTask creates a symbolic link.
type LinkTask struct {
	ty.Attributes
	Root   string
	Target string
	Utils  ext.Ext
	Log    l.Logger
}

// NewLinkTask returns a new LinkTask for creating a symbolic link.
func NewLinkTask(attributes *ty.Attributes, root, target string) (*LinkTask, error) {
	if attributes == nil {
		return nil, fmt.Errorf("attributes cannot be nil")
	}
	if attributes.ID == "" {
		return nil, fmt.Errorf("task ID cannot be empty")
	}
	return &LinkTask{
		Attributes: *attributes,
		Root:       root,
		Target:     target,
		Utils:      ext.NewExt(),
		Log:        l.Log(),
	}, nil
}

// GetAttributes returns the Attributes.
func (t *LinkTask) GetAttributes() ty.Attributes {
	return t.Attributes
}

// GetID returns the ID.
func (t *LinkTask) GetID() string {
	return t.ID
}

// IsEnabled returns whether the task is enabled.
func (t *LinkTask) IsEnabled() bool {
	return t.Enabled
}

// Install installs the task.
func (t *LinkTask) Install(ctx context.Context, input *ty.TaskInput) (*ty.TaskResult, error) {
	t.Log.Debugf("%s START ... play: %s taskID: %s", TaskEllipsis, input.Play, t.ID)
	if t.Root == "" {
		result, err := t.CreateTaskResult(input, false)
		if err != nil {
			return nil, fmt.Errorf("failed to create task result: %w", err)
		}
		return result, nil
	}
	target := t.ProvideTarget()
	if target == "" {
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
	rootPath, err := t.Utils.ExpandUser(t.Root)
	if err != nil {
		return nil, fmt.Errorf("failed to expand root path %s: %w", t.Root, err)
	}
	targetPath := filepath.Join(rootPath, target)
	if t.Utils.IsSymlink(targetPath) {
		t.Log.Infof("%v %v: %v -> %s/%s already linked", TaskEllipsis, input.Play, t.ID, t.Root, target)
	}
	t.Log.Infof("%v %v: planning to link %v -> %s/%s", TaskEllipsis, input.Play, t.ID, t.Root, target)
	if !input.DryRun {
		if t.Utils.SoftLink(ctx, rootPath, t.ID, target, t.Sudo) != nil {
			return nil, fmt.Errorf("failed to link %v -> %s/%s: %v", t.ID, t.Root, target, err)
		}
	}
	result, err := t.CreateTaskResult(input, true)
	if err != nil {
		return nil, fmt.Errorf("failed to create task result: %w", err)
	}
	return result, nil
}

// List lists the task.
func (t *LinkTask) List(_ context.Context, input *ty.TaskInput) (*ty.TaskResult, error) {
	target := t.ProvideTarget()
	err := t.Log.Printlnf("%v linking %v -> %s/%s", TaskEllipsis, t.ID, t.Root, target)
	if err != nil {
		return nil, fmt.Errorf("failed to list link: %v", err)
	}
	result, err := t.CreateTaskResult(input, true)
	if err != nil {
		return nil, fmt.Errorf("failed to create task result: %w", err)
	}
	return result, nil
}

// ProvideTarget returns the target path for the link.
func (t *LinkTask) ProvideTarget() string {
	if t.Target != "" {
		return t.Target
	}
	return filepath.Base(t.ID)
}
