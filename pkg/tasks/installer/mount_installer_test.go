package installer

import (
	"context"
	"testing"

	ty "awong/dotfiles/pkg/types"

	"github.com/stretchr/testify/assert"
	mock "github.com/stretchr/testify/mock"
)

func TestMountInstaller_Install_NotOSX(t *testing.T) {
	mockExt := &MockExt{}
	mockExt.On("IsOSX").Return(false)

	mi := &MountInstaller{
		Utils: mockExt,
		Log:   &MockLogger{},
	}

	input := &ty.TaskInput{
		Play:   "test-play",
		DryRun: false,
	}

	result, err := mi.Install(context.Background(), input)
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "not on a Mac OSX")
	mockExt.AssertExpectations(t)
}

func TestMountInstaller_Install_DryRun(t *testing.T) {
	mockExt := &MockExt{}
	mockExt.On("IsOSX").Return(true)

	mockLog := NewMockLogger(t)
	mockLog.On("Infof", "%v %v: mounting a case-sensitive file system", []any{"...........", "test-play"}).Return()

	mi := &MountInstaller{
		Utils: mockExt,
		Log:   mockLog,
	}

	input := &ty.TaskInput{
		Play:   "test-play",
		DryRun: true,
	}

	result, err := mi.Install(context.Background(), input)
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.True(t, result.Success)
	mockExt.AssertExpectations(t)
	mockLog.AssertExpectations(t)
}

func TestMountInstaller_Install_Success(t *testing.T) {
	mockExt := &MockExt{}
	mockExt.On("IsOSX").Return(true)

	mockRunner := &mockMountRunner{}
	mockRunner.On("Run", context.Background()).Return(nil)

	mockLog := NewMockLogger(t)
	mockLog.On("Infof", "%v %v: mounting a case-sensitive file system", []any{"...........", "test-play"}).Return()

	mi := &MountInstaller{
		Utils:  mockExt,
		Runner: mockRunner,
		Log:    mockLog,
	}

	input := &ty.TaskInput{
		Play:   "test-play",
		DryRun: false,
	}

	result, err := mi.Install(context.Background(), input)
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.True(t, result.Success)
	mockExt.AssertExpectations(t)
	mockRunner.AssertExpectations(t)
	mockLog.AssertExpectations(t)
}

func TestMountInstaller_Install_RunError(t *testing.T) {
	mockExt := &MockExt{}
	mockExt.On("IsOSX").Return(true)

	mockRunner := &mockMountRunner{}
	mockRunner.On("Run", context.Background()).Return(assert.AnError)

	mockLog := NewMockLogger(t)
	mockLog.On("Infof", "%v %v: mounting a case-sensitive file system", []any{"...........", "test-play"}).Return()

	mi := &MountInstaller{
		Utils:  mockExt,
		Runner: mockRunner,
		Log:    mockLog,
	}

	input := &ty.TaskInput{
		Play:   "test-play",
		DryRun: false,
	}

	result, err := mi.Install(context.Background(), input)
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "case-sensitive file system")
	mockExt.AssertExpectations(t)
	mockRunner.AssertExpectations(t)
}

type mockMountRunner struct {
	mock.Mock
}

func (m *mockMountRunner) Run(ctx context.Context) error {
	ret := m.Called(ctx)
	return ret.Error(0)
}
