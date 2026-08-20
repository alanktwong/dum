package tasks

import (
	"context"
	"testing"

	i "github.com/alanktwong/dum/internal/tasks/installer"
	ty "github.com/alanktwong/dum/internal/types"

	"github.com/stretchr/testify/assert"
	mock "github.com/stretchr/testify/mock"
)

func newDirTestLogger(t *testing.T) *i.MockLogger {
	l := i.NewMockLogger(t)
	l.On("Debugf", mock.Anything, mock.Anything, mock.Anything).Return().Maybe()
	l.On("Infof", mock.Anything, mock.Anything, mock.Anything).Return().Maybe()
	l.On("Printlnf", mock.Anything, mock.Anything, mock.Anything).Return(nil).Maybe()
	return l
}

func TestNewDirTask_NilAttributes(t *testing.T) {
	task, err := NewDirTask(nil)
	assert.Error(t, err)
	assert.Nil(t, task)
	assert.Contains(t, err.Error(), "attributes cannot be nil")
}

func TestNewDirTask_EmptyID(t *testing.T) {
	attrs := &ty.Attributes{
		ID:          "",
		Description: "test",
	}
	task, err := NewDirTask(attrs)
	assert.Error(t, err)
	assert.Nil(t, task)
	assert.Contains(t, err.Error(), "task ID cannot be empty")
}

func TestNewDirTask_Success(t *testing.T) {
	attrs := &ty.Attributes{
		ID:          "~/test-dir",
		Description: "test dir task",
		Enabled:     true,
	}
	task, err := NewDirTask(attrs)
	assert.NoError(t, err)
	assert.NotNil(t, task)
	assert.Equal(t, "~/test-dir", task.ID)
}

func TestDirTask_Install_Disabled(t *testing.T) {
	attrs := &ty.Attributes{
		ID:          "~/test-dir",
		Description: "test dir task",
		Enabled:     false,
	}
	task := &DirTask{
		Attributes: *attrs,
		Utils:      new(MockExt),
		Log:        newDirTestLogger(t),
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

func TestDirTask_Install_AlreadyExists(t *testing.T) {
	attrs := &ty.Attributes{
		ID:          "~/test-dir",
		Description: "test dir task",
		Enabled:     true,
	}
	mockExt := new(MockExt)
	mockExt.On("ExpandUser", "~/test-dir").Return("/Users/testuser/test-dir", nil)
	mockExt.On("IsDir", "/Users/testuser/test-dir").Return(true)

	task := &DirTask{
		Attributes: *attrs,
		Utils:      mockExt,
		Log:        newDirTestLogger(t),
	}

	input := &ty.TaskInput{
		Play:   "test-play",
		DryRun: false,
	}

	result, err := task.Install(context.Background(), input)
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.False(t, result.Success)
	mockExt.AssertExpectations(t)
}

func TestDirTask_Install_DryRun(t *testing.T) {
	attrs := &ty.Attributes{
		ID:          "~/test-dir",
		Description: "test dir task",
		Enabled:     true,
	}
	mockExt := new(MockExt)
	mockExt.On("ExpandUser", "~/test-dir").Return("/Users/testuser/test-dir", nil)
	mockExt.On("IsDir", "/Users/testuser/test-dir").Return(false)

	task := &DirTask{
		Attributes: *attrs,
		Utils:      mockExt,
		Log:        newDirTestLogger(t),
	}

	input := &ty.TaskInput{
		Play:   "test-play",
		DryRun: true,
	}

	result, err := task.Install(context.Background(), input)
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.True(t, result.Success)
	mockExt.AssertExpectations(t)
}

func TestDirTask_Install_Success(t *testing.T) {
	attrs := &ty.Attributes{
		ID:          "~/test-dir",
		Description: "test dir task",
		Enabled:     true,
	}
	mockExt := new(MockExt)
	mockExt.On("ExpandUser", "~/test-dir").Return("/Users/testuser/test-dir", nil)
	mockExt.On("IsDir", "/Users/testuser/test-dir").Return(false)
	mockExt.On("CreateDirectory", context.Background(), "/Users/testuser/test-dir", false).Return(nil)

	task := &DirTask{
		Attributes: *attrs,
		Utils:      mockExt,
		Log:        newDirTestLogger(t),
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
}

func TestDirTask_Install_ExpandUserError(t *testing.T) {
	attrs := &ty.Attributes{
		ID:          "~/test-dir",
		Description: "test dir task",
		Enabled:     true,
	}
	mockExt := new(MockExt)
	mockExt.On("ExpandUser", "~/test-dir").Return("", assert.AnError)

	task := &DirTask{
		Attributes: *attrs,
		Utils:      mockExt,
		Log:        newDirTestLogger(t),
	}

	input := &ty.TaskInput{
		Play:   "test-play",
		DryRun: false,
	}

	result, err := task.Install(context.Background(), input)
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "failed to expand user")
	mockExt.AssertExpectations(t)
}

func TestDirTask_Install_CreateDirectoryError(t *testing.T) {
	attrs := &ty.Attributes{
		ID:          "~/test-dir",
		Description: "test dir task",
		Enabled:     true,
	}
	mockExt := new(MockExt)
	mockExt.On("ExpandUser", "~/test-dir").Return("/Users/testuser/test-dir", nil)
	mockExt.On("IsDir", "/Users/testuser/test-dir").Return(false)
	mockExt.On("CreateDirectory", context.Background(), "/Users/testuser/test-dir", false).Return(assert.AnError)

	task := &DirTask{
		Attributes: *attrs,
		Utils:      mockExt,
		Log:        newDirTestLogger(t),
	}

	input := &ty.TaskInput{
		Play:   "test-play",
		DryRun: false,
	}

	result, err := task.Install(context.Background(), input)
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "failed to mkdir -p")
	mockExt.AssertExpectations(t)
}

func TestDirTask_GetAttributes(t *testing.T) {
	attrs := &ty.Attributes{
		ID:          "~/test-dir",
		Description: "test dir task",
		Enabled:     true,
	}
	task, err := NewDirTask(attrs)
	assert.NoError(t, err)

	result := task.GetAttributes()
	assert.Equal(t, "~/test-dir", result.ID)
}

func TestDirTask_GetID(t *testing.T) {
	attrs := &ty.Attributes{
		ID:          "~/test-dir",
		Description: "test dir task",
		Enabled:     true,
	}
	task, err := NewDirTask(attrs)
	assert.NoError(t, err)

	assert.Equal(t, "~/test-dir", task.GetID())
}

func TestDirTask_IsEnabled(t *testing.T) {
	attrs := &ty.Attributes{
		ID:          "~/test-dir",
		Description: "test dir task",
		Enabled:     true,
	}
	task, err := NewDirTask(attrs)
	assert.NoError(t, err)

	assert.True(t, task.IsEnabled())
}

func TestDirTask_Install_Disabled_DryRun(t *testing.T) {
	attrs := &ty.Attributes{
		ID:          "~/test-dir",
		Description: "test dir task",
		Enabled:     false,
	}
	mockLog := i.NewMockLogger(t)
	mockLog.On(
		"Debugf",
		"%s %s START ... play: %s taskID: %s",
		[]any{TaskEllipsis, "DirTask", "test-play", "~/test-dir"},
	).Return().Maybe()
	mockLog.On(
		"Printlnf",
		"%v %s: mkdir -p %s",
		[]any{TaskEllipsis, "DirTask", "~/test-dir"},
	).Return(nil)

	task := &DirTask{
		Attributes: *attrs,
		Utils:      new(MockExt),
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
