package plays

import (
	ext "awong/dotfiles/pkg/external"
	l "awong/dotfiles/pkg/logging"
	t "awong/dotfiles/pkg/types"
	"context"
	"fmt"
)

// PlayExecutor performs the work of executing a play or task.
type PlayExecutor struct {
	Log l.Logger
	Ext ext.Ext
}

// NewPlayExecutor constructs a PlayExecutot.
func NewPlayExecutor() *PlayExecutor {
	return &PlayExecutor{
		Log: l.Log(),
		Ext: ext.NewExt(),
	}
}

func (e *PlayExecutor) Initialize(input *PlayInput) (*PlayResult, error) {
	if !e.Ext.IsInstalled("tar") {
		return nil, fmt.Errorf("tar is not installed, please install it to continue")
	}
	if !e.Ext.IsInstalled("zip") {
		return nil, fmt.Errorf("zip is not installed, please install it to continue")
	}
	if !e.Ext.IsInstalled("unzip") {
		return nil, fmt.Errorf("unzip is not installed, please install it to continue")
	}
	empty := make([]*t.TaskResult, 0)
	return NewPlayResult(input, empty), nil
}

// InstallPlay installs a play.
func (e *PlayExecutor) InstallPlay(ctx context.Context, p *Play, input *PlayInput) (*PlayResult, error) {
	var taskResults []*t.TaskResult
	err := e.Log.Printlnf("%v Installing play (%v) ... %v", PlayEllipsis, p.ID, p.Description)
	if err != nil {
		return nil, fmt.Errorf("failed to install play %s: %w", p.ID, err)
	}

	taskMap := p.GetTasks(false)
	for id, task := range taskMap.AllFromFront() {
		if attributable, ok := task.(t.Attributable); ok {
			taskInput := t.NewTaskInput(
				input.DryRun,
				attributable.GetAttributes().Sudo,
				id,
				p.ID,
				input.PlayBook.GetID(),
				input.PlayBook.GetJetBrainsApps(),
			)
			taskResult, err := task.Install(ctx, taskInput)
			if err != nil {
				return nil, fmt.Errorf("failed to install task %s in play %v: %w", id, p.ID, err)
			}
			taskResults = append(taskResults, taskResult)
		}
	}
	return NewPlayResult(input, taskResults), nil
}

// ListPlay lists a play.
func (e *PlayExecutor) ListPlay(ctx context.Context, p *Play, input *PlayInput) (*PlayResult, error) {
	var taskResults []*t.TaskResult
	err := e.Log.Printlnf("%v Listing play (%v) ... %v", PlayEllipsis, p.ID, p.Description)
	if err != nil {
		return nil, fmt.Errorf("failed to list play %s: %w", p.ID, err)
	}
	taskMap := p.GetTasks(false)
	for id, task := range taskMap.AllFromFront() {
		if attributable, ok := task.(t.Attributable); ok {
			taskInput := t.NewTaskInput(
				input.DryRun,
				attributable.GetAttributes().Sudo,
				id,
				p.ID,
				input.PlayBook.GetID(),
				input.PlayBook.GetJetBrainsApps(),
			)
			taskResult, err := task.List(ctx, taskInput)
			if err != nil {
				return nil, fmt.Errorf("failed to list task %s in play %v: %w", id, p.ID, err)
			}
			taskResults = append(taskResults, taskResult)
		}
	}
	return NewPlayResult(input, taskResults), nil
}
