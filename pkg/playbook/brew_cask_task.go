package playbook

import (
	"awong/dotfiles/pkg/external"
	"awong/dotfiles/pkg/logging"
	"context"
	"fmt"
)

// BrewCaskTask uses homebrew to install a Cask package.
// Use 'id' to identify the package.
type BrewCaskTask struct {
	Attributes
	Tap   Installer
	Brew  external.Brew
	Utils external.Ext
	Log   logging.Logger
}

// NewBrewCaskTask constructs BrewCaskTask.
func NewBrewCaskTask(attributes *Attributes, tap string) (*BrewCaskTask, error) {
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
	return &BrewCaskTask{
		Attributes: *attributes,
		Tap:        brewTap,
		Brew:       external.NewBrew(),
		Utils:      external.NewExt(),
		Log:        logging.Log(),
	}, nil
}

// GetAttributes implements Attributable.
func (t *BrewCaskTask) GetAttributes() Attributes {
	return t.Attributes
}

// GetID implements Identifiable.
func (t *BrewCaskTask) GetID() string {
	return t.ID
}

// IsEnabled implements Enableable.
func (t *BrewCaskTask) IsEnabled() bool {
	return t.Enabled
}

// List implements Lister.
func (t *BrewCaskTask) List(_ context.Context, input *Input) (*TaskResult, error) {
	err := t.Log.Printlnf("%v brew install (Cask) %s", TaskEllipsis, t.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to list brew cask: %v", err)
	}
	return t.CreateTaskResult(input, true)
}

// Install implements Installer.
func (t *BrewCaskTask) Install(ctx context.Context, input *Input) (*TaskResult, error) {
	t.Log.Debugf("%s START ... play: %s taskID: %s", TaskEllipsis, input.Play, t.ID)
	if !t.Enabled {
		return t.CreateTaskResult(input, false)
	}
	if !t.Utils.IsOSX() {
		return nil, fmt.Errorf("cannot install cask %s outside of macOS", t.ID)
	}
	if t.Tap != nil {
		_, err := t.Tap.Install(ctx, input)
		if err != nil {
			return nil, fmt.Errorf("failed to tap for %s: %v", t.ID, err)
		}
	}
	if t.Brew.InPath(ctx, "Caskroom", t.ID) {
		t.Log.Infof("%s %s: %s is already installed", TaskEllipsis, input.Play, t.ID)
		return t.CreateTaskResult(input, false)
	}
	t.Log.Infof("%s %s: brew install (Cask) %s", TaskEllipsis, input.Play, t.ID)
	if !input.DryRun {
		if err := t.Brew.InstallCask(ctx, t.ID); err != nil {
			return nil, fmt.Errorf("failed to brew install --cask %v: %v", t.ID, err)
		}
	}
	return t.CreateTaskResult(input, true)
}
