package tasks

import (
	ext "awong/dotfiles/pkg/external"
	l "awong/dotfiles/pkg/logging"
	ty "awong/dotfiles/pkg/types"
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type GitTask struct {
	ty.Attributes
	Root string
	Name string
	Git  ext.Git
	Log  l.Logger
}

func NewGitTask(attributes *ty.Attributes, root, name string) (*GitTask, error) {
	if attributes == nil {
		return nil, fmt.Errorf("attributes cannot be nil")
	}
	if attributes.ID == "" {
		return nil, fmt.Errorf("task ID cannot be empty")
	}
	if root == "" {
		root = "~/projects"
	}
	providedName := provideName(attributes.ID, name)
	return &GitTask{
		Attributes: *attributes,
		Root:       root,
		Name:       providedName,
		Git:        ext.NewGit(),
		Log:        l.Log(),
	}, nil
}

func (t *GitTask) GetAttributes() ty.Attributes {
	return t.Attributes
}

func (t *GitTask) GetID() string {
	return t.ID
}

func (t *GitTask) IsEnabled() bool {
	return t.Enabled
}

func (t *GitTask) Install(ctx context.Context, input *ty.TaskInput) (*ty.TaskResult, error) {
	t.Log.Debugf("%s START ... play: %s taskID: %s", TaskEllipsis, input.Play, t.ID)
	if t.Root == "" {
		return nil, fmt.Errorf("task root cannot be empty")
	}
	if !t.Enabled {
		return t.CreateTaskResult(input, false)
	}
	providedName := provideName(t.ID, t.Name)
	rootPath := filepath.Join(os.ExpandEnv(t.Root))
	targetPath := filepath.Join(rootPath, providedName)

	if t.Git.AlreadyExists(targetPath) {
		t.Log.Infof("%s %s git dir %v/%v exists", TaskEllipsis, input.Play, t.Root, providedName)
		return t.CreateTaskResult(input, false)
	}
	t.Log.Infof("%s %s: git clone %s %s", TaskEllipsis, input.Play, t.ID, providedName)
	if !input.DryRun {
		if err := t.Git.Clone(ctx, t.ID, t.Name, rootPath, t.Sudo); err != nil {
			return nil, fmt.Errorf("failed to git clone %v: %v", t.ID, err)
		}
	}
	return t.CreateTaskResult(input, true)
}

func (t *GitTask) List(_ context.Context, input *ty.TaskInput) (*ty.TaskResult, error) {
	name := provideName(t.ID, t.Name)
	err := t.Log.Printlnf("%v at %v, git clone %v %v", TaskEllipsis, t.Root, t.ID, name)
	if err != nil {
		return nil, fmt.Errorf("failed to list git: %v", err)
	}
	return t.CreateTaskResult(input, true)
}

func provideName(id, name string) string {
	if name != "" {
		return name
	}
	parts := strings.Split(strings.TrimSuffix(id, ".git"), "/")
	return parts[len(parts)-1]
}
