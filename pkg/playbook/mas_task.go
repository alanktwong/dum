package playbook

import (
	"awong/dotfiles/pkg/external"
	"awong/dotfiles/pkg/logging"
	"context"
	"fmt"
	"strings"
)

// MasTask uses the mas command to install an app from the Mac App Store.
// Use 'id' as the app store id.
type MasTask struct {
	Attributes
	Mas external.Mas
	Log logging.Logger
}

// NewMasTask constructs a MasTask.
func NewMasTask(attributes *Attributes) (*MasTask, error) {
	if attributes == nil {
		return nil, fmt.Errorf("attributes cannot be nil")
	}
	if attributes.ID == "" {
		return nil, fmt.Errorf("task ID cannot be empty")
	}
	return &MasTask{
		Attributes: *attributes,
		Mas:        external.NewMas(),
		Log:        logging.Log(),
	}, nil
}

// GetAttributes implements Attributable.
func (t *MasTask) GetAttributes() Attributes {
	return t.Attributes
}

// GetID implements Identifiable.
func (t *MasTask) GetID() string {
	return t.ID
}

// IsEnabled implements Enableable.
func (t *MasTask) IsEnabled() bool {
	return t.Enabled
}

// Install implements Installer.
func (t *MasTask) Install(ctx context.Context, input *Input) (*TaskResult, error) {
	t.Log.Debugf("%s START ... play: %s taskID: %s", TaskEllipsis, input.Play, t.ID)
	if t.Description == "" {
		return t.CreateTaskResult(input, false)
	}
	if !t.Enabled {
		return t.CreateTaskResult(input, false)
	}
	masList, err := t.Mas.List(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to list mas apps: %w", err)
	}
	count := strings.Count(masList, t.ID)
	if count > 0 {
		t.Log.Infof("%s %s: %s is already installed", TaskEllipsis, input.Play, t.ID)
		return t.CreateTaskResult(input, false)
	}
	t.Log.Infof("%s %s: mas install %s", TaskEllipsis, input.Play, t.ID)
	if !input.DryRun {
		if err := t.Mas.Install(ctx, t.ID); err != nil {
			return nil, fmt.Errorf("%s %s: mas install %s failed: %w", TaskEllipsis, input.Play, t.ID, err)
		}
	}
	return t.CreateTaskResult(input, true)
}

// List implements Lister.
func (t *MasTask) List(_ context.Context, input *Input) (*TaskResult, error) {
	err := t.Log.Printlnf("%v mas install %s ... desc: %s", TaskEllipsis, t.ID, t.Description)
	if err != nil {
		return nil, fmt.Errorf("failed to list mas: %v", err)
	}
	return t.CreateTaskResult(input, true)
}
