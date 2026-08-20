// Shared Cobra flag and verbosity helpers for the dum CLI.
package main

import (
	"fmt"
	"math"
	"os"

	clog "github.com/charmbracelet/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// CLI flag names and shared logging defaults.
const (
	DRYRUN                  = "dryrun"
	FILE                    = "file"
	GROUP                   = "group"
	VERBOSE                 = "verbose"
	REPLACE                 = "replace"
	SOURCE                  = "source"
	LOWERCASE               = "lowercase"
	TRIMFRONT               = "trimfront"
	TRIMBACK                = "trimback"
	ITERATE                 = "iterate"
	OUTPUT                  = "output"
	PREFIX                  = "prefix"
	SuccessLevel clog.Level = math.MaxInt
	VerboseUsage            = "increase output verbosity. E.g. --vv is debug and -v is info"
)

// Verbosity represents the repeated --verbose flag.
type Verbosity struct{ Verbose uint8 }

// GetVerbosityFromCommand reads the command's verbosity count.
func GetVerbosityFromCommand(cmd *cobra.Command) Verbosity {
	count, err := cmd.Flags().GetCount(VERBOSE)
	if err != nil || count < 0 || count >= math.MaxUint8 {
		count = 0
	}
	return Verbosity{Verbose: uint8(count)}
}

// Level converts verbosity to a logging level.
func (v *Verbosity) Level() clog.Level {
	if v.Verbose >= 2 {
		return clog.DebugLevel
	}
	if v.Verbose >= 1 {
		return clog.InfoLevel
	}
	return clog.WarnLevel
}

// AddGroupFlag adds the --group selector.
func AddGroupFlag(cmd *cobra.Command) { cmd.Flags().StringP(GROUP, "g", "", "select GROUP to run") }

// AddFileFlag adds the --file selector.
func AddFileFlag(cmd *cobra.Command) {
	cmd.Flags().StringP(FILE, "f", GetDefaultConfig(), "select config FILE to use")
}

// AddVerboseFlag adds the repeatable --verbose flag.
func AddVerboseFlag(cmd *cobra.Command) { cmd.Flags().CountP(VERBOSE, "v", VerboseUsage) }

// AddDryRunFlag adds the --dryrun flag.
func AddDryRunFlag(cmd *cobra.Command) {
	cmd.Flags().BoolP(DRYRUN, "d", false, fmt.Sprintf("dry run of %s", cmd.Use))
}

// AddPrefixFlag adds the --prefix flag.
func AddPrefixFlag(cmd *cobra.Command) {
	cmd.Flags().StringP(PREFIX, "p", "", "add a prefix to the logger")
}

// AddOutputFlag adds the --output flag.
func AddOutputFlag(cmd *cobra.Command) {
	cmd.Flags().StringP(OUTPUT, "o", "", "output file path (default: stdout)")
}

// GetDefaultConfig returns the configured dum configuration file path.
//
// Resolution order, first match wins:
//  1. DUM_CONFIG environment variable (returned verbatim).
//  2. ./dum.yml if it exists.
//  3. $XDG_CONFIG_HOME/dum/dum.yml if it exists.
//  4. Legacy $XDG_CONFIG_HOME/dum/installer.yml if it exists.
//  5. The canonical $XDG_CONFIG_HOME/dum/dum.yml path otherwise.
func GetDefaultConfig() string {
	config := viper.New()
	_ = config.BindEnv("DUM_CONFIG")
	_ = config.BindEnv("XDG_CONFIG_HOME")

	if path := config.GetString("DUM_CONFIG"); path != "" {
		return path
	}

	base := config.GetString("XDG_CONFIG_HOME")
	if base == "" {
		base = "~/.config"
	}
	candidates := []string{
		"dum.yml",
		fmt.Sprintf("%s/dum/dum.yml", base),
		fmt.Sprintf("%s/dum/installer.yml", base),
	}
	for _, candidate := range candidates {
		if _, err := os.Stat(candidate); err == nil {
			return candidate
		}
	}
	return candidates[1]
}
