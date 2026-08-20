package main

import (
	"context"
	"fmt"
	"testing"

	fy "alanktwong/dum/internal/factory"
	pb "alanktwong/dum/internal/playbook"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func newTestDumForInstall(t *testing.T) (*InstallCommand, *mockLogger, *mockFactoryProvider, *mockInstallExecutor) {
	t.Helper()
	logger := &mockLogger{}
	factory := &mockFactoryProvider{}
	executor := &mockInstallExecutor{}
	dum := &InstallCommand{
		Log:             logger,
		FactoryProvider: factory,
		Executor:        executor,
	}
	return dum, logger, factory, executor
}

func newTestInstallCmd(t *testing.T, dum *InstallCommand) *cobra.Command {
	t.Helper()
	return NewInstallCommand("test", dum)
}

func TestNewInstallCommand_Structure(t *testing.T) {
	dum, _, _, _ := newTestDumForInstall(t)
	cmd := newTestInstallCmd(t, dum)

	assert.Equal(t, "install", cmd.Use)
	assert.Equal(t, []string{"i"}, cmd.Aliases)

	assert.NotNil(t, cmd.Flags().Lookup("verbose"))
	assert.NotNil(t, cmd.Flags().Lookup("file"))
	assert.NotNil(t, cmd.Flags().Lookup("dryrun"))
	assert.NotNil(t, cmd.Flags().Lookup("group"))
}

func TestRunInstall_Success(t *testing.T) {
	dum, logger, factory, executor := newTestDumForInstall(t)
	cmd := newTestInstallCmd(t, dum)
	cmd.SetContext(context.Background())

	assert.NoError(t, cmd.Flags().Set("file", "/test/path"))
	assert.NoError(t, cmd.Flags().Set("dryrun", "false"))
	assert.NoError(t, cmd.Flags().Set("group", ""))

	logger.On("SetLevel", mock.Anything).Return()
	logger.On("Debug",
		"Running install command:",
		DRYRUN, false,
		VERBOSE, uint8(0),
		GROUP, "",
		FILE, "/test/path",
	).Return()

	expectedOpts := fy.InputOptions{
		File:   "/test/path",
		Group:  "",
		DryRun: false,
	}
	factory.On("Provide", expectedOpts).Return(&pb.Input{}, nil)

	executor.On("Install", mock.Anything, mock.Anything).Return(&pb.Result{}, nil)

	err := dum.runInstall(cmd)

	assert.NoError(t, err)
	factory.AssertExpectations(t)
	executor.AssertExpectations(t)
	logger.AssertExpectations(t)
}

func TestRunInstall_NilContext(t *testing.T) {
	dum, _, _, _ := newTestDumForInstall(t)
	cmd := &cobra.Command{}

	err := dum.runInstall(cmd)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "nil context")
}

func TestRunInstall_FactoryError(t *testing.T) {
	dum, logger, factory, _ := newTestDumForInstall(t)
	cmd := newTestInstallCmd(t, dum)
	cmd.SetContext(context.Background())

	assert.NoError(t, cmd.Flags().Set("file", "/test/path"))
	assert.NoError(t, cmd.Flags().Set("dryrun", "false"))
	assert.NoError(t, cmd.Flags().Set("group", ""))

	logger.On("SetLevel", mock.Anything).Return()
	logger.On("Debug",
		"Running install command:",
		DRYRUN, false,
		VERBOSE, uint8(0),
		GROUP, "",
		FILE, "/test/path",
	).Return()

	factory.On("Provide", mock.Anything).Return(nil, fmt.Errorf("file not found"))

	err := dum.runInstall(cmd)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "error providing context from file /test/path")
	assert.Contains(t, err.Error(), "file not found")
	factory.AssertExpectations(t)
}

func TestRunInstall_ExecutorError(t *testing.T) {
	dum, logger, factory, executor := newTestDumForInstall(t)
	cmd := newTestInstallCmd(t, dum)
	cmd.SetContext(context.Background())

	assert.NoError(t, cmd.Flags().Set("file", "/test/path"))
	assert.NoError(t, cmd.Flags().Set("dryrun", "false"))
	assert.NoError(t, cmd.Flags().Set("group", ""))

	logger.On("SetLevel", mock.Anything).Return()
	logger.On("Debug",
		"Running install command:",
		DRYRUN, false,
		VERBOSE, uint8(0),
		GROUP, "",
		FILE, "/test/path",
	).Return()

	factory.On("Provide", mock.Anything).Return(&pb.Input{}, nil)
	executor.On("Install", mock.Anything, mock.Anything).Return(nil, fmt.Errorf("install failed"))

	err := dum.runInstall(cmd)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "error installing config file /test/path")
	assert.Contains(t, err.Error(), "install failed")
	factory.AssertExpectations(t)
	executor.AssertExpectations(t)
}

func TestRunInstall_ViaCobra(t *testing.T) {
	t.Setenv("DUM_CONFIG", "/env/default.yml")
	dum, logger, factory, executor := newTestDumForInstall(t)
	cmd := newTestInstallCmd(t, dum)
	cmd.SetContext(context.Background())

	logger.On("SetLevel", mock.Anything).Return()
	logger.On("Debug",
		"Running install command:",
		DRYRUN, false,
		VERBOSE, uint8(0),
		GROUP, "",
		FILE, "/test/path",
	).Return()

	factory.On("Provide", fy.InputOptions{
		File:   "/test/path",
		Group:  "",
		DryRun: false,
	}).Return(&pb.Input{}, nil)
	executor.On("Install", mock.Anything, mock.Anything).Return(&pb.Result{}, nil)

	cmd.SetArgs([]string{"--file", "/test/path"})
	err := cmd.Execute()

	assert.NoError(t, err)
	factory.AssertExpectations(t)
	executor.AssertExpectations(t)
}
