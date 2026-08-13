package cli

import (
	"strconv"

	clog "github.com/charmbracelet/log"
	"github.com/spf13/viper"
)

// GetLogLevel resolves ZSH_LOG_LEVEL using a fresh Viper instance.
func GetLogLevel() clog.Level {
	config := viper.New()
	_ = config.BindEnv("ZSH_LOG_LEVEL")
	return parseLogLevel(config.GetString("ZSH_LOG_LEVEL"))
}

func parseLogLevel(value string) clog.Level {
	switch value {
	case "debug":
		return clog.DebugLevel
	case "info":
		return clog.InfoLevel
	case "warn":
		return clog.WarnLevel
	case "error":
		return clog.ErrorLevel
	case "fatal":
		return clog.FatalLevel
	}

	valueInt, err := strconv.Atoi(value)
	if err != nil {
		return clog.WarnLevel
	}
	if valueInt <= 10 {
		return clog.DebugLevel
	}
	if valueInt <= 20 {
		return clog.InfoLevel
	}
	if valueInt <= 30 {
		return clog.WarnLevel
	}
	if valueInt <= 40 {
		return clog.ErrorLevel
	}
	if valueInt <= 50 {
		return clog.FatalLevel
	}
	return clog.WarnLevel
}
