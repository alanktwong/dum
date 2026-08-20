package main

import (
	"context"
	"fmt"
	"strings"

	lg "alanktwong/dum/internal/logging"
	rn "alanktwong/dum/internal/rename"

	"github.com/spf13/cobra"
)

// RenameCommand holds Cobra command dependencies for rename operations.
type RenameCommand struct {
	Log     lg.Logger
	Service *rn.Renamer
}

// newRenameDeps constructs rename command dependencies.
func newRenameDeps(logger lg.Logger) *RenameCommand {
	return &RenameCommand{Log: logger, Service: rn.NewRenamer(logger)}
}

// Options is retained as the command-to-service request type.
type Options = rn.Options

// NewRenameCommand constructs the rename Cobra command.
func NewRenameCommand(rootUse string, dum *RenameCommand) *cobra.Command {
	use, alias := "rename", "r"
	cmd := &cobra.Command{
		Use:     use,
		Aliases: []string{alias},
		Short:   fmt.Sprintf("%v-%v (or %v) files based on input options", rootUse, use, alias),
		Long: strings.Join(
			[]string{
				fmt.Sprintf(
					"%v-%v (or %v) is a command line tool that renames files based on input flags.",
					rootUse,
					use,
					alias,
				),
				"",
				"Logging:",
				"  Set ZSH_LOG_LEVEL env var to 'info' or 'debug', or use --verbose flag",
			},
			"\n",
		),
		Example: strings.Join(
			[]string{
				"  # Replace '_DSC' with 'photo' in all jpg files",
				fmt.Sprintf("  %v %v -s _DSC -r photo *.jpg", rootUse, alias),
				"",
				"  # Rename files to lowercase",
				fmt.Sprintf("  %v %v --lowercase *.jpg", rootUse, use),
				"",
				"  # Trim 4 characters from front of filename",
				fmt.Sprintf("  %v %v --trim-front 4 IMG_1234.jpg", rootUse, use),
				"",
				"  # Trim 3 characters from back of filename (before extension)",
				fmt.Sprintf("  %v %v --trim-back 4 file.bak.txt", rootUse, use),
				"",
				"  # Add sequential numbers starting from 1",
				fmt.Sprintf("  %v %v --iterate 1 photo.jpg vacation.png", rootUse, use),
				"  # Result: photo_1.jpg, vacation_2.png",
				"",
				"  # Preview changes with verbose output",
				fmt.Sprintf("  %v %v --v -s old -r new *.txt", rootUse, use),
				"",
				"  # Preview changes without actually renaming",
				fmt.Sprintf("  %v %v -s old -r new -d *.txt", rootUse, use),
			},
			"\n",
		),
		Args: cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error { return dum.run(cmd, args) },
	}
	cmd.Flags().StringP(REPLACE, "r", "", "replaces SOURCE with REPLACE in the PAT. See example.")
	cmd.Flags().StringP(SOURCE, "s", "", "replaces SOURCE with REPLACE in the PAT. See example.")
	AddVerboseFlag(cmd)
	cmd.Flags().BoolP(LOWERCASE, "l", false, "convert the filename to lowercase. E.g. FILE.JPG -> file.jpg")
	AddDryRunFlag(cmd)
	cmd.Flags().Uint16P(
		TRIMFRONT,
		"f",
		0,
		"trims N characters from the front of the filename. E.g. file.jpg -> le.jpg if N = 2",
	)
	cmd.Flags().Uint16P(
		TRIMBACK,
		"b",
		0,
		"trims N characters from the back of the filename. E.g. file.jpg -> fi.jpg if N = 2",
	)
	cmd.Flags().Uint16P(
		ITERATE,
		"i",
		0,
		"suffix an integer to the filename starting from N. E.g. file.jpg -> file_1.jpg if N = 1",
	)
	return cmd
}

func (d *RenameCommand) run(cmd *cobra.Command, args []string) error {
	ctx := cmd.Context()
	if ctx == nil {
		return fmt.Errorf("nil context")
	}
	if len(args) < 1 {
		return fmt.Errorf("dum requires more than 1 argument")
	}
	v := GetVerbosityFromCommand(cmd)
	d.Log.SetLevel(v.Level())
	dr, _ := cmd.Flags().GetBool(DRYRUN)
	src, err := cmd.Flags().GetString(SOURCE)
	if err != nil {
		return fmt.Errorf("error getting %v flag: %w", SOURCE, err)
	}
	rep, err := cmd.Flags().GetString(REPLACE)
	if err != nil {
		return fmt.Errorf("error getting %v flag: %w", REPLACE, err)
	}
	lc, _ := cmd.Flags().GetBool(LOWERCASE)
	tf, _ := cmd.Flags().GetUint16(TRIMFRONT)
	tb, _ := cmd.Flags().GetUint16(TRIMBACK)
	it, _ := cmd.Flags().GetUint16(ITERATE)
	d.Log.Debug(
		"Running the rename command:",
		SOURCE,
		src,
		REPLACE,
		rep,
		"PATTERN",
		args,
		DRYRUN,
		dr,
		LOWERCASE,
		lc,
		TRIMFRONT,
		tf,
		TRIMBACK,
		tb,
		ITERATE,
		it,
		VERBOSE,
		v.Verbose,
	)
	if src != "" && rep == "" {
		return fmt.Errorf("cannot replace %s without replace", src)
	} else if src == "" && rep != "" {
		return fmt.Errorf("cannot replace for %s without source", rep)
	}
	for idx, p := range args {
		if err := d.rename(
			ctx,
			Options{
				Index:       idx,
				Source:      src,
				Replace:     rep,
				FilePattern: p,
				Lowercase:   lc,
				DryRun:      dr,
				TrimFront:   int(tf),
				TrimBack:    int(tb),
				Iterate:     int(it),
			},
		); err != nil {
			d.Log.Errorf("Error renaming file %s: %v", p, err)
		}
	}
	return nil
}

func (d *RenameCommand) rename(ctx context.Context, options Options) error {
	if d.Service == nil {
		d.Service = rn.NewRenamer(d.Log)
	}
	if err := d.Service.Rename(ctx, options); err != nil {
		return fmt.Errorf("rename service: %w", err)
	}
	return nil
}
