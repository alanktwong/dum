// Package tasks provides Task types and related utilities.
package tasks

import (
	"awong/dotfiles/pkg/types"
	"context"
)

// Installer can install given an input.
type Installer interface {
	Install(ctx context.Context, input *types.TaskInput) (*types.TaskResult, error)
}
