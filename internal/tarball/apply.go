package tarball

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"os"
	"path"
	"path/filepath"
	"sort"
	"strings"
)

func preflightDestination(destination string, members []archiveMember, overwrite bool) error {
	rootInfo, err := os.Lstat(destination)
	if err == nil {
		if rootInfo.Mode()&os.ModeSymlink != 0 || !rootInfo.IsDir() {
			return fmt.Errorf("destination root %q is not a real directory", destination)
		}
	} else if !os.IsNotExist(err) {
		return fmt.Errorf("inspect destination root %q: %w", destination, err)
	}

	var conflicts []error
	for _, member := range members {
		components := strings.Split(member.Path, "/")
		currentPath := destination
		for index, component := range components {
			currentPath = filepath.Join(currentPath, filepath.FromSlash(component))
			info, err := os.Lstat(currentPath)
			if os.IsNotExist(err) {
				break
			}
			if err != nil {
				conflicts = append(conflicts, fmt.Errorf("inspect destination path %q: %w", currentPath, err))
				break
			}
			isFinal := index == len(components)-1
			if !isFinal {
				if info.Mode()&os.ModeSymlink != 0 {
					conflicts = append(conflicts, fmt.Errorf("destination path %q traverses a symlink", currentPath))
					break
				}
				if !info.IsDir() {
					conflicts = append(conflicts, fmt.Errorf("destination parent %q is not a directory", currentPath))
					break
				}
				continue
			}

			if member.Type == EntryTypeDirectory {
				if info.Mode()&os.ModeSymlink != 0 || !info.IsDir() {
					conflicts = append(conflicts, fmt.Errorf("destination path %q has a file/directory type conflict", currentPath))
				}
				continue
			}
			if info.IsDir() {
				conflicts = append(conflicts, fmt.Errorf("destination directory %q cannot be replaced", currentPath))
				continue
			}
			if info.Mode()&os.ModeSymlink == 0 && !info.Mode().IsRegular() {
				conflicts = append(conflicts, fmt.Errorf("destination path %q has unsupported type %s", currentPath, info.Mode().Type()))
				continue
			}
			if !overwrite {
				conflicts = append(conflicts, fmt.Errorf("destination path already exists: %q", currentPath))
			}
		}
	}
	return errors.Join(conflicts...)
}

func (s *Service) applyStagedRestore(ctx context.Context, destination *os.Root, staged *stagedRestore, options RestoreOptions) error {
	manifestEntries := make(map[string]ManifestEntry, len(staged.Manifest.Entries))
	for _, entry := range staged.Manifest.Entries {
		manifestEntries[entry.Path] = entry
	}

	for _, member := range staged.Members {
		if err := ctx.Err(); err != nil {
			return fmt.Errorf("apply restore: %w", err)
		}
		parent := path.Dir(member.Path)
		if parent != "." {
			if err := destination.MkdirAll(filepath.FromSlash(parent), 0o755); err != nil {
				return fmt.Errorf("create destination parent for %q: %w", member.Path, err)
			}
		}
		if member.Type == EntryTypeDirectory {
			if err := destination.MkdirAll(filepath.FromSlash(member.Path), 0o755); err != nil {
				return fmt.Errorf("create destination directory %q: %w", member.Path, err)
			}
		}
	}

	for _, member := range staged.Members {
		if err := ctx.Err(); err != nil {
			return fmt.Errorf("apply restore: %w", err)
		}
		if member.Type != EntryTypeFile {
			continue
		}
		manifestEntry := manifestEntries[member.Path]
		if err := s.applyStagedFile(destination, staged.StageDir, member, manifestEntry, options); err != nil {
			return fmt.Errorf("apply file %q: %w", member.Path, err)
		}
	}

	for _, member := range staged.Members {
		if err := ctx.Err(); err != nil {
			return fmt.Errorf("apply restore: %w", err)
		}
		if member.Type != EntryTypeSymlink {
			continue
		}
		manifestEntry := manifestEntries[member.Path]
		if err := s.applyStagedSymlink(destination, member, manifestEntry, options); err != nil {
			return fmt.Errorf("apply symlink %q: %w", member.Path, err)
		}
	}

	directories := make([]archiveMember, 0)
	for _, member := range staged.Members {
		if member.Type == EntryTypeDirectory {
			directories = append(directories, member)
		}
	}
	sort.Slice(directories, func(left, right int) bool {
		return strings.Count(directories[left].Path, "/") > strings.Count(directories[right].Path, "/")
	})
	for _, member := range directories {
		if err := ctx.Err(); err != nil {
			return fmt.Errorf("apply restore: %w", err)
		}
		manifestEntry := manifestEntries[member.Path]
		if err := s.restoreOwner(destination, member.Path, manifestEntry, false, options.PreserveOwner); err != nil {
			return fmt.Errorf("restore directory owner for %q: %w", member.Path, err)
		}
		if err := destination.Chmod(filepath.FromSlash(member.Path), os.FileMode(member.Header.Mode).Perm()); err != nil {
			return fmt.Errorf("restore directory mode for %q: %w", member.Path, err)
		}
		if err := destination.Chtimes(filepath.FromSlash(member.Path), member.Header.ModTime, member.Header.ModTime); err != nil {
			return fmt.Errorf("restore directory timestamp for %q: %w", member.Path, err)
		}
	}
	return nil
}

func (s *Service) applyStagedFile(destination *os.Root, stageDir string, member archiveMember, manifestEntry ManifestEntry, options RestoreOptions) error {
	tempPath, target, err := createTempFileInRoot(destination, member.Path)
	if err != nil {
		return err
	}
	defer func() {
		_ = target.Close()
		_ = destination.Remove(filepath.FromSlash(tempPath))
	}()

	stagedPath := filepath.Join(stageDir, filepath.FromSlash(member.Path))
	stagedFile, err := os.Open(stagedPath)
	if err != nil {
		return fmt.Errorf("open staged file: %w", err)
	}
	written, copyErr := io.CopyN(target, stagedFile, member.Header.Size)
	stagedCloseErr := stagedFile.Close()
	if err := errors.Join(copyErr, stagedCloseErr); err != nil {
		return fmt.Errorf("copy staged file: %w", err)
	}
	if written != member.Header.Size {
		return fmt.Errorf("copied %d bytes, expected %d", written, member.Header.Size)
	}
	if err := target.Sync(); err != nil {
		return fmt.Errorf("sync staged destination file: %w", err)
	}
	if err := target.Close(); err != nil {
		return fmt.Errorf("close staged destination file: %w", err)
	}
	if err := s.restoreOwner(destination, tempPath, manifestEntry, false, options.PreserveOwner); err != nil {
		return err
	}
	if err := destination.Chmod(filepath.FromSlash(tempPath), os.FileMode(member.Header.Mode).Perm()); err != nil {
		return fmt.Errorf("set file mode: %w", err)
	}
	if err := destination.Chtimes(filepath.FromSlash(tempPath), member.Header.ModTime, member.Header.ModTime); err != nil {
		return fmt.Errorf("set file timestamp: %w", err)
	}
	if options.Overwrite {
		if err := destination.Rename(filepath.FromSlash(tempPath), filepath.FromSlash(member.Path)); err != nil {
			return fmt.Errorf("replace destination file: %w", err)
		}
		return nil
	}
	if err := destination.Link(filepath.FromSlash(tempPath), filepath.FromSlash(member.Path)); err != nil {
		return fmt.Errorf("publish destination file: %w", err)
	}
	if err := destination.Remove(filepath.FromSlash(tempPath)); err != nil {
		return fmt.Errorf("remove temporary destination link: %w", err)
	}
	return nil
}

func (s *Service) applyStagedSymlink(destination *os.Root, member archiveMember, manifestEntry ManifestEntry, options RestoreOptions) error {
	if options.Overwrite {
		tempPath, err := createTempSymlinkInRoot(destination, member.Path, member.Header.Linkname)
		if err != nil {
			return err
		}
		if err := s.restoreOwner(destination, tempPath, manifestEntry, true, options.PreserveOwner); err != nil {
			_ = destination.Remove(filepath.FromSlash(tempPath))
			return err
		}
		if err := destination.Rename(filepath.FromSlash(tempPath), filepath.FromSlash(member.Path)); err != nil {
			_ = destination.Remove(filepath.FromSlash(tempPath))
			return fmt.Errorf("replace destination symlink: %w", err)
		}
		return nil
	}
	if err := destination.Symlink(member.Header.Linkname, filepath.FromSlash(member.Path)); err != nil {
		return fmt.Errorf("create destination symlink: %w", err)
	}
	if err := s.restoreOwner(destination, member.Path, manifestEntry, true, options.PreserveOwner); err != nil {
		_ = destination.Remove(filepath.FromSlash(member.Path))
		return err
	}
	return nil
}

func (s *Service) restoreOwner(root *os.Root, memberPath string, entry ManifestEntry, symlink, preserve bool) error {
	if !preserve {
		return nil
	}
	if entry.UID == nil || entry.GID == nil {
		s.warnOwnerFallback(memberPath, "uid/gid metadata is absent")
		return nil
	}
	uid := int(*entry.UID)
	gid := int(*entry.GID)
	if int64(uid) != *entry.UID || int64(gid) != *entry.GID {
		return fmt.Errorf("uid/gid values are outside the platform range")
	}
	if uid == os.Geteuid() && gid == os.Getegid() {
		return nil
	}
	changeOwner := s.changeOwner
	if changeOwner == nil {
		changeOwner = changeOwnerDefault
	}
	if err := changeOwner(root, filepath.FromSlash(memberPath), uid, gid, symlink); err != nil {
		if os.IsPermission(err) {
			s.warnOwnerFallback(memberPath, err.Error())
			return nil
		}
		return fmt.Errorf("change owner: %w", err)
	}
	return nil
}

func (s *Service) warnOwnerFallback(memberPath, reason string) {
	if s.Log != nil {
		s.Log.Warn("archive ownership not restored; using current owner", "path", memberPath, "reason", reason)
	}
}

func changeOwnerDefault(root *os.Root, memberPath string, uid, gid int, symlink bool) error {
	if symlink {
		return root.Lchown(memberPath, uid, gid)
	}
	return root.Chown(memberPath, uid, gid)
}

func createTempFileInRoot(root *os.Root, memberPath string) (string, *os.File, error) {
	parent := path.Dir(memberPath)
	if parent == "." {
		parent = ""
	}
	for attempt := 0; attempt < 16; attempt++ {
		name, err := randomTempName(parent)
		if err != nil {
			return "", nil, err
		}
		file, err := root.OpenFile(filepath.FromSlash(name), os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0o600)
		if os.IsExist(err) {
			continue
		}
		if err != nil {
			return "", nil, fmt.Errorf("create temporary destination file: %w", err)
		}
		return name, file, nil
	}
	return "", nil, fmt.Errorf("could not allocate temporary destination file for %q", memberPath)
}

func createTempSymlinkInRoot(root *os.Root, memberPath, target string) (string, error) {
	parent := path.Dir(memberPath)
	if parent == "." {
		parent = ""
	}
	for attempt := 0; attempt < 16; attempt++ {
		name, err := randomTempName(parent)
		if err != nil {
			return "", err
		}
		err = root.Symlink(target, filepath.FromSlash(name))
		if os.IsExist(err) {
			continue
		}
		if err != nil {
			return "", fmt.Errorf("create temporary symlink: %w", err)
		}
		return name, nil
	}
	return "", fmt.Errorf("could not allocate temporary symlink for %q", memberPath)
}

func randomTempName(parent string) (string, error) {
	var random [12]byte
	if _, err := rand.Read(random[:]); err != nil {
		return "", fmt.Errorf("generate temporary name: %w", err)
	}
	name := ".dum-tarball-" + hex.EncodeToString(random[:]) + ".tmp"
	return path.Join(parent, name), nil
}
