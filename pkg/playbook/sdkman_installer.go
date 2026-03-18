package playbook

import (
	"awong/dotfiles/pkg/external"
	"awong/dotfiles/pkg/logging"
	"context"
	"fmt"
	"os/exec"
)

// SdkmanInstaller installs sdkman.
type SdkmanInstaller struct {
	Utils external.Ext
	Curl  external.Runner
	Log   logging.Logger
}

// NewSdkmanInstaller constructs a SdkmanInstaller.
func NewSdkmanInstaller() *SdkmanInstaller {
	return &SdkmanInstaller{
		Utils: external.NewExt(),
		Curl:  &SdkCurlRunner{},
		Log:   logging.Log(),
	}
}

// SdkCurlRunner implements the Runner interface.
type SdkCurlRunner struct{}

// Run implements the Runner interface.
func (r *SdkCurlRunner) Run(ctx context.Context) error {
	var cmd *exec.Cmd
	sudo := GetSudo(ctx)
	if sudo {
		cmd = exec.CommandContext(ctx, "bash", "-c", "$(curl -s https://get.sdkman.io)")
	} else {
		cmd = exec.CommandContext(ctx, "sudo", "bash", "-c", "$(curl -s https://get.sdkman.io)")
	}
	return cmd.Run() //nolint:wrapcheck
}

// Install implements Installer.
func (t *SdkmanInstaller) Install(ctx context.Context, input *Input) (*TaskResult, error) {
	if t.Utils.IsInstalled("sdk") {
		t.Log.Infof("%v %v: sdkman already installed", TaskEllipsis, input.Play)
		return NewTaskResult(input, false)
	}
	t.Log.Infof("%v %v: install sdkman", TaskEllipsis, input.Play)
	if !input.DryRun {
		if t.Utils.IsOSX() || t.Utils.IsLinux() {
			if err := t.Curl.Run(WithSudo(ctx, input.Sudo)); err != nil {
				return nil, fmt.Errorf("failed to run sdkman install: %w", err)
			}
		} else {
			return nil, fmt.Errorf("fail to install sdkman since this is neither linux nor OSX")
		}
	}
	return NewTaskResult(input, true)
}
