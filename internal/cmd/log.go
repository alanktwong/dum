package cmd

import (
	l "awong/dotfiles/internal/logging"
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
				"",
				"Log levels:",
				"  debug  - debug messages (lowest)",
				"  info   - informational messages",
				"  warn   - warning messages",
				"  error  - error messages",
				"  success - success messages (printed without log prefix)",
				"",
				"Logging:",
				"  Set ZSH_LOG_LEVEL env var to 'info' or 'debug' to filter output",
				"  Messages below the threshold are suppressed",
			},
			"\n"),
		Example: strings.Join([]string{
			"# Print debug message",
			"# Print info message",
			"# Print success message",
			"# Print with prefix",
		}, "\n"),
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
				return fmt.Errorf("failed to println %v: %w", message, err)
			}
			return nil
		}
		if err := d.Log.Printlnf("%v", message); err != nil {
			return fmt.Errorf("failed to println %v: %w", message, err)
		}
		return nil
	default:
		return fmt.Errorf("invalid log level: %v", level)
	}
}
