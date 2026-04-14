package cmd

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/log"
	"github.com/spf13/cobra"
)

func newInfoCommand(rootUse, logUse string, dum *Dum) *cobra.Command {
	use := "info"
	alias := "i"
	return &cobra.Command{
		Use:     use,
		Aliases: []string{alias},
		Args:    cobra.MinimumNArgs(1),
		Short:   fmt.Sprintf("%v-%v-%v (or %v) logs [messages] at %v level", rootUse, logUse, use, alias, use),
		Long: strings.Join([]string{
			fmt.Sprintf("%v-%v-%v (or %v) logs [messages] at info level.", rootUse, logUse, use, alias),
			"",
			"Messages are only printed if ZSH_LOG_LEVEL is set to 'info' or 'debug'.",
		}, "\n"),
		Example: strings.Join([]string{
			"  # Print info message",
			fmt.Sprintf("  %v %v %v \"Installation complete\"", rootUse, logUse, use),
			"",
			"  # With prefix",
			fmt.Sprintf("  %v %v %v -p myapp \"Config loaded\"", rootUse, logUse, use),
		}, "\n"),
		RunE: func(cmd *cobra.Command, args []string) error {
			return dum.executeLogging(cmd, args, log.InfoLevel)
		},
	}
}
