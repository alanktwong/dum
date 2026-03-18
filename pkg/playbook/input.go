package playbook

// Input holds the current input of an execution of a play or task.
type Input struct {
	DryRun   bool
	Task     string
	Play     string
	PlayBook *PlayBook
	Sudo     bool
}

// NewInput constructs an Input.
func NewInput(dryRun bool, play string, playBook *PlayBook) *Input {
	return &Input{
		DryRun:   dryRun,
		Task:     "",
		Play:     play,
		PlayBook: playBook,
		Sudo:     playBook.Sudo,
	}
}
