// Package types provides common types used throughout the application.
package types

// TaskInput holds the input parameters for executing a task.
type TaskInput struct {
	DryRun        bool
	Task          string
	Play          string
	PlayBook      string
	Sudo          bool
	JetBrainsApps map[string]string
}

// NewTaskInput creates a new TaskInput instance.
func NewTaskInput(dryRun, sudo bool, task, play, playbook string, jetbrainsApps map[string]string) *TaskInput {
	return &TaskInput{
		DryRun:        dryRun,
		Task:          task,
		Play:          play,
		PlayBook:      playbook,
		Sudo:          sudo,
		JetBrainsApps: jetbrainsApps,
	}
}
