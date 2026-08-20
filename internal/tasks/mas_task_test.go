package tasks

import (
	"context"
	"testing"

	ty "github.com/alanktwong/dum/internal/types"

	"github.com/stretchr/testify/assert"
)

func TestNewMasTask_NilAttributes(t *testing.T) {
	task, err := NewMasTask(nil)
	assert.Error(t, err)
	assert.Nil(t, task)
	assert.Contains(t, err.Error(), "attributes cannot be nil")
}

func TestNewMasTask_EmptyID(t *testing.T) {
	attrs := &ty.Attributes{
		ID:          "",
		Description: "test",
	}
	task, err := NewMasTask(attrs)
	assert.Error(t, err)
	assert.Nil(t, task)
	assert.Contains(t, err.Error(), "task ID cannot be empty")
}

func TestNewMasTask_Success(t *testing.T) {
	attrs := &ty.Attributes{
		ID:          "123456789",
		Description: "test mas task",
		Enabled:     true,
	}
	task, err := NewMasTask(attrs)
	assert.NoError(t, err)
	assert.NotNil(t, task)
	assert.Equal(t, "123456789", task.ID)
}

func TestMasTask_Install_EmptyDescription(t *testing.T) {
	attrs := &ty.Attributes{
		ID:          "123456789",
		Description: "",
		Enabled:     true,
	}
	task := &MasTask{
		Attributes: *attrs,
		Mas:        new(MockMas),
		Log:        NewTestLogger(t),
	}

	input := &ty.TaskInput{
		Play:   "test-play",
		DryRun: false,
	}

	result, err := task.Install(context.Background(), input)
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.False(t, result.Success)
}

func TestMasTask_Install_Disabled(t *testing.T) {
	attrs := &ty.Attributes{
		ID:          "123456789",
		Description: "test mas task",
		Enabled:     false,
	}
	task := &MasTask{
		Attributes: *attrs,
		Mas:        new(MockMas),
		Log:        NewTestLogger(t),
	}

	input := &ty.TaskInput{
		Play:   "test-play",
		DryRun: false,
	}

	result, err := task.Install(context.Background(), input)
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.False(t, result.Success)
}

func TestMasTask_Install_AlreadyInstalled(t *testing.T) {
	attrs := &ty.Attributes{
		ID:          "123456789",
		Description: "test mas task",
		Enabled:     true,
	}
	mockMas := new(MockMas)
	mockMas.On("List", context.Background()).Return("123456789 SomeApp\n987654321 AnotherApp", nil)

	task := &MasTask{
		Attributes: *attrs,
		Mas:        mockMas,
		Log:        NewTestLogger(t),
	}

	input := &ty.TaskInput{
		Play:   "test-play",
		DryRun: false,
	}

	result, err := task.Install(context.Background(), input)
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.False(t, result.Success)
	mockMas.AssertExpectations(t)
}

func TestMasTask_Install_DryRun(t *testing.T) {
	attrs := &ty.Attributes{
		ID:          "123456789",
		Description: "test mas task",
		Enabled:     true,
	}
	mockMas := new(MockMas)
	mockMas.On("List", context.Background()).Return("987654321 AnotherApp", nil)

	task := &MasTask{
		Attributes: *attrs,
		Mas:        mockMas,
		Log:        NewTestLogger(t),
	}

	input := &ty.TaskInput{
		Play:   "test-play",
		DryRun: true,
	}

	result, err := task.Install(context.Background(), input)
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.True(t, result.Success)
	mockMas.AssertExpectations(t)
}

func TestMasTask_Install_Success(t *testing.T) {
	attrs := &ty.Attributes{
		ID:          "123456789",
		Description: "test mas task",
		Enabled:     true,
	}
	mockMas := new(MockMas)
	mockMas.On("List", context.Background()).Return("987654321 AnotherApp", nil)
	mockMas.On("Install", context.Background(), "123456789").Return(nil)

	task := &MasTask{
		Attributes: *attrs,
		Mas:        mockMas,
		Log:        NewTestLogger(t),
	}

	input := &ty.TaskInput{
		Play:   "test-play",
		DryRun: false,
	}

	result, err := task.Install(context.Background(), input)
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.True(t, result.Success)
	mockMas.AssertExpectations(t)
}

func TestMasTask_Install_ListError(t *testing.T) {
	attrs := &ty.Attributes{
		ID:          "123456789",
		Description: "test mas task",
		Enabled:     true,
	}
	mockMas := new(MockMas)
	mockMas.On("List", context.Background()).Return("", assert.AnError)

	task := &MasTask{
		Attributes: *attrs,
		Mas:        mockMas,
		Log:        NewTestLogger(t),
	}

	input := &ty.TaskInput{
		Play:   "test-play",
		DryRun: false,
	}

	result, err := task.Install(context.Background(), input)
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "failed to list mas apps")
	mockMas.AssertExpectations(t)
}

func TestMasTask_Install_InstallError(t *testing.T) {
	attrs := &ty.Attributes{
		ID:          "123456789",
		Description: "test mas task",
		Enabled:     true,
	}
	mockMas := new(MockMas)
	mockMas.On("List", context.Background()).Return("987654321 AnotherApp", nil)
	mockMas.On("Install", context.Background(), "123456789").Return(assert.AnError)

	task := &MasTask{
		Attributes: *attrs,
		Mas:        mockMas,
		Log:        NewTestLogger(t),
	}

	input := &ty.TaskInput{
		Play:   "test-play",
		DryRun: false,
	}

	result, err := task.Install(context.Background(), input)
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "mas install")
	mockMas.AssertExpectations(t)
}

func TestMasTask_GetAttributes(t *testing.T) {
	attrs := &ty.Attributes{
		ID:          "123456789",
		Description: "test mas task",
		Enabled:     true,
	}
	task, err := NewMasTask(attrs)
	assert.NoError(t, err)

	result := task.GetAttributes()
	assert.Equal(t, "123456789", result.ID)
}

func TestMasTask_GetID(t *testing.T) {
	attrs := &ty.Attributes{
		ID:          "123456789",
		Description: "test mas task",
		Enabled:     true,
	}
	task, err := NewMasTask(attrs)
	assert.NoError(t, err)

	assert.Equal(t, "123456789", task.GetID())
}

func TestMasTask_IsEnabled(t *testing.T) {
	attrs := &ty.Attributes{
		ID:          "123456789",
		Description: "test mas task",
		Enabled:     true,
	}
	task, err := NewMasTask(attrs)
	assert.NoError(t, err)

	assert.True(t, task.IsEnabled())
}

func TestMasTask_Install_Disabled_DryRun(t *testing.T) {
	attrs := &ty.Attributes{
		ID:          "123456789",
		Description: "test mas task",
		Enabled:     false,
	}
	mockLog := NewTestLogger(t)
	mockLog.On(
		"Printlnf",
		"%v %s: mas install %s ... desc: %s",
		[]any{TaskEllipsis, "MasTask", "123456789", "test mas task"},
	).Return(nil)

	task := &MasTask{
		Attributes: *attrs,
		Mas:        new(MockMas),
		Log:        mockLog,
	}

	input := &ty.TaskInput{
		Play:   "test-play",
		DryRun: true,
	}

	result, err := task.Install(context.Background(), input)
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.False(t, result.Success)
	mockLog.AssertExpectations(t)
}
