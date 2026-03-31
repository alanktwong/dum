// Package yaml provides typed structs for parsing and validating
// the installer.yml configuration file.
package yaml

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

// Config represents the top-level configuration loaded from a YAML file.
type Config struct {
	PlayBook PlayBookYAML `yaml:"playbook"`
}

// Load reads and parses a YAML configuration file.
func Load(filePath string) (*Config, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}

	var config Config
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("failed to unmarshal YAML: %w", err)
	}

	if err := config.PlayBook.Validate(); err != nil {
		return nil, fmt.Errorf("validation failed: %w", err)
	}

	return &config, nil
}
