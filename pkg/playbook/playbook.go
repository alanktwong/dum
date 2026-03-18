// Package playbook contains the abstractions for installing and managing software
// configurations and installations via playbooks, plays and tasks.
package playbook

import (
	"context"
	"fmt"

	om "github.com/elliotchance/orderedmap/v3"
)

const (
	// Ellipsis is prefix used in logs and println.
	Ellipsis = "..."
	// PlayEllipsis is prefix used in logs and println.
	PlayEllipsis = "......."
	// TaskEllipsis is prefix used in logs and println.
	TaskEllipsis = "..........."
)

// PlayBook consists of a single install execution of many plays.
type PlayBook struct {
	Attributes
	Plays         []*Play
	JetBrainsApps map[string]string
}

// NewPlayBook constructs a PlayBook.
func NewPlayBook(attributes *Attributes, plays []*Play, apps map[string]string) (*PlayBook, error) {
	if attributes == nil {
		return nil, fmt.Errorf("attributes cannot be nil")
	}
	if attributes.ID == "" {
		return nil, fmt.Errorf("playbook ID cannot be empty")
	}
	return &PlayBook{
		Attributes:    *attributes,
		Plays:         plays,
		JetBrainsApps: apps,
	}, nil
}

// GetAttributes implements Attributable.
func (p *PlayBook) GetAttributes() Attributes {
	return p.Attributes
}

// GetID implements Identifiable.
func (p *PlayBook) GetID() string {
	return p.ID
}

// IsEnabled implements Enableable.
func (p *PlayBook) IsEnabled() bool {
	return p.Enabled
}

// GetPlays indexes all the Plays by its ID in an ordered map. If active is true it finds
// all the active plays. Otherwise, it finds them all.
func (p *PlayBook) GetPlays(active bool) *om.OrderedMap[string, *Play] {
	plays := om.NewOrderedMap[string, *Play]()
	if !active {
		for _, play := range p.Plays {
			plays.Set(play.GetID(), play)
		}
	} else if p.Enabled {
		for _, play := range p.Plays {
			if play.IsEnabled() {
				plays.Set(play.GetID(), play)
			}
		}
	}
	return plays
}

// sudoCtxKey is a private type key for the sudo value that may be in context.Context.
type sudoCtxKey struct{}

// WithSudo adds sudo to context.
func WithSudo(ctx context.Context, sudo bool) context.Context {
	return context.WithValue(ctx, sudoCtxKey{}, sudo)
}

// GetSudo gets sudo from context.
func GetSudo(ctx context.Context) bool {
	sudo, ok := ctx.Value(sudoCtxKey{}).(bool)
	if !ok {
		return false
	}
	return sudo
}
