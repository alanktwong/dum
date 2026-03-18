package playbook

import (
	util "awong/dotfiles/pkg/external"
	"awong/dotfiles/pkg/logging"
	"context"
	"fmt"
)

// JetBrainsPluginTask uses the jetbrains IDE command to install a plugin.
// See https://www.jetbrains.com/help/idea/install-plugins-from-the-command-line.html#macos
// Use 'id' as the plugin identifier.
type JetBrainsPluginTask struct {
	Attributes
	Apps      []string
	Utils     util.Ext
	JetBrains util.JetBrainsApp
	Log       logging.Logger
}

// NewJetBrainsPluginTask constructs a JetBrainsPluginTask.
func NewJetBrainsPluginTask(attributes *Attributes, apps []string) (*JetBrainsPluginTask, error) {
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
		Utils:      util.NewExt(),
		JetBrains:  util.NewJetBrains(),
		Log:        logging.Log(),
	}, nil
}

// GetAttributes implements Attributable.
func (t *JetBrainsPluginTask) GetAttributes() Attributes {
	return t.Attributes
}

// GetID implements Identifiable.
func (t *JetBrainsPluginTask) GetID() string {
	return t.ID
}

// IsEnabled implements Enableable.
func (t *JetBrainsPluginTask) IsEnabled() bool {
	return t.Enabled
}

// Install implements Installer.
func (t *JetBrainsPluginTask) Install(ctx context.Context, input *Input) (*TaskResult, error) {
	t.Log.Debugf("%s START ... play: %s taskID: %s", TaskEllipsis, input.Play, t.ID)
	if !t.Enabled {
		return t.CreateTaskResult(input, false)
	}

	success := true
	for _, app := range t.activeApps() {
		apps := input.PlayBook.JetBrainsApps
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

// List implements Lister.
func (t *JetBrainsPluginTask) List(_ context.Context, input *Input) (*TaskResult, error) {
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
