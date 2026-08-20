package plays

import (
	"context"

	ty "awong/dotfiles/internal/types"

	"github.com/stretchr/testify/mock"
)

// MockTask is a mock for the t.Task interface that also implements ty.Attributable.
type MockTask struct {
	mock.Mock
	Attr ty.Attributes
}

// GetAttributes implements ty.Attributable.
func (m *MockTask) GetAttributes() ty.Attributes {
	return m.Attr
}

// Install implements i.Installer.
func (m *MockTask) Install(ctx context.Context, input *ty.TaskInput) (*ty.TaskResult, error) {
	ret := m.Called(ctx, input)
	var r0 *ty.TaskResult
	if ret.Get(0) != nil {
		r0 = ret.Get(0).(*ty.TaskResult)
	}
	//nolint:wrapcheck
	return r0, ret.Error(1)
}
