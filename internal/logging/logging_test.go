package logging

import (
	"bytes"
	"sync"
	"testing"

	"github.com/charmbracelet/log"
	"github.com/stretchr/testify/assert"

	gen "alanktwong/dum/internal/logging/gen"
)

// resetSingleton resets the logger singleton for testing.
func resetSingleton() {
	loggerSingleton = nil
	loggerLock = sync.Once{}
}

func TestNewLogger(t *testing.T) {
	resetSingleton()
	options := Options{
		Prefix: "test",
		Level:  log.InfoLevel,
	}
	logger := NewLogger(options)
	loggerImpl, ok := logger.(*LoggerImpl)
	assert.True(t, ok)
	assert.NotNil(t, logger)
	assert.NotNil(t, loggerImpl.Log)
	assert.Equal(t, "test", loggerImpl.GetPrefix())
}

func TestNewLogger_NoPrefix(t *testing.T) {
	resetSingleton()
	options := Options{
		Prefix: "",
		Level:  log.DebugLevel,
	}
	logger := NewLogger(options)
	assert.NotNil(t, logger)
}

func TestNewLogger_Singleton(t *testing.T) {
	resetSingleton()
	options := Options{
		Prefix: "first",
		Level:  log.InfoLevel,
	}
	first := NewLogger(options)
	second := NewLogger(Options{Prefix: "second", Level: log.DebugLevel})
	assert.Same(t, first, second)
}

func TestLog_CreatesSingleton(t *testing.T) {
	resetSingleton()
	logger := Log()
	assert.NotNil(t, logger)
	loggerImpl, ok := logger.(*LoggerImpl)
	assert.True(t, ok)
	assert.NotNil(t, loggerImpl.Log)
}

func TestLog_ReturnsExistingSingleton(t *testing.T) {
	resetSingleton()
	first := NewLogger(Options{Prefix: "test", Level: log.InfoLevel})
	second := Log()
	assert.Same(t, first, second)
}

func TestLoggerImpl_Debug(t *testing.T) {
	var buf bytes.Buffer
	l := log.New(&buf)
	l.SetLevel(log.DebugLevel)
	logger := &LoggerImpl{Log: l}

	logger.Debug("debug message")
	assert.Contains(t, buf.String(), "debug message")
}

func TestLoggerImpl_Debugf(t *testing.T) {
	var buf bytes.Buffer
	l := log.New(&buf)
	l.SetLevel(log.DebugLevel)
	logger := &LoggerImpl{Log: l}

	logger.Debugf("debug %s", "formatted")
	assert.Contains(t, buf.String(), "debug formatted")
}

func TestLoggerImpl_Info(t *testing.T) {
	var buf bytes.Buffer
	l := log.New(&buf)
	l.SetLevel(log.InfoLevel)
	logger := &LoggerImpl{Log: l}

	logger.Info("info message")
	assert.Contains(t, buf.String(), "info message")
}

func TestLoggerImpl_Infof(t *testing.T) {
	var buf bytes.Buffer
	l := log.New(&buf)
	l.SetLevel(log.InfoLevel)
	logger := &LoggerImpl{Log: l}

	logger.Infof("info %d", 42)
	assert.Contains(t, buf.String(), "info 42")
}

func TestLoggerImpl_Warn(t *testing.T) {
	var buf bytes.Buffer
	l := log.New(&buf)
	l.SetLevel(log.WarnLevel)
	logger := &LoggerImpl{Log: l}

	logger.Warn("warn message")
	assert.Contains(t, buf.String(), "warn message")
}

func TestLoggerImpl_Warnf(t *testing.T) {
	var buf bytes.Buffer
	l := log.New(&buf)
	l.SetLevel(log.WarnLevel)
	logger := &LoggerImpl{Log: l}

	logger.Warnf("warn %s", "test")
	assert.Contains(t, buf.String(), "warn test")
}

func TestLoggerImpl_Error(t *testing.T) {
	var buf bytes.Buffer
	l := log.New(&buf)
	l.SetLevel(log.ErrorLevel)
	logger := &LoggerImpl{Log: l}

	logger.Error("error message")
	assert.Contains(t, buf.String(), "error message")
}

func TestLoggerImpl_Errorf(t *testing.T) {
	var buf bytes.Buffer
	l := log.New(&buf)
	l.SetLevel(log.ErrorLevel)
	logger := &LoggerImpl{Log: l}

	logger.Errorf("error %s", "occurred")
	assert.Contains(t, buf.String(), "error occurred")
}

func TestLoggerImpl_SetLevel(t *testing.T) {
	var buf bytes.Buffer
	l := log.New(&buf)
	l.SetLevel(log.InfoLevel)
	logger := &LoggerImpl{Log: l}

	logger.SetLevel(log.DebugLevel)
	logger.Debug("should appear")
	assert.Contains(t, buf.String(), "should appear")
}

func TestLoggerImpl_WithPrefix(t *testing.T) {
	var buf bytes.Buffer
	l := log.New(&buf)
	l.SetLevel(log.InfoLevel)
	logger := &LoggerImpl{Log: l}

	child := logger.WithPrefix("child")
	assert.NotNil(t, child)
	assert.Equal(t, "child", child.GetPrefix())
}

func TestLoggerImpl_GetPrefix(t *testing.T) {
	var buf bytes.Buffer
	l := log.New(&buf)
	l.SetPrefix("myprefix")
	logger := &LoggerImpl{Log: l}

	assert.Equal(t, "myprefix", logger.GetPrefix())
}

func TestLoggerImpl_Printlnf(t *testing.T) {
	var buf bytes.Buffer
	l := log.New(&buf)
	logger := &LoggerImpl{Log: l}

	err := logger.Printlnf("hello %s", "world")
	assert.NoError(t, err)
}

func TestLoggerImpl_Printlnf_NoArgs(t *testing.T) {
	var buf bytes.Buffer
	l := log.New(&buf)
	logger := &LoggerImpl{Log: l}

	err := logger.Printlnf("hello")
	assert.NoError(t, err)
}

func TestMockLogger_Debug(t *testing.T) {
	m := gen.NewMockLogger(t)
	m.EXPECT().Debug("msg", []any{"key", "val"}).Return()
	m.Debug("msg", "key", "val")
}

func TestMockLogger_Info(t *testing.T) {
	m := gen.NewMockLogger(t)
	m.EXPECT().Info("msg").Return()
	m.Info("msg")
}

func TestMockLogger_Warn(t *testing.T) {
	m := gen.NewMockLogger(t)
	m.EXPECT().Warn("msg").Return()
	m.Warn("msg")
}

func TestMockLogger_Error(t *testing.T) {
	m := gen.NewMockLogger(t)
	m.EXPECT().Error("msg").Return()
	m.Error("msg")
}

func TestMockLogger_Fatal(t *testing.T) {
	m := gen.NewMockLogger(t)
	m.EXPECT().Fatal("msg").Return()
	m.Fatal("msg")
}

func TestMockLogger_Debugf(t *testing.T) {
	m := gen.NewMockLogger(t)
	m.EXPECT().Debugf("format %s", []any{"arg"}).Return()
	m.Debugf("format %s", "arg")
}

func TestMockLogger_Infof(t *testing.T) {
	m := gen.NewMockLogger(t)
	m.EXPECT().Infof("format %s", []any{"arg"}).Return()
	m.Infof("format %s", "arg")
}

func TestMockLogger_Warnf(t *testing.T) {
	m := gen.NewMockLogger(t)
	m.EXPECT().Warnf("format %s", []any{"arg"}).Return()
	m.Warnf("format %s", "arg")
}

func TestMockLogger_Errorf(t *testing.T) {
	m := gen.NewMockLogger(t)
	m.EXPECT().Errorf("format %s", []any{"arg"}).Return()
	m.Errorf("format %s", "arg")
}

func TestMockLogger_Fatalf(t *testing.T) {
	m := gen.NewMockLogger(t)
	m.EXPECT().Fatalf("format %s", []any{"arg"}).Return()
	m.Fatalf("format %s", "arg")
}

func TestMockLogger_Printlnf(t *testing.T) {
	m := gen.NewMockLogger(t)
	m.EXPECT().Printlnf("hello %s", []any{"world"}).Return(nil)
	err := m.Printlnf("hello %s", "world")
	assert.NoError(t, err)
}

func TestMockLogger_SetLevel(t *testing.T) {
	m := gen.NewMockLogger(t)
	m.EXPECT().SetLevel(log.DebugLevel).Return()
	m.SetLevel(log.DebugLevel)
}

func TestMockLogger_WithPrefix(t *testing.T) {
	m := gen.NewMockLogger(t)
	child := gen.NewMockLogger(t)
	m.EXPECT().WithPrefix("child").Return(child)
	result := m.WithPrefix("child")
	assert.Same(t, child, result)
}

func TestMockLogger_GetPrefix(t *testing.T) {
	m := gen.NewMockLogger(t)
	m.EXPECT().GetPrefix().Return("myprefix")
	assert.Equal(t, "myprefix", m.GetPrefix())
}
