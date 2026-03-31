package cmd

import (
	"testing"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"

	gen "awong/dotfiles/internal/logging/gen"
)

func TestNewDum_RootCommand(t *testing.T) {
	dum := NewDum()

	assert.NotNil(t, dum)
	assert.NotNil(t, dum.Cmd)
	assert.Equal(t, "dum", dum.Cmd.Use)
	assert.Equal(t, version, dum.Cmd.Version)
}

func TestNewDum_Subcommands(t *testing.T) {
	dum := NewDum()

	subcommands := make(map[string]bool)
	for _, cmd := range dum.Cmd.Commands() {
		subcommands[cmd.Use] = true
	}

	assert.True(t, subcommands["install"])
	assert.True(t, subcommands["list"])
	assert.True(t, subcommands["rename"])
	assert.True(t, subcommands["log"])
	assert.True(t, subcommands["schema"])
}

func TestNewDum_SubcommandAliases(t *testing.T) {
	dum := NewDum()

	tests := []struct {
		name          string
		expectedAlias string
	}{
		{"install", "i"},
		{"list", "ls"},
		{"rename", "r"},
		{"log", "lg"},
		{"schema", "s"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := findCommand(dum.Cmd.Commands(), tt.name)
			assert.NotNil(t, cmd)
			assert.Contains(t, cmd.Aliases, tt.expectedAlias)
		})
	}
}

func TestNewDum_LogSubcommands(t *testing.T) {
	dum := NewDum()

	logCmd := findCommand(dum.Cmd.Commands(), "log")
	assert.NotNil(t, logCmd)

	logSubcommands := make(map[string]bool)
	for _, cmd := range logCmd.Commands() {
		logSubcommands[cmd.Use] = true
	}

	assert.True(t, logSubcommands["debug"])
	assert.True(t, logSubcommands["info"])
	assert.True(t, logSubcommands["success"])
	assert.True(t, logSubcommands["warn"])
	assert.True(t, logSubcommands["error"])
}

func TestNewDum_LogSubcommandAliases(t *testing.T) {
	dum := NewDum()

	logCmd := findCommand(dum.Cmd.Commands(), "log")
	assert.NotNil(t, logCmd)

	tests := []struct {
		name          string
		expectedAlias string
	}{
		{"debug", "d"},
		{"info", "i"},
		{"success", "s"},
		{"warn", "w"},
		{"error", "e"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := findCommand(logCmd.Commands(), tt.name)
			assert.NotNil(t, cmd)
			assert.Contains(t, cmd.Aliases, tt.expectedAlias)
		})
	}
}

func TestNewDum_DependenciesSet(t *testing.T) {
	dum := NewDum()

	assert.NotNil(t, dum.FactoryProvider)
	assert.NotNil(t, dum.InstallExecutor)
	assert.NotNil(t, dum.ListExecutor)
	assert.NotNil(t, dum.Log)
}

func TestExec_Success(t *testing.T) {
	logger := gen.NewMockLogger(t)
	cmd := &cobra.Command{Use: "test"}
	dum := &Dum{Log: logger, Cmd: cmd}

	// Should not panic; Exec() calls Cmd.Execute() which succeeds with no RunE.
	dum.Exec()
}

func findCommand(cmds []*cobra.Command, name string) *cobra.Command {
	for _, cmd := range cmds {
		if cmd.Use == name {
			return cmd
		}
	}
	return nil
}
