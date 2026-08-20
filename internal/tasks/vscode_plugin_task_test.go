package tasks

import (
	"context"
	"testing"

	ty "github.com/alanktwong/dum/internal/types"

	"github.com/stretchr/testify/assert"
)

func TestNewVsCodePluginTask_NilAttributes(t *testing.T) {
	task, err := NewVsCodePluginTask(nil)
	assert.Error(t, err)
	assert.Nil(t, task)
	assert.Contains(t, err.Error(), "attributes cannot be nil")
}

func TestNewVsCodePluginTask_EmptyID(t *testing.T) {
	attrs := &ty.Attributes{
		ID:          "",
		Description: "test",
	}
	task, err := NewVsCodePluginTask(attrs)
	assert.Error(t, err)
	assert.Nil(t, task)
	assert.Contains(t, err.Error(), "task ID cannot be empty")
}

func TestNewVsCodePluginTask_Success(t *testing.T) {
	attrs := &ty.Attributes{
		ID:          "ms-python.python",
		Description: "test vscode task",
		Enabled:     true,
	}
	task, err := NewVsCodePluginTask(attrs)
	assert.NoError(t, err)
	assert.NotNil(t, task)
	assert.Equal(t, "ms-python.python", task.ID)
}

func TestVsCodePluginTask_Install_Disabled(t *testing.T) {
	attrs := &ty.Attributes{
		ID:          "ms-python.python",
		Description: "test vscode task",
		Enabled:     false,
	}
	task := &VsCodePluginTask{
		Attributes: *attrs,
		Code:       new(MockCode),
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

func TestVsCodePluginTask_Install_AlreadyInstalled(t *testing.T) {
	attrs := &ty.Attributes{
		ID:          "ms-python.python",
		Description: "test vscode task",
		Enabled:     true,
	}
	mockCode := new(MockCode)
	mockCode.On("ListExtensions", context.Background()).Return("ms-python.python\nms-vscode.cpptools", nil)

	task := &VsCodePluginTask{
		Attributes: *attrs,
		Code:       mockCode,
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
	mockCode.AssertExpectations(t)
}

func TestVsCodePluginTask_Install_DryRun(t *testing.T) {
	attrs := &ty.Attributes{
		ID:          "ms-python.python",
		Description: "test vscode task",
		Enabled:     true,
	}
	mockCode := new(MockCode)
	mockCode.On("ListExtensions", context.Background()).Return("ms-vscode.cpptools", nil)

	task := &VsCodePluginTask{
		Attributes: *attrs,
		Code:       mockCode,
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
	mockCode.AssertExpectations(t)
}

func TestVsCodePluginTask_Install_Success(t *testing.T) {
	attrs := &ty.Attributes{
		ID:          "ms-python.python",
		Description: "test vscode task",
		Enabled:     true,
	}
	mockCode := new(MockCode)
	mockCode.On("ListExtensions", context.Background()).Return("ms-vscode.cpptools", nil)
	mockCode.On("InstallExtension", context.Background(), "ms-python.python").Return(nil)

	task := &VsCodePluginTask{
		Attributes: *attrs,
		Code:       mockCode,
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
	mockCode.AssertExpectations(t)
}

func TestVsCodePluginTask_Install_ListError(t *testing.T) {
	attrs := &ty.Attributes{
		ID:          "ms-python.python",
		Description: "test vscode task",
		Enabled:     true,
	}
	mockCode := new(MockCode)
	mockCode.On("ListExtensions", context.Background()).Return("", assert.AnError)

	task := &VsCodePluginTask{
		Attributes: *attrs,
		Code:       mockCode,
		Log:        NewTestLogger(t),
	}

	input := &ty.TaskInput{
		Play:   "test-play",
		DryRun: false,
	}

	result, err := task.Install(context.Background(), input)
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "failed to list VS Code extensions")
	mockCode.AssertExpectations(t)
}

func TestVsCodePluginTask_Install_InstallError(t *testing.T) {
	attrs := &ty.Attributes{
		ID:          "ms-python.python",
		Description: "test vscode task",
		Enabled:     true,
	}
	mockCode := new(MockCode)
	mockCode.On("ListExtensions", context.Background()).Return("ms-vscode.cpptools", nil)
	mockCode.On("InstallExtension", context.Background(), "ms-python.python").Return(assert.AnError)

	task := &VsCodePluginTask{
		Attributes: *attrs,
		Code:       mockCode,
		Log:        NewTestLogger(t),
	}

	input := &ty.TaskInput{
		Play:   "test-play",
		DryRun: false,
	}

	result, err := task.Install(context.Background(), input)
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "code --install-extension")
	mockCode.AssertExpectations(t)
}

func TestVsCodePluginTask_GetAttributes(t *testing.T) {
	attrs := &ty.Attributes{
		ID:          "ms-python.python",
		Description: "test vscode task",
		Enabled:     true,
	}
	task, err := NewVsCodePluginTask(attrs)
	assert.NoError(t, err)

	result := task.GetAttributes()
	assert.Equal(t, "ms-python.python", result.ID)
}

func TestVsCodePluginTask_GetID(t *testing.T) {
	attrs := &ty.Attributes{
		ID:          "ms-python.python",
		Description: "test vscode task",
		Enabled:     true,
	}
	task, err := NewVsCodePluginTask(attrs)
	assert.NoError(t, err)

	assert.Equal(t, "ms-python.python", task.GetID())
}

func TestVsCodePluginTask_IsEnabled(t *testing.T) {
	attrs := &ty.Attributes{
		ID:          "ms-python.python",
		Description: "test vscode task",
		Enabled:     true,
	}
	task, err := NewVsCodePluginTask(attrs)
	assert.NoError(t, err)

	assert.True(t, task.IsEnabled())
}

func TestVsCodePluginTask_Install_Disabled_DryRun(t *testing.T) {
	attrs := &ty.Attributes{
		ID:          "ms-python.python",
		Description: "test vscode task",
		Enabled:     false,
	}
	mockLog := NewTestLogger(t)
	mockLog.On(
		"Printlnf",
		"%v %s: code --install-extensions %s",
		[]any{TaskEllipsis, "VsCodePluginTask", "ms-python.python"},
	).Return(nil)

	task := &VsCodePluginTask{
		Attributes: *attrs,
		Code:       new(MockCode),
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
