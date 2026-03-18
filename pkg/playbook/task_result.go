package playbook

// TaskResult expresses the result of executing a task.
type TaskResult struct {
	PlayBook string
	Play     string
	Task     string
	Success  bool
	DryRun   bool
}

// NewTaskResult constructs a TaskResult.
func NewTaskResult(input *Input, success bool) (*TaskResult, error) {
	return &TaskResult{
		PlayBook: input.PlayBook.ID,
		Play:     input.Play,
		Task:     input.Task,
		Success:  success,
		DryRun:   input.DryRun,
	}, nil
}

// CreateTaskResult constructs a TaskResult from an attribute.
func (a *Attributes) CreateTaskResult(input *Input, success bool) (*TaskResult, error) {
	return &TaskResult{
		Task:     a.ID,
		Play:     input.Play,
		PlayBook: input.PlayBook.ID,
		DryRun:   input.DryRun,
		Success:  success,
	}, nil
}
