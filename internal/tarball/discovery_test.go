package tarball

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDiscoverArchives_ReturnsZeroOneAndMultipleMatches(t *testing.T) {
	target := t.TempDir()
	for _, name := range []string{
		"Documents-2026-01-01-000001.tgz",
		"Documents-2026-01-02-000001.tgz",
		"Photos-2026-01-01-000001.tgz",
		"unrelated.txt",
	} {
		if err := os.WriteFile(filepath.Join(target, name), []byte("archive"), 0o600); err != nil {
			t.Fatal(err)
		}
	}

	candidates, err := NewService(nil).DiscoverArchives(context.Background(), []string{
		"/home/user/Documents",
		"/home/user/Missing",
	}, target)
	if !assert.NoError(t, err) {
		return
	}
	if !assert.Len(t, candidates, 2) {
		return
	}
	assert.Equal(t, []string{
		filepath.Join(target, "Documents-2026-01-01-000001.tgz"),
		filepath.Join(target, "Documents-2026-01-02-000001.tgz"),
	}, candidates[0].Archives)
	assert.Empty(t, candidates[1].Archives)
}

func TestDiscoverArchives_MatchesBasenamesCaseInsensitively(t *testing.T) {
	target := t.TempDir()
	archive := filepath.Join(target, "Documents-2026-01-01-000001.tgz")
	if err := os.WriteFile(archive, []byte("archive"), 0o600); err != nil {
		t.Fatal(err)
	}

	candidates, err := NewService(nil).DiscoverArchives(
		context.Background(), []string{"/home/user/documents"}, target,
	)
	if !assert.NoError(t, err) || !assert.Len(t, candidates, 1) {
		return
	}
	assert.Equal(t, []string{archive}, candidates[0].Archives)
}

func TestDiscoverArchives_IncludesArchivesWithoutSidecarsForRestoreError(t *testing.T) {
	target := t.TempDir()
	archive := filepath.Join(target, "Documents-legacy.tgz")
	if err := os.WriteFile(archive, []byte("legacy archive"), 0o600); err != nil {
		t.Fatal(err)
	}

	candidates, err := NewService(nil).DiscoverArchives(context.Background(), []string{"/home/user/Documents"}, target)
	if !assert.NoError(t, err) {
		return
	}
	if !assert.Len(t, candidates, 1) {
		return
	}
	assert.Equal(t, []string{archive}, candidates[0].Archives)
	assert.NoFileExists(t, archive+".sha256")
	assert.NoFileExists(t, archive+".manifest.yaml")
}

func TestDiscoverArchives_RejectsBasenameCollisions(t *testing.T) {
	_, err := NewService(nil).DiscoverArchives(context.Background(), []string{
		"/home/user/Shared",
		"/mnt/data/shared",
	}, t.TempDir())
	assert.Error(t, err)
}

func TestDiscoverArchives_MissingTargetReturnsNoMatches(t *testing.T) {
	target := filepath.Join(t.TempDir(), "not-created")
	candidates, err := NewService(nil).DiscoverArchives(context.Background(), []string{"/home/user/Documents"}, target)
	if !assert.NoError(t, err) {
		return
	}
	if !assert.Len(t, candidates, 1) {
		return
	}
	assert.Empty(t, candidates[0].Archives)
}

func TestDiscoverArchives_IgnoresDirectoriesAndSymlinks(t *testing.T) {
	target := t.TempDir()
	archiveName := "Documents-2026-01-01-000001.tgz"
	if err := os.Mkdir(filepath.Join(target, archiveName), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.Symlink("missing.tgz", filepath.Join(target, "Documents-2026-01-02-000001.tgz")); err != nil {
		t.Fatal(err)
	}

	candidates, err := NewService(nil).DiscoverArchives(context.Background(), []string{"/home/user/Documents"}, target)
	if !assert.NoError(t, err) {
		return
	}
	if !assert.Len(t, candidates, 1) {
		return
	}
	assert.Empty(t, candidates[0].Archives)
}
