package main

import (
	"bytes"
	"os"
	"os/exec"
	"strings"
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

func TestNewDum_LogLevelSubprocess(t *testing.T) {
	const helperEnv = "DUM_ROOT_LOG_LEVEL_HELPER"

	for _, tc := range []struct {
		name      string
		level     string
		wantDebug bool
	}{
		{name: "debug threshold", level: "debug", wantDebug: true},
		{name: "warn threshold", level: "warn", wantDebug: false},
	} {
		t.Run(tc.name, func(t *testing.T) {
			command := exec.Command(os.Args[0], "-test.run=^TestNewDum_LogLevelHelper$")
			env := make([]string, 0, len(os.Environ())+2)
			for _, value := range os.Environ() {
				if strings.HasPrefix(value, "ZSH_LOG_LEVEL=") || strings.HasPrefix(value, helperEnv+"=") {
					continue
				}
				env = append(env, value)
			}
			command.Env = append(env, helperEnv+"=1", "ZSH_LOG_LEVEL="+tc.level)

			var stderr bytes.Buffer
			command.Stderr = &stderr
			assert.NoError(t, command.Run())

			if tc.wantDebug {
				assert.Contains(t, stderr.String(), "debug-message")
			} else {
				assert.NotContains(t, stderr.String(), "debug-message")
			}
			assert.Contains(t, stderr.String(), "warn-message")
		})
	}
}

func TestNewDum_LogLevelHelper(t *testing.T) {
	if os.Getenv("DUM_ROOT_LOG_LEVEL_HELPER") != "1" {
		return
	}

	dum := NewDum()
	dum.Cmd.SetArgs([]string{"log", "debug", "debug-message"})
	if err := dum.Cmd.Execute(); err != nil {
		t.Fatal(err)
	}
	dum.Cmd.SetArgs([]string{"log", "warn", "warn-message"})
	if err := dum.Cmd.Execute(); err != nil {
		t.Fatal(err)
	}
}

func findCommand(cmds []*cobra.Command, name string) *cobra.Command {
	for _, cmd := range cmds {
		if cmd.Use == name {
			return cmd
		}
	}
	return nil
}
