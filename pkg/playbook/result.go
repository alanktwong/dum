package playbook

// Result expresses the result of a PlayBook execution.
type Result struct {
	PlayBook string
	Success  bool
}

// NewResult constructs a PlayBook Result.
func NewResult(input *Input, success bool) *Result {
	return &Result{PlayBook: input.PlayBook.ID, Success: success}
}
