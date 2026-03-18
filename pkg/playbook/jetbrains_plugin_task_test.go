package playbook

import (
	e "awong/dotfiles/pkg/enums"
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func findJetBrainsTask(t *testing.T, ctx *Input) *JetBrainsPluginTask {
	id := "org.asciidoctor.intellij.asciidoc"
	tk, err := findTask(ctx.PlayBook, id)
	assert.NoError(t, err)
	task, ok := tk.(*JetBrainsPluginTask)
	assert.True(t, ok)
	return task
}

func TestJetBrainsPluginTask_Install(t *testing.T) {
	input := createTestInput(t)
	input.DryRun = false
	task := findJetBrainsTask(t, input)
	apps := input.PlayBook.JetBrainsApps
	ctx := context.Background()

	mockUtils := NewMockExt(t)
	mockUtils.EXPECT().IsInstalled(e.JetBrainsTypeIdea.String()).Return(true).Once()
	mockUtils.EXPECT().IsInstalled(e.JetBrainsTypeGoland.String()).Return(true).Once()

	mockJetBrains := NewMockJetBrainsApp(t)
	intellij, ok := apps[e.JetBrainsTypeIdea.String()]
	assert.True(t, ok)
	mockJetBrains.EXPECT().IsInstalled(intellij, task.ID).Return(false).Once()
	goland, ok := apps[e.JetBrainsTypeGoland.String()]
	assert.True(t, ok)
	mockJetBrains.EXPECT().IsInstalled(goland, task.ID).Return(false).Once()
	mockJetBrains.EXPECT().Install(ctx, e.JetBrainsTypeIdea.String(), task.ID).Return(nil).Once()
	mockJetBrains.EXPECT().Install(ctx, e.JetBrainsTypeGoland.String(), task.ID).Return(nil).Once()
	// when
	task.JetBrains = mockJetBrains
	task.Utils = mockUtils
	got, err := task.Install(ctx, input)
	// then
	assert.NoError(t, err)
	wantedResult := expectTaskResult(t, true, task.Attributes, input)
	assert.Equalf(t, wantedResult, got, "Install(%v) should mkdir", input)
}
