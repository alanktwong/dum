// Package types provides common types used throughout the application.
package types

// TaskResult holds the result of executing a task.
type TaskResult struct {
	PlayBook string
	Play     string
	Task     string
	Success  bool
	DryRun   bool
}

// NewTaskResult creates a new TaskResult instance.
func NewTaskResult(input *TaskInput, success bool) (*TaskResult, error) {
	return &TaskResult{
		PlayBook: input.PlayBook,
		Play:     input.Play,
		Task:     input.Task,
		Success:  success,
		DryRun:   input.DryRun,
	}, nil
}
