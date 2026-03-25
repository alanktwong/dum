// Package gen provides generated types for external package.
package gen

import (
	"context"
)

// Runner is an abstraction around things such as exec.Command.
type Runner interface {
	Run(ctx context.Context) error
}

// Ext is an OS abstraction.
type Ext interface {
	IsInstalled(command string) bool
	IsOSX() bool
	IsLinux() bool
	IsUserInFileGroup(filePath string) (bool, error)
	CreateDirectory(ctx context.Context, path string, sudo bool) error
	SoftLink(ctx context.Context, rootPath, src, target string, sudo bool) error
	ExpandUser(path string) (string, error)
	ToAbsolutePath(path string) (string, error)
	IsDir(path string) bool
	IsSymlink(path string) bool
	RunCommand(ctx context.Context, command string, sudo bool) error
	GetString(data map[string]any, key string, def string) string
	GetStrings(data map[string]any, key string, def []string) []string
	GetBool(data map[string]any, key string, def bool) bool
}

// Brew is an abstraction for homebrew.
type Brew interface {
	Install(ctx context.Context, formula string) error
	InstallCask(ctx context.Context, formula string) error
	Tap(ctx context.Context, formula string) error
	Prefix(ctx context.Context) (string, error)
	InPath(ctx context.Context, prefix, id string) bool
}

// Git is the abstraction for cloning a repository.
type Git interface {
	Clone(ctx context.Context, url, name, path string, sudo bool) error
	AlreadyExists(targetPath string) bool
}

// Code is an abstraction for vscode.
type Code interface {
	InstallExtension(ctx context.Context, formula string) error
	ListExtensions(ctx context.Context) (string, error)
}

// Mas is an abstraction for Apple app store.
type Mas interface {
	Install(ctx context.Context, app string) error
	List(ctx context.Context) (string, error)
}

// JetBrainsApp is an abstraction for jetbrains app.
type JetBrainsApp interface {
	Install(ctx context.Context, app, plugin string) error
	IsInstalled(ideName, plugin string) bool
}
