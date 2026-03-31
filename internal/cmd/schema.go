package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

// NewSchemaCommand creates a new schema command for outputting the JSON schema.
func NewSchemaCommand(rootUse string, dum *Dum) *cobra.Command {
	use := "schema"
	alias := "s"
	cmd := &cobra.Command{
		Use:     use,
		Aliases: []string{alias},
		Short:   fmt.Sprintf("%v-%v (or %v) output or extract the JSON schema for installer.yml", rootUse, use, alias),
		Long: strings.Join([]string{
			fmt.Sprintf("%v-%v (or %v) is a command line tool that outputs the JSON schema for installer.yml configuration.", rootUse, use, alias),
			"",
			"This can be used for IDE autocomplete and validation in editors like",
			"VS Code and JetBrains IDEs.",
			"",
			"For more info about playbooks, plays and tasks, see the help the install or list command.",
		}, "\n"),
		Example: strings.Join([]string{
			"# Print schema to stdout",
			fmt.Sprintf("%v %v", rootUse, use),
			"",
			"# Save schema to a file",
			fmt.Sprintf("%v %v --output installer.schema.json", rootUse, use),
		}, "\n"),
		RunE: func(cmd *cobra.Command, _ []string) error {
			return dum.runSchema(cmd)
		},
	}

	addOutputFlag(cmd)
	return cmd
}

func (d *Dum) runSchema(cmd *cobra.Command) error {
	// ctx := cmd.Context()
	outputPath, err := cmd.Flags().GetString("output")
	if err != nil {
		return fmt.Errorf("error getting output flag: %w", err)
	}

	if outputPath == "" {
		_, _ = os.Stdout.Write(schemaData)

		return nil
	}

	if err := os.WriteFile(outputPath, schemaData, 0o600); err != nil {
		return fmt.Errorf("error writing schema to %s: %w", outputPath, err)
	}

	return nil
}
