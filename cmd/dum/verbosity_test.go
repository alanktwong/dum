package main

import (
	"awong/dotfiles/cmd/dum/cli"
	"testing"

	clog "github.com/charmbracelet/log"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
)

func TestVerbosity(t *testing.T) {
	v := cli.Verbosity{}
	assert.Equal(t, uint8(0), v.Verbose)
	assert.Equal(t, clog.WarnLevel, v.Level())

	// Test verbosity with different levels
	v = cli.Verbosity{Verbose: 1}
	assert.Equal(t, uint8(1), v.Verbose)
	assert.Equal(t, clog.InfoLevel, v.Level())

	v = cli.Verbosity{Verbose: 2}
	assert.Equal(t, uint8(2), v.Verbose)
	assert.Equal(t, clog.DebugLevel, v.Level())

	v = cli.Verbosity{Verbose: 3}
	assert.Equal(t, uint8(3), v.Verbose)
	assert.Equal(t, clog.DebugLevel, v.Level())
}

func TestGetVerbosityFromCommand_NoFlag(t *testing.T) {
	cmd := &cobra.Command{Use: "test"}
	v := cli.GetVerbosityFromCommand(cmd)
	assert.Equal(t, uint8(0), v.Verbose)
}

func TestGetVerbosityFromCommand_FlagNotChanged(t *testing.T) {
	cmd := &cobra.Command{Use: "test"}
	cli.AddVerboseFlag(cmd)
	v := cli.GetVerbosityFromCommand(cmd)
	assert.Equal(t, uint8(0), v.Verbose)
}

func TestGetVerbosityFromCommand_SingleVerbose(t *testing.T) {
	cmd := &cobra.Command{Use: "test"}
	cli.AddVerboseFlag(cmd)
	cmd.SetArgs([]string{"-v"})
	_ = cmd.Execute()
	v := cli.GetVerbosityFromCommand(cmd)
	assert.Equal(t, uint8(1), v.Verbose)
	assert.Equal(t, clog.InfoLevel, v.Level())
}

func TestGetVerbosityFromCommand_DoubleVerbose(t *testing.T) {
	cmd := &cobra.Command{Use: "test"}
	cli.AddVerboseFlag(cmd)
	cmd.SetArgs([]string{"-vv"})
	_ = cmd.Execute()
	v := cli.GetVerbosityFromCommand(cmd)
	assert.Equal(t, uint8(2), v.Verbose)
	assert.Equal(t, clog.DebugLevel, v.Level())
}
