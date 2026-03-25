// Package tasks provides Task types and related utilities.
package tasks

import (
	ext "awong/dotfiles/pkg/external"
	l "awong/dotfiles/pkg/logging"
	ty "awong/dotfiles/pkg/types"
	"context"
	"fmt"
)

// BrewTap taps a Homebrew repository.
type BrewTap struct {
	ty.Attributes
	Brew ext.Brew
	Log  l.Logger
}

// NewBrewTap returns a new BrewTap for tapping a Homebrew repository.
func NewBrewTap(attributes *ty.Attributes) (*BrewTap, error) {
	if attributes == nil {
		return nil, fmt.Errorf("attributes cannot be nil")
	}
	if attributes.ID == "" {
		return nil, fmt.Errorf("tap ID cannot be empty")
	}
	return &BrewTap{
		Attributes: *attributes,
		Brew:       ext.NewBrew(),
		Log:        l.Log(),
	}, nil
}

// Install installs the task.
func (t *BrewTap) Install(ctx context.Context, input *ty.TaskInput) (*ty.TaskResult, error) {
	if !t.IsEnabled() {
		result, err := t.CreateTaskResult(input, false)
		if err != nil {
			return nil, fmt.Errorf("failed to create task result: %w", err)
		}
		return result, nil
	}
	t.Log.Infof("%s %s: checking tap %s", TaskEllipsis, input.Play, t.ID)
	if !input.DryRun {
		if err := t.Brew.Tap(ctx, t.ID); err != nil {
			return nil, fmt.Errorf("failed to tap %s: %v", t.ID, err)
		}
	}
	result, err := t.CreateTaskResult(input, true)
	if err != nil {
		return nil, fmt.Errorf("failed to create task result: %w", err)
	}
	return result, nil
}
