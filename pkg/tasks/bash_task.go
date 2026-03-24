package tasks

import (
	ext "awong/dotfiles/pkg/external"
	l "awong/dotfiles/pkg/logging"
	ty "awong/dotfiles/pkg/types"
	"context"
	"fmt"
)

type BashTask struct {
	ty.Attributes
	Command string
	Script  string
	Utils   ext.Ext
	Log     l.Logger
}

func NewBashTask(attributes *ty.Attributes, command, script string) (*BashTask, error) {
	if attributes == nil {
		return nil, fmt.Errorf("attributes cannot be nil")
	}
	if attributes.ID == "" {
		return nil, fmt.Errorf("task ID cannot be empty")
	}
	if command == "" && script == "" {
		return nil, fmt.Errorf("either command or script must be specified")
	}
	return &BashTask{
		Attributes: *attributes,
		Command:    command,
		Script:     script,
		Utils:      ext.NewExt(),
		Log:        l.Log(),
	}, nil
}

func (t *BashTask) GetAttributes() ty.Attributes {
	return t.Attributes
}

func (t *BashTask) GetID() string {
	return t.ID
}

func (t *BashTask) IsEnabled() bool {
	return t.Enabled
}

func (t *BashTask) Install(ctx context.Context, input *ty.TaskInput) (*ty.TaskResult, error) {
	t.Log.Debugf("%s START ... play: %s taskID: %s", TaskEllipsis, input.Play, t.ID)
	if !t.Enabled {
		return t.CreateTaskResult(input, false)
	}

	cmd := t.Command
	if cmd == "" {
		cmd = t.Script
	}

	t.Log.Infof("%s %s: %s", TaskEllipsis, input.Play, cmd)
	if !input.DryRun {
		err := t.Utils.RunCommand(ctx, cmd, t.Sudo)
		if err != nil {
			return nil, fmt.Errorf("failed to run command %s: %w", t.ID, err)
		}
	}
	return t.CreateTaskResult(input, true)
}

func (t *BashTask) List(_ context.Context, input *ty.TaskInput) (*ty.TaskResult, error) {
	cmd := t.Command
	if cmd == "" {
		cmd = t.Script
	}
	err := t.Log.Printlnf("%s %s", TaskEllipsis, cmd)
	if err != nil {
		return nil, fmt.Errorf("failed to list command: %w", err)
	}
	return t.CreateTaskResult(input, true)
}
