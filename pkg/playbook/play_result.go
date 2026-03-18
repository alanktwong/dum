package playbook

// PlayResult expresses the result of a Play execution.
type PlayResult struct {
	PlayBook string
	Play     string
	Success  bool
	Tasks    []*TaskResult
}

// NewPlayResult constructs a PlayResult.
func NewPlayResult(input *Input, tasksResults []*TaskResult) *PlayResult {
	return &PlayResult{
		PlayBook: input.PlayBook.ID,
		Play:     input.Play,
		Success:  calculateSuccess(tasksResults),
	}
}

func calculateSuccess(tasksResults []*TaskResult) bool {
	for _, result := range tasksResults {
		if !result.Success {
			return false
		}
	}
	return true
}
