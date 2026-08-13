package main

import (
	"awong/dotfiles/cmd/dum/install"
	"awong/dotfiles/cmd/dum/list"
	"awong/dotfiles/cmd/dum/rename"
	"awong/dotfiles/cmd/dum/schema"
	"fmt"

	logcmd "awong/dotfiles/cmd/dum/log"

	lg "awong/dotfiles/internal/logging"

	clog "github.com/charmbracelet/log"
	"github.com/spf13/cobra"
)

var (
	version = "dev"
	commit  = ""
)

// Dum is the executable application for the dum CLI.
type Dum struct {
	Log             lg.Logger
	Cmd             *cobra.Command
	FactoryProvider install.FactoryProvider
	InstallExecutor install.Executor
	ListExecutor    list.Executor
}

// NewDum constructs the root command and wires command dependencies.
func NewDum() *Dum {
	_ = commit
	rootUse := "dum"
	rootCmd := &cobra.Command{
		Use:     rootUse,
		Short:   fmt.Sprintf("%v is a command line tool for your utilities", rootUse),
		Long:    fmt.Sprintf("%v is a command line tool for your utilites such as managing software installations and configurations.", rootUse),
		Version: version,
	}
	logger := lg.NewLogger(lg.Options{Prefix: "", Level: clog.InfoLevel})
	installDum := install.NewCommand(logger)
	listDum := list.NewCommand(logger)
	logDum := logcmd.NewCommand(logger)
	renameDum := rename.NewCommand(logger)
	schemaDum := schema.NewCommand()

	rootCmd.AddCommand(list.NewListCommand(rootUse, listDum))
	rootCmd.AddCommand(install.NewInstallCommand(rootUse, installDum))
	rootCmd.AddCommand(rename.NewRenameCommand(rootUse, renameDum))
	rootCmd.AddCommand(logcmd.NewLogCommand(rootUse, logDum))
	rootCmd.AddCommand(schema.NewSchemaCommand(rootUse, schemaDum))

	return &Dum{
		Log:             logger,
		Cmd:             rootCmd,
		FactoryProvider: installDum.FactoryProvider,
		InstallExecutor: installDum.Executor,
		ListExecutor:    listDum.Executor,
	}
}

// Exec executes the root command and reports CLI errors through the logger.
func (d *Dum) Exec() {
	if err := d.Cmd.Execute(); err != nil {
		d.Log.Fatal(err)
	}
}
