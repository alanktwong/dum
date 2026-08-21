// Package tasks provides Task types and related utilities.
package tasks

import (
	"fmt"

	ext "github.com/alanktwong/dum/internal/external"
	l "github.com/alanktwong/dum/internal/logging"
	ty "github.com/alanktwong/dum/internal/types"
	tyg "github.com/alanktwong/dum/internal/types/gen"
	yml "github.com/alanktwong/dum/internal/yaml"
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

// ProvideTasks creates multiple Task instances from typed YML configuration.
func (f *TaskFactory) ProvideTasks(tasksYAML []yml.TaskYAML) ([]Task, error) {
	var tasks []Task
	for _, taskYAML := range tasksYAML {
		task, err := f.provideTask(taskYAML)
		if err != nil {
			return nil, fmt.Errorf("failed to provide task: %w", err)
		}
		tasks = append(tasks, task)
	}
	return tasks, nil
}

func (f *TaskFactory) provideTask(taskYAML yml.TaskYAML) (Task, error) {
	attributes, err := ty.NewAttributes(
		taskYAML.ID,
		taskYAML.Description,
		taskYAML.Enabled,
		taskYAML.Root != "",
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create attributes: %w", err)
	}

	taskType, err := tyg.ParseTaskType(taskYAML.Type)
	if err != nil {
		return nil, fmt.Errorf("unknown task type %s for task %s", taskYAML.Type, taskYAML.ID)
	}
	switch taskType {
	case tyg.TaskTypeDir:
		return NewDirTask(attributes)
	case tyg.TaskTypeMas:
		return NewMasTask(attributes)
	case tyg.TaskTypeVscode:
		return NewVsCodePluginTask(attributes)
	case tyg.TaskTypeLink:
		return NewLinkTask(attributes,
			taskYAML.Root,
			taskYAML.Target)
	case tyg.TaskTypeGit:
		return NewGitTask(attributes,
			taskYAML.Root,
			taskYAML.Name)
	case tyg.TaskTypeJetbrains:
		return NewJetBrainsPluginTask(attributes, taskYAML.Apps)
	case tyg.TaskTypeBrew:
		return NewBrewTask(attributes,
			taskYAML.Tap)
	case tyg.TaskTypeCask:
		return NewBrewCaskTask(attributes,
			taskYAML.Tap)
	case tyg.TaskTypeCellar:
		return NewBrewCellarTask(attributes,
			taskYAML.Tap)
	case tyg.TaskTypeBash:
		return NewBashTask(attributes,
			taskYAML.Command,
			taskYAML.Script)
	default:
		return nil, fmt.Errorf("unknown task type %s for task %s", taskType, taskYAML.ID)
	}
}
