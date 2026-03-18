package playbook

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func findFunctionTask(t *testing.T, input *Input) *FunctionTask {
	id := "install_test"
	tk, err := findTask(input.PlayBook, id)
	assert.NoError(t, err)
	task, ok := tk.(*FunctionTask)
	assert.True(t, ok)
	return task
}

func TestFunctionTask_Install(t *testing.T) {
	input := createTestInput(t)
	input.DryRun = false
	task := findFunctionTask(t, input)
	input.Task = task.ID
	ctx := context.Background()

	mockInstaller := NewMockInstaller(t)
	mockRegistry := make(map[string]Installer)
	wantedResult := expectTaskResult(t, true, task.Attributes, input)
	mockInstaller.EXPECT().Install(ctx, input).Return(wantedResult, nil).Once()
	mockRegistry[task.ID] = mockInstaller
	task.Registry = mockRegistry

	// when
	got, err := task.Install(ctx, input)
	// then
	assert.NoError(t, err)
	assert.Equalf(t, wantedResult, got, "Install(%v): should install for happy path", input)
}

func TestFunctionTask_MissRegistry(t *testing.T) {
	input := createTestInput(t)
	input.DryRun = false
	task := findFunctionTask(t, input)
	input.Task = task.ID
	ctx := context.Background()

	mockInstaller := NewMockInstaller(t)
	mockRegistry := make(map[string]Installer)
	task.Registry = mockRegistry

	// when
	got, err := task.Install(ctx, input)
	// then
	assert.Errorf(t, err, "Install(%v): should install fail if task.ID is not in registry", input)
	assert.Nilf(t, got, "Install(%v): should install be nil", input)
	mockInstaller.AssertNotCalled(t, "Install")
}
