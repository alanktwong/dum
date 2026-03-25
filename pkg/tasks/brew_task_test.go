package tasks

import (
	"context"
	"testing"

	i "awong/dotfiles/pkg/tasks/installer"
	ty "awong/dotfiles/pkg/types"

	"github.com/stretchr/testify/assert"
	mock "github.com/stretchr/testify/mock"
)

func newTestLogger(t *testing.T) *i.MockLogger {
	l := i.NewMockLogger(t)
	l.On("Debugf", mock.Anything, mock.Anything, mock.Anything).Return().Maybe()
	l.On("Infof", mock.Anything, mock.Anything, mock.Anything).Return().Maybe()
	l.On("Printlnf", mock.Anything, mock.Anything, mock.Anything).Return(nil).Maybe()
	return l
}

func TestNewBrewTask_NilAttributes(t *testing.T) {
	task, err := NewBrewTask(nil, "")
	assert.Error(t, err)
	assert.Nil(t, task)
	assert.Contains(t, err.Error(), "attributes cannot be nil")
}

func TestNewBrewTask_EmptyID(t *testing.T) {
	attrs := &ty.Attributes{
		ID:          "",
		Description: "test",
	}
	task, err := NewBrewTask(attrs, "")
	assert.Error(t, err)
	assert.Nil(t, task)
	assert.Contains(t, err.Error(), "attribute ID cannot be empty")
}

func TestNewBrewTask_Success(t *testing.T) {
	attrs := &ty.Attributes{
		ID:          "test-brew",
		Description: "test brew task",
		Enabled:     true,
	}
	task, err := NewBrewTask(attrs, "")
	assert.NoError(t, err)
	assert.NotNil(t, task)
	assert.Equal(t, "test-brew", task.ID)
}

func TestNewBrewTask_WithTap(t *testing.T) {
	attrs := &ty.Attributes{
		ID:          "test-brew",
		Description: "test brew task",
		Enabled:     true,
	}
	task, err := NewBrewTask(attrs, "homebrew/cask")
	assert.NoError(t, err)
	assert.NotNil(t, task)
	assert.NotNil(t, task.Tap)
}

func TestBrewTask_Install_Disabled(t *testing.T) {
	attrs := &ty.Attributes{
		ID:          "test-brew",
		Description: "test brew task",
		Enabled:     false,
	}
	task := &BrewTask{
		Attributes: *attrs,
		Brew:       new(MockBrew),
		Utils:      new(MockExt),
		Log:        newTestLogger(t),
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

func TestBrewTask_Install_AlreadyInstalled(t *testing.T) {
	attrs := &ty.Attributes{
		ID:          "test-brew",
		Description: "test brew task",
		Enabled:     true,
	}
	mockBrew := new(MockBrew)
	mockExt := new(MockExt)
	mockBrew.On("InPath", context.Background(), "opt", "test-brew").Return(true)
	mockExt.On("IsInstalled", "test-brew").Return(false)

	task := &BrewTask{
		Attributes: *attrs,
		Brew:       mockBrew,
		Utils:      mockExt,
		Log:        newTestLogger(t),
	}

	input := &ty.TaskInput{
		Play:   "test-play",
		DryRun: false,
	}

	result, err := task.Install(context.Background(), input)
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.False(t, result.Success)
	mockBrew.AssertExpectations(t)
	mockExt.AssertExpectations(t)
}

func TestBrewTask_Install_DryRun(t *testing.T) {
	attrs := &ty.Attributes{
		ID:          "test-brew",
		Description: "test brew task",
		Enabled:     true,
	}
	mockBrew := new(MockBrew)
	mockExt := new(MockExt)
	mockBrew.On("InPath", context.Background(), "opt", "test-brew").Return(false)
	mockExt.On("IsInstalled", "test-brew").Return(false)

	task := &BrewTask{
		Attributes: *attrs,
		Brew:       mockBrew,
		Utils:      mockExt,
		Log:        newTestLogger(t),
	}

	input := &ty.TaskInput{
		Play:   "test-play",
		DryRun: true,
	}

	result, err := task.Install(context.Background(), input)
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.True(t, result.Success)
	mockBrew.AssertExpectations(t)
	mockExt.AssertExpectations(t)
}

func TestBrewTask_Install_Success(t *testing.T) {
	attrs := &ty.Attributes{
		ID:          "test-brew",
		Description: "test brew task",
		Enabled:     true,
	}
	mockBrew := new(MockBrew)
	mockExt := new(MockExt)
	mockBrew.On("InPath", context.Background(), "opt", "test-brew").Return(false)
	mockExt.On("IsInstalled", "test-brew").Return(false)
	mockBrew.On("Install", context.Background(), "test-brew").Return(nil)

	task := &BrewTask{
		Attributes: *attrs,
		Brew:       mockBrew,
		Utils:      mockExt,
		Log:        newTestLogger(t),
	}

	input := &ty.TaskInput{
		Play:   "test-play",
		DryRun: false,
	}

	result, err := task.Install(context.Background(), input)
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.True(t, result.Success)
	mockBrew.AssertExpectations(t)
	mockExt.AssertExpectations(t)
}

func TestBrewTask_Install_Failure(t *testing.T) {
	attrs := &ty.Attributes{
		ID:          "test-brew",
		Description: "test brew task",
		Enabled:     true,
	}
	mockBrew := new(MockBrew)
	mockExt := new(MockExt)
	mockBrew.On("InPath", context.Background(), "opt", "test-brew").Return(false)
	mockExt.On("IsInstalled", "test-brew").Return(false)
	mockBrew.On("Install", context.Background(), "test-brew").Return(assert.AnError)

	task := &BrewTask{
		Attributes: *attrs,
		Brew:       mockBrew,
		Utils:      mockExt,
		Log:        newTestLogger(t),
	}

	input := &ty.TaskInput{
		Play:   "test-play",
		DryRun: false,
	}

	result, err := task.Install(context.Background(), input)
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "failed to brew install")
	mockBrew.AssertExpectations(t)
	mockExt.AssertExpectations(t)
}

func TestBrewTask_List(t *testing.T) {
	attrs := &ty.Attributes{
		ID:          "test-brew",
		Description: "test brew task",
		Enabled:     true,
	}
	mockLog := newTestLogger(t)

	task := &BrewTask{
		Attributes: *attrs,
		Log:        mockLog,
	}

	input := &ty.TaskInput{
		Play: "test-play",
	}

	result, err := task.List(context.Background(), input)
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.True(t, result.Success)
}

func TestBrewTask_GetAttributes(t *testing.T) {
	attrs := &ty.Attributes{
		ID:          "test-brew",
		Description: "test brew task",
		Enabled:     true,
	}
	task, err := NewBrewTask(attrs, "")
	assert.NoError(t, err)

	result := task.GetAttributes()
	assert.Equal(t, "test-brew", result.ID)
}

func TestBrewTask_GetID(t *testing.T) {
	attrs := &ty.Attributes{
		ID:          "test-brew",
		Description: "test brew task",
		Enabled:     true,
	}
	task, err := NewBrewTask(attrs, "")
	assert.NoError(t, err)

	assert.Equal(t, "test-brew", task.GetID())
}

func TestBrewTask_IsEnabled(t *testing.T) {
	attrs := &ty.Attributes{
		ID:          "test-brew",
		Description: "test brew task",
		Enabled:     true,
	}
	task, err := NewBrewTask(attrs, "")
	assert.NoError(t, err)

	assert.True(t, task.IsEnabled())
}
