package tasks

import (
	ext "awong/dotfiles/pkg/external"
	l "awong/dotfiles/pkg/logging"
	ty "awong/dotfiles/pkg/types"
	"context"
	"fmt"
	"os/exec"
)

type SdkmanInstaller struct {
	Utils ext.Ext
	Curl  ext.Runner
	Log   l.Logger
}

func NewSdkmanInstaller() *SdkmanInstaller {
	return &SdkmanInstaller{
		Utils: ext.NewExt(),
		Curl:  &SdkCurlRunner{},
		Log:   l.Log(),
	}
}

type SdkCurlRunner struct{}

func (r *SdkCurlRunner) Run(ctx context.Context) error {
	var cmd *exec.Cmd
	sudo := ty.GetSudo(ctx)
	if sudo {
		cmd = exec.CommandContext(ctx, "bash", "-c", "$(curl -s https://get.sdkman.io)")
	} else {
		cmd = exec.CommandContext(ctx, "sudo", "bash", "-c", "$(curl -s https://get.sdkman.io)")
	}
	return cmd.Run()
}

func (t *SdkmanInstaller) Install(ctx context.Context, input *ty.TaskInput) (*ty.TaskResult, error) {
	if t.Utils.IsInstalled("sdk") {
		t.Log.Infof("%v %v: sdkman already installed", TaskEllipsis, input.Play)
		return ty.NewTaskResult(input, false)
	}
	t.Log.Infof("%v %v: install sdkman", TaskEllipsis, input.Play)
	if !input.DryRun {
		if t.Utils.IsOSX() || t.Utils.IsLinux() {
			if err := t.Curl.Run(ty.WithSudo(ctx, input.Sudo)); err != nil {
				return nil, fmt.Errorf("failed to run sdkman install: %w", err)
			}
		} else {
			return nil, fmt.Errorf("fail to install sdkman since this is neither linux nor OSX")
		}
	}
	return ty.NewTaskResult(input, true)
}
