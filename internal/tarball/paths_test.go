package tarball

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestArchivePathForSource(t *testing.T) {
	root := t.TempDir()
	home := filepath.Join(root, "home")
	if err := os.MkdirAll(home, 0o755); err != nil {
		t.Fatal(err)
	}

	tests := []struct {
		name   string
		source string
		want   string
	}{
		{
			name:   "nested home path",
			source: filepath.Join(home, "projects", "dum"),
			want:   "projects/dum",
		},
		{
			name:   "home root",
			source: home,
			want:   "",
		},
		{
			name:   "external path uses basename",
			source: filepath.Join(root, "external", "Photos"),
			want:   "Photos",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, err := archivePathForSource(home, test.source)
			if !assert.NoError(t, err) {
				return
			}
			assert.Equal(t, test.want, got)
		})
	}
}

func TestValidateMemberPath(t *testing.T) {
	assert.NoError(t, validateMemberPath("Documents/project/file.txt"))

	for _, member := range []string{"", ".", "../escape", "a/../../escape", "/etc/passwd", "a/./b", "a//b", "a/"} {
		t.Run(member, func(t *testing.T) {
			assert.Error(t, validateMemberPath(member))
		})
	}
}

func TestClassifySymlinkTarget(t *testing.T) {
	tests := []struct {
		name       string
		memberPath string
		target     string
		unsafe     bool
		warning    string
	}{
		{name: "relative target stays in archive", memberPath: "Documents/link", target: "settings", unsafe: false},
		{name: "relative parent stays in archive", memberPath: "Documents/link", target: "../Pictures", unsafe: false},
		{
			name:       "absolute target",
			memberPath: "Documents/link",
			target:     "/etc/passwd",
			unsafe:     true,
			warning:    "absolute symlink target",
		},
		{
			name:       "relative target escapes archive",
			memberPath: "link",
			target:     "../../outside",
			unsafe:     true,
			warning:    "symlink target escapes archive root",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			unsafe, warning := classifySymlinkTarget(test.memberPath, test.target)
			assert.Equal(t, test.unsafe, unsafe)
			assert.Equal(t, test.warning, warning)
		})
	}
}

func TestValidateUniqueBasenames(t *testing.T) {
	tests := []struct {
		name    string
		paths   []string
		wantErr bool
	}{
		{name: "unique", paths: []string{"/home/user/Documents", "/mnt/data/Photos"}},
		{name: "duplicate", paths: []string{"/home/user/Shared", "/mnt/data/Shared"}, wantErr: true},
		{name: "case-insensitive duplicate", paths: []string{"/home/user/Photos", "/mnt/data/photos"}, wantErr: true},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := validateUniqueBasenames(test.paths)
			if test.wantErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
		})
	}
}

func TestValidateTargetDoesNotOverlapSources(t *testing.T) {
	root := t.TempDir()
	source := filepath.Join(root, "source")
	if err := os.MkdirAll(source, 0o755); err != nil {
		t.Fatal(err)
	}

	assert.Error(t, validateTargetDoesNotOverlap(filepath.Join(source, "backups"), []string{source}))
	assert.Error(t, validateTargetDoesNotOverlap(source, []string{source}))
	assert.NoError(t, validateTargetDoesNotOverlap(filepath.Join(root, "backups"), []string{source}))
}

func TestValidateTargetDoesNotOverlapSourcesThroughSymlinkedAncestor(t *testing.T) {
	root := t.TempDir()
	physicalSource := filepath.Join(root, "real", "source")
	if err := os.MkdirAll(physicalSource, 0o755); err != nil {
		t.Fatal(err)
	}

	symlinkRoot := filepath.Join(root, "linked")
	if err := os.Symlink(filepath.Join(root, "real"), symlinkRoot); err != nil {
		t.Fatal(err)
	}

	logicalSource := filepath.Join(symlinkRoot, "source")
	target := filepath.Join(physicalSource, "backups")
	assert.Error(t, validateTargetDoesNotOverlap(target, []string{logicalSource}))

	logicalTarget := filepath.Join(symlinkRoot, "source", "backups")
	assert.Error(t, validateTargetDoesNotOverlap(logicalTarget, []string{physicalSource}))
}
