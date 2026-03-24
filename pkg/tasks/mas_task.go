package tasks

import (
	ext "awong/dotfiles/pkg/external"
	l "awong/dotfiles/pkg/logging"
	ty "awong/dotfiles/pkg/types"
	"context"
	"fmt"
	"strings"
)

type MasTask struct {
	ty.Attributes
	Mas ext.Mas
	Log l.Logger
}

func NewMasTask(attributes *ty.Attributes) (*MasTask, error) {
	if attributes == nil {
		return nil, fmt.Errorf("attributes cannot be nil")
	}
	if attributes.ID == "" {
		return nil, fmt.Errorf("task ID cannot be empty")
	}
	return &MasTask{
		Attributes: *attributes,
		Mas:        ext.NewMas(),
		Log:        l.Log(),
	}, nil
}

func (t *MasTask) GetAttributes() ty.Attributes {
	return t.Attributes
}

func (t *MasTask) GetID() string {
	return t.ID
}

func (t *MasTask) IsEnabled() bool {
	return t.Enabled
}

func (t *MasTask) Install(ctx context.Context, input *ty.TaskInput) (*ty.TaskResult, error) {
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

func (t *MasTask) List(_ context.Context, input *ty.TaskInput) (*ty.TaskResult, error) {
	err := t.Log.Printlnf("%v mas install %s ... desc: %s", TaskEllipsis, t.ID, t.Description)
	if err != nil {
		return nil, fmt.Errorf("failed to list mas: %v", err)
	}
	return t.CreateTaskResult(input, true)
}
