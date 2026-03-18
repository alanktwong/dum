package playbook

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStarshipInstaller_Install_Osx(t *testing.T) {
	input := createTestInput(t)
	input.DryRun = false
	task := findFunctionTask(t, input)
	task.ID = "install_starship"
	input.Task = task.ID
	ctx := context.Background()
	runCtx := WithSudo(ctx, input.Sudo)

	installer := NewStarshipInstaller()
	mockCurl := NewMockRunner(t)
	mockCurl.EXPECT().Run(runCtx).Return(nil).Once()
	mockLinuxCurl := NewMockRunner(t)

	mockUtils := NewMockExt(t)
	mockUtils.EXPECT().IsInstalled("starship").Return(false).Once()
	mockUtils.EXPECT().IsOSX().Return(true).Once()

	installer.Curl = mockCurl
	installer.LinuxCurl = mockLinuxCurl
	installer.Utils = mockUtils
	task.Registry[task.ID] = installer
	// when
	got, err := installer.Install(ctx, input)
	// then
	assert.NoError(t, err)
	wantedResult := expectTaskResult(t, true, task.Attributes, input)
	assert.Equalf(t, wantedResult, got, "Install(%v)", input)
}

func TestStarshipInstaller_Install_Osx_Sudo(t *testing.T) {
	input := createTestInput(t)
	input.DryRun = false
	input.Sudo = true
	task := findFunctionTask(t, input)
	task.ID = "install_starship"
	input.Task = task.ID
	ctx := context.Background()
	runCtx := WithSudo(ctx, input.Sudo)

	installer := NewStarshipInstaller()
	mockCurl := NewMockRunner(t)
	mockCurl.EXPECT().Run(runCtx).Return(nil).Once()
	mockLinuxCurl := NewMockRunner(t)

	mockUtils := NewMockExt(t)
	mockUtils.EXPECT().IsInstalled("starship").Return(false).Once()
	mockUtils.EXPECT().IsOSX().Return(true).Once()

	installer.Curl = mockCurl
	installer.LinuxCurl = mockLinuxCurl
	installer.Utils = mockUtils
	task.Registry[task.ID] = installer
	// when
	got, err := installer.Install(ctx, input)
	// then
	assert.NoError(t, err)
	wantedResult := expectTaskResult(t, true, task.Attributes, input)
	assert.Equalf(t, wantedResult, got, "Install(%v)", input)
}

func TestStarshipInstaller_Install_Linux(t *testing.T) {
	input := createTestInput(t)
	input.DryRun = false
	task := findFunctionTask(t, input)
	task.ID = "install_starship"
	input.Task = task.ID
	ctx := context.Background()
	runCtx := WithSudo(ctx, input.Sudo)

	installer := NewStarshipInstaller()
	mockCurl := NewMockRunner(t)
	mockLinuxCurl := NewMockRunner(t)
	mockLinuxCurl.EXPECT().Run(runCtx).Return(nil).Once()
	// mockLinuxSudoCurl := NewMockRunner(t)

	mockUtils := NewMockExt(t)
	mockUtils.EXPECT().IsInstalled("starship").Return(false).Once()
	mockUtils.EXPECT().IsOSX().Return(false).Once()
	mockUtils.EXPECT().IsLinux().Return(true).Once()

	installer.Curl = mockCurl
	installer.LinuxCurl = mockLinuxCurl
	installer.Utils = mockUtils
	task.Registry[task.ID] = installer
	// when
	got, err := installer.Install(ctx, input)
	// then
	assert.NoError(t, err)
	wantedResult := expectTaskResult(t, true, task.Attributes, input)
	assert.Equalf(t, wantedResult, got, "Install(%v)", input)
}

func TestStarshipInstaller_Install_Linux_Sudo(t *testing.T) {
	input := createTestInput(t)
	input.DryRun = false
	input.Sudo = true
	task := findFunctionTask(t, input)
	task.ID = "install_starship"
	input.Task = task.ID
	ctx := context.Background()
	runCtx := WithSudo(ctx, input.Sudo)

	installer := NewStarshipInstaller()
	mockCurl := NewMockRunner(t)
	mockLinuxCurl := NewMockRunner(t)
	mockLinuxCurl.EXPECT().Run(runCtx).Return(nil).Once()

	mockUtils := NewMockExt(t)
	mockUtils.EXPECT().IsInstalled("starship").Return(false).Once()
	mockUtils.EXPECT().IsOSX().Return(false).Once()
	mockUtils.EXPECT().IsLinux().Return(true).Once()

	installer.Curl = mockCurl
	installer.LinuxCurl = mockLinuxCurl
	installer.Utils = mockUtils
	task.Registry[task.ID] = installer
	// when
	got, err := installer.Install(ctx, input)
	// then
	assert.NoError(t, err)
	wantedResult := expectTaskResult(t, true, task.Attributes, input)
	assert.Equalf(t, wantedResult, got, "Install(%v)", input)
}
