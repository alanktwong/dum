package playbook

import (
	"fmt"

	om "github.com/elliotchance/orderedmap/v3"
)

// Play consists of a single install execution of many tasks.
type Play struct {
	Attributes
	Tasks []Task
}

// NewPlay constructs a Play.
func NewPlay(attributes *Attributes, tasks []Task) (*Play, error) {
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
func (p *Play) GetAttributes() Attributes {
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
func (p *Play) GetTasks(active bool) *om.OrderedMap[string, Task] {
	tasks := om.NewOrderedMap[string, Task]()
	for _, task := range p.Tasks {
		if attributable, ok := task.(Attributable); ok {
			attr := attributable.GetAttributes()
			id := attr.GetID()
			if !active {
				tasks.Set(id, task)
			} else if p.Enabled && attr.IsEnabled() {
				tasks.Set(id, task)
			}
		}
	}
	return tasks
}
