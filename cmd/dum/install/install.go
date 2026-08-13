// Package install defines the install Cobra command and its typed dependencies.
package install

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

// CLI flag names used by the install command.
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

// Executor runs a prepared playbook.
type Executor interface {
	Install(context.Context, *pb.Input) (*pb.Result, error)
}

// Command holds dependencies for the install command.
type Command struct {
	Log             lg.Logger
	FactoryProvider FactoryProvider
	Executor        Executor
}

// NewCommand constructs an install command dependency set.
func NewCommand(logger lg.Logger) *Command {
	return &Command{Log: logger, FactoryProvider: &defaultFactoryProvider{}, Executor: &defaultInstallExecutor{}}
}

type defaultFactoryProvider struct{}

func (*defaultFactoryProvider) Provide(opts fy.InputOptions) (*pb.Input, error) {
	input, err := fy.NewFactory().Provide(opts)
	if err != nil {
		return nil, fmt.Errorf("factory provide: %w", err)
	}
	return input, nil
}

type defaultInstallExecutor struct{}

func (*defaultInstallExecutor) Install(ctx context.Context, input *pb.Input) (*pb.Result, error) {
	result, err := pb.NewExecutor().Install(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("executor install: %w", err)
	}
	return result, nil
}

// NewInstallCommand provides a command that installs plays and tasks.
func NewInstallCommand(rootUse string, dum *Command) *cobra.Command {
	use, alias := "install", "i"
	cmd := &cobra.Command{Use: use, Aliases: []string{alias}, Short: fmt.Sprintf("%v-%v (or %v) runs plays and tasks", rootUse, use, alias), Long: strings.Join([]string{
		fmt.Sprintf("%v-%v (or %v) runs plays and tasks for software installations and configurations.", rootUse, use, alias), "",
		"A playbook is a group of plays specified by a YAML file. A play is a group of tasks.", "A task is specific work to install or configure something on the computer.", "",
		"Task types:", "  dir       - runs 'mkdir -p {ID}'", "  link      - runs 'ln -s {ID} {target}'", "  git       - runs 'git clone {ID} {target}'", "  brew      - runs 'brew install {ID}'", "  cask      - runs 'brew install --cask {ID}'", "  cellar    - runs 'brew cellar'", "  bash      - runs a bash command or script", "  vscode    - runs 'code --install-extension {ID}'", "  mas       - runs 'mas install {ID}'", "  jetbrains - runs '{app} installPlugin {ID}' for JetBrains apps", "  function  - runs a custom function from the playbook", "",
		"Config:", "  Default:   $XDG_CONFIG_HOME/dum/installer.yml (or ~/.config/dum/installer.yml)", "  Override:  --file flag or INSTALLER_CONFIG env var", "",
		"Logging:", "  Set ZSH_LOG_LEVEL env var to 'info' or 'debug', or use --verbose flag", "",
		"Configuring the installer.yml:", "  # dir - create a directory", "  - type: \"dir\"", "    task: \"~/.dotfiles\"", "    description: \"Create dotfiles directory\"", "",
		"  # link - create a symbolic link", "  - type: \"link\"", "    task: \"../projects/dotfiles\"", "    root: \"~/.dotfiles\"", "    description: \"Link dotfiles folder\"", "",
		"  # git - clone a repository", "  - type: \"git\"", "    task: \"https://github.com/tmux-plugins/tpm.git\"", "    description: \"Clone TPM\"", "",
		"  # brew - install a Homebrew formula", "  - type: \"brew\"", "    task: \"gh\"", "    description: \"Install GitHub CLI\"", "",
		"  # cask - install a Homebrew cask", "  - type: \"cask\"", "    task: \"visual-studio-code\"", "    description: \"Install VS Code\"", "",
		"  # cellar - install a Homebrew cellar", "  - type: \"cellar\"", "    task: \"boost\"", "    description: \"Install Boost\"", "",
		"  # bash - run a bash command", "  - type: \"bash\"", "    task: \"hello\"", "    command: \"echo 'Hello, World!'\"", "    description: \"Print greeting\"", "",
		"  # bash - run a bash script", "  - type: \"bash\"", "    task: \"setup-script\"", "    script: |", "      echo 'Running setup...'", "      ./configure && make", "    description: \"Run setup script\"", "",
		"  # vscode - install a VS Code extension", "  - type: \"vscode\"", "    task: \"vscodevim.vim\"", "    description: \"Install Vim extension\"", "",
		"  # mas - install a Mac App Store app", "  - type: \"mas\"", "    task: \"462058435\"", "    description: \"Install Microsoft Excel\"", "",
		"  # jetbrains - install a JetBrains plugin", "  - type: \"jetbrains\"", "    task: \"org.asciidoctor.intellij.asciidoc\"", "    apps: [\"goland\", \"idea\"]", "    description: \"Install AsciiDoc plugin\"", "",
		"  # function - call a custom function", "  - type: \"function\"", "    task: \"my_custom_function\"", "    description: \"Run custom function\"",
	}, "\n"), Example: strings.Join([]string{
		"  # Run all plays and tasks from default config", fmt.Sprintf("  %v %v", rootUse, use), "", "  # Run with verbose output", fmt.Sprintf("  %v %v -v", rootUse, use), "", "  # Dry run (preview what would happen)", fmt.Sprintf("  %v %v --dryrun", rootUse, use), "", "  # Run a specific group of plays", fmt.Sprintf("  %v %v --group work", rootUse, use), "", "  # Use a custom config file", fmt.Sprintf("  %v %v --file ~/my-install.yml", rootUse, use), "", "  # Combine flags", fmt.Sprintf("  %v %v --group work --dryrun -v", rootUse, use)}, "\n"), RunE: func(cmd *cobra.Command, _ []string) error { return dum.runInstall(cmd) }}
	cli.AddVerboseFlag(cmd)
	cli.AddFileFlag(cmd)
	cli.AddDryRunFlag(cmd)
	cli.AddGroupFlag(cmd)
	return cmd
}

func (d *Command) runInstall(cmd *cobra.Command) error {
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
	dry, err := cmd.Flags().GetBool(cli.DRYRUN)
	if err != nil {
		return fmt.Errorf("error getting %v flag: %w", cli.DRYRUN, err)
	}
	config, err := cmd.Flags().GetString(cli.FILE)
	if err != nil {
		return fmt.Errorf("error getting %v flag: %w", cli.FILE, err)
	}
	d.Log.Debug("Running install command:", cli.DRYRUN, dry, cli.VERBOSE, verbosity.Verbose, cli.GROUP, groupName, cli.FILE, config)
	input, err := d.FactoryProvider.Provide(fy.InputOptions{File: config, Group: groupName, DryRun: dry})
	if err != nil {
		return fmt.Errorf("error providing context from file %s: %w", config, err)
	}
	if _, err = d.Executor.Install(ctx, input); err != nil {
		return fmt.Errorf("error installing config file %s: %w", config, err)
	}
	return nil
}
