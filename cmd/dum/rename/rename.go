// Package rename defines the file rename Cobra command.
package rename

import (
	"awong/dotfiles/cmd/dum/cli"
	"context"
	"fmt"
	"strings"

	lg "awong/dotfiles/internal/logging"
	rn "awong/dotfiles/internal/rename"

	"github.com/spf13/cobra"
)

const (
	replace   = "replace"
	source    = "source"
	lowercase = "lowercase"
	trimFront = "trimfront"
	trimBack  = "trimback"
	iterate   = "iterate"
)

// Command holds Cobra command dependencies for rename operations.
type Command struct {
	Log     lg.Logger
	Service *rn.Renamer
}

// NewCommand constructs rename command dependencies.
func NewCommand(logger lg.Logger) *Command {
	return &Command{Log: logger, Service: rn.NewRenamer(logger)}
}

// Options is retained as the command-to-service request type.
type Options = rn.Options

// NewRenameCommand constructs the rename Cobra command.
func NewRenameCommand(rootUse string, dum *Command) *cobra.Command {
	use, alias := "rename", "r"
	cmd := &cobra.Command{Use: use, Aliases: []string{alias}, Short: fmt.Sprintf("%v-%v (or %v) files based on input options", rootUse, use, alias), Long: strings.Join([]string{fmt.Sprintf("%v-%v (or %v) is a command line tool that renames files based on input flags.", rootUse, use, alias), "", "Logging:", "  Set ZSH_LOG_LEVEL env var to 'info' or 'debug', or use --verbose flag"}, "\n"), Example: strings.Join([]string{"  # Replace '_DSC' with 'photo' in all jpg files", fmt.Sprintf("  %v %v -s _DSC -r photo *.jpg", rootUse, alias), "", "  # Rename files to lowercase", fmt.Sprintf("  %v %v --lowercase *.jpg", rootUse, use), "", "  # Trim 4 characters from front of filename", fmt.Sprintf("  %v %v --trim-front 4 IMG_1234.jpg", rootUse, use), "", "  # Trim 3 characters from back of filename (before extension)", fmt.Sprintf("  %v %v --trim-back 4 file.bak.txt", rootUse, use), "", "  # Add sequential numbers starting from 1", fmt.Sprintf("  %v %v --iterate 1 photo.jpg vacation.png", rootUse, use), "  # Result: photo_1.jpg, vacation_2.png", "", "  # Preview changes with verbose output", fmt.Sprintf("  %v %v --v -s old -r new *.txt", rootUse, use), "", "  # Preview changes without actually renaming", fmt.Sprintf("  %v %v -s old -r new -d *.txt", rootUse, use)}, "\n"), Args: cobra.MinimumNArgs(1), RunE: func(cmd *cobra.Command, args []string) error { return dum.run(cmd, args) }}
	cmd.Flags().StringP(replace, "r", "", "replaces SOURCE with REPLACE in the PAT. See example.")
	cmd.Flags().StringP(source, "s", "", "replaces SOURCE with REPLACE in the PAT. See example.")
	cli.AddVerboseFlag(cmd)
	cmd.Flags().BoolP(lowercase, "l", false, "convert the filename to lowercase. E.g. FILE.JPG -> file.jpg")
	cli.AddDryRunFlag(cmd)
	cmd.Flags().Uint16P(trimFront, "f", 0, "trims N characters from the front of the filename. E.g. file.jpg -> le.jpg if N = 2")
	cmd.Flags().Uint16P(trimBack, "b", 0, "trims N characters from the back of the filename. E.g. file.jpg -> fi.jpg if N = 2")
	cmd.Flags().Uint16P(iterate, "i", 0, "suffix an integer to the filename starting from N. E.g. file.jpg -> file_1.jpg if N = 1")
	return cmd
}

func (d *Command) run(cmd *cobra.Command, args []string) error {
	ctx := cmd.Context()
	if ctx == nil {
		return fmt.Errorf("nil context")
	}
	if len(args) < 1 {
		return fmt.Errorf("dum requires more than 1 argument")
	}
	v := cli.GetVerbosityFromCommand(cmd)
	d.Log.SetLevel(v.Level())
	dr, _ := cmd.Flags().GetBool(cli.DRYRUN)
	src, err := cmd.Flags().GetString(source)
	if err != nil {
		return fmt.Errorf("error getting %v flag: %w", source, err)
	}
	rep, err := cmd.Flags().GetString(replace)
	if err != nil {
		return fmt.Errorf("error getting %v flag: %w", replace, err)
	}
	lc, _ := cmd.Flags().GetBool(lowercase)
	tf, _ := cmd.Flags().GetUint16(trimFront)
	tb, _ := cmd.Flags().GetUint16(trimBack)
	it, _ := cmd.Flags().GetUint16(iterate)
	d.Log.Debug("Running the rename command:", source, src, replace, rep, "PATTERN", args, cli.DRYRUN, dr, lowercase, lc, trimFront, tf, trimBack, tb, iterate, it, cli.VERBOSE, v.Verbose)
	if src != "" && rep == "" {
		return fmt.Errorf("cannot replace %s without replace", src)
	} else if src == "" && rep != "" {
		return fmt.Errorf("cannot replace for %s without source", rep)
	}
	for idx, p := range args {
		if err := d.rename(ctx, Options{Index: idx, Source: src, Replace: rep, FilePattern: p, Lowercase: lc, DryRun: dr, TrimFront: int(tf), TrimBack: int(tb), Iterate: int(it)}); err != nil {
			d.Log.Errorf("Error renaming file %s: %v", p, err)
		}
	}
	return nil
}

func (d *Command) rename(ctx context.Context, options Options) error {
	if d.Service == nil {
		d.Service = rn.NewRenamer(d.Log)
	}
	if err := d.Service.Rename(ctx, options); err != nil {
		return fmt.Errorf("rename service: %w", err)
	}
	return nil
}
