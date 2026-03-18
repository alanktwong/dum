package cmd

import (
	pb "awong/dotfiles/pkg/playbook"
	"fmt"
	"strings"

	"github.com/spf13/cobra"
)

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
			ctx := cmd.Context()
			if ctx == nil {
				return fmt.Errorf("nil context")
			}
			verbosity := GetVerbosityFromCommand(cmd)
			dum.Log.SetLevel(verbosity.Level())

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
			dum.Log.Debug("Running install command:",
				DRYRUN, dryrun,
				VERBOSE, verbosity.Verbose,
				GROUP, group,
				FILE, file)

			factory := pb.NewFactory()
			input, err := factory.Provide(pb.InputOptions{
				File:   file,
				Group:  group,
				DryRun: dryrun,
			})
			if err != nil {
				return fmt.Errorf("error providing context from file %s: %w", file, err)
			}

			executor := pb.NewExecutor()
			_, err = executor.Install(ctx, input)
			if err != nil {
				return fmt.Errorf("error installing config file %s: %w", file, err)
			}
			return nil
		},
	}
	addVerboseFlag(cmd)
	addFileFlag(cmd)
	addDryRunFlag(cmd)
	addGroupFlag(cmd)
	return cmd
}
