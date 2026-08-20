package main

import (
	"context"

	fy "awong/dotfiles/internal/factory"
	lg "awong/dotfiles/internal/logging"
	pb "awong/dotfiles/internal/playbook"

	clog "github.com/charmbracelet/log"
	"github.com/stretchr/testify/mock"
)

// mockFactoryProvider implements factoryProvider for testing.
type mockFactoryProvider struct {
	mock.Mock
}

func (m *mockFactoryProvider) Provide(opts fy.InputOptions) (*pb.Input, error) {
	ret := m.Called(opts)
	var r0 *pb.Input
	if ret.Get(0) != nil {
		r0 = ret.Get(0).(*pb.Input)
	}
	return r0, ret.Error(1)
}

// mockInstallExecutor implements installExecutor for testing.
type mockInstallExecutor struct {
	mock.Mock
}

func (m *mockInstallExecutor) Install(ctx context.Context, input *pb.Input) (*pb.Result, error) {
	ret := m.Called(ctx, input)
	var r0 *pb.Result
	if ret.Get(0) != nil {
		r0 = ret.Get(0).(*pb.Result)
	}
	return r0, ret.Error(1)
}

// mockListExecutor implements listExecutor for testing.
type mockListExecutor struct {
	mock.Mock
}

func (m *mockListExecutor) List(ctx context.Context, input *pb.Input) (*pb.Result, error) {
	ret := m.Called(ctx, input)
	var r0 *pb.Result
	if ret.Get(0) != nil {
		r0 = ret.Get(0).(*pb.Result)
	}
	return r0, ret.Error(1)
}

// mockLogger implements lg.Logger for testing (only implements used methods).
type mockLogger struct {
	mock.Mock
}

func (m *mockLogger) Debug(msg any, keyvals ...any) {
	args := []any{msg}
	args = append(args, keyvals...)
	m.Called(args...)
}

func (m *mockLogger) Info(msg any, keyvals ...any)  { panic("not implemented") }
func (m *mockLogger) Warn(msg any, keyvals ...any)  { panic("not implemented") }
func (m *mockLogger) Error(msg any, keyvals ...any) { panic("not implemented") }
func (m *mockLogger) Fatal(msg any, keyvals ...any) { panic("not implemented") }

func (m *mockLogger) Debugf(format string, args ...any) {
	callArgs := []any{format}
	callArgs = append(callArgs, args...)
	m.Called(callArgs...)
}

func (m *mockLogger) Infof(format string, args ...any) {
	callArgs := []any{format}
	callArgs = append(callArgs, args...)
	m.Called(callArgs...)
}

func (m *mockLogger) Warnf(format string, args ...any) {
	callArgs := []any{format}
	callArgs = append(callArgs, args...)
	m.Called(callArgs...)
}

func (m *mockLogger) Errorf(format string, args ...any) {
	callArgs := []any{format}
	callArgs = append(callArgs, args...)
	m.Called(callArgs...)
}

func (m *mockLogger) Fatalf(format string, args ...any) {
	callArgs := []any{format}
	callArgs = append(callArgs, args...)
	m.Called(callArgs...)
}

func (m *mockLogger) Printlnf(format string, a ...any) error {
	args := []any{format}
	args = append(args, a...)
	ret := m.Called(args...)
	return ret.Error(0)
}

func (m *mockLogger) SetLevel(level clog.Level) {
	m.Called(level)
}

func (m *mockLogger) WithPrefix(prefix string) lg.Logger {
	ret := m.Called(prefix)
	return ret.Get(0).(lg.Logger)
}

func (m *mockLogger) GetPrefix() string {
	ret := m.Called()
	return ret.Get(0).(string)
}
