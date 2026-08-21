// Package tasks provides Task types and related utilities.
package tasks

import (
	"reflect"

	i "github.com/alanktwong/dum/internal/tasks/installer"
)

const (
	// TaskEllipsis is prefix used in logs and println.
	TaskEllipsis = i.TaskEllipsis
)

// Task represents a unit of work that can be installed.
type Task interface {
	i.Installer
}

// TaskTypeName returns the struct name of a task.
func TaskTypeName(t Task) string {
	return reflect.TypeOf(t).Elem().Name()
}
