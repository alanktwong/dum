package logging

import (
	"testing"

	"github.com/charmbracelet/log"
	"github.com/stretchr/testify/assert"
)

func TestNewLogger(t *testing.T) {
	options := Options{
		Prefix: "test",
		Level:  log.InfoLevel,
	}
	// when
	logger := NewLogger(options)
	loggerImpl, ok := logger.(*LoggerImpl)
	assert.True(t, ok)
	// then
	assert.NotNil(t, logger)
	assert.NotNil(t, loggerImpl.Log)
}
