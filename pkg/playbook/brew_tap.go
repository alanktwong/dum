package playbook

import (
	"awong/dotfiles/pkg/external"
	"awong/dotfiles/pkg/logging"
	"context"
	"fmt"
)

// BrewTap installs a brew tap.
type BrewTap struct {
	Attributes
	Brew external.Brew
	Log  logging.Logger
}

// NewBrewTap constructs a BrewTap.
func NewBrewTap(attributes *Attributes) (*BrewTap, error) {
	if attributes == nil {
		return nil, fmt.Errorf("attributes cannot be nil")
	}
	if attributes.ID == "" {
		return nil, fmt.Errorf("tap ID cannot be empty")
	}
	return &BrewTap{
		Attributes: *attributes,
		Brew:       external.NewBrew(),
		Log:        logging.Log(),
	}, nil
}

// Install implements Installer.
func (t *BrewTap) Install(ctx context.Context, input *Input) (*TaskResult, error) {
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
