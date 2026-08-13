// Package tasks provides Task types and related utilities.
package tasks

import (
	i "awong/dotfiles/internal/tasks/installer"
	ty "awong/dotfiles/internal/types"
	"context"
	"reflect"
)

const (
	// TaskEllipsis is prefix used in logs and println.
	TaskEllipsis = i.TaskEllipsis
)

// Task represents a unit of work that can be installed or listed.
type Task interface {
	i.Installer
	List(ctx context.Context, input *ty.TaskInput) (*ty.TaskResult, error)
}

// TaskTypeName returns the struct name of a task.
func TaskTypeName(t Task) string {
	return reflect.TypeOf(t).Elem().Name()
}
