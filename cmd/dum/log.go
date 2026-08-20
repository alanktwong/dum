package main

import (
	"fmt"
	"strings"

	lg "github.com/alanktwong/dum/internal/logging"

	clog "github.com/charmbracelet/log"
	"github.com/spf13/cobra"
)

// LogCommand holds the logger used by log subcommands.
type LogCommand struct{ Log lg.Logger }

// newLogDeps constructs logger dependencies for log commands.
func newLogDeps(logger lg.Logger) *LogCommand { return &LogCommand{Log: logger} }

// NewLogCommand constructs the log Cobra command.
func NewLogCommand(rootUse string, dum *LogCommand) *cobra.Command {
	use, alias := "log", "lg"
	cmd := &cobra.Command{
		Use:     use,
		Aliases: []string{alias},
		Short:   fmt.Sprintf("%v-%v (or %v) logs for the shell", rootUse, use, alias),
		Long: strings.Join(
			[]string{
				fmt.Sprintf(
					"%v-%v (or %v) is a command line interface (cli) that helps your shell log.",
					rootUse,
					use,
					alias,
				),
				"",
				"Log levels:",
				"  debug (or d)   - debug messages (lowest)",
				"  info (or i)    - informational messages",
				"  warn (or w)    - warning messages",
				"  error (or e)   - error messages",
				"  success (or s) - success messages (printed without log prefix)",
				"",
				"Logging:",
				"  Set ZSH_LOG_LEVEL env var to 'info' or 'debug' to filter output",
				"  Messages below the threshold are suppressed",
			},
			"\n",
		),
		Example: strings.Join(
			[]string{
				"  # Print a debug message",
				fmt.Sprintf("  %v %v debug \"Starting process\"", rootUse, use),
				"",
				"  # Print an info message",
				fmt.Sprintf("  %v %v i \"Installation complete\"", rootUse, use),
				"",
				"  # Print a success message (no log prefix)",
				fmt.Sprintf("  %v %v success \"All done!\"", rootUse, use),
				"",
				"  # Print with a prefix",
				fmt.Sprintf("  %v %v warn -p myapp \"Disk space low\"", rootUse, use),
				"",
				"  # Print an error message",
				fmt.Sprintf("  %v %v e \"Connection failed\"", rootUse, use),
			},
			"\n",
		),
	}
	debugCmd := newDebugCommand(rootUse, use, dum)
	AddPrefixFlag(debugCmd)
	cmd.AddCommand(debugCmd)
	infoCmd := newInfoCommand(rootUse, use, dum)
	AddPrefixFlag(infoCmd)
	cmd.AddCommand(infoCmd)
	successCmd := newSuccessCommand(rootUse, use, dum)
	AddPrefixFlag(successCmd)
	cmd.AddCommand(successCmd)
	warnCmd := newWarnCommand(rootUse, use, dum)
	AddPrefixFlag(warnCmd)
	cmd.AddCommand(warnCmd)
	errorCmd := newErrorCommand(rootUse, use, dum)
	AddPrefixFlag(errorCmd)
	cmd.AddCommand(errorCmd)
	return cmd
}

func newDebugCommand(rootUse, logUse string, dum *LogCommand) *cobra.Command {
	use := "debug"
	alias := "d"
	return &cobra.Command{
		Use:     use,
		Aliases: []string{alias},
		Args:    cobra.MinimumNArgs(1),
		Short:   fmt.Sprintf("%v-%v-%v (or %v) logs [messages] at %v level", rootUse, logUse, use, alias, use),
		Long: strings.Join([]string{
			fmt.Sprintf("%v-%v-%v (or %v) logs [messages] at debug level.", rootUse, logUse, use, alias),
			"",
			"Messages are only printed if ZSH_LOG_LEVEL is set to 'debug'.",
		}, "\n"),
		Example: strings.Join([]string{
			"  # Print debug message",
			fmt.Sprintf("  %v %v %v \"Starting process\"", rootUse, logUse, use),
			"",
			"  # With prefix",
			fmt.Sprintf("  %v %v %v -p myapp \"Loading config\"", rootUse, logUse, use),
		}, "\n"),
		RunE: func(cmd *cobra.Command, args []string) error {
			return dum.executeLogging(cmd, args, clog.DebugLevel)
		},
	}
}

func newInfoCommand(rootUse, logUse string, dum *LogCommand) *cobra.Command {
	use := "info"
	alias := "i"
	return &cobra.Command{
		Use:     use,
		Aliases: []string{alias},
		Args:    cobra.MinimumNArgs(1),
		Short:   fmt.Sprintf("%v-%v-%v (or %v) logs [messages] at %v level", rootUse, logUse, use, alias, use),
		Long: strings.Join([]string{
			fmt.Sprintf("%v-%v-%v (or %v) logs [messages] at info level.", rootUse, logUse, use, alias),
			"",
			"Messages are only printed if ZSH_LOG_LEVEL is set to 'info' or 'debug'.",
		}, "\n"),
		Example: strings.Join([]string{
			"  # Print info message",
			fmt.Sprintf("  %v %v %v \"Installation complete\"", rootUse, logUse, use),
			"",
			"  # With prefix",
			fmt.Sprintf("  %v %v %v -p myapp \"Config loaded\"", rootUse, logUse, use),
		}, "\n"),
		RunE: func(cmd *cobra.Command, args []string) error {
			return dum.executeLogging(cmd, args, clog.InfoLevel)
		},
	}
}

func newSuccessCommand(rootUse, logUse string, dum *LogCommand) *cobra.Command {
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

func newWarnCommand(rootUse, logUse string, dum *LogCommand) *cobra.Command {
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

func newErrorCommand(rootUse, logUse string, dum *LogCommand) *cobra.Command {
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

func (d *LogCommand) executeLogging(cmd *cobra.Command, args []string, level clog.Level) error {
	logger := d.Log
	p, err := cmd.Flags().GetString(PREFIX)
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
