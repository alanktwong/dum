// Package tasks provides Task types and related utilities.
package tasks

import (
	ext "alanktwong/dum/internal/external"
	l "alanktwong/dum/internal/logging"
	ty "alanktwong/dum/internal/types"
	"context"
	"fmt"
)

// JetBrainsPluginTask installs a JetBrains IDE plugin.
type JetBrainsPluginTask struct {
	ty.Attributes
	Apps      []string
	Utils     ext.Ext
	JetBrains ext.JetBrainsApp
	Log       l.Logger
}

// NewJetBrainsPluginTask returns a new JetBrainsPluginTask for installing a JetBrains IDE plugin.
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

// GetAttributes returns the Attributes.
func (t *JetBrainsPluginTask) GetAttributes() ty.Attributes {
	return t.Attributes
}

// GetID returns the ID.
func (t *JetBrainsPluginTask) GetID() string {
	return t.ID
}

// IsEnabled returns whether the task is enabled.
func (t *JetBrainsPluginTask) IsEnabled() bool {
	return t.Enabled
}

// Install installs the task.
func (t *JetBrainsPluginTask) Install(ctx context.Context, input *ty.TaskInput) (*ty.TaskResult, error) {
	t.Log.Debugf("%s %s START ... play: %s taskID: %s", TaskEllipsis, TaskTypeName(t), input.Play, t.ID)
	if !t.Enabled {
		if input.DryRun {
			for _, app := range t.activeApps() {
				err := t.Log.Printlnf("%v %s: %s installPlugins %s", TaskEllipsis, TaskTypeName(t), app, t.ID)
				if err != nil {
					return nil, fmt.Errorf("failed to log disabled jetbrains: %w", err)
				}
			}
		}
		result, err := t.CreateTaskResult(input, false)
		if err != nil {
			return nil, fmt.Errorf("failed to create task result: %w", err)
		}
		return result, nil
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
	result, err := t.CreateTaskResult(input, success)
	if err != nil {
		return nil, fmt.Errorf("failed to create task result: %w", err)
	}
	return result, nil
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
