package types

import (
	"context"
	"fmt"
)

type Attributable interface {
	GetAttributes() Attributes
}

type Attributes struct {
	ID          string
	Description string
	Enabled     bool
	Sudo        bool
}

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

func (a *Attributes) IsEnabled() bool {
	return a.Enabled
}

func (a *Attributes) GetID() string {
	return a.ID
}

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

func WithSudo(ctx context.Context, sudo bool) context.Context {
	return context.WithValue(ctx, sudoCtxKey{}, sudo)
}

func GetSudo(ctx context.Context) bool {
	sudo, ok := ctx.Value(sudoCtxKey{}).(bool)
	if !ok {
		return false
	}
	return sudo
}
