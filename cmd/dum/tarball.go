package main

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	lg "github.com/alanktwong/dum/internal/logging"
	tb "github.com/alanktwong/dum/internal/tarball"
	ca "github.com/spf13/cobra"
)

const (
	tarballInputFlag            = "input"
	tarballTargetFlag           = "target"
	tarballTargetHelpDefault    = "TGZ_TARGET, else $XDG_DATA_HOME/tarballs, else ~/.local/share/tarballs"
	tarballDestinationFlag      = "dest"
	tarballOverwriteFlag        = "overwrite"
	tarballAllowUnsafeLinksFlag = "allow-unsafe-links"
	tarballPreserveOwnerFlag    = "preserve-owner"
	tarballDryRunFlag           = "dry-run"
)

type archiveRestorer func(context.Context, tb.RestoreOptions) (tb.RestoreResult, error)

// TarballCommand contains the dependencies used by tarball subcommands.
type TarballCommand struct {
	Log            lg.Logger
	Service        *tb.Service
	pickArchive    archivePicker
	restoreArchive archiveRestorer
}

// NewTarballCommand constructs the tarball command and its subcommands.
func NewTarballCommand(rootUse string, dum *TarballCommand) *ca.Command {
	if dum == nil {
		dum = newTarballDeps(nil)
	}
	if dum.Service == nil {
		dum.Service = tb.NewService(dum.Log)
	}
	if dum.pickArchive == nil {
		dum.pickArchive = selectArchive
	}
	if dum.restoreArchive == nil {
		dum.restoreArchive = dum.Service.Restore
	}

	use, alias := "tarball", "t"
	cmd := &ca.Command{
		Aliases: []string{alias},
		Use:     use,
		Short:   fmt.Sprintf("%v-%v (or %v) create and restore checksummed tar archives", rootUse, use, alias),
		Long: fmt.Sprintf(
			"%s-%s (or %s) creates gzip tar archives with SHA-256 and YAML sidecars, then verifies and safely restores them.",
			rootUse,
			use,
			alias,
		),
	}
	cmd.AddCommand(dum.newBackupCommand(rootUse, use))
	cmd.AddCommand(dum.newRestoreCommand(rootUse, use))
	return cmd
}

// tarballInputFileDoc joins command-specific help lines around the shared
// INPUT.yml format block used by list-based backup and restore.
func tarballInputFileDoc(head []string, tail ...string) string {
	const docLines = 12
	lines := make([]string, 0, len(head)+docLines+len(tail))
	lines = append(lines, head...)
	lines = append(lines,
		"INPUT.yml is a single YAML document with one required key, paths:",
		"",
		"  paths:",
		"    - Documents              # relative to the input.yml directory",
		"    - /mnt/archive/Photos    # absolute paths are used as-is",
		"",
		"Input rules:",
		"  - paths must contain at least one non-empty entry",
		"  - unknown keys and multiple YAML documents are rejected",
		"  - relative entries resolve from the directory containing INPUT.yml",
		"  - duplicate paths are rejected",
	)
	lines = append(lines, tail...)
	return strings.Join(lines, "\n")
}

func (d *TarballCommand) newBackupCommand(rootUse, tarballUse string) *ca.Command {
	use, alias := "backup", "b"
	cmd := &ca.Command{
		Use:     fmt.Sprintf("%s [PATH...]", use),
		Aliases: []string{alias},
		Short: fmt.Sprintf("%v-%v-%v (or %v) create one tarball and sidecars per path",
			rootUse, tarballUse, use, alias),
		Long: tarballInputFileDoc(
			[]string{
				fmt.Sprintf(
					"%s %s %s (or %s) creates one gzip tarball per source path with SHA-256 and YAML sidecars.",
					rootUse,
					tarballUse,
					use,
					alias,
				),
				"",
				"When PATH is omitted, --input is required; sources are read from that file:",
				"",
			},
			"",
			"PATH arguments replace the list from INPUT.yml; relative PATH arguments resolve",
			"from the current directory.",
		),
		Example: strings.Join([]string{
			"  # Backup by PATH",
			fmt.Sprintf("  %v %v %v ~/Documents", rootUse, tarballUse, use),
			"",
			"  # Backup in batch using an input yml",
			fmt.Sprintf("  %v %v %v --input input.yml", rootUse, tarballUse, use),
			"",
			"  # Backup in batch using an input yml using the dry-run flag",
			fmt.Sprintf("  %v %v %v --input input.yml --dry-run", rootUse, tarballUse, use),
		}, "\n"),
		Args: ca.ArbitraryArgs,
		RunE: func(cmd *ca.Command, args []string) error { return d.runBackup(cmd, args) },
	}
	cmd.Flags().String(tarballInputFlag, "", "read source paths from INPUT.yml (required when PATH is omitted)")
	addTarballTargetFlag(cmd, "directory for tarballs and sidecars")
	cmd.Flags().Bool(tarballDryRunFlag, false, "show planned outputs without writing files")
	return cmd
}

func (d *TarballCommand) newRestoreCommand(rootUse, tarballUse string) *ca.Command {
	use, alias := "restore", "r"
	cmd := &ca.Command{
		Use:     fmt.Sprintf("%s [ARCHIVE]", use),
		Aliases: []string{alias},
		Short: fmt.Sprintf("%v-%v-%v (or %v) verify and restore one tarball (from input.yml)",
			rootUse, tarballUse, use, alias),
		Long: tarballInputFileDoc(
			[]string{
				fmt.Sprintf(
					"%s %s %s (or %s) verifies checksums and sidecars, then restores one archive.",
					rootUse,
					tarballUse,
					use,
					alias,
				),
				"",
				"Without an ARCHIVE argument, --input is required. Each entry selects archives",
				"from --target:",
				"",
			},
			"",
			"List restore matches <basename>-*.tgz under --target for each entry, where",
			"<basename> is the final element of that entry. One match restores automatically;",
			"multiple matches use fzf or gum when installed. Entries must have unique",
			"basenames. Missing or invalid entries are reported while later entries continue.",
			"",
			"Passing ARCHIVE explicitly ignores INPUT.yml.",
		),
		Example: strings.Join([]string{
			"  # Restore by ARCHIVE",
			fmt.Sprintf("  %v %v %v /path/to/backup.tgz", rootUse, tarballUse, use),
			"",
			"  # Restore in batch using an input yml",
			fmt.Sprintf("  %v %v %v --input input.yml --dest /path/to/restore", rootUse, tarballUse, use),
			"",
			"  # Restore ARCHIVE allowing unsafe links",
			fmt.Sprintf("  %v %v %v backup.tgz --allow-unsafe-links", rootUse, tarballUse, use),
		}, "\n"),
		Args: ca.MaximumNArgs(1),
		RunE: func(cmd *ca.Command, args []string) error { return d.runRestore(cmd, args) },
	}

	cmd.Flags().String(tarballInputFlag, "", "read list-based restore paths from INPUT.yml (required without ARCHIVE)")
	addTarballTargetFlag(cmd, "directory containing tarballs and sidecars")
	cmd.Flags().String(tarballDestinationFlag, "", "restore under DIR (default: $HOME)")
	cmd.Flags().Bool(tarballOverwriteFlag, false, "replace conflicting files and links; never replace directories")
	cmd.Flags().Bool(tarballAllowUnsafeLinksFlag, false, "allow recreating absolute or escaping symlinks")
	cmd.Flags().Bool(tarballPreserveOwnerFlag, false, "attempt archived UID/GID restoration")
	cmd.Flags().Bool(tarballDryRunFlag, false, "verify and report planned paths without changing destination")
	return cmd
}

func (d *TarballCommand) runBackup(cmd *ca.Command, args []string) error {
	if cmd.Context() == nil {
		return fmt.Errorf("nil context")
	}
	inputPath, err := cmd.Flags().GetString(tarballInputFlag)
	if err != nil {
		return fmt.Errorf("read --input flag: %w", err)
	}
	target, err := cmd.Flags().GetString(tarballTargetFlag)
	if err != nil {
		return fmt.Errorf("read --target flag: %w", err)
	}
	dryRun, err := cmd.Flags().GetBool(tarballDryRunFlag)
	if err != nil {
		return fmt.Errorf("read --dry-run flag: %w", err)
	}

	paths := make([]string, 0, len(args))
	if len(args) > 0 {
		for _, source := range args {
			absolute, err := filepath.Abs(source)
			if err != nil {
				return fmt.Errorf("resolve source path %q: %w", source, err)
			}
			paths = append(paths, filepath.Clean(absolute))
		}
	} else {
		if inputPath == "" {
			return fmt.Errorf("--input is required when no backup paths are supplied")
		}
		input, err := tb.LoadInput(cmd.Context(), inputPath)
		if err != nil {
			return fmt.Errorf("load backup input: %w", err)
		}
		paths = input.Paths
	}

	home, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("resolve home directory: %w", err)
	}
	result, backupErr := d.Service.Backup(cmd.Context(), tb.BackupOptions{
		Paths:  paths,
		Target: target,
		Home:   home,
		DryRun: dryRun,
	})
	if err := printBackupResult(cmd, result, dryRun); err != nil {
		return err
	}
	if backupErr != nil {
		return fmt.Errorf("backup: %w", backupErr)
	}
	return nil
}

func (d *TarballCommand) runRestore(cmd *ca.Command, args []string) error {
	if cmd.Context() == nil {
		return fmt.Errorf("nil context")
	}
	inputPath, err := cmd.Flags().GetString(tarballInputFlag)
	if err != nil {
		return fmt.Errorf("read --input flag: %w", err)
	}
	target, err := cmd.Flags().GetString(tarballTargetFlag)
	if err != nil {
		return fmt.Errorf("read --target flag: %w", err)
	}
	destination, err := cmd.Flags().GetString(tarballDestinationFlag)
	if err != nil {
		return fmt.Errorf("read --dest flag: %w", err)
	}
	overwrite, err := cmd.Flags().GetBool(tarballOverwriteFlag)
	if err != nil {
		return fmt.Errorf("read --overwrite flag: %w", err)
	}
	allowUnsafeLinks, err := cmd.Flags().GetBool(tarballAllowUnsafeLinksFlag)
	if err != nil {
		return fmt.Errorf("read --allow-unsafe-links flag: %w", err)
	}
	preserveOwner, err := cmd.Flags().GetBool(tarballPreserveOwnerFlag)
	if err != nil {
		return fmt.Errorf("read --preserve-owner flag: %w", err)
	}
	dryRun, err := cmd.Flags().GetBool(tarballDryRunFlag)
	if err != nil {
		return fmt.Errorf("read --dry-run flag: %w", err)
	}

	options := tb.RestoreOptions{
		Destination:      destination,
		Overwrite:        overwrite,
		AllowUnsafeLinks: allowUnsafeLinks,
		PreserveOwner:    preserveOwner,
		DryRun:           dryRun,
	}
	if len(args) == 1 {
		options.Archive = args[0]
		result, err := d.restoreArchive(cmd.Context(), options)
		if err != nil {
			return fmt.Errorf("restore %s: %w", args[0], err)
		}
		return printRestoreResult(cmd, result, dryRun)
	}
	if inputPath == "" {
		return fmt.Errorf("--input is required when no archive is supplied")
	}
	input, err := tb.LoadInput(cmd.Context(), inputPath)
	if err != nil {
		return fmt.Errorf("load restore input: %w", err)
	}
	candidates, err := d.Service.DiscoverArchives(cmd.Context(), input.Paths, target)
	if err != nil {
		return fmt.Errorf("discover restore archives: %w", err)
	}

	var failures []error
	for _, candidateSet := range candidates {
		if len(candidateSet.Archives) == 0 {
			err := fmt.Errorf("no archive for %s in %s", candidateSet.Source, target)
			_, _ = fmt.Fprintf(cmd.ErrOrStderr(), "warning: %v\n", err)
			failures = append(failures, err)
			continue
		}

		archivePath := candidateSet.Archives[0]
		if len(candidateSet.Archives) > 1 {
			archivePath, err = d.pickArchive(cmd.Context(), candidateSet.Archives)
			if err != nil {
				wrapped := fmt.Errorf("select archive for %s: %w", candidateSet.Source, err)
				_, _ = fmt.Fprintf(cmd.ErrOrStderr(), "error: %v\n", wrapped)
				failures = append(failures, wrapped)
				continue
			}
		}

		options.Archive = archivePath
		result, restoreErr := d.restoreArchive(cmd.Context(), options)
		if restoreErr != nil {
			wrapped := fmt.Errorf("restore %s: %w", archivePath, restoreErr)
			_, _ = fmt.Fprintf(cmd.ErrOrStderr(), "error: %v\n", wrapped)
			failures = append(failures, wrapped)
			continue
		}
		if err := printRestoreResult(cmd, result, dryRun); err != nil {
			failures = append(failures, fmt.Errorf("print restore result: %w", err))
		}
	}
	return errors.Join(failures...)
}

func defaultTarballTarget() string {
	if target := os.Getenv("TGZ_TARGET"); target != "" {
		return target
	}
	dataHome := os.Getenv("XDG_DATA_HOME")
	if dataHome == "" {
		home, err := os.UserHomeDir()
		if err != nil {
			return ""
		}
		dataHome = filepath.Join(home, ".local", "share")
	}
	return filepath.Join(dataHome, "tarballs")
}

func addTarballTargetFlag(cmd *ca.Command, usage string) {
	cmd.Flags().String(tarballTargetFlag, defaultTarballTarget(), usage)
	cmd.Flags().Lookup(tarballTargetFlag).DefValue = tarballTargetHelpDefault
}

func printBackupResult(cmd *ca.Command, result tb.BackupResult, dryRun bool) error {
	for _, artifacts := range result.Artifacts {
		if dryRun {
			for _, path := range []string{artifacts.ArchivePath, artifacts.ChecksumPath, artifacts.ManifestPath} {
				if _, err := fmt.Fprintf(cmd.OutOrStdout(), "would create: %s\n", path); err != nil {
					return fmt.Errorf("write backup plan: %w", err)
				}
			}
			continue
		}
		if _, err := fmt.Fprintf(cmd.OutOrStdout(), "backed up: %s\n", artifacts.ArchivePath); err != nil {
			return fmt.Errorf("write backup result: %w", err)
		}
	}
	return nil
}

func printRestoreResult(cmd *ca.Command, result tb.RestoreResult, dryRun bool) error {
	verb := "restored"
	if dryRun {
		verb = "would restore"
	}
	if _, err := fmt.Fprintf(
		cmd.OutOrStdout(),
		"%s: %s -> %s\n",
		verb,
		result.ArchivePath,
		result.Destination,
	); err != nil {
		return fmt.Errorf("write restore result: %w", err)
	}
	if dryRun {
		for _, path := range result.Paths {
			destinationPath := filepath.Join(result.Destination, filepath.FromSlash(path))
			if _, err := fmt.Fprintf(cmd.OutOrStdout(), "  %s\n", destinationPath); err != nil {
				return fmt.Errorf("write restore plan: %w", err)
			}
		}
	}
	return nil
}

func newTarballDeps(logger lg.Logger) *TarballCommand {
	return &TarballCommand{Log: logger, Service: tb.NewService(logger), pickArchive: selectArchive}
}
