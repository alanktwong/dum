package cmd

import (
	"context"
	"fmt"
	"testing"

	f "awong/dotfiles/pkg/factory"
	pb "awong/dotfiles/pkg/playbook"

	"github.com/charmbracelet/log"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func newTestDumForList(t *testing.T) (*Dum, *mockLogger, *mockFactoryProvider, *mockListExecutor) {
	t.Helper()
	logger := &mockLogger{}
	factory := &mockFactoryProvider{}
	executor := &mockListExecutor{}
	dum := &Dum{
		Log:             logger,
		FactoryProvider: factory,
		ListExecutor:    executor,
	}
	return dum, logger, factory, executor
}

func TestNewListCommand_Structure(t *testing.T) {
	dum := &Dum{}
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

	logger.On("SetLevel", log.WarnLevel).Return()
	logger.On("Debug",
		"Running list command:",
		DRYRUN, false,
		VERBOSE, uint8(0),
		GROUP, "",
		FILE, "/test/path",
	).Return()

	input := &pb.Input{}
	factory.On("Provide", f.InputOptions{
		File:   "/test/path",
		Group:  "",
		DryRun: false,
	}).Return(input, nil)

	executor.On("List", mock.Anything, input).Return(&pb.Result{}, nil)

	cmd := &cobra.Command{Use: "list"}
	cmd.SetContext(context.Background())
	addVerboseFlag(cmd)
	addFileFlag(cmd)
	addGroupFlag(cmd)
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
	addVerboseFlag(cmd)
	addFileFlag(cmd)
	addGroupFlag(cmd)

	err := dum.runList(cmd)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "nil context")
}

func TestRunList_FactoryError(t *testing.T) {
	dum, logger, factory, _ := newTestDumForList(t)

	logger.On("SetLevel", log.WarnLevel).Return()
	logger.On("Debug",
		"Running list command:",
		DRYRUN, false,
		VERBOSE, uint8(0),
		GROUP, "",
		FILE, "/test/path",
	).Return()

	factory.On("Provide", f.InputOptions{
		File:   "/test/path",
		Group:  "",
		DryRun: false,
	}).Return(nil, fmt.Errorf("factory failed"))

	cmd := &cobra.Command{Use: "list"}
	cmd.SetContext(context.Background())
	addVerboseFlag(cmd)
	addFileFlag(cmd)
	addGroupFlag(cmd)
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

	logger.On("SetLevel", log.WarnLevel).Return()
	logger.On("Debug",
		"Running list command:",
		DRYRUN, false,
		VERBOSE, uint8(0),
		GROUP, "",
		FILE, "/test/path",
	).Return()

	input := &pb.Input{}
	factory.On("Provide", f.InputOptions{
		File:   "/test/path",
		Group:  "",
		DryRun: false,
	}).Return(input, nil)

	executor.On("List", mock.Anything, input).Return(nil, fmt.Errorf("list failed"))

	cmd := &cobra.Command{Use: "list"}
	cmd.SetContext(context.Background())
	addVerboseFlag(cmd)
	addFileFlag(cmd)
	addGroupFlag(cmd)
	_ = cmd.Flags().Set(FILE, "/test/path")

	err := dum.runList(cmd)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "error listing")
	assert.Contains(t, err.Error(), "list failed")

	factory.AssertExpectations(t)
	executor.AssertExpectations(t)
	logger.AssertExpectations(t)
}
