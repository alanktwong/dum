// Package tasks provides Task types and related utilities.
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

// GitTask clones a git repository.
type GitTask struct {
	ty.Attributes
	Root string
	Name string
	Git  ext.Git
	Log  l.Logger
}

// NewGitTask returns a new GitTask for cloning a git repository.
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

// GetAttributes returns the Attributes.
func (t *GitTask) GetAttributes() ty.Attributes {
	return t.Attributes
}

// GetID returns the ID.
func (t *GitTask) GetID() string {
	return t.ID
}

// IsEnabled returns whether the task is enabled.
func (t *GitTask) IsEnabled() bool {
	return t.Enabled
}

// Install installs the task.
func (t *GitTask) Install(ctx context.Context, input *ty.TaskInput) (*ty.TaskResult, error) {
	t.Log.Debugf("%s START ... play: %s taskID: %s", TaskEllipsis, input.Play, t.ID)
	if t.Root == "" {
		return nil, fmt.Errorf("task root cannot be empty")
	}
	if !t.Enabled {
		result, err := t.CreateTaskResult(input, false)
		if err != nil {
			return nil, fmt.Errorf("failed to create task result: %w", err)
		}
		return result, nil
	}
	providedName := provideName(t.ID, t.Name)
	rootPath := filepath.Join(os.ExpandEnv(t.Root))
	targetPath := filepath.Join(rootPath, providedName)

	if t.Git.AlreadyExists(targetPath) {
		t.Log.Infof("%s %s git dir %v/%v exists", TaskEllipsis, input.Play, t.Root, providedName)
		result, err := t.CreateTaskResult(input, false)
		if err != nil {
			return nil, fmt.Errorf("failed to create task result: %w", err)
		}
		return result, nil
	}
	t.Log.Infof("%s %s: git clone %s %s", TaskEllipsis, input.Play, t.ID, providedName)
	if !input.DryRun {
		if err := t.Git.Clone(ctx, t.ID, t.Name, rootPath, t.Sudo); err != nil {
			return nil, fmt.Errorf("failed to git clone %v: %v", t.ID, err)
		}
	}
	result, err := t.CreateTaskResult(input, true)
	if err != nil {
		return nil, fmt.Errorf("failed to create task result: %w", err)
	}
	return result, nil
}

// List lists the task.
func (t *GitTask) List(_ context.Context, input *ty.TaskInput) (*ty.TaskResult, error) {
	name := provideName(t.ID, t.Name)
	err := t.Log.Printlnf("%v at %v, git clone %v %v", TaskEllipsis, t.Root, t.ID, name)
	if err != nil {
		return nil, fmt.Errorf("failed to list git: %v", err)
	}
	result, err := t.CreateTaskResult(input, true)
	if err != nil {
		return nil, fmt.Errorf("failed to create task result: %w", err)
	}
	return result, nil
}

func provideName(id, name string) string {
	if name != "" {
		return name
	}
	parts := strings.Split(strings.TrimSuffix(id, ".git"), "/")
	return parts[len(parts)-1]
}
