// Package plays provides Play types and related utilities.
package plays

import pg "awong/dotfiles/internal/plays/gen"

// PlayInput holds the input parameters for executing a play.
type PlayInput struct {
	DryRun   bool
	Play     string
	PlayBook pg.PlayBookInfo
	Sudo     bool
}
