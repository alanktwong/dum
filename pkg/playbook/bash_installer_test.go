package playbook

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBashInstaller_Install(t *testing.T) {
	input := createTestInput(t)
	input.DryRun = false
	task := findFunctionTask(t, input)
	task.ID = "install_bash"
	input.Task = task.ID
	ctx := context.Background()

	installer := NewBashInstaller()
	mockBrew := NewMockBrew(t)
	mockBrew.EXPECT().Install(ctx, "bash").Return(nil).Once()
	mockBrew.EXPECT().Prefix(ctx).Return("/opt/homebrew", nil)

	mockUtils := NewMockExt(t)
	mockUtils.EXPECT().IsOSX().Return(true).Once()
	mockUtils.EXPECT().IsInstalled("bash").Return(false).Once()
	mockUtils.EXPECT().SoftLink(ctx, "/usr/local/bin", "/opt/homebrew/bin/bash", "bash", input.Sudo).
		Return(nil).Once()

	installer.Brew = mockBrew
	installer.Utils = mockUtils
	task.Registry[task.ID] = installer
	// when
	got, err := installer.Install(ctx, input)
	// then
	assert.NoError(t, err)
	wantedResult := expectTaskResult(t, true, task.Attributes, input)
	assert.Equalf(t, wantedResult, got, "Install(%v)", input)
}
