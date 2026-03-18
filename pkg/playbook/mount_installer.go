package playbook

import (
	"awong/dotfiles/pkg/external"
	"awong/dotfiles/pkg/logging"
	"context"
	"fmt"
	"os/exec"
)

// MountInstaller creates a case-sensitive file system for MacOSX development.
type MountInstaller struct {
	Utils  external.Ext
	Runner external.Runner
	Log    logging.Logger
}

// NewMountInstaller constructs a MountInstaller.
func NewMountInstaller() *MountInstaller {
	return &MountInstaller{
		Utils:  external.NewExt(),
		Runner: &MountRunner{},
		Log:    logging.Log(),
	}
}

// MountRunner implements the Runner interface.
type MountRunner struct{}

// Run implements the Runner interface.
func (t *MountRunner) Run(ctx context.Context) error {
	cmd := exec.CommandContext(ctx, "/bin/bash", "-c", "~/projects/dotfiles/bin/make_case_sensitive_fs.sh")
	return cmd.Run() //nolint:wrapcheck
}

// Install implements Installer.
func (t *MountInstaller) Install(ctx context.Context, input *Input) (*TaskResult, error) {
	if !t.Utils.IsOSX() {
		return nil, fmt.Errorf("cannot install a case-sensitive file system; %v is not on a Mac OSX", input.Play)
	}
	t.Log.Infof("%v %v: mounting a case-sensitive file system", TaskEllipsis, input.Play)
	if !input.DryRun {
		if err := t.Runner.Run(ctx); err != nil {
			return nil, fmt.Errorf("fail to install a case-sensitive file system: %v", err)
		}
	}
	return NewTaskResult(input, true)
}
