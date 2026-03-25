// Package plays provides Play types and related utilities.
package plays

import "awong/dotfiles/pkg/plays/gen"

// PlayInput holds the input parameters for executing a play.
type PlayInput struct {
	DryRun   bool
	Play     string
	PlayBook gen.PlayBookInfo
	Sudo     bool
}
