package cmd

import (
	"testing"

	"github.com/charmbracelet/log"
	"github.com/spf13/cobra"
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

func TestGetVerbosityFromCommand_NoFlag(t *testing.T) {
	cmd := &cobra.Command{Use: "test"}
	v := GetVerbosityFromCommand(cmd)
	assert.Equal(t, uint8(0), v.Verbose)
}

func TestGetVerbosityFromCommand_FlagNotChanged(t *testing.T) {
	cmd := &cobra.Command{Use: "test"}
	addVerboseFlag(cmd)
	v := GetVerbosityFromCommand(cmd)
	assert.Equal(t, uint8(0), v.Verbose)
}

func TestGetVerbosityFromCommand_SingleVerbose(t *testing.T) {
	cmd := &cobra.Command{Use: "test"}
	addVerboseFlag(cmd)
	cmd.SetArgs([]string{"-v"})
	_ = cmd.Execute()
	v := GetVerbosityFromCommand(cmd)
	assert.Equal(t, uint8(1), v.Verbose)
	assert.Equal(t, log.InfoLevel, v.Level())
}

func TestGetVerbosityFromCommand_DoubleVerbose(t *testing.T) {
	cmd := &cobra.Command{Use: "test"}
	addVerboseFlag(cmd)
	cmd.SetArgs([]string{"-vv"})
	_ = cmd.Execute()
	v := GetVerbosityFromCommand(cmd)
	assert.Equal(t, uint8(2), v.Verbose)
	assert.Equal(t, log.DebugLevel, v.Level())
}
