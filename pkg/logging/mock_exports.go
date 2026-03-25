package logging

import (
	"awong/dotfiles/pkg/logging/gen"
	"testing"
)

// MockLogger is a mock of the Logger interface.
type MockLogger = gen.MockLogger

// NewMockLogger creates a new mock Logger for testing.
func NewMockLogger(t *testing.T) *gen.MockLogger {
	return gen.NewMockLogger(t)
}
