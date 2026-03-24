package tasks

import (
	"context"
	"testing"

	i "awong/dotfiles/pkg/tasks/installer"
	ty "awong/dotfiles/pkg/types"

	"github.com/stretchr/testify/assert"
)

func TestNewBashTask_NilAttributes(t *testing.T) {
	task, err := NewBashTask(nil, "echo hello", "")
	assert.Error(t, err)
	assert.Nil(t, task)
	assert.Contains(t, err.Error(), "attributes cannot be nil")
}

func TestNewBashTask_EmptyID(t *testing.T) {
	attrs := &ty.Attributes{
		ID:          "",
		Description: "test",
	}
	task, err := NewBashTask(attrs, "echo hello", "")
	assert.Error(t, err)
	assert.Nil(t, task)
	assert.Contains(t, err.Error(), "task ID cannot be empty")
}

func TestNewBashTask_EmptyCommandAndScript(t *testing.T) {
	attrs := &ty.Attributes{
		ID:          "test-bash",
		Description: "test",
	}
	task, err := NewBashTask(attrs, "", "")
	assert.Error(t, err)
	assert.Nil(t, task)
	assert.Contains(t, err.Error(), "either command or script must be specified")
}

func TestNewBashTask_SuccessWithCommand(t *testing.T) {
	attrs := &ty.Attributes{
		ID:          "test-bash",
		Description: "test bash task",
		Enabled:     true,
	}
	task, err := NewBashTask(attrs, "echo hello", "")
	assert.NoError(t, err)
	assert.NotNil(t, task)
	assert.Equal(t, "test-bash", task.ID)
	assert.Equal(t, "echo hello", task.Command)
}

func TestNewBashTask_SuccessWithScript(t *testing.T) {
	attrs := &ty.Attributes{
		ID:          "test-bash",
		Description: "test bash task",
		Enabled:     true,
	}
	task, err := NewBashTask(attrs, "", "echo hello")
	assert.NoError(t, err)
	assert.NotNil(t, task)
	assert.Equal(t, "test-bash", task.ID)
	assert.Equal(t, "echo hello", task.Script)
}

func TestBashTask_Install_Disabled(t *testing.T) {
	attrs := &ty.Attributes{
		ID:          "test-bash",
		Description: "test bash task",
		Enabled:     false,
	}
	task, err := NewBashTask(attrs, "echo hello", "")
	assert.NoError(t, err)

	input := &ty.TaskInput{
		Play:   "test-play",
		DryRun: false,
	}

	result, err := task.Install(context.Background(), input)
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.False(t, result.Success)
}

func TestBashTask_Install_DryRun(t *testing.T) {
	attrs := &ty.Attributes{
		ID:          "test-bash",
		Description: "test bash task",
		Enabled:     true,
	}
	task, err := NewBashTask(attrs, "echo hello", "")
	assert.NoError(t, err)

	input := &ty.TaskInput{
		Play:   "test-play",
		DryRun: true,
	}

	result, err := task.Install(context.Background(), input)
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.True(t, result.Success)
}

func TestBashTask_Install_Success(t *testing.T) {
	attrs := &ty.Attributes{
		ID:          "test-bash",
		Description: "test bash task",
		Enabled:     true,
	}
	mockExt := new(MockExt)
	mockExt.On("RunCommand", context.Background(), "echo hello", false).Return(nil)

	mockLog := i.NewMockLogger(t)
	mockLog.On("Debugf", "%s START ... play: %s taskID: %s", []any{"...........", "test-play", "test-bash"}).Return()
	mockLog.On("Infof", "%s %s: %s", []any{"...........", "test-play", "echo hello"}).Return()

	task := &BashTask{
		Attributes: *attrs,
		Command:    "echo hello",
		Utils:      mockExt,
		Log:        mockLog,
	}

	input := &ty.TaskInput{
		Play:   "test-play",
		DryRun: false,
	}

	result, err := task.Install(context.Background(), input)
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.True(t, result.Success)
	mockExt.AssertExpectations(t)
	mockLog.AssertExpectations(t)
}

func TestBashTask_Install_Failure(t *testing.T) {
	attrs := &ty.Attributes{
		ID:          "test-bash",
		Description: "test bash task",
		Enabled:     true,
	}
	mockExt := new(MockExt)
	mockExt.On("RunCommand", context.Background(), "echo hello", false).Return(assert.AnError)

	mockLog := i.NewMockLogger(t)
	mockLog.On("Debugf", "%s START ... play: %s taskID: %s", []any{"...........", "test-play", "test-bash"}).Return()
	mockLog.On("Infof", "%s %s: %s", []any{"...........", "test-play", "echo hello"}).Return()

	task := &BashTask{
		Attributes: *attrs,
		Command:    "echo hello",
		Utils:      mockExt,
		Log:        mockLog,
	}

	input := &ty.TaskInput{
		Play:   "test-play",
		DryRun: false,
	}

	result, err := task.Install(context.Background(), input)
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "failed to run command")
	mockExt.AssertExpectations(t)
	mockLog.AssertExpectations(t)
}

func TestBashTask_List(t *testing.T) {
	attrs := &ty.Attributes{
		ID:          "test-bash",
		Description: "test bash task",
		Enabled:     true,
	}
	mockLog := i.NewMockLogger(t)
	mockLog.On("Printlnf", "%s %s", []any{"...........", "echo hello"}).Return(nil).Maybe()

	task := &BashTask{
		Attributes: *attrs,
		Command:    "echo hello",
		Log:        mockLog,
	}

	input := &ty.TaskInput{
		Play: "test-play",
	}

	result, err := task.List(context.Background(), input)
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.True(t, result.Success)
	mockLog.AssertExpectations(t)
}

func TestBashTask_GetAttributes(t *testing.T) {
	attrs := &ty.Attributes{
		ID:          "test-bash",
		Description: "test bash task",
		Enabled:     true,
	}
	task, err := NewBashTask(attrs, "echo hello", "")
	assert.NoError(t, err)

	result := task.GetAttributes()
	assert.Equal(t, "test-bash", result.ID)
}

func TestBashTask_GetID(t *testing.T) {
	attrs := &ty.Attributes{
		ID:          "test-bash",
		Description: "test bash task",
		Enabled:     true,
	}
	task, err := NewBashTask(attrs, "echo hello", "")
	assert.NoError(t, err)

	assert.Equal(t, "test-bash", task.GetID())
}

func TestBashTask_IsEnabled(t *testing.T) {
	attrs := &ty.Attributes{
		ID:          "test-bash",
		Description: "test bash task",
		Enabled:     true,
	}
	task, err := NewBashTask(attrs, "echo hello", "")
	assert.NoError(t, err)

	assert.True(t, task.IsEnabled())
}
