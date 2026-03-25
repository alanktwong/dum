package playbook

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func findBashTask(t *testing.T, input *Input, id string) *BashTask {
	tk, err := findTask(input.PlayBook, id)
	assert.NoError(t, err)
	task, ok := tk.(*BashTask)
	assert.True(t, ok)
	return task
}

func TestBashTask_Install(t *testing.T) {
	input := createTestInput(t)
	input.DryRun = false
	task := findBashTask(t, input, "test-echo")
	ctx := context.Background()

	mockUtils := NewMockExt(t)
	mockUtils.EXPECT().RunCommand(ctx, task.Command, task.Sudo).Return(nil).Once()
	task.Utils = mockUtils
	// when
	got, err := task.Install(ctx, input)
	// then
	assert.NoError(t, err)
	wantedResult := expectTaskResult(t, true, task.Attributes, input)
	assert.Equalf(t, wantedResult, got, "Install(%v) should run command", input)
}

func TestBashTask_Install_WithScript(t *testing.T) {
	input := createTestInput(t)
	input.DryRun = false
	task := findBashTask(t, input, "test-script")
	ctx := context.Background()

	mockUtils := NewMockExt(t)
	mockUtils.EXPECT().RunCommand(ctx, task.Script, task.Sudo).Return(nil).Once()
	task.Utils = mockUtils
	// when
	got, err := task.Install(ctx, input)
	// then
	assert.NoError(t, err)
	wantedResult := expectTaskResult(t, true, task.Attributes, input)
	assert.Equalf(t, wantedResult, got, "Install(%v) should run script", input)
}

func TestBashTask_Fail(t *testing.T) {
	input := createTestInput(t)
	input.DryRun = false
	task := findBashTask(t, input, "test-echo")
	ctx := context.Background()

	mockUtils := NewMockExt(t)
	mockUtils.EXPECT().RunCommand(ctx, task.Command, task.Sudo).Return(fmt.Errorf("command failed")).Once()
	task.Utils = mockUtils
	// when
	got, err := task.Install(ctx, input)
	// then
	assert.Errorf(t, err, "Install(%v) can fail to run command", input)
	assert.Nil(t, got)
}

func TestBashTask_Disabled(t *testing.T) {
	input := createTestInput(t)
	input.DryRun = false
	task := findBashTask(t, input, "test-echo")
	task.Enabled = false
	ctx := context.Background()

	// when
	got, err := task.Install(ctx, input)
	// then
	assert.NoError(t, err)
	wantedResult := expectTaskResult(t, false, task.Attributes, input)
	assert.Equalf(t, wantedResult, got, "Install(%v) should skip disabled task", input)
}

func TestBashTask_List(t *testing.T) {
	input := createTestInput(t)
	task := findBashTask(t, input, "test-echo")
	ctx := context.Background()

	// when
	got, err := task.List(ctx, input)
	// then
	assert.NoError(t, err)
	wantedResult := expectTaskResult(t, true, task.Attributes, input)
	assert.Equalf(t, wantedResult, got, "List(%v) should log command", input)
}
