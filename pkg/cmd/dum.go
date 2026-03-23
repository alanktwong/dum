package cmd

import (
	l "awong/dotfiles/pkg/logging"
	"fmt"

	"github.com/charmbracelet/log"
	"github.com/spf13/cobra"
)

var version = "dev"

// Dum is the application for the dum CLI.
type Dum struct {
	Log l.Logger
	Cmd *cobra.Command
}

// NewDum constructs Dum.
func NewDum() *Dum {
	rootUse := "dum"
	rootCmd := &cobra.Command{
		Use:     rootUse,
		Short:   fmt.Sprintf("%v is a command line tool for your utilities", rootUse),
		Long:    fmt.Sprintf("%v is a command line tool for your utilites such as managing software installations and configurations.", rootUse),
		Version: version,
	}
	dum := &Dum{
		Log: l.NewLogger(l.Options{
			Prefix: "",
			Level:  log.InfoLevel,
		}),
		Cmd: rootCmd,
	}
	listCmd := NewListCommand(rootUse, dum)
	rootCmd.AddCommand(listCmd)

	installCmd := NewInstallCommand(rootUse, dum)
	rootCmd.AddCommand(installCmd)

	renameCmd := NewRenameCommand(rootUse, dum)
	rootCmd.AddCommand(renameCmd)

	logCmd := NewLogCommand(rootUse, dum)
	rootCmd.AddCommand(logCmd)

	return dum
}

// Exec is the main function for dum.
func (d *Dum) Exec() {
	err := d.Cmd.Execute()
	if err != nil {
		d.Log.Fatal(err)
	}
}
