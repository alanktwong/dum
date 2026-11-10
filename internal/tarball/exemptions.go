package tarball

import (
	"bufio"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strings"
)

// compiledPattern is a single compiled exemption pattern with matching rules.
type compiledPattern struct {
	pattern  string // original pattern string
	matcher  func(string) bool
	dirOnly  bool // true when the pattern ends with '/' (matches directories only)
	anchored bool // true when the pattern contains '/' (anchored to source root)
}

// ExemptionMatcher holds the compiled exemption patterns for one backup source.
type ExemptionMatcher struct {
	patterns     []compiledPattern
	useGitignore bool
}

// NewExemptionMatcher compiles raw glob patterns into a matcher. Patterns are
// validated up front; an invalid glob returns an error before any walk begins.
func NewExemptionMatcher(patterns []string, useGitignore bool) (*ExemptionMatcher, error) {
	compiled := make([]compiledPattern, 0, len(patterns))
	for _, pattern := range patterns {
		cp, err := compilePattern(pattern)
		if err != nil {
			return nil, fmt.Errorf("invalid exclude pattern %q: %w", pattern, err)
		}
		compiled = append(compiled, cp)
	}
	return &ExemptionMatcher{patterns: compiled, useGitignore: useGitignore}, nil
}

// compilePattern turns a raw glob string into a compiledPattern.
func compilePattern(pattern string) (compiledPattern, error) {
	if pattern == "" {
		return compiledPattern{}, fmt.Errorf("empty pattern")
	}

	dirOnly := strings.HasSuffix(pattern, "/")
	patternToMatch := pattern
	if dirOnly {
		patternToMatch = strings.TrimSuffix(pattern, "/")
	}

	anchored := strings.Contains(patternToMatch, "/")

	var matcher func(string) bool
	if anchored {
		// Anchored patterns match against the full relative path.
		matchFn, err := pathMatch(patternToMatch)
		if err != nil {
			return compiledPattern{}, err
		}
		matcher = matchFn
	} else {
		// Unanchored patterns match any path segment's basename at any depth.
		baseMatch, err := pathMatch(patternToMatch)
		if err != nil {
			return compiledPattern{}, err
		}
		matcher = func(relPath string) bool {
			// Split the path into segments and test each segment's basename.
			for _, segment := range strings.Split(relPath, "/") {
				if baseMatch(segment) {
					return true
				}
			}
			return false
		}
	}

	return compiledPattern{
		pattern:  pattern,
		matcher:  matcher,
		dirOnly:  dirOnly,
		anchored: anchored,
	}, nil
}

// pathMatch compiles a glob using path.Match semantics.
func pathMatch(pattern string) (func(string) bool, error) {
	if _, err := path.Match(pattern, ""); err != nil {
		return nil, fmt.Errorf("invalid glob pattern %q: %w", pattern, err)
	}
	return func(candidate string) bool {
		matched, err := path.Match(pattern, candidate)
		if err != nil {
			return false
		}
		return matched
	}, nil
}

// Match reports whether the given relative path (slash-separated, relative to
// the source root) is exempted. When isDir is true, directory-only patterns
// also apply; when false, directory-only patterns are ignored.
func (m *ExemptionMatcher) Match(relPath string, isDir bool) bool {
	for _, cp := range m.patterns {
		if cp.dirOnly && !isDir {
			continue
		}
		if cp.matcher(relPath) {
			return true
		}
	}
	return false
}

// UseGitignore reports whether .gitignore patterns were merged into this
// matcher.
func (m *ExemptionMatcher) UseGitignore() bool {
	return m.useGitignore
}

// Patterns returns the original raw patterns.
func (m *ExemptionMatcher) Patterns() []string {
	patterns := make([]string, 0, len(m.patterns))
	for _, cp := range m.patterns {
		patterns = append(patterns, cp.pattern)
	}
	return patterns
}

// LoadGitignorePatterns reads a .gitignore file from the given directory and
// returns its patterns. Comments, blank lines, and the negation prefix (!)
// are handled: comments and blanks are skipped, negation prefixes are
// stripped (negation itself is not supported in v1 and is rejected).
// Patterns are returned relative to the .gitignore file's own directory.
func LoadGitignorePatterns(gitignorePath string) ([]string, error) {
	file, err := os.Open(gitignorePath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, fmt.Errorf("read .gitignore %q: %w", gitignorePath, err)
	}
	defer func() { _ = file.Close() }()

	var patterns []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		if strings.HasPrefix(line, "!") {
			return nil, fmt.Errorf("negation patterns are not supported in .gitignore %q", gitignorePath)
		}
		patterns = append(patterns, line)
	}
	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("read .gitignore %q: %w", gitignorePath, err)
	}
	return patterns, nil
}

// MergeGitignore loads .gitignore from each source root and merges its
// patterns into the matcher. Nested .gitignore files are resolved relative
// to their own directory. This is a no-op when useGitignore is false.
func (m *ExemptionMatcher) MergeGitignore(sourceRoots []string) error {
	if !m.useGitignore {
		return nil
	}
	seen := make(map[string]struct{})
	for _, root := range sourceRoots {
		gitignorePath := filepath.Join(root, ".gitignore")
		patterns, err := LoadGitignorePatterns(gitignorePath)
		if err != nil {
			return err
		}
		for _, pattern := range patterns {
			if _, ok := seen[pattern]; ok {
				continue
			}
			seen[pattern] = struct{}{}
			cp, err := compilePattern(pattern)
			if err != nil {
				return fmt.Errorf("invalid pattern in .gitignore %q: %w", gitignorePath, err)
			}
			m.patterns = append(m.patterns, cp)
		}
	}
	return nil
}
