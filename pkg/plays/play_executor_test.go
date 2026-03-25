package plays

import (
	"context"
	"errors"
	"fmt"

	"awong/dotfiles/pkg/plays/gen"
	tk "awong/dotfiles/pkg/tasks"
	i "awong/dotfiles/pkg/tasks/installer"
	ty "awong/dotfiles/pkg/types"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func newPlayExecutorWithMocks(t *testing.T) (*PlayExecutor, *i.MockLogger, *i.MockExt) {
	logger := i.NewMockLogger(t)
	ext := &i.MockExt{}
	return &PlayExecutor{
		Log: logger,
		Ext: ext,
	}, logger, ext
}

func newPlayInput(pbi *gen.MockPlayBookInfo) *PlayInput {
	return &PlayInput{
		DryRun:   false,
		Play:     "test-play",
		PlayBook: pbi,
		Sudo:     false,
	}
}

func TestNewPlayExecutor_NoArgs(t *testing.T) {
	e := NewPlayExecutor()
	assert.NotNil(t, e)
	assert.NotNil(t, e.Log)
	assert.NotNil(t, e.Ext)
}

func TestPlayExecutor_Initialize_Success(t *testing.T) {
	e, _, ext := newPlayExecutorWithMocks(t)
	ext.On("IsInstalled", "tar").Return(true)
	ext.On("IsInstalled", "zip").Return(true)
	ext.On("IsInstalled", "unzip").Return(true)

	pbi := gen.NewMockPlayBookInfo(t)
	pbi.On("GetID").Return("test-book")
	input := newPlayInput(pbi)
	result, err := e.Initialize(input)
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.True(t, result.Success)
}

func TestPlayExecutor_Initialize_MissingTar(t *testing.T) {
	e, _, ext := newPlayExecutorWithMocks(t)
	ext.On("IsInstalled", "tar").Return(false)

	input := newPlayInput(gen.NewMockPlayBookInfo(t))
	result, err := e.Initialize(input)
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "tar is not installed")
}

func TestPlayExecutor_Initialize_MissingZip(t *testing.T) {
	e, _, ext := newPlayExecutorWithMocks(t)
	ext.On("IsInstalled", "tar").Return(true)
	ext.On("IsInstalled", "zip").Return(false)

	input := newPlayInput(gen.NewMockPlayBookInfo(t))
	result, err := e.Initialize(input)
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "zip is not installed")
}

func TestPlayExecutor_Initialize_MissingUnzip(t *testing.T) {
	e, _, ext := newPlayExecutorWithMocks(t)
	ext.On("IsInstalled", "tar").Return(true)
	ext.On("IsInstalled", "zip").Return(true)
	ext.On("IsInstalled", "unzip").Return(false)

	input := newPlayInput(gen.NewMockPlayBookInfo(t))
	result, err := e.Initialize(input)
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "unzip is not installed")
}

func TestPlayExecutor_InstallPlay_Success(t *testing.T) {
	e, logger, _ := newPlayExecutorWithMocks(t)
	ctx := context.Background()

	pbi := gen.NewMockPlayBookInfo(t)
	pbi.On("GetID").Return("test-book")
	pbi.On("GetJetBrainsApps").Return(map[string]string{})
	input := newPlayInput(pbi)

	logger.On("Printlnf", "%v Installing play (%v) ... %v", []any{PlayEllipsis, "test-play", ""}).Return(nil)

	taskResult := &ty.TaskResult{Task: "task-1", Play: "test-play", PlayBook: "test-book", Success: true}
	mockTask := &MockTask{
		Attr: ty.Attributes{ID: "task-1", Enabled: true},
	}
	mockTask.On("Install", ctx, mock.AnythingOfType("*types.TaskInput")).Return(taskResult, nil)

	playAttr := &ty.Attributes{ID: "test-play", Enabled: true}
	play, err := NewPlay(playAttr, []tk.Task{mockTask})
	assert.NoError(t, err)

	result, err := e.InstallPlay(ctx, play, input)
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.True(t, result.Success)
	assert.Equal(t, "test-play", result.Play)
	assert.Equal(t, "test-book", result.PlayBook)
}

func TestPlayExecutor_InstallPlay_PrintlnfError(t *testing.T) {
	e, logger, _ := newPlayExecutorWithMocks(t)
	ctx := context.Background()

	pbi := gen.NewMockPlayBookInfo(t)
	input := newPlayInput(pbi)

	logger.On("Printlnf", "%v Installing play (%v) ... %v", []any{PlayEllipsis, "test-play", ""}).Return(errors.New("log failed"))

	playAttr := &ty.Attributes{ID: "test-play", Enabled: true}
	play, err := NewPlay(playAttr, []tk.Task{})
	assert.NoError(t, err)

	result, err := e.InstallPlay(ctx, play, input)
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "failed to install play")
}

func TestPlayExecutor_InstallPlay_TaskError(t *testing.T) {
	e, logger, _ := newPlayExecutorWithMocks(t)
	ctx := context.Background()

	pbi := gen.NewMockPlayBookInfo(t)
	pbi.On("GetID").Return("test-book")
	pbi.On("GetJetBrainsApps").Return(map[string]string{})
	input := newPlayInput(pbi)

	logger.On("Printlnf", "%v Installing play (%v) ... %v", []any{PlayEllipsis, "test-play", ""}).Return(nil)

	mockTask := &MockTask{
		Attr: ty.Attributes{ID: "task-1", Enabled: true},
	}
	mockTask.On("Install", ctx, mock.AnythingOfType("*types.TaskInput")).Return((*ty.TaskResult)(nil), fmt.Errorf("install failed"))

	playAttr := &ty.Attributes{ID: "test-play", Enabled: true}
	play, err := NewPlay(playAttr, []tk.Task{mockTask})
	assert.NoError(t, err)

	result, err := e.InstallPlay(ctx, play, input)
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "failed to install task")
}

func TestPlayExecutor_ListPlay_Success(t *testing.T) {
	e, logger, _ := newPlayExecutorWithMocks(t)
	ctx := context.Background()

	pbi := gen.NewMockPlayBookInfo(t)
	pbi.On("GetID").Return("test-book")
	pbi.On("GetJetBrainsApps").Return(map[string]string{})
	input := newPlayInput(pbi)

	logger.On("Printlnf", "%v Listing play (%v) ... %v", []any{PlayEllipsis, "test-play", ""}).Return(nil)

	taskResult := &ty.TaskResult{Task: "task-1", Play: "test-play", PlayBook: "test-book", Success: true}
	mockTask := &MockTask{
		Attr: ty.Attributes{ID: "task-1", Enabled: true},
	}
	mockTask.On("List", ctx, mock.AnythingOfType("*types.TaskInput")).Return(taskResult, nil)

	playAttr := &ty.Attributes{ID: "test-play", Enabled: true}
	play, err := NewPlay(playAttr, []tk.Task{mockTask})
	assert.NoError(t, err)

	result, err := e.ListPlay(ctx, play, input)
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.True(t, result.Success)
}

func TestPlayExecutor_ListPlay_PrintlnfError(t *testing.T) {
	e, logger, _ := newPlayExecutorWithMocks(t)
	ctx := context.Background()

	pbi := gen.NewMockPlayBookInfo(t)
	input := newPlayInput(pbi)

	logger.On("Printlnf", "%v Listing play (%v) ... %v", []any{PlayEllipsis, "test-play", ""}).Return(errors.New("log failed"))

	playAttr := &ty.Attributes{ID: "test-play", Enabled: true}
	play, err := NewPlay(playAttr, []tk.Task{})
	assert.NoError(t, err)

	result, err := e.ListPlay(ctx, play, input)
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "failed to list play")
}

func TestPlayExecutor_ListPlay_TaskError(t *testing.T) {
	e, logger, _ := newPlayExecutorWithMocks(t)
	ctx := context.Background()

	pbi := gen.NewMockPlayBookInfo(t)
	pbi.On("GetID").Return("test-book")
	pbi.On("GetJetBrainsApps").Return(map[string]string{})
	input := newPlayInput(pbi)

	logger.On("Printlnf", "%v Listing play (%v) ... %v", []any{PlayEllipsis, "test-play", ""}).Return(nil)

	mockTask := &MockTask{
		Attr: ty.Attributes{ID: "task-1", Enabled: true},
	}
	mockTask.On("List", ctx, mock.AnythingOfType("*types.TaskInput")).Return((*ty.TaskResult)(nil), fmt.Errorf("list failed"))

	playAttr := &ty.Attributes{ID: "test-play", Enabled: true}
	play, err := NewPlay(playAttr, []tk.Task{mockTask})
	assert.NoError(t, err)

	result, err := e.ListPlay(ctx, play, input)
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "failed to list task")
}
