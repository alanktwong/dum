package tasks

import (
	"awong/dotfiles/pkg/types"
	"context"
)

const (
	// TaskEllipsis is prefix used in logs and println.
	TaskEllipsis = "..........."
)

type Task interface {
	Installer
	List(ctx context.Context, input *types.TaskInput) (*types.TaskResult, error)
}
