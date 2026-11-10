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
	assert.Equal(t, []string{"Documents"}, manifestPaths(insideManifest))
	assert.NotContains(t, manifestPaths(insideManifest), "Documents/notes.txt")
	stat, err := os.Stat(insideFile)
	if assert.NoError(t, err) {
		owner, ok := stat.Sys().(*syscall.Stat_t)
		if assert.True(t, ok) {
			assert.Equal(t, int(owner.Uid), insideEntries["Documents/notes.txt"].Header.Uid)
			assert.Equal(t, int(owner.Gid), insideEntries["Documents/notes.txt"].Header.Gid)
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

func TestBackup_CreatesTargetWithRestrictedPermissions(t *testing.T) {
	root := t.TempDir()
	home := filepath.Join(root, "home")
	source := filepath.Join(home, "Documents")
	writeTestFile(t, filepath.Join(source, "file.txt"), []byte("data"), 0o600, time.Time{})
	target := filepath.Join(root, "archives")

	_, err := newBackupTestService().Backup(context.Background(), BackupOptions{
		Paths: []string{source}, Target: target, Home: home,
	})
	if !assert.NoError(t, err) {
		return
	}
	info, err := os.Stat(target)
	if !assert.NoError(t, err) {
		return
	}
	permissions := info.Mode().Perm()
	assert.Equal(t, os.FileMode(0o700), permissions&0o700)
	assert.Zero(t, permissions&^os.FileMode(0o750))
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

func TestBackup_ExcludesFilesByPattern(t *testing.T) {
	t.Parallel()
	root := t.TempDir()
	home := filepath.Join(root, "home")
	source := filepath.Join(home, "Project")
	writeTestFile(t, filepath.Join(source, "main.go"), []byte("code"), 0o644, time.Now())
	writeTestFile(t, filepath.Join(source, "debug.log"), []byte("log"), 0o644, time.Now())
	writeTestFile(t, filepath.Join(source, "notes.txt"), []byte("notes"), 0o644, time.Now())

	target := filepath.Join(root, "archives")
	service := newBackupTestService()
	result, err := service.Backup(context.Background(), BackupOptions{
		Paths:        []string{source},
		Exemptions:   []string{"*.log"},
		UseGitignore: false,
		Target:       target,
		Home:         home,
	})
	if !assert.NoError(t, err) {
		return
	}
	assert.Len(t, result.Artifacts, 1)

	entries := readTestArchive(t, result.Artifacts[0].ArchivePath)
	assert.Contains(t, entries, "Project/main.go")
	assert.Contains(t, entries, "Project/notes.txt")
	assert.NotContains(t, entries, "Project/debug.log")

	manifest := readTestManifest(t, result.Artifacts[0].ManifestPath)
	assert.Equal(t, []string{"*.log"}, manifest.Exemptions)
	assert.False(t, manifest.UseGitignore)
	manifestPathSet := make(map[string]struct{}, len(manifest.Entries))
	for _, entry := range manifest.Entries {
		manifestPathSet[entry.Path] = struct{}{}
	}
	assert.Equal(t, map[string]struct{}{"Project": {}}, manifestPathSet)
	assert.NotContains(t, manifestPathSet, "Project/main.go")
	assert.NotContains(t, manifestPathSet, "Project/notes.txt")
	assert.NotContains(t, manifestPathSet, "Project/debug.log")
}

func TestBackup_ManifestListsDirectoriesAndSymlinks(t *testing.T) {
	t.Parallel()
	root := t.TempDir()
	home := filepath.Join(root, "home")
	source := filepath.Join(home, "Project")
	writeTestFile(t, filepath.Join(source, "notes.txt"), []byte("notes"), 0o644, time.Now())
	writeTestFile(t, filepath.Join(source, "sub", "readme.txt"), []byte("readme"), 0o644, time.Now())
	if err := os.Mkdir(filepath.Join(source, "empty"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.Symlink("readme.txt", filepath.Join(source, "sub", "readme-link")); err != nil {
		t.Fatal(err)
	}

	target := filepath.Join(root, "archives")
	result, err := newBackupTestService().Backup(context.Background(), BackupOptions{
		Paths:  []string{source},
		Target: target,
		Home:   home,
	})
	if !assert.NoError(t, err) {
		return
	}
	if !assert.Len(t, result.Artifacts, 1) {
		return
	}

	archiveEntries := readTestArchive(t, result.Artifacts[0].ArchivePath)
	assert.Contains(t, archiveEntries, "Project/notes.txt")
	assert.Contains(t, archiveEntries, "Project/sub/readme.txt")
	assert.Contains(t, archiveEntries, "Project/empty")
	assert.Contains(t, archiveEntries, "Project/sub/readme-link")

	manifest := readTestManifest(t, result.Artifacts[0].ManifestPath)
	assert.Equal(t, manifestSchemaVersion, manifest.SchemaVersion)
	assert.Equal(t, []string{
		"Project",
		"Project/empty",
		"Project/sub",
		"Project/sub/readme-link",
	}, manifestPaths(manifest))
}

func TestBackup_RegularFileOwnerMetadataMarker(t *testing.T) {
	t.Parallel()
	for _, test := range []struct {
		name       string
		uid        int64
		gid        int64
		available  bool
		wantMarker string
	}{
		{name: "unavailable", available: false, wantMarker: "unavailable"},
		{name: "known root owner", uid: 0, gid: 0, available: true},
	} {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			root := t.TempDir()
			home := filepath.Join(root, "home")
			source := filepath.Join(home, "Project")
			writeTestFile(t, filepath.Join(source, "notes.txt"), []byte("notes"), 0o644, time.Now())

			service := newBackupTestService()
			service.ownerOf = func(os.FileInfo) (int64, int64, bool) {
				return test.uid, test.gid, test.available
			}
			result, err := service.Backup(context.Background(), BackupOptions{
				Paths:  []string{source},
				Target: filepath.Join(root, "archives"),
				Home:   home,
			})
			if !assert.NoError(t, err) || !assert.Len(t, result.Artifacts, 1) {
				return
			}

			entry := readTestArchive(t, result.Artifacts[0].ArchivePath)["Project/notes.txt"]
			assert.Equal(t, int(test.uid), entry.Header.Uid)
			assert.Equal(t, int(test.gid), entry.Header.Gid)
			assert.Equal(t, test.wantMarker, entry.Header.PAXRecords["dum.owner-metadata"])
		})
	}
}

func TestBackup_ExemptsDirectoryAndPrunesChildren(t *testing.T) {
	t.Parallel()
	root := t.TempDir()
	home := filepath.Join(root, "home")
	source := filepath.Join(home, "Project")
	writeTestFile(t, filepath.Join(source, "main.go"), []byte("code"), 0o644, time.Now())
	writeTestFile(t, filepath.Join(source, "node_modules", "dep.js"), []byte("dep"), 0o644, time.Now())
	writeTestFile(t, filepath.Join(source, "node_modules", "nested", "deep.js"), []byte("deep"), 0o644, time.Now())

	target := filepath.Join(root, "archives")
	service := newBackupTestService()
	result, err := service.Backup(context.Background(), BackupOptions{
		Paths:        []string{source},
		Exemptions:   []string{"node_modules"},
		UseGitignore: false,
		Target:       target,
		Home:         home,
	})
	if !assert.NoError(t, err) {
		return
	}
	assert.Len(t, result.Artifacts, 1)

	entries := readTestArchive(t, result.Artifacts[0].ArchivePath)
	assert.Contains(t, entries, "Project/main.go")
	assert.NotContains(t, entries, "Project/node_modules")
	assert.NotContains(t, entries, "Project/node_modules/dep.js")
	assert.NotContains(t, entries, "Project/node_modules/nested/deep.js")
}

func TestBackup_UsesGitignore(t *testing.T) {
	t.Parallel()
	root := t.TempDir()
	home := filepath.Join(root, "home")
	source := filepath.Join(home, "Project")
	writeTestFile(t, filepath.Join(source, ".gitignore"), []byte("*.log\nnode_modules\n"), 0o644, time.Now())
	writeTestFile(t, filepath.Join(source, "main.go"), []byte("code"), 0o644, time.Now())
	writeTestFile(t, filepath.Join(source, "debug.log"), []byte("log"), 0o644, time.Now())
	writeTestFile(t, filepath.Join(source, "node_modules", "dep.js"), []byte("dep"), 0o644, time.Now())

	target := filepath.Join(root, "archives")
	service := newBackupTestService()
	result, err := service.Backup(context.Background(), BackupOptions{
		Paths:        []string{source},
		Exemptions:   []string{},
		UseGitignore: true,
		Target:       target,
		Home:         home,
	})
	if !assert.NoError(t, err) {
		return
	}
	assert.Len(t, result.Artifacts, 1)

	entries := readTestArchive(t, result.Artifacts[0].ArchivePath)
	assert.Contains(t, entries, "Project/.gitignore")
	assert.Contains(t, entries, "Project/main.go")
	assert.NotContains(t, entries, "Project/debug.log")
	assert.NotContains(t, entries, "Project/node_modules")
	assert.NotContains(t, entries, "Project/node_modules/dep.js")

	manifest := readTestManifest(t, result.Artifacts[0].ManifestPath)
	assert.True(t, manifest.UseGitignore)
	assert.Contains(t, manifest.Exemptions, "*.log")
	assert.Contains(t, manifest.Exemptions, "node_modules")
}

func TestBackup_UsesNestedGitignoreFilesInScope(t *testing.T) {
	t.Parallel()
	root := t.TempDir()
	home := filepath.Join(root, "home")
	source := filepath.Join(home, "Project")
	writeTestFile(t, filepath.Join(source, ".gitignore"), []byte("node_modules/\n"), 0o644, time.Now())
	writeTestFile(t, filepath.Join(source, "src", ".gitignore"), []byte("*.log\ncache/\n"), 0o644, time.Now())
	writeTestFile(t, filepath.Join(source, "src", "main.log"), []byte("ignored"), 0o644, time.Now())
	writeTestFile(t, filepath.Join(source, "src", "keep.txt"), []byte("kept"), 0o644, time.Now())
	writeTestFile(t, filepath.Join(source, "src", "cache", "secret.txt"), []byte("ignored"), 0o644, time.Now())
	writeTestFile(t, filepath.Join(source, "other", "debug.log"), []byte("kept"), 0o644, time.Now())
	writeTestFile(t, filepath.Join(source, "node_modules", "pkg", "index.js"), []byte("ignored"), 0o644, time.Now())

	result, err := newBackupTestService().Backup(context.Background(), BackupOptions{
		Paths:        []string{source},
		UseGitignore: true,
		Target:       filepath.Join(root, "archives"),
		Home:         home,
	})
	if !assert.NoError(t, err) || !assert.Len(t, result.Artifacts, 1) {
		return
	}

	entries := readTestArchive(t, result.Artifacts[0].ArchivePath)
	assert.NotContains(t, entries, "Project/src/main.log")
	assert.NotContains(t, entries, "Project/src/cache")
	assert.NotContains(t, entries, "Project/src/cache/secret.txt")
	assert.NotContains(t, entries, "Project/node_modules")
	assert.NotContains(t, entries, "Project/node_modules/pkg/index.js")
	assert.Contains(t, entries, "Project/src/keep.txt")
	assert.Contains(t, entries, "Project/other/debug.log")

	manifest := readTestManifest(t, result.Artifacts[0].ManifestPath)
	assert.True(t, manifest.UseGitignore)
	assert.Contains(t, manifest.Exemptions, "node_modules/")
	assert.Contains(t, manifest.Exemptions, "*.log")
	assert.Contains(t, manifest.Exemptions, "cache/")
}

func TestBackup_GitignoreRulesArePerSource(t *testing.T) {
	t.Parallel()
	root := t.TempDir()
	home := filepath.Join(root, "home")
	sourceA := filepath.Join(home, "Alpha")
	sourceB := filepath.Join(home, "Beta")
	writeTestFile(t, filepath.Join(sourceA, ".gitignore"), []byte("*.log\n"), 0o644, time.Now())
	writeTestFile(t, filepath.Join(sourceA, "ignored.log"), []byte("ignored"), 0o644, time.Now())
	writeTestFile(t, filepath.Join(sourceA, "kept.tmp"), []byte("kept"), 0o644, time.Now())
	writeTestFile(t, filepath.Join(sourceB, ".gitignore"), []byte("*.tmp\n"), 0o644, time.Now())
	writeTestFile(t, filepath.Join(sourceB, "kept.log"), []byte("kept"), 0o644, time.Now())
	writeTestFile(t, filepath.Join(sourceB, "ignored.tmp"), []byte("ignored"), 0o644, time.Now())

	result, err := newBackupTestService().Backup(context.Background(), BackupOptions{
		Paths:        []string{sourceA, sourceB},
		UseGitignore: true,
		Target:       filepath.Join(root, "archives"),
		Home:         home,
	})
	if !assert.NoError(t, err) || !assert.Len(t, result.Artifacts, 2) {
		return
	}

	alphaEntries := readTestArchive(t, result.Artifacts[0].ArchivePath)
	assert.NotContains(t, alphaEntries, "Alpha/ignored.log")
	assert.Contains(t, alphaEntries, "Alpha/kept.tmp")
	alphaManifest := readTestManifest(t, result.Artifacts[0].ManifestPath)
	assert.Equal(t, []string{"*.log"}, alphaManifest.Exemptions)

	betaEntries := readTestArchive(t, result.Artifacts[1].ArchivePath)
	assert.Contains(t, betaEntries, "Beta/kept.log")
	assert.NotContains(t, betaEntries, "Beta/ignored.tmp")
	betaManifest := readTestManifest(t, result.Artifacts[1].ManifestPath)
	assert.Equal(t, []string{"*.tmp"}, betaManifest.Exemptions)
}

func TestBackup_InvalidExcludeFailsBeforeArtifacts(t *testing.T) {
	t.Parallel()
	root := t.TempDir()
	home := filepath.Join(root, "home")
	source := filepath.Join(home, "Project")
	writeTestFile(t, filepath.Join(source, "main.go"), []byte("code"), 0o644, time.Now())

	target := filepath.Join(root, "archives")
	service := newBackupTestService()
	_, err := service.Backup(context.Background(), BackupOptions{
		Paths:        []string{source},
		Exemptions:   []string{"[unclosed"},
		UseGitignore: false,
		Target:       target,
		Home:         home,
	})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid exclude pattern")
	assert.NoDirExists(t, target)
}

func TestBackup_ExemptionsRecordedInManifest(t *testing.T) {
	t.Parallel()
	root := t.TempDir()
	home := filepath.Join(root, "home")
	source := filepath.Join(home, "Project")
	writeTestFile(t, filepath.Join(source, "main.go"), []byte("code"), 0o644, time.Now())

	target := filepath.Join(root, "archives")
	service := newBackupTestService()
	result, err := service.Backup(context.Background(), BackupOptions{
		Paths:        []string{source},
		Exemptions:   []string{"dist/", "*.tmp"},
		UseGitignore: true,
		Target:       target,
		Home:         home,
	})
	if !assert.NoError(t, err) {
		return
	}
	assert.Len(t, result.Artifacts, 1)

	manifest := readTestManifest(t, result.Artifacts[0].ManifestPath)
	assert.Equal(t, []string{"dist/", "*.tmp"}, manifest.Exemptions)
	assert.True(t, manifest.UseGitignore)
}
