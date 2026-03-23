package playbook

import (
	"awong/dotfiles/pkg/external"
	"awong/dotfiles/pkg/logging"
	"context"
	"fmt"
)

// BashTask executes a bash command.
// Use 'command' or 'script' in attributes to specify what to run.
type BashTask struct {
	Attributes
	Utils external.Ext
	Log   logging.Logger
}

// NewBashTask constructs a BashTask.
func NewBashTask(attributes *Attributes) (*BashTask, error) {
	if attributes == nil {
		return nil, fmt.Errorf("attributes cannot be nil")
	}
	if attributes.ID == "" {
		return nil, fmt.Errorf("task ID cannot be empty")
	}
	if attributes.Command == "" && attributes.Script == "" {
		return nil, fmt.Errorf("either command or script must be specified")
	}
	return &BashTask{
		Attributes: *attributes,
		Utils:      external.NewExt(),
		Log:        logging.Log(),
	}, nil
}

// GetAttributes implements Attributable.
func (t *BashTask) GetAttributes() Attributes {
	return t.Attributes
}

// GetID implements Identifiable.
func (t *BashTask) GetID() string {
	return t.ID
}

// IsEnabled implements Enableable.
func (t *BashTask) IsEnabled() bool {
	return t.Enabled
}

// Install implements Installer.
func (t *BashTask) Install(ctx context.Context, input *Input) (*TaskResult, error) {
	t.Log.Debugf("%s START ... play: %s taskID: %s", TaskEllipsis, input.Play, t.ID)
	if !t.Enabled {
		return t.CreateTaskResult(input, false)
	}

	cmd := t.Command
	if cmd == "" {
		cmd = t.Script
	}

	t.Log.Infof("%s %s: %s", TaskEllipsis, input.Play, cmd)
	if !input.DryRun {
		err := t.Utils.RunCommand(ctx, cmd, t.Sudo)
		if err != nil {
			return nil, fmt.Errorf("failed to run command %s: %w", t.ID, err)
		}
	}
	return t.CreateTaskResult(input, true)
}

// List implements Lister.
func (t *BashTask) List(_ context.Context, input *Input) (*TaskResult, error) {
	cmd := t.Command
	if cmd == "" {
		cmd = t.Script
	}
	err := t.Log.Printlnf("%s %s", TaskEllipsis, cmd)
	if err != nil {
		return nil, fmt.Errorf("failed to list command: %w", err)
	}
	return t.CreateTaskResult(input, true)
}
