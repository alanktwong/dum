package playbook

type PlaybookInput struct {
	DryRun   bool
	Play     string
	PlayBook *PlayBook
	Sudo     bool
}
