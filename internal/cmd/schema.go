package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// NewSchemaCommand creates a new schema command for outputting the JSON schema.
func NewSchemaCommand(_ string, dum *Dum) *cobra.Command {
	use := "schema"
	alias := "s"
	cmd := &cobra.Command{
		Use:     use,
		Aliases: []string{alias},
		Short:   "Output or extract the JSON schema for installer.yml",
		Long: `Outputs the JSON schema for installer.yml configuration file

This can be used for IDE autocomplete and validation in editors like
VS Code and JetBrains IDEs.`,
		Example: `
  # Print schema to stdout
  dum schema

  # Save schema to a file
  dum schema --output installer.schema.json
`,
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
