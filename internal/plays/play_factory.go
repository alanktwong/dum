package plays

import (
	ext "awong/dotfiles/internal/external"
	l "awong/dotfiles/internal/logging"
	t "awong/dotfiles/internal/tasks"
	"awong/dotfiles/internal/types"
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

// ProvidePlays creates multiple Play instances from YML configuration.
func (f *PlayFactory) ProvidePlays(yml map[string]any) ([]*Play, error) {
	var plays []*Play
	if arr, ok := yml["plays"].([]any); ok {
		for _, each := range arr {
			if m, ok := each.(map[string]any); ok {
				play, err := f.ProvidePlay(m)
				if err != nil {
					return nil, fmt.Errorf("failed to provide play: %w", err)
				}
				plays = append(plays, play)
			} else {
				return nil, fmt.Errorf("play is not a map[string]interface{}")
			}
		}
	}
	return plays, nil
}

// ProvidePlay creates a single Play instance from YML configuration.
func (f *PlayFactory) ProvidePlay(yml map[string]any) (*Play, error) {
	tasks, err := f.TaskFactory.ProvideTasks(yml)
	if err != nil {
		return nil, fmt.Errorf("failed to provide tasks: %w", err)
	}
	id := f.getString(yml, "id", "")
	description := f.getString(yml, "description", id)
	if len(tasks) == 0 {
		return nil, fmt.Errorf("play %v has no tasks", id)
	}
	attributes, err := types.NewAttributes(
		id,
		description,
		f.getBool(yml, "enabled", true),
		f.getBool(yml, "sudo", false),
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

func (f *PlayFactory) getString(data map[string]any, key string, def string) string {
	return f.Utils.GetString(data, key, def)
}

func (f *PlayFactory) getBool(data map[string]any, key string, def bool) bool {
	return f.Utils.GetBool(data, key, def)
}
