package playbook

import (
	"context"
	"errors"
	"fmt"
	"testing"

	l "awong/dotfiles/internal/logging"
	pl "awong/dotfiles/internal/plays"
	ty "awong/dotfiles/internal/types"

	"github.com/charmbracelet/log"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockPlayExecutor implements PlayExec for testing.
type MockPlayExecutor struct {
	mock.Mock
}

func (m *MockPlayExecutor) Initialize(input *pl.PlayInput) (*pl.PlayResult, error) {
	ret := m.Called(input)
	var r0 *pl.PlayResult
	if ret.Get(0) != nil {
		r0 = ret.Get(0).(*pl.PlayResult)
	}
	return r0, ret.Error(1)
}

func (m *MockPlayExecutor) InstallPlay(ctx context.Context, p *pl.Play, input *pl.PlayInput) (*pl.PlayResult, error) {
	ret := m.Called(ctx, p, input)
	var r0 *pl.PlayResult
	if ret.Get(0) != nil {
		r0 = ret.Get(0).(*pl.PlayResult)
	}
	return r0, ret.Error(1)
}

// MockPlayBookInfo implements plays.PlayBookInfo for testing.
type MockPlayBookInfo struct {
	mock.Mock
}

func (m *MockPlayBookInfo) GetID() string {
	ret := m.Called()
	return ret.Get(0).(string)
}

func (m *MockPlayBookInfo) GetJetBrainsApps() map[string]string {
	ret := m.Called()
	if ret.Get(0) == nil {
		return nil
	}
	return ret.Get(0).(map[string]string)
}

// mockLogger satisfies logging.Logger for testing (only implements used methods).
type mockLogger struct {
	mock.Mock
}

func (m *mockLogger) Printlnf(format string, a ...any) error {
	args := []any{format}
	args = append(args, a...)
	ret := m.Called(args...)
	return ret.Error(0)
}

func (m *mockLogger) Infof(format string, args ...any) {
	callArgs := []any{format}
	callArgs = append(callArgs, args...)
	m.Called(callArgs...)
}

func (m *mockLogger) Debug(msg any, keyvals ...any)     { panic("not implemented") }
func (m *mockLogger) Info(msg any, keyvals ...any)      { panic("not implemented") }
func (m *mockLogger) Warn(msg any, keyvals ...any)      { panic("not implemented") }
func (m *mockLogger) Error(msg any, keyvals ...any)     { panic("not implemented") }
func (m *mockLogger) Fatal(msg any, keyvals ...any)     { panic("not implemented") }
func (m *mockLogger) Debugf(format string, args ...any) { panic("not implemented") }
func (m *mockLogger) Warnf(format string, args ...any)  { panic("not implemented") }
func (m *mockLogger) Errorf(format string, args ...any) { panic("not implemented") }
func (m *mockLogger) Fatalf(format string, args ...any) { panic("not implemented") }
func (m *mockLogger) SetLevel(level log.Level)          { panic("not implemented") }
func (m *mockLogger) WithPrefix(prefix string) l.Logger { panic("not implemented") }
func (m *mockLogger) GetPrefix() string                 { panic("not implemented") }

func newTestExecutor(t *testing.T) (*Executor, *mockLogger, *MockPlayExecutor) {
	t.Helper()
	logger := &mockLogger{}
	playExec := &MockPlayExecutor{}
	return &Executor{
		Log:          logger,
		PlayExecutor: playExec,
	}, logger, playExec
}

func newTestPlayBook(id string, playNames ...string) *PlayBook {
	var plays []*pl.Play
	for _, name := range playNames {
		plays = append(plays, &pl.Play{Attributes: ty.Attributes{ID: name, Enabled: true}})
	}
	pb := &PlayBook{
		Attributes: ty.Attributes{ID: id, Description: "test", Enabled: true},
		Plays:      plays,
	}
	return pb
}

func TestNewExecutor_NoArgs(t *testing.T) {
	e := NewExecutor()
	assert.NotNil(t, e)
	assert.NotNil(t, e.Log)
	assert.NotNil(t, e.Ext)
	assert.NotNil(t, e.PlayExecutor)
}

func TestExecutor_Install_GroupPlay_Success(t *testing.T) {
	e, logger, playExec := newTestExecutor(t)
	ctx := context.Background()

	pb := newTestPlayBook("test-book", "play-1")
	input := &Input{Play: "play-1", PlayBook: pb}

	logger.On("Printlnf", "...Installing playbook (%v) ... %v", "test-book", "test").Return(nil)
	logger.On("Infof", "...END: installing playbook (%v) ... %v", "test-book", "test").Return(nil)

	initResult := &pl.PlayResult{PlayBook: "test-book", Play: "initialize", Success: true}
	playExec.On("Initialize", mock.AnythingOfType("*plays.PlayInput")).Return(initResult, nil)

	playResult := &pl.PlayResult{PlayBook: "test-book", Play: "play-1", Success: true}
	playExec.On(
		"InstallPlay",
		ctx,
		mock.AnythingOfType("*plays.Play"),
		mock.AnythingOfType("*plays.PlayInput"),
	).Return(
		playResult,
		nil,
	)

	result, err := e.Install(ctx, input)
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.True(t, result.Success)
	assert.Equal(t, "test-book", result.PlayBook)
}

func TestExecutor_Install_AllPlays_Success(t *testing.T) {
	e, logger, playExec := newTestExecutor(t)
	ctx := context.Background()

	pb := newTestPlayBook("test-book", "play-1")
	input := &Input{Play: "", PlayBook: pb}

	logger.On("Printlnf", "...Installing playbook (%v) ... %v", "test-book", "test").Return(nil)
	logger.On("Infof", "%v %v", PlayEllipsis, "play-1").Return(nil)
	logger.On("Infof", "...END: installing playbook (%v) ... %v", "test-book", "test").Return(nil)

	initResult := &pl.PlayResult{PlayBook: "test-book", Play: "initialize", Success: true}
	playExec.On("Initialize", mock.AnythingOfType("*plays.PlayInput")).Return(initResult, nil)

	playResult := &pl.PlayResult{PlayBook: "test-book", Play: "play-1", Success: true}
	playExec.On(
		"InstallPlay",
		ctx,
		mock.AnythingOfType("*plays.Play"),
		mock.AnythingOfType("*plays.PlayInput"),
	).Return(
		playResult,
		nil,
	)

	result, err := e.Install(ctx, input)
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.True(t, result.Success)
}

func TestExecutor_Install_PrintlnfError(t *testing.T) {
	e, logger, _ := newTestExecutor(t)
	ctx := context.Background()

	pb := newTestPlayBook("test-book")
	input := &Input{PlayBook: pb}

	logger.On("Printlnf", "...Installing playbook (%v) ... %v", "test-book", "test").Return(errors.New("log failed"))

	result, err := e.Install(ctx, input)
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "failed to install playbook")
}

func TestExecutor_Install_InitializeError(t *testing.T) {
	e, logger, playExec := newTestExecutor(t)
	ctx := context.Background()

	pb := newTestPlayBook("test-book")
	input := &Input{PlayBook: pb}

	logger.On("Printlnf", "...Installing playbook (%v) ... %v", "test-book", "test").Return(nil)
	playExec.On("Initialize", mock.AnythingOfType("*plays.PlayInput")).Return(nil, fmt.Errorf("init failed"))

	result, err := e.Install(ctx, input)
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "failed to Initialize")
}

func TestExecutor_Install_InstallPlayError(t *testing.T) {
	e, logger, playExec := newTestExecutor(t)
	ctx := context.Background()

	pb := newTestPlayBook("test-book", "play-1")
	input := &Input{Play: "play-1", PlayBook: pb}

	logger.On("Printlnf", "...Installing playbook (%v) ... %v", "test-book", "test").Return(nil)

	initResult := &pl.PlayResult{PlayBook: "test-book", Play: "initialize", Success: true}
	playExec.On("Initialize", mock.AnythingOfType("*plays.PlayInput")).Return(initResult, nil)
	playExec.On(
		"InstallPlay",
		ctx,
		mock.AnythingOfType("*plays.Play"),
		mock.AnythingOfType("*plays.PlayInput"),
	).Return(
		nil,
		fmt.Errorf("install failed"),
	)

	result, err := e.Install(ctx, input)
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "failed to install play")
}

func TestExecutor_Install_DryRunIncludesDisabledPlays(t *testing.T) {
	e, logger, playExec := newTestExecutor(t)
	ctx := context.Background()

	pb := newTestPlayBook("test-book", "play-1")
	disabled := &pl.Play{Attributes: ty.Attributes{ID: "play-disabled", Enabled: false}}
	pb.Plays = append(pb.Plays, disabled)
	input := &Input{Play: "", PlayBook: pb, DryRun: true}

	logger.On("Printlnf", "...Installing playbook (%v) ... %v", "test-book", "test").Return(nil)
	logger.On("Infof", "%v %v", PlayEllipsis, mock.AnythingOfType("string")).Return(nil)
	logger.On("Infof", "...END: installing playbook (%v) ... %v", "test-book", "test").Return(nil)

	initResult := &pl.PlayResult{PlayBook: "test-book", Play: "initialize", Success: true}
	playExec.On("Initialize", mock.AnythingOfType("*plays.PlayInput")).Return(initResult, nil)

	playResult := &pl.PlayResult{PlayBook: "test-book", Success: true}
	playExec.On(
		"InstallPlay",
		ctx,
		mock.AnythingOfType("*plays.Play"),
		mock.AnythingOfType("*plays.PlayInput"),
	).Return(
		playResult,
		nil,
	)

	result, err := e.Install(ctx, input)
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.True(t, result.Success)
	playExec.AssertNumberOfCalls(t, "InstallPlay", 2)
}

func TestExecutor_Install_NonDryRunExcludesDisabledPlays(t *testing.T) {
	e, logger, playExec := newTestExecutor(t)
	ctx := context.Background()

	pb := newTestPlayBook("test-book", "play-1")
	disabled := &pl.Play{Attributes: ty.Attributes{ID: "play-disabled", Enabled: false}}
	pb.Plays = append(pb.Plays, disabled)
	input := &Input{Play: "", PlayBook: pb, DryRun: false}

	logger.On("Printlnf", "...Installing playbook (%v) ... %v", "test-book", "test").Return(nil)
	logger.On("Infof", "%v %v", PlayEllipsis, "play-1").Return(nil)
	logger.On("Infof", "...END: installing playbook (%v) ... %v", "test-book", "test").Return(nil)

	initResult := &pl.PlayResult{PlayBook: "test-book", Play: "initialize", Success: true}
	playExec.On("Initialize", mock.AnythingOfType("*plays.PlayInput")).Return(initResult, nil)

	playResult := &pl.PlayResult{PlayBook: "test-book", Play: "play-1", Success: true}
	playExec.On(
		"InstallPlay",
		ctx,
		mock.AnythingOfType("*plays.Play"),
		mock.AnythingOfType("*plays.PlayInput"),
	).Return(
		playResult,
		nil,
	)

	result, err := e.Install(ctx, input)
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.True(t, result.Success)
	playExec.AssertNumberOfCalls(t, "InstallPlay", 1)
}
