package types

type TaskInput struct {
	DryRun        bool
	Task          string
	Play          string
	PlayBook      string
	Sudo          bool
	JetBrainsApps map[string]string
}

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
