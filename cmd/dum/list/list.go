// Package list defines the list Cobra command and its typed dependencies.
package list

import (
	"awong/dotfiles/cmd/dum/cli"
	"context"
	"fmt"
	"strings"

	fy "awong/dotfiles/internal/factory"
	lg "awong/dotfiles/internal/logging"
	pb "awong/dotfiles/internal/playbook"

	"github.com/spf13/cobra"
)

// CLI flag names used by the list command.
const (
	DRYRUN  = cli.DRYRUN
	FILE    = cli.FILE
	GROUP   = cli.GROUP
	VERBOSE = cli.VERBOSE
)

// FactoryProvider builds playbook input from installer options.
type FactoryProvider interface {
	Provide(fy.InputOptions) (*pb.Input, error)
}

// Executor lists a prepared playbook.
type Executor interface {
	List(context.Context, *pb.Input) (*pb.Result, error)
}

// Command holds dependencies for the list command.
type Command struct {
	Log             lg.Logger
	FactoryProvider FactoryProvider
	Executor        Executor
}

// NewCommand constructs a list command dependency set.
func NewCommand(logger lg.Logger) *Command {
	return &Command{Log: logger, FactoryProvider: &defaultFactoryProvider{}, Executor: &defaultListExecutor{}}
}

type defaultFactoryProvider struct{}

func (*defaultFactoryProvider) Provide(opts fy.InputOptions) (*pb.Input, error) {
	input, err := fy.NewFactory().Provide(opts)
	if err != nil {
		return nil, fmt.Errorf("factory provide: %w", err)
	}
	return input, nil
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
func NewListCommand(rootUse string, dum *Command) *cobra.Command {
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
	cli.AddVerboseFlag(cmd)
	cli.AddFileFlag(cmd)
	cli.AddGroupFlag(cmd)
	return cmd
}

func (d *Command) runList(cmd *cobra.Command) error {
	ctx := cmd.Context()
	if ctx == nil {
		return fmt.Errorf("nil context")
	}
	verbosity := cli.GetVerbosityFromCommand(cmd)
	d.Log.SetLevel(verbosity.Level())
	groupName, err := cmd.Flags().GetString(cli.GROUP)
	if err != nil {
		groupName = ""
	}
	config, err := cmd.Flags().GetString(cli.FILE)
	if err != nil {
		return fmt.Errorf("error getting %v flag: %w", cli.FILE, err)
	}
	d.Log.Debug(
		"Running list command:",
		cli.DRYRUN,
		false,
		cli.VERBOSE,
		verbosity.Verbose,
		cli.GROUP,
		groupName,
		cli.FILE,
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
