package yaml

import (
	"fmt"
)

// PlayYAML represents a single play within a playbook.
type PlayYAML struct {
	ID          string     `yaml:"play"`
	Description string     `yaml:"description"`
	Enabled     bool       `yaml:"enabled"`
	Tasks       []TaskYAML `yaml:"tasks"`
}

// Validate checks that the PlayYAML is valid.
func (p *PlayYAML) Validate() error {
	if p.ID == "" {
		return fmt.Errorf("play id cannot be empty")
	}
	for _, task := range p.Tasks {
		if err := task.Validate(); err != nil {
			return fmt.Errorf("task %s: %w", task.ID, err)
		}
	}
	return nil
}
