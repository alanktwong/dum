// Package cmd provides factory methods for subcommands
package cmd

import (
	"fmt"
	"math"
	"os"

	"github.com/charmbracelet/log"
	"github.com/spf13/cobra"
)

const (
	// DRYRUN is key to 'dry run' flag in Cobra.
	DRYRUN = "dryrun"
	// FILE is key to 'file' flag in Cobra.
	FILE = "file"
	// GROUP is key to 'group' flag in Cobra.
	GROUP = "group"
	// VERBOSE is key to 'verbose' flag in Cobra.
	VERBOSE = "verbose"
	// REPLACE is key to 'replace' flag in Cobra.
	REPLACE = "replace"
	// SOURCE is key to 'source' flag in Cobra.
	SOURCE = "source"
	// LOWERCASE is the key to 'lowercase' flag in Cobra.
	LOWERCASE = "lowercase"
	// TRIMFRONT is the key to 'trimfront' flag in Cobra.
	TRIMFRONT = "trimfront"
	// TRIMBACK is the key to 'trimback' flag in Cobra.
	TRIMBACK = "trimback"
	// ITERATE is the key to 'iterate' flag in Cobra.
	ITERATE = "iterate"
	// OUTPUT is the key to 'output' flag in Cobra.
	OUTPUT = "output"
	// PREFIX is the key to 'prefix' flag in Cobra.
	PREFIX = "prefix"
	// SuccessLevel is the level for success.
	SuccessLevel log.Level = math.MaxInt
)

func addGroupFlag(cmd *cobra.Command) {
	groupUsage := "Select GROUP to run"
	cmd.Flags().StringP(GROUP, "g", "", groupUsage)
}

func addFileFlag(cmd *cobra.Command) {
	defaultConfig := getDefaultConfig()
	fileUsage := "Select config FILE to use"
	cmd.Flags().StringP(FILE, "f", defaultConfig, fileUsage)
}

func addVerboseFlag(cmd *cobra.Command) {
	verboseUsage := "Increase output verbosity. E.g. --vv is debug and -v is info"
	cmd.Flags().CountP(VERBOSE, "v", verboseUsage)
}

func addDryRunFlag(cmd *cobra.Command) {
	dryRunUsage := fmt.Sprintf("Dry run of %s", cmd.Use)
	cmd.Flags().BoolP(DRYRUN, "d", false, dryRunUsage)
}

func addPrefixFlag(debugCmd *cobra.Command) {
	prefixUsage := "add a prefix to the logger"
	debugCmd.Flags().StringP(PREFIX, "p", "", prefixUsage)
}

func addOutputFlag(cmd *cobra.Command) {
	outputUsage := "Output file path (default: stdout)"
	cmd.Flags().StringP(OUTPUT, "o", "", outputUsage)
}

// getDefaultConfig fetches the default config file akin to how starship looks for its config file.
func getDefaultConfig() string {
	installerConfig := os.Getenv("INSTALLER_CONFIG")
	if installerConfig == "" {
		xdgConfig := os.Getenv("XDG_CONFIG_HOME")
		if xdgConfig == "" {
			xdgConfig = "~/.config"
		}
		installerConfig = fmt.Sprintf("%s/dum/%v", xdgConfig, "installer.yml")
	}
	return installerConfig
}
