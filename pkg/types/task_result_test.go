package types

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewTaskResult_Success(t *testing.T) {
	input := &TaskInput{
		Task:     "task-1",
		Play:     "play-1",
		PlayBook: "book-1",
		DryRun:   false,
	}

	result, err := NewTaskResult(input, true)
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "task-1", result.Task)
	assert.Equal(t, "play-1", result.Play)
	assert.Equal(t, "book-1", result.PlayBook)
	assert.False(t, result.DryRun)
	assert.True(t, result.Success)
}

func TestNewTaskResult_Failure(t *testing.T) {
	input := &TaskInput{
		Task:     "task-2",
		Play:     "play-2",
		PlayBook: "book-2",
		DryRun:   false,
	}

	result, err := NewTaskResult(input, false)
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.False(t, result.Success)
}

func TestNewTaskResult_DryRun(t *testing.T) {
	input := &TaskInput{
		Task:     "dry-task",
		Play:     "dry-play",
		PlayBook: "dry-book",
		DryRun:   true,
	}

	result, err := NewTaskResult(input, true)
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.True(t, result.DryRun)
}
