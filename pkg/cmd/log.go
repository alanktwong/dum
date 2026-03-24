package cmd

import (
	l "awong/dotfiles/pkg/logging"
	"fmt"
	"strings"

	"github.com/charmbracelet/log"
	"github.com/spf13/cobra"
)

// NewLogCommand provides a command that logs.
func NewLogCommand(rootUse string, dum *Dum) *cobra.Command {
	use := "log"
	alias := "lg"
	cmd := &cobra.Command{
		Use:     use,
		Aliases: []string{alias},
		Short:   fmt.Sprintf("%v-%v (or %v) logs for the shell", rootUse, use, alias),
		Long: strings.Join(
			[]string{
				fmt.Sprintf("%v-%v (or %v) is a command line interface (cli) that helps your shell log.",
					rootUse, use, alias),
				"The log level is printed out to standard output if it exceed the environment log level.",
				"The environment log level is the ZSH_LOG_LEVEL.",
				"ZSH_LOG_LEVEL can be set to 'info' or 'debug'.",
			},
			"\n"),
		Example: fmt.Sprintf("'%v %v %v hi world' produces 'hi world' at %v level", rootUse, use, "d", "debug"),
	}

	debugCmd := newDebugCommand(rootUse, use, dum)
	addPrefixFlag(debugCmd)
	cmd.AddCommand(debugCmd)

	infoCmd := newInfoCommand(rootUse, use, dum)
	addPrefixFlag(infoCmd)
	cmd.AddCommand(infoCmd)

	successCmd := newSuccessCommand(rootUse, use, dum)
	addPrefixFlag(successCmd)
	cmd.AddCommand(successCmd)

	warnCmd := newWarnCommand(rootUse, use, dum)
	addPrefixFlag(warnCmd)
	cmd.AddCommand(warnCmd)

	errorCmd := newErrorCommand(rootUse, use, dum)
	addPrefixFlag(errorCmd)
	cmd.AddCommand(errorCmd)

	return cmd
}

func (d *Dum) executeLogging(cmd *cobra.Command, args []string, level log.Level) error {
	logger := d.Log
	logger.SetLevel(l.EnvLevel())
	prefix, err := cmd.Flags().GetString(PREFIX)
	if err != nil {
		prefix = ""
	}

	if prefix != "" {
		d.Log = logger.WithPrefix(prefix)
	}
	message := strings.Join(args, " ")
	switch level {
	case log.DebugLevel:
		d.Log.Debugf("%v", message)
		return nil
	case log.InfoLevel:
		d.Log.Infof("%v", message)
		return nil
	case log.WarnLevel:
		d.Log.Warnf("%v", message)
		return nil
	case log.ErrorLevel:
		d.Log.Errorf("%v", message)
		return nil
	case log.FatalLevel:
		d.Log.Fatalf("%v", message)
		return nil
	case SuccessLevel:
		if prefix != "" {
			// this is not in the implementation of Printlnf b/c we only want prefix
			// for log4sh.
			if err := d.Log.Printlnf("%v %v", prefix, message); err != nil {
				return fmt.Errorf("failed to println %v: %v", message, err)
			}
			return nil
		}
		if err := d.Log.Printlnf("%v", message); err != nil {
			return fmt.Errorf("failed to println %v: %v", message, err)
		}
		return nil
	default:
		return fmt.Errorf("invalid log level: %v", level)
	}
}
