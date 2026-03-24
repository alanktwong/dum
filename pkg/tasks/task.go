// Package tasks provides Task types and related utilities.
package tasks

import (
	"awong/dotfiles/pkg/types"
	"context"
)

const (
	// TaskEllipsis is prefix used in logs and println.
	TaskEllipsis = "..........."
)

// Task represents a unit of work that can be installed or listed.
type Task interface {
	Installer
	List(ctx context.Context, input *types.TaskInput) (*types.TaskResult, error)
}
