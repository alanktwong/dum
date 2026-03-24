package tasks

import (
	ext "awong/dotfiles/pkg/external"
	l "awong/dotfiles/pkg/logging"
	ty "awong/dotfiles/pkg/types"
	"context"
	"fmt"
)

type BrewTap struct {
	ty.Attributes
	Brew ext.Brew
	Log  l.Logger
}

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

func (t *BrewTap) Install(ctx context.Context, input *ty.TaskInput) (*ty.TaskResult, error) {
	if !t.IsEnabled() {
		return t.CreateTaskResult(input, false)
	}
	t.Log.Infof("%s %s: checking tap %s", TaskEllipsis, input.Play, t.ID)
	if !input.DryRun {
		if err := t.Brew.Tap(ctx, t.ID); err != nil {
			return nil, fmt.Errorf("failed to tap %s: %v", t.ID, err)
		}
	}
	return t.CreateTaskResult(input, true)
}
