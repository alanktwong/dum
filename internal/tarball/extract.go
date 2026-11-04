package tarball

import (
	"archive/tar"
	"compress/gzip"
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"path"
	"path/filepath"
	"strings"
)

type archiveMember struct {
	Path   string
	Type   EntryType
	Header tar.Header
}

func inspectArchive(snapshotPath string, manifest Manifest, allowUnsafeLinks bool) ([]archiveMember, error) {
	file, err := os.Open(snapshotPath)
	if err != nil {
		return nil, fmt.Errorf("open verified snapshot: %w", err)
	}
	defer func() { _ = file.Close() }()

	compressed, err := gzip.NewReader(file)
	if err != nil {
		return nil, fmt.Errorf("open gzip stream: %w", err)
	}
	defer func() { _ = compressed.Close() }()

	expected := make(map[string]ManifestEntry, len(manifest.Entries))
	for _, entry := range manifest.Entries {
		expected[entry.Path] = entry
	}
	seen := make(map[string]struct{}, len(manifest.Entries))
	caseFolded := make(map[string]string, len(manifest.Entries))
	actual := make([]archiveMember, 0, len(manifest.Entries))
	reader := tar.NewReader(compressed)
	for {
		header, err := reader.Next()
		if errors.Is(err, io.EOF) {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("read tar header: %w", err)
		}

		entryType, err := archiveEntryType(header.Typeflag)
		if err != nil {
			return nil, fmt.Errorf("archive member %q: %w", header.Name, err)
		}
		memberPath := header.Name
		if entryType == EntryTypeDirectory {
			memberPath = strings.TrimSuffix(memberPath, "/")
		}
		if err := validateMemberPath(memberPath); err != nil {
			return nil, fmt.Errorf("archive member %q: %w", header.Name, err)
		}
		if _, ok := seen[memberPath]; ok {
			return nil, fmt.Errorf("duplicate archive member path %q", memberPath)
		}
		seen[memberPath] = struct{}{}
		folded := strings.ToLower(memberPath)
		if previous, ok := caseFolded[folded]; ok && previous != memberPath {
			return nil, fmt.Errorf("case-insensitive archive path collision between %q and %q", previous, memberPath)
		}
		caseFolded[folded] = memberPath

		if header.Uid < 0 || header.Gid < 0 {
			return nil, fmt.Errorf("archive member %q has negative uid or gid", memberPath)
		}
		if entryType != EntryTypeFile && header.Size != 0 {
			return nil, fmt.Errorf("non-file archive member %q has size %d", memberPath, header.Size)
		}
		if entryType == EntryTypeFile && header.Size < 0 {
			return nil, fmt.Errorf("archive member %q has negative size", memberPath)
		}
		if entryType != EntryTypeSymlink && header.Linkname != "" {
			return nil, fmt.Errorf("non-symlink archive member %q has a link target", memberPath)
		}

		manifestEntry, ok := expected[memberPath]
		if !ok {
			return nil, fmt.Errorf("archive member %q is missing from manifest", memberPath)
		}
		if manifestEntry.Type != entryType {
			return nil, fmt.Errorf("archive member %q type does not match manifest", memberPath)
		}
		if manifestEntry.UID == nil {
			if header.Uid != 0 || header.Gid != 0 {
				return nil, fmt.Errorf("archive member %q has owner ids omitted from manifest", memberPath)
			}
		} else if int64(header.Uid) != *manifestEntry.UID || int64(header.Gid) != *manifestEntry.GID {
			return nil, fmt.Errorf("archive member %q owner ids do not match manifest", memberPath)
		}
		if entryType == EntryTypeSymlink {
			if manifestEntry.Target != header.Linkname {
				return nil, fmt.Errorf("archive member %q link target does not match manifest", memberPath)
			}
			unsafe, warning := classifySymlinkTarget(memberPath, header.Linkname)
			if unsafe != manifestEntry.Unsafe || (unsafe && !containsUnsafeWarning(manifestEntry.Warning, warning)) {
				return nil, fmt.Errorf("archive member %q symlink safety metadata does not match tar header", memberPath)
			}
			if unsafe && !allowUnsafeLinks {
				return nil, fmt.Errorf("archive contains unsafe symlink %q; pass --allow-unsafe-links to restore it", memberPath)
			}
		}

		headerCopy := *header
		actual = append(actual, archiveMember{Path: memberPath, Type: entryType, Header: headerCopy})
	}
	if _, err := io.Copy(io.Discard, compressed); err != nil {
		return nil, fmt.Errorf("verify gzip trailer: %w", err)
	}
	if len(actual) != len(manifest.Entries) {
		return nil, fmt.Errorf("archive has %d entries; manifest lists %d", len(actual), len(manifest.Entries))
	}
	if err := validateArchiveParents(actual); err != nil {
		return nil, err
	}
	return actual, nil
}

func archiveEntryType(typeflag byte) (EntryType, error) {
	switch typeflag {
	case tar.TypeReg, tar.TypeRegA:
		return EntryTypeFile, nil
	case tar.TypeDir:
		return EntryTypeDirectory, nil
	case tar.TypeSymlink:
		return EntryTypeSymlink, nil
	default:
		return "", fmt.Errorf("unsupported tar type %d", typeflag)
	}
}

func validateArchiveParents(members []archiveMember) error {
	byPath := make(map[string]EntryType, len(members))
	for _, member := range members {
		byPath[member.Path] = member.Type
	}
	for _, member := range members {
		parent := path.Dir(member.Path)
		for parent != "." && parent != "/" {
			if parentType, ok := byPath[parent]; ok && parentType != EntryTypeDirectory {
				return fmt.Errorf("archive member %q traverses non-directory member %q", member.Path, parent)
			}
			parent = path.Dir(parent)
		}
	}
	return nil
}

func extractArchiveToStage(ctx context.Context, snapshotPath, stageDir string, members []archiveMember) error {
	stageRoot, err := os.OpenRoot(stageDir)
	if err != nil {
		return fmt.Errorf("open staging root: %w", err)
	}
	defer func() { _ = stageRoot.Close() }()

	for _, member := range members {
		if err := ctx.Err(); err != nil {
			return fmt.Errorf("extract archive: %w", err)
		}
		if member.Type == EntryTypeDirectory {
			if err := makeStageDirectory(stageRoot, member.Path); err != nil {
				return err
			}
		}
		if err := makeStageParents(stageRoot, member.Path); err != nil {
			return err
		}
	}

	file, err := os.Open(snapshotPath)
	if err != nil {
		return fmt.Errorf("open verified snapshot for extraction: %w", err)
	}
	defer func() { _ = file.Close() }()
	compressed, err := gzip.NewReader(file)
	if err != nil {
		return fmt.Errorf("open gzip stream for extraction: %w", err)
	}
	defer func() { _ = compressed.Close() }()
	reader := tar.NewReader(compressed)

	symlinks := make([]archiveMember, 0)
	index := 0
	for {
		header, err := reader.Next()
		if errors.Is(err, io.EOF) {
			break
		}
		if err != nil {
			return fmt.Errorf("read archive while extracting: %w", err)
		}
		if index >= len(members) {
			return fmt.Errorf("archive gained an unexpected member during extraction")
		}
		member := members[index]
		index++
		memberPath := header.Name
		if member.Type == EntryTypeDirectory {
			memberPath = strings.TrimSuffix(memberPath, "/")
		}
		if memberPath != member.Path {
			return fmt.Errorf("archive changed between validation and extraction at %q", member.Path)
		}
		if err := ctx.Err(); err != nil {
			return fmt.Errorf("extract archive: %w", err)
		}

		switch member.Type {
		case EntryTypeDirectory:
			continue
		case EntryTypeSymlink:
			symlinks = append(symlinks, member)
		case EntryTypeFile:
			if err := extractRegularFile(stageRoot, member.Path, reader, header.Size); err != nil {
				return fmt.Errorf("extract archive member %q: %w", member.Path, err)
			}
		default:
			return fmt.Errorf("unexpected archive member type %q", member.Type)
		}
	}
	if index != len(members) {
		return fmt.Errorf("archive ended after %d of %d validated members", index, len(members))
	}
	if _, err := io.Copy(io.Discard, compressed); err != nil {
		return fmt.Errorf("verify gzip trailer during extraction: %w", err)
	}
	for _, member := range symlinks {
		targetPath := filepath.Join(stageDir, filepath.FromSlash(member.Path))
		if err := os.Symlink(member.Header.Linkname, targetPath); err != nil {
			return fmt.Errorf("create staged symlink %q: %w", member.Path, err)
		}
	}
	return nil
}

func makeStageDirectory(stageRoot *os.Root, memberPath string) error {
	if err := stageRoot.MkdirAll(filepath.FromSlash(memberPath), 0o700); err != nil {
		return fmt.Errorf("create staged directory %q: %w", memberPath, err)
	}
	return nil
}

func makeStageParents(stageRoot *os.Root, memberPath string) error {
	parent := path.Dir(memberPath)
	if parent == "." {
		return nil
	}
	if err := stageRoot.MkdirAll(filepath.FromSlash(parent), 0o700); err != nil {
		return fmt.Errorf("create staged parent for %q: %w", memberPath, err)
	}
	return nil
}

func extractRegularFile(stageRoot *os.Root, memberPath string, source io.Reader, size int64) error {
	if size < 0 {
		return fmt.Errorf("negative file size %d", size)
	}
	target, err := stageRoot.OpenFile(filepath.FromSlash(memberPath), os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0o600)
	if err != nil {
		return fmt.Errorf("create staged file: %w", err)
	}
	written, copyErr := io.CopyN(target, source, size)
	closeErr := target.Close()
	if err := errors.Join(copyErr, closeErr); err != nil {
		return err
	}
	if written != size {
		return fmt.Errorf("copied %d bytes, expected %d", written, size)
	}
	return nil
}
