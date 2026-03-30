package cmd

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/log"
	"github.com/spf13/cobra"
)

func newDebugCommand(rootUse, logUse string, dum *Dum) *cobra.Command {
	use := "debug"
	alias := "d"
	return &cobra.Command{
		Use:     use,
		Aliases: []string{alias},
		Args:    cobra.MinimumNArgs(1),
		Short:   fmt.Sprintf("%v-%v-%v (or %v) logs [messages] at %v level", rootUse, logUse, use, alias, use),
		Long: fmt.Sprintf(`%v-%v-%v (or %v) logs [messages] at debug level.

Messages are only printed if ZSH_LOG_LEVEL is set to 'debug'.

Example:
  %v %v-%v -s myapp hi world`,
			rootUse, logUse, use, alias, rootUse, logUse, alias),
		Example: strings.Join([]string{
			"# Print debug message",
			"# With prefix",
		}, "\n"),
		RunE: func(cmd *cobra.Command, args []string) error {
			return dum.executeLogging(cmd, args, log.DebugLevel)
		},
	}
}
