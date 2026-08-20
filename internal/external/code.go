package external

import (
	"alanktwong/dum/internal/external/gen"
	"context"
	"fmt"
	"os/exec"
)

// Code is an abstraction for vscode.
type Code = gen.Code

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
		return fmt.Errorf("failed to code --install-extension %s: %w", formula, err)
	}
	return nil
}

// ListExtensions implements Code.
func (m *CodeImpl) ListExtensions(ctx context.Context) (string, error) {
	cmd := exec.CommandContext(ctx, "code", "--list-extensions")
	out, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("failed to code --list-extensions: %w", err)
	}
	return string(out), nil
}
