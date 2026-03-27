package cmd

import (
	"fmt"

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
		Long:    fmt.Sprintf("%v-%v-%v (or %v) logs [messages] at %v level joined by spaces", rootUse, logUse, use, alias, use),
		Example: fmt.Sprintf("'%v %v-%v hi world' produces 'hi world' at %v level", rootUse, logUse, alias, use),
		RunE: func(cmd *cobra.Command, args []string) error {
			return dum.executeLogging(cmd, args, log.InfoLevel)
		},
	}
}
