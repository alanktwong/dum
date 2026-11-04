package tarball

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strings"
)

func archivePathForSource(home, source string) (string, error) {
	absoluteHome, err := filepath.Abs(home)
	if err != nil {
		return "", fmt.Errorf("resolve home path %q: %w", home, err)
	}
	absoluteSource, err := filepath.Abs(source)
	if err != nil {
		return "", fmt.Errorf("resolve source path %q: %w", source, err)
	}

	relative, err := filepath.Rel(absoluteHome, absoluteSource)
	if err != nil {
		return "", fmt.Errorf("make source path %q relative to home: %w", source, err)
	}
	if relative == "." {
		return "", nil
	}
	if relative != ".." && !strings.HasPrefix(relative, ".."+string(filepath.Separator)) {
		return filepath.ToSlash(relative), nil
	}

	base := filepath.Base(absoluteSource)
	if base == "" || base == "." || base == string(filepath.Separator) {
		return "", fmt.Errorf("cannot determine archive root for source %q", source)
	}
	return base, nil
}

func classifySymlinkTarget(memberPath, target string) (bool, string) {
	if path.IsAbs(target) {
		return true, "absolute symlink target"
	}

	resolved := path.Clean(path.Join(path.Dir(memberPath), target))
	if resolved == ".." || strings.HasPrefix(resolved, "../") {
		return true, "symlink target escapes archive root"
	}
	return false, ""
}

func validateUniqueBasenames(paths []string) error {
	seen := make(map[string]string, len(paths))
	for _, source := range paths {
		base := filepath.Base(filepath.Clean(source))
		if base == "" || base == "." || base == string(filepath.Separator) {
			return fmt.Errorf("cannot determine source basename for %q", source)
		}
		folded := strings.ToLower(base)
		if previous, ok := seen[folded]; ok {
			return fmt.Errorf("duplicate source basename %q conflicts with %q", base, previous)
		}
		seen[folded] = source
	}
	return nil
}

func validateTargetDoesNotOverlap(target string, sources []string) error {
	absoluteTarget, err := canonicalPath(target)
	if err != nil {
		return fmt.Errorf("resolve target path %q: %w", target, err)
	}

	for _, source := range sources {
		absoluteSource, err := canonicalPath(source)
		if err != nil {
			return fmt.Errorf("resolve source path %q: %w", source, err)
		}
		relativeTarget, err := filepath.Rel(absoluteSource, absoluteTarget)
		if err != nil {
			return fmt.Errorf("compare target %q with source %q: %w", target, source, err)
		}
		if relativeTarget == "." || (relativeTarget != ".." && !strings.HasPrefix(relativeTarget, ".."+string(filepath.Separator))) {
			return fmt.Errorf("target %q overlaps source %q", target, source)
		}
	}
	return nil
}

func canonicalPath(input string) (string, error) {
	absolute, err := filepath.Abs(input)
	if err != nil {
		return "", fmt.Errorf("resolve absolute path: %w", err)
	}
	cursor := filepath.Clean(absolute)
	var suffix []string

	for {
		_, err := os.Lstat(cursor)
		if err == nil {
			resolved, err := filepath.EvalSymlinks(cursor)
			if err != nil {
				return "", fmt.Errorf("resolve symlinks for %q: %w", cursor, err)
			}
			if len(suffix) > 0 {
				info, err := os.Stat(resolved)
				if err != nil {
					return "", fmt.Errorf("inspect existing parent %q: %w", resolved, err)
				}
				if !info.IsDir() {
					return "", fmt.Errorf("existing parent %q is not a directory", resolved)
				}
			}
			for index := len(suffix) - 1; index >= 0; index-- {
				resolved = filepath.Join(resolved, suffix[index])
			}
			return filepath.Clean(resolved), nil
		}
		if !os.IsNotExist(err) {
			return "", fmt.Errorf("inspect path %q: %w", cursor, err)
		}

		parent := filepath.Dir(cursor)
		if parent == cursor {
			return "", fmt.Errorf("cannot find existing parent for %q", input)
		}
		suffix = append(suffix, filepath.Base(cursor))
		cursor = parent
	}
}
