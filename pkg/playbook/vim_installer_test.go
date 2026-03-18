package playbook

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestVimInstaller_Install(t *testing.T) {
	input := createTestInput(t)
	input.DryRun = false
	task := findFunctionTask(t, input)
	task.ID = "install_vim"
	input.Task = task.ID
	ctx := context.Background()

	installer := NewVimInstaller()
	mockBrew := NewMockBrew(t)
	mockBrew.EXPECT().Install(ctx, "vim").Return(nil).Once()

	mockUtils := NewMockExt(t)
	mockUtils.EXPECT().IsInstalled("ovim").Return(false).Once()

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
