package tarball

import (
	"archive/tar"
	"compress/gzip"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"os"
	"path"
	"path/filepath"
	"syscall"
	"time"
)

type backupPlan struct {
	source     string
	memberRoot string
	artifacts  ArtifactSet
}

type publishedArtifact struct {
	stagedPath string
	finalPath  string
}

// Backup creates one archive and two sidecars for each selected path.
func (s *Service) Backup(ctx context.Context, options BackupOptions) (BackupResult, error) {
	if ctx == nil {
		return BackupResult{}, fmt.Errorf("backup: nil context")
	}
	if err := ctx.Err(); err != nil {
		return BackupResult{}, fmt.Errorf("backup: %w", err)
	}
	plans, target, err := s.planBackup(options)
	if err != nil {
		return BackupResult{}, err
	}
	planned := backupResult(plans)
	if options.DryRun {
		return planned, nil
	}
	if err := os.MkdirAll(target, 0o750); err != nil {
		return BackupResult{}, fmt.Errorf("create backup target %q: %w", target, err)
	}
	return s.writeBackupPlans(ctx, plans)
}

func (s *Service) planBackup(options BackupOptions) ([]backupPlan, string, error) {
	if len(options.Paths) == 0 {
		return nil, "", fmt.Errorf("backup requires at least one path")
	}
	if options.Target == "" {
		return nil, "", fmt.Errorf("backup target cannot be empty")
	}

	home := options.Home
	if home == "" {
		var err error
		home, err = os.UserHomeDir()
		if err != nil {
			return nil, "", fmt.Errorf("resolve home directory: %w", err)
		}
	}
	absoluteHome, err := filepath.Abs(home)
	if err != nil {
		return nil, "", fmt.Errorf("resolve home path %q: %w", home, err)
	}
	absoluteTarget, err := filepath.Abs(options.Target)
	if err != nil {
		return nil, "", fmt.Errorf("resolve target path %q: %w", options.Target, err)
	}

	sources, err := normalizeBackupSources(options.Paths)
	if err != nil {
		return nil, "", err
	}
	if err := validateUniqueBasenames(sources); err != nil {
		return nil, "", fmt.Errorf("validate backup paths: %w", err)
	}
	if err := validateTargetDoesNotOverlap(absoluteTarget, sources); err != nil {
		return nil, "", fmt.Errorf("validate backup target: %w", err)
	}

	now := time.Now
	if s != nil && s.now != nil {
		now = s.now
	}
	timestamp := now().Format("2006-01-02-150405")
	plans := make([]backupPlan, len(sources))
	for index, source := range sources {
		memberRoot, err := archivePathForSource(absoluteHome, source)
		if err != nil {
			return nil, "", fmt.Errorf("resolve archive root for %q: %w", source, err)
		}
		base := filepath.Base(source)
		archivePath := filepath.Join(absoluteTarget, fmt.Sprintf("%s-%s.tgz", base, timestamp))
		artifacts := ArtifactSet{
			ArchivePath:  archivePath,
			ChecksumPath: archivePath + ".sha256",
			ManifestPath: archivePath + ".manifest.yaml",
		}
		for _, outputPath := range []string{artifacts.ArchivePath, artifacts.ChecksumPath, artifacts.ManifestPath} {
			if err := ensurePathAbsent(outputPath); err != nil {
				return nil, "", fmt.Errorf("preflight output for %q: %w", source, err)
			}
		}
		plans[index] = backupPlan{source: source, memberRoot: memberRoot, artifacts: artifacts}
	}
	return plans, filepath.Clean(absoluteTarget), nil
}

func normalizeBackupSources(paths []string) ([]string, error) {
	sources := make([]string, len(paths))
	for index, source := range paths {
		if source == "" {
			return nil, fmt.Errorf("backup path at index %d is empty", index)
		}
		absoluteSource, err := filepath.Abs(source)
		if err != nil {
			return nil, fmt.Errorf("resolve source path %q: %w", source, err)
		}
		info, err := os.Lstat(absoluteSource)
		if err != nil {
			return nil, fmt.Errorf("inspect source path %q: %w", absoluteSource, err)
		}
		if !info.IsDir() && !info.Mode().IsRegular() && info.Mode()&os.ModeSymlink == 0 {
			return nil, fmt.Errorf("unsupported source type at %q: %s", absoluteSource, info.Mode().Type())
		}
		sources[index] = filepath.Clean(absoluteSource)
	}
	return sources, nil
}

func backupResult(plans []backupPlan) BackupResult {
	planned := BackupResult{Artifacts: make([]ArtifactSet, 0, len(plans))}
	for _, plan := range plans {
		planned.Artifacts = append(planned.Artifacts, plan.artifacts)
	}
	return planned
}

func (s *Service) writeBackupPlans(ctx context.Context, plans []backupPlan) (BackupResult, error) {
	result := BackupResult{Artifacts: make([]ArtifactSet, 0, len(plans))}
	var failures []error
	for _, plan := range plans {
		if err := ctx.Err(); err != nil {
			failures = append(failures, fmt.Errorf("backup stopped: %w", err))
			break
		}
		if err := s.backupOne(ctx, plan); err != nil {
			failures = append(failures, fmt.Errorf("backup %q: %w", plan.source, err))
			continue
		}
		result.Artifacts = append(result.Artifacts, plan.artifacts)
	}
	if err := errors.Join(failures...); err != nil {
		return result, err
	}
	return result, nil
}

func (s *Service) backupOne(ctx context.Context, plan backupPlan) error {
	target := filepath.Dir(plan.artifacts.ArchivePath)
	archiveFile, err := os.CreateTemp(target, ".dum-tarball-archive-*.tmp")
	if err != nil {
		return fmt.Errorf("create temporary archive: %w", err)
	}
	archiveTempPath := archiveFile.Name()
	defer func() {
		_ = archiveFile.Close()
		_ = os.Remove(archiveTempPath)
	}()

	manifest, err := s.writeArchive(ctx, archiveFile, plan.source, plan.memberRoot)
	if err != nil {
		return fmt.Errorf("write archive: %w", err)
	}
	if err := archiveFile.Sync(); err != nil {
		return fmt.Errorf("sync temporary archive: %w", err)
	}
	if err := archiveFile.Close(); err != nil {
		return fmt.Errorf("close temporary archive: %w", err)
	}

	digest, err := hashFile(archiveTempPath)
	if err != nil {
		return fmt.Errorf("hash temporary archive: %w", err)
	}
	checksumTempPath, err := writeTempFile(target, ".dum-tarball-checksum-*.tmp", []byte(digest+"\n"))
	if err != nil {
		return fmt.Errorf("write checksum sidecar: %w", err)
	}
	defer func() { _ = os.Remove(checksumTempPath) }()

	manifestData, err := marshalManifest(manifest)
	if err != nil {
		return fmt.Errorf("marshal archive manifest: %w", err)
	}
	manifestTempPath, err := writeTempFile(target, ".dum-tarball-manifest-*.tmp", manifestData)
	if err != nil {
		return fmt.Errorf("write manifest sidecar: %w", err)
	}
	defer func() { _ = os.Remove(manifestTempPath) }()

	publish := s.publish
	if publish == nil {
		publish = os.Link
	}
	published := make([]publishedArtifact, 0, 2)
	for _, artifact := range []publishedArtifact{
		{stagedPath: checksumTempPath, finalPath: plan.artifacts.ChecksumPath},
		{stagedPath: manifestTempPath, finalPath: plan.artifacts.ManifestPath},
	} {
		if err := publish(artifact.stagedPath, artifact.finalPath); err != nil {
			cleanupErr := removePublishedArtifacts(published)
			return errors.Join(fmt.Errorf("publish sidecar %q: %w", artifact.finalPath, err), cleanupErr)
		}
		published = append(published, artifact)
	}

	archiveArtifact := publishedArtifact{stagedPath: archiveTempPath, finalPath: plan.artifacts.ArchivePath}
	if err := publish(archiveArtifact.stagedPath, archiveArtifact.finalPath); err != nil {
		cleanupErr := removePublishedArtifacts(published)
		return errors.Join(fmt.Errorf("publish archive %q: %w", archiveArtifact.finalPath, err), cleanupErr)
	}
	return nil
}

func (s *Service) writeArchive(ctx context.Context, output io.Writer, source, memberRoot string) (Manifest, error) {
	compressed := gzip.NewWriter(output)
	archive := tar.NewWriter(compressed)
	manifest := Manifest{SchemaVersion: manifestSchemaVersion}

	walkErr := filepath.WalkDir(source, func(currentPath string, _ os.DirEntry, walkErr error) error {
		if walkErr != nil {
			return fmt.Errorf("walk source path %q: %w", currentPath, walkErr)
		}
		if err := ctx.Err(); err != nil {
			return fmt.Errorf("backup source walk canceled: %w", err)
		}

		info, err := os.Lstat(currentPath)
		if err != nil {
			return fmt.Errorf("inspect source entry %q: %w", currentPath, err)
		}
		memberPath, skip, err := archiveMemberPath(source, currentPath, memberRoot, info)
		if err != nil {
			return err
		}
		if skip {
			return nil
		}
		if err := validateMemberPath(memberPath); err != nil {
			return fmt.Errorf("invalid archive member %q: %w", memberPath, err)
		}

		uid, gid, hasOwner := fileOwner(info)
		header := &tar.Header{
			Name:    memberPath,
			Mode:    int64(info.Mode().Perm() & 0o777),
			ModTime: info.ModTime(),
			Uid:     0,
			Gid:     0,
			Format:  tar.FormatPAX,
		}
		manifestEntry := ManifestEntry{Path: memberPath}
		if hasOwner {
			header.Uid = int(uid)
			header.Gid = int(gid)
			manifestEntry.UID = &uid
			manifestEntry.GID = &gid
		}

		mode := info.Mode()
		switch {
		case mode&os.ModeSymlink != 0:
			target, err := os.Readlink(currentPath)
			if err != nil {
				return fmt.Errorf("read symlink %q: %w", currentPath, err)
			}
			unsafe, warning := classifySymlinkTarget(memberPath, target)
			header.Typeflag = tar.TypeSymlink
			header.Linkname = target
			manifestEntry.Type = EntryTypeSymlink
			manifestEntry.Target = target
			manifestEntry.Unsafe = unsafe
			if unsafe {
				manifestEntry.Warning = warning + "; restore requires --allow-unsafe-links"
				if s.Log != nil {
					s.Log.Warn("unsafe symlink in backup", "path", currentPath, "target", target, "warning", warning)
				}
			}
			if err := archive.WriteHeader(header); err != nil {
				return fmt.Errorf("write symlink header %q: %w", currentPath, err)
			}
		case info.IsDir():
			header.Typeflag = tar.TypeDir
			header.Name = memberPath + "/"
			manifestEntry.Type = EntryTypeDirectory
			if err := archive.WriteHeader(header); err != nil {
				return fmt.Errorf("write directory header %q: %w", currentPath, err)
			}
		case mode.IsRegular():
			manifestEntry.Type = EntryTypeFile
			header.Typeflag = tar.TypeReg
			if err := writeRegularFile(archive, header, info, currentPath); err != nil {
				return err
			}
		default:
			return fmt.Errorf("unsupported source entry %q with mode %s", currentPath, mode.Type())
		}

		manifest.Entries = append(manifest.Entries, manifestEntry)
		return nil
	})
	archiveCloseErr := archive.Close()
	gzipCloseErr := compressed.Close()
	if err := errors.Join(walkErr, archiveCloseErr, gzipCloseErr); err != nil {
		return Manifest{}, fmt.Errorf("finish archive: %w", err)
	}
	return manifest, nil
}

func writeRegularFile(archive *tar.Writer, header *tar.Header, info os.FileInfo, sourcePath string) error {
	file, err := os.Open(sourcePath)
	if err != nil {
		return fmt.Errorf("open source file %q: %w", sourcePath, err)
	}
	openedInfo, statErr := file.Stat()
	if statErr != nil {
		_ = file.Close()
		return fmt.Errorf("stat opened source file %q: %w", sourcePath, statErr)
	}
	if !os.SameFile(info, openedInfo) || !openedInfo.Mode().IsRegular() {
		_ = file.Close()
		return fmt.Errorf("source file changed during backup: %q", sourcePath)
	}
	header.Size = openedInfo.Size()
	if err := archive.WriteHeader(header); err != nil {
		_ = file.Close()
		return fmt.Errorf("write file header %q: %w", sourcePath, err)
	}
	written, copyErr := io.CopyN(archive, file, header.Size)
	closeErr := file.Close()
	if copyErr != nil {
		return fmt.Errorf("archive source file %q: %w", sourcePath, errors.Join(copyErr, closeErr))
	}
	if written != header.Size {
		return fmt.Errorf("archive source file %q: copied %d bytes, expected %d", sourcePath, written, header.Size)
	}
	if closeErr != nil {
		return fmt.Errorf("close source file %q: %w", sourcePath, closeErr)
	}
	return nil
}

func archiveMemberPath(source, currentPath, memberRoot string, info os.FileInfo) (string, bool, error) {
	relative, err := filepath.Rel(source, currentPath)
	if err != nil {
		return "", false, fmt.Errorf("make source entry %q relative to %q: %w", currentPath, source, err)
	}
	if relative == "." {
		if memberRoot == "" && info.IsDir() {
			return "", true, nil
		}
		if memberRoot == "" {
			return "", false, fmt.Errorf("cannot determine archive member for source root %q", source)
		}
		return memberRoot, false, nil
	}
	member := filepath.ToSlash(relative)
	if memberRoot != "" {
		member = path.Join(memberRoot, member)
	}
	return member, false, nil
}

func fileOwner(info os.FileInfo) (int64, int64, bool) {
	stat, ok := info.Sys().(*syscall.Stat_t)
	if !ok {
		return 0, 0, false
	}
	return int64(stat.Uid), int64(stat.Gid), true
}

func hashFile(filePath string) (string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return "", fmt.Errorf("open file %q: %w", filePath, err)
	}
	hash := sha256.New()
	_, copyErr := io.Copy(hash, file)
	closeErr := file.Close()
	if err := errors.Join(copyErr, closeErr); err != nil {
		return "", fmt.Errorf("hash file %q: %w", filePath, err)
	}
	return hex.EncodeToString(hash.Sum(nil)), nil
}

func writeTempFile(directory, pattern string, data []byte) (string, error) {
	file, err := os.CreateTemp(directory, pattern)
	if err != nil {
		return "", fmt.Errorf("create temporary file: %w", err)
	}
	filePath := file.Name()
	_, writeErr := file.Write(data)
	syncErr := file.Sync()
	closeErr := file.Close()
	if err := errors.Join(writeErr, syncErr, closeErr); err != nil {
		_ = os.Remove(filePath)
		return "", fmt.Errorf("write temporary file %q: %w", filePath, err)
	}
	return filePath, nil
}

func ensurePathAbsent(filePath string) error {
	_, err := os.Lstat(filePath)
	if err == nil {
		return fmt.Errorf("output already exists")
	}
	if !os.IsNotExist(err) {
		return fmt.Errorf("inspect output path: %w", err)
	}
	return nil
}

func removePublishedArtifacts(published []publishedArtifact) error {
	var failures []error
	for index := len(published) - 1; index >= 0; index-- {
		artifact := published[index]
		stagedInfo, err := os.Stat(artifact.stagedPath)
		if err != nil {
			if !os.IsNotExist(err) {
				failures = append(failures, fmt.Errorf("inspect staged artifact %q: %w", artifact.stagedPath, err))
			}
			continue
		}
		publishedInfo, err := os.Lstat(artifact.finalPath)
		if err != nil {
			if !os.IsNotExist(err) {
				failures = append(failures, fmt.Errorf("inspect published artifact %q: %w", artifact.finalPath, err))
			}
			continue
		}
		if !os.SameFile(stagedInfo, publishedInfo) {
			continue
		}
		if err := os.Remove(artifact.finalPath); err != nil && !os.IsNotExist(err) {
			failures = append(failures, fmt.Errorf("remove published artifact %q: %w", artifact.finalPath, err))
		}
	}
	return errors.Join(failures...)
}
