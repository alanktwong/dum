package cmd

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

// NewRenameCommand provides a command that renames files based on the input options.
// This was originally based on:
// See http://www.mattweber.org/2007/03/04/python-script-renamepy/
// See https://blog.thestateofme.com/2010/10/05/renumbering-media-files/
// NewInstallCommand provides a command that installs plays and task.
func NewRenameCommand(rootUse string, dum *Dum) *cobra.Command {
	use := "rename"
	alias := "r"
	cmd := &cobra.Command{
		Use:     use,
		Aliases: []string{alias},
		Short:   fmt.Sprintf("%v-%v (or %v) files based on input options", rootUse, use, alias),
		Long: strings.Join([]string{
			fmt.Sprintf("%v-%v (or %v) is a command line tool that renames files based on input flags.", rootUse, use, alias),
			"",
			"Flags:",
			"  -s, --source      String to find in filename",
			"  -r, --replace     String to replace source with",
			"  -l, --lowercase   Convert filename to lowercase",
			"  -f, --trim-front   Trim N characters from front of filename",
			"  -b, --trim-back    Trim N characters from back of filename",
			"  -i, --iterate      Add sequential suffix starting from N",
			"  -d, --dryrun       Preview changes without renaming",
			"",
			"Logging:",
			"  Set ZSH_LOG_LEVEL env var to 'info' or 'debug', or use --verbose flag",
		}, "\n"),
		Example: strings.Join([]string{
			"# Replace '_DSC' with 'photo' in all jpg files",
			fmt.Sprintf("%v %v -s _DSC -r photo *.jpg", rootUse, alias),
			"",
			"# Rename files to lowercase",
			fmt.Sprintf("%v %v --lowercase *.jpg", rootUse, use),
			"",
			"# Trim 4 characters from front of filename",
			fmt.Sprintf("%v %v --trim-front 4 IMG_1234.jpg", rootUse, use),
			"",
			"# Trim 3 characters from back of filename (before extension)",
			fmt.Sprintf("%v %v --trim-back 4 file.bak.txt", rootUse, use),
			"",
			"# Add sequential numbers starting from 1",
			fmt.Sprintf("%v %v --iterate 1 photo.jpg vacation.png", rootUse, use),
			"# Result: photo_1.jpg, vacation_2.png",
			"",
			"# Preview changes with verbose output",
			fmt.Sprintf("%v %v --v -s old -r new *.txt", rootUse, use),
			"",
			"# Preview changes without actually renaming",
			fmt.Sprintf("%v %v -s old -r new -d *.txt", rootUse, use),
		}, "\n"),
		Args: cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			if ctx == nil {
				return fmt.Errorf("nil context")
			}
			if len(args) < 1 {
				return fmt.Errorf("%v requires more than 1 argument", rootUse)
			}
			verbosity := GetVerbosityFromCommand(cmd)
			dum.Log.SetLevel(verbosity.Level())
			dryrun, err := cmd.Flags().GetBool(DRYRUN)
			if err != nil {
				dryrun = false
			}
			source, err := cmd.Flags().GetString(SOURCE)
			if err != nil {
				return fmt.Errorf("error getting %v flag: %w", SOURCE, err)
			}
			replace, err := cmd.Flags().GetString(REPLACE)
			if err != nil {
				return fmt.Errorf("error getting %v flag: %w", REPLACE, err)
			}
			lowercase, err := cmd.Flags().GetBool(LOWERCASE)
			if err != nil {
				lowercase = false
			}
			trimfront, err := cmd.Flags().GetUint16(TRIMFRONT)
			if err != nil {
				trimfront = 0
			}
			trimback, err := cmd.Flags().GetUint16(TRIMBACK)
			if err != nil {
				trimback = 0
			}
			iterate, err := cmd.Flags().GetUint16(ITERATE)
			if err != nil {
				iterate = 0
			}
			dum.Log.Debug("Running the rename command:",
				SOURCE, source,
				REPLACE, replace,
				"PATTERN", args,
				DRYRUN, dryrun,
				LOWERCASE, lowercase,
				TRIMFRONT, trimfront,
				TRIMBACK, trimback,
				ITERATE, iterate,
				VERBOSE, verbosity.Verbose,
			)
			if source != "" && replace == "" {
				return fmt.Errorf("cannot replace %s without replace", source)
			} else if source == "" && replace != "" {
				return fmt.Errorf("cannot replace for %s without source", replace)
			}
			for idx, filePattern := range args {
				err := dum.rename(ctx, RenameOptions{
					Index:       idx,
					Source:      source,
					Replace:     replace,
					FilePattern: filePattern,
					Lowercase:   lowercase,
					DryRun:      dryrun,
					TrimFront:   int(trimfront),
					TrimBack:    int(trimback),
					Iterate:     int(iterate),
					Verbose:     verbosity,
				})
				if err != nil {
					dum.Log.Errorf("Error renaming file %s: %v", filePattern, err)
				}
			}

			return nil
		},
	}

	fundamental := "Replaces SOURCE with REPLACE in the PAT. See example."
	cmd.Flags().StringP(REPLACE, "r", "", fundamental)
	cmd.Flags().StringP(SOURCE, "s", "", fundamental)

	addVerboseFlag(cmd)

	lowerCaseUsage := "Convert the filename to lowercase. E.g. FILE.JPG -> file.jpg"
	cmd.Flags().BoolP(LOWERCASE, "l", false, lowerCaseUsage)

	addDryRunFlag(cmd)

	frontUsage := "Trims N characters from the front of the filename. E.g. file.jpg -> le.jpg if N = 2"
	cmd.Flags().Uint16P(TRIMFRONT, "f", 0, frontUsage)

	backUsage := "Trims N characters from the back of the filename. E.g. file.jpg -> fi.jpg if N = 2"
	cmd.Flags().Uint16P(TRIMBACK, "b", 0, backUsage)

	iterateUsage := "Suffix an integer to the filename starting from N. E.g. file.jpg -> file_1.jpg if N = 1"
	cmd.Flags().Uint16P(ITERATE, "i", 0, iterateUsage)
	return cmd
}

// RenameOptions provides parameterized options to rename func.
type RenameOptions struct {
	Index       int
	FilePattern string
	Source      string
	Replace     string
	Lowercase   bool
	DryRun      bool
	TrimFront   int
	TrimBack    int
	Iterate     int
	Verbose     Verbosity
}

// rename renames a file with the given options.
func (d *Dum) rename(_ context.Context, options RenameOptions) error {
	ellipsis := "..."
	// split the pathname, filename, and extension
	pathname := filepath.Dir(options.FilePattern)
	d.Log.Debugf("%v pathname = %s", ellipsis, pathname)

	filename := strings.TrimSuffix(filepath.Base(options.FilePattern), filepath.Ext(options.FilePattern))
	d.Log.Debugf("%v filename = %s", ellipsis, filename)
	extension := filepath.Ext(options.FilePattern)
	d.Log.Debugf("%v extension = %s", ellipsis, extension)

	// trim characters from front
	if options.TrimFront > 0 && options.TrimFront < len(filename) {
		filename = filename[options.TrimFront:]
		d.Log.Debugf("%v trimmed-front filename = %s", ellipsis, filename)
	}
	// trim characters from back
	if options.TrimBack > 0 && options.TrimBack < len(filename) {
		filename = filename[:len(filename)-options.TrimBack]
		d.Log.Debugf("%v trimmed-back filename = %s", ellipsis, filename)
	}
	// replace values if any
	if options.Source != "" && options.Replace != "" {
		filename = strings.ReplaceAll(filename, options.Source, options.Replace)
		d.Log.Debugf("%v replaced filename = %s", ellipsis, filename)
	}

	// convert to lower case if flag is set
	if options.Lowercase {
		filename = strings.ToLower(filename)
		extension = strings.ToLower(extension)
		d.Log.Debugf("%v lowercased filename = %s", ellipsis, filename)
		d.Log.Debugf("%v lowercased extension = %s", ellipsis, extension)
	}

	// append to filename starting from Iterate
	if options.Iterate > 0 {
		filename = fmt.Sprintf("%s_%d", filename, options.Index+options.Iterate)
		d.Log.Debugf("%v iterated filename = %s", ellipsis, filename)
	}
	newFilepath := filepath.Join(pathname, filename+extension)
	d.Log.Infof("%v renaming %s -> %s", ellipsis, options.FilePattern, newFilepath)
	if !options.DryRun {
		if err := os.Rename(options.FilePattern, newFilepath); err != nil {
			return fmt.Errorf("error renaming file %s: %w", options.FilePattern, err)
		}
		return nil
	}
	return nil
}
