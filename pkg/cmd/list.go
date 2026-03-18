package cmd

import (
	pb "awong/dotfiles/pkg/playbook"
	"fmt"

	"github.com/spf13/cobra"
)

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
			ctx := cmd.Context()
			if ctx == nil {
				return fmt.Errorf("nil context")
			}
			verbosity := GetVerbosityFromCommand(cmd)
			dum.Log.SetLevel(verbosity.Level())

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
			dum.Log.Debug("Running list command:",
				DRYRUN, dryrun,
				VERBOSE, verbosity.Verbose,
				GROUP, group,
				FILE, file)

			factory := pb.NewFactory()
			input, err := factory.Provide(pb.InputOptions{
				File:   file,
				Group:  group,
				DryRun: dryrun,
			})
			if err != nil {
				return fmt.Errorf("error providing playbook from file %s: %w", file, err)
			}

			executor := pb.NewExecutor()
			_, err = executor.List(ctx, input)
			if err != nil {
				return fmt.Errorf("error listing config file %s: %w", file, err)
			}
			return nil
		},
	}

	addVerboseFlag(cmd)
	addFileFlag(cmd)
	addGroupFlag(cmd)
	return cmd
}
