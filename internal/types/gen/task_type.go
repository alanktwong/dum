// Package gen contains generated enum types.
package gen

//go:generate go run github.com/abice/go-enum@latest --file=$GOFILE --names

// TaskType is an enum source.
// ENUM(cask, cellar, brew, dir, git, link, jetbrains, mas, vscode, bash).
type TaskType string
