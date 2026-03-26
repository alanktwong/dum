package cmd

import (
	f "awong/dotfiles/internal/factory"
	pb "awong/dotfiles/internal/playbook"
	"context"
	"fmt"
	"strings"

	"github.com/spf13/cobra"
)

// factoryProvider provides Input from options.
type factoryProvider interface {
	Provide(opts f.InputOptions) (*pb.Input, error)
}

// installExecutor installs a playbook.
type installExecutor interface {
	Install(ctx context.Context, input *pb.Input) (*pb.Result, error)
}

// defaultFactoryProvider wraps f.NewFactory().Provide().
type defaultFactoryProvider struct{}

func (d *defaultFactoryProvider) Provide(opts f.InputOptions) (*pb.Input, error) {
	input, err := f.NewFactory().Provide(opts)
	if err != nil {
		return nil, fmt.Errorf("factory provide: %w", err)
	}
	return input, nil
}

// defaultInstallExecutor wraps pb.NewExecutor().Install().
type defaultInstallExecutor struct{}

func (d *defaultInstallExecutor) Install(ctx context.Context, input *pb.Input) (*pb.Result, error) {
	result, err := pb.NewExecutor().Install(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("executor install: %w", err)
	}
	return result, nil
}

// NewInstallCommand provides a command that installs plays and task.
func NewInstallCommand(rootUse string, dum *Dum) *cobra.Command {
	use := "install"
	alias := "i"
	cmd := &cobra.Command{
		Use:     use,
		Aliases: []string{alias},
		Short:   fmt.Sprintf("%v-%v (or %v) runs plays and tasks", rootUse, use, alias),
		Long: strings.Join(
			[]string{
				fmt.Sprintf("%v-%v (or %v) runs plays and tasks for software installations and configurations.", rootUse, use, alias),
				"A playbook is a group of plays specified by a YAML file. A play is a group of tasks.",
				"A task is specific work that can be done to install or configure something on the computer.",
				"Tasks have types and an ID.",
				"The 'dir' task runs 'mkdir -p {ID}'",
				"The 'brew', 'cellar' and 'cask' task types runs 'brew install' with the formula specified by its ID.",
				"The 'link' task runs 'ln -s {ID}'",
				"The 'git' task runs 'git clone {ID}'",
				"The 'vscode' task runs 'code --install-extension {ID}'",
				"The 'mas' task runs the Apple Store CLI: 'mas install {ID}'",
				"The 'jetbrains' task installs JetBrains plugins by its ID given a collection of apps such as `idea`.",
				"Each app must be executable from the command line as the install command will run `idea installPlugin`.",
				"",
				"A playbook file by default is at '~/.config/installer.yml`.",
				"This command can load the file by its 'file' flag or by environment variable INSTALLER_CONFIG.",
				"",
				"The log level is initialized by an environment variable but overridden by by the verbose flag.",
				"The environment log level is the ZSH_LOG_LEVEL.",
				"ZSH_LOG_LEVEL can be set to 'info' or 'debug'.",
			},
			"\n"),
		RunE: func(cmd *cobra.Command, _ []string) error {
			return dum.runInstall(cmd)
		},
	}
	addVerboseFlag(cmd)
	addFileFlag(cmd)
	addDryRunFlag(cmd)
	addGroupFlag(cmd)
	return cmd
}

func (d *Dum) runInstall(cmd *cobra.Command) error {
	ctx := cmd.Context()
	if ctx == nil {
		return fmt.Errorf("nil context")
	}
	verbosity := GetVerbosityFromCommand(cmd)
	d.Log.SetLevel(verbosity.Level())

	group, err := cmd.Flags().GetString(GROUP)
	if err != nil {
		group = ""
	}
	dryrun, err := cmd.Flags().GetBool(DRYRUN)
	if err != nil {
		return fmt.Errorf("error getting %v flag: %w", DRYRUN, err)
	}
	file, err := cmd.Flags().GetString(FILE)
	if err != nil {
		return fmt.Errorf("error getting %v flag: %w", FILE, err)
	}
	d.Log.Debug("Running install command:",
		DRYRUN, dryrun,
		VERBOSE, verbosity.Verbose,
		GROUP, group,
		FILE, file)

	input, err := d.FactoryProvider.Provide(f.InputOptions{
		File:   file,
		Group:  group,
		DryRun: dryrun,
	})
	if err != nil {
		return fmt.Errorf("error providing context from file %s: %w", file, err)
	}

	_, err = d.InstallExecutor.Install(ctx, input)
	if err != nil {
		return fmt.Errorf("error installing config file %s: %w", file, err)
	}
	return nil
}
