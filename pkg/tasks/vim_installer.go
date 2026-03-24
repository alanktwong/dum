package tasks

import (
	ext "awong/dotfiles/pkg/external"
	l "awong/dotfiles/pkg/logging"
	ty "awong/dotfiles/pkg/types"
	"context"
	"fmt"
)

type VimInstaller struct {
	Brew  ext.Brew
	Utils ext.Ext
	Log   l.Logger
}

func NewVimInstaller() *VimInstaller {
	return &VimInstaller{
		Brew:  ext.NewBrew(),
		Utils: ext.NewExt(),
		Log:   l.Log(),
	}
}

func (t *VimInstaller) Install(ctx context.Context, input *ty.TaskInput) (*ty.TaskResult, error) {
	if t.Utils.IsInstalled("ovim") {
		t.Log.Infof("%v %v: vim is already installed", TaskEllipsis, input.Play)
		return ty.NewTaskResult(input, false)
	}
	t.Log.Infof("%v %v: install vim", TaskEllipsis, input.Play)
	if !input.DryRun {
		if err := t.Brew.Install(ctx, "vim"); err != nil {
			return nil, fmt.Errorf("fail to install vim: %v", err)
		}
	}
	return ty.NewTaskResult(input, true)
}
