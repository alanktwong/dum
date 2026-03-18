package playbook

import (
	"awong/dotfiles/pkg/external"
	"awong/dotfiles/pkg/logging"
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// GitTask uses git-clone for installation.
// Use 'id' as the git repo URL.
// Use 'name' as the explicit directory.
type GitTask struct {
	Attributes
	Root string
	Name string
	Git  external.Git
	Log  logging.Logger
}

// NewGitTask constructs a GitTask.
func NewGitTask(attributes *Attributes, root, name string) (*GitTask, error) {
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
		Git:        external.NewGit(),
		Log:        logging.Log(),
	}, nil
}

// GetAttributes implements Attributable.
func (t *GitTask) GetAttributes() Attributes {
	return t.Attributes
}

// GetID implements Identifiable.
func (t *GitTask) GetID() string {
	return t.ID
}

// IsEnabled implements Enableable.
func (t *GitTask) IsEnabled() bool {
	return t.Enabled
}

// Install implements Installer.
func (t *GitTask) Install(ctx context.Context, input *Input) (*TaskResult, error) {
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

// List implements Lister.
func (t *GitTask) List(_ context.Context, input *Input) (*TaskResult, error) {
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
