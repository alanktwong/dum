// Package installer provides Task types and related utilities.
package installer

import (
	ext "alanktwong/dum/internal/external"
	l "alanktwong/dum/internal/logging"
	ty "alanktwong/dum/internal/types"
	"context"
	"fmt"
	"os/exec"
)

// SdkmanInstaller installs the SDKMAN! SDK management tool.
type SdkmanInstaller struct {
	Utils ext.Ext
	Curl  ext.Runner
	Log   l.Logger
}

// NewSdkmanInstaller returns a new SdkmanInstaller for installing SDKMAN!.
func NewSdkmanInstaller() *SdkmanInstaller {
	return &SdkmanInstaller{
		Utils: ext.NewExt(),
		Curl:  &SdkCurlRunner{},
		Log:   l.Log(),
	}
}

// SdkCurlRunner runs the SDKMAN! install script.
type SdkCurlRunner struct{}

// Run runs the SDKMAN! install script.
func (r *SdkCurlRunner) Run(ctx context.Context) error {
	var cmd *exec.Cmd
	sudo := GetSudo(ctx)
	if sudo {
		cmd = exec.CommandContext(ctx, "bash", "-c", "$(curl -s https://get.sdkman.io)")
	} else {
		cmd = exec.CommandContext(ctx, "sudo", "bash", "-c", "$(curl -s https://get.sdkman.io)")
	}
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to run sdkman install: %w", err)
	}
	return nil
}

// Install installs the task.
func (t *SdkmanInstaller) Install(ctx context.Context, input *ty.TaskInput) (*ty.TaskResult, error) {
	if t.Utils.IsInstalled("sdk") {
		t.Log.Infof("%v %v: sdkman already installed", TaskEllipsis, input.Play)
		result, err := ty.NewTaskResult(input, false)
		if err != nil {
			return nil, fmt.Errorf("failed to create task result: %w", err)
		}
		return result, nil
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
	result, err := ty.NewTaskResult(input, true)
	if err != nil {
		return nil, fmt.Errorf("failed to create task result: %w", err)
	}
	return result, nil
}
