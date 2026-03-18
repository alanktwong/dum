package cmd

import (
	"testing"

	"github.com/charmbracelet/log"
	"github.com/stretchr/testify/assert"
)

func TestVerbosity(t *testing.T) {
	v := Verbosity{}
	assert.Equal(t, uint8(0), v.Verbose)
	assert.Equal(t, log.WarnLevel, v.Level())

	// Test verbosity with different levels
	v = Verbosity{Verbose: 1}
	assert.Equal(t, uint8(1), v.Verbose)
	assert.Equal(t, log.InfoLevel, v.Level())

	v = Verbosity{Verbose: 2}
	assert.Equal(t, uint8(2), v.Verbose)
	assert.Equal(t, log.DebugLevel, v.Level())

	v = Verbosity{Verbose: 3}
	assert.Equal(t, uint8(3), v.Verbose)
	assert.Equal(t, log.DebugLevel, v.Level())
}
