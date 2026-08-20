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

// StarshipInstaller installs the Starship prompt.
type StarshipInstaller struct {
	Utils     ext.Ext
	Curl      ext.Runner
	LinuxCurl ext.Runner
	Log       l.Logger
}

// NewStarshipInstaller returns a new StarshipInstaller for installing the Starship prompt.
func NewStarshipInstaller() *StarshipInstaller {
	return &StarshipInstaller{
		Utils:     ext.NewExt(),
		Curl:      &StarshipCurlRunner{},
		LinuxCurl: &StarshipLinuxCurlRunner{},
		Log:       l.Log(),
	}
}

// StarshipCurlRunner runs the Starship install script on macOS.
type StarshipCurlRunner struct{}

// Run runs the Starship install script on macOS.
func (t *StarshipCurlRunner) Run(ctx context.Context) error {
	sudo := GetSudo(ctx)
	var cmd *exec.Cmd
	if sudo {
		cmd = exec.CommandContext(ctx, "sudo", "sh", "-c", "$(curl -sS https://starship.rs/install.sh)")
	} else {
		cmd = exec.CommandContext(ctx, "sh", "-c", "$(curl -sS https://starship.rs/install.sh)")
	}
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to run starship install: %w", err)
	}
	return nil
}

// StarshipLinuxCurlRunner runs the Starship install script on Linux.
type StarshipLinuxCurlRunner struct{}

// Run runs the Starship install script on Linux.
func (t *StarshipLinuxCurlRunner) Run(ctx context.Context) error {
	sudo := GetSudo(ctx)
	var cmd *exec.Cmd
	if sudo {
		cmd = exec.CommandContext(ctx, "sudo", "bash", "-c", "$(curl -sS https://starship.rs/install.sh)")
	} else {
		cmd = exec.CommandContext(ctx, "bash", "-c", "$(curl -sS https://starship.rs/install.sh)")
	}

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to run starship install: %w", err)
	}
	return nil
}

// Install installs the task.
func (t *StarshipInstaller) Install(ctx context.Context, input *ty.TaskInput) (*ty.TaskResult, error) {
	if t.Utils.IsInstalled("starship") {
		t.Log.Infof("%v %v: starship already installed", TaskEllipsis, input.Play)
		result, err := ty.NewTaskResult(input, false)
		if err != nil {
			return nil, fmt.Errorf("failed to create task result: %w", err)
		}
		return result, nil
	}
	t.Log.Infof("%v %v: install starship", TaskEllipsis, input.Play)
	if !input.DryRun {
		if err := t.Run(WithSudo(ctx, input.Sudo)); err != nil {
			return nil, err
		}
	}
	result, err := ty.NewTaskResult(input, true)
	if err != nil {
		return nil, fmt.Errorf("failed to create task result: %w", err)
	}
	return result, nil
}

// Run runs the Starship installer based on the OS.
func (t *StarshipInstaller) Run(ctx context.Context) error {
	switch {
	case t.Utils.IsOSX():
		if err := t.Curl.Run(ctx); err != nil {
			return fmt.Errorf("failed to install starship for OSX: %w", err)
		}
	case t.Utils.IsLinux():
		if err := t.LinuxCurl.Run(ctx); err != nil {
			return fmt.Errorf("failed to install starship for linux: %w", err)
		}
	default:
		return fmt.Errorf("fail to install starship since this is neither linux nor OSX")
	}
	return nil
}
