package external

import (
	"context"
	"fmt"
	"os/exec"

	"github.com/alanktwong/dum/internal/external/gen"
)

// Mas is an abstraction for Apple app store.
type Mas = gen.Mas

// MasImpl implements Mas.
type MasImpl struct {
	ext Ext
}

// NewMas constructs MasImpl.
func NewMas() *MasImpl {
	return &MasImpl{ext: NewExt()}
}

// Install implements Mas.
func (m *MasImpl) Install(ctx context.Context, app string) error {
	if !m.ext.IsOSX() {
		return fmt.Errorf("cannot install mas app %s outside of macOS", app)
	}
	cmd := exec.CommandContext(ctx, "mas", "install", app)
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to mas install %s: %w", app, err)
	}
	return nil
}

// List implements Mas.
func (m *MasImpl) List(ctx context.Context) (string, error) {
	if !m.ext.IsOSX() {
		return "", fmt.Errorf("cannot list mas apps outside of macOS")
	}
	cmd := exec.CommandContext(ctx, "mas", "list")
	out, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("failed to mas list: %w", err)
	}
	return string(out), nil
}
