package playbook

import (
	"awong/dotfiles/pkg/external"
	"awong/dotfiles/pkg/logging"
	"context"
	"fmt"
)

// BrewCellarTask uses homebrew to install a Cellar package.
// Use 'id' to identify the package.
type BrewCellarTask struct {
	Attributes
	Tap  Installer
	Brew external.Brew
	Log  logging.Logger
}

// NewBrewCellarTask constructs BrewCellarTask.
func NewBrewCellarTask(attributes *Attributes, tap string) (*BrewCellarTask, error) {
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
	return &BrewCellarTask{
		Attributes: *attributes,
		Tap:        brewTap,
		Brew:       external.NewBrew(),
		Log:        logging.Log(),
	}, nil
}

// GetAttributes implements Attributable.
func (t *BrewCellarTask) GetAttributes() Attributes {
	return t.Attributes
}

// GetID implements Identifiable.
func (t *BrewCellarTask) GetID() string {
	return t.ID
}

// IsEnabled implements Enableable.
func (t *BrewCellarTask) IsEnabled() bool {
	return t.Enabled
}

// List implements Lister.
func (t *BrewCellarTask) List(_ context.Context, input *Input) (*TaskResult, error) {
	err := t.Log.Printlnf("%v brew install (Cellar) %s", TaskEllipsis, t.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to list brew cellar: %v", err)
	}
	return t.CreateTaskResult(input, true)
}

// Install implements Installer.
func (t *BrewCellarTask) Install(ctx context.Context, input *Input) (*TaskResult, error) {
	t.Log.Infof("%s START ... play: %s taskID: %s", TaskEllipsis, input.Play, t.ID)
	if !t.Enabled {
		return t.CreateTaskResult(input, false)
	}
	if t.Tap != nil {
		_, err := t.Tap.Install(ctx, input)
		if err != nil {
			return nil, fmt.Errorf("failed to tap for %s: %v", t.ID, err)
		}
	}
	if t.Brew.InPath(ctx, "Cellar", t.ID) {
		t.Log.Debugf("%s %s: %s is already installed", TaskEllipsis, input.Play, t.ID)
		return t.CreateTaskResult(input, false)
	}
	t.Log.Infof("%s %s: brew install (cellar) %s", TaskEllipsis, input.Play, t.ID)
	if !input.DryRun {
		if err := t.Brew.Install(ctx, t.ID); err != nil {
			return nil, fmt.Errorf("failed to brew install (cellar) %v: %v", t.ID, err)
		}
	}
	return t.CreateTaskResult(input, true)
}
