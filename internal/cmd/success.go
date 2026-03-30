package cmd

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
)

func newSuccessCommand(rootUse, logUse string, dum *Dum) *cobra.Command {
	use := "success"
	alias := "s"
	return &cobra.Command{
		Use:     use,
		Aliases: []string{alias},
		Args:    cobra.MinimumNArgs(1),
		Short:   fmt.Sprintf("%v-%v-%v (or %v) logs [messages] at %v level", rootUse, logUse, use, alias, use),
		Long: fmt.Sprintf(`%v-%v-%v (or %v) logs [messages] at success level.

Success messages are printed without the standard log prefix.

Example:
  %v %v-%v -s myapp hi world`,
			rootUse, logUse, use, alias, rootUse, logUse, alias),
		Example: strings.Join([]string{
			"# Print success message",
			"# With prefix",
		}, "\n"),
		RunE: func(cmd *cobra.Command, args []string) error {
			return dum.executeLogging(cmd, args, SuccessLevel)
		},
	}
}
