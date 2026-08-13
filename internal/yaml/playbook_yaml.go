package yaml

import (
	"fmt"
)

// PlayBookYAML represents the top-level playbook configuration.
type PlayBookYAML struct {
	ID          string              `yaml:"book"`
	Description string              `yaml:"description"`
	Enabled     bool                `yaml:"enabled"`
	Sudo        bool                `yaml:"sudo"`
	JetBrains   []map[string]string `yaml:"jetbrains"`
	Plays       []PlayYAML          `yaml:"plays"`
}

// Validate checks that the PlayBookYAML is valid.
func (p *PlayBookYAML) Validate() error {
	if p.ID == "" {
		return fmt.Errorf("playbook id cannot be empty")
	}
	for _, play := range p.Plays {
		if err := play.Validate(); err != nil {
			return fmt.Errorf("play %s: %w", play.ID, err)
		}
	}
	return nil
}
