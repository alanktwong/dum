package tasks

import (
	"context"
	"testing"

	i "alanktwong/dum/internal/tasks/installer"
	ty "alanktwong/dum/internal/types"

	"github.com/stretchr/testify/assert"
	mock "github.com/stretchr/testify/mock"
)

func TestNewBrewTap_NilAttributes(t *testing.T) {
	task, err := NewBrewTap(nil)
	assert.Error(t, err)
	assert.Nil(t, task)
	assert.Contains(t, err.Error(), "attributes cannot be nil")
}

func TestNewBrewTap_EmptyID(t *testing.T) {
	attrs := &ty.Attributes{
		ID:          "",
		Description: "test",
	}
	task, err := NewBrewTap(attrs)
	assert.Error(t, err)
	assert.Nil(t, task)
	assert.Contains(t, err.Error(), "tap ID cannot be empty")
}

func TestNewBrewTap_Success(t *testing.T) {
	attrs := &ty.Attributes{
		ID:          "homebrew/cask",
		Description: "test tap",
		Enabled:     true,
	}
	task, err := NewBrewTap(attrs)
	assert.NoError(t, err)
	assert.NotNil(t, task)
	assert.Equal(t, "homebrew/cask", task.ID)
}

func TestBrewTap_Install_Disabled(t *testing.T) {
	attrs := &ty.Attributes{
		ID:          "homebrew/cask",
		Description: "test tap",
		Enabled:     false,
	}
	task := &BrewTap{
		Attributes: *attrs,
		Brew:       new(MockBrew),
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

func TestBrewTap_Install_DryRun(t *testing.T) {
	attrs := &ty.Attributes{
		ID:          "homebrew/cask",
		Description: "test tap",
		Enabled:     true,
	}
	mockBrew := new(MockBrew)

	mockLog := i.NewMockLogger(t)
	mockLog.On("Infof", mock.Anything, mock.Anything, mock.Anything).Return().Maybe()

	task := &BrewTap{
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
}

func TestBrewTap_Install_Success(t *testing.T) {
	attrs := &ty.Attributes{
		ID:          "homebrew/cask",
		Description: "test tap",
		Enabled:     true,
	}
	mockBrew := new(MockBrew)
	mockBrew.On("Tap", context.Background(), "homebrew/cask").Return(nil)

	mockLog := i.NewMockLogger(t)
	mockLog.On("Infof", mock.Anything, mock.Anything, mock.Anything).Return().Maybe()

	task := &BrewTap{
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

func TestBrewTap_Install_Failure(t *testing.T) {
	attrs := &ty.Attributes{
		ID:          "homebrew/cask",
		Description: "test tap",
		Enabled:     true,
	}
	mockBrew := new(MockBrew)
	mockBrew.On("Tap", context.Background(), "homebrew/cask").Return(assert.AnError)

	mockLog := i.NewMockLogger(t)
	mockLog.On("Infof", mock.Anything, mock.Anything, mock.Anything).Return().Maybe()

	task := &BrewTap{
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
	assert.Contains(t, err.Error(), "failed to tap")
	mockBrew.AssertExpectations(t)
}

func TestBrewTap_IsEnabled(t *testing.T) {
	attrs := &ty.Attributes{
		ID:          "homebrew/cask",
		Description: "test tap",
		Enabled:     true,
	}
	task, err := NewBrewTap(attrs)
	assert.NoError(t, err)

	assert.True(t, task.IsEnabled())
}
