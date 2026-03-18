package playbook

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func findLinkTask(t *testing.T, input *Input) *LinkTask {
	id := "../projects/dotfiles"
	tk, err := findTask(input.PlayBook, id)
	assert.NoError(t, err)
	task, ok := tk.(*LinkTask)
	assert.True(t, ok)
	return task
}

func TestLinkTask_Install(t *testing.T) {
	input := createTestInput(t)
	input.DryRun = false
	task := findLinkTask(t, input)
	ctx := context.Background()

	mockUtils := NewMockExt(t)
	rel := "projects"
	rootPath := fmt.Sprintf("/Users/user/%v", rel)
	providedTarget := task.provideTarget()
	mockUtils.EXPECT().ExpandUser(task.Root).Return(rootPath, nil).Once()

	targetPath := rootPath + "/dotfiles"
	mockUtils.EXPECT().IsSymlink(targetPath).Return(false).Once()
	mockUtils.EXPECT().SoftLink(ctx, rootPath, task.ID, providedTarget, task.Sudo).Return(nil).Once()
	task.Utils = mockUtils
	// when
	got, err := task.Install(ctx, input)
	// then
	assert.NoError(t, err)
	wantedResult := expectTaskResult(t, true, task.Attributes, input)
	assert.Equalf(t, wantedResult, got, "Install(%v): should install for happy path", input)
}
