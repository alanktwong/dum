package tasks

import (
	ext "awong/dotfiles/pkg/external"
	l "awong/dotfiles/pkg/logging"
	ty "awong/dotfiles/pkg/types"
	"context"
	"fmt"
	"strings"
)

type VsCodePluginTask struct {
	ty.Attributes
	Code ext.Code
	Log  l.Logger
}

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

func (t *VsCodePluginTask) GetAttributes() ty.Attributes {
	return t.Attributes
}

func (t *VsCodePluginTask) GetID() string {
	return t.ID
}

func (t *VsCodePluginTask) IsEnabled() bool {
	return t.Enabled
}

func (t *VsCodePluginTask) Install(ctx context.Context, input *ty.TaskInput) (*ty.TaskResult, error) {
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

func (t *VsCodePluginTask) List(_ context.Context, input *ty.TaskInput) (*ty.TaskResult, error) {
	err := t.Log.Printlnf("%v code --install-extensions %s", TaskEllipsis, t.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to list code: %v", err)
	}
	return t.CreateTaskResult(input, t.Enabled)
}
