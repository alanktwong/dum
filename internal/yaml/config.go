// Package yaml provides typed structs for parsing and validating
// the installer.yml configuration file.

//go:generate go run -tags gen_schema ./gen_schema/main.go
package yaml

import (
	"fmt"

	tyg "awong/dotfiles/internal/types/gen"
)

// PlayBookYAML represents the top-level playbook configuration.
type PlayBookYAML struct {
	ID          string              `yaml:"id"`
	Description string              `yaml:"description"`
	Enabled     bool                `yaml:"enabled"`
	Sudo        bool                `yaml:"sudo"`
	JetBrains   []map[string]string `yaml:"jetbrains"`
	Plays       []PlayYAML          `yaml:"plays"`
}

// PlayYAML represents a single play within a playbook.
type PlayYAML struct {
	ID          string     `yaml:"id"`
	Description string     `yaml:"description"`
	Enabled     bool       `yaml:"enabled"`
	Tasks       []TaskYAML `yaml:"tasks"`
}

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
