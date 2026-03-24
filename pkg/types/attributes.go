// Package types provides common types used throughout the application.
package types

import (
	"context"
	"fmt"
)

// Attributable provides access to Attributes.
type Attributable interface {
	GetAttributes() Attributes
}

// Attributes holds common attributes for tasks and plays.
type Attributes struct {
	ID          string
	Description string
	Enabled     bool
	Sudo        bool
}

// NewAttributes creates a new Attributes instance.
func NewAttributes(id, description string, enabled, sudo bool) (*Attributes, error) {
	if id == "" {
		return nil, fmt.Errorf("attribute ID cannot be empty")
	}
	return &Attributes{
		ID:          id,
		Description: description,
		Enabled:     enabled,
		Sudo:        sudo,
	}, nil
}

// IsEnabled returns whether the attribute is enabled.
func (a *Attributes) IsEnabled() bool {
	return a.Enabled
}

// GetID returns the ID of the attribute.
func (a *Attributes) GetID() string {
	return a.ID
}

// CreateTaskResult creates a new TaskResult from the attribute.
func (a *Attributes) CreateTaskResult(input *TaskInput, success bool) (*TaskResult, error) {
	return &TaskResult{
		Task:     a.ID,
		Play:     input.Play,
		PlayBook: input.PlayBook,
		DryRun:   input.DryRun,
		Success:  success,
	}, nil
}

type sudoCtxKey struct{}

// WithSudo adds sudo privilege to the context.
func WithSudo(ctx context.Context, sudo bool) context.Context {
	return context.WithValue(ctx, sudoCtxKey{}, sudo)
}

// GetSudo retrieves sudo privilege from the context.
func GetSudo(ctx context.Context) bool {
	sudo, ok := ctx.Value(sudoCtxKey{}).(bool)
	if !ok {
		return false
	}
	return sudo
}
