package rename

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	lg "awong/dotfiles/internal/logging"
	clog "github.com/charmbracelet/log"
	"github.com/stretchr/testify/assert"
)

func newTestRenamer() *Renamer {
	return NewRenamer(lg.NewLogger(lg.Options{Level: clog.ErrorLevel}))
}

func TestRenamer_RenameTransformsFilename(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		output  string
		options Options
	}{
		{name: "basic", input: "source.txt", output: "renamed.txt", options: Options{Source: "source", Replace: "renamed"}},
		{name: "trim front", input: "abcdef.txt", output: "cdef.txt", options: Options{TrimFront: 2}},
		{name: "trim back", input: "abcdef.txt", output: "abcd.txt", options: Options{TrimBack: 2}},
		{name: "lowercase", input: "FILE.JPG", output: "file.jpg", options: Options{Lowercase: true}},
		{name: "iterate", input: "photo.jpg", output: "photo_3.jpg", options: Options{Index: 2, Iterate: 1}},
		{name: "combined", input: "_DSC1234.NEF", output: "bar1234.nef", options: Options{Source: "_DSC", Replace: "bar", Lowercase: true}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dir := t.TempDir()
			oldPath := filepath.Join(dir, tt.input)
			newPath := filepath.Join(dir, tt.output)
			assert.NoError(t, os.WriteFile(oldPath, []byte("content"), 0o644))

			tt.options.FilePattern = oldPath
			err := newTestRenamer().Rename(context.Background(), tt.options)

			assert.NoError(t, err)
			_, err = os.Stat(newPath)
			assert.NoError(t, err)
		})
	}
}

func TestRenamer_RenameDryRunPreservesFile(t *testing.T) {
	dir := t.TempDir()
	oldPath := filepath.Join(dir, "old.txt")
	newPath := filepath.Join(dir, "new.txt")
	assert.NoError(t, os.WriteFile(oldPath, []byte("content"), 0o644))

	err := newTestRenamer().Rename(context.Background(), Options{
		FilePattern: oldPath,
		Source:      "old",
		Replace:     "new",
		DryRun:      true,
	})

	assert.NoError(t, err)
	_, err = os.Stat(oldPath)
	assert.NoError(t, err)
	_, err = os.Stat(newPath)
	assert.ErrorIs(t, err, os.ErrNotExist)
}

func TestRenamer_RenameReturnsFilesystemError(t *testing.T) {
	err := newTestRenamer().Rename(context.Background(), Options{
		FilePattern: filepath.Join(t.TempDir(), "missing.txt"),
	})

	assert.Error(t, err)
	assert.ErrorContains(t, err, "error renaming file")
}
