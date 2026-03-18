package playbook

import (
	e "awong/dotfiles/pkg/enums"
	"context"
	"fmt"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestExecutor_Install(t *testing.T) {
	// given: context
	ctx := context.Background()
	input := createTestInput(t)
	input.Play = ""
	input.DryRun = false
	// and mocks
	mockBrew := NewMockBrew(t)
	mockUtils := NewMockExt(t)
	mockUtils.EXPECT().IsInstalled("tar").Return(true)
	mockUtils.EXPECT().IsInstalled("zip").Return(true)
	mockUtils.EXPECT().IsInstalled("unzip").Return(true)
	mockTapInstaller := NewMockInstaller(t)
	// and: brew cask
	brewCaskTask := findBrewCaskTask(t, input)

	mockUtils.EXPECT().IsOSX().Return(true)
	wantedTapResult := expectTaskResult(t, false, brewCaskTask.Attributes, input)
	mockTapInstaller.EXPECT().Install(ctx, input).Return(wantedTapResult, nil)

	mockBrew.EXPECT().InPath(ctx, "Caskroom", brewCaskTask.ID).Return(false)
	mockBrew.EXPECT().InstallCask(ctx, brewCaskTask.ID).Return(nil).Once()

	brewCaskTask.Brew = mockBrew
	brewCaskTask.Utils = mockUtils
	brewCaskTask.Tap = mockTapInstaller
	brewCaskTask.Enabled = true
	// and: brew cellar
	brewCellarTask := findBrewCellarTask(t, input)
	mockBrew.EXPECT().InPath(ctx, "Cellar", brewCellarTask.ID).Return(false)
	mockBrew.EXPECT().Install(ctx, brewCellarTask.ID).Return(nil).Once()

	brewCellarTask.Brew = mockBrew
	brewCellarTask.Tap = mockTapInstaller
	brewCellarTask.Enabled = true
	// and: brew
	brewTask := findBrewTask(t, input)
	mockBrew.EXPECT().InPath(ctx, "opt", brewTask.ID).Return(false)
	mockUtils.EXPECT().IsInstalled(brewTask.ID).Return(false)
	mockBrew.EXPECT().Install(ctx, brewTask.ID).Return(nil).Once()

	brewTask.Brew = mockBrew
	brewTask.Utils = mockUtils
	brewTask.Tap = mockTapInstaller
	brewTask.Enabled = true
	// and: dir
	dirTask := findDirTask(t, input)
	expanded := fmt.Sprintf("/Users/user/%v", dirTask.ID)
	mockUtils.EXPECT().ExpandUser(dirTask.ID).Return(expanded, nil)
	mockUtils.EXPECT().IsDir(expanded).Return(false)
	mockUtils.EXPECT().CreateDirectory(ctx, expanded, dirTask.Sudo).Return(nil).Once()

	dirTask.Utils = mockUtils
	dirTask.Enabled = true
	// and: link
	linkTask := findLinkTask(t, input)
	rel := "projects"
	rootPath := fmt.Sprintf("/Users/user/%v", rel)
	linkUtils := NewMockExt(t)
	providedTarget := linkTask.provideTarget()
	linkUtils.EXPECT().ExpandUser(linkTask.Root).Return(rootPath, nil)
	targetPath := rootPath + "/dotfiles"
	linkUtils.EXPECT().IsSymlink(targetPath).Return(false)
	linkUtils.EXPECT().SoftLink(ctx, rootPath, linkTask.ID, providedTarget, linkTask.Sudo).Return(nil).Once()

	linkTask.Utils = linkUtils
	linkTask.Enabled = true
	// and: function
	functionTask := findFunctionTask(t, input)
	mockRegistry := make(map[string]Installer)
	mockFnInstaller := NewMockInstaller(t)
	wantedResult := expectTaskResult(t, true, functionTask.Attributes, input)
	mockFnInstaller.EXPECT().Install(ctx, input).Return(wantedResult, nil).Once()
	mockRegistry[functionTask.ID] = mockFnInstaller
	functionTask.Registry = mockRegistry
	functionTask.Enabled = true
	// and: git
	gitTask := findGitTask(t, input)

	mockGit := NewMockGit(t)
	mockGit.EXPECT().Clone(ctx, gitTask.ID, gitTask.Name, gitTask.Root, gitTask.Sudo).Return(nil).Once()
	gitPath := filepath.Join(gitTask.Root, gitTask.Name)
	mockGit.EXPECT().AlreadyExists(gitPath).Return(false)
	gitTask.Git = mockGit
	gitTask.Enabled = true
	// and: jetbrains
	jetBrainsTask := findJetBrainsTask(t, input)
	apps := input.PlayBook.JetBrainsApps
	mockUtils.EXPECT().IsInstalled(e.JetBrainsTypeIdea.String()).Return(true)
	mockUtils.EXPECT().IsInstalled(e.JetBrainsTypeGoland.String()).Return(true)

	mockJetBrains := NewMockJetBrainsApp(t)
	intellij, ok := apps[e.JetBrainsTypeIdea.String()]
	assert.True(t, ok)
	mockJetBrains.EXPECT().IsInstalled(intellij, jetBrainsTask.ID).Return(false)
	goland, ok := apps[e.JetBrainsTypeGoland.String()]
	assert.True(t, ok)
	mockJetBrains.EXPECT().IsInstalled(goland, jetBrainsTask.ID).Return(false)
	mockJetBrains.EXPECT().Install(ctx, e.JetBrainsTypeIdea.String(), jetBrainsTask.ID).Return(nil).Once()
	mockJetBrains.EXPECT().Install(ctx, e.JetBrainsTypeGoland.String(), jetBrainsTask.ID).Return(nil).Once()

	jetBrainsTask.JetBrains = mockJetBrains
	jetBrainsTask.Utils = mockUtils
	jetBrainsTask.Enabled = true
	// and: mas
	masTask := findMasTask(t, input)

	mockMas := NewMockMas(t)
	mockMas.EXPECT().List(ctx).Return("", nil)
	mockMas.EXPECT().Install(ctx, masTask.ID).Return(nil).Once()

	masTask.Mas = mockMas
	masTask.Enabled = true
	// and: vscode
	vsCodeTask := findVsCodeTask(t, input)

	mockCode := NewMockCode(t)
	mockCode.EXPECT().ListExtensions(ctx).Return("", nil)
	mockCode.EXPECT().InstallExtension(ctx, vsCodeTask.ID).Return(nil).Once()
	vsCodeTask.Code = mockCode
	vsCodeTask.Enabled = true

	// when
	executor := NewExecutor()
	executor.Utils = mockUtils
	got, err := executor.Install(ctx, input)
	// then
	assert.NoError(t, err)
	assert.NotNilf(t, got, "executor should install(%v) all active tasks", input)
}
