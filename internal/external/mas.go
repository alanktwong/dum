package external

import (
	"awong/dotfiles/internal/external/gen"
	"context"
	"fmt"
	"os/exec"
)

// Mas is an abstraction for Apple app store.
type Mas = gen.Mas

// MasImpl implements Mas.
type MasImpl struct{}

// NewMas constructs MasImpl.
func NewMas() *MasImpl {
	return &MasImpl{}
}

// Install implements Mas.
func (m *MasImpl) Install(ctx context.Context, app string) error {
	cmd := exec.CommandContext(ctx, "mas", "install", app)
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to mas install %s: %w", app, err)
	}
	return nil
}

// List implements Mas.
func (m *MasImpl) List(ctx context.Context) (string, error) {
	cmd := exec.CommandContext(ctx, "mas", "list")
	out, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("failed to mas list: %w", err)
	}
	return string(out), nil
}
