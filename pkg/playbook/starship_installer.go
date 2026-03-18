package playbook

import (
	"awong/dotfiles/pkg/external"
	"awong/dotfiles/pkg/logging"
	"context"
	"fmt"
	"os/exec"
)

// StarshipInstaller installs starship, a cross-platform command line prompt.
type StarshipInstaller struct {
	Utils     external.Ext
	Curl      external.Runner
	LinuxCurl external.Runner
	Log       logging.Logger
}

// NewStarshipInstaller constructs a StarshipInstaller.
func NewStarshipInstaller() *StarshipInstaller {
	return &StarshipInstaller{
		Utils:     external.NewExt(),
		Curl:      &StarshipCurlRunner{},
		LinuxCurl: &StarshipLinuxCurlRunner{},
		Log:       logging.Log(),
	}
}

// StarshipCurlRunner implements Runner.
type StarshipCurlRunner struct{}

// Run implements Runner.
func (t *StarshipCurlRunner) Run(ctx context.Context) error {
	sudo := GetSudo(ctx)
	var cmd *exec.Cmd
	if sudo {
		cmd = exec.CommandContext(ctx, "sudo", "sh", "-c", "$(curl -sS https://starship.rs/install.sh)")
	} else {
		cmd = exec.CommandContext(ctx, "sh", "-c", "$(curl -sS https://starship.rs/install.sh)")
	}
	return cmd.Run() //nolint:wrapcheck
}

// StarshipLinuxCurlRunner implements Runner.
type StarshipLinuxCurlRunner struct{}

// Run implements Runner.
func (t *StarshipLinuxCurlRunner) Run(ctx context.Context) error {
	sudo := GetSudo(ctx)
	var cmd *exec.Cmd
	if sudo {
		cmd = exec.CommandContext(ctx, "sudo", "bash", "-c", "$(curl -sS https://starship.rs/install.sh)")
	} else {
		cmd = exec.CommandContext(ctx, "bash", "-c", "$(curl -sS https://starship.rs/install.sh)")
	}

	return cmd.Run() //nolint:wrapcheck
}

// Install implements Installer.
func (t *StarshipInstaller) Install(ctx context.Context, input *Input) (*TaskResult, error) {
	if t.Utils.IsInstalled("starship") {
		t.Log.Infof("%v %v: starship already installed", TaskEllipsis, input.Play)
		return NewTaskResult(input, false)
	}
	t.Log.Infof("%v %v: install starship", TaskEllipsis, input.Play)
	if !input.DryRun {
		if err := t.run(WithSudo(ctx, input.Sudo)); err != nil {
			return nil, err
		}
	}
	return NewTaskResult(input, true)
}

func (t *StarshipInstaller) run(ctx context.Context) error {
	if t.Utils.IsOSX() {
		if err := t.Curl.Run(ctx); err != nil {
			return fmt.Errorf("failed to install starship for OSX: %w", err)
		}
	} else if t.Utils.IsLinux() {
		if err := t.LinuxCurl.Run(ctx); err != nil {
			return fmt.Errorf("failed to install starship for linux: %w", err)
		}
	} else {
		return fmt.Errorf("fail to install starship since this is neither linux nor OSX")
	}
	return nil
}
