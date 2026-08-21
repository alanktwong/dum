package external

import (
	"context"
	"fmt"
	"os/exec"
	"path/filepath"

	"github.com/alanktwong/dum/internal/external/gen"
)

// Brew is an abstraction for homebrew.
type Brew = gen.Brew

// BrewImpl implements Brew.
type BrewImpl struct {
	Utils Ext
}

// NewBrew constructs BrewImpl.
func NewBrew() *BrewImpl {
	return &BrewImpl{
		Utils: NewExt(),
	}
}

// Install implements Brew.
func (b *BrewImpl) Install(ctx context.Context, formula string) error {
	cmd := exec.CommandContext(ctx, "brew", "install", formula)
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to brew install %s: %w", formula, err)
	}
	return nil
}

// InstallCask implements Brew.
func (b *BrewImpl) InstallCask(ctx context.Context, formula string) error {
	cmd := exec.CommandContext(ctx, "brew", "install", "--cask", formula)
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to brew install --cask %s: %w", formula, err)
	}
	return nil
}

// Prefix implements Brew.
func (b *BrewImpl) Prefix(ctx context.Context) (string, error) {
	cmd := exec.CommandContext(ctx, "brew", "--prefix")
	out, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("failed to brew --prefix: %w", err)
	}
	return string(out), nil
}

// Tap implements Brew.
func (b *BrewImpl) Tap(ctx context.Context, tap string) error {
	cmd := exec.CommandContext(ctx, "brew", "tap", tap)
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to brew tap %s: %w", tap, err)
	}
	return nil
}

// InPath implements Brew.
func (b *BrewImpl) InPath(ctx context.Context, prefix, id string) bool {
	homebrewPrefix, err := b.Prefix(ctx)
	if err != nil {
		return false
	}
	path := filepath.Join(homebrewPrefix, prefix, id)
	return b.Utils.IsDir(path)
}
