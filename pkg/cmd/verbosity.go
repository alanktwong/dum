package cmd

import (
	"math"

	"github.com/charmbracelet/log"
	"github.com/spf13/cobra"
)

// Verbosity is a struct for expressing the verbosity of a CLI flag.
type Verbosity struct {
	Verbose uint8
}

// GetVerbosityFromCommand infers the verbosity level from a Cobra command
// in order to set log levels.
func GetVerbosityFromCommand(cmd *cobra.Command) Verbosity {
	count, err := cmd.Flags().GetCount(VERBOSE)
	if err != nil {
		count = 0
	}
	if count >= 0 && count < math.MaxUint8 {
		return Verbosity{Verbose: uint8(count)}
	}
	return Verbosity{Verbose: 0}
}

// Level converts verbosity to a log level.
func (v *Verbosity) Level() log.Level {
	if v.Verbose >= 2 {
		return log.DebugLevel
	} else if v.Verbose >= 1 {
		return log.InfoLevel
	}
	return log.WarnLevel
}
