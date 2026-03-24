package tasks

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewTaskFactory(t *testing.T) {
	factory := NewTaskFactory()
	assert.NotNil(t, factory)
	assert.NotNil(t, factory.Log)
	assert.NotNil(t, factory.Utils)
}

func TestTaskFactory_ProvideTasks_EmptyInput(t *testing.T) {
	factory := NewTaskFactory()
	tasks, err := factory.ProvideTasks(map[string]interface{}{})
	assert.NoError(t, err)
	assert.Empty(t, tasks)
}

func TestTaskFactory_ProvideTasks_NoTasksKey(t *testing.T) {
	factory := NewTaskFactory()
	tasks, err := factory.ProvideTasks(map[string]interface{}{
		"other": "data",
	})
	assert.NoError(t, err)
	assert.Empty(t, tasks)
}

func TestTaskFactory_ProvideTasks_InvalidTaskFormat(t *testing.T) {
	factory := NewTaskFactory()
	_, err := factory.ProvideTasks(map[string]interface{}{
		"tasks": []interface{}{
			"not a map",
		},
	})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "not a map[string]interface{}")
}

func TestTaskFactory_ProvideTasks_UnknownTaskType(t *testing.T) {
	factory := NewTaskFactory()
	_, err := factory.ProvideTasks(map[string]interface{}{
		"tasks": []interface{}{
			map[string]interface{}{
				"id":   "test-task",
				"type": "unknown-type",
			},
		},
	})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "unknown task type")
}

func TestTaskFactory_ProvideTasks_DirTask(t *testing.T) {
	factory := NewTaskFactory()
	tasks, err := factory.ProvideTasks(map[string]interface{}{
		"tasks": []interface{}{
			map[string]interface{}{
				"id":          "test-dir",
				"type":        "dir",
				"description": "test directory",
			},
		},
	})
	assert.NoError(t, err)
	assert.Len(t, tasks, 1)
	dirTask, ok := tasks[0].(*DirTask)
	assert.True(t, ok)
	assert.Equal(t, "test-dir", dirTask.ID)
}

func TestTaskFactory_ProvideTasks_LinkTask(t *testing.T) {
	factory := NewTaskFactory()
	tasks, err := factory.ProvideTasks(map[string]interface{}{
		"tasks": []interface{}{
			map[string]interface{}{
				"id":     "test-link",
				"type":   "link",
				"root":   "/tmp/root",
				"target": "/tmp/target",
			},
		},
	})
	assert.NoError(t, err)
	assert.Len(t, tasks, 1)
	linkTask, ok := tasks[0].(*LinkTask)
	assert.True(t, ok)
	assert.Equal(t, "test-link", linkTask.ID)
}

func TestTaskFactory_ProvideTasks_BrewTask(t *testing.T) {
	factory := NewTaskFactory()
	tasks, err := factory.ProvideTasks(map[string]interface{}{
		"tasks": []interface{}{
			map[string]interface{}{
				"id":   "test-brew",
				"type": "brew",
				"tap":  "homebrew/cask",
			},
		},
	})
	assert.NoError(t, err)
	assert.Len(t, tasks, 1)
	brewTask, ok := tasks[0].(*BrewTask)
	assert.True(t, ok)
	assert.Equal(t, "test-brew", brewTask.ID)
}

func TestTaskFactory_ProvideTasks_BashTask(t *testing.T) {
	factory := NewTaskFactory()
	tasks, err := factory.ProvideTasks(map[string]interface{}{
		"tasks": []interface{}{
			map[string]interface{}{
				"id":      "test-bash",
				"type":    "bash",
				"command": "echo hello",
			},
		},
	})
	assert.NoError(t, err)
	assert.Len(t, tasks, 1)
	bashTask, ok := tasks[0].(*BashTask)
	assert.True(t, ok)
	assert.Equal(t, "test-bash", bashTask.ID)
}

func TestTaskFactory_ProvideTasks_DisabledTask(t *testing.T) {
	factory := NewTaskFactory()
	tasks, err := factory.ProvideTasks(map[string]interface{}{
		"tasks": []interface{}{
			map[string]interface{}{
				"id":      "test-disabled",
				"type":    "dir",
				"enabled": false,
			},
		},
	})
	assert.NoError(t, err)
	assert.Len(t, tasks, 1)
	dirTask, ok := tasks[0].(*DirTask)
	assert.True(t, ok)
	assert.False(t, dirTask.Enabled)
}

func TestTaskFactory_ProvideTasks_MultipleTasks(t *testing.T) {
	factory := NewTaskFactory()
	tasks, err := factory.ProvideTasks(map[string]interface{}{
		"tasks": []interface{}{
			map[string]interface{}{
				"id":   "task-1",
				"type": "dir",
			},
			map[string]interface{}{
				"id":   "task-2",
				"type": "link",
				"root": "/tmp",
			},
			map[string]interface{}{
				"id":      "task-3",
				"type":    "bash",
				"command": "echo hello",
			},
		},
	})
	assert.NoError(t, err)
	assert.Len(t, tasks, 3)
	dirTask, ok := tasks[0].(*DirTask)
	assert.True(t, ok)
	assert.Equal(t, "task-1", dirTask.ID)
	linkTask, ok := tasks[1].(*LinkTask)
	assert.True(t, ok)
	assert.Equal(t, "task-2", linkTask.ID)
	bashTask, ok := tasks[2].(*BashTask)
	assert.True(t, ok)
	assert.Equal(t, "task-3", bashTask.ID)
}
