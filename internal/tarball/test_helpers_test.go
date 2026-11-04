package tarball

import (
	"compress/gzip"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"archive/tar"
)

type testArchiveEntry struct {
	Header  tar.Header
	Content []byte
}

func writeTestFile(t *testing.T, path string, content []byte, mode os.FileMode, modified time.Time) {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, content, mode); err != nil {
		t.Fatal(err)
	}
	if err := os.Chmod(path, mode); err != nil {
		t.Fatal(err)
	}
	if !modified.IsZero() {
		if err := os.Chtimes(path, modified, modified); err != nil {
			t.Fatal(err)
		}
	}
}

func readTestArchive(t *testing.T, archivePath string) map[string]testArchiveEntry {
	t.Helper()
	archiveFile, err := os.Open(archivePath)
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		if err := archiveFile.Close(); err != nil {
			t.Errorf("close archive: %v", err)
		}
	}()

	compressed, err := gzip.NewReader(archiveFile)
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		if err := compressed.Close(); err != nil {
			t.Errorf("close gzip reader: %v", err)
		}
	}()

	entries := make(map[string]testArchiveEntry)
	reader := tar.NewReader(compressed)
	for {
		header, err := reader.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			t.Fatal(err)
		}

		name := strings.TrimSuffix(header.Name, "/")
		entry := testArchiveEntry{Header: *header}
		if header.Typeflag == tar.TypeReg || header.Typeflag == tar.TypeRegA {
			entry.Content, err = io.ReadAll(reader)
			if err != nil {
				t.Fatal(err)
			}
		}
		entries[name] = entry
	}
	return entries
}

func readTestManifest(t *testing.T, manifestPath string) Manifest {
	t.Helper()
	data, err := os.ReadFile(manifestPath)
	if err != nil {
		t.Fatal(err)
	}
	manifest, err := parseManifest(data)
	if err != nil {
		t.Fatal(err)
	}
	return manifest
}
