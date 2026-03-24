package tasks

import (
	ext "awong/dotfiles/pkg/external"
	l "awong/dotfiles/pkg/logging"
	ty "awong/dotfiles/pkg/types"
	"context"
	"fmt"
)

type BrewCellarTask struct {
	ty.Attributes
	Tap  Installer
	Brew ext.Brew
	Log  l.Logger
}

func NewBrewCellarTask(attributes *ty.Attributes, tap string) (*BrewCellarTask, error) {
	if attributes == nil {
		return nil, fmt.Errorf("attributes cannot be nil")
	}
	if attributes.ID == "" {
		return nil, fmt.Errorf("attribute ID cannot be empty")
	}
	var brewTap Installer
	if tap != "" {
		aTap, err := NewBrewTap(&ty.Attributes{
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
		Brew:       ext.NewBrew(),
		Log:        l.Log(),
	}, nil
}

func (t *BrewCellarTask) GetAttributes() ty.Attributes {
	return t.Attributes
}

func (t *BrewCellarTask) GetID() string {
	return t.ID
}

func (t *BrewCellarTask) IsEnabled() bool {
	return t.Enabled
}

func (t *BrewCellarTask) List(_ context.Context, input *ty.TaskInput) (*ty.TaskResult, error) {
	err := t.Log.Printlnf("%v brew install (Cellar) %s", TaskEllipsis, t.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to list brew cellar: %v", err)
	}
	return t.CreateTaskResult(input, true)
}

func (t *BrewCellarTask) Install(ctx context.Context, input *ty.TaskInput) (*ty.TaskResult, error) {
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
