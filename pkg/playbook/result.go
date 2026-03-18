package playbook

// Result expresses the result of a PlayBook execution.
type Result struct {
	PlayBook string
	Plays    []*PlayResult
}

// NewResult constructs a PlayBook Result.
func NewResult(input *Input, plays []*PlayResult) *Result {
	return &Result{PlayBook: input.PlayBook.ID, Plays: plays}
}
