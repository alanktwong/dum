package tarball

import (
	"archive/tar"
	"context"
	"io"
	"os"
	"path/filepath"
	"strings"
	"syscall"
	"testing"

	asrt "github.com/stretchr/testify/assert"
)

func TestBackupRejectsInvalidOptionsAndSourceTypes(t *testing.T) {
	root := t.TempDir()
	validSource := filepath.Join(root, "source")
	if err := os.Mkdir(validSource, 0o700); err != nil {
		t.Fatal(err)
	}
	missingSource := filepath.Join(root, "missing")
	pipePath := filepath.Join(root, "pipe")
	if err := syscall.Mkfifo(pipePath, 0o600); err != nil {
		t.Fatal(err)
	}
	canceled, cancel := context.WithCancel(context.Background())
	cancel()

	tests := []struct {
		name    string
		ctx     context.Context
		options BackupOptions
	}{
		{name: "nil context", options: BackupOptions{Paths: []string{validSource}, Target: root}},
		{name: "canceled context", ctx: canceled, options: BackupOptions{Paths: []string{validSource}, Target: root}},
		{name: "no paths", ctx: context.Background(), options: BackupOptions{Target: root}},
		{name: "empty target", ctx: context.Background(), options: BackupOptions{Paths: []string{validSource}}},
		{name: "empty source", ctx: context.Background(), options: BackupOptions{Paths: []string{""}, Target: root}},
		{
			name:    "missing source",
			ctx:     context.Background(),
			options: BackupOptions{Paths: []string{missingSource}, Target: root},
		},
		{
			name:    "unsupported source type",
			ctx:     context.Background(),
			options: BackupOptions{Paths: []string{pipePath}, Target: root},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			_, err := NewService(nil).Backup(test.ctx, test.options)
			asrt.Error(t, err)
		})
	}
}

func TestRestoreRejectsInvalidContextAndArchiveOptions(t *testing.T) {
	canceled, cancel := context.WithCancel(context.Background())
	cancel()
	tests := []struct {
		name    string
		ctx     context.Context
		options RestoreOptions
	}{
		{name: "nil context", options: RestoreOptions{Archive: "archive.tgz"}},
		{name: "canceled context", ctx: canceled, options: RestoreOptions{Archive: "archive.tgz"}},
		{name: "empty archive", ctx: context.Background()},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			_, err := NewService(nil).Restore(test.ctx, test.options)
			asrt.Error(t, err)
		})
	}
}

func TestParseManifestRejectsMultipleDocuments(t *testing.T) {
	data := []byte("schema_version: 1\nentries: []\n---\nschema_version: 1\nentries: []\n")
	_, err := parseManifest(data)
	asrt.Error(t, err)
}

func TestReadChecksumSidecarRejectsUppercaseDigest(t *testing.T) {
	path := filepath.Join(t.TempDir(), "archive.tgz.sha256")
	data := []byte(strings.Repeat("A", 64) + "\n")
	if err := os.WriteFile(path, data, 0o600); err != nil {
		t.Fatal(err)
	}
	_, err := readChecksumSidecar(path)
	asrt.Error(t, err)
}

func TestValidateArchiveHeaderRejectsInvalidMetadata(t *testing.T) {
	tests := []struct {
		name   string
		header tar.Header
		type_  EntryType
	}{
		{name: "negative uid", header: tar.Header{Uid: -1}, type_: EntryTypeFile},
		{name: "negative gid", header: tar.Header{Gid: -1}, type_: EntryTypeFile},
		{name: "negative mode", header: tar.Header{Mode: -1}, type_: EntryTypeFile},
		{name: "non-file with size", header: tar.Header{Size: 1}, type_: EntryTypeDirectory},
		{name: "negative file size", header: tar.Header{Size: -1}, type_: EntryTypeFile},
		{name: "non-symlink link target", header: tar.Header{Linkname: "target"}, type_: EntryTypeFile},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			asrt.Error(t, validateArchiveHeader(&test.header, "member", test.type_))
		})
	}
}

func TestExtractRegularFileRejectsInvalidSizeAndExistingPath(t *testing.T) {
	stage := t.TempDir()
	root, err := os.OpenRoot(stage)
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = root.Close() }()

	asrt.Error(t, extractRegularFile(root, "negative", strings.NewReader("x"), -1))
	asrt.Error(t, extractRegularFile(root, "short", strings.NewReader("x"), 2))
	if err := os.WriteFile(filepath.Join(stage, "exists"), []byte("old"), 0o600); err != nil {
		t.Fatal(err)
	}
	asrt.Error(t, extractRegularFile(root, "exists", strings.NewReader("new"), 3))
}

func TestFilesystemHelpersWrapMissingPathErrors(t *testing.T) {
	root := t.TempDir()
	_, err := hashFile(filepath.Join(root, "missing"))
	asrt.Error(t, err)
	_, err = writeTempFile(filepath.Join(root, "missing-directory"), "file-*.tmp", []byte("data"))
	asrt.Error(t, err)
}

func TestRemovePublishedArtifactsLeavesCompetingFileUntouched(t *testing.T) {
	root := t.TempDir()
	stagedPath := filepath.Join(root, "staged")
	finalPath := filepath.Join(root, "final")
	if err := os.WriteFile(stagedPath, []byte("staged"), 0o600); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(finalPath, []byte("competitor"), 0o600); err != nil {
		t.Fatal(err)
	}

	err := removePublishedArtifacts([]publishedArtifact{{stagedPath: stagedPath, finalPath: finalPath}})
	if !asrt.NoError(t, err) {
		return
	}
	data, err := os.ReadFile(finalPath)
	if asrt.NoError(t, err) {
		asrt.Equal(t, []byte("competitor"), data)
	}
}

func TestStagedRestoreCleanupHandlesNilAndRemovesTemporaryDirectory(t *testing.T) {
	var nilRestore *stagedRestore
	asrt.NoError(t, nilRestore.cleanup())

	temporaryDir := t.TempDir()
	restore := &stagedRestore{TemporaryDir: temporaryDir}
	if !asrt.NoError(t, restore.cleanup()) {
		return
	}
	_, err := os.Stat(temporaryDir)
	asrt.ErrorIs(t, err, os.ErrNotExist)
}

func TestChangeOwnerDefaultWrapsMissingFileErrors(t *testing.T) {
	rootDir := t.TempDir()
	root, err := os.OpenRoot(rootDir)
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = root.Close() }()

	asrt.Error(t, changeOwnerDefault(root, "missing-file", os.Geteuid(), os.Getegid(), false))
	asrt.Error(t, changeOwnerDefault(root, "missing-link", os.Geteuid(), os.Getegid(), true))
}

func TestCanonicalPathRejectsFileAsParentAndBrokenSymlink(t *testing.T) {
	root := t.TempDir()
	regularFile := filepath.Join(root, "file")
	if err := os.WriteFile(regularFile, []byte("data"), 0o600); err != nil {
		t.Fatal(err)
	}
	_, err := canonicalPath(filepath.Join(regularFile, "child"))
	asrt.Error(t, err)

	brokenLink := filepath.Join(root, "broken")
	if err := os.Symlink(filepath.Join(root, "missing-target"), brokenLink); err != nil {
		t.Fatal(err)
	}
	_, err = canonicalPath(brokenLink)
	asrt.Error(t, err)
}

func TestChangeOwnerDefaultHandlesFilesAndSymlinks(t *testing.T) {
	directory := t.TempDir()
	if err := os.WriteFile(filepath.Join(directory, "file"), []byte("data"), 0o600); err != nil {
		t.Fatal(err)
	}
	if err := os.Symlink("file", filepath.Join(directory, "link")); err != nil {
		t.Fatal(err)
	}
	root, err := os.OpenRoot(directory)
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = root.Close() }()

	if err := changeOwnerDefault(root, "file", os.Geteuid(), os.Getegid(), false); err != nil {
		t.Skipf("filesystem does not permit owner-setting no-op: %v", err)
	}
	if err := changeOwnerDefault(root, "link", os.Geteuid(), os.Getegid(), true); err != nil {
		t.Skipf("filesystem does not permit symlink owner-setting no-op: %v", err)
	}
}

func TestContextReaderPreservesEOF(t *testing.T) {
	reader := contextReader{ctx: context.Background(), reader: strings.NewReader("")}
	_, err := reader.Read(make([]byte, 1))
	asrt.ErrorIs(t, err, io.EOF)
}
