package playbook

import (
	ext "alanktwong/dum/internal/external"
	l "alanktwong/dum/internal/logging"
	pl "alanktwong/dum/internal/plays"
	"context"
	"fmt"

	om "github.com/elliotchance/orderedmap/v3"
)

// Executor performs the work of executing a play or task.
type Executor struct {
	Log          l.Logger
	Ext          ext.Ext
	PlayExecutor pl.PlayExec
}

// NewExecutor constructs an Executor.
func NewExecutor() *Executor {
	return &Executor{
		Log:          l.Log(),
		Ext:          ext.NewExt(),
		PlayExecutor: pl.NewPlayExecutor(),
	}
}

// Install installs a playbook.
func (e *Executor) Install(ctx context.Context, input *Input) (*Result, error) {
	group := input.Play
	pb := input.PlayBook
	err := e.Log.Printlnf(Ellipsis+"Installing playbook (%v) ... %v", pb.ID, pb.Description)
	if err != nil {
		return nil, fmt.Errorf("failed to install playbook %w", err)
	}
	var playResults []*pl.PlayResult
	initResult, err := e.PlayExecutor.Initialize(&pl.PlayInput{
		DryRun:   input.DryRun,
		Play:     "initialize",
		PlayBook: pb,
		Sudo:     input.Sudo,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to Initialize %w", err)
	}
	playResults = append(playResults, initResult)

	// Dry runs include disabled plays so the output shows the full playbook manifest.
	playMap := pb.GetPlays(!input.DryRun)
	if group != "" {
		results, err := e.installGroupPlay(ctx, group, playResults, playMap, input)
		if err != nil {
			return nil, fmt.Errorf("failed to install group %w", err)
		}
		playResults = results
	} else {
		results, err := e.installAllPlays(ctx, playResults, playMap, input)
		if err != nil {
			return nil, fmt.Errorf("failed to install group %w", err)
		}
		playResults = results
	}
	e.Log.Infof(Ellipsis+"END: installing playbook (%v) ... %v", pb.ID, pb.Description)
	success := pl.CalculateSuccess(playResults)
	return NewResult(input, success), nil
}

func (e *Executor) installGroupPlay(ctx context.Context,
	group string,
	playResults []*pl.PlayResult,
	playMap *om.OrderedMap[string, *pl.Play],
	input *Input,
) ([]*pl.PlayResult, error) {
	if play, ok := playMap.Get(group); ok {
		playResult, err := e.PlayExecutor.InstallPlay(ctx, play, &pl.PlayInput{
			DryRun:   input.DryRun,
			Play:     play.ID,
			PlayBook: input.PlayBook,
			Sudo:     input.Sudo,
		})
		if err != nil {
			return nil, fmt.Errorf("failed to install play %s: %w", play.GetID(), err)
		}
		playResults = append(playResults, playResult)
	}
	return playResults, nil
}

func (e *Executor) installAllPlays(ctx context.Context,
	playResults []*pl.PlayResult,
	playMap *om.OrderedMap[string, *pl.Play],
	input *Input,
) ([]*pl.PlayResult, error) {
	for current, play := range playMap.AllFromFront() {
		e.Log.Infof("%v %v", PlayEllipsis, current)
		playResult, err := e.PlayExecutor.InstallPlay(ctx, play, &pl.PlayInput{
			DryRun:   input.DryRun,
			Play:     play.ID,
			PlayBook: input.PlayBook,
			Sudo:     input.Sudo,
		})
		if err != nil {
			return nil, fmt.Errorf("failed to install play %s: %w", play.GetID(), err)
		}
		playResults = append(playResults, playResult)
	}
	return playResults, nil
}
