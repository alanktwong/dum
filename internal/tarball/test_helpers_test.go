package tarball

import (
	"archive/tar"
	"compress/gzip"
	"crypto/sha256"
	"encoding/hex"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	yamlv3 "gopkg.in/yaml.v3"
)

type testArchiveEntry struct {
	Header  tar.Header
	Content []byte
}

type testArchiveMember struct {
	Header tar.Header
	Data   []byte
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
		if header.Typeflag == tar.TypeReg || header.Typeflag == 0 {
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

func writeTestArchiveFixture(
	t *testing.T,
	archivePath string,
	members []testArchiveMember,
	manifest Manifest,
) ArtifactSet {
	t.Helper()
	manifestData, err := yamlv3.Marshal(manifest)
	if err != nil {
		t.Fatal(err)
	}
	return writeTestArchiveFixtureWithManifestData(t, archivePath, members, manifestData)
}

func writeTestArchiveFixtureWithManifestData(
	t *testing.T,
	archivePath string,
	members []testArchiveMember,
	manifestData []byte,
) ArtifactSet {
	return writeTestArchiveFixtureWithManifestDataAndTrailingBytes(
		t, archivePath, members, manifestData, nil,
	)
}

func writeTestArchiveFixtureWithManifestDataAndTrailingBytes(
	t *testing.T,
	archivePath string,
	members []testArchiveMember,
	manifestData, trailingBytes []byte,
) ArtifactSet {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(archivePath), 0o755); err != nil {
		t.Fatal(err)
	}
	archiveFile, err := os.Create(archivePath)
	if err != nil {
		t.Fatal(err)
	}
	compressed := gzip.NewWriter(archiveFile)
	archive := tar.NewWriter(compressed)
	for _, member := range members {
		header := member.Header
		if err := archive.WriteHeader(&header); err != nil {
			t.Fatal(err)
		}
		if header.Typeflag == tar.TypeReg || header.Typeflag == 0 {
			if _, err := archive.Write(member.Data); err != nil {
				t.Fatal(err)
			}
		}
	}
	if err := archive.Close(); err != nil {
		t.Fatal(err)
	}
	if len(trailingBytes) > 0 {
		if _, err := compressed.Write(trailingBytes); err != nil {
			t.Fatal(err)
		}
	}
	if err := compressed.Close(); err != nil {
		t.Fatal(err)
	}
	if err := archiveFile.Close(); err != nil {
		t.Fatal(err)
	}
	return writeTestSidecars(t, archivePath, manifestData)
}

func writeTestSidecars(t *testing.T, archivePath string, manifestData []byte) ArtifactSet {
	t.Helper()
	archiveData, err := os.ReadFile(archivePath)
	if err != nil {
		t.Fatal(err)
	}
	digest := sha256.Sum256(archiveData)
	artifacts := ArtifactSet{
		ArchivePath:  archivePath,
		ChecksumPath: archivePath + ".sha256",
		ManifestPath: archivePath + ".manifest.yaml",
	}
	if err := os.WriteFile(artifacts.ChecksumPath, []byte(hex.EncodeToString(digest[:])+"\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(artifacts.ManifestPath, manifestData, 0o600); err != nil {
		t.Fatal(err)
	}
	return artifacts
}
