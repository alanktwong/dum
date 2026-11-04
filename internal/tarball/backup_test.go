package tarball

import (
	"archive/tar"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"syscall"
	"testing"
	"time"

	gen "github.com/alanktwong/dum/internal/logging/gen"
	"github.com/stretchr/testify/assert"
	mck "github.com/stretchr/testify/mock"
)

func TestBackup_WritesArchiveAndSidecarsPerPath(t *testing.T) {
	root := t.TempDir()
	home := filepath.Join(root, "home")
	inside := filepath.Join(home, "Documents")
	external := filepath.Join(root, "external", "Photos")
	insideFile := filepath.Join(inside, "notes.txt")
	externalFile := filepath.Join(external, "image.bin")
	modified := time.Date(2024, time.April, 5, 6, 7, 8, 0, time.Local)
	writeTestFile(t, insideFile, []byte("notes"), 0o640, modified)
	writeTestFile(t, externalFile, []byte{0x01, 0x02, 0x03}, 0o600, modified)

	target := filepath.Join(root, "archives")
	service := newBackupTestService()
	result, err := service.Backup(context.Background(), BackupOptions{
		Paths:  []string{inside, external},
		Target: target,
		Home:   home,
	})
	if !assert.NoError(t, err) {
		return
	}
	assert.Len(t, result.Artifacts, 2)

	for index, name := range []string{"Documents", "Photos"} {
		artifacts := result.Artifacts[index]
		assert.Equal(t, filepath.Join(target, name+"-2026-11-03-101112.tgz"), artifacts.ArchivePath)
		assert.Equal(t, artifacts.ArchivePath+".sha256", artifacts.ChecksumPath)
		assert.Equal(t, artifacts.ArchivePath+".manifest.yaml", artifacts.ManifestPath)

		archiveBytes, readErr := os.ReadFile(artifacts.ArchivePath)
		if !assert.NoError(t, readErr) {
			return
		}
		digest := sha256.Sum256(archiveBytes)
		checksum, readErr := os.ReadFile(artifacts.ChecksumPath)
		if !assert.NoError(t, readErr) {
			return
		}
		assert.Equal(t, hex.EncodeToString(digest[:])+"\n", string(checksum))
		assert.FileExists(t, artifacts.ManifestPath)
	}

	insideEntries := readTestArchive(t, result.Artifacts[0].ArchivePath)
	assert.Equal(t, []byte("notes"), insideEntries["Documents/notes.txt"].Content)
	assert.Equal(t, int64(0o640), insideEntries["Documents/notes.txt"].Header.Mode)
	assert.True(t, insideEntries["Documents/notes.txt"].Header.ModTime.Equal(modified))
	assert.Contains(t, insideEntries, "Documents")

	externalEntries := readTestArchive(t, result.Artifacts[1].ArchivePath)
	assert.Equal(t, []byte{0x01, 0x02, 0x03}, externalEntries["Photos/image.bin"].Content)

	insideManifest := readTestManifest(t, result.Artifacts[0].ManifestPath)
	assert.Equal(t, 1, insideManifest.SchemaVersion)
	assert.Equal(t, []string{"Documents", "Documents/notes.txt"}, manifestPaths(insideManifest))
	fileEntry := manifestEntryByPath(insideManifest, "Documents/notes.txt")
	assert.NotNil(t, fileEntry.UID)
	assert.NotNil(t, fileEntry.GID)
	stat, err := os.Stat(insideFile)
	if assert.NoError(t, err) {
		owner, ok := stat.Sys().(*syscall.Stat_t)
		if assert.True(t, ok) {
			assert.Equal(t, int64(owner.Uid), *fileEntry.UID)
			assert.Equal(t, int64(owner.Gid), *fileEntry.GID)
		}
	}
}

func TestBackup_DryRunReportsArtifactsWithoutWriting(t *testing.T) {
	root := t.TempDir()
	home := filepath.Join(root, "home")
	source := filepath.Join(home, "Documents")
	writeTestFile(t, filepath.Join(source, "notes.txt"), []byte("notes"), 0o600, time.Time{})
	target := filepath.Join(root, "archives")

	service := newBackupTestService()
	result, err := service.Backup(context.Background(), BackupOptions{
		Paths:  []string{source},
		Target: target,
		Home:   home,
		DryRun: true,
	})
	if !assert.NoError(t, err) {
		return
	}
	assert.Len(t, result.Artifacts, 1)
	assert.NoDirExists(t, target)
}

func TestBackup_RejectsBasenameCollisionBeforeWriting(t *testing.T) {
	root := t.TempDir()
	home := filepath.Join(root, "home")
	first := filepath.Join(home, "A", "Shared")
	second := filepath.Join(root, "external", "shared")
	writeTestFile(t, filepath.Join(first, "first.txt"), []byte("first"), 0o600, time.Time{})
	writeTestFile(t, filepath.Join(second, "second.txt"), []byte("second"), 0o600, time.Time{})
	target := filepath.Join(root, "archives")

	service := newBackupTestService()
	_, err := service.Backup(context.Background(), BackupOptions{
		Paths:  []string{first, second},
		Target: target,
		Home:   home,
	})
	assert.Error(t, err)
	assert.NoDirExists(t, target)
}

func TestBackup_RejectsTargetInsideSourceBeforeWriting(t *testing.T) {
	root := t.TempDir()
	home := filepath.Join(root, "home")
	source := filepath.Join(home, "Documents")
	writeTestFile(t, filepath.Join(source, "notes.txt"), []byte("notes"), 0o600, time.Time{})
	target := filepath.Join(source, "backups")

	service := newBackupTestService()
	_, err := service.Backup(context.Background(), BackupOptions{
		Paths:  []string{source},
		Target: target,
		Home:   home,
	})
	assert.Error(t, err)
	assert.NoDirExists(t, target)
}

func TestBackup_DoesNotOverwriteExistingArtifacts(t *testing.T) {
	root := t.TempDir()
	home := filepath.Join(root, "home")
	source := filepath.Join(home, "Documents")
	writeTestFile(t, filepath.Join(source, "notes.txt"), []byte("notes"), 0o600, time.Time{})
	target := filepath.Join(root, "archives")
	service := newBackupTestService()

	first, err := service.Backup(context.Background(), BackupOptions{
		Paths:  []string{source},
		Target: target,
		Home:   home,
	})
	if !assert.NoError(t, err) {
		return
	}
	before, err := os.ReadFile(first.Artifacts[0].ArchivePath)
	if !assert.NoError(t, err) {
		return
	}

	_, err = service.Backup(context.Background(), BackupOptions{
		Paths:  []string{source},
		Target: target,
		Home:   home,
	})
	assert.Error(t, err)
	after, readErr := os.ReadFile(first.Artifacts[0].ArchivePath)
	if assert.NoError(t, readErr) {
		assert.Equal(t, before, after)
	}
}

func TestBackup_PreflightsAllSourcesBeforeWriting(t *testing.T) {
	root := t.TempDir()
	home := filepath.Join(root, "home")
	first := filepath.Join(home, "First")
	missing := filepath.Join(home, "Missing")
	writeTestFile(t, filepath.Join(first, "file.txt"), []byte("first"), 0o600, time.Time{})
	target := filepath.Join(root, "archives")

	service := newBackupTestService()
	_, err := service.Backup(context.Background(), BackupOptions{
		Paths:  []string{first, missing},
		Target: target,
		Home:   home,
	})
	assert.Error(t, err)
	assert.NoDirExists(t, target)
}

func TestBackup_LogsAndRecordsUnsafeSymlink(t *testing.T) {
	root := t.TempDir()
	home := filepath.Join(root, "home")
	source := filepath.Join(home, "Documents")
	writeTestFile(t, filepath.Join(source, "settings"), []byte("data"), 0o600, time.Time{})
	if err := os.Symlink("settings", filepath.Join(source, "relative-link")); err != nil {
		t.Fatal(err)
	}
	if err := os.Symlink("/etc/passwd", filepath.Join(source, "absolute-link")); err != nil {
		t.Fatal(err)
	}

	logger := gen.NewMockLogger(t)
	logger.On("Warn", mck.Anything, mck.Anything).Return().Once()
	service := newBackupTestService()
	service.Log = logger
	target := filepath.Join(root, "archives")
	result, err := service.Backup(context.Background(), BackupOptions{
		Paths:  []string{source},
		Target: target,
		Home:   home,
	})
	if !assert.NoError(t, err) {
		return
	}

	entries := readTestArchive(t, result.Artifacts[0].ArchivePath)
	assert.Equal(t, byte(tar.TypeSymlink), entries["Documents/absolute-link"].Header.Typeflag)
	assert.Equal(t, "/etc/passwd", entries["Documents/absolute-link"].Header.Linkname)
	assert.Equal(t, byte(tar.TypeSymlink), entries["Documents/relative-link"].Header.Typeflag)
	assert.Equal(t, "settings", entries["Documents/relative-link"].Header.Linkname)

	manifest := readTestManifest(t, result.Artifacts[0].ManifestPath)
	absoluteLink := manifestEntryByPath(manifest, "Documents/absolute-link")
	assert.True(t, absoluteLink.Unsafe)
	assert.Equal(t, "/etc/passwd", absoluteLink.Target)
	assert.Contains(t, absoluteLink.Warning, "absolute symlink target")
	relativeLink := manifestEntryByPath(manifest, "Documents/relative-link")
	assert.False(t, relativeLink.Unsafe)
	assert.Empty(t, relativeLink.Warning)
}

func TestBackup_CleansPublishedSidecarsAfterPublishFailure(t *testing.T) {
	root := t.TempDir()
	home := filepath.Join(root, "home")
	source := filepath.Join(home, "Documents")
	writeTestFile(t, filepath.Join(source, "notes.txt"), []byte("notes"), 0o600, time.Time{})
	target := filepath.Join(root, "archives")
	service := newBackupTestService()
	service.publish = func(sourcePath, targetPath string) error {
		if strings.HasSuffix(targetPath, ".manifest.yaml") {
			return errors.New("injected publication failure")
		}
		return os.Link(sourcePath, targetPath)
	}

	_, err := service.Backup(context.Background(), BackupOptions{
		Paths:  []string{source},
		Target: target,
		Home:   home,
	})
	assert.Error(t, err)
	entries, readErr := os.ReadDir(target)
	if assert.NoError(t, readErr) {
		assert.Empty(t, entries)
	}
}

func newBackupTestService() *Service {
	service := NewService(nil)
	service.now = func() time.Time {
		return time.Date(2026, time.November, 3, 10, 11, 12, 0, time.Local)
	}
	return service
}

func manifestPaths(manifest Manifest) []string {
	paths := make([]string, 0, len(manifest.Entries))
	for _, entry := range manifest.Entries {
		paths = append(paths, entry.Path)
	}
	return paths
}

func manifestEntryByPath(manifest Manifest, memberPath string) ManifestEntry {
	for _, entry := range manifest.Entries {
		if entry.Path == memberPath {
			return entry
		}
	}
	return ManifestEntry{}
}
