package tasks

import (
	"context"
	"testing"

	ty "alanktwong/dum/internal/types"

	"github.com/stretchr/testify/assert"
	mock "github.com/stretchr/testify/mock"
)

func TestNewLinkTask_NilAttributes(t *testing.T) {
	task, err := NewLinkTask(nil, "", "")
	assert.Error(t, err)
	assert.Nil(t, task)
	assert.Contains(t, err.Error(), "attributes cannot be nil")
}

func TestNewLinkTask_EmptyID(t *testing.T) {
	attrs := &ty.Attributes{
		ID:          "",
		Description: "test",
	}
	task, err := NewLinkTask(attrs, "", "")
	assert.Error(t, err)
	assert.Nil(t, task)
	assert.Contains(t, err.Error(), "task ID cannot be empty")
}

func TestNewLinkTask_Success(t *testing.T) {
	attrs := &ty.Attributes{
		ID:          "~/test-link",
		Description: "test link task",
		Enabled:     true,
	}
	task, err := NewLinkTask(attrs, "~/dotfiles", "")
	assert.NoError(t, err)
	assert.NotNil(t, task)
	assert.Equal(t, "~/test-link", task.ID)
	assert.Equal(t, "~/dotfiles", task.Root)
}

func TestNewLinkTask_WithTarget(t *testing.T) {
	attrs := &ty.Attributes{
		ID:          "~/test-link",
		Description: "test link task",
		Enabled:     true,
	}
	task, err := NewLinkTask(attrs, "~/dotfiles", "custom-target")
	assert.NoError(t, err)
	assert.NotNil(t, task)
	assert.Equal(t, "custom-target", task.Target)
}

func TestLinkTask_Install_EmptyRoot(t *testing.T) {
	attrs := &ty.Attributes{
		ID:          "~/test-link",
		Description: "test link task",
		Enabled:     true,
	}
	task := &LinkTask{
		Attributes: *attrs,
		Root:       "",
		Utils:      new(MockExt),
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

func TestLinkTask_Install_EmptyTarget(t *testing.T) {
	attrs := &ty.Attributes{
		ID:          "~/test-link",
		Description: "test link task",
		Enabled:     true,
	}
	mockExt := new(MockExt)
	mockExt.On("ExpandUser", "~/dotfiles").Return("/Users/testuser/dotfiles", nil)
	mockExt.On("IsSymlink", mock.Anything).Return(false)
	mockExt.On("SoftLink", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)

	task := &LinkTask{
		Attributes: *attrs,
		Root:       "~/dotfiles",
		Target:     "",
		Utils:      mockExt,
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
}

func TestLinkTask_Install_Disabled(t *testing.T) {
	attrs := &ty.Attributes{
		ID:          "~/test-link",
		Description: "test link task",
		Enabled:     false,
	}
	mockExt := new(MockExt)
	mockExt.On("ExpandUser", "~/dotfiles").Return("/Users/testuser/dotfiles", nil)

	task := &LinkTask{
		Attributes: *attrs,
		Root:       "~/dotfiles",
		Target:     "test-link",
		Utils:      mockExt,
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

func TestLinkTask_Install_DryRun(t *testing.T) {
	attrs := &ty.Attributes{
		ID:          "~/test-link",
		Description: "test link task",
		Enabled:     true,
	}
	mockExt := new(MockExt)
	mockExt.On("ExpandUser", "~/dotfiles").Return("/Users/testuser/dotfiles", nil)
	mockExt.On("IsSymlink", "/Users/testuser/dotfiles/test-link").Return(false)

	task := &LinkTask{
		Attributes: *attrs,
		Root:       "~/dotfiles",
		Target:     "test-link",
		Utils:      mockExt,
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
	mockExt.AssertExpectations(t)
}

func TestLinkTask_Install_Success(t *testing.T) {
	attrs := &ty.Attributes{
		ID:          "~/test-link",
		Description: "test link task",
		Enabled:     true,
	}
	mockExt := new(MockExt)
	mockExt.On("ExpandUser", "~/dotfiles").Return("/Users/testuser/dotfiles", nil)
	mockExt.On("IsSymlink", mock.Anything).Return(false)
	mockExt.On("SoftLink", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)

	task := &LinkTask{
		Attributes: *attrs,
		Root:       "~/dotfiles",
		Target:     "test-link",
		Utils:      mockExt,
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
	mockExt.AssertExpectations(t)
}

func TestLinkTask_Install_AlreadyLinked(t *testing.T) {
	attrs := &ty.Attributes{
		ID:          "~/test-link",
		Description: "test link task",
		Enabled:     true,
	}
	mockExt := new(MockExt)
	mockExt.On("ExpandUser", "~/dotfiles").Return("/Users/testuser/dotfiles", nil)
	mockExt.On("IsSymlink", mock.Anything).Return(true)
	mockExt.On("SoftLink", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)

	task := &LinkTask{
		Attributes: *attrs,
		Root:       "~/dotfiles",
		Target:     "test-link",
		Utils:      mockExt,
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
	mockExt.AssertExpectations(t)
}

func TestLinkTask_Install_Error(t *testing.T) {
	attrs := &ty.Attributes{
		ID:          "~/test-link",
		Description: "test link task",
		Enabled:     true,
	}
	mockExt := new(MockExt)
	mockExt.On("ExpandUser", "~/dotfiles").Return("/Users/testuser/dotfiles", nil)
	mockExt.On("IsSymlink", mock.Anything).Return(false)
	mockExt.On(
		"SoftLink",
		mock.Anything,
		mock.Anything,
		mock.Anything,
		mock.Anything,
		mock.Anything,
	).Return(
		assert.AnError,
	)

	task := &LinkTask{
		Attributes: *attrs,
		Root:       "~/dotfiles",
		Target:     "test-link",
		Utils:      mockExt,
		Log:        NewTestLogger(t),
	}

	input := &ty.TaskInput{
		Play:   "test-play",
		DryRun: false,
	}

	result, err := task.Install(context.Background(), input)
	assert.Error(t, err)
	assert.Nil(t, result)
	mockExt.AssertExpectations(t)
}

func TestLinkTask_ProvideTarget(t *testing.T) {
	task := &LinkTask{
		Target: "custom-target",
	}
	assert.Equal(t, "custom-target", task.ProvideTarget())

	task2 := &LinkTask{
		Attributes: ty.Attributes{
			ID: "~/some/path/file.txt",
		},
	}
	assert.Equal(t, "file.txt", task2.ProvideTarget())
}

func TestLinkTask_GetAttributes(t *testing.T) {
	attrs := &ty.Attributes{
		ID:          "~/test-link",
		Description: "test link task",
		Enabled:     true,
	}
	task, err := NewLinkTask(attrs, "~/dotfiles", "")
	assert.NoError(t, err)

	result := task.GetAttributes()
	assert.Equal(t, "~/test-link", result.ID)
}

func TestLinkTask_GetID(t *testing.T) {
	attrs := &ty.Attributes{
		ID:          "~/test-link",
		Description: "test link task",
		Enabled:     true,
	}
	task, err := NewLinkTask(attrs, "~/dotfiles", "")
	assert.NoError(t, err)

	assert.Equal(t, "~/test-link", task.GetID())
}

func TestLinkTask_IsEnabled(t *testing.T) {
	attrs := &ty.Attributes{
		ID:          "~/test-link",
		Description: "test link task",
		Enabled:     true,
	}
	task, err := NewLinkTask(attrs, "~/dotfiles", "")
	assert.NoError(t, err)

	assert.True(t, task.IsEnabled())
}

func TestLinkTask_Install_Disabled_DryRun(t *testing.T) {
	attrs := &ty.Attributes{
		ID:          "~/test-link",
		Description: "test link task",
		Enabled:     false,
	}
	mockExt := new(MockExt)

	mockLog := NewTestLogger(t)
	mockLog.On(
		"Printlnf",
		"%v %s: linking %v -> %s/%s",
		[]any{TaskEllipsis, "LinkTask", "~/test-link", "~/dotfiles", "test-link"},
	).Return(nil)

	task := &LinkTask{
		Attributes: *attrs,
		Root:       "~/dotfiles",
		Target:     "test-link",
		Utils:      mockExt,
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
