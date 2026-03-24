// Package tasks provides Task types and related utilities.
package tasks

import (
	ext "awong/dotfiles/pkg/external"
	l "awong/dotfiles/pkg/logging"
	ty "awong/dotfiles/pkg/types"
	"context"
	"fmt"
	"path/filepath"
)

// BashInstaller installs bash on macOS.
type BashInstaller struct {
	Brew  ext.Brew
	Utils ext.Ext
	Log   l.Logger
}

// NewBashInstaller returns a new BashInstaller for installing bash.
func NewBashInstaller() *BashInstaller {
	return &BashInstaller{
		Brew:  ext.NewBrew(),
		Utils: ext.NewExt(),
		Log:   l.Log(),
	}
}

// Install installs the task.
func (t *BashInstaller) Install(ctx context.Context, input *ty.TaskInput) (*ty.TaskResult, error) {
	if !t.Utils.IsOSX() {
		t.Log.Infof("%v %v: there is no need to reinstall bash outside of Mac OSX", TaskEllipsis, input.Play)
		result, err := ty.NewTaskResult(input, false)
		if err != nil {
			return nil, fmt.Errorf("failed to create task result: %w", err)
		}
		return result, nil
	}
	if t.Utils.IsInstalled("bash") {
		t.Log.Infof("%v %v: updated bash is already installed", TaskEllipsis, input.Play)
		result, err := ty.NewTaskResult(input, true)
		if err != nil {
			return nil, fmt.Errorf("failed to create task result: %w", err)
		}
		return result, nil
	}

	t.Log.Infof("%v %v: install bash", TaskEllipsis, input.Play)
	if !input.DryRun {
		formula := "bash"
		if err := t.Brew.Install(ctx, formula); err != nil {
			return nil, fmt.Errorf("failed to install bash: %w", err)
		}
		if err := t.buildLink(ctx, formula, input.Sudo); err != nil {
			return nil, err
		}
	}
	result, err := ty.NewTaskResult(input, true)
	if err != nil {
		return nil, fmt.Errorf("failed to create task result: %w", err)
	}
	return result, nil
}

func (t *BashInstaller) buildLink(ctx context.Context, formula string, sudo bool) error {
	homebrewPrefix, err := t.Brew.Prefix(ctx)
	if err != nil {
		return fmt.Errorf("failed to get homebrew prefix: %w", err)
	}
	src := filepath.Join(homebrewPrefix, "bin", formula)
	if err := t.Utils.SoftLink(ctx, "/usr/local/bin", src, formula, sudo); err != nil {
		return fmt.Errorf("failed to make soft link for bash: %w", err)
	}
	return nil
}
