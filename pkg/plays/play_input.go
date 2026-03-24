// Package plays provides Play types and related utilities.
package plays

// PlayInput holds the input parameters for executing a play.
type PlayInput struct {
	DryRun   bool
	Play     string
	PlayBook PlayBookInfo
	Sudo     bool
}
