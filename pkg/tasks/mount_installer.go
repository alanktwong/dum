package tasks

import (
	ext "awong/dotfiles/pkg/external"
	l "awong/dotfiles/pkg/logging"
	ty "awong/dotfiles/pkg/types"
	"context"
	"fmt"
	"os/exec"
)

type MountInstaller struct {
	Utils  ext.Ext
	Runner ext.Runner
	Log    l.Logger
}

func NewMountInstaller() *MountInstaller {
	return &MountInstaller{
		Utils:  ext.NewExt(),
		Runner: &MountRunner{},
		Log:    l.Log(),
	}
}

type MountRunner struct{}

func (t *MountRunner) Run(ctx context.Context) error {
	cmd := exec.CommandContext(ctx, "/bin/bash", "-c", "~/projects/dotfiles/bin/make_case_sensitive_fs.sh")
	return cmd.Run()
}

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
	return ty.NewTaskResult(input, true)
}
