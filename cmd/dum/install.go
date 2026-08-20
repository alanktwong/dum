package main

import (
	"context"
	"fmt"
	"strings"

	fy "github.com/alanktwong/dum/internal/factory"
	lg "github.com/alanktwong/dum/internal/logging"
	pb "github.com/alanktwong/dum/internal/playbook"

	"github.com/spf13/cobra"
)

// InstallCommand holds dependencies for the install command.
type InstallCommand struct {
	Log             lg.Logger
	FactoryProvider FactoryProvider
	Executor        InstallExecutor
}

// newInstallDeps constructs an install command dependency set.
func newInstallDeps(logger lg.Logger) *InstallCommand {
	return &InstallCommand{Log: logger, FactoryProvider: &defaultFactoryProvider{}, Executor: &defaultInstallExecutor{}}
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
func NewInstallCommand(rootUse string, dum *InstallCommand) *cobra.Command {
	use, alias := "install", "i"
	longHunk := []string{
		fmt.Sprintf(
			"%v-%v (or %v) runs plays and tasks for software installations and configurations.",
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
		"",
		"Config:",
		"  Default:   ./dum.yml, $XDG_CONFIG_HOME/dum/dum.yml, or legacy $XDG_CONFIG_HOME/dum/installer.yml",
		"  Override:  --file flag or DUM_CONFIG env var",
		"",
		"Logging:",
		"  Set ZSH_LOG_LEVEL env var to 'info' or 'debug', or use --verbose flag",
		"",
		"Dry run:",
		"  Use --dryrun to preview what would run without making changes.",
		"  Dry run output includes disabled plays and tasks, showing the full playbook manifest.",
		"",
		"Configuring dum.yml:",
		"  # dir - create a directory",
		"  - type: \"dir\"",
		"    task: \"~/.dotfiles\"",
		"    description: \"Create dotfiles directory\"",
		"",
		"  # link - create a symbolic link",
		"  - type: \"link\"",
		"    task: \"../projects/dotfiles\"",
		"    root: \"~/.dotfiles\"",
		"    description: \"Link dotfiles folder\"",
		"",
		"  # git - clone a repository",
		"  - type: \"git\"",
		"    task: \"https://github.com/tmux-plugins/tpm.git\"",
		"    description: \"Clone TPM\"",
		"",
		"  # brew - install a Homebrew formula",
		"  - type: \"brew\"",
		"    task: \"gh\"",
		"    description: \"Install GitHub CLI\"",
		"",
		"  # cask - install a Homebrew cask",
		"  - type: \"cask\"",
		"    task: \"visual-studio-code\"",
		"    description: \"Install VS Code\"",
		"",
		"  # cellar - install a Homebrew cellar",
		"  - type: \"cellar\"",
		"    task: \"boost\"",
		"    description: \"Install Boost\"",
		"",
		"  # bash - run a bash command",
		"  - type: \"bash\"",
		"    task: \"hello\"",
		"    command: \"echo 'Hello, World!'\"",
		"    description: \"Print greeting\"",
		"",
		"  # bash - run a bash script",
		"  - type: \"bash\"",
		"    task: \"setup-script\"",
		"    script: |",
		"      echo 'Running setup...'",
		"      ./configure && make",
		"    description: \"Run setup script\"",
		"",
		"  # vscode - install a VS Code extension",
		"  - type: \"vscode\"",
		"    task: \"vscodevim.vim\"",
		"    description: \"Install Vim extension\"",
		"",
		"  # mas - install a Mac App Store app",
		"  - type: \"mas\"",
		"    task: \"462058435\"",
		"    description: \"Install Microsoft Excel\"",
		"",
		"  # jetbrains - install a JetBrains plugin",
		"  - type: \"jetbrains\"",
		"    task: \"org.asciidoctor.intellij.asciidoc\"",
		"    apps: [\"goland\", \"idea\"]",
		"    description: \"Install AsciiDoc plugin\"",
	}
	exampleHunk := []string{
		"  # Run all plays and tasks from default config",
		fmt.Sprintf("  %v %v", rootUse, use),
		"",
		"  # Run with verbose output",
		fmt.Sprintf("  %v %v -v", rootUse, use),
		"", "  # Dry run (preview what would run, including disabled plays and tasks)",
		fmt.Sprintf("  %v %v --dryrun", rootUse, use),
		"",
		"  # Run a specific group of plays",
		fmt.Sprintf("  %v %v --group work", rootUse, use),
		"",
		"  # Use a custom config file",
		fmt.Sprintf("  %v %v --file ~/my-install.yml", rootUse, use),
		"",
		"  # Combine flags",
		fmt.Sprintf("  %v %v --group work --dryrun -v", rootUse, use),
	}
	cmd := &cobra.Command{
		Use:     use,
		Aliases: []string{alias},
		Short:   fmt.Sprintf("%v-%v (or %v) runs plays and tasks", rootUse, use, alias),
		Long:    strings.Join(longHunk, "\n"),
		Example: strings.Join(exampleHunk, "\n"),
		RunE:    func(cmd *cobra.Command, _ []string) error { return dum.runInstall(cmd) },
	}
	AddVerboseFlag(cmd)
	AddFileFlag(cmd)
	AddDryRunFlag(cmd)
	AddGroupFlag(cmd)
	return cmd
}

func (d *InstallCommand) runInstall(cmd *cobra.Command) error {
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
	dry, err := cmd.Flags().GetBool(DRYRUN)
	if err != nil {
		return fmt.Errorf("error getting %v flag: %w", DRYRUN, err)
	}
	config, err := cmd.Flags().GetString(FILE)
	if err != nil {
		return fmt.Errorf("error getting %v flag: %w", FILE, err)
	}
	d.Log.Debug(
		"Running install command:",
		DRYRUN,
		dry,
		VERBOSE,
		verbosity.Verbose,
		GROUP,
		groupName,
		FILE,
		config,
	)
	input, err := d.FactoryProvider.Provide(fy.InputOptions{File: config, Group: groupName, DryRun: dry})
	if err != nil {
		return fmt.Errorf("error providing context from file %s: %w", config, err)
	}
	if _, err = d.Executor.Install(ctx, input); err != nil {
		return fmt.Errorf("error installing config file %s: %w", config, err)
	}
	return nil
}
