package plays

type PlayInput struct {
	DryRun   bool
	Play     string
	PlayBook PlayBookInfo
	Sudo     bool
}
