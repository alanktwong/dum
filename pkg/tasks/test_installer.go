package tasks

import (
	ext "awong/dotfiles/pkg/external"
	l "awong/dotfiles/pkg/logging"
	ty "awong/dotfiles/pkg/types"
	"context"
)

type TestInstaller struct {
	Utils ext.Ext
	Log   l.Logger
}

func NewTestInstaller() *TestInstaller {
	return &TestInstaller{
		Log: l.Log(),
	}
}

func (t *TestInstaller) Install(_ context.Context, input *ty.TaskInput) (*ty.TaskResult, error) {
	t.Log.Debugf("%v %v: install test", TaskEllipsis, input.Play)
	return ty.NewTaskResult(input, true)
}
