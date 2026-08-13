package log

import (
	"fmt"
	"strings"

	clog "github.com/charmbracelet/log"
	"github.com/spf13/cobra"
)

func newErrorCommand(rootUse, logUse string, dum *Command) *cobra.Command {
	use := "error"
	alias := "e"
	return &cobra.Command{
		Use:     use,
		Aliases: []string{alias},
		Args:    cobra.MinimumNArgs(1),
		Short:   fmt.Sprintf("%v-%v-%v (or %v) logs [messages] at %v level", rootUse, logUse, use, alias, use),
		Long: strings.Join([]string{
			fmt.Sprintf("%v-%v-%v (or %v) logs [messages] at error level.", rootUse, logUse, use, alias),
			"",
			"Messages are only printed if ZSH_LOG_LEVEL is set to 'info' or 'debug'.",
		}, "\n"),
		Example: strings.Join([]string{
			"  # Print error message",
			fmt.Sprintf("  %v %v %v \"Connection refused\"", rootUse, logUse, use),
			"",
			"  # With prefix",
			fmt.Sprintf("  %v %v %v -p myapp \"Failed to start server\"", rootUse, logUse, use),
		}, "\n"),
		RunE: func(cmd *cobra.Command, args []string) error {
			return dum.executeLogging(cmd, args, clog.ErrorLevel)
		},
	}
}
