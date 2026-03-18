package playbook

import (
	"awong/dotfiles/pkg/external"
	"awong/dotfiles/pkg/logging"
	"context"
	"fmt"
	"path/filepath"
)

// BashInstaller installs an updated version of bash for Mac OSX.
type BashInstaller struct {
	Brew  external.Brew
	Utils external.Ext
	Log   logging.Logger
}

// NewBashInstaller constructs a BashInstaller.
func NewBashInstaller() *BashInstaller {
	return &BashInstaller{
		Brew:  external.NewBrew(),
		Utils: external.NewExt(),
		Log:   logging.Log(),
	}
}

// Install implements Installer.
func (t *BashInstaller) Install(ctx context.Context, input *Input) (*TaskResult, error) {
	if !t.Utils.IsOSX() {
		t.Log.Infof("%v %v: there is no need to reinstall bash outside of Mac OSX", TaskEllipsis, input.Play)
		return NewTaskResult(input, false)
	}
	if t.Utils.IsInstalled("bash") {
		t.Log.Infof("%v %v: updated bash is already installed", TaskEllipsis, input.Play)
		return NewTaskResult(input, true)
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
	return NewTaskResult(input, true)
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
