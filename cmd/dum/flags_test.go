package main

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
)

func TestGetDefaultConfig_Default(t *testing.T) {
	t.Setenv("DUM_CONFIG", "")
	t.Setenv("XDG_CONFIG_HOME", "")
	t.Chdir(t.TempDir())
	result := GetDefaultConfig()
	assert.Equal(t, "~/.config/dum/dum.yml", result)
}

func TestGetDefaultConfig_XDGConfigHome(t *testing.T) {
	t.Setenv("DUM_CONFIG", "")
	t.Setenv("XDG_CONFIG_HOME", "/custom/config")
	t.Chdir(t.TempDir())
	result := GetDefaultConfig()
	assert.Equal(t, "/custom/config/dum/dum.yml", result)
}

func TestGetDefaultConfig_FromEnv(t *testing.T) {
	t.Setenv("DUM_CONFIG", "/custom/path.yml")
	result := GetDefaultConfig()
	assert.Equal(t, "/custom/path.yml", result)
}

func TestGetDefaultConfig_EnvTakesPrecedence(t *testing.T) {
	t.Setenv("DUM_CONFIG", "/custom/path.yml")
	t.Setenv("XDG_CONFIG_HOME", "/custom/config")

	assert.Equal(t, "/custom/path.yml", GetDefaultConfig())
}

func TestGetDefaultConfig_CWDFile(t *testing.T) {
	t.Setenv("DUM_CONFIG", "")
	t.Setenv("XDG_CONFIG_HOME", "")
	dir := t.TempDir()
	t.Chdir(dir)
	assert.NoError(t, os.WriteFile(filepath.Join(dir, "dum.yml"), []byte("playbook:\n  book: test\n"), 0o600))

	assert.Equal(t, "dum.yml", GetDefaultConfig())
}

func TestGetDefaultConfig_XDGFile(t *testing.T) {
	t.Setenv("DUM_CONFIG", "")
	xdg := t.TempDir()
	t.Setenv("XDG_CONFIG_HOME", xdg)
	t.Chdir(t.TempDir())
	assert.NoError(t, os.MkdirAll(filepath.Join(xdg, "dum"), 0o700))
	assert.NoError(t, os.WriteFile(filepath.Join(xdg, "dum", "dum.yml"), []byte("playbook:\n  book: test\n"), 0o600))

	assert.Equal(t, filepath.Join(xdg, "dum", "dum.yml"), GetDefaultConfig())
}

func TestGetDefaultConfig_LegacyInstallerFallback(t *testing.T) {
	t.Setenv("DUM_CONFIG", "")
	xdg := t.TempDir()
	t.Setenv("XDG_CONFIG_HOME", xdg)
	t.Chdir(t.TempDir())
	assert.NoError(t, os.MkdirAll(filepath.Join(xdg, "dum"), 0o700))
	assert.NoError(
		t,
		os.WriteFile(filepath.Join(xdg, "dum", "installer.yml"), []byte("playbook:\n  book: test\n"), 0o600),
	)

	assert.Equal(t, filepath.Join(xdg, "dum", "installer.yml"), GetDefaultConfig())
}

func TestGetDefaultConfig_CWDBeatsXDGAndLegacy(t *testing.T) {
	t.Setenv("DUM_CONFIG", "")
	xdg := t.TempDir()
	t.Setenv("XDG_CONFIG_HOME", xdg)
	dir := t.TempDir()
	t.Chdir(dir)
	assert.NoError(t, os.MkdirAll(filepath.Join(xdg, "dum"), 0o700))
	assert.NoError(t, os.WriteFile(filepath.Join(xdg, "dum", "dum.yml"), []byte("playbook:\n  book: xdg\n"), 0o600))
	assert.NoError(
		t,
		os.WriteFile(filepath.Join(xdg, "dum", "installer.yml"), []byte("playbook:\n  book: legacy\n"), 0o600),
	)
	assert.NoError(t, os.WriteFile(filepath.Join(dir, "dum.yml"), []byte("playbook:\n  book: cwd\n"), 0o600))

	assert.Equal(t, "dum.yml", GetDefaultConfig())
}

func TestGetDefaultConfig_RepeatedCallsAreEnvironmentIsolated(t *testing.T) {
	t.Setenv("DUM_CONFIG", "/first/path.yml")
	t.Setenv("XDG_CONFIG_HOME", "")
	t.Chdir(t.TempDir())
	assert.Equal(t, "/first/path.yml", GetDefaultConfig())

	t.Setenv("DUM_CONFIG", "")
	t.Setenv("XDG_CONFIG_HOME", "/second/config")
	assert.Equal(t, "/second/config/dum/dum.yml", GetDefaultConfig())

	t.Setenv("XDG_CONFIG_HOME", "")
	assert.Equal(t, "~/.config/dum/dum.yml", GetDefaultConfig())
}

func TestAddGroupFlag(t *testing.T) {
	cmd := &cobra.Command{Use: "test"}
	AddGroupFlag(cmd)

	flag := cmd.Flags().Lookup("group")
	assert.NotNil(t, flag)
	assert.Equal(t, "g", flag.Shorthand)
}

func TestAddFileFlag(t *testing.T) {
	t.Setenv("DUM_CONFIG", "")
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
