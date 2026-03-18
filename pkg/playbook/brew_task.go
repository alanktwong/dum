package playbook

import (
	"awong/dotfiles/pkg/external"
	"awong/dotfiles/pkg/logging"
	"context"
	"fmt"
)

// BrewTask uses homebrew to install a package.
// Use 'id' to identify the package.
type BrewTask struct {
	Attributes
	Tap   Installer
	Brew  external.Brew
	Utils external.Ext
	Log   logging.Logger
}

// NewBrewTask constructs BrewTask.
func NewBrewTask(attributes *Attributes, tap string) (*BrewTask, error) {
	if attributes == nil {
		return nil, fmt.Errorf("attributes cannot be nil")
	}
	if attributes.ID == "" {
		return nil, fmt.Errorf("attribute ID cannot be empty")
	}
	var brewTap Installer
	if tap != "" {
		aTap, err := NewBrewTap(&Attributes{
			ID:          tap,
			Description: fmt.Sprintf("brew tap %s", tap),
			Enabled:     attributes.Enabled,
			Sudo:        attributes.Sudo,
		})
		if err != nil {
			return nil, fmt.Errorf("failed to create brew tap: %v", err)
		}
		brewTap = aTap
	}
	return &BrewTask{
		Attributes: *attributes,
		Tap:        brewTap,
		Brew:       external.NewBrew(),
		Utils:      external.NewExt(),
		Log:        logging.Log(),
	}, nil
}

// GetAttributes implements Attributable.
func (t *BrewTask) GetAttributes() Attributes {
	return t.Attributes
}

// GetID implements Identifiable.
func (t *BrewTask) GetID() string {
	return t.ID
}

// IsEnabled implements Enableable.
func (t *BrewTask) IsEnabled() bool {
	return t.Enabled
}

// List implements Lister.
func (t *BrewTask) List(_ context.Context, input *Input) (*TaskResult, error) {
	err := t.Log.Printlnf("%v brew install %s", TaskEllipsis, t.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to list brew: %v", err)
	}
	return t.CreateTaskResult(input, true)
}

// Install implements Installer.
func (t *BrewTask) Install(ctx context.Context, input *Input) (*TaskResult, error) {
	t.Log.Debugf("%s START ... play: %s taskID: %s", TaskEllipsis, input.Play, t.ID)
	if !t.Enabled {
		return t.CreateTaskResult(input, false)
	}
	if t.Tap != nil {
		_, err := t.Tap.Install(ctx, input)
		if err != nil {
			return nil, fmt.Errorf("failed to tap for %s: %v", t.ID, err)
		}
	}
	if t.Utils.IsInstalled(t.ID) || t.Brew.InPath(ctx, "opt", t.ID) {
		t.Log.Infof("%s %s: %s is already installed", TaskEllipsis, input.Play, t.ID)
		return t.CreateTaskResult(input, false)
	}
	t.Log.Infof("%s %s: brew install %s", TaskEllipsis, input.Play, t.ID)
	if !input.DryRun {
		if err := t.Brew.Install(ctx, t.ID); err != nil {
			return nil, fmt.Errorf("failed to brew install %v: %v", t.ID, err)
		}
	}
	return t.CreateTaskResult(input, true)
}
