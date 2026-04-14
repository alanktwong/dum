package cmd

import (
	f "awong/dotfiles/internal/factory"
	pb "awong/dotfiles/internal/playbook"
	"context"
	"fmt"
	"strings"

	"github.com/spf13/cobra"
)

// listExecutor lists a playbook.
type listExecutor interface {
	List(ctx context.Context, input *pb.Input) (*pb.Result, error)
}

// defaultListExecutor wraps pb.NewExecutor().List().
type defaultListExecutor struct{}

func (d *defaultListExecutor) List(ctx context.Context, input *pb.Input) (*pb.Result, error) {
	result, err := pb.NewExecutor().List(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("executor list: %w", err)
	}
	return result, nil
}

// NewListCommand provides a command that lists plays and task.
func NewListCommand(rootUse string, dum *Dum) *cobra.Command {
	use := "list"
	alias := "ls"
	cmd := &cobra.Command{
		Use:     use,
		Aliases: []string{alias},
		Short:   fmt.Sprintf("%v-%v (or %v) lists plays and tasks", rootUse, use, alias),
		Long: strings.Join([]string{
			fmt.Sprintf("%v-%v (or %v) lists plays and tasks for software installations and configurations.", rootUse, use, alias),
			"",
			"A playbook is a group of plays specified by a YAML file. A play is a group of tasks.",
			"A task is specific work to install or configure something on the computer.",
			"",
			"Task types:",
			"  dir       - runs 'mkdir -p {ID}'",
			"  link      - runs 'ln -s {ID} {target}'",
			"  git       - runs 'git clone {ID} {target}'",
			"  brew      - runs 'brew install {ID}'",
			"  cask      - runs 'brew install --cask {ID}'",
			"  cellar    - runs 'brew cellar'",
			"  bash      - runs a bash command or script",
			"  vscode    - runs 'code --install-extension {ID}'",
			"  mas       - runs 'mas install {ID}'",
			"  jetbrains - runs '{app} installPlugin {ID}' for JetBrains apps",
			"  function  - runs a custom function from the playbook",
			"",
			"Config:",
			"  Default:   $XDG_CONFIG_HOME/dum/installer.yml (or ~/.config/dum/installer.yml)",
			"  Override:  --file flag or INSTALLER_CONFIG env var",
			"",
			"Logging:",
			"  Set ZSH_LOG_LEVEL env var to 'info' or 'debug', or use --verbose flag",
		}, "\n"),
		Example: strings.Join([]string{
			"  # List all plays and tasks from default config",
			fmt.Sprintf("  %v %v", rootUse, use),
			"",
			"  # List with verbose output",
			fmt.Sprintf("  %v %v -v", rootUse, use),
			"",
			"  # List a specific group of plays",
			fmt.Sprintf("  %v %v -v --group work", rootUse, use),
			"",
			"  # Use a custom config file",
			fmt.Sprintf("  %v %v -v --file ~/my-install.yml", rootUse, use),
		}, "\n"),
		RunE: func(cmd *cobra.Command, _ []string) error {
			return dum.runList(cmd)
		},
	}

	addVerboseFlag(cmd)
	addFileFlag(cmd)
	addGroupFlag(cmd)
	return cmd
}

func (d *Dum) runList(cmd *cobra.Command) error {
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
		dryrun = false
	}
	file, err := cmd.Flags().GetString(FILE)
	if err != nil {
		return fmt.Errorf("error getting %v flag: %w", FILE, err)
	}
	d.Log.Debug("Running list command:",
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
		return fmt.Errorf("error providing playbook from file %s: %w", file, err)
	}

	_, err = d.ListExecutor.List(ctx, input)
	if err != nil {
		return fmt.Errorf("error listing config file %s: %w", file, err)
	}
	return nil
}
