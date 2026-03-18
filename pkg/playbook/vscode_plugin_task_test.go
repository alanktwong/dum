package playbook

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func findVsCodeTask(t *testing.T, input *Input) *VsCodePluginTask {
	id := "vscodevim.vim"
	tk, err := findTask(input.PlayBook, id)
	assert.NoError(t, err)
	task, ok := tk.(*VsCodePluginTask)
	assert.True(t, ok)
	return task
}

func TestVsCodePluginTask_Install(t *testing.T) {
	input := createTestInput(t)
	input.DryRun = false
	task := findVsCodeTask(t, input)
	ctx := context.Background()

	mockCode := NewMockCode(t)
	mockCode.EXPECT().ListExtensions(ctx).Return("", nil).Once()
	mockCode.EXPECT().InstallExtension(ctx, task.ID).Return(nil).Once()
	task.Code = mockCode
	// when
	got, err := task.Install(ctx, input)
	// then
	assert.NoError(t, err)
	wanted := expectTaskResult(t, true, task.Attributes, input)
	assert.Equalf(t, wanted, got, "Install(%v) should succeed", input)
}

func TestVsCodePluginTask_Install_Fail(t *testing.T) {
	input := createTestInput(t)
	input.DryRun = false
	task := findVsCodeTask(t, input)
	ctx := context.Background()

	mockCode := NewMockCode(t)
	mockCode.EXPECT().ListExtensions(ctx).Return("", nil).Once()
	mockCode.EXPECT().InstallExtension(ctx, task.ID).Return(fmt.Errorf("VS Code plugin failed")).Once()
	task.Code = mockCode
	// when
	got, err := task.Install(ctx, input)
	// then
	assert.Errorf(t, err, "Install(%v) should fail", input)
	assert.Nil(t, got)
}

func TestVsCodePluginTask_Install_FailList(t *testing.T) {
	input := createTestInput(t)
	input.DryRun = false
	task := findVsCodeTask(t, input)
	ctx := context.Background()

	mockCode := NewMockCode(t)
	mockCode.EXPECT().ListExtensions(ctx).Return("", fmt.Errorf("fail")).Once()
	task.Code = mockCode
	// when
	got, err := task.Install(ctx, input)
	// then
	assert.Errorf(t, err, "Install(%v) should fail", input)
	assert.Nil(t, got)
}

func TestVsCodePluginTask_Disabled(t *testing.T) {
	input := createTestInput(t)
	input.DryRun = false
	task := findVsCodeTask(t, input)
	task.Enabled = false
	ctx := context.Background()

	mockCode := NewMockCode(t)
	task.Code = mockCode
	// when
	got, err := task.Install(ctx, input)
	// then
	assert.NoError(t, err)
	wanted := expectTaskResult(t, false, task.Attributes, input)
	assert.Equalf(t, wanted, got, "Install(%v) should be disabled", input)
	mockCode.AssertNotCalled(t, "ListExtensions")
	mockCode.AssertNotCalled(t, "InstallExtension")
}

func TestVsCodePluginTask_AlreadyInstalled(t *testing.T) {
	input := createTestInput(t)
	input.DryRun = false
	task := findVsCodeTask(t, input)
	ctx := context.Background()

	mockCode := NewMockCode(t)
	mockCode.EXPECT().ListExtensions(ctx).Return(task.ID, nil).Once()
	task.Code = mockCode
	// when
	got, err := task.Install(ctx, input)
	// then
	assert.NoError(t, err)
	wanted := expectTaskResult(t, true, task.Attributes, input)
	assert.Equalf(t, wanted, got, "Install(%v) should succeed when already installed", input)
	mockCode.AssertNotCalled(t, "InstallExtension")
}

func TestVsCodePluginTask_DryRun(t *testing.T) {
	input := createTestInput(t)
	input.DryRun = true
	task := findVsCodeTask(t, input)
	ctx := context.Background()

	mockCode := NewMockCode(t)
	mockCode.EXPECT().ListExtensions(ctx).Return("", nil).Once()
	task.Code = mockCode
	// when
	got, err := task.Install(ctx, input)
	// then
	assert.NoError(t, err)
	wanted := expectTaskResult(t, true, task.Attributes, input)
	assert.Equalf(t, wanted, got, "Install(%v) should succeed for dry run", input)
}
