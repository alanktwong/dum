package playbook

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func findMasTask(t *testing.T, input *Input) *MasTask {
	id := "462058435"
	tk, err := findTask(input.PlayBook, id)
	assert.NoError(t, err)
	task, ok := tk.(*MasTask)
	assert.True(t, ok)
	return task
}

func TestMasTask_Install(t *testing.T) {
	input := createTestInput(t)
	input.DryRun = false
	task := findMasTask(t, input)
	ctx := context.Background()

	mockMas := NewMockMas(t)
	mockMas.EXPECT().List(ctx).Return("", nil).Once()
	mockMas.EXPECT().Install(ctx, task.ID).Return(nil).Once()

	task.Mas = mockMas
	// when
	got, err := task.Install(ctx, input)
	// then
	assert.NoError(t, err)
	wantedResult := expectTaskResult(t, true, task.Attributes, input)
	assert.Equalf(t, wantedResult, got, "Install(%v) should mkdir", input)
}
