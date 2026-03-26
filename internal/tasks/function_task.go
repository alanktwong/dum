// Package tasks provides Task types and related utilities.
package tasks

import (
	l "awong/dotfiles/internal/logging"
	i "awong/dotfiles/internal/tasks/installer"
	ty "awong/dotfiles/internal/types"
	"context"
	"fmt"
)

// FunctionTask executes a function from a predefined registry.
type FunctionTask struct {
	ty.Attributes
	Registry map[string]i.Installer
	Log      l.Logger
}

// NewFunctionTask returns a new FunctionTask for executing a function from the registry.
func NewFunctionTask(attributes *ty.Attributes) (*FunctionTask, error) {
	if attributes == nil {
		return nil, fmt.Errorf("attributes cannot be nil")
	}
	if attributes.ID == "" {
		return nil, fmt.Errorf("task ID cannot be empty")
	}
	return &FunctionTask{
		Attributes: *attributes,
		Registry:   NewFunctionRegistry(),
		Log:        l.Log(),
	}, nil
}

// NewFunctionRegistry returns a map of available functions (installers) that can be executed.
func NewFunctionRegistry() map[string]i.Installer {
	return map[string]i.Installer{
		"install_bash":                 i.NewBashInstaller(),
		"install_starship":             i.NewStarshipInstaller(),
		"install_sdkman":               i.NewSdkmanInstaller(),
		"install_vim":                  i.NewVimInstaller(),
		"install_case_sensitive_mount": i.NewMountInstaller(),
		"install_test":                 i.NewTestInstaller(),
	}
}

// GetAttributes returns the Attributes.
func (t *FunctionTask) GetAttributes() ty.Attributes {
	return t.Attributes
}

// GetID returns the ID.
func (t *FunctionTask) GetID() string {
	return t.ID
}

// IsEnabled returns whether the task is enabled.
func (t *FunctionTask) IsEnabled() bool {
	return t.Enabled
}

// Install installs the task.
func (t *FunctionTask) Install(ctx context.Context, input *ty.TaskInput) (*ty.TaskResult, error) {
	t.Log.Debugf("%s START ... play: %s taskID: %s", TaskEllipsis, input.Play, t.ID)
	if !t.Enabled {
		result, err := t.CreateTaskResult(input, false)
		if err != nil {
			return nil, fmt.Errorf("failed to create task result: %w", err)
		}
		return result, nil
	}
	fn, ok := t.Registry[t.ID]
	if !ok {
		return nil, fmt.Errorf("function %v is not in the registry", t.ID)
	}
	t.Log.Infof("%s %s: executing function %s", TaskEllipsis, input.Play, t.ID)
	if !input.DryRun {
		res, err := fn.Install(ctx, input)
		if err != nil {
			return nil, fmt.Errorf("error executing function %s: %w", t.ID, err)
		}
		return res, nil
	}
	result, err := t.CreateTaskResult(input, true)
	if err != nil {
		return nil, fmt.Errorf("failed to create task result: %w", err)
	}
	return result, nil
}

// List lists the task.
func (t *FunctionTask) List(_ context.Context, input *ty.TaskInput) (*ty.TaskResult, error) {
	err := t.Log.Printlnf("%v function -> %s", TaskEllipsis, t.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to list function: %v", err)
	}
	result, err := t.CreateTaskResult(input, true)
	if err != nil {
		return nil, fmt.Errorf("failed to create task result: %w", err)
	}
	return result, nil
}
