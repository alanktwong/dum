package tarball

import (
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestExemptionMatcher_CompileAndMatch(t *testing.T) {
	t.Parallel()

	// Unanchored pattern matches basename at any depth.
	m, err := NewExemptionMatcher([]string{"node_modules"}, false)
	assert.NoError(t, err)
	assert.True(t, m.Match("node_modules", true))
	assert.True(t, m.Match("app/node_modules", true))
	assert.True(t, m.Match("a/b/node_modules", true))
	assert.False(t, m.Match("not_node_modules", true))
	assert.False(t, m.Match("my_node_modules", true))

	// Unanchored pattern with glob.
	m, err = NewExemptionMatcher([]string{"*.log"}, false)
	assert.NoError(t, err)
	assert.True(t, m.Match("debug.log", false))
	assert.True(t, m.Match("sub/debug.log", false))
	assert.False(t, m.Match("debug.txt", false))

	// Anchored pattern matches from source root only.
	m, err = NewExemptionMatcher([]string{"build/output"}, false)
	assert.NoError(t, err)
	assert.True(t, m.Match("build/output", true))
	assert.True(t, m.Match("build/output", false))
	assert.False(t, m.Match("build/output/file.txt", false)) // children handled by SkipDir
	assert.False(t, m.Match("other/build/output", true))
	assert.False(t, m.Match("src/build/output", true))

	// Directory-only trailing slash matches directories, not files.
	m, err = NewExemptionMatcher([]string{"dist/"}, false)
	assert.NoError(t, err)
	assert.True(t, m.Match("dist", true))
	assert.False(t, m.Match("dist", false))          // file named "dist" is not exempted
	assert.False(t, m.Match("dist/file.txt", false)) // dir match + SkipDir handles children

	// Multiple patterns.
	m, err = NewExemptionMatcher([]string{"node_modules", "*.tmp", "secrets/"}, false)
	assert.NoError(t, err)
	assert.True(t, m.Match("node_modules", true))
	assert.True(t, m.Match("app/node_modules", true))
	assert.True(t, m.Match("cache.tmp", false))
	assert.True(t, m.Match("sub/cache.tmp", false))
	assert.True(t, m.Match("secrets", true))
	assert.False(t, m.Match("secrets.txt", false))
}

func TestExemptionMatcher_InvalidPatterns(t *testing.T) {
	t.Parallel()

	_, err := NewExemptionMatcher([]string{"[unclosed"}, false)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid exclude pattern")

	_, err = NewExemptionMatcher([]string{""}, false)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "empty pattern")
}

func TestExemptionMatcher_PatternsAndGitignore(t *testing.T) {
	t.Parallel()

	m, err := NewExemptionMatcher([]string{"*.log"}, true)
	assert.NoError(t, err)
	assert.True(t, m.UseGitignore())
	assert.Equal(t, []string{"*.log"}, m.Patterns())

	m2, err := NewExemptionMatcher([]string{}, false)
	assert.NoError(t, err)
	assert.False(t, m2.UseGitignore())
	assert.Empty(t, m2.Patterns())
}

func TestLoadGitignorePatterns(t *testing.T) {
	t.Parallel()

	// Non-existent file returns nil, nil.
	patterns, err := LoadGitignorePatterns("/nonexistent/path/.gitignore")
	assert.NoError(t, err)
	assert.Nil(t, patterns)

	// Valid .gitignore: comments and blanks skipped.
	dir := t.TempDir()
	writeTestFile(t, filepath.Join(dir, ".gitignore"), []byte(
		"# comment\n"+
			"\n"+
			"node_modules\n"+
			"  dist/  \n"+
			"# another comment\n"+
			"*.log\n",
	), 0o644, time.Now())

	patterns, err = LoadGitignorePatterns(filepath.Join(dir, ".gitignore"))
	assert.NoError(t, err)
	assert.Equal(t, []string{"node_modules", "dist/", "*.log"}, patterns)

	// Negation prefix is rejected.
	writeTestFile(t, filepath.Join(dir, "bad.gitignore"), []byte("!important.txt\n"), 0o644, time.Now())
	_, err = LoadGitignorePatterns(filepath.Join(dir, "bad.gitignore"))
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "negation patterns are not supported")
}

func TestExemptionMatcher_MergeGitignore(t *testing.T) {
	t.Parallel()

	root := t.TempDir()
	writeTestFile(t, filepath.Join(root, ".gitignore"), []byte("node_modules\n*.log\n"), 0o644, time.Now())

	m, err := NewExemptionMatcher([]string{"dist/"}, true)
	assert.NoError(t, err)
	err = m.MergeGitignore([]string{root})
	assert.NoError(t, err)

	assert.True(t, m.Match("node_modules", true))
	assert.True(t, m.Match("sub/node_modules", true))
	assert.True(t, m.Match("app.log", false))
	assert.True(t, m.Match("dist", true))
	assert.False(t, m.Match("keep.txt", false))
}

func TestExemptionMatcher_MergeGitignore_Disabled(t *testing.T) {
	t.Parallel()

	root := t.TempDir()
	writeTestFile(t, filepath.Join(root, ".gitignore"), []byte("node_modules\n"), 0o644, time.Now())

	m, err := NewExemptionMatcher([]string{}, false)
	assert.NoError(t, err)
	err = m.MergeGitignore([]string{root})
	assert.NoError(t, err)
	assert.False(t, m.Match("node_modules", true))
}
