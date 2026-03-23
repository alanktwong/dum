// Package enums contains enum definitions.
package enums

// TaskType is an enum source.
// ENUM(cask, cellar, brew, dir, function,  git, link, jetbrains, mas, vscode, bash).
//
//go:generate ../../bin/go-enum --ptr --marshal --flag --nocase --mustparse --names --values --nocomments
type TaskType string
