package tasks

import (
	"context"
	"testing"

	i "awong/dotfiles/internal/tasks/installer"
	ty "awong/dotfiles/internal/types"

	"github.com/stretchr/testify/assert"
	mock "github.com/stretchr/testify/mock"
)

func TestNewBrewCaskTask_NilAttributes(t *testing.T) {
	task, err := NewBrewCaskTask(nil, "")
	assert.Error(t, err)
	assert.Nil(t, task)
	assert.Contains(t, err.Error(), "attributes cannot be nil")
}

func TestNewBrewCaskTask_EmptyID(t *testing.T) {
	attrs := &ty.Attributes{
		ID:          "",
		Description: "test",
	}
	task, err := NewBrewCaskTask(attrs, "")
	assert.Error(t, err)
	assert.Nil(t, task)
	assert.Contains(t, err.Error(), "attribute ID cannot be empty")
}

func TestNewBrewCaskTask_Success(t *testing.T) {
	attrs := &ty.Attributes{
		ID:          "test-cask",
		Description: "test cask task",
		Enabled:     true,
	}
	task, err := NewBrewCaskTask(attrs, "")
	assert.NoError(t, err)
	assert.NotNil(t, task)
	assert.Equal(t, "test-cask", task.ID)
}

func TestBrewCaskTask_Install_Disabled(t *testing.T) {
	attrs := &ty.Attributes{
		ID:          "test-cask",
		Description: "test cask task",
		Enabled:     false,
	}
	mockLog := i.NewMockLogger(t)
	mockLog.On("Debugf", "%s %s START ... play: %s taskID: %s", []any{"...........", "BrewCaskTask", "test-play", "test-cask"}).Return().Maybe()

	task := &BrewCaskTask{
		Attributes: *attrs,
		Brew:       new(MockBrew),
		Utils:      new(MockExt),
		Log:        mockLog,
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

func TestBrewCaskTask_Install_NotOSX(t *testing.T) {
	attrs := &ty.Attributes{
		ID:          "test-cask",
		Description: "test cask task",
		Enabled:     true,
	}
	mockExt := new(MockExt)
	mockExt.On("IsOSX").Return(false)

	mockLog := i.NewMockLogger(t)
	mockLog.On("Debugf", mock.Anything, mock.Anything, mock.Anything).Return().Maybe()

	task := &BrewCaskTask{
		Attributes: *attrs,
		Brew:       new(MockBrew),
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
	assert.Contains(t, err.Error(), "outside of macOS")
}

func TestBrewCaskTask_Install_AlreadyInstalled(t *testing.T) {
	attrs := &ty.Attributes{
		ID:          "test-cask",
		Description: "test cask task",
		Enabled:     true,
	}
	mockBrew := new(MockBrew)
	mockExt := new(MockExt)
	mockExt.On("IsOSX").Return(true)
	mockBrew.On("InPath", context.Background(), "Caskroom", "test-cask").Return(true)

	mockLog := i.NewMockLogger(t)
	mockLog.On("Debugf", mock.Anything, mock.Anything, mock.Anything).Return().Maybe()
	mockLog.On("Infof", mock.Anything, mock.Anything, mock.Anything).Return().Maybe()

	task := &BrewCaskTask{
		Attributes: *attrs,
		Brew:       mockBrew,
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
	assert.False(t, result.Success)
	mockBrew.AssertExpectations(t)
	mockExt.AssertExpectations(t)
}

func TestBrewCaskTask_Install_DryRun(t *testing.T) {
	attrs := &ty.Attributes{
		ID:          "test-cask",
		Description: "test cask task",
		Enabled:     true,
	}
	mockBrew := new(MockBrew)
	mockExt := new(MockExt)
	mockExt.On("IsOSX").Return(true)
	mockBrew.On("InPath", context.Background(), "Caskroom", "test-cask").Return(false)

	mockLog := i.NewMockLogger(t)
	mockLog.On("Debugf", mock.Anything, mock.Anything, mock.Anything).Return().Maybe()
	mockLog.On("Infof", mock.Anything, mock.Anything, mock.Anything).Return().Maybe()

	task := &BrewCaskTask{
		Attributes: *attrs,
		Brew:       mockBrew,
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
	assert.True(t, result.Success)
	mockBrew.AssertExpectations(t)
	mockExt.AssertExpectations(t)
}

func TestBrewCaskTask_Install_Success(t *testing.T) {
	attrs := &ty.Attributes{
		ID:          "test-cask",
		Description: "test cask task",
		Enabled:     true,
	}
	mockBrew := new(MockBrew)
	mockExt := new(MockExt)
	mockExt.On("IsOSX").Return(true)
	mockBrew.On("InPath", context.Background(), "Caskroom", "test-cask").Return(false)
	mockBrew.On("InstallCask", context.Background(), "test-cask").Return(nil)

	mockLog := i.NewMockLogger(t)
	mockLog.On("Debugf", mock.Anything, mock.Anything, mock.Anything).Return().Maybe()
	mockLog.On("Infof", mock.Anything, mock.Anything, mock.Anything).Return().Maybe()

	task := &BrewCaskTask{
		Attributes: *attrs,
		Brew:       mockBrew,
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
	mockBrew.AssertExpectations(t)
	mockExt.AssertExpectations(t)
}

func TestBrewCaskTask_Install_Failure(t *testing.T) {
	attrs := &ty.Attributes{
		ID:          "test-cask",
		Description: "test cask task",
		Enabled:     true,
	}
	mockBrew := new(MockBrew)
	mockExt := new(MockExt)
	mockExt.On("IsOSX").Return(true)
	mockBrew.On("InPath", context.Background(), "Caskroom", "test-cask").Return(false)
	mockBrew.On("InstallCask", context.Background(), "test-cask").Return(assert.AnError)

	mockLog := i.NewMockLogger(t)
	mockLog.On("Debugf", mock.Anything, mock.Anything, mock.Anything).Return().Maybe()
	mockLog.On("Infof", mock.Anything, mock.Anything, mock.Anything).Return().Maybe()

	task := &BrewCaskTask{
		Attributes: *attrs,
		Brew:       mockBrew,
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
	assert.Contains(t, err.Error(), "failed to brew install --cask")
	mockBrew.AssertExpectations(t)
	mockExt.AssertExpectations(t)
}

func TestBrewCaskTask_List(t *testing.T) {
	attrs := &ty.Attributes{
		ID:          "test-cask",
		Description: "test cask task",
		Enabled:     true,
	}
	mockLog := i.NewMockLogger(t)
	mockLog.On("Printlnf", mock.Anything, mock.Anything, mock.Anything).Return(nil).Maybe()

	task := &BrewCaskTask{
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
	mockLog.AssertExpectations(t)
}

func TestBrewCaskTask_GetAttributes(t *testing.T) {
	attrs := &ty.Attributes{
		ID:          "test-cask",
		Description: "test cask task",
		Enabled:     true,
	}
	task, err := NewBrewCaskTask(attrs, "")
	assert.NoError(t, err)

	result := task.GetAttributes()
	assert.Equal(t, "test-cask", result.ID)
}

func TestBrewCaskTask_GetID(t *testing.T) {
	attrs := &ty.Attributes{
		ID:          "test-cask",
		Description: "test cask task",
		Enabled:     true,
	}
	task, err := NewBrewCaskTask(attrs, "")
	assert.NoError(t, err)

	assert.Equal(t, "test-cask", task.GetID())
}

func TestBrewCaskTask_IsEnabled(t *testing.T) {
	attrs := &ty.Attributes{
		ID:          "test-cask",
		Description: "test cask task",
		Enabled:     true,
	}
	task, err := NewBrewCaskTask(attrs, "")
	assert.NoError(t, err)

	assert.True(t, task.IsEnabled())
}
