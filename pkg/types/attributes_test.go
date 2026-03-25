package types

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewAttributes_Success(t *testing.T) {
	attrs, err := NewAttributes("test-id", "test description", true, false)
	assert.NoError(t, err)
	assert.NotNil(t, attrs)
	assert.Equal(t, "test-id", attrs.ID)
	assert.Equal(t, "test description", attrs.Description)
	assert.True(t, attrs.Enabled)
	assert.False(t, attrs.Sudo)
}

func TestNewAttributes_SudoEnabled(t *testing.T) {
	attrs, err := NewAttributes("sudo-task", "needs sudo", false, true)
	assert.NoError(t, err)
	assert.NotNil(t, attrs)
	assert.False(t, attrs.Enabled)
	assert.True(t, attrs.Sudo)
}

func TestNewAttributes_EmptyID(t *testing.T) {
	attrs, err := NewAttributes("", "description", true, false)
	assert.Error(t, err)
	assert.Nil(t, attrs)
	assert.Contains(t, err.Error(), "attribute ID cannot be empty")
}

func TestAttributes_IsEnabled(t *testing.T) {
	tests := []struct {
		name    string
		enabled bool
	}{
		{"enabled", true},
		{"disabled", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			attrs := &Attributes{ID: "test", Enabled: tt.enabled}
			assert.Equal(t, tt.enabled, attrs.IsEnabled())
		})
	}
}

func TestAttributes_GetID(t *testing.T) {
	attrs := &Attributes{ID: "my-id"}
	assert.Equal(t, "my-id", attrs.GetID())
}

func TestAttributes_CreateTaskResult_Success(t *testing.T) {
	attrs := &Attributes{ID: "task-1"}
	input := &TaskInput{
		Play:     "my-play",
		PlayBook: "my-book",
		DryRun:   false,
	}

	result, err := attrs.CreateTaskResult(input, true)
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "task-1", result.Task)
	assert.Equal(t, "my-play", result.Play)
	assert.Equal(t, "my-book", result.PlayBook)
	assert.False(t, result.DryRun)
	assert.True(t, result.Success)
}

func TestAttributes_CreateTaskResult_Failure(t *testing.T) {
	attrs := &Attributes{ID: "task-2"}
	input := &TaskInput{
		Play:     "play-a",
		PlayBook: "book-b",
		DryRun:   true,
	}

	result, err := attrs.CreateTaskResult(input, false)
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "task-2", result.Task)
	assert.Equal(t, "play-a", result.Play)
	assert.Equal(t, "book-b", result.PlayBook)
	assert.True(t, result.DryRun)
	assert.False(t, result.Success)
}
