// Package playbook provides Playbook types and related utilities.
package playbook

// Input holds the input parameters for executing a playbook.
type Input struct {
	DryRun   bool
	Play     string
	PlayBook *PlayBook
	Sudo     bool
}
