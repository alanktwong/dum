package log

import (
	"fmt"
	"strings"

	clog "github.com/charmbracelet/log"
	"github.com/spf13/cobra"
)

func newWarnCommand(rootUse, logUse string, dum *Command) *cobra.Command {
	use := "warn"
	alias := "w"
	return &cobra.Command{
		Use:     use,
		Aliases: []string{alias},
		Args:    cobra.MinimumNArgs(1),
		Short:   fmt.Sprintf("%v-%v-%v (or %v) logs [messages] at %v level", rootUse, logUse, use, alias, use),
		Long: strings.Join([]string{
			fmt.Sprintf("%v-%v-%v (or %v) logs [messages] at warn level.", rootUse, logUse, use, alias),
			"",
			"Messages are only printed if ZSH_LOG_LEVEL is set to 'info' or 'debug'.",
		}, "\n"),
		Example: strings.Join([]string{
			"  # Print warning message",
			fmt.Sprintf("  %v %v %v \"Disk space running low\"", rootUse, logUse, use),
			"",
			"  # With prefix",
			fmt.Sprintf("  %v %v %v -p myapp \"Deprecated featured used\"", rootUse, logUse, use),
		}, "\n"),
		RunE: func(cmd *cobra.Command, args []string) error {
			return dum.executeLogging(cmd, args, clog.WarnLevel)
		},
	}
}
