package plays

import (
	t "awong/dotfiles/pkg/tasks"
	ty "awong/dotfiles/pkg/types"
	"fmt"

	om "github.com/elliotchance/orderedmap/v3"
)

const (
	// PlayEllipsis is prefix used in logs and println.
	PlayEllipsis = "......."
)

// Play consists of a single install execution of many tasks.
type Play struct {
	ty.Attributes
	Tasks []t.Task
}

// NewPlay constructs a Play.
func NewPlay(attributes *ty.Attributes, tasks []t.Task) (*Play, error) {
	if attributes == nil {
		return nil, fmt.Errorf("attributes cannot be nil")
	}
	if attributes.ID == "" {
		return nil, fmt.Errorf("play ID cannot be empty")
	}
	return &Play{
		Attributes: *attributes,
		Tasks:      tasks,
	}, nil
}

// GetAttributes implements Attributable.
func (p *Play) GetAttributes() ty.Attributes {
	return p.Attributes
}

// GetID implements Identifiable.
func (p *Play) GetID() string {
	return p.ID
}

// IsEnabled implements Enableable.
func (p *Play) IsEnabled() bool {
	return p.Enabled
}

// GetTasks indexes all the Tasks by its ID in an ordered map. If active is true it finds
// all the active plays. Otherwise, it finds them all.
func (p *Play) GetTasks(active bool) *om.OrderedMap[string, t.Task] {
	tasks := om.NewOrderedMap[string, t.Task]()
	for _, it := range p.Tasks {
		if attributable, ok := it.(ty.Attributable); ok {
			attr := attributable.GetAttributes()
			id := attr.GetID()
			if !active {
				tasks.Set(id, it)
			} else if p.Enabled && attr.IsEnabled() {
				tasks.Set(id, it)
			}
		}
	}
	return tasks
}
