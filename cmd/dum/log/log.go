// Package log defines the shell logging Cobra commands.
package log

import (
	"awong/dotfiles/cmd/dum/cli"
	"fmt"
	"strings"

	lg "awong/dotfiles/internal/logging"

	clog "github.com/charmbracelet/log"
	"github.com/spf13/cobra"
)

// SuccessLevel and PREFIX preserve the historical log command contract.
const (
	SuccessLevel clog.Level = cli.SuccessLevel
	PREFIX                  = cli.PREFIX
)

// Command holds the logger used by log subcommands.
type Command struct{ Log lg.Logger }

// NewCommand constructs logger dependencies for log commands.
func NewCommand(logger lg.Logger) *Command { return &Command{Log: logger} }
func addPrefixFlag(cmd *cobra.Command)     { cli.AddPrefixFlag(cmd) }

// NewLogCommand constructs the log Cobra command.
func NewLogCommand(rootUse string, dum *Command) *cobra.Command {
	use, alias := "log", "lg"
	cmd := &cobra.Command{Use: use, Aliases: []string{alias}, Short: fmt.Sprintf("%v-%v (or %v) logs for the shell", rootUse, use, alias), Long: strings.Join([]string{fmt.Sprintf("%v-%v (or %v) is a command line interface (cli) that helps your shell log.", rootUse, use, alias), "", "Log levels:", "  debug (or d)   - debug messages (lowest)", "  info (or i)    - informational messages", "  warn (or w)    - warning messages", "  error (or e)   - error messages", "  success (or s) - success messages (printed without log prefix)", "", "Logging:", "  Set ZSH_LOG_LEVEL env var to 'info' or 'debug' to filter output", "  Messages below the threshold are suppressed"}, "\n"), Example: strings.Join([]string{"  # Print a debug message", fmt.Sprintf("  %v %v debug \"Starting process\"", rootUse, use), "", "  # Print an info message", fmt.Sprintf("  %v %v i \"Installation complete\"", rootUse, use), "", "  # Print a success message (no log prefix)", fmt.Sprintf("  %v %v success \"All done!\"", rootUse, use), "", "  # Print with a prefix", fmt.Sprintf("  %v %v warn -p myapp \"Disk space low\"", rootUse, use), "", "  # Print an error message", fmt.Sprintf("  %v %v e \"Connection failed\"", rootUse, use)}, "\n")}
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

func (d *Command) executeLogging(cmd *cobra.Command, args []string, level clog.Level) error {
	logger := d.Log
	p, err := cmd.Flags().GetString(cli.PREFIX)
	if err != nil {
		p = ""
	}
	if p != "" {
		d.Log = logger.WithPrefix(p)
	}
	message := strings.Join(args, " ")
	switch level {
	case clog.DebugLevel:
		d.Log.Debugf("%v", message)
	case clog.InfoLevel:
		d.Log.Infof("%v", message)
	case clog.WarnLevel:
		d.Log.Warnf("%v", message)
	case clog.ErrorLevel:
		d.Log.Errorf("%v", message)
	case clog.FatalLevel:
		d.Log.Fatalf("%v", message)
	case SuccessLevel:
		if p != "" {
			if err := d.Log.Printlnf("%v %v", p, message); err != nil {
				return fmt.Errorf("failed to println %v: %w", message, err)
			}
		} else if err := d.Log.Printlnf("%v", message); err != nil {
			return fmt.Errorf("failed to println %v: %w", message, err)
		}
	default:
		return fmt.Errorf("invalid log level: %v", level)
	}
	return nil
}
