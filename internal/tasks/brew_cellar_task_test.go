package tasks

import (
	"context"
	"testing"

	i "awong/dotfiles/internal/tasks/installer"
	ty "awong/dotfiles/internal/types"

	"github.com/stretchr/testify/assert"
	mock "github.com/stretchr/testify/mock"
)

func TestNewBrewCellarTask_NilAttributes(t *testing.T) {
	task, err := NewBrewCellarTask(nil, "")
	assert.Error(t, err)
	assert.Nil(t, task)
	assert.Contains(t, err.Error(), "attributes cannot be nil")
}

func TestNewBrewCellarTask_EmptyID(t *testing.T) {
	attrs := &ty.Attributes{
		ID:          "",
		Description: "test",
	}
	task, err := NewBrewCellarTask(attrs, "")
	assert.Error(t, err)
	assert.Nil(t, task)
	assert.Contains(t, err.Error(), "attribute ID cannot be empty")
}

func TestNewBrewCellarTask_Success(t *testing.T) {
	attrs := &ty.Attributes{
		ID:          "test-cellar",
		Description: "test cellar task",
		Enabled:     true,
	}
	task, err := NewBrewCellarTask(attrs, "")
	assert.NoError(t, err)
	assert.NotNil(t, task)
	assert.Equal(t, "test-cellar", task.ID)
}

func TestBrewCellarTask_Install_Disabled(t *testing.T) {
	attrs := &ty.Attributes{
		ID:          "test-cellar",
		Description: "test cellar task",
		Enabled:     false,
	}
	mockLog := i.NewMockLogger(t)
	mockLog.On("Infof", mock.Anything, mock.Anything, mock.Anything).Return().Maybe()

	task := &BrewCellarTask{
		Attributes: *attrs,
		Brew:       new(MockBrew),
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

func TestBrewCellarTask_Install_AlreadyInstalled(t *testing.T) {
	attrs := &ty.Attributes{
		ID:          "test-cellar",
		Description: "test cellar task",
		Enabled:     true,
	}
	mockBrew := new(MockBrew)
	mockBrew.On("InPath", context.Background(), "Cellar", "test-cellar").Return(true)

	mockLog := i.NewMockLogger(t)
	mockLog.On("Debugf", mock.Anything, mock.Anything, mock.Anything).Return().Maybe()
	mockLog.On("Infof", mock.Anything, mock.Anything, mock.Anything).Return().Maybe()
	mockLog.On("Printlnf", mock.Anything, mock.Anything, mock.Anything).Return(nil).Maybe()

	task := &BrewCellarTask{
		Attributes: *attrs,
		Brew:       mockBrew,
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
}

func TestBrewCellarTask_Install_DryRun(t *testing.T) {
	attrs := &ty.Attributes{
		ID:          "test-cellar",
		Description: "test cellar task",
		Enabled:     true,
	}
	mockBrew := new(MockBrew)
	mockBrew.On("InPath", context.Background(), "Cellar", "test-cellar").Return(false)

	mockLog := i.NewMockLogger(t)
	mockLog.On("Debugf", mock.Anything, mock.Anything, mock.Anything).Return().Maybe()
	mockLog.On("Infof", mock.Anything, mock.Anything, mock.Anything).Return().Maybe()
	mockLog.On("Printlnf", mock.Anything, mock.Anything, mock.Anything).Return(nil).Maybe()

	task := &BrewCellarTask{
		Attributes: *attrs,
		Brew:       mockBrew,
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
}

func TestBrewCellarTask_Install_Success(t *testing.T) {
	attrs := &ty.Attributes{
		ID:          "test-cellar",
		Description: "test cellar task",
		Enabled:     true,
	}
	mockBrew := new(MockBrew)
	mockBrew.On("InPath", context.Background(), "Cellar", "test-cellar").Return(false)
	mockBrew.On("Install", context.Background(), "test-cellar").Return(nil)

	mockLog := i.NewMockLogger(t)
	mockLog.On("Debugf", mock.Anything, mock.Anything, mock.Anything).Return().Maybe()
	mockLog.On("Infof", mock.Anything, mock.Anything, mock.Anything).Return().Maybe()

	task := &BrewCellarTask{
		Attributes: *attrs,
		Brew:       mockBrew,
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
}

func TestBrewCellarTask_Install_Failure(t *testing.T) {
	attrs := &ty.Attributes{
		ID:          "test-cellar",
		Description: "test cellar task",
		Enabled:     true,
	}
	mockBrew := new(MockBrew)
	mockBrew.On("InPath", context.Background(), "Cellar", "test-cellar").Return(false)
	mockBrew.On("Install", context.Background(), "test-cellar").Return(assert.AnError)

	mockLog := i.NewMockLogger(t)
	mockLog.On("Debugf", mock.Anything, mock.Anything, mock.Anything).Return().Maybe()
	mockLog.On("Infof", mock.Anything, mock.Anything, mock.Anything).Return().Maybe()

	task := &BrewCellarTask{
		Attributes: *attrs,
		Brew:       mockBrew,
		Log:        mockLog,
	}

	input := &ty.TaskInput{
		Play:   "test-play",
		DryRun: false,
	}

	result, err := task.Install(context.Background(), input)
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "failed to brew install (cellar)")
	mockBrew.AssertExpectations(t)
}

func TestBrewCellarTask_List(t *testing.T) {
	attrs := &ty.Attributes{
		ID:          "test-cellar",
		Description: "test cellar task",
		Enabled:     true,
	}
	mockLog := i.NewMockLogger(t)
	mockLog.On("Printlnf", mock.Anything, mock.Anything, mock.Anything).Return(nil).Maybe()

	task := &BrewCellarTask{
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

func TestBrewCellarTask_GetAttributes(t *testing.T) {
	attrs := &ty.Attributes{
		ID:          "test-cellar",
		Description: "test cellar task",
		Enabled:     true,
	}
	task, err := NewBrewCellarTask(attrs, "")
	assert.NoError(t, err)

	result := task.GetAttributes()
	assert.Equal(t, "test-cellar", result.ID)
}

func TestBrewCellarTask_GetID(t *testing.T) {
	attrs := &ty.Attributes{
		ID:          "test-cellar",
		Description: "test cellar task",
		Enabled:     true,
	}
	task, err := NewBrewCellarTask(attrs, "")
	assert.NoError(t, err)

	assert.Equal(t, "test-cellar", task.GetID())
}

func TestBrewCellarTask_IsEnabled(t *testing.T) {
	attrs := &ty.Attributes{
		ID:          "test-cellar",
		Description: "test cellar task",
		Enabled:     true,
	}
	task, err := NewBrewCellarTask(attrs, "")
	assert.NoError(t, err)

	assert.True(t, task.IsEnabled())
}
