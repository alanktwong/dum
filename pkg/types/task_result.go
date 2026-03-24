package types

type TaskResult struct {
	PlayBook string
	Play     string
	Task     string
	Success  bool
	DryRun   bool
}

func NewTaskResult(input *TaskInput, success bool) (*TaskResult, error) {
	return &TaskResult{
		PlayBook: input.PlayBook,
		Play:     input.Play,
		Task:     input.Task,
		Success:  success,
		DryRun:   input.DryRun,
	}, nil
}
