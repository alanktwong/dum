package playbook

import (
	"awong/dotfiles/pkg/external"
	"awong/dotfiles/pkg/logging"
	"context"
	"fmt"
	"strings"
)

// VsCodePluginTask installs VS Code extensions from the command line
// Use 'id' as the extension/plugin identifier.
//
// code --list-extensions
// code --install-extension ms-vscode.cpptools
// code --uninstall-extension ms-vscode.csharp
// See https://stackoverflow.com/a/34339780/17205
type VsCodePluginTask struct {
	Attributes
	Code external.Code
	Log  logging.Logger
}

// NewVsCodePluginTask constructs a VsCodePluginTask.
func NewVsCodePluginTask(attributes *Attributes) (*VsCodePluginTask, error) {
	if attributes == nil {
		return nil, fmt.Errorf("attributes cannot be nil")
	}
	if attributes.ID == "" {
		return nil, fmt.Errorf("task ID cannot be empty")
	}
	return &VsCodePluginTask{
		Attributes: *attributes,
		Code:       external.NewCode(),
		Log:        logging.Log(),
	}, nil
}

// GetAttributes implements Attributable.
func (t *VsCodePluginTask) GetAttributes() Attributes {
	return t.Attributes
}

// GetID implements Identifiable.
func (t *VsCodePluginTask) GetID() string {
	return t.ID
}

// IsEnabled implements Enableable.
func (t *VsCodePluginTask) IsEnabled() bool {
	return t.Enabled
}

// Install implements Installer.
func (t *VsCodePluginTask) Install(ctx context.Context, input *Input) (*TaskResult, error) {
	t.Log.Debugf("%s START ... play: %s taskID: %s", TaskEllipsis, input.Play, t.ID)
	if !t.Enabled {
		return t.CreateTaskResult(input, false)
	}
	extList, err := t.Code.ListExtensions(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to list VS Code extensions: %w", err)
	}
	if strings.Contains(extList, t.ID) {
		t.Log.Infof("%s %s: %s is already installed", TaskEllipsis, input.Play, t.ID)
		return t.CreateTaskResult(input, true)
	}
	t.Log.Infof("%s %s: code --install-extension %s", TaskEllipsis, input.Play, t.ID)
	if !input.DryRun {
		if err = t.Code.InstallExtension(ctx, t.ID); err != nil {
			return nil, fmt.Errorf("%s %s: code --install-extension %s failed: %w", TaskEllipsis, input.Play, t.ID, err)
		}
	}
	return t.CreateTaskResult(input, true)
}

// List implements Lister.
func (t *VsCodePluginTask) List(_ context.Context, input *Input) (*TaskResult, error) {
	err := t.Log.Printlnf("%v code --install-extensions %s", TaskEllipsis, t.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to list code: %v", err)
	}
	return t.CreateTaskResult(input, t.Enabled)
}
