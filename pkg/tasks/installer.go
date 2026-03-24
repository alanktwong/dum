package tasks

import (
	"awong/dotfiles/pkg/types"
	"context"
)

type Installer interface {
	Install(ctx context.Context, input *types.TaskInput) (*types.TaskResult, error)
}
