package external

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/alanktwong/dum/internal/external/gen"
)

// Git is the abstraction for cloning a repository.
type Git = gen.Git

// GitImpl implements Git.
type GitImpl struct {
	Utils Ext
}

// NewGit constructs GitImpl.
func NewGit() *GitImpl {
	return &GitImpl{
		Utils: NewExt(),
	}
}

// Clone implements Git.
func (g *GitImpl) Clone(ctx context.Context, url, name, path string, sudo bool) error {
	if err := g.Utils.CreateDirectory(ctx, path, sudo); err != nil {
		return fmt.Errorf("failed to create root dir: %w", err)
	}
	cmd := g.createCloneCommand(ctx, url, name, sudo)
	cmd.Dir = path
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to git clone %s: %w", url, err)
	}
	return nil
}

// AlreadyExists implements Git.
func (g *GitImpl) AlreadyExists(targetPath string) bool {
	if _, err := os.Stat(targetPath); err == nil {
		return true
	}
	return false
}

func (g *GitImpl) createCloneCommand(ctx context.Context, url, name string, sudo bool) *exec.Cmd {
	if g.isGithub(url) && g.Utils.IsInstalled("gh") {
		return g.createGithubClone(ctx, url, sudo)
	}
	return g.createGitClone(ctx, url, name, sudo)
}

func (g *GitImpl) createGitClone(ctx context.Context, url, name string, sudo bool) *exec.Cmd {
	if sudo {
		return exec.CommandContext(ctx, "sudo", "git", "clone", url, name)
	}
	return exec.CommandContext(ctx, "git", "clone", url, name)
}

func (g *GitImpl) createGithubClone(ctx context.Context, url string, sudo bool) *exec.Cmd {
	org, repo := g.githubOrgRepo(url)
	target := fmt.Sprintf("%s/%s", org, repo)
	if sudo {
		return exec.CommandContext(ctx, "sudo", "gh", "repo", "clone", target)
	}
	return exec.CommandContext(ctx, "gh", "repo", "clone", target)
}

func (g *GitImpl) isGithub(url string) bool {
	return strings.HasPrefix(url, "https://github.com")
}

func (g *GitImpl) githubOrgRepo(url string) (string, string) {
	parts := strings.Split(strings.TrimPrefix(url, "https://github.com"), "/")
	org := parts[0]
	repo := strings.TrimSuffix(parts[1], ".git")
	return org, repo
}
