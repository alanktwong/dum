// Package playbook contains the abstractions for installing and managing software
// configurations and installations via playbooks, plays and tasks.
package playbook

import (
	"fmt"

	pl "github.com/alanktwong/dum/internal/plays"
	t "github.com/alanktwong/dum/internal/types"

	om "github.com/elliotchance/orderedmap/v3"
)

const (
	// Ellipsis is prefix used in logs and println.
	Ellipsis = "..."
	// PlayEllipsis is prefix used in logs and println.
	PlayEllipsis = pl.PlayEllipsis
)

// PlayBook consists of a single install execution of many plays.
type PlayBook struct {
	t.Attributes
	Plays         []*pl.Play
	JetBrainsApps map[string]string
}

// NewPlayBook constructs a PlayBook.
func NewPlayBook(attributes *t.Attributes, plays []*pl.Play, apps map[string]string) (*PlayBook, error) {
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
func (p *PlayBook) GetAttributes() t.Attributes {
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

// GetJetBrainsApps implements PlayBookInfo.
func (p *PlayBook) GetJetBrainsApps() map[string]string {
	return p.JetBrainsApps
}

// GetPlays indexes all the Plays by its ID in an ordered map. If active is true it finds
// all the active plays. Otherwise, it finds them all.
func (p *PlayBook) GetPlays(active bool) *om.OrderedMap[string, *pl.Play] {
	plays := om.NewOrderedMap[string, *pl.Play]()
	if !active {
		for _, play := range p.Plays {
			plays.Set(play.GetID(), play)
		}
	} else {
		for _, play := range p.Plays {
			if play.IsEnabled() {
				plays.Set(play.GetID(), play)
			}
		}
	}
	return plays
}
