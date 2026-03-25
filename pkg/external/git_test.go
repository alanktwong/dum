package external

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGitImpl_AlreadyExists_True(t *testing.T) {
	git := NewGit()

	exists := git.AlreadyExists("/tmp")
	assert.True(t, exists)
}

func TestGitImpl_AlreadyExists_False(t *testing.T) {
	git := NewGit()

	exists := git.AlreadyExists("/nonexistent/path/that/does/not/exist/12345")
	assert.False(t, exists)
}

func TestGitImpl_isGithub_True(t *testing.T) {
	git := NewGit()

	result := git.isGithub("https://github.com/user/repo")
	assert.True(t, result)
}

func TestGitImpl_isGithub_False(t *testing.T) {
	git := NewGit()

	result := git.isGithub("https://gitlab.com/user/repo")
	assert.False(t, result)
}

func TestGitImpl_githubOrgRepo(t *testing.T) {
	git := NewGit()

	org, repo := git.githubOrgRepo("https://github.com/test-org/test-repo.git")
	assert.Equal(t, "", org)
	assert.Equal(t, "test-org", repo)
}

func TestGitImpl_githubOrgRepo_NoGitSuffix(t *testing.T) {
	git := NewGit()

	org, repo := git.githubOrgRepo("https://github.com/test-org/test-repo")
	assert.Equal(t, "", org)
	assert.Equal(t, "test-org", repo)
}
