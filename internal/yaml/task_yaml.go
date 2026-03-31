package yaml

import (
	"fmt"

	tyg "awong/dotfiles/internal/types/gen"
)

// TaskYAML represents a single task within a play.
type TaskYAML struct {
	Type        string   `yaml:"type"`
	ID          string   `yaml:"id"`
	Description string   `yaml:"description"`
	Enabled     bool     `yaml:"enabled"`
	Command     string   `yaml:"command,omitempty"`
	Script      string   `yaml:"script,omitempty"`
	Root        string   `yaml:"root,omitempty"`
	Target      string   `yaml:"target,omitempty"`
	Name        string   `yaml:"name,omitempty"`
	Tap         string   `yaml:"tap,omitempty"`
	Apps        []string `yaml:"apps,omitempty"`
}

// Validate checks that the TaskYAML is valid.
func (t *TaskYAML) Validate() error {
	if t.ID == "" {
		return fmt.Errorf("task id cannot be empty")
	}
	if t.Type == "" {
		return fmt.Errorf("task type cannot be empty")
	}
	if !tyg.TaskType(t.Type).IsValid() {
		return fmt.Errorf("invalid task type: %s", t.Type)
	}
	if t.Type == string(tyg.TaskTypeBash) && t.Command == "" && t.Script == "" {
		return fmt.Errorf("bash task must have either command or script")
	}
	return nil
}
