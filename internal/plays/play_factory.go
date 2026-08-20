package plays

import (
	ext "alanktwong/dum/internal/external"
	l "alanktwong/dum/internal/logging"
	t "alanktwong/dum/internal/tasks"
	"alanktwong/dum/internal/types"
	yml "alanktwong/dum/internal/yaml"
	"fmt"
)

// PlayFactory creates Play instances from YML configuration.
type PlayFactory struct {
	Log         l.Logger
	Utils       ext.Ext
	TaskFactory *t.TaskFactory
}

// NewPlayFactory returns a new PlayFactory for creating Play instances.
func NewPlayFactory() *PlayFactory {
	return &PlayFactory{
		Log:         l.Log(),
		Utils:       ext.NewExt(),
		TaskFactory: t.NewTaskFactory(),
	}
}

// ProvidePlays creates multiple Play instances from typed YML configuration.
func (f *PlayFactory) ProvidePlays(playsYAML []yml.PlayYAML) ([]*Play, error) {
	var plays []*Play
	for _, playYAML := range playsYAML {
		play, err := f.providePlay(playYAML)
		if err != nil {
			return nil, fmt.Errorf("failed to provide play: %w", err)
		}
		plays = append(plays, play)
	}
	return plays, nil
}

// ProvidePlay creates a single Play instance from typed YML configuration.
func (f *PlayFactory) providePlay(playYAML yml.PlayYAML) (*Play, error) {
	tasks, err := f.TaskFactory.ProvideTasks(playYAML.Tasks)
	if err != nil {
		return nil, fmt.Errorf("failed to provide tasks: %w", err)
	}
	if len(tasks) == 0 {
		return nil, fmt.Errorf("play %v has no tasks", playYAML.ID)
	}
	attributes, err := types.NewAttributes(
		playYAML.ID,
		playYAML.Description,
		playYAML.Enabled,
		false,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create attributes: %w", err)
	}
	play, err := NewPlay(attributes, tasks)
	if err != nil {
		return nil, fmt.Errorf("failed to create play: %w", err)
	}
	return play, nil
}
