package cmd

import (
	f "awong/dotfiles/internal/factory"
	pb "awong/dotfiles/internal/playbook"
	"context"
	"fmt"

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
		Long:    fmt.Sprintf("%v-%v (or %v) lists plays and tasks for software installations and configurations.", rootUse, use, alias),
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
