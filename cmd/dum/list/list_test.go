package list

import (
	"awong/dotfiles/cmd/dum/cli"
	"context"
	"fmt"
	"testing"

	fy "awong/dotfiles/internal/factory"
	pb "awong/dotfiles/internal/playbook"

	clog "github.com/charmbracelet/log"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func newTestDumForList(t *testing.T) (*Command, *mockLogger, *mockFactoryProvider, *mockListExecutor) {
	t.Helper()
	logger := &mockLogger{}
	factory := &mockFactoryProvider{}
	executor := &mockListExecutor{}
	dum := &Command{
		Log:             logger,
		FactoryProvider: factory,
		Executor:        executor,
	}
	return dum, logger, factory, executor
}

func TestNewListCommand_Structure(t *testing.T) {
	dum := &Command{}
	cmd := NewListCommand("dum", dum)

	assert.Equal(t, "list", cmd.Use)
	assert.Contains(t, cmd.Aliases, "ls")

	assert.NotNil(t, cmd.Flags().Lookup("verbose"))
	assert.NotNil(t, cmd.Flags().Lookup("file"))
	assert.NotNil(t, cmd.Flags().Lookup("group"))
	assert.Nil(t, cmd.Flags().Lookup("dryrun"))
}

func TestRunList_Success(t *testing.T) {
	dum, logger, factory, executor := newTestDumForList(t)

	logger.On("SetLevel", clog.WarnLevel).Return()
	logger.On("Debug",
		"Running list command:",
		DRYRUN, false,
		VERBOSE, uint8(0),
		GROUP, "",
		FILE, "/test/path",
	).Return()

	input := &pb.Input{}
	factory.On("Provide", fy.InputOptions{
		File:   "/test/path",
		Group:  "",
		DryRun: false,
	}).Return(input, nil)

	executor.On("List", mock.Anything, input).Return(&pb.Result{}, nil)

	cmd := &cobra.Command{Use: "list"}
	cmd.SetContext(context.Background())
	cli.AddVerboseFlag(cmd)
	cli.AddFileFlag(cmd)
	cli.AddGroupFlag(cmd)
	_ = cmd.Flags().Set(FILE, "/test/path")

	err := dum.runList(cmd)
	assert.NoError(t, err)

	factory.AssertExpectations(t)
	executor.AssertExpectations(t)
	logger.AssertExpectations(t)
}

func TestRunList_NilContext(t *testing.T) {
	dum, _, _, _ := newTestDumForList(t)

	cmd := &cobra.Command{Use: "list"}
	cli.AddVerboseFlag(cmd)
	cli.AddFileFlag(cmd)
	cli.AddGroupFlag(cmd)

	err := dum.runList(cmd)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "nil context")
}

func TestRunList_FactoryError(t *testing.T) {
	dum, logger, factory, _ := newTestDumForList(t)

	logger.On("SetLevel", clog.WarnLevel).Return()
	logger.On("Debug",
		"Running list command:",
		DRYRUN, false,
		VERBOSE, uint8(0),
		GROUP, "",
		FILE, "/test/path",
	).Return()

	factory.On("Provide", fy.InputOptions{
		File:   "/test/path",
		Group:  "",
		DryRun: false,
	}).Return(nil, fmt.Errorf("factory failed"))

	cmd := &cobra.Command{Use: "list"}
	cmd.SetContext(context.Background())
	cli.AddVerboseFlag(cmd)
	cli.AddFileFlag(cmd)
	cli.AddGroupFlag(cmd)
	_ = cmd.Flags().Set(FILE, "/test/path")

	err := dum.runList(cmd)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "error providing playbook")
	assert.Contains(t, err.Error(), "factory failed")

	factory.AssertExpectations(t)
	logger.AssertExpectations(t)
}

func TestRunList_ExecutorError(t *testing.T) {
	dum, logger, factory, executor := newTestDumForList(t)

	logger.On("SetLevel", clog.WarnLevel).Return()
	logger.On("Debug",
		"Running list command:",
		DRYRUN, false,
		VERBOSE, uint8(0),
		GROUP, "",
		FILE, "/test/path",
	).Return()

	input := &pb.Input{}
	factory.On("Provide", fy.InputOptions{
		File:   "/test/path",
		Group:  "",
		DryRun: false,
	}).Return(input, nil)

	executor.On("List", mock.Anything, input).Return(nil, fmt.Errorf("list failed"))

	cmd := &cobra.Command{Use: "list"}
	cmd.SetContext(context.Background())
	cli.AddVerboseFlag(cmd)
	cli.AddFileFlag(cmd)
	cli.AddGroupFlag(cmd)
	_ = cmd.Flags().Set(FILE, "/test/path")

	err := dum.runList(cmd)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "error listing")
	assert.Contains(t, err.Error(), "list failed")

	factory.AssertExpectations(t)
	executor.AssertExpectations(t)
	logger.AssertExpectations(t)
}

func TestRunList_ViaCobraExplicitFileReachesFactory(t *testing.T) {
	t.Setenv("INSTALLER_CONFIG", "/env/default.yml")

	dum, logger, factory, executor := newTestDumForList(t)
	cmd := NewListCommand("test", dum)
	cmd.SetContext(context.Background())

	logger.On("SetLevel", clog.WarnLevel).Return()
	logger.On("Debug",
		"Running list command:",
		DRYRUN, false,
		VERBOSE, uint8(0),
		GROUP, "",
		FILE, "/explicit/list.yml",
	).Return()
	factory.On("Provide", fy.InputOptions{File: "/explicit/list.yml", Group: ""}).Return(&pb.Input{}, nil)
	executor.On("List", mock.Anything, mock.Anything).Return(&pb.Result{}, nil)

	cmd.SetArgs([]string{"--file", "/explicit/list.yml"})
	err := cmd.Execute()

	assert.NoError(t, err)
	factory.AssertExpectations(t)
	executor.AssertExpectations(t)
	logger.AssertExpectations(t)
}
