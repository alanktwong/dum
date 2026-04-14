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
				"",
				"Configuring the installer.yml:",
				"  # dir - create a directory",
				"  - type: \"dir\"",
				"    id: \"~/.dotfiles\"",
				"    description: \"Create dotfiles directory\"",
				"",
				"  # link - create a symbolic link",
				"  - type: \"link\"",
				"    id: \"../projects/dotfiles\"",
				"    root: \"~/.dotfiles\"",
				"    description: \"Link dotfiles folder\"",
				"",
				"  # git - clone a repository",
				"  - type: \"git\"",
				"    id: \"https://github.com/tmux-plugins/tpm.git\"",
				"    description: \"Clone TPM\"",
				"",
				"  # brew - install a Homebrew formula",
				"  - type: \"brew\"",
				"    id: \"gh\"",
				"    description: \"Install GitHub CLI\"",
				"",
				"  # cask - install a Homebrew cask",
				"  - type: \"cask\"",
				"    id: \"visual-studio-code\"",
				"    description: \"Install VS Code\"",
				"",
				"  # cellar - install a Homebrew cellar",
				"  - type: \"cellar\"",
				"    id: \"boost\"",
				"    description: \"Install Boost\"",
				"",
				"  # bash - run a bash command",
				"  - type: \"bash\"",
				"    id: \"hello\"",
				"    command: \"echo 'Hello, World!'\"",
				"    description: \"Print greeting\"",
				"",
				"  # bash - run a bash script",
				"  - type: \"bash\"",
				"    id: \"setup-script\"",
				"    script: |",
				"      echo 'Running setup...'",
				"      ./configure && make",
				"    description: \"Run setup script\"",
				"",
				"  # vscode - install a VS Code extension",
				"  - type: \"vscode\"",
				"    id: \"vscodevim.vim\"",
				"    description: \"Install Vim extension\"",
				"",
				"  # mas - install a Mac App Store app",
				"  - type: \"mas\"",
				"    id: \"462058435\"",
				"    description: \"Install Microsoft Excel\"",
				"",
				"  # jetbrains - install a JetBrains plugin",
				"  - type: \"jetbrains\"",
				"    id: \"org.asciidoctor.intellij.asciidoc\"",
				"    apps: [\"goland\", \"idea\"]",
				"    description: \"Install AsciiDoc plugin\"",
				"",
				"  # function - call a custom function",
				"  - type: \"function\"",
				"    id: \"my_custom_function\"",
				"    description: \"Run custom function\"",
			},
			"\n"),
		Example: strings.Join([]string{
			"  # Run all plays and tasks from default config",
			fmt.Sprintf("  %v %v", rootUse, use),
			"",
			"  # Run with verbose output",
			fmt.Sprintf("  %v %v -v", rootUse, use),
			"",
			"  # Dry run (preview what would happen)",
			fmt.Sprintf("  %v %v --dry-run", rootUse, use),
			"",
			"  # Run a specific group of plays",
			fmt.Sprintf("  %v %v --group work", rootUse, use),
			"",
			"  # Use a custom config file",
			fmt.Sprintf("  %v %v --file ~/my-install.yml", rootUse, use),
			"",
			"  # Combine flags",
			fmt.Sprintf("  %v %v --group work -vv --dry-run", rootUse, use),
		}, "\n"),
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
