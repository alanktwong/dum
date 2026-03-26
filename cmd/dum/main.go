// Main package for installer
package main

import (
	"awong/dotfiles/internal/cmd"
)

func main() {
	dum := cmd.NewDum()
	dum.Exec()
}
