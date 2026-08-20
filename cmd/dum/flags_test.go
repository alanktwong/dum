package main

import (
	"os"
	"testing"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
)

func TestGetDefaultConfig_Default(t *testing.T) {
	t.Setenv("INSTALLER_CONFIG", "")
	t.Setenv("XDG_CONFIG_HOME", "")
	result := GetDefaultConfig()
	assert.Equal(t, "~/.config/dum/installer.yml", result)
}

func TestGetDefaultConfig_XDGConfigHome(t *testing.T) {
	t.Setenv("INSTALLER_CONFIG", "")
	t.Setenv("XDG_CONFIG_HOME", "/custom/config")
	result := GetDefaultConfig()
	assert.Equal(t, "/custom/config/dum/installer.yml", result)
}

func TestGetDefaultConfig_FromEnv(t *testing.T) {
	t.Setenv("INSTALLER_CONFIG", "/custom/path.yml")
	result := GetDefaultConfig()
	assert.Equal(t, "/custom/path.yml", result)
}

func TestGetDefaultConfig_InstallerConfigTakesPrecedence(t *testing.T) {
	t.Setenv("INSTALLER_CONFIG", "/custom/path.yml")
	t.Setenv("XDG_CONFIG_HOME", "/custom/config")

	assert.Equal(t, "/custom/path.yml", GetDefaultConfig())
}

func TestGetDefaultConfig_RepeatedCallsAreEnvironmentIsolated(t *testing.T) {
	t.Setenv("INSTALLER_CONFIG", "/first/path.yml")
	t.Setenv("XDG_CONFIG_HOME", "")
	assert.Equal(t, "/first/path.yml", GetDefaultConfig())

	t.Setenv("INSTALLER_CONFIG", "")
	t.Setenv("XDG_CONFIG_HOME", "/second/config")
	assert.Equal(t, "/second/config/dum/installer.yml", GetDefaultConfig())

	t.Setenv("XDG_CONFIG_HOME", "")
	assert.Equal(t, "~/.config/dum/installer.yml", GetDefaultConfig())
}

func TestAddGroupFlag(t *testing.T) {
	cmd := &cobra.Command{Use: "test"}
	AddGroupFlag(cmd)

	flag := cmd.Flags().Lookup("group")
	assert.NotNil(t, flag)
	assert.Equal(t, "g", flag.Shorthand)
}

func TestAddFileFlag(t *testing.T) {
	t.Setenv("INSTALLER_CONFIG", "")
	cmd := &cobra.Command{Use: "test"}
	AddFileFlag(cmd)

	flag := cmd.Flags().Lookup("file")
	assert.NotNil(t, flag)
	assert.Equal(t, "f", flag.Shorthand)
}

func TestAddVerboseFlag(t *testing.T) {
	cmd := &cobra.Command{Use: "test"}
	AddVerboseFlag(cmd)

	flag := cmd.Flags().Lookup("verbose")
	assert.NotNil(t, flag)
	assert.Equal(t, "v", flag.Shorthand)
}

func TestAddDryRunFlag(t *testing.T) {
	cmd := &cobra.Command{Use: "test"}
	AddDryRunFlag(cmd)

	flag := cmd.Flags().Lookup("dryrun")
	assert.NotNil(t, flag)
	assert.Equal(t, "d", flag.Shorthand)
}

func TestAddPrefixFlag(t *testing.T) {
	cmd := &cobra.Command{Use: "test"}
	AddPrefixFlag(cmd)

	flag := cmd.Flags().Lookup("prefix")
	assert.NotNil(t, flag)
	assert.Equal(t, "p", flag.Shorthand)
}

func TestMain(m *testing.M) {
	os.Exit(m.Run())
}
