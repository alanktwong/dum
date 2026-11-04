package tarball

import (
	"archive/tar"
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStageRestore_VerifiesAndStagesValidArchive(t *testing.T) {
	root := t.TempDir()
	archivePath := filepath.Join(root, "Documents-2026-11-03-101112.tgz")
	artifacts := writeValidRestoreFixture(t, archivePath)
	destination := filepath.Join(root, "destination")

	staged, err := NewService(nil).stageRestore(context.Background(), RestoreOptions{
		Archive:     artifacts.ArchivePath,
		Destination: destination,
	})
	if !assert.NoError(t, err) {
		return
	}
	defer func() { _ = staged.cleanup() }()

	content, err := os.ReadFile(filepath.Join(staged.StageDir, "Documents", "file.txt"))
	if !assert.NoError(t, err) {
		return
	}
	assert.Equal(t, []byte("payload"), content)
	assert.NoDirExists(t, destination)
}

func TestStageRestore_RequiresBothSidecars(t *testing.T) {
	root := t.TempDir()
	archivePath := filepath.Join(root, "legacy.tgz")
	if err := os.WriteFile(archivePath, []byte("archive"), 0o600); err != nil {
		t.Fatal(err)
	}
	if _, err := NewService(nil).stageRestore(context.Background(), RestoreOptions{Archive: archivePath}); err == nil {
		t.Fatal("expected missing sidecars to fail")
	}
}

func TestStageRestore_RejectsMalformedAndMismatchedChecksums(t *testing.T) {
	t.Run("malformed checksum", func(t *testing.T) {
		root := t.TempDir()
		artifacts := writeValidRestoreFixture(t, filepath.Join(root, "archive.tgz"))
		if err := os.WriteFile(artifacts.ChecksumPath, []byte("not-a-digest\n"), 0o600); err != nil {
			t.Fatal(err)
		}
		assertStageRestoreFails(t, artifacts.ArchivePath)
	})

	t.Run("checksum mismatch", func(t *testing.T) {
		root := t.TempDir()
		artifacts := writeValidRestoreFixture(t, filepath.Join(root, "archive.tgz"))
		if err := os.WriteFile(artifacts.ArchivePath, []byte("changed archive"), 0o600); err != nil {
			t.Fatal(err)
		}
		assertStageRestoreFails(t, artifacts.ArchivePath)
	})
}

func TestStageRestore_RejectsUnsupportedManifestVersion(t *testing.T) {
	root := t.TempDir()
	artifacts := writeValidRestoreFixture(t, filepath.Join(root, "archive.tgz"))
	data := []byte("schema_version: 2\nentries: []\n")
	if err := os.WriteFile(artifacts.ManifestPath, data, 0o600); err != nil {
		t.Fatal(err)
	}
	assertStageRestoreFails(t, artifacts.ArchivePath)
}

func TestStageRestore_RejectsManifestArchiveMismatches(t *testing.T) {
	tests := []struct {
		name       string
		header     tar.Header
		manifest   ManifestEntry
		manifestID string
	}{
		{
			name: "path mismatch",
			header: tar.Header{
				Name: "Documents/file.txt", Typeflag: tar.TypeReg, Mode: 0o644, Size: 7,
			},
			manifest:   ManifestEntry{Path: "Documents/other.txt", Type: EntryTypeFile},
			manifestID: "path-mismatch",
		},
		{
			name: "type mismatch",
			header: tar.Header{
				Name: "Documents/file.txt", Typeflag: tar.TypeReg, Mode: 0o644, Size: 7,
			},
			manifest:   ManifestEntry{Path: "Documents/file.txt", Type: EntryTypeDirectory},
			manifestID: "type-mismatch",
		},
		{
			name: "owner mismatch",
			header: tar.Header{
				Name: "Documents/file.txt", Typeflag: tar.TypeReg, Mode: 0o644, Size: 7, Uid: 501, Gid: 20,
			},
			manifest:   ManifestEntry{Path: "Documents/file.txt", Type: EntryTypeFile, UID: int64Pointer(502), GID: int64Pointer(20)},
			manifestID: "owner-mismatch",
		},
		{
			name: "symlink target mismatch",
			header: tar.Header{
				Name: "Documents/link", Typeflag: tar.TypeSymlink, Linkname: "file.txt", Mode: 0o777,
			},
			manifest:   ManifestEntry{Path: "Documents/link", Type: EntryTypeSymlink, Target: "other.txt"},
			manifestID: "link-mismatch",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			root := t.TempDir()
			archivePath := filepath.Join(root, test.manifestID+".tgz")
			manifest := Manifest{SchemaVersion: 1, Entries: []ManifestEntry{test.manifest}}
			writeTestArchiveFixture(t, archivePath, []testArchiveMember{{Header: test.header, Data: []byte("payload")}}, manifest)
			assertStageRestoreFails(t, archivePath)
		})
	}
}

func TestStageRestore_RejectsArchivePathAttacks(t *testing.T) {
	tests := []struct {
		name     string
		members  []testArchiveMember
		manifest string
	}{
		{
			name: "parent traversal",
			members: []testArchiveMember{{
				Header: tar.Header{Name: "../escape", Typeflag: tar.TypeReg, Mode: 0o600, Size: 4},
				Data:   []byte("evil"),
			}},
			manifest: "schema_version: 1\nentries:\n  - path: ../escape\n    type: file\n",
		},
		{
			name: "absolute member",
			members: []testArchiveMember{{
				Header: tar.Header{Name: "/tmp/escape", Typeflag: tar.TypeReg, Mode: 0o600, Size: 4},
				Data:   []byte("evil"),
			}},
			manifest: "schema_version: 1\nentries:\n  - path: /tmp/escape\n    type: file\n",
		},
		{
			name: "hard link",
			members: []testArchiveMember{{
				Header: tar.Header{Name: "Documents/link", Typeflag: tar.TypeLink, Linkname: "Documents/file.txt", Mode: 0o600},
			}},
			manifest: "schema_version: 1\nentries:\n  - path: Documents/link\n    type: file\n",
		},
		{
			name: "device entry",
			members: []testArchiveMember{{
				Header: tar.Header{Name: "Documents/device", Typeflag: tar.TypeChar, Devmajor: 1, Devminor: 3, Mode: 0o600},
			}},
			manifest: "schema_version: 1\nentries:\n  - path: Documents/device\n    type: file\n",
		},
		{
			name: "duplicate archive path",
			members: []testArchiveMember{
				{Header: tar.Header{Name: "Documents/file.txt", Typeflag: tar.TypeReg, Mode: 0o600, Size: 1}, Data: []byte("a")},
				{Header: tar.Header{Name: "Documents/file.txt", Typeflag: tar.TypeReg, Mode: 0o600, Size: 1}, Data: []byte("b")},
			},
			manifest: "schema_version: 1\nentries:\n  - path: Documents/file.txt\n    type: file\n",
		},
		{
			name: "member below archive symlink",
			members: []testArchiveMember{
				{Header: tar.Header{Name: "Documents/link", Typeflag: tar.TypeSymlink, Linkname: "target", Mode: 0o777}},
				{Header: tar.Header{Name: "Documents/link/file.txt", Typeflag: tar.TypeReg, Mode: 0o600, Size: 1}, Data: []byte("a")},
			},
			manifest: "schema_version: 1\nentries:\n  - path: Documents/link\n    type: symlink\n    target: target\n  - path: Documents/link/file.txt\n    type: file\n",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			root := t.TempDir()
			archivePath := filepath.Join(root, "malicious.tgz")
			writeTestArchiveFixtureWithManifestData(t, archivePath, test.members, []byte(test.manifest))
			assertStageRestoreFails(t, archivePath)
			assert.NoFileExists(t, filepath.Join(root, "escape"))
		})
	}
}

func TestStageRestore_UnsafeSymlinkRequiresOptIn(t *testing.T) {
	root := t.TempDir()
	archivePath := filepath.Join(root, "unsafe.tgz")
	manifest := Manifest{
		SchemaVersion: 1,
		Entries: []ManifestEntry{{
			Path:    "Documents/link",
			Type:    EntryTypeSymlink,
			Target:  "/etc/passwd",
			Unsafe:  true,
			Warning: "absolute symlink target; restore requires --allow-unsafe-links",
		}},
	}
	artifacts := writeTestArchiveFixture(t, archivePath, []testArchiveMember{{
		Header: tar.Header{Name: "Documents/link", Typeflag: tar.TypeSymlink, Linkname: "/etc/passwd", Mode: 0o777},
	}}, manifest)

	assertStageRestoreFails(t, artifacts.ArchivePath)

	staged, err := NewService(nil).stageRestore(context.Background(), RestoreOptions{
		Archive:          artifacts.ArchivePath,
		AllowUnsafeLinks: true,
	})
	if !assert.NoError(t, err) {
		return
	}
	defer func() { _ = staged.cleanup() }()
	linkPath := filepath.Join(staged.StageDir, "Documents", "link")
	linkInfo, err := os.Lstat(linkPath)
	if !assert.NoError(t, err) {
		return
	}
	assert.True(t, linkInfo.Mode()&os.ModeSymlink != 0)
	target, err := os.Readlink(linkPath)
	assert.NoError(t, err)
	assert.Equal(t, "/etc/passwd", target)
}

func TestStageRestore_RejectsInvalidGzipAndTruncatedArchive(t *testing.T) {
	tests := []struct {
		name  string
		bytes []byte
	}{
		{name: "invalid gzip", bytes: []byte("not gzip")},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			root := t.TempDir()
			archivePath := filepath.Join(root, "invalid.tgz")
			if err := os.WriteFile(archivePath, test.bytes, 0o600); err != nil {
				t.Fatal(err)
			}
			writeTestSidecars(t, archivePath, []byte("schema_version: 1\nentries: []\n"))
			assertStageRestoreFails(t, archivePath)
		})
	}

	t.Run("truncated gzip trailer", func(t *testing.T) {
		root := t.TempDir()
		archivePath := filepath.Join(root, "truncated.tgz")
		artifacts := writeValidRestoreFixture(t, archivePath)
		archiveData, err := os.ReadFile(archivePath)
		if !assert.NoError(t, err) {
			return
		}
		if len(archiveData) <= 8 {
			t.Fatal("archive fixture unexpectedly short")
		}
		if err := os.WriteFile(archivePath, archiveData[:len(archiveData)-8], 0o600); err != nil {
			t.Fatal(err)
		}
		writeTestSidecars(t, archivePath, mustMarshalTestManifest(t, validRestoreManifest()))
		assertStageRestoreFails(t, artifacts.ArchivePath)
	})
}

func writeValidRestoreFixture(t *testing.T, archivePath string) ArtifactSet {
	t.Helper()
	return writeTestArchiveFixture(t, archivePath, validRestoreMembers(), validRestoreManifest())
}

func validRestoreMembers() []testArchiveMember {
	return []testArchiveMember{
		{Header: tar.Header{Name: "Documents/", Typeflag: tar.TypeDir, Mode: 0o755}},
		{Header: tar.Header{Name: "Documents/file.txt", Typeflag: tar.TypeReg, Mode: 0o644, Size: 7}, Data: []byte("payload")},
	}
}

func validRestoreManifest() Manifest {
	return Manifest{
		SchemaVersion: 1,
		Entries: []ManifestEntry{
			{Path: "Documents", Type: EntryTypeDirectory},
			{Path: "Documents/file.txt", Type: EntryTypeFile},
		},
	}
}

func mustMarshalTestManifest(t *testing.T, manifest Manifest) []byte {
	t.Helper()
	data, err := marshalManifest(manifest)
	if err != nil {
		t.Fatal(err)
	}
	return data
}

func int64Pointer(value int64) *int64 {
	return &value
}

func assertStageRestoreFails(t *testing.T, archivePath string) {
	t.Helper()
	destination := filepath.Join(t.TempDir(), "destination")
	_, err := NewService(nil).stageRestore(context.Background(), RestoreOptions{
		Archive:     archivePath,
		Destination: destination,
	})
	assert.Error(t, err)
	assert.NoDirExists(t, destination)
}

func TestStageRestore_RejectsManifestOwnerMismatchWhenFieldsOmitted(t *testing.T) {
	root := t.TempDir()
	archivePath := filepath.Join(root, "owner-mismatch.tgz")
	manifest := Manifest{
		SchemaVersion: 1,
		Entries:       []ManifestEntry{{Path: "file.txt", Type: EntryTypeFile}},
	}
	writeTestArchiveFixture(t, archivePath, []testArchiveMember{{
		Header: tar.Header{Name: "file.txt", Typeflag: tar.TypeReg, Mode: 0o600, Size: 1, Uid: 501, Gid: 20},
		Data:   []byte("x"),
	}}, manifest)
	assertStageRestoreFails(t, archivePath)
}

func TestStageRestore_RejectsUnsafeSymlinkManifestMarkerMismatch(t *testing.T) {
	root := t.TempDir()
	archivePath := filepath.Join(root, "unsafe-marker-mismatch.tgz")
	manifest := []byte("schema_version: 1\nentries:\n  - path: link\n    type: symlink\n    target: /etc/passwd\n    unsafe: false\n")
	writeTestArchiveFixtureWithManifestData(t, archivePath, []testArchiveMember{{
		Header: tar.Header{Name: "link", Typeflag: tar.TypeSymlink, Linkname: "/etc/passwd", Mode: 0o777},
	}}, manifest)
	assertStageRestoreFails(t, archivePath)
}
