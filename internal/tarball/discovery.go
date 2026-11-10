package tarball

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

// ArchiveCandidates contains the timestamped archives matching one source.
type ArchiveCandidates struct {
	Source   string
	Archives []string
}

// DiscoverArchives finds matching archive files for list-based restore.
func (s *Service) DiscoverArchives(ctx context.Context, sources []string, target string) ([]ArchiveCandidates, error) {
	if ctx == nil {
		return nil, fmt.Errorf("discover archives: nil context")
	}
	if err := ctx.Err(); err != nil {
		return nil, fmt.Errorf("discover archives: %w", err)
	}
	if len(sources) == 0 {
		return nil, fmt.Errorf("archive discovery requires at least one source path")
	}
	if target == "" {
		return nil, fmt.Errorf("archive target cannot be empty")
	}
	if err := validateUniqueBasenames(sources); err != nil {
		return nil, fmt.Errorf("validate restore sources: %w", err)
	}
	absoluteTarget, err := filepath.Abs(target)
	if err != nil {
		return nil, fmt.Errorf("resolve archive target %q: %w", target, err)
	}

	candidates := make([]ArchiveCandidates, len(sources))
	for index, source := range sources {
		candidates[index] = ArchiveCandidates{Source: source}
	}

	entries, err := os.ReadDir(absoluteTarget)
	if os.IsNotExist(err) {
		return candidates, nil
	}
	if err != nil {
		return nil, fmt.Errorf("read archive target %q: %w", absoluteTarget, err)
	}

	for index, source := range sources {
		if err := ctx.Err(); err != nil {
			return nil, fmt.Errorf("discover archives: %w", err)
		}
		base := filepath.Base(filepath.Clean(source))
		prefix := base + "-"
		foldedPrefix := strings.ToLower(prefix)
		for _, entry := range entries {
			if !strings.HasPrefix(strings.ToLower(entry.Name()), foldedPrefix) ||
				!strings.HasSuffix(entry.Name(), ".tgz") {
				continue
			}
			archivePath := filepath.Join(absoluteTarget, entry.Name())
			info, err := os.Lstat(archivePath)
			if os.IsNotExist(err) {
				continue
			}
			if err != nil {
				return nil, fmt.Errorf("inspect archive candidate %q: %w", archivePath, err)
			}
			if !info.Mode().IsRegular() {
				continue
			}
			candidates[index].Archives = append(candidates[index].Archives, archivePath)
		}
		sort.Strings(candidates[index].Archives)
	}
	return candidates, nil
}
