package main

import (
	"context"
	"fmt"

	fy "awong/dotfiles/internal/factory"
	lg "awong/dotfiles/internal/logging"
	pb "awong/dotfiles/internal/playbook"

	"github.com/spf13/cobra"
)

var (
	version = "dev"
	commit  = ""
)

// FactoryProvider builds playbook input from installer options.
type FactoryProvider interface {
	Provide(fy.InputOptions) (*pb.Input, error)
}

// InstallExecutor runs a prepared playbook.
type InstallExecutor interface {
	Install(context.Context, *pb.Input) (*pb.Result, error)
}

type defaultFactoryProvider struct{}

func (*defaultFactoryProvider) Provide(opts fy.InputOptions) (*pb.Input, error) {
	input, err := fy.NewFactory().Provide(opts)
	if err != nil {
		return nil, fmt.Errorf("factory provide: %w", err)
	}
	return input, nil
}

// Dum is the executable application for the dum CLI.
type Dum struct {
	Log             lg.Logger
	Cmd             *cobra.Command
	FactoryProvider FactoryProvider
	InstallExecutor InstallExecutor
}

// NewDum constructs the root command and wires command dependencies.
func NewDum() *Dum {
	_ = commit
	rootUse := "dum"
	rootCmd := &cobra.Command{
		Use:   rootUse,
		Short: fmt.Sprintf("%v is a command line tool for your utilities", rootUse),
		Long: fmt.Sprintf(
			"%v is a command line tool for your utilites such as managing software installations and configurations.",
			rootUse,
		),
		Version: version,
	}
	logger := lg.NewLogger(lg.Options{Prefix: "", Level: GetLogLevel()})
	installDum := newInstallDeps(logger)
	logDum := newLogDeps(logger)
	renameDum := newRenameDeps(logger)
	schemaDum := newSchemaDeps()

	rootCmd.AddCommand(NewInstallCommand(rootUse, installDum))
	rootCmd.AddCommand(NewRenameCommand(rootUse, renameDum))
	rootCmd.AddCommand(NewLogCommand(rootUse, logDum))
	rootCmd.AddCommand(NewSchemaCommand(rootUse, schemaDum))

	return &Dum{
		Log:             logger,
		Cmd:             rootCmd,
		FactoryProvider: installDum.FactoryProvider,
		InstallExecutor: installDum.Executor,
	}
}

// Exec executes the root command and reports CLI errors through the logger.
func (d *Dum) Exec() {
	if err := d.Cmd.Execute(); err != nil {
		d.Log.Fatal(err)
	}
}
