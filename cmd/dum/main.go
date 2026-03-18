// Main package for installer
package main

import (
	"awong/dotfiles/pkg/cmd"
)

func main() {
	dum := cmd.NewDum()
	dum.Exec()
}
