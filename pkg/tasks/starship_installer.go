package tasks

import (
	ext "awong/dotfiles/pkg/external"
	l "awong/dotfiles/pkg/logging"
	ty "awong/dotfiles/pkg/types"
	"context"
	"fmt"
	"os/exec"
)

type StarshipInstaller struct {
	Utils     ext.Ext
	Curl      ext.Runner
	LinuxCurl ext.Runner
	Log       l.Logger
}

func NewStarshipInstaller() *StarshipInstaller {
	return &StarshipInstaller{
		Utils:     ext.NewExt(),
		Curl:      &StarshipCurlRunner{},
		LinuxCurl: &StarshipLinuxCurlRunner{},
		Log:       l.Log(),
	}
}

type StarshipCurlRunner struct{}

func (t *StarshipCurlRunner) Run(ctx context.Context) error {
	sudo := ty.GetSudo(ctx)
	var cmd *exec.Cmd
	if sudo {
		cmd = exec.CommandContext(ctx, "sudo", "sh", "-c", "$(curl -sS https://starship.rs/install.sh)")
	} else {
		cmd = exec.CommandContext(ctx, "sh", "-c", "$(curl -sS https://starship.rs/install.sh)")
	}
	return cmd.Run()
}

type StarshipLinuxCurlRunner struct{}

func (t *StarshipLinuxCurlRunner) Run(ctx context.Context) error {
	sudo := ty.GetSudo(ctx)
	var cmd *exec.Cmd
	if sudo {
		cmd = exec.CommandContext(ctx, "sudo", "bash", "-c", "$(curl -sS https://starship.rs/install.sh)")
	} else {
		cmd = exec.CommandContext(ctx, "bash", "-c", "$(curl -sS https://starship.rs/install.sh)")
	}

	return cmd.Run()
}

func (t *StarshipInstaller) Install(ctx context.Context, input *ty.TaskInput) (*ty.TaskResult, error) {
	if t.Utils.IsInstalled("starship") {
		t.Log.Infof("%v %v: starship already installed", TaskEllipsis, input.Play)
		return ty.NewTaskResult(input, false)
	}
	t.Log.Infof("%v %v: install starship", TaskEllipsis, input.Play)
	if !input.DryRun {
		if err := t.Run(ty.WithSudo(ctx, input.Sudo)); err != nil {
			return nil, err
		}
	}
	return ty.NewTaskResult(input, true)
}

func (t *StarshipInstaller) Run(ctx context.Context) error {
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
