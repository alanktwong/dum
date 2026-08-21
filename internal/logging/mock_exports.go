package logging

import (
	"testing"

	"github.com/alanktwong/dum/internal/logging/gen"
)

// MockLogger is a mock of the Logger interface.
type MockLogger = gen.MockLogger

// NewMockLogger creates a new mock Logger for testing.
func NewMockLogger(t *testing.T) *gen.MockLogger {
	return gen.NewMockLogger(t)
}
