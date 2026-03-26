package installer

import (
	"context"
	"testing"

	ty "awong/dotfiles/internal/types"

	"github.com/stretchr/testify/assert"
)

func TestBashInstaller_Install_NonOSX(t *testing.T) {
	mockBrew := &MockBrew{}
	mockExt := &MockExt{}
	mockLog := &MockLogger{}

	bi := &BashInstaller{
		Brew:  mockBrew,
		Utils: mockExt,
		Log:   mockLog,
	}

	mockExt.On("IsOSX").Return(false).Maybe()
	mockLog.On("Infof", "%v %v: there is no need to reinstall bash outside of Mac OSX", []any{"...........", "test-play"}).Return().Maybe()

	input := &ty.TaskInput{
		Play:   "test-play",
		Sudo:   false,
		DryRun: false,
	}

	result, err := bi.Install(context.Background(), input)
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.False(t, result.Success)
}

func TestBashInstaller_Install_AlreadyInstalled(t *testing.T) {
	mockBrew := &MockBrew{}
	mockExt := &MockExt{}
	mockLog := &MockLogger{}

	bi := &BashInstaller{
		Brew:  mockBrew,
		Utils: mockExt,
		Log:   mockLog,
	}

	mockExt.On("IsOSX").Return(true).Maybe()
	mockExt.On("IsInstalled", "bash").Return(true).Maybe()
	mockLog.On("Infof", "%v %v: updated bash is already installed", []any{"...........", "test-play"}).Return().Maybe()

	input := &ty.TaskInput{
		Play:   "test-play",
		Sudo:   false,
		DryRun: false,
	}

	result, err := bi.Install(context.Background(), input)
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.True(t, result.Success)
}

func TestBashInstaller_Install_Success(t *testing.T) {
	mockBrew := &MockBrew{}
	mockExt := &MockExt{}
	mockLog := &MockLogger{}

	bi := &BashInstaller{
		Brew:  mockBrew,
		Utils: mockExt,
		Log:   mockLog,
	}

	mockExt.On("IsOSX").Return(true).Maybe()
	mockExt.On("IsInstalled", "bash").Return(false).Maybe()
	mockLog.On("Infof", "%v %v: install bash", []any{"...........", "test-play"}).Return().Maybe()
	mockBrew.On("Install", context.Background(), "bash").Return(nil)
	mockBrew.On("Prefix", context.Background()).Return("/usr/local", nil)
	mockExt.On("SoftLink", context.Background(), "/usr/local/bin", "/usr/local/bin/bash", "bash", false).Return(nil)

	input := &ty.TaskInput{
		Play:   "test-play",
		Sudo:   false,
		DryRun: false,
	}

	result, err := bi.Install(context.Background(), input)
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.True(t, result.Success)
}

func TestBashInstaller_Install_DryRun(t *testing.T) {
	mockBrew := &MockBrew{}
	mockExt := &MockExt{}
	mockLog := &MockLogger{}

	bi := &BashInstaller{
		Brew:  mockBrew,
		Utils: mockExt,
		Log:   mockLog,
	}

	mockExt.On("IsOSX").Return(true).Maybe()
	mockExt.On("IsInstalled", "bash").Return(false).Maybe()
	mockLog.On("Infof", "%v %v: install bash", []any{"...........", "test-play"}).Return().Maybe()

	input := &ty.TaskInput{
		Play:   "test-play",
		Sudo:   false,
		DryRun: true,
	}

	result, err := bi.Install(context.Background(), input)
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.True(t, result.Success)
}

func TestBashInstaller_Install_BrewInstallFails(t *testing.T) {
	mockBrew := &MockBrew{}
	mockExt := &MockExt{}
	mockLog := &MockLogger{}

	bi := &BashInstaller{
		Brew:  mockBrew,
		Utils: mockExt,
		Log:   mockLog,
	}

	mockExt.On("IsOSX").Return(true).Maybe()
	mockExt.On("IsInstalled", "bash").Return(false).Maybe()
	mockLog.On("Infof", "%v %v: install bash", []any{"...........", "test-play"}).Return().Maybe()
	mockBrew.On("Install", context.Background(), "bash").Return(assert.AnError)

	input := &ty.TaskInput{
		Play:   "test-play",
		Sudo:   false,
		DryRun: false,
	}

	result, err := bi.Install(context.Background(), input)
	assert.Error(t, err)
	assert.Nil(t, result)
}
