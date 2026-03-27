// Package installer provides Task types and related utilities.
package installer

import (
	ty "awong/dotfiles/internal/types"
	"context"
)

const (
	// TaskEllipsis is prefix used in logs and println.
	TaskEllipsis = "..........."
)

// Installer can install given an input.
type Installer interface {
	Install(ctx context.Context, input *ty.TaskInput) (*ty.TaskResult, error)
}
