// Package installer provides Task types and related utilities.
package installer

import (
	ext "awong/dotfiles/pkg/external"
	l "awong/dotfiles/pkg/logging"
	ty "awong/dotfiles/pkg/types"
	"context"
	"fmt"
)

// VimInstaller installs vim via Homebrew.
type VimInstaller struct {
	Brew  ext.Brew
	Utils ext.Ext
	Log   l.Logger
}

// NewVimInstaller returns a new VimInstaller for installing vim.
func NewVimInstaller() *VimInstaller {
	return &VimInstaller{
		Brew:  ext.NewBrew(),
		Utils: ext.NewExt(),
		Log:   l.Log(),
	}
}

// Install installs the task.
func (t *VimInstaller) Install(ctx context.Context, input *ty.TaskInput) (*ty.TaskResult, error) {
	if t.Utils.IsInstalled("ovim") {
		t.Log.Infof("%v %v: vim is already installed", TaskEllipsis, input.Play)
		result, err := ty.NewTaskResult(input, false)
		if err != nil {
			return nil, fmt.Errorf("failed to create task result: %w", err)
		}
		return result, nil
	}
	t.Log.Infof("%v %v: install vim", TaskEllipsis, input.Play)
	if !input.DryRun {
		if err := t.Brew.Install(ctx, "vim"); err != nil {
			return nil, fmt.Errorf("fail to install vim: %v", err)
		}
	}
	result, err := ty.NewTaskResult(input, true)
	if err != nil {
		return nil, fmt.Errorf("failed to create task result: %w", err)
	}
	return result, nil
}
