package types

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewTaskInput_AllFields(t *testing.T) {
	apps := map[string]string{"idea": "/path/to/idea", "goland": "/path/to/goland"}
	input := NewTaskInput(true, true, "my-task", "my-play", "my-book", apps)

	assert.NotNil(t, input)
	assert.True(t, input.DryRun)
	assert.True(t, input.Sudo)
	assert.Equal(t, "my-task", input.Task)
	assert.Equal(t, "my-play", input.Play)
	assert.Equal(t, "my-book", input.PlayBook)
	assert.Equal(t, apps, input.JetBrainsApps)
}

func TestNewTaskInput_NilJetBrainsApps(t *testing.T) {
	input := NewTaskInput(false, false, "task", "play", "book", nil)

	assert.NotNil(t, input)
	assert.False(t, input.DryRun)
	assert.False(t, input.Sudo)
	assert.Nil(t, input.JetBrainsApps)
}

func TestNewTaskInput_EmptyStrings(t *testing.T) {
	input := NewTaskInput(false, false, "", "", "", nil)

	assert.NotNil(t, input)
	assert.Empty(t, input.Task)
	assert.Empty(t, input.Play)
	assert.Empty(t, input.PlayBook)
}
