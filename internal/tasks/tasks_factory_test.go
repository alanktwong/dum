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
	tasks, err := factory.ProvideTasks(map[string]any{})
	assert.NoError(t, err)
	assert.Empty(t, tasks)
}

func TestTaskFactory_ProvideTasks_NoTasksKey(t *testing.T) {
	factory := NewTaskFactory()
	tasks, err := factory.ProvideTasks(map[string]any{
		"other": "data",
	})
	assert.NoError(t, err)
	assert.Empty(t, tasks)
}

func TestTaskFactory_ProvideTasks_InvalidTaskFormat(t *testing.T) {
	factory := NewTaskFactory()
	_, err := factory.ProvideTasks(map[string]any{
		"tasks": []any{
			"not a map",
		},
	})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "not a map[string]interface{}")
}

func TestTaskFactory_ProvideTasks_UnknownTaskType(t *testing.T) {
	factory := NewTaskFactory()
	_, err := factory.ProvideTasks(map[string]any{
		"tasks": []any{
			map[string]any{
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
	tasks, err := factory.ProvideTasks(map[string]any{
		"tasks": []any{
			map[string]any{
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
	tasks, err := factory.ProvideTasks(map[string]any{
		"tasks": []any{
			map[string]any{
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
	tasks, err := factory.ProvideTasks(map[string]any{
		"tasks": []any{
			map[string]any{
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
	tasks, err := factory.ProvideTasks(map[string]any{
		"tasks": []any{
			map[string]any{
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
	tasks, err := factory.ProvideTasks(map[string]any{
		"tasks": []any{
			map[string]any{
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
	tasks, err := factory.ProvideTasks(map[string]any{
		"tasks": []any{
			map[string]any{
				"id":   "task-1",
				"type": "dir",
			},
			map[string]any{
				"id":   "task-2",
				"type": "link",
				"root": "/tmp",
			},
			map[string]any{
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

func TestTaskFactory_ProvideTasks_MasTask(t *testing.T) {
	factory := NewTaskFactory()
	tasks, err := factory.ProvideTasks(map[string]any{
		"tasks": []any{
			map[string]any{
				"id":   "test-mas",
				"type": "mas",
			},
		},
	})
	assert.NoError(t, err)
	assert.Len(t, tasks, 1)
	masTask, ok := tasks[0].(*MasTask)
	assert.True(t, ok)
	assert.Equal(t, "test-mas", masTask.ID)
}

func TestTaskFactory_ProvideTasks_VsCodeTask(t *testing.T) {
	factory := NewTaskFactory()
	tasks, err := factory.ProvideTasks(map[string]any{
		"tasks": []any{
			map[string]any{
				"id":   "test-vscode",
				"type": "vscode",
			},
		},
	})
	assert.NoError(t, err)
	assert.Len(t, tasks, 1)
	vsTask, ok := tasks[0].(*VsCodePluginTask)
	assert.True(t, ok)
	assert.Equal(t, "test-vscode", vsTask.ID)
}

func TestTaskFactory_ProvideTasks_FunctionTask(t *testing.T) {
	factory := NewTaskFactory()
	tasks, err := factory.ProvideTasks(map[string]any{
		"tasks": []any{
			map[string]any{
				"id":   "test-function",
				"type": "function",
			},
		},
	})
	assert.NoError(t, err)
	assert.Len(t, tasks, 1)
	fnTask, ok := tasks[0].(*FunctionTask)
	assert.True(t, ok)
	assert.Equal(t, "test-function", fnTask.ID)
}

func TestTaskFactory_ProvideTasks_GitTask(t *testing.T) {
	factory := NewTaskFactory()
	tasks, err := factory.ProvideTasks(map[string]any{
		"tasks": []any{
			map[string]any{
				"id":   "test-git",
				"type": "git",
				"name": "dotfiles",
			},
		},
	})
	assert.NoError(t, err)
	assert.Len(t, tasks, 1)
	gitTask, ok := tasks[0].(*GitTask)
	assert.True(t, ok)
	assert.Equal(t, "test-git", gitTask.ID)
}

func TestTaskFactory_ProvideTasks_JetBrainsTask(t *testing.T) {
	factory := NewTaskFactory()
	tasks, err := factory.ProvideTasks(map[string]any{
		"tasks": []any{
			map[string]any{
				"id":   "test-jetbrains",
				"type": "jetbrains",
				"apps": []any{"IntelliJ IDEA", "GoLand"},
			},
		},
	})
	assert.NoError(t, err)
	assert.Len(t, tasks, 1)
	jbTask, ok := tasks[0].(*JetBrainsPluginTask)
	assert.True(t, ok)
	assert.Equal(t, "test-jetbrains", jbTask.ID)
}

func TestTaskFactory_ProvideTasks_CaskTask(t *testing.T) {
	factory := NewTaskFactory()
	tasks, err := factory.ProvideTasks(map[string]any{
		"tasks": []any{
			map[string]any{
				"id":   "test-cask",
				"type": "cask",
			},
		},
	})
	assert.NoError(t, err)
	assert.Len(t, tasks, 1)
	caskTask, ok := tasks[0].(*BrewCaskTask)
	assert.True(t, ok)
	assert.Equal(t, "test-cask", caskTask.ID)
}

func TestTaskFactory_ProvideTasks_CellarTask(t *testing.T) {
	factory := NewTaskFactory()
	tasks, err := factory.ProvideTasks(map[string]any{
		"tasks": []any{
			map[string]any{
				"id":   "test-cellar",
				"type": "cellar",
			},
		},
	})
	assert.NoError(t, err)
	assert.Len(t, tasks, 1)
	cellarTask, ok := tasks[0].(*BrewCellarTask)
	assert.True(t, ok)
	assert.Equal(t, "test-cellar", cellarTask.ID)
}

func TestTaskFactory_ProvideTasks_SudoTask(t *testing.T) {
	factory := NewTaskFactory()
	tasks, err := factory.ProvideTasks(map[string]any{
		"tasks": []any{
			map[string]any{
				"id":   "test-sudo",
				"type": "dir",
				"sudo": true,
			},
		},
	})
	assert.NoError(t, err)
	assert.Len(t, tasks, 1)
	dirTask, ok := tasks[0].(*DirTask)
	assert.True(t, ok)
	assert.True(t, dirTask.Sudo)
}
