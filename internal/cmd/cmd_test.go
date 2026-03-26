package cmd

import (
	"os"
	"testing"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
)

func TestGetDefaultConfig_Default(t *testing.T) {
	t.Setenv("INSTALLER_CONFIG", "")
	result := getDefaultConfig()
	assert.Equal(t, "~/.config/installer.yml", result)
}

func TestGetDefaultConfig_FromEnv(t *testing.T) {
	t.Setenv("INSTALLER_CONFIG", "/custom/path.yml")
	result := getDefaultConfig()
	assert.Equal(t, "/custom/path.yml", result)
}

func TestAddGroupFlag(t *testing.T) {
	cmd := &cobra.Command{Use: "test"}
	addGroupFlag(cmd)

	flag := cmd.Flags().Lookup("group")
	assert.NotNil(t, flag)
	assert.Equal(t, "g", flag.Shorthand)
}

func TestAddFileFlag(t *testing.T) {
	t.Setenv("INSTALLER_CONFIG", "")
	cmd := &cobra.Command{Use: "test"}
	addFileFlag(cmd)

	flag := cmd.Flags().Lookup("file")
	assert.NotNil(t, flag)
	assert.Equal(t, "f", flag.Shorthand)
}

func TestAddVerboseFlag(t *testing.T) {
	cmd := &cobra.Command{Use: "test"}
	addVerboseFlag(cmd)

	flag := cmd.Flags().Lookup("verbose")
	assert.NotNil(t, flag)
	assert.Equal(t, "v", flag.Shorthand)
}

func TestAddDryRunFlag(t *testing.T) {
	cmd := &cobra.Command{Use: "test"}
	addDryRunFlag(cmd)

	flag := cmd.Flags().Lookup("dryrun")
	assert.NotNil(t, flag)
	assert.Equal(t, "d", flag.Shorthand)
}

func TestAddPrefixFlag(t *testing.T) {
	cmd := &cobra.Command{Use: "test"}
	addPrefixFlag(cmd)

	flag := cmd.Flags().Lookup("prefix")
	assert.NotNil(t, flag)
	assert.Equal(t, "p", flag.Shorthand)
}

func TestMain(m *testing.M) {
	os.Exit(m.Run())
}
