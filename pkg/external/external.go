// Package external contains the low-level functions interfaces and structs
// to handle the OS and network .
package external

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"os/user"
	"path/filepath"
	"runtime"
	"strings"
	"syscall"
)

// Runner is an abstraction around things such as exec.Command.
type Runner interface {
	// Run runs something in a context.
	Run(ctx context.Context) error
}

// Ext is an OS abstraction.
type Ext interface {
	// IsInstalled determines whether the 'command' is installed in the PATH.
	IsInstalled(command string) bool
	// IsOSX determines whether the current OS is a MacOSX.
	IsOSX() bool
	// IsLinux determines whether the current OS is Linux.
	IsLinux() bool
	// IsUserInFileGroup - Checks if the current user belongs to the group of the given file.
	IsUserInFileGroup(filePath string) (bool, error)
	// CreateDirectory does a mkdir -p.
	CreateDirectory(ctx context.Context, path string, sudo bool) error
	// SoftLink creates a soft link at rootPath from src to target.
	SoftLink(ctx context.Context, rootPath, src, target string, sudo bool) error
	// ExpandUser expands that path if it contains home directory aliases.
	ExpandUser(path string) (string, error)
	// ToAbsolutePath converts a path given its environment variables.
	ToAbsolutePath(path string) (string, error)
	// IsDir determines whether a path is a directory.
	IsDir(path string) bool
	// IsSymlink determines whether a path is a symbolic link.
	IsSymlink(path string) bool
	// RunCommand runs a bash command.
	RunCommand(ctx context.Context, command string, sudo bool) error
}

// DefaultExt implements Ext.
type DefaultExt struct{}

// NewExt constructs DefaultExt.
func NewExt() *DefaultExt {
	return &DefaultExt{}
}

// IsInstalled implements Ext.
func (u *DefaultExt) IsInstalled(command string) bool {
	if command == "bash" {
		return u.isBashInstalled()
	}
	if command == "ovim" {
		return u.isVimInstalled()
	}
	if command == "tar" {
		return u.checkPath(command, filepath.Join("/usr", "bin", command))
	}
	if command == "zip" {
		return u.checkPath(command, filepath.Join("/usr", "bin", command))
	}
	if command == "unzip" {
		return u.checkPath(command, filepath.Join("/usr", "bin", command))
	}
	err := u.lookPath(command)
	return err == nil
}

func (u *DefaultExt) lookPath(command string) error {
	_, err := exec.LookPath(command)
	if err != nil {
		return fmt.Errorf("failed to lookup path %s with error: %v", command, err)
	}
	return nil
}

func (u *DefaultExt) isBashInstalled() bool {
	bashPath := filepath.Join("/usr", "local", "bin", "bash")
	if _, err := os.Stat(bashPath); err == nil {
		return true
	}
	return false
}

func (u *DefaultExt) isVimInstalled() bool {
	vimPath := filepath.Join("/usr", "bin", "vim")
	if _, err := os.Stat(vimPath); err == nil {
		return true
	}
	err := u.lookPath(vimPath)
	return err == nil
}

func (u *DefaultExt) checkPath(command, filepath string) bool {
	if _, err := os.Stat(filepath); err == nil {
		return true
	}
	err := u.lookPath(command)
	return err == nil
}

// IsOSX implements Ext.
func (u *DefaultExt) IsOSX() bool {
	return runtime.GOOS == "darwin"
}

// IsLinux implements Ext.
func (u *DefaultExt) IsLinux() bool {
	return runtime.GOOS == "linux"
}

// IsUserInFileGroup implements Ext.
func (u *DefaultExt) IsUserInFileGroup(filePath string) (bool, error) {
	info, err := os.Stat(filePath)
	if err != nil {
		return false, fmt.Errorf("failed to stat file: %w", err)
	}
	stat := info.Sys().(*syscall.Stat_t)
	fileGid := stat.Gid

	currentUser, err := user.Current()
	if err != nil {
		return false, fmt.Errorf("failed to get current user: %w", err)
	}
	groups, err := currentUser.GroupIds()
	if err != nil {
		return false, fmt.Errorf("failed to get user groups: %w", err)
	}
	for _, gid := range groups {
		if fmt.Sprintf("%d", fileGid) == gid {
			return true, nil
		}
	}
	return false, fmt.Errorf("not yet implemented")
}

// CreateDirectory implements Ext.
func (u *DefaultExt) CreateDirectory(ctx context.Context, path string, sudo bool) error {
	_, err := os.Stat(path)
	if os.IsExist(err) {
		return nil
	}
	runner := func() error {
		return os.MkdirAll(path, os.ModeDir)
	}
	if sudo {
		runner = func() error {
			cmd := exec.CommandContext(ctx, "sudo", "mkdir", "-p", path)
			return cmd.Run()
		}
	}
	if err := runner(); err != nil {
		return fmt.Errorf("failed to create directory with: %w", err)
	}
	return nil
}

// SoftLink implements Ext.
func (u *DefaultExt) SoftLink(ctx context.Context, rootPath, src, target string, sudo bool) error {
	return u.pushd(rootPath, func() error {
		if sudo {
			cmd := exec.CommandContext(ctx, "sudo", "ln", "-s", src, target)
			return cmd.Run()
		}
		return os.Symlink(src, target)
	})
}

func (u *DefaultExt) pushd(newDir string, fn func() error) error {
	previousDir, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("failed to os.Getwd: %v", err)
	}
	if err := os.Chdir(newDir); err != nil {
		return fmt.Errorf("failed to os.Chdir: %v", err)
	}
	defer func() {
		_ = os.Chdir(previousDir)
	}()
	return fn()
}

// ExpandUser implements Ext.
func (u *DefaultExt) ExpandUser(path string) (string, error) {
	if strings.HasPrefix(path, "~") {
		home, err := os.UserHomeDir()
		if err != nil {
			return "", fmt.Errorf("failed to get user home directory: %w", err)
		}
		return home + path[1:], nil
	}
	return path, nil
}

// IsDir implements Ext.
func (u *DefaultExt) IsDir(path string) bool {
	info, err := os.Stat(path)
	return err == nil && info.IsDir()
}

// IsSymlink implements Ext.
func (u *DefaultExt) IsSymlink(path string) bool {
	info, err := os.Lstat(path)
	return err == nil && info.Mode()&os.ModeSymlink != 0
}

// RunCommand implements Ext.
func (u *DefaultExt) RunCommand(ctx context.Context, command string, sudo bool) error {
	runner := func() error {
		cmd := exec.CommandContext(ctx, "bash", "-c", command)
		return cmd.Run()
	}
	if sudo {
		runner = func() error {
			cmd := exec.CommandContext(ctx, "sudo", "bash", "-c", command)
			return cmd.Run()
		}
	}
	return runner()
}

// ToAbsolutePath implements Ext.
func (u *DefaultExt) ToAbsolutePath(path string) (string, error) {
	return u.absolutePath(path)
}
