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
		Enabled: true,
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
		Enabled: true,
		Sudo: true,
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
		Enabled: true,
		Sudo: true,
	}
	f := NewPlayFactory()

	tasks, err := f.TaskFactory.ProvideTasks(map[string]interface{}{
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
	play, err := NewPlay(attr, tasks)
	assert.NoError(t, err)
	assert.NotNil(t, play)
	assert.Equal(t, attr.ID, play.ID)
	assert.Equal(t,  tasks, play.Tasks)
}
