package playbook

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetAllPlays(t *testing.T) {
	pb := createTestPlayBook(t)
	playMap := pb.GetPlays(false)
	assert.Equal(t, 2, playMap.Len(), "should return 2 plays")
}

func TestGetActivePlays(t *testing.T) {
	pb := createTestPlayBook(t)
	allPlayMap := pb.GetPlays(false)
	play, ok := allPlayMap.Get("play-1")
	assert.True(t, ok, "play should exist in playbook")
	play.Enabled = true

	activePlayMap := pb.GetPlays(true)
	assert.Equal(t, 1, activePlayMap.Len(), "should return 1 active plays")
}

func expectTaskResult(t *testing.T, success bool, attr Attributes, ctx *Input) *TaskResult {
	res, err := attr.CreateTaskResult(ctx, success)
	assert.NoError(t, err)
	return res
}

func createTestPlayBook(t *testing.T) *PlayBook {
	input := createTestInput(t)
	return input.PlayBook
}

func createTestInput(t *testing.T) *Input {
	f := NewFactory()
	input, err := f.Provide(InputOptions{
		File:   "../../pkg/testdata/test_installer.yml",
		Group:  "my-test-play",
		DryRun: true,
	})
	assert.NoError(t, err)
	return input
}

func findTask(pb *PlayBook, id string) (Task, error) {
	playMap := pb.GetPlays(false)
	for _, play := range playMap.AllFromFront() {
		taskMap := play.GetTasks(false)
		if task, ok := taskMap.Get(id); ok {
			return task, nil
		}
	}
	return nil, fmt.Errorf("cannot find task %s in playbook %s", id, pb.ID)
}

func createTestBrewTask(t *testing.T, attr Attributes) *BrewTask {
	brewTask, err := NewBrewTask(&attr, "")
	assert.NoError(t, err)
	return brewTask
}

func createTestBrewCaskTask(t *testing.T, attr Attributes) *BrewCaskTask {
	brewTask, err := NewBrewCaskTask(&attr, "")
	assert.NoError(t, err)
	return brewTask
}

func createTestBrewCellarTask(t *testing.T, attr Attributes) *BrewCellarTask {
	brewTask, err := NewBrewCellarTask(&attr, "")
	assert.NoError(t, err)
	return brewTask
}

func createTestDirTask(t *testing.T, attr Attributes) *DirTask {
	task, err := NewDirTask(&attr)
	assert.NoError(t, err)
	return task
}

func createTestFunctionTask(t *testing.T, attr Attributes) *FunctionTask {
	task, err := NewFunctionTask(&attr)
	assert.NoError(t, err)
	return task
}

func createTestGitTask(t *testing.T, attr Attributes, root, name string) *GitTask {
	task, err := NewGitTask(&attr, root, name)
	assert.NoError(t, err)
	return task
}

func createTestJetBrainsTask(t *testing.T, attr Attributes, apps []string) *JetBrainsPluginTask {
	task, err := NewJetBrainsPluginTask(&attr, apps)
	assert.NoError(t, err)
	return task
}

func createTestLinkTask(t *testing.T, attr Attributes, root, target string) *LinkTask {
	task, err := NewLinkTask(&attr, root, target)
	assert.NoError(t, err)
	return task
}

func createTestMasTask(t *testing.T, attr Attributes) *MasTask {
	task, err := NewMasTask(&attr)
	assert.NoError(t, err)
	return task
}

func createTestVsCodeTask(t *testing.T, attr Attributes) *VsCodePluginTask {
	task, err := NewVsCodePluginTask(&attr)
	assert.NoError(t, err)
	return task
}

func createTestBrewTap(t *testing.T, attr Attributes) *BrewTap {
	installer, err := NewBrewTap(&attr)
	assert.NoError(t, err)
	return installer
}
