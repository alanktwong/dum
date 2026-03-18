package playbook

import (
	"context"
	"fmt"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func findGitTask(t *testing.T, ctx *Input) *GitTask {
	id := "https://github.com/tmux-plugins/tpm.git"
	tk, err := findTask(ctx.PlayBook, id)
	assert.NoError(t, err)
	task, ok := tk.(*GitTask)
	task.Name = "tpm"
	assert.True(t, ok)
	return task
}

func TestGitTask_Install(t *testing.T) {
	input := createTestInput(t)
	input.DryRun = false
	task := findGitTask(t, input)
	ctx := context.Background()

	mockGit := NewMockGit(t)
	path := filepath.Join(task.Root, task.Name)
	mockGit.EXPECT().AlreadyExists(path).Return(false).Once()
	mockGit.EXPECT().Clone(ctx, task.ID, task.Name, task.Root, task.Sudo).Return(nil).Once()
	task.Git = mockGit

	// when
	got, err := task.Install(ctx, input)
	// then
	assert.NoError(t, err)
	wantedResult := expectTaskResult(t, true, task.Attributes, input)
	assert.Equalf(t, wantedResult, got, "Install(%v) should git clone", input)
}

func TestGitTask_Install_Fail(t *testing.T) {
	input := createTestInput(t)
	input.DryRun = false
	task := findGitTask(t, input)
	ctx := context.Background()

	mockGit := NewMockGit(t)
	path := filepath.Join(task.Root, task.Name)
	mockGit.EXPECT().AlreadyExists(path).Return(false).Once()
	mockGit.EXPECT().Clone(ctx, task.ID, task.Name, task.Root, task.Sudo).Return(fmt.Errorf("git clone failed")).Once()
	task.Git = mockGit

	// when
	got, err := task.Install(ctx, input)
	// then
	assert.Error(t, err)
	assert.Nilf(t, got, "Install(%v) can fail git clone", input)
}

func TestGitTask_Disabled(t *testing.T) {
	input := createTestInput(t)
	input.DryRun = false
	task := findGitTask(t, input)
	task.Enabled = false
	ctx := context.Background()

	mockGit := NewMockGit(t)
	task.Git = mockGit

	// when
	got, err := task.Install(ctx, input)
	// then
	assert.NoError(t, err)
	wantedResult := expectTaskResult(t, false, task.Attributes, input)
	assert.Equalf(t, wantedResult, got, "Install(%v) should git clone", input)
	mockGit.AssertNotCalled(t, "Clone")
}
