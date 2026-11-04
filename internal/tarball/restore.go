package tarball

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

type stagedRestore struct {
	TemporaryDir string
	SnapshotPath string
	StageDir     string
	Manifest     Manifest
	Members      []archiveMember
}

func (restore *stagedRestore) cleanup() error {
	if restore == nil || restore.TemporaryDir == "" {
		return nil
	}
	if err := os.RemoveAll(restore.TemporaryDir); err != nil {
		return fmt.Errorf("remove restore staging directory %q: %w", restore.TemporaryDir, err)
	}
	return nil
}

func (s *Service) stageRestore(ctx context.Context, options RestoreOptions) (*stagedRestore, error) {
	if ctx == nil {
		return nil, fmt.Errorf("restore: nil context")
	}
	if err := ctx.Err(); err != nil {
		return nil, fmt.Errorf("restore: %w", err)
	}
	if options.Archive == "" {
		return nil, fmt.Errorf("restore archive cannot be empty")
	}

	archivePath, err := filepath.Abs(options.Archive)
	if err != nil {
		return nil, fmt.Errorf("resolve archive path %q: %w", options.Archive, err)
	}
	checksumPath := archivePath + ".sha256"
	manifestPath := archivePath + ".manifest.yaml"
	expectedChecksum, err := readChecksumSidecar(checksumPath)
	if err != nil {
		return nil, fmt.Errorf("read checksum sidecar for %q: %w", archivePath, err)
	}

	temporaryDir, err := os.MkdirTemp("", "dum-tarball-restore-*")
	if err != nil {
		return nil, fmt.Errorf("create private restore directory: %w", err)
	}
	keepTemporaryDir := false
	defer func() {
		if !keepTemporaryDir {
			_ = os.RemoveAll(temporaryDir)
		}
	}()

	snapshotPath, actualChecksum, err := snapshotAndHash(ctx, archivePath, temporaryDir)
	if err != nil {
		return nil, fmt.Errorf("verify archive %q: %w", archivePath, err)
	}
	if actualChecksum != expectedChecksum {
		return nil, fmt.Errorf("checksum mismatch for %q", archivePath)
	}

	manifestData, err := os.ReadFile(manifestPath)
	if err != nil {
		return nil, fmt.Errorf("read manifest sidecar for %q: %w", archivePath, err)
	}
	manifest, err := parseManifest(manifestData)
	if err != nil {
		return nil, fmt.Errorf("parse manifest sidecar for %q: %w", archivePath, err)
	}
	members, err := inspectArchive(snapshotPath, manifest, options.AllowUnsafeLinks)
	if err != nil {
		return nil, fmt.Errorf("validate archive %q: %w", archivePath, err)
	}

	stageDir, err := os.MkdirTemp(temporaryDir, "stage-*")
	if err != nil {
		return nil, fmt.Errorf("create private extraction stage: %w", err)
	}
	if err := extractArchiveToStage(ctx, snapshotPath, stageDir, members); err != nil {
		return nil, fmt.Errorf("stage archive %q: %w", archivePath, err)
	}

	keepTemporaryDir = true
	return &stagedRestore{
		TemporaryDir: temporaryDir,
		SnapshotPath: snapshotPath,
		StageDir:     stageDir,
		Manifest:     manifest,
		Members:      members,
	}, nil
}

func readChecksumSidecar(checksumPath string) (string, error) {
	data, err := os.ReadFile(checksumPath)
	if err != nil {
		return "", fmt.Errorf("read file: %w", err)
	}
	if len(data) != sha256.Size*2+1 || data[len(data)-1] != '\n' {
		return "", fmt.Errorf("checksum must be 64 lowercase hexadecimal characters followed by a newline")
	}
	checksum := string(data[:len(data)-1])
	decoded, err := hex.DecodeString(checksum)
	if err != nil || hex.EncodeToString(decoded) != checksum {
		return "", fmt.Errorf("checksum is not lowercase hexadecimal")
	}
	return checksum, nil
}

func snapshotAndHash(ctx context.Context, archivePath, temporaryDir string) (string, string, error) {
	source, err := os.Open(archivePath)
	if err != nil {
		return "", "", fmt.Errorf("open archive: %w", err)
	}
	defer func() { _ = source.Close() }()

	snapshot, err := os.CreateTemp(temporaryDir, "archive-*.tgz")
	if err != nil {
		return "", "", fmt.Errorf("create verified snapshot: %w", err)
	}
	snapshotPath := snapshot.Name()
	hash := sha256.New()
	_, copyErr := io.Copy(io.MultiWriter(snapshot, hash), contextReader{ctx: ctx, reader: source})
	if copyErr == nil {
		copyErr = snapshot.Sync()
	}
	closeErr := snapshot.Close()
	if err := errors.Join(copyErr, closeErr); err != nil {
		_ = os.Remove(snapshotPath)
		return "", "", fmt.Errorf("copy archive snapshot: %w", err)
	}
	if err := ctx.Err(); err != nil {
		_ = os.Remove(snapshotPath)
		return "", "", fmt.Errorf("copy archive snapshot: %w", err)
	}
	return snapshotPath, hex.EncodeToString(hash.Sum(nil)), nil
}

type contextReader struct {
	ctx    context.Context
	reader io.Reader
}

func (reader contextReader) Read(data []byte) (int, error) {
	if err := reader.ctx.Err(); err != nil {
		return 0, err
	}
	return reader.reader.Read(data)
}

func containsUnsafeWarning(warning, reason string) bool {
	return warning != "" && strings.Contains(warning, reason)
}
