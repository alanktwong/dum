package tasks

import (
	ext "awong/dotfiles/pkg/external"
	l "awong/dotfiles/pkg/logging"
	ty "awong/dotfiles/pkg/types"
	"context"
	"fmt"
)

type JetBrainsPluginTask struct {
	ty.Attributes
	Apps      []string
	Utils     ext.Ext
	JetBrains ext.JetBrainsApp
	Log       l.Logger
}

func NewJetBrainsPluginTask(attributes *ty.Attributes, apps []string) (*JetBrainsPluginTask, error) {
	if attributes == nil {
		return nil, fmt.Errorf("attributes cannot be nil")
	}
	if attributes.ID == "" {
		return nil, fmt.Errorf("task ID cannot be empty")
	}
	if len(apps) == 0 {
		return nil, fmt.Errorf("app cannot be empty")
	}
	return &JetBrainsPluginTask{
		Attributes: *attributes,
		Apps:       apps,
		Utils:      ext.NewExt(),
		JetBrains:  ext.NewJetBrains(),
		Log:        l.Log(),
	}, nil
}

func (t *JetBrainsPluginTask) GetAttributes() ty.Attributes {
	return t.Attributes
}

func (t *JetBrainsPluginTask) GetID() string {
	return t.ID
}

func (t *JetBrainsPluginTask) IsEnabled() bool {
	return t.Enabled
}

func (t *JetBrainsPluginTask) Install(ctx context.Context, input *ty.TaskInput) (*ty.TaskResult, error) {
	t.Log.Debugf("%s START ... play: %s taskID: %s", TaskEllipsis, input.Play, t.ID)
	if !t.Enabled {
		return t.CreateTaskResult(input, false)
	}

	success := true
	for _, app := range t.activeApps() {
		apps := input.JetBrainsApps
		ideName, ok := apps[app]
		if !ok {
			t.Log.Debugf("app %s not found in playbook %s", app, t.ID)
			success = false
			continue
		}
		if t.JetBrains.IsInstalled(ideName, t.ID) {
			t.Log.Infof("%v %v: plugin %v already installed for %s", TaskEllipsis, input.Play, app, t.ID)
			continue
		}
		if input.DryRun {
			t.Log.Infof("%v %v: %v installPlugins  %s", TaskEllipsis, input.Play, app, t.ID)
			continue
		}
		if err := t.JetBrains.Install(ctx, app, t.ID); err != nil {
			t.Log.Errorf("%v %v: %v installPlugins %s failed: %v", TaskEllipsis, input.Play, app, t.ID, err)
			success = false
		}
	}
	return t.CreateTaskResult(input, success)
}

func (t *JetBrainsPluginTask) List(_ context.Context, input *ty.TaskInput) (*ty.TaskResult, error) {
	for _, app := range t.activeApps() {
		err := t.Log.Printlnf("%v %s installPlugins  %s", TaskEllipsis, app, t.ID)
		if err != nil {
			return nil, fmt.Errorf("failed to list jetbrains: %v", err)
		}
	}
	return t.CreateTaskResult(input, t.Enabled)
}

func (t *JetBrainsPluginTask) activeApps() []string {
	result := make([]string, 0)
	for _, app := range t.Apps {
		if t.Utils.IsInstalled(app) {
			result = append(result, app)
		}
	}
	return result
}
