package playbook

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func findBrewCellarTask(t *testing.T, input *Input) *BrewCellarTask {
	id := "boost"
	tk, err := findTask(input.PlayBook, id)
	assert.NoError(t, err)
	task, ok := tk.(*BrewCellarTask)
	assert.True(t, ok)
	return task
}

func TestBrewCellarTask_Install(t *testing.T) {
	input := createTestInput(t)
	input.DryRun = false
	task := findBrewCellarTask(t, input)
	ctx := context.Background()

	mockBrew := NewMockBrew(t)
	mockBrew.EXPECT().InPath(ctx, "Cellar", task.ID).Return(false).Once()
	mockBrew.EXPECT().Install(ctx, task.ID).Return(nil)

	task.Brew = mockBrew
	// when
	got, err := task.Install(ctx, input)
	// then
	assert.NoError(t, err)
	wantedResult := expectTaskResult(t, true, task.Attributes, input)
	assert.Equalf(t, wantedResult, got, "Install(%v)", input)
}

func TestBrewCellarTask_Install_Fail(t *testing.T) {
	input := createTestInput(t)
	input.DryRun = false
	task := findBrewCellarTask(t, input)
	ctx := context.Background()

	mockBrew := NewMockBrew(t)
	mockBrew.EXPECT().InPath(ctx, "Cellar", task.ID).Return(false).Once()
	mockBrew.EXPECT().Install(ctx, task.ID).Return(fmt.Errorf("some error"))

	task.Brew = mockBrew
	// when
	got, err := task.Install(ctx, input)
	// then
	assert.Error(t, err)
	assert.Nilf(t, got, "Install(%v) can fail", input)
}
