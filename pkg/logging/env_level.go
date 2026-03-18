package logging

import (
	"os"
	"strconv"

	"github.com/charmbracelet/log"
)

// EnvLevel infers log level from an environment variable.
func EnvLevel() log.Level {
	envLevel := os.Getenv("ZSH_LOG_LEVEL")
	defaultLevel := log.WarnLevel
	if envLevel == "" {
		return defaultLevel
	}
	switch envLevel {
	case "debug":
		return log.DebugLevel
	case "info":
		return log.InfoLevel
	case "warn":
		return log.WarnLevel
	case "error":
		return log.ErrorLevel
	case "fatal":
		return log.FatalLevel
	default:
		envInt, err := strconv.Atoi(envLevel)
		if err != nil {
			return defaultLevel
		}
		if envInt <= 10 {
			return log.DebugLevel
		}
		if envInt <= 20 {
			return log.InfoLevel
		}
		if envInt <= 30 {
			return log.WarnLevel
		}
		if envInt <= 40 {
			return log.ErrorLevel
		}
		if envInt <= 50 {
			return log.FatalLevel
		}
		return defaultLevel
	}
}
