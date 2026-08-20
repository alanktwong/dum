package main

import (
	"context"
	"fmt"
	"strings"

	fy "awong/dotfiles/internal/factory"
	lg "awong/dotfiles/internal/logging"
	pb "awong/dotfiles/internal/playbook"

	"github.com/spf13/cobra"
)

// ListCommand holds dependencies for the list command.
type ListCommand struct {
	Log             lg.Logger
	FactoryProvider FactoryProvider
	Executor        ListExecutor
}

// newListDeps constructs a list command dependency set.
func newListDeps(logger lg.Logger) *ListCommand {
	return &ListCommand{Log: logger, FactoryProvider: &defaultFactoryProvider{}, Executor: &defaultListExecutor{}}
}

type defaultListExecutor struct{}

func (*defaultListExecutor) List(ctx context.Context, input *pb.Input) (*pb.Result, error) {
	result, err := pb.NewExecutor().List(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("executor list: %w", err)
	}
	return result, nil
}

// NewListCommand constructs the list Cobra command.
func NewListCommand(rootUse string, dum *ListCommand) *cobra.Command {
	use, alias := "list", "ls"
	cmd := &cobra.Command{
		Use:     use,
		Aliases: []string{alias},
		Short:   fmt.Sprintf("%v-%v (or %v) lists plays and tasks", rootUse, use, alias),
		Long: strings.Join([]string{
			fmt.Sprintf(
				"%v-%v (or %v) lists plays and tasks for software installations and configurations.",
				rootUse,
				use,
				alias,
			),
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
		Example: strings.Join(
			[]string{
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
			},
			"\n",
		),
		RunE: func(cmd *cobra.Command, _ []string) error { return dum.runList(cmd) },
	}
	AddVerboseFlag(cmd)
	AddFileFlag(cmd)
	AddGroupFlag(cmd)
	return cmd
}

func (d *ListCommand) runList(cmd *cobra.Command) error {
	ctx := cmd.Context()
	if ctx == nil {
		return fmt.Errorf("nil context")
	}
	verbosity := GetVerbosityFromCommand(cmd)
	d.Log.SetLevel(verbosity.Level())
	groupName, err := cmd.Flags().GetString(GROUP)
	if err != nil {
		groupName = ""
	}
	config, err := cmd.Flags().GetString(FILE)
	if err != nil {
		return fmt.Errorf("error getting %v flag: %w", FILE, err)
	}
	d.Log.Debug(
		"Running list command:",
		DRYRUN,
		false,
		VERBOSE,
		verbosity.Verbose,
		GROUP,
		groupName,
		FILE,
		config,
	)
	input, err := d.FactoryProvider.Provide(fy.InputOptions{File: config, Group: groupName})
	if err != nil {
		return fmt.Errorf("error providing playbook from file %s: %w", config, err)
	}
	if _, err = d.Executor.List(ctx, input); err != nil {
		return fmt.Errorf("error listing config file %s: %w", config, err)
	}
	return nil
}
