package main

import (
	"awong/dotfiles/cmd/dum/cli"
	"os"
	"testing"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
)

func TestGetDefaultConfig_Default(t *testing.T) {
	t.Setenv("INSTALLER_CONFIG", "")
	t.Setenv("XDG_CONFIG_HOME", "")
	result := cli.GetDefaultConfig()
	assert.Equal(t, "~/.config/dum/installer.yml", result)
}

func TestGetDefaultConfig_XDGConfigHome(t *testing.T) {
	t.Setenv("INSTALLER_CONFIG", "")
	t.Setenv("XDG_CONFIG_HOME", "/custom/config")
	result := cli.GetDefaultConfig()
	assert.Equal(t, "/custom/config/dum/installer.yml", result)
}

func TestGetDefaultConfig_FromEnv(t *testing.T) {
	t.Setenv("INSTALLER_CONFIG", "/custom/path.yml")
	result := cli.GetDefaultConfig()
	assert.Equal(t, "/custom/path.yml", result)
}

func TestAddGroupFlag(t *testing.T) {
	cmd := &cobra.Command{Use: "test"}
	cli.AddGroupFlag(cmd)

	flag := cmd.Flags().Lookup("group")
	assert.NotNil(t, flag)
	assert.Equal(t, "g", flag.Shorthand)
}

func TestAddFileFlag(t *testing.T) {
	t.Setenv("INSTALLER_CONFIG", "")
	cmd := &cobra.Command{Use: "test"}
	cli.AddFileFlag(cmd)

	flag := cmd.Flags().Lookup("file")
	assert.NotNil(t, flag)
	assert.Equal(t, "f", flag.Shorthand)
}

func TestAddVerboseFlag(t *testing.T) {
	cmd := &cobra.Command{Use: "test"}
	cli.AddVerboseFlag(cmd)

	flag := cmd.Flags().Lookup("verbose")
	assert.NotNil(t, flag)
	assert.Equal(t, "v", flag.Shorthand)
}

func TestAddDryRunFlag(t *testing.T) {
	cmd := &cobra.Command{Use: "test"}
	cli.AddDryRunFlag(cmd)

	flag := cmd.Flags().Lookup("dryrun")
	assert.NotNil(t, flag)
	assert.Equal(t, "d", flag.Shorthand)
}

func TestAddPrefixFlag(t *testing.T) {
	cmd := &cobra.Command{Use: "test"}
	cli.AddPrefixFlag(cmd)

	flag := cmd.Flags().Lookup("prefix")
	assert.NotNil(t, flag)
	assert.Equal(t, "p", flag.Shorthand)
}

func TestMain(m *testing.M) {
	os.Exit(m.Run())
}
