package tasks

import (
	l "awong/dotfiles/pkg/logging"
	ty "awong/dotfiles/pkg/types"
	"context"
	"fmt"
)

type FunctionTask struct {
	ty.Attributes
	Registry map[string]Installer
	Log      l.Logger
}

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

func NewFunctionRegistry() map[string]Installer {
	return map[string]Installer{
		"install_bash":                 NewBashInstaller(),
		"install_starship":             NewStarshipInstaller(),
		"install_sdkman":               NewSdkmanInstaller(),
		"install_vim":                  NewVimInstaller(),
		"install_case_sensitive_mount": NewMountInstaller(),
		"install_test":                 NewTestInstaller(),
	}
}

func (t *FunctionTask) GetAttributes() ty.Attributes {
	return t.Attributes
}

func (t *FunctionTask) GetID() string {
	return t.ID
}

func (t *FunctionTask) IsEnabled() bool {
	return t.Enabled
}

func (t *FunctionTask) Install(ctx context.Context, input *ty.TaskInput) (*ty.TaskResult, error) {
	t.Log.Debugf("%s START ... play: %s taskID: %s", TaskEllipsis, input.Play, t.ID)
	if !t.Enabled {
		return t.CreateTaskResult(input, false)
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
	return t.CreateTaskResult(input, true)
}

func (t *FunctionTask) List(_ context.Context, input *ty.TaskInput) (*ty.TaskResult, error) {
	err := t.Log.Printlnf("%v function -> %s", TaskEllipsis, t.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to list function: %v", err)
	}
	return t.CreateTaskResult(input, true)
}
