// Package logging provides a logging abstraction to wrap charmbracelet/log.
package logging

import (
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/charmbracelet/log"
)

// Note charmbracelet/log uses atomic.Pointer, sync.Once and sync.RWMutex.
var (
	// logger lock provides a mutex for the logger singleton.
	loggerLock = sync.Once{}
	// loggerSingleton is a singleton of the logger.
	loggerSingleton Logger
)

// Logger is the interface for CLIs logging.
type Logger interface {
	// Debug prints a debug message.
	Debug(msg any, keyvals ...any)
	// Info prints an info message.
	Info(msg any, keyvals ...any)
	// Warn prints a warn message.
	Warn(msg any, keyvals ...any)
	// Error prints an error message.
	Error(msg any, keyvals ...any)
	// Fatal prints a fatal message.
	Fatal(msg any, keyvals ...any)
	// Debugf prints a debug message with formatting.
	Debugf(format string, args ...any)
	// Infof prints an info message with formatting.
	Infof(format string, args ...any)
	// Warnf prints a warn message with formatting.
	Warnf(format string, args ...any)
	// Errorf prints an error message with formatting.
	Errorf(format string, args ...any)
	// Fatalf prints a fatal message with formatting.
	Fatalf(format string, args ...any)
	// Printlnf prints a message with formatting.
	Printlnf(format string, a ...any) error
	// SetLevel sets the current level.
	SetLevel(level log.Level)
	// WithPrefix returns a new logger with the given prefix.
	WithPrefix(prefix string) Logger
	// GetPrefix returns the current prefix.
	GetPrefix() string
}

// LoggerImpl is the default implementation of Logger.
type LoggerImpl struct {
	Log *log.Logger
}

// Options is the struct for constructing LoggerImpl.
type Options struct {
	Prefix string
	Level  log.Level
}

// NewLogger is the factory method for constructing LoggerImpl.
func NewLogger(options Options) Logger {
	if loggerSingleton == nil {
		loggerLock.Do(func() {
			logger := log.NewWithOptions(os.Stderr, log.Options{
				ReportCaller:    false,
				ReportTimestamp: true,
				TimeFormat:      time.RFC1123Z,
			})
			logger.SetLevel(options.Level)
			if options.Prefix != "" {
				logger.SetPrefix(options.Prefix)
			}
			loggerSingleton = &LoggerImpl{
				logger,
			}
		})
	}
	return loggerSingleton
}

// Log provides the logger singleton.
func Log() Logger {
	if loggerSingleton == nil {
		return NewLogger(Options{
			Prefix: "",
			Level:  log.InfoLevel,
		})
	}
	return loggerSingleton
}

// Debug implements the Logger interface.
func (l *LoggerImpl) Debug(msg any, keyvals ...any) {
	l.Log.Debug(msg, keyvals...)
}

// Info implements the Logger interface.
func (l *LoggerImpl) Info(msg any, keyvals ...any) {
	l.Log.Info(msg, keyvals...)
}

// Warn implements the Logger interface.
func (l *LoggerImpl) Warn(msg any, keyvals ...any) {
	l.Log.Warn(msg, keyvals...)
}

// Error implements the Logger interface.
func (l *LoggerImpl) Error(msg any, keyvals ...any) {
	l.Log.Error(msg, keyvals...)
}

// Fatal implements the Logger interface.
func (l *LoggerImpl) Fatal(msg any, keyvals ...any) {
	l.Log.Fatal(msg, keyvals...)
}

// Debugf implements the Logger interface.
func (l *LoggerImpl) Debugf(format string, args ...any) {
	l.Log.Debugf(format, args...)
}

// Infof implements the Logger interface.
func (l *LoggerImpl) Infof(format string, args ...any) {
	l.Log.Infof(format, args...)
}

// Warnf implements the Logger interface.
func (l *LoggerImpl) Warnf(format string, args ...any) {
	l.Log.Warnf(format, args...)
}

// Errorf implements the Logger interface.
func (l *LoggerImpl) Errorf(format string, args ...any) {
	l.Log.Errorf(format, args...)
}

// Fatalf implements the Logger interface.
func (l *LoggerImpl) Fatalf(format string, args ...any) {
	l.Log.Fatalf(format, args...)
}

// SetLevel implements the Logger interface.
func (l *LoggerImpl) SetLevel(level log.Level) {
	l.Log.SetLevel(level)
}

// WithPrefix implements the Logger interface.
func (l *LoggerImpl) WithPrefix(prefix string) Logger {
	return &LoggerImpl{
		Log: l.Log.WithPrefix(prefix),
	}
}

// GetPrefix implements the Logger interface.
func (l *LoggerImpl) GetPrefix() string {
	return l.Log.GetPrefix()
}

// Printlnf implements the Logger interface.
func (l *LoggerImpl) Printlnf(format string, a ...any) error {
	formatted := fmt.Sprintf(format, a...)
	_, err := fmt.Println(formatted)
	if err != nil {
		return fmt.Errorf("failed to println: %v", err)
	}
	return nil
}
