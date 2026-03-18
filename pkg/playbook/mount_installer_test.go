package playbook

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMountInstaller_Install(t *testing.T) {
	input := createTestInput(t)
	input.DryRun = false
	task := findFunctionTask(t, input)
	task.ID = "install_case_sensitive_mount"
	input.Task = task.ID
	ctx := context.Background()

	installer := NewMountInstaller()
	mockRunner := NewMockRunner(t)
	mockRunner.EXPECT().Run(ctx).Return(nil).Once()

	mockUtils := NewMockExt(t)
	mockUtils.EXPECT().IsOSX().Return(true).Once()

	installer.Utils = mockUtils
	installer.Runner = mockRunner
	task.Registry[task.ID] = installer
	// when
	got, err := installer.Install(ctx, input)
	// then
	assert.NoError(t, err)
	wantedResult := expectTaskResult(t, true, task.Attributes, input)
	assert.Equalf(t, wantedResult, got, "Install(%v)", input)
}
