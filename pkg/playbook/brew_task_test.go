package playbook

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func findBrewTask(t *testing.T, input *Input) *BrewTask {
	tk, err := findTask(input.PlayBook, "scc")
	assert.NoError(t, err)
	task, ok := tk.(*BrewTask)
	assert.True(t, ok)
	return task
}

func TestBrewTask_Install(t *testing.T) {
	input := createTestInput(t)
	input.DryRun = false
	task := findBrewTask(t, input)
	ctx := context.Background()

	mockBrew := NewMockBrew(t)
	mockBrew.EXPECT().InPath(ctx, "opt", task.ID).Return(false).Once()
	mockBrew.EXPECT().Install(ctx, task.ID).Return(nil).Once()
	mockUtils := NewMockExt(t)
	mockUtils.EXPECT().IsInstalled(task.ID).Return(false).Once()

	mockInstaller := NewMockInstaller(t)
	mockInstaller.EXPECT().Install(ctx, input).Return(NewTaskResult(input, false))

	task.Brew = mockBrew
	task.Utils = mockUtils
	task.Tap = mockInstaller
	// when
	got, err := task.Install(ctx, input)
	// then
	assert.NoError(t, err)
	wantedResult := expectTaskResult(t, true, task.Attributes, input)
	assert.Equalf(t, wantedResult, got, "Install(%v)", input)
	mockInstaller.AssertExpectations(t)
}

func TestBrewTask_Install_Fail(t *testing.T) {
	input := createTestInput(t)
	input.DryRun = false
	task := findBrewTask(t, input)
	ctx := context.Background()

	mockBrew := NewMockBrew(t)
	mockBrew.EXPECT().InPath(ctx, "opt", task.ID).Return(false).Once()
	mockBrew.EXPECT().Install(ctx, task.ID).Return(fmt.Errorf("some error")).Once()
	task.Brew = mockBrew

	mockUtils := NewMockExt(t)
	mockUtils.EXPECT().IsInstalled(task.ID).Return(false).Once()
	task.Utils = mockUtils

	mockInstaller := NewMockInstaller(t)
	mockInstaller.EXPECT().Install(ctx, input).Return(NewTaskResult(input, false))
	task.Tap = mockInstaller

	// when
	got, err := task.Install(ctx, input)
	// then
	assert.Nil(t, got)
	assert.Errorf(t, err, "should fail when install(%v) fails", input)
}

func TestBrewTask_Install_FailTap(t *testing.T) {
	input := createTestInput(t)
	input.DryRun = false
	task := findBrewTask(t, input)
	ctx := context.Background()

	mockBrew := NewMockBrew(t)
	task.Brew = mockBrew
	mockUtils := NewMockExt(t)
	task.Utils = mockUtils
	mockInstaller := NewMockInstaller(t)
	mockInstaller.EXPECT().Install(ctx, input).Return(nil, fmt.Errorf("some error"))
	task.Tap = mockInstaller
	// when
	got, err := task.Install(ctx, input)
	// then
	assert.Errorf(t, err, "install(%v) should fail when tap fails", input)
	assert.Nil(t, got)
	mockBrew.AssertNotCalled(t, "Tap")
	mockBrew.AssertNotCalled(t, "InPath")
	mockUtils.AssertNotCalled(t, "IsInstalled")
}

func TestBrewTask_Install_Prior(t *testing.T) {
	input := createTestInput(t)
	input.DryRun = false
	task := findBrewTask(t, input)
	ctx := context.Background()

	mockBrew := NewMockBrew(t)
	mockBrew.EXPECT().InPath(ctx, "opt", task.ID).Return(true).Once()
	task.Brew = mockBrew

	mockUtils := NewMockExt(t)
	mockUtils.EXPECT().IsInstalled(task.ID).Return(false).Once()
	task.Utils = mockUtils
	// when
	got, err := task.Install(ctx, input)
	// then
	assert.NoError(t, err)
	wantedResult := expectTaskResult(t, false, task.Attributes, input)
	assert.Equalf(t, wantedResult, got, "should succeed but not install(%v) what is already installed", input)
	mockBrew.AssertNotCalled(t, "Install")
}

func TestBrewTask_Install_Disabled(t *testing.T) {
	input := createTestInput(t)
	input.DryRun = false
	task := findBrewTask(t, input)
	task.Enabled = false
	ctx := context.Background()

	mockBrew := NewMockBrew(t)
	task.Brew = mockBrew
	mockUtils := NewMockExt(t)
	task.Utils = mockUtils
	// when
	got, err := task.Install(ctx, input)
	// then
	assert.NoError(t, err)
	wantedResult := expectTaskResult(t, false, task.Attributes, input)
	assert.Equalf(t, wantedResult, got, "should succeed since it is disabled")
	mockBrew.AssertNotCalled(t, "Install")
	mockBrew.AssertNotCalled(t, "InPath")
	mockUtils.AssertNotCalled(t, "IsInstalled")
}

func TestBrewTask_Install_DryRun(t *testing.T) {
	input := createTestInput(t)
	input.DryRun = true
	task := findBrewTask(t, input)
	ctx := context.Background()

	mockBrew := NewMockBrew(t)
	mockBrew.EXPECT().InPath(ctx, "opt", task.ID).Return(false).Once()

	task.Brew = mockBrew
	mockUtils := NewMockExt(t)
	mockUtils.EXPECT().IsInstalled(task.ID).Return(false).Once()
	task.Utils = mockUtils
	// when
	got, err := task.Install(ctx, input)
	// then
	assert.NoError(t, err)
	wantedResult := expectTaskResult(t, true, task.Attributes, input)
	assert.Equal(t, wantedResult, got, "should succeed since it is a dry run")
	mockBrew.AssertNotCalled(t, "Install")
}
