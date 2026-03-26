package cmd

import (
	"context"
	"fmt"
	"testing"

	"github.com/charmbracelet/log"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	gen "awong/dotfiles/pkg/logging/gen"
)

func newTestDumWithCmd(t *testing.T) (*Dum, *gen.MockLogger, *cobra.Command) {
	t.Helper()
	logger := gen.NewMockLogger(t)
	cmd := &cobra.Command{Use: "test"}
	addPrefixFlag(cmd)
	dum := &Dum{Log: logger}
	return dum, logger, cmd
}

func TestExecuteLogging_Debug(t *testing.T) {
	dum, logger, cmd := newTestDumWithCmd(t)

	logger.On("SetLevel", mock.Anything).Return()
	logger.On("Debugf", "%v", mock.Anything).Return()

	err := dum.executeLogging(cmd, []string{"hello", "world"}, log.DebugLevel)

	assert.NoError(t, err)
}

func TestExecuteLogging_Info(t *testing.T) {
	dum, logger, cmd := newTestDumWithCmd(t)

	logger.On("SetLevel", mock.Anything).Return()
	logger.On("Infof", "%v", mock.Anything).Return()

	err := dum.executeLogging(cmd, []string{"hello"}, log.InfoLevel)

	assert.NoError(t, err)
}

func TestExecuteLogging_Warn(t *testing.T) {
	dum, logger, cmd := newTestDumWithCmd(t)

	logger.On("SetLevel", mock.Anything).Return()
	logger.On("Warnf", "%v", mock.Anything).Return()

	err := dum.executeLogging(cmd, []string{"warning"}, log.WarnLevel)

	assert.NoError(t, err)
}

func TestExecuteLogging_Error(t *testing.T) {
	dum, logger, cmd := newTestDumWithCmd(t)

	logger.On("SetLevel", mock.Anything).Return()
	logger.On("Errorf", "%v", mock.Anything).Return()

	err := dum.executeLogging(cmd, []string{"error"}, log.ErrorLevel)

	assert.NoError(t, err)
}

func TestExecuteLogging_Fatal(t *testing.T) {
	dum, logger, cmd := newTestDumWithCmd(t)

	logger.On("SetLevel", mock.Anything).Return()
	logger.On("Fatalf", "%v", mock.Anything).Return()

	err := dum.executeLogging(cmd, []string{"fatal"}, log.FatalLevel)

	assert.NoError(t, err)
}

func TestExecuteLogging_Success(t *testing.T) {
	dum, logger, cmd := newTestDumWithCmd(t)

	logger.On("SetLevel", mock.Anything).Return()
	logger.On("Printlnf", "%v", mock.Anything).Return(nil)

	err := dum.executeLogging(cmd, []string{"done"}, SuccessLevel)

	assert.NoError(t, err)
}

func TestExecuteLogging_SuccessWithPrefix(t *testing.T) {
	dum, logger, cmd := newTestDumWithCmd(t)

	err := cmd.Flags().Set(PREFIX, "my-prefix")
	assert.NoError(t, err)

	prefixLogger := gen.NewMockLogger(t)

	logger.On("SetLevel", mock.Anything).Return()
	logger.On("WithPrefix", "my-prefix").Return(prefixLogger)
	prefixLogger.On("Printlnf", "%v %v", mock.Anything).Return(nil)

	err = dum.executeLogging(cmd, []string{"done"}, SuccessLevel)

	assert.NoError(t, err)
}

func TestExecuteLogging_SuccessPrintlnfError(t *testing.T) {
	dum, logger, cmd := newTestDumWithCmd(t)

	logger.On("SetLevel", mock.Anything).Return()
	logger.On("Printlnf", "%v", mock.Anything).Return(fmt.Errorf("write error"))

	err := dum.executeLogging(cmd, []string{"done"}, SuccessLevel)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to println")
}

func TestExecuteLogging_InvalidLevel(t *testing.T) {
	dum, logger, cmd := newTestDumWithCmd(t)

	logger.On("SetLevel", mock.Anything).Return()

	err := dum.executeLogging(cmd, []string{"test"}, log.Level(99))

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid log level")
}

func TestNewDebugCommand_RunE(t *testing.T) {
	logger := gen.NewMockLogger(t)
	dum := &Dum{Log: logger}
	logCmd := NewLogCommand("dum", dum)
	debugCmd, _, err := logCmd.Find([]string{"debug"})
	assert.NoError(t, err)
	debugCmd.SetContext(context.Background())

	logger.On("SetLevel", mock.Anything).Return()
	logger.On("Debugf", "%v", mock.Anything).Return()

	err = debugCmd.RunE(debugCmd, []string{"hello", "world"})
	assert.NoError(t, err)
}

func TestNewInfoCommand_RunE(t *testing.T) {
	logger := gen.NewMockLogger(t)
	dum := &Dum{Log: logger}
	logCmd := NewLogCommand("dum", dum)
	infoCmd, _, err := logCmd.Find([]string{"info"})
	assert.NoError(t, err)
	infoCmd.SetContext(context.Background())

	logger.On("SetLevel", mock.Anything).Return()
	logger.On("Infof", "%v", mock.Anything).Return()

	err = infoCmd.RunE(infoCmd, []string{"hello", "world"})
	assert.NoError(t, err)
}

func TestNewWarnCommand_RunE(t *testing.T) {
	logger := gen.NewMockLogger(t)
	dum := &Dum{Log: logger}
	logCmd := NewLogCommand("dum", dum)
	warnCmd, _, err := logCmd.Find([]string{"warn"})
	assert.NoError(t, err)
	warnCmd.SetContext(context.Background())

	logger.On("SetLevel", mock.Anything).Return()
	logger.On("Warnf", "%v", mock.Anything).Return()

	err = warnCmd.RunE(warnCmd, []string{"hello", "world"})
	assert.NoError(t, err)
}

func TestNewErrorCommand_RunE(t *testing.T) {
	logger := gen.NewMockLogger(t)
	dum := &Dum{Log: logger}
	logCmd := NewLogCommand("dum", dum)
	errorCmd, _, err := logCmd.Find([]string{"error"})
	assert.NoError(t, err)
	errorCmd.SetContext(context.Background())

	logger.On("SetLevel", mock.Anything).Return()
	logger.On("Errorf", "%v", mock.Anything).Return()

	err = errorCmd.RunE(errorCmd, []string{"hello", "world"})
	assert.NoError(t, err)
}

func TestNewSuccessCommand_RunE(t *testing.T) {
	logger := gen.NewMockLogger(t)
	dum := &Dum{Log: logger}
	logCmd := NewLogCommand("dum", dum)
	successCmd, _, err := logCmd.Find([]string{"success"})
	assert.NoError(t, err)
	successCmd.SetContext(context.Background())

	logger.On("SetLevel", mock.Anything).Return()
	logger.On("Printlnf", "%v", mock.Anything).Return(nil)

	err = successCmd.RunE(successCmd, []string{"hello", "world"})
	assert.NoError(t, err)
}
