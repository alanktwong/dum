package tasks

import (
	"context"
	"testing"

	ty "alanktwong/dum/internal/types"

	"github.com/stretchr/testify/assert"
)

func TestNewJetBrainsPluginTask_NilAttributes(t *testing.T) {
	task, err := NewJetBrainsPluginTask(nil, []string{"idea"})
	assert.Error(t, err)
	assert.Nil(t, task)
	assert.Contains(t, err.Error(), "attributes cannot be nil")
}

func TestNewJetBrainsPluginTask_EmptyID(t *testing.T) {
	attrs := &ty.Attributes{
		ID:          "",
		Description: "test",
	}
	task, err := NewJetBrainsPluginTask(attrs, []string{"idea"})
	assert.Error(t, err)
	assert.Nil(t, task)
	assert.Contains(t, err.Error(), "task ID cannot be empty")
}

func TestNewJetBrainsPluginTask_EmptyApps(t *testing.T) {
	attrs := &ty.Attributes{
		ID:          "test-plugin",
		Description: "test",
	}
	task, err := NewJetBrainsPluginTask(attrs, []string{})
	assert.Error(t, err)
	assert.Nil(t, task)
	assert.Contains(t, err.Error(), "app cannot be empty")
}

func TestNewJetBrainsPluginTask_Success(t *testing.T) {
	attrs := &ty.Attributes{
		ID:          "test-plugin",
		Description: "test jetbrains task",
		Enabled:     true,
	}
	task, err := NewJetBrainsPluginTask(attrs, []string{"idea", "pycharm"})
	assert.NoError(t, err)
	assert.NotNil(t, task)
	assert.Equal(t, "test-plugin", task.ID)
	assert.Len(t, task.Apps, 2)
}

func TestJetBrainsPluginTask_Install_Disabled(t *testing.T) {
	attrs := &ty.Attributes{
		ID:          "test-plugin",
		Description: "test jetbrains task",
		Enabled:     false,
	}
	task := &JetBrainsPluginTask{
		Attributes: *attrs,
		Apps:       []string{"idea"},
		Utils:      new(MockExt),
		JetBrains:  new(MockJetBrainsApp),
		Log:        NewTestLogger(t),
	}

	input := &ty.TaskInput{
		Play:          "test-play",
		DryRun:        false,
		JetBrainsApps: map[string]string{},
	}

	result, err := task.Install(context.Background(), input)
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.False(t, result.Success)
}

func TestJetBrainsPluginTask_Install_NoAppsInInput(t *testing.T) {
	attrs := &ty.Attributes{
		ID:          "test-plugin",
		Description: "test jetbrains task",
		Enabled:     true,
	}
	mockExt := new(MockExt)
	mockExt.On("IsInstalled", "idea").Return(true)

	task := &JetBrainsPluginTask{
		Attributes: *attrs,
		Apps:       []string{"idea"},
		Utils:      mockExt,
		JetBrains:  new(MockJetBrainsApp),
		Log:        NewTestLogger(t),
	}

	input := &ty.TaskInput{
		Play:          "test-play",
		DryRun:        false,
		JetBrainsApps: map[string]string{},
	}

	result, err := task.Install(context.Background(), input)
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.False(t, result.Success)
	mockExt.AssertExpectations(t)
}

func TestJetBrainsPluginTask_Install_AlreadyInstalled(t *testing.T) {
	attrs := &ty.Attributes{
		ID:          "test-plugin",
		Description: "test jetbrains task",
		Enabled:     true,
	}
	mockExt := new(MockExt)
	mockExt.On("IsInstalled", "idea").Return(true)
	mockJetBrains := new(MockJetBrainsApp)
	mockJetBrains.On("IsInstalled", "ideaIU", "test-plugin").Return(true)

	task := &JetBrainsPluginTask{
		Attributes: *attrs,
		Apps:       []string{"idea"},
		Utils:      mockExt,
		JetBrains:  mockJetBrains,
		Log:        NewTestLogger(t),
	}

	input := &ty.TaskInput{
		Play:          "test-play",
		DryRun:        false,
		JetBrainsApps: map[string]string{"idea": "ideaIU"},
	}

	result, err := task.Install(context.Background(), input)
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.True(t, result.Success)
	mockExt.AssertExpectations(t)
	mockJetBrains.AssertExpectations(t)
}

func TestJetBrainsPluginTask_Install_DryRun(t *testing.T) {
	attrs := &ty.Attributes{
		ID:          "test-plugin",
		Description: "test jetbrains task",
		Enabled:     true,
	}
	mockExt := new(MockExt)
	mockExt.On("IsInstalled", "idea").Return(true)
	mockJetBrains := new(MockJetBrainsApp)
	mockJetBrains.On("IsInstalled", "ideaIU", "test-plugin").Return(false)

	task := &JetBrainsPluginTask{
		Attributes: *attrs,
		Apps:       []string{"idea"},
		Utils:      mockExt,
		JetBrains:  mockJetBrains,
		Log:        NewTestLogger(t),
	}

	input := &ty.TaskInput{
		Play:          "test-play",
		DryRun:        true,
		JetBrainsApps: map[string]string{"idea": "ideaIU"},
	}

	result, err := task.Install(context.Background(), input)
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.True(t, result.Success)
	mockExt.AssertExpectations(t)
	mockJetBrains.AssertExpectations(t)
}

func TestJetBrainsPluginTask_Install_Success(t *testing.T) {
	attrs := &ty.Attributes{
		ID:          "test-plugin",
		Description: "test jetbrains task",
		Enabled:     true,
	}
	mockExt := new(MockExt)
	mockExt.On("IsInstalled", "idea").Return(true)
	mockJetBrains := new(MockJetBrainsApp)
	mockJetBrains.On("IsInstalled", "ideaIU", "test-plugin").Return(false)
	mockJetBrains.On("Install", context.Background(), "idea", "test-plugin").Return(nil)

	task := &JetBrainsPluginTask{
		Attributes: *attrs,
		Apps:       []string{"idea"},
		Utils:      mockExt,
		JetBrains:  mockJetBrains,
		Log:        NewTestLogger(t),
	}

	input := &ty.TaskInput{
		Play:          "test-play",
		DryRun:        false,
		JetBrainsApps: map[string]string{"idea": "ideaIU"},
	}

	result, err := task.Install(context.Background(), input)
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.True(t, result.Success)
	mockExt.AssertExpectations(t)
	mockJetBrains.AssertExpectations(t)
}

func TestJetBrainsPluginTask_Install_Error(t *testing.T) {
	attrs := &ty.Attributes{
		ID:          "test-plugin",
		Description: "test jetbrains task",
		Enabled:     true,
	}
	mockExt := new(MockExt)
	mockExt.On("IsInstalled", "idea").Return(true)
	mockJetBrains := new(MockJetBrainsApp)
	mockJetBrains.On("IsInstalled", "ideaIU", "test-plugin").Return(false)
	mockJetBrains.On("Install", context.Background(), "idea", "test-plugin").Return(assert.AnError)

	task := &JetBrainsPluginTask{
		Attributes: *attrs,
		Apps:       []string{"idea"},
		Utils:      mockExt,
		JetBrains:  mockJetBrains,
		Log:        NewTestLogger(t),
	}

	input := &ty.TaskInput{
		Play:          "test-play",
		DryRun:        false,
		JetBrainsApps: map[string]string{"idea": "ideaIU"},
	}

	result, err := task.Install(context.Background(), input)
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.False(t, result.Success)
	mockExt.AssertExpectations(t)
	mockJetBrains.AssertExpectations(t)
}

func TestJetBrainsPluginTask_activeApps(t *testing.T) {
	mockExt := new(MockExt)
	mockExt.On("IsInstalled", "idea").Return(true)
	mockExt.On("IsInstalled", "pycharm").Return(false)

	task := &JetBrainsPluginTask{
		Apps:  []string{"idea", "pycharm"},
		Utils: mockExt,
	}

	activeApps := task.activeApps()
	assert.Len(t, activeApps, 1)
	assert.Equal(t, "idea", activeApps[0])
	mockExt.AssertExpectations(t)
}

func TestJetBrainsPluginTask_GetAttributes(t *testing.T) {
	attrs := &ty.Attributes{
		ID:          "test-plugin",
		Description: "test jetbrains task",
		Enabled:     true,
	}
	task, err := NewJetBrainsPluginTask(attrs, []string{"idea"})
	assert.NoError(t, err)

	result := task.GetAttributes()
	assert.Equal(t, "test-plugin", result.ID)
}

func TestJetBrainsPluginTask_GetID(t *testing.T) {
	attrs := &ty.Attributes{
		ID:          "test-plugin",
		Description: "test jetbrains task",
		Enabled:     true,
	}
	task, err := NewJetBrainsPluginTask(attrs, []string{"idea"})
	assert.NoError(t, err)

	assert.Equal(t, "test-plugin", task.GetID())
}

func TestJetBrainsPluginTask_IsEnabled(t *testing.T) {
	attrs := &ty.Attributes{
		ID:          "test-plugin",
		Description: "test jetbrains task",
		Enabled:     true,
	}
	task, err := NewJetBrainsPluginTask(attrs, []string{"idea"})
	assert.NoError(t, err)

	assert.True(t, task.IsEnabled())
}

func TestJetBrainsPluginTask_Install_Disabled_DryRun(t *testing.T) {
	attrs := &ty.Attributes{
		ID:          "test-plugin",
		Description: "test jetbrains task",
		Enabled:     false,
	}
	mockExt := new(MockExt)
	mockExt.On("IsInstalled", "idea").Return(true)

	mockLog := NewTestLogger(t)
	mockLog.On(
		"Printlnf",
		"%v %s: %s installPlugins %s",
		[]any{TaskEllipsis, "JetBrainsPluginTask", "idea", "test-plugin"},
	).Return(nil)

	task := &JetBrainsPluginTask{
		Attributes: *attrs,
		Apps:       []string{"idea"},
		Utils:      mockExt,
		JetBrains:  new(MockJetBrainsApp),
		Log:        mockLog,
	}

	input := &ty.TaskInput{
		Play:          "test-play",
		DryRun:        true,
		JetBrainsApps: map[string]string{},
	}

	result, err := task.Install(context.Background(), input)
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.False(t, result.Success)
	mockExt.AssertExpectations(t)
	mockLog.AssertExpectations(t)
}
