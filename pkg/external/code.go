package external

import (
	"context"
	"fmt"
	"os/exec"
)

// Code is an abstraction for vscode.
type Code interface {
	// InstallExtension installs a VS Code extension.
	InstallExtension(ctx context.Context, formula string) error
	// ListExtensions lists the installed  VS Code extensions.
	ListExtensions(ctx context.Context) (string, error)
}

// CodeImpl implements Code.
type CodeImpl struct{}

// NewCode constructs CodeImpl.
func NewCode() *CodeImpl {
	return &CodeImpl{}
}

// InstallExtension implements Code.
func (m *CodeImpl) InstallExtension(ctx context.Context, formula string) error {
	cmd := exec.CommandContext(ctx, "code", "--install-extension", formula)
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to code --install-extension %s: %v", formula, err)
	}
	return nil
}

// ListExtensions implements Code.
func (m *CodeImpl) ListExtensions(ctx context.Context) (string, error) {
	cmd := exec.CommandContext(ctx, "code", "--list-extensions")
	out, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("failed to code --list-extensions: %v", err)
	}
	return string(out), nil
}
