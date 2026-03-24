package plays

import (
	"awong/dotfiles/pkg/types"
)

// PlayResult expresses the result of a Play execution.
type PlayResult struct {
	PlayBook string
	Play     string
	Success  bool
	Tasks    []*types.TaskResult
}

// NewPlayResult constructs a PlayResult.
func NewPlayResult(input *PlayInput, tasksResults []*types.TaskResult) *PlayResult {
	return &PlayResult{
		PlayBook: input.PlayBook.GetID(),
		Play:     input.Play,
		Success:  calculateSuccess(tasksResults),
	}
}

func calculateSuccess(tasksResults []*types.TaskResult) bool {
	for _, result := range tasksResults {
		if !result.Success {
			return false
		}
	}
	return true
}

// CalculateSuccess determines if all plays were successful.
func CalculateSuccess(playResults []*PlayResult) bool {
	for _, result := range playResults {
		if !result.Success {
			return false
		}
	}
	return true
}
