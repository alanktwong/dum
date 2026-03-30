package cmd

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/log"
	"github.com/spf13/cobra"
)

func newWarnCommand(rootUse, logUse string, dum *Dum) *cobra.Command {
	use := "warn"
	alias := "w"
	return &cobra.Command{
		Use:     use,
		Aliases: []string{alias},
		Args:    cobra.MinimumNArgs(1),
		Short:   fmt.Sprintf("%v-%v-%v (or %v) logs [messages] at %v level", rootUse, logUse, use, alias, use),
		Long: fmt.Sprintf(`%v-%v-%v (or %v) logs [messages] at warn level.

Messages are only printed if ZSH_LOG_LEVEL is set to 'info' or 'debug'.

Example:
  %v %v-%v -s myapp hi world`,
			rootUse, logUse, use, alias, rootUse, logUse, alias),
		Example: strings.Join([]string{
			"# Print warning message",
			"# With prefix",
		}, "\n"),
		RunE: func(cmd *cobra.Command, args []string) error {
			return dum.executeLogging(cmd, args, log.WarnLevel)
		},
	}
}
