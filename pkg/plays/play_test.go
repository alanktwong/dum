package plays

import (
	tk "awong/dotfiles/pkg/tasks"
	ty "awong/dotfiles/pkg/types"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewPlay_NilAttributes(t *testing.T) {
	tasks := *new([]tk.Task)
	play, err := NewPlay(nil, tasks)
	assert.Error(t, err)
	assert.Nil(t, play)
	assert.Contains(t, err.Error(), "attributes cannot be nil")
}

func TestNewPlay_EmptyID(t *testing.T) {
	attr := &ty.Attributes{
		ID:          "",
		Description: "test",
	}
	tasks := *new([]tk.Task)
	play, err := NewPlay(attr, tasks)
	assert.Error(t, err)
	assert.Nil(t, play)
	assert.Contains(t, err.Error(), "play ID cannot be empty")
}

func TestNewPlay_Success(t *testing.T) {
	attr := &ty.Attributes{
		ID:          "test-play",
		Description: "",
	}
	tasks := *new([]tk.Task)
	play, err := NewPlay(attr, tasks)
	assert.NoError(t, err)
	assert.NotNil(t, play)
	assert.Equal(t, attr.ID, play.ID)
}

func TestNewPlay_Success_Enabled(t *testing.T) {
	attr := &ty.Attributes{
		ID:          "test-play",
		Description: "",
		Enabled:     true,
	}
	tasks := *new([]tk.Task)
	play, err := NewPlay(attr, tasks)
	assert.NoError(t, err)
	assert.NotNil(t, play)
	assert.Equal(t, attr.ID, play.ID)
	assert.Equal(t, attr.Enabled, play.IsEnabled())
}

func TestNewPlay_Success_Sudo(t *testing.T) {
	attr := &ty.Attributes{
		ID:          "test-play",
		Description: "",
		Enabled:     true,
		Sudo:        true,
	}
	tasks := *new([]tk.Task)
	play, err := NewPlay(attr, tasks)
	assert.NoError(t, err)
	assert.NotNil(t, play)
	assert.Equal(t, attr.ID, play.ID)
	assert.Equal(t, attr.Enabled, play.IsEnabled())
	assert.Equal(t, attr.Sudo, play.Sudo)
}

func TestNewPlay_Success_Tasks(t *testing.T) {
	attr := &ty.Attributes{
		ID:          "test-play",
		Description: "",
		Enabled:     true,
		Sudo:        true,
	}
	f := NewPlayFactory()

	tasks, err := f.TaskFactory.ProvideTasks(map[string]any{
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
	play, err := NewPlay(attr, tasks)
	assert.NoError(t, err)
	assert.NotNil(t, play)
	assert.Equal(t, attr.ID, play.ID)
	assert.Equal(t, tasks, play.Tasks)
}

func TestPlay_GetTasks_All(t *testing.T) {
	attr := &ty.Attributes{
		ID:      "test-play",
		Enabled: true,
	}
	enabledTask := &MockTask{
		Attr: ty.Attributes{ID: "task-1", Enabled: true},
	}
	disabledTask := &MockTask{
		Attr: ty.Attributes{ID: "task-2", Enabled: false},
	}
	play, err := NewPlay(attr, []tk.Task{enabledTask, disabledTask})
	assert.NoError(t, err)

	tasks := play.GetTasks(false)
	assert.Equal(t, 2, tasks.Len())
	assert.True(t, tasks.Has("task-1"))
	assert.True(t, tasks.Has("task-2"))
}

func TestPlay_GetTasks_Filtered(t *testing.T) {
	attr := &ty.Attributes{
		ID:      "test-play",
		Enabled: true,
	}
	enabledTask := &MockTask{
		Attr: ty.Attributes{ID: "task-1", Enabled: true},
	}
	disabledTask := &MockTask{
		Attr: ty.Attributes{ID: "task-2", Enabled: false},
	}
	play, err := NewPlay(attr, []tk.Task{enabledTask, disabledTask})
	assert.NoError(t, err)

	tasks := play.GetTasks(true)
	assert.Equal(t, 1, tasks.Len())
	assert.True(t, tasks.Has("task-1"))
	assert.False(t, tasks.Has("task-2"))
}

func TestPlay_GetTasks_Filtered_DisabledPlay(t *testing.T) {
	attr := &ty.Attributes{
		ID:      "test-play",
		Enabled: false,
	}
	enabledTask := &MockTask{
		Attr: ty.Attributes{ID: "task-1", Enabled: true},
	}
	play, err := NewPlay(attr, []tk.Task{enabledTask})
	assert.NoError(t, err)

	tasks := play.GetTasks(true)
	assert.Equal(t, 0, tasks.Len())
}
