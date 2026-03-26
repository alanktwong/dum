// Package gen provides generated types for logging package.
package gen

import "github.com/charmbracelet/log"

// Logger is the interface for CLIs logging.
type Logger interface {
	Debug(msg any, keyvals ...any)
	Info(msg any, keyvals ...any)
	Warn(msg any, keyvals ...any)
	Error(msg any, keyvals ...any)
	Fatal(msg any, keyvals ...any)
	Debugf(format string, args ...any)
	Infof(format string, args ...any)
	Warnf(format string, args ...any)
	Errorf(format string, args ...any)
	Fatalf(format string, args ...any)
	Printlnf(format string, a ...any) error
	SetLevel(level log.Level)
	WithPrefix(prefix string) Logger
	GetPrefix() string
}
