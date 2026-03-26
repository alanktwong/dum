// Package tasks provides Task types and related utilities.
package tasks

import (
	ext "awong/dotfiles/internal/external"
	l "awong/dotfiles/internal/logging"
	i "awong/dotfiles/internal/tasks/installer"
	ty "awong/dotfiles/internal/types"
	"context"
	"fmt"
)

// BrewCaskTask installs a Homebrew Cask package.
type BrewCaskTask struct {
	ty.Attributes
	Tap   i.Installer
	Brew  ext.Brew
	Utils ext.Ext
	Log   l.Logger
}

// NewBrewCaskTask returns a new BrewCaskTask for installing a Homebrew Cask package.
func NewBrewCaskTask(attributes *ty.Attributes, tap string) (*BrewCaskTask, error) {
	if attributes == nil {
		return nil, fmt.Errorf("attributes cannot be nil")
	}
	if attributes.ID == "" {
		return nil, fmt.Errorf("attribute ID cannot be empty")
	}
	var brewTap i.Installer
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
	return &BrewCaskTask{
		Attributes: *attributes,
		Tap:        brewTap,
		Brew:       ext.NewBrew(),
		Utils:      ext.NewExt(),
		Log:        l.Log(),
	}, nil
}

// GetAttributes returns the Attributes.
func (t *BrewCaskTask) GetAttributes() ty.Attributes {
	return t.Attributes
}

// GetID returns the ID.
func (t *BrewCaskTask) GetID() string {
	return t.ID
}

// IsEnabled returns whether the task is enabled.
func (t *BrewCaskTask) IsEnabled() bool {
	return t.Enabled
}

// List lists the task.
func (t *BrewCaskTask) List(_ context.Context, input *ty.TaskInput) (*ty.TaskResult, error) {
	err := t.Log.Printlnf("%v brew install (Cask) %s", TaskEllipsis, t.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to list brew cask: %v", err)
	}
	result, err := t.CreateTaskResult(input, true)
	if err != nil {
		return nil, fmt.Errorf("failed to create task result: %w", err)
	}
	return result, nil
}

// Install installs the task.
func (t *BrewCaskTask) Install(ctx context.Context, input *ty.TaskInput) (*ty.TaskResult, error) {
	t.Log.Debugf("%s START ... play: %s taskID: %s", TaskEllipsis, input.Play, t.ID)
	if !t.Enabled {
		result, err := t.CreateTaskResult(input, false)
		if err != nil {
			return nil, fmt.Errorf("failed to create task result: %w", err)
		}
		return result, nil
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
		result, err := t.CreateTaskResult(input, false)
		if err != nil {
			return nil, fmt.Errorf("failed to create task result: %w", err)
		}
		return result, nil
	}
	t.Log.Infof("%s %s: brew install (Cask) %s", TaskEllipsis, input.Play, t.ID)
	if !input.DryRun {
		if err := t.Brew.InstallCask(ctx, t.ID); err != nil {
			return nil, fmt.Errorf("failed to brew install --cask %v: %v", t.ID, err)
		}
	}
	result, err := t.CreateTaskResult(input, true)
	if err != nil {
		return nil, fmt.Errorf("failed to create task result: %w", err)
	}
	return result, nil
}
