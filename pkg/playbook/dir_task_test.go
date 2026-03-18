package playbook

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func findDirTask(t *testing.T, input *Input) *DirTask {
	id := "~/.dotfiles"
	tk, err := findTask(input.PlayBook, id)
	assert.NoError(t, err)
	task, ok := tk.(*DirTask)
	assert.True(t, ok)
	return task
}

func TestDirTask_Install(t *testing.T) {
	input := createTestInput(t)
	input.DryRun = false
	task := findDirTask(t, input)
	ctx := context.Background()

	mockUtils := NewMockExt(t)
	expanded := fmt.Sprintf("/Users/user/%v", task.ID)
	mockUtils.EXPECT().ExpandUser(task.ID).Return(expanded, nil).Once()
	mockUtils.EXPECT().IsDir(expanded).Return(false).Once()
	mockUtils.EXPECT().CreateDirectory(ctx, expanded, task.Sudo).Return(nil).Once()
	task.Utils = mockUtils
	// when
	got, err := task.Install(ctx, input)
	// then
	assert.NoError(t, err)
	wantedResult := expectTaskResult(t, true, task.Attributes, input)
	assert.Equalf(t, wantedResult, got, "Install(%v) should mkdir", input)
}

func TestDirTask_Fail(t *testing.T) {
	input := createTestInput(t)
	input.DryRun = false
	task := findDirTask(t, input)
	ctx := context.Background()

	mockUtils := NewMockExt(t)
	expanded := fmt.Sprintf("/Users/user/%v", task.ID)
	mockUtils.EXPECT().ExpandUser(task.ID).Return(expanded, nil).Once()
	mockUtils.EXPECT().IsDir(expanded).Return(false).Once()
	mockUtils.EXPECT().CreateDirectory(ctx, expanded, task.Sudo).Return(fmt.Errorf("fail")).Once()

	task.Utils = mockUtils
	// when
	got, err := task.Install(ctx, input)
	// then
	assert.Errorf(t, err, "Install(%v) can fail to mkdir", input)
	assert.Nil(t, got)
}

func TestDirTask_IsDir(t *testing.T) {
	input := createTestInput(t)
	input.DryRun = false
	task := findDirTask(t, input)
	ctx := context.Background()

	mockUtils := NewMockExt(t)
	expanded := fmt.Sprintf("/Users/user/%v", task.ID)
	mockUtils.EXPECT().ExpandUser(task.ID).Return(expanded, nil).Once()
	mockUtils.EXPECT().IsDir(expanded).Return(true).Once()
	task.Utils = mockUtils
	// when
	got, err := task.Install(ctx, input)
	// then
	assert.NoError(t, err)
	wantedResult := expectTaskResult(t, false, task.Attributes, input)
	assert.Equalf(t, wantedResult, got, "Install(%v) should mkdir", input)
	mockUtils.AssertNotCalled(t, "CreateDirectory")
}

func TestDirTask_Fail_Expand(t *testing.T) {
	input := createTestInput(t)
	input.DryRun = false
	task := findDirTask(t, input)
	ctx := context.Background()

	mockUtils := NewMockExt(t)
	mockUtils.EXPECT().ExpandUser(task.ID).Return("", fmt.Errorf("fail")).Once()

	task.Utils = mockUtils
	// when
	got, err := task.Install(ctx, input)
	// then
	assert.Errorf(t, err, "Install(%v) can fail to mkdir", input)
	assert.Nil(t, got)
	mockUtils.AssertNotCalled(t, "IsDir", 1)
	mockUtils.AssertNotCalled(t, "CreateDirectory")
}
