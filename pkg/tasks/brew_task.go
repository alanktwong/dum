package tasks

import (
	ext "awong/dotfiles/pkg/external"
	l "awong/dotfiles/pkg/logging"
	ty "awong/dotfiles/pkg/types"
	"context"
	"fmt"
)

type BrewTask struct {
	ty.Attributes
	Tap   Installer
	Brew  ext.Brew
	Utils ext.Ext
	Log   l.Logger
}

func NewBrewTask(attributes *ty.Attributes, tap string) (*BrewTask, error) {
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
	return &BrewTask{
		Attributes: *attributes,
		Tap:        brewTap,
		Brew:       ext.NewBrew(),
		Utils:      ext.NewExt(),
		Log:        l.Log(),
	}, nil
}

func (t *BrewTask) GetAttributes() ty.Attributes {
	return t.Attributes
}

func (t *BrewTask) GetID() string {
	return t.ID
}

func (t *BrewTask) IsEnabled() bool {
	return t.Enabled
}

func (t *BrewTask) List(_ context.Context, input *ty.TaskInput) (*ty.TaskResult, error) {
	err := t.Log.Printlnf("%v brew install %s", TaskEllipsis, t.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to list brew: %v", err)
	}
	return t.CreateTaskResult(input, true)
}

func (t *BrewTask) Install(ctx context.Context, input *ty.TaskInput) (*ty.TaskResult, error) {
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
