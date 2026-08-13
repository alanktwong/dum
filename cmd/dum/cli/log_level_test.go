package cli

import (
	"testing"

	clog "github.com/charmbracelet/log"
	"github.com/stretchr/testify/assert"
)

func TestGetLogLevel(t *testing.T) {
	tests := []struct {
		name  string
		value string
		want  clog.Level
	}{
		{name: "empty", value: "", want: clog.WarnLevel},
		{name: "invalid", value: "invalid", want: clog.WarnLevel},
		{name: "debug", value: "debug", want: clog.DebugLevel},
		{name: "info", value: "info", want: clog.InfoLevel},
		{name: "warn", value: "warn", want: clog.WarnLevel},
		{name: "error", value: "error", want: clog.ErrorLevel},
		{name: "fatal", value: "fatal", want: clog.FatalLevel},
		{name: "negative", value: "-1", want: clog.DebugLevel},
		{name: "debug upper boundary", value: "10", want: clog.DebugLevel},
		{name: "info lower boundary", value: "11", want: clog.InfoLevel},
		{name: "info upper boundary", value: "20", want: clog.InfoLevel},
		{name: "warn lower boundary", value: "21", want: clog.WarnLevel},
		{name: "warn upper boundary", value: "30", want: clog.WarnLevel},
		{name: "error lower boundary", value: "31", want: clog.ErrorLevel},
		{name: "error upper boundary", value: "40", want: clog.ErrorLevel},
		{name: "fatal lower boundary", value: "41", want: clog.FatalLevel},
		{name: "fatal upper boundary", value: "50", want: clog.FatalLevel},
		{name: "over 50", value: "51", want: clog.WarnLevel},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Setenv("ZSH_LOG_LEVEL", tt.value)
			assert.Equal(t, tt.want, GetLogLevel())
		})
	}
}

func TestGetLogLevel_RepeatedCallsAreEnvironmentIsolated(t *testing.T) {
	t.Setenv("ZSH_LOG_LEVEL", "debug")
	assert.Equal(t, clog.DebugLevel, GetLogLevel())

	t.Setenv("ZSH_LOG_LEVEL", "warn")
	assert.Equal(t, clog.WarnLevel, GetLogLevel())
}
