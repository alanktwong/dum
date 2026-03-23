package playbook

import (
	"fmt"
)

// Attributable is the provider interface to get attributes.
type Attributable interface {
	GetAttributes() Attributes
}

// Attributes is an embedded struct common to playbooks, plays and tasks.
type Attributes struct {
	ID          string
	Description string
	Enabled     bool
	Sudo        bool
	Command     string
	Script      string
}

// Identifiable is the provider interface to get an id.
type Identifiable interface {
	GetID() string
}

// Enableable is the provider interface to get the enabled flag.
type Enableable interface {
	IsEnabled() bool
}

// NewAttributes constructs Attributes.
func NewAttributes(id, description string, enabled, sudo bool) (*Attributes, error) {
	if id == "" {
		return nil, fmt.Errorf("attribute ID cannot be empty")
	}
	return &Attributes{
		ID:          id,
		Description: description,
		Enabled:     enabled,
		Sudo:        sudo, // Default to false, can be overridden
	}, nil
}

// IsEnabled implements Enableable.
func (a *Attributes) IsEnabled() bool {
	return a.Enabled
}

// GetID implements Identifiable.
func (a *Attributes) GetID() string {
	return a.ID
}
