package tasks

import (
	ext "awong/dotfiles/pkg/external"
	l "awong/dotfiles/pkg/logging"
	ty "awong/dotfiles/pkg/types"
	"context"
	"fmt"
)

type DirTask struct {
	ty.Attributes
	Utils ext.Ext
	Log   l.Logger
}

func NewDirTask(attributes *ty.Attributes) (*DirTask, error) {
	if attributes == nil {
		return nil, fmt.Errorf("attributes cannot be nil")
	}
	if attributes.ID == "" {
		return nil, fmt.Errorf("task ID cannot be empty")
	}
	return &DirTask{
		Attributes: *attributes,
		Utils:      ext.NewExt(),
		Log:        l.Log(),
	}, nil
}

func (t *DirTask) GetAttributes() ty.Attributes {
	return t.Attributes
}

func (t *DirTask) GetID() string {
	return t.ID
}

func (t *DirTask) IsEnabled() bool {
	return t.Enabled
}

func (t *DirTask) Install(ctx context.Context, input *ty.TaskInput) (*ty.TaskResult, error) {
	t.Log.Debugf("%s START ... play: %s taskID: %s", TaskEllipsis, input.Play, t.ID)
	if !t.Enabled {
		return t.CreateTaskResult(input, false)
	}
	path, err := t.Utils.ExpandUser(t.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to expand user for %s: %v", t.ID, err)
	}
	if t.Utils.IsDir(path) {
		t.Log.Infof("%s %s: %s already exists", TaskEllipsis, input.Play, t.ID)
		return t.CreateTaskResult(input, false)
	}
	t.Log.Infof("%s %s: mkdir -p %s", TaskEllipsis, input.Play, t.ID)
	if !input.DryRun {
		err = t.Utils.CreateDirectory(ctx, path, t.Sudo)
		if err != nil {
			return nil, fmt.Errorf("failed to mkdir -p %v: %v", t.ID, err)
		}
	}
	return t.CreateTaskResult(input, true)
}

func (t *DirTask) List(_ context.Context, input *ty.TaskInput) (*ty.TaskResult, error) {
	err := t.Log.Printlnf("%v mkdir -p %s", TaskEllipsis, t.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to list dir: %v", err)
	}
	return t.CreateTaskResult(input, true)
}
