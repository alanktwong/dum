package playbook

import (
	"awong/dotfiles/pkg/external"
	l "awong/dotfiles/pkg/logging"
	"context"
	"fmt"

	om "github.com/elliotchance/orderedmap/v3"
)

// Executor performs the work of executing a play or task.
type Executor struct {
	Log   l.Logger
	Utils external.Ext
}

// NewExecutor constructs an Executor.
func NewExecutor() *Executor {
	return &Executor{
		Log:   l.Log(),
		Utils: external.NewExt(),
	}
}

// Install installs a playbook.
func (e *Executor) Install(ctx context.Context, input *Input) (*Result, error) {
	group := input.Play
	pb := input.PlayBook
	err := e.Log.Printlnf(Ellipsis+"Installing playbook (%v) ... %v", pb.ID, pb.Description)
	if err != nil {
		return nil, fmt.Errorf("failed to install playbook %v", err)
	}
	var playResults []*PlayResult
	initResult, err := e.initialize(input)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize %v", err)
	}
	playResults = append(playResults, initResult)

	playMap := pb.GetPlays(true)
	if group != "" {
		results, err := e.installGroupPlay(ctx, group, playResults, playMap, input)
		if err != nil {
			return nil, fmt.Errorf("failed to install group %v", err)
		}
		playResults = results
	} else {
		results, err := e.installAllPlays(ctx, playResults, playMap, input)
		if err != nil {
			return nil, fmt.Errorf("failed to install group %v", err)
		}
		playResults = results
	}
	e.Log.Infof(Ellipsis+"END: installing playbook (%v) ... %v", pb.ID, pb.Description)
	return NewResult(input, playResults), nil
}

func (e *Executor) installGroupPlay(ctx context.Context,
	group string,
	playResults []*PlayResult,
	playMap *om.OrderedMap[string, *Play],
	input *Input,
) ([]*PlayResult, error) {
	if play, ok := playMap.Get(group); ok {
		playResult, err := e.InstallPlay(ctx, play, input)
		if err != nil {
			return nil, fmt.Errorf("failed to install play %s: %w", play.GetID(), err)
		}
		playResults = append(playResults, playResult)
	}
	return playResults, nil
}

func (e *Executor) installAllPlays(ctx context.Context,
	playResults []*PlayResult,
	playMap *om.OrderedMap[string, *Play],
	input *Input,
) ([]*PlayResult, error) {
	for current, play := range playMap.AllFromFront() {
		e.Log.Infof("%v %v", PlayEllipsis, current)
		playResult, err := e.InstallPlay(ctx, play, input)
		if err != nil {
			return nil, fmt.Errorf("failed to install play %s: %w", play.GetID(), err)
		}
		playResults = append(playResults, playResult)
	}
	return playResults, nil
}

// List lists a playbook.
func (e *Executor) List(ctx context.Context, input *Input) (*Result, error) {
	group := input.Play
	pb := input.PlayBook
	err := e.Log.Printlnf(Ellipsis+"Listing playbook (%v) ... %v", pb.ID, pb.Description)
	if err != nil {
		return nil, fmt.Errorf("failed to list playbook: %v", err)
	}
	var playResults []*PlayResult
	playMap := pb.GetPlays(false)
	if group != "" {
		results, err := e.listGroupPlay(ctx, group, playResults, playMap, input)
		if err != nil {
			return nil, fmt.Errorf("failed to list group %v", err)
		}
		playResults = results
	} else {
		results, err := e.listAllPlays(ctx, playResults, playMap, input)
		if err != nil {
			return nil, fmt.Errorf("failed to list %v", err)
		}
		playResults = results
	}
	e.Log.Infof(Ellipsis+"END: listing playbook (%v) ... %v", pb.ID, pb.Description)
	return NewResult(input, playResults), nil
}

func (e *Executor) listGroupPlay(ctx context.Context,
	group string,
	playResults []*PlayResult,
	playMap *om.OrderedMap[string, *Play],
	input *Input,
) ([]*PlayResult, error) {
	if play, ok := playMap.Get(group); ok {
		playResult, err := e.ListPlay(ctx, play, input)
		if err != nil {
			return nil, fmt.Errorf("failed to list play %s: %w", play.GetID(), err)
		}
		playResults = append(playResults, playResult)
	}
	return playResults, nil
}

func (e *Executor) listAllPlays(ctx context.Context,
	playResults []*PlayResult,
	playMap *om.OrderedMap[string, *Play],
	input *Input,
) ([]*PlayResult, error) {
	for _, play := range playMap.AllFromFront() {
		playResult, err := e.ListPlay(ctx, play, input)
		if err != nil {
			return nil, fmt.Errorf("failed to list play %s: %w", play.GetID(), err)
		}
		playResults = append(playResults, playResult)
	}
	return playResults, nil
}

// InstallPlay installs a play.
func (e *Executor) InstallPlay(ctx context.Context, p *Play, input *Input) (*PlayResult, error) {
	input.Play = p.ID
	var taskResults []*TaskResult
	err := e.Log.Printlnf("%v Installing play (%v) ... %v", PlayEllipsis, p.ID, p.Description)
	if err != nil {
		return nil, fmt.Errorf("failed to install play %s: %w", p.ID, err)
	}

	taskMap := p.GetTasks(false)
	for id, task := range taskMap.AllFromFront() {
		input.Task = id
		if attributable, ok := task.(Attributable); ok {
			input.Sudo = attributable.GetAttributes().Sudo
			taskResult, err := task.Install(ctx, input)
			if err != nil {
				return nil, fmt.Errorf("failed to install task %s in play %v: %w", id, p.ID, err)
			}
			taskResults = append(taskResults, taskResult)
		}
	}
	return NewPlayResult(input, taskResults), nil
}

// ListPlay lists a play.
func (e *Executor) ListPlay(ctx context.Context, p *Play, input *Input) (*PlayResult, error) {
	input.Play = p.ID
	var taskResults []*TaskResult
	err := e.Log.Printlnf("%v Listing play (%v) ... %v", PlayEllipsis, p.ID, p.Description)
	if err != nil {
		return nil, fmt.Errorf("failed to list play %s: %w", p.ID, err)
	}
	taskMap := p.GetTasks(false)
	for id, task := range taskMap.AllFromFront() {
		input.Task = id
		if attributable, ok := task.(Attributable); ok {
			input.Sudo = attributable.GetAttributes().Sudo
			taskResult, err := task.List(ctx, input)
			if err != nil {
				return nil, fmt.Errorf("failed to list task %s in play %v: %w", id, p.ID, err)
			}
			taskResults = append(taskResults, taskResult)
		}
	}
	return NewPlayResult(input, taskResults), nil
}

func (e *Executor) initialize(input *Input) (*PlayResult, error) {
	if !e.Utils.IsInstalled("tar") {
		return nil, fmt.Errorf("tar is not installed, please install it to continue")
	}
	if !e.Utils.IsInstalled("zip") {
		return nil, fmt.Errorf("zip is not installed, please install it to continue")
	}
	if !e.Utils.IsInstalled("unzip") {
		return nil, fmt.Errorf("unzip is not installed, please install it to continue")
	}
	empty := make([]*TaskResult, 0)
	return NewPlayResult(input, empty), nil
}
