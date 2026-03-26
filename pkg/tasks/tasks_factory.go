// Package tasks provides Task types and related utilities.
package tasks

import (
	ext "awong/dotfiles/pkg/external"
	l "awong/dotfiles/pkg/logging"
	ty "awong/dotfiles/pkg/types"
	tyg "awong/dotfiles/pkg/types/gen"
	"fmt"
)

// TaskFactory creates Task instances from YML configuration.
type TaskFactory struct {
	Log   l.Logger
	Utils ext.Ext
}

// NewTaskFactory returns a new TaskFactory for creating Task instances.
func NewTaskFactory() *TaskFactory {
	return &TaskFactory{
		Log:   l.Log(),
		Utils: ext.NewExt(),
	}
}

// ProvideTasks creates multiple Task instances from YML configuration.
func (f *TaskFactory) ProvideTasks(yml map[string]any) ([]Task, error) {
	var tasks []Task
	if arr, ok := yml["tasks"].([]any); ok {
		for _, it := range arr {
			if m, ok := it.(map[string]any); ok {
				task, err := f.provideTask(m)
				if err != nil {
					return nil, fmt.Errorf("failed to provide task: %w", err)
				}
				tasks = append(tasks, task)
			} else {
				return nil, fmt.Errorf("task is not a map[string]interface{}")
			}
		}
	}
	return tasks, nil
}

func (f *TaskFactory) provideTask(yml map[string]any) (Task, error) {
	id := f.Utils.GetString(yml, "id", "")
	description := f.Utils.GetString(yml, "description", id)
	attributes, err := ty.NewAttributes(
		id,
		description,
		f.Utils.GetBool(yml, "enabled", true),
		f.Utils.GetBool(yml, "sudo", false),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create attributes: %w", err)
	}

	typeStr := f.Utils.GetString(yml, "type", "")
	taskType, err := tyg.ParseTaskType(typeStr)
	if err != nil {
		return nil, fmt.Errorf("unknown task type %s for task %s", typeStr, id)
	}
	switch taskType {
	case tyg.TaskTypeDir:
		return NewDirTask(attributes)
	case tyg.TaskTypeMas:
		return NewMasTask(attributes)
	case tyg.TaskTypeVscode:
		return NewVsCodePluginTask(attributes)
	case tyg.TaskTypeFunction:
		return NewFunctionTask(attributes)
	case tyg.TaskTypeLink:
		return NewLinkTask(attributes,
			f.Utils.GetString(yml, "root", ""),
			f.Utils.GetString(yml, "target", ""))
	case tyg.TaskTypeGit:
		return NewGitTask(attributes,
			f.Utils.GetString(yml, "root", ""),
			f.Utils.GetString(yml, "name", ""))
	case tyg.TaskTypeJetbrains:
		apps := f.Utils.GetStrings(yml, "apps", make([]string, 0))
		return NewJetBrainsPluginTask(attributes, apps)
	case tyg.TaskTypeBrew:
		return NewBrewTask(attributes,
			f.Utils.GetString(yml, "tap", ""))
	case tyg.TaskTypeCask:
		return NewBrewCaskTask(attributes,
			f.Utils.GetString(yml, "tap", ""))
	case tyg.TaskTypeCellar:
		return NewBrewCellarTask(attributes,
			f.Utils.GetString(yml, "tap", ""))
	case tyg.TaskTypeBash:
		return NewBashTask(attributes,
			f.Utils.GetString(yml, "command", ""),
			f.Utils.GetString(yml, "script", ""))
	default:
		return nil, fmt.Errorf("unknown task type %s for task %s", taskType, id)
	}
}
