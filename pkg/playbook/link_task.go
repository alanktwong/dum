package playbook

import (
	"awong/dotfiles/pkg/external"
	"awong/dotfiles/pkg/logging"
	"context"
	"fmt"
	"path/filepath"
)

// LinkTask creates a soft link for configuration or installation.
// Use 'id' to specify source of the link.
type LinkTask struct {
	Attributes
	Root   string
	Target string
	Utils  external.Ext
	Log    logging.Logger
}

// NewLinkTask constructs a LinkTask.
func NewLinkTask(attributes *Attributes, root, target string) (*LinkTask, error) {
	if attributes == nil {
		return nil, fmt.Errorf("attributes cannot be nil")
	}
	if attributes.ID == "" {
		return nil, fmt.Errorf("task ID cannot be empty")
	}
	return &LinkTask{
		Attributes: *attributes,
		Root:       root,
		Target:     target,
		Utils:      external.NewExt(),
		Log:        logging.Log(),
	}, nil
}

// GetAttributes implements Attributable.
func (t *LinkTask) GetAttributes() Attributes {
	return t.Attributes
}

// GetID implements Identifiable.
func (t *LinkTask) GetID() string {
	return t.ID
}

// IsEnabled implements Enableable.
func (t *LinkTask) IsEnabled() bool {
	return t.Enabled
}

// Install implements Installer.
func (t *LinkTask) Install(ctx context.Context, input *Input) (*TaskResult, error) {
	t.Log.Debugf("%s START ... play: %s taskID: %s", TaskEllipsis, input.Play, t.ID)
	if t.Root == "" {
		return t.CreateTaskResult(input, false)
	}
	target := t.provideTarget()
	if target == "" {
		return t.CreateTaskResult(input, false)
	}
	if !t.Enabled {
		return t.CreateTaskResult(input, false)
	}
	rootPath, err := t.Utils.ExpandUser(t.Root)
	if err != nil {
		return nil, fmt.Errorf("failed to expand root path %s: %w", t.Root, err)
	}
	targetPath := filepath.Join(rootPath, target)
	if t.Utils.IsSymlink(targetPath) {
		t.Log.Infof("%v %v: %v -> %s/%s already linked", TaskEllipsis, input.Play, t.ID, t.Root, target)
	}
	t.Log.Infof("%v %v: planning to link %v -> %s/%s", TaskEllipsis, input.Play, t.ID, t.Root, target)
	if !input.DryRun {
		if t.Utils.SoftLink(ctx, rootPath, t.ID, target, t.Sudo) != nil {
			return nil, fmt.Errorf("failed to link %v -> %s/%s: %v", t.ID, t.Root, target, err)
		}
	}
	return t.CreateTaskResult(input, true)
}

// List implements Lister.
func (t *LinkTask) List(_ context.Context, input *Input) (*TaskResult, error) {
	target := t.provideTarget()
	err := t.Log.Printlnf("%v linking %v -> %s/%s", TaskEllipsis, t.ID, t.Root, target)
	if err != nil {
		return nil, fmt.Errorf("failed to list link: %v", err)
	}
	return t.CreateTaskResult(input, true)
}

func (t *LinkTask) provideTarget() string {
	if t.Target != "" {
		return t.Target
	}
	return filepath.Base(t.ID)
}
