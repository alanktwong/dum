package tasks

import (
	"testing"

	yml "github.com/alanktwong/dum/internal/yaml"

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
	tasks, err := factory.ProvideTasks([]yml.TaskYAML{})
	assert.NoError(t, err)
	assert.Empty(t, tasks)
}

func TestTaskFactory_ProvideTasks_DirTask(t *testing.T) {
	factory := NewTaskFactory()
	tasks, err := factory.ProvideTasks([]yml.TaskYAML{
		{
			ID:          "test-dir",
			Type:        "dir",
			Description: "test directory",
			Enabled:     true,
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
	tasks, err := factory.ProvideTasks([]yml.TaskYAML{
		{
			ID:     "test-link",
			Type:   "link",
			Root:   "/tmp/root",
			Target: "/tmp/target",
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
	tasks, err := factory.ProvideTasks([]yml.TaskYAML{
		{
			ID:   "test-brew",
			Type: "brew",
			Tap:  "homebrew/cask",
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
	tasks, err := factory.ProvideTasks([]yml.TaskYAML{
		{
			ID:      "test-bash",
			Type:    "bash",
			Command: "echo hello",
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
	tasks, err := factory.ProvideTasks([]yml.TaskYAML{
		{
			ID:      "test-disabled",
			Type:    "dir",
			Enabled: false,
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
	tasks, err := factory.ProvideTasks([]yml.TaskYAML{
		{
			ID:   "task-1",
			Type: "dir",
		},
		{
			ID:   "task-2",
			Type: "link",
			Root: "/tmp",
		},
		{
			ID:      "task-3",
			Type:    "bash",
			Command: "echo hello",
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

func TestTaskFactory_ProvideTasks_UnknownTaskType(t *testing.T) {
	factory := NewTaskFactory()
	_, err := factory.ProvideTasks([]yml.TaskYAML{
		{
			ID:   "test-task",
			Type: "unknown-type",
		},
	})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "unknown task type")
}

func TestTaskFactory_ProvideTasks_MasTask(t *testing.T) {
	factory := NewTaskFactory()
	tasks, err := factory.ProvideTasks([]yml.TaskYAML{
		{
			ID:   "test-mas",
			Type: "mas",
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
	tasks, err := factory.ProvideTasks([]yml.TaskYAML{
		{
			ID:   "test-vscode",
			Type: "vscode",
		},
	})
	assert.NoError(t, err)
	assert.Len(t, tasks, 1)
	vsTask, ok := tasks[0].(*VsCodePluginTask)
	assert.True(t, ok)
	assert.Equal(t, "test-vscode", vsTask.ID)
}

func TestTaskFactory_ProvideTasks_UnknownType(t *testing.T) {
	factory := NewTaskFactory()
	tasks, err := factory.ProvideTasks([]yml.TaskYAML{
		{
			ID:   "test-function",
			Type: "function",
		},
	})
	assert.Nil(t, tasks)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "unknown task type function")
}

func TestTaskFactory_ProvideTasks_GitTask(t *testing.T) {
	factory := NewTaskFactory()
	tasks, err := factory.ProvideTasks([]yml.TaskYAML{
		{
			ID:   "test-git",
			Type: "git",
			Name: "dotfiles",
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
	tasks, err := factory.ProvideTasks([]yml.TaskYAML{
		{
			ID:   "test-jetbrains",
			Type: "jetbrains",
			Apps: []string{"IntelliJ IDEA", "GoLand"},
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
	tasks, err := factory.ProvideTasks([]yml.TaskYAML{
		{
			ID:   "test-cask",
			Type: "cask",
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
	tasks, err := factory.ProvideTasks([]yml.TaskYAML{
		{
			ID:   "test-cellar",
			Type: "cellar",
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
	tasks, err := factory.ProvideTasks([]yml.TaskYAML{
		{
			ID:   "test-sudo",
			Type: "dir",
			Root: "/some/path",
		},
	})
	assert.NoError(t, err)
	assert.Len(t, tasks, 1)
	dirTask, ok := tasks[0].(*DirTask)
	assert.True(t, ok)
	assert.True(t, dirTask.Sudo)
}
