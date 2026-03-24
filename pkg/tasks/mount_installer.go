// Package tasks provides Task types and related utilities.
package tasks

import (
	ext "awong/dotfiles/pkg/external"
	l "awong/dotfiles/pkg/logging"
	ty "awong/dotfiles/pkg/types"
	"context"
	"fmt"
	"os/exec"
)

// MountInstaller installs a case-sensitive file system mount on macOS.
type MountInstaller struct {
	Utils  ext.Ext
	Runner ext.Runner
	Log    l.Logger
}

// NewMountInstaller returns a new MountInstaller for installing a case-sensitive file system.
func NewMountInstaller() *MountInstaller {
	return &MountInstaller{
		Utils:  ext.NewExt(),
		Runner: &MountRunner{},
		Log:    l.Log(),
	}
}

// MountRunner runs the mount script.
type MountRunner struct{}

// Run runs the mount script.
func (t *MountRunner) Run(ctx context.Context) error {
	cmd := exec.CommandContext(ctx, "/bin/bash", "-c", "~/projects/dotfiles/bin/make_case_sensitive_fs.sh")
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to run mount script: %w", err)
	}
	return nil
}

// Install installs the task.
func (t *MountInstaller) Install(ctx context.Context, input *ty.TaskInput) (*ty.TaskResult, error) {
	if !t.Utils.IsOSX() {
		return nil, fmt.Errorf("cannot install a case-sensitive file system; %v is not on a Mac OSX", input.Play)
	}
	t.Log.Infof("%v %v: mounting a case-sensitive file system", TaskEllipsis, input.Play)
	if !input.DryRun {
		if err := t.Runner.Run(ctx); err != nil {
			return nil, fmt.Errorf("fail to install a case-sensitive file system: %v", err)
		}
	}
	result, err := ty.NewTaskResult(input, true)
	if err != nil {
		return nil, fmt.Errorf("failed to create task result: %w", err)
	}
	return result, nil
}
