package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// NewSchemaCommand creates a new schema command for outputting the JSON schema.
func NewSchemaCommand(_ string, _ *Dum) *cobra.Command {
	use := "schema"
	cmd := &cobra.Command{
		Use:   use,
		Short: "Output or extract the JSON schema for installer.yml",
		Long: `Outputs the JSON schema for installer.yml configuration file.

This can be used for IDE autocomplete and validation in editors like
VS Code and JetBrains IDEs.

Examples:
  # Print schema to stdout
  dum schema

  # Save schema to a file
  dum schema --output installer.schema.json
`,
		RunE: func(cmd *cobra.Command, _ []string) error {
			outputPath, err := cmd.Flags().GetString("output")
			if err != nil {
				return fmt.Errorf("error getting output flag: %w", err)
			}

			if outputPath == "" {
				_, _ = os.Stdout.Write(schemaData)

				return nil
			}

			return os.WriteFile(outputPath, schemaData, 0o600)
		},
	}

	cmd.Flags().StringP("output", "o", "", "Output file path (default: stdout)")
	return cmd
}
