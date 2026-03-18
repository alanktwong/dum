package playbook

import (
	"awong/dotfiles/pkg/external"
	"awong/dotfiles/pkg/logging"
	"context"
	"fmt"
)

// DirTask creates a directory configuration or installation.
// Use 'id' to specify the directory.
type DirTask struct {
	Attributes
	Utils external.Ext
	Log   logging.Logger
}

// NewDirTask constructs a DirTask.
func NewDirTask(attributes *Attributes) (*DirTask, error) {
	if attributes == nil {
		return nil, fmt.Errorf("attributes cannot be nil")
	}
	if attributes.ID == "" {
		return nil, fmt.Errorf("task ID cannot be empty")
	}
	return &DirTask{
		Attributes: *attributes,
		Utils:      external.NewExt(),
		Log:        logging.Log(),
	}, nil
}

// GetAttributes implements Attributable.
func (t *DirTask) GetAttributes() Attributes {
	return t.Attributes
}

// GetID implements Identifiable.
func (t *DirTask) GetID() string {
	return t.ID
}

// IsEnabled implements Enableable.
func (t *DirTask) IsEnabled() bool {
	return t.Enabled
}

// Install implements Installer.
func (t *DirTask) Install(ctx context.Context, input *Input) (*TaskResult, error) {
	t.Log.Debugf("%s START ... play: %s taskID: %s", TaskEllipsis, input.Play, t.ID)
	if !t.Enabled {
		return t.CreateTaskResult(input, false)
	}
	path, err := t.Utils.ExpandUser(t.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to expand user for %s: %v", t.ID, err)
	}
	if t.Utils.IsDir(path) {
		t.Log.Infof("%s %s: %s already exists", TaskEllipsis, input.Play, t.ID)
		return t.CreateTaskResult(input, false)
	}
	t.Log.Infof("%s %s: mkdir -p %s", TaskEllipsis, input.Play, t.ID)
	if !input.DryRun {
		err = t.Utils.CreateDirectory(ctx, path, t.Sudo)
		if err != nil {
			return nil, fmt.Errorf("failed to mkdir -p %v: %v", t.ID, err)
		}
	}
	return t.CreateTaskResult(input, true)
}

// List implements Lister.
func (t *DirTask) List(_ context.Context, input *Input) (*TaskResult, error) {
	err := t.Log.Printlnf("%v mkdir -p %s", TaskEllipsis, t.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to list dir: %v", err)
	}
	return t.CreateTaskResult(input, true)
}
