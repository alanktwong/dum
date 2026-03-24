package tasks

import (
	ext "awong/dotfiles/pkg/external"
	l "awong/dotfiles/pkg/logging"
	ty "awong/dotfiles/pkg/types"
	"context"
	"fmt"
	"path/filepath"
)

type LinkTask struct {
	ty.Attributes
	Root   string
	Target string
	Utils  ext.Ext
	Log    l.Logger
}

func NewLinkTask(attributes *ty.Attributes, root, target string) (*LinkTask, error) {
	if attributes == nil {
		return nil, fmt.Errorf("attributes cannot be nil")
	}
	if attributes.ID == "" {
		return nil, fmt.Errorf("task ID cannot be empty")
	}
	return &LinkTask{
		Attributes: *attributes,
		Root:       root,
		Target:     target,
		Utils:      ext.NewExt(),
		Log:        l.Log(),
	}, nil
}

func (t *LinkTask) GetAttributes() ty.Attributes {
	return t.Attributes
}

func (t *LinkTask) GetID() string {
	return t.ID
}

func (t *LinkTask) IsEnabled() bool {
	return t.Enabled
}

func (t *LinkTask) Install(ctx context.Context, input *ty.TaskInput) (*ty.TaskResult, error) {
	t.Log.Debugf("%s START ... play: %s taskID: %s", TaskEllipsis, input.Play, t.ID)
	if t.Root == "" {
		return t.CreateTaskResult(input, false)
	}
	target := t.ProvideTarget()
	if target == "" {
		return t.CreateTaskResult(input, false)
	}
	if !t.Enabled {
		return t.CreateTaskResult(input, false)
	}
	rootPath, err := t.Utils.ExpandUser(t.Root)
	if err != nil {
		return nil, fmt.Errorf("failed to expand root path %s: %w", t.Root, err)
	}
	targetPath := filepath.Join(rootPath, target)
	if t.Utils.IsSymlink(targetPath) {
		t.Log.Infof("%v %v: %v -> %s/%s already linked", TaskEllipsis, input.Play, t.ID, t.Root, target)
	}
	t.Log.Infof("%v %v: planning to link %v -> %s/%s", TaskEllipsis, input.Play, t.ID, t.Root, target)
	if !input.DryRun {
		if t.Utils.SoftLink(ctx, rootPath, t.ID, target, t.Sudo) != nil {
			return nil, fmt.Errorf("failed to link %v -> %s/%s: %v", t.ID, t.Root, target, err)
		}
	}
	return t.CreateTaskResult(input, true)
}

func (t *LinkTask) List(_ context.Context, input *ty.TaskInput) (*ty.TaskResult, error) {
	target := t.ProvideTarget()
	err := t.Log.Printlnf("%v linking %v -> %s/%s", TaskEllipsis, t.ID, t.Root, target)
	if err != nil {
		return nil, fmt.Errorf("failed to list link: %v", err)
	}
	return t.CreateTaskResult(input, true)
}

func (t *LinkTask) ProvideTarget() string {
	if t.Target != "" {
		return t.Target
	}
	return filepath.Base(t.ID)
}
