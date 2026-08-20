// Package logging provides a logging abstraction to wrap charmbracelet/log.
package logging

import (
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/alanktwong/dum/internal/logging/gen"

	"github.com/charmbracelet/log"
)

// Logger is the interface for CLIs logging.
type Logger = gen.Logger

// Note charmbracelet/log uses atomic.Pointer, sync.Once and sync.RWMutex.
var (
	// logger lock provides a mutex for the logger singleton.
	loggerLock = sync.Once{}
	// loggerSingleton is a singleton of the logger.
	loggerSingleton Logger
)

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
		return fmt.Errorf("failed to println: %w", err)
	}
	return nil
}
