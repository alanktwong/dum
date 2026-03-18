package playbook

import (
	"awong/dotfiles/pkg/external"
	"awong/dotfiles/pkg/logging"
	"context"
	"fmt"
)

// VimInstaller installs vim so that it stays out of the way of neovim.
type VimInstaller struct {
	Brew  external.Brew
	Utils external.Ext
	Log   logging.Logger
}

// NewVimInstaller constructs a VimInstaller.
func NewVimInstaller() *VimInstaller {
	return &VimInstaller{
		Brew:  external.NewBrew(),
		Utils: external.NewExt(),
		Log:   logging.Log(),
	}
}

// Install implements Installer.
func (t *VimInstaller) Install(ctx context.Context, input *Input) (*TaskResult, error) {
	if t.Utils.IsInstalled("ovim") {
		t.Log.Infof("%v %v: vim is already installed", TaskEllipsis, input.Play)
		return NewTaskResult(input, false)
	}
	t.Log.Infof("%v %v: install vim", TaskEllipsis, input.Play)
	if !input.DryRun {
		if err := t.Brew.Install(ctx, "vim"); err != nil {
			return nil, fmt.Errorf("fail to install vim: %v", err)
		}
	}
	return NewTaskResult(input, true)
}
