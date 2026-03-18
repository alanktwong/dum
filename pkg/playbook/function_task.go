package playbook

import (
	"awong/dotfiles/pkg/logging"
	"context"
	"fmt"
)

// FunctionTask delegates to another installer to install a package.
// Each installer used for installation should be registered in the registry.
//
// Use 'id' to specify the function name.
type FunctionTask struct {
	Attributes
	// Registry holds a map of installers by their ID
	Registry map[string]Installer
	Log      logging.Logger
}

// NewFunctionTask constructs a FunctionTask.
func NewFunctionTask(attributes *Attributes) (*FunctionTask, error) {
	if attributes == nil {
		return nil, fmt.Errorf("attributes cannot be nil")
	}
	if attributes.ID == "" {
		return nil, fmt.Errorf("task ID cannot be empty")
	}
	return &FunctionTask{
		Attributes: *attributes,
		Registry:   NewFunctionRegistry(),
		Log:        logging.Log(),
	}, nil
}

// NewFunctionRegistry constructs a registry of installers keyed by its id.
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

// GetAttributes implements Attributable.
func (t *FunctionTask) GetAttributes() Attributes {
	return t.Attributes
}

// GetID implements Identifiable.
func (t *FunctionTask) GetID() string {
	return t.ID
}

// IsEnabled implements Enableable.
func (t *FunctionTask) IsEnabled() bool {
	return t.Enabled
}

// Install implements Installer.
func (t *FunctionTask) Install(ctx context.Context, input *Input) (*TaskResult, error) {
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

// List implements Lister.
func (t *FunctionTask) List(_ context.Context, input *Input) (*TaskResult, error) {
	err := t.Log.Printlnf("%v function -> %s", TaskEllipsis, t.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to list function: %v", err)
	}
	return t.CreateTaskResult(input, true)
}
