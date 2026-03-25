package types

// JetBrainsType is an enum source.
// ENUM(clion, datagrip, goland, idea, phpstorm, pycharm, rider, rubymine, rustrover, webstorm).
//
//go:generate ../../bin/go-enum --ptr --marshal --flag --nocase --mustparse --names --values --nocomments
type JetBrainsType string
