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
		Long: strings.Join([]string{
			fmt.Sprintf("%v-%v-%v (or %v) logs [messages] at success level.", rootUse, logUse, use, alias),
			"",
			"Success messages are printed without the standard log prefix.",
		}, "\n"),
		Example: strings.Join([]string{
			"  # Print success message",
			fmt.Sprintf("  %v %v %v \"All tasks complete\"", rootUse, logUse, use),
			"",
			"  # With prefix",
			fmt.Sprintf("  %v %v %v -p myapp \"Build finished\"", rootUse, logUse, use),
		}, "\n"),
		RunE: func(cmd *cobra.Command, args []string) error {
			return dum.executeLogging(cmd, args, SuccessLevel)
		},
	}
}
