package playbook

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPlay_GetAllTasks(t *testing.T) {
	pb := createTestPlayBook(t)
	allPlayMap := pb.GetPlays(false)
	play2, ok := allPlayMap.Get("play-2")
	assert.True(t, ok, "play-2 should exist in playbook")

	allTasks := play2.GetTasks(false)
	assert.Equal(t, 3, allTasks.Len(), "should return 3 tasks")
}

func TestPlay_GetActiveTasks(t *testing.T) {
	pb := createTestPlayBook(t)
	allPlayMap := pb.GetPlays(false)
	play1, ok := allPlayMap.Get("play-1")
	assert.True(t, ok, "play-1 should exist in playbook")

	assert.True(t, play1.Enabled, "assume play is enabled")
	play1.Enabled = false
	activeTasks := play1.GetTasks(true)
	assert.Equal(t, 0, activeTasks.Len(), "should no active tasks when play is disabled")

	play1.Enabled = true
	againActiveTasks := play1.GetTasks(true)
	assert.Equal(t, 10, againActiveTasks.Len(), "should return 10 active plays when play is enabled")
}
