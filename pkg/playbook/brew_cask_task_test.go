package playbook

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func findBrewCaskTask(t *testing.T, input *Input) *BrewCaskTask {
	id := "visual-studio-code"
	tk, err := findTask(input.PlayBook, id)
	assert.NoError(t, err)
	task, ok := tk.(*BrewCaskTask)
	assert.True(t, ok)
	return task
}

func TestBrewCaskTask_Install(t *testing.T) {
	input := createTestInput(t)
	input.DryRun = false
	task := findBrewCaskTask(t, input)
	ctx := context.Background()

	mockBrew := NewMockBrew(t)
	mockBrew.EXPECT().InPath(ctx, "Caskroom", task.ID).Return(false).Once()
	mockBrew.EXPECT().InstallCask(ctx, task.ID).Return(nil).Once()
	mockUtils := NewMockExt(t)
	mockUtils.EXPECT().IsOSX().Return(true).Once()

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

func TestBrewCaskTask_FailTap(t *testing.T) {
	input := createTestInput(t)
	input.DryRun = false
	task := findBrewCaskTask(t, input)
	ctx := context.Background()

	mockBrew := NewMockBrew(t)
	mockUtils := NewMockExt(t)
	mockUtils.EXPECT().IsOSX().Return(true).Once()

	mockInstaller := NewMockInstaller(t)
	mockInstaller.EXPECT().Install(ctx, input).
		Return(nil, fmt.Errorf("failed to install")).Once()

	task.Brew = mockBrew
	task.Utils = mockUtils
	task.Tap = mockInstaller
	// when
	got, err := task.Install(ctx, input)
	// then
	assert.Errorf(t, err, "Install(%v) can fail", input)
	assert.Nilf(t, got, "Install(%v)", input)
	mockInstaller.AssertExpectations(t)
	mockBrew.AssertNotCalled(t, "InPath")
	mockBrew.AssertNotCalled(t, "InstallCask")
}

func TestBrewCaskTask_Fail(t *testing.T) {
	input := createTestInput(t)
	input.DryRun = false
	task := findBrewCaskTask(t, input)
	ctx := context.Background()

	mockBrew := NewMockBrew(t)
	mockBrew.EXPECT().InPath(ctx, "Caskroom", task.ID).Return(false).Once()
	mockBrew.EXPECT().InstallCask(ctx, task.ID).Return(fmt.Errorf("some error")).Once()
	mockUtils := NewMockExt(t)
	mockUtils.EXPECT().IsOSX().Return(true).Once()

	mockInstaller := new(MockInstaller)
	mockInstaller.EXPECT().Install(ctx, input).Return(NewTaskResult(input, false))

	task.Brew = mockBrew
	task.Utils = mockUtils
	task.Tap = mockInstaller
	// when
	got, err := task.Install(ctx, input)
	// then
	assert.Errorf(t, err, "Install(%v) can fail", input)
	assert.Nilf(t, got, "Install(%v)", input)
	mockInstaller.AssertExpectations(t)
}

func TestBrewCaskTask_Disabled(t *testing.T) {
	input := createTestInput(t)
	input.DryRun = false
	task := findBrewCaskTask(t, input)
	task.Enabled = false
	ctx := context.Background()

	mockBrew := NewMockBrew(t)
	mockUtils := NewMockExt(t)
	mockInstaller := new(MockInstaller)

	task.Brew = mockBrew
	task.Utils = mockUtils
	task.Tap = mockInstaller
	// when
	got, err := task.Install(ctx, input)
	// then
	assert.NoError(t, err)
	wantedResult := expectTaskResult(t, false, task.Attributes, input)
	assert.Equalf(t, wantedResult, got, "Install(%v)", input)
	mockInstaller.AssertExpectations(t)
	mockUtils.AssertNotCalled(t, "IsOSX")
	mockBrew.AssertNotCalled(t, "InPath")
	mockBrew.AssertNotCalled(t, "InstallCask")
}

func TestBrewCaskTask_DryRun(t *testing.T) {
	input := createTestInput(t)
	input.DryRun = true
	task := findBrewCaskTask(t, input)
	ctx := context.Background()

	mockBrew := NewMockBrew(t)
	mockBrew.EXPECT().InPath(ctx, "Caskroom", task.ID).Return(false).Once()
	mockUtils := NewMockExt(t)
	mockUtils.On("IsOSX").Return(true).Once()
	mockInstaller := new(MockInstaller)
	mockInstaller.EXPECT().Install(ctx, input).Return(NewTaskResult(input, false))

	task.Brew = mockBrew
	task.Utils = mockUtils
	task.Tap = mockInstaller
	// when
	got, err := task.Install(ctx, input)
	// then
	assert.NoError(t, err)
	wantedResult := expectTaskResult(t, true, task.Attributes, input)
	assert.Equalf(t, wantedResult, got, "Install(%v) should succeed for dry run", input)
	mockInstaller.AssertExpectations(t)
	mockBrew.AssertNotCalled(t, "InstallCask")
}

func TestBrewCaskTask_NotOSX(t *testing.T) {
	input := createTestInput(t)
	input.DryRun = false
	task := findBrewCaskTask(t, input)
	ctx := context.Background()

	mockBrew := NewMockBrew(t)
	mockUtils := NewMockExt(t)
	mockUtils.EXPECT().IsOSX().Return(false).Once()
	mockInstaller := NewMockInstaller(t)

	task.Brew = mockBrew
	task.Utils = mockUtils
	task.Tap = mockInstaller
	// when
	got, err := task.Install(ctx, input)
	// then
	assert.Errorf(t, err, "Install(%v) wont work outside OSX", input)
	assert.Nilf(t, got, "Install(%v)", input)
	mockInstaller.AssertExpectations(t)
	mockBrew.AssertNotCalled(t, "InPath")
	mockBrew.AssertNotCalled(t, "InstallCask")
}
