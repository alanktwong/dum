package playbook

import (
	"awong/dotfiles/pkg/external"
	"awong/dotfiles/pkg/logging"
	"context"
)

// TestInstaller is a no-ope installer.
type TestInstaller struct {
	Utils external.Ext
	Log   logging.Logger
}

// NewTestInstaller constructs a TestInstaller.
func NewTestInstaller() *TestInstaller {
	return &TestInstaller{
		Log: logging.Log(),
	}
}

// Install implements Installer.
func (t *TestInstaller) Install(_ context.Context, input *Input) (*TaskResult, error) {
	t.Log.Debugf("%v %v: install test", TaskEllipsis, input.Play)
	return NewTaskResult(input, true)
}
