// Package installer provides Task types and related utilities.
package installer

import (
	ext "awong/dotfiles/internal/external"
	l "awong/dotfiles/internal/logging"
	ty "awong/dotfiles/internal/types"
	"context"
	"fmt"
)

// TestInstaller is a test installer that always succeeds.
type TestInstaller struct {
	Utils ext.Ext
	Log   l.Logger
}

// NewTestInstaller returns a new TestInstaller for testing purposes.
func NewTestInstaller() *TestInstaller {
	return &TestInstaller{
		Log: l.Log(),
	}
}

// Install installs the task.
func (t *TestInstaller) Install(_ context.Context, input *ty.TaskInput) (*ty.TaskResult, error) {
	t.Log.Debugf("%v %v: install test", TaskEllipsis, input.Play)
	result, err := ty.NewTaskResult(input, true)
	if err != nil {
		return nil, fmt.Errorf("failed to create task result: %w", err)
	}
	return result, nil
}
