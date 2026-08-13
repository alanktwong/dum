package rename

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func setupMockLogger(t *testing.T) *mockLogger {
	t.Helper()
	logger := &mockLogger{}
	logger.Test(t)

	// Debugf and Infof have variable argument counts depending on code path.
	// Use multiple On() calls with increasing mock.Anything args to cover all paths.
	logger.On("Debugf", mock.Anything, mock.Anything, mock.Anything).Return()
	logger.On("Debugf", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return()
	logger.On("Debugf", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return()
	logger.On("Infof", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return()

	return logger
}

func TestRename_BasicRename(t *testing.T) {
	tmpDir := t.TempDir()
	oldPath := filepath.Join(tmpDir, "source.txt")
	newPath := filepath.Join(tmpDir, "renamed.txt")

	err := os.WriteFile(oldPath, []byte("hello"), 0o644)
	assert.NoError(t, err)

	dum := &Command{Log: setupMockLogger(t)}
	err = dum.rename(context.Background(), Options{
		FilePattern: oldPath,
		Source:      "source",
		Replace:     "renamed",
	})
	assert.NoError(t, err)

	_, err = os.Stat(oldPath)
	assert.True(t, os.IsNotExist(err), "old file should not exist")

	_, err = os.Stat(newPath)
	assert.NoError(t, err, "new file should exist")
}

func TestRename_TrimFront(t *testing.T) {
	tmpDir := t.TempDir()
	oldPath := filepath.Join(tmpDir, "abcdef.txt")
	newPath := filepath.Join(tmpDir, "cdef.txt")

	err := os.WriteFile(oldPath, []byte("hello"), 0o644)
	assert.NoError(t, err)

	dum := &Command{Log: setupMockLogger(t)}
	err = dum.rename(context.Background(), Options{
		FilePattern: oldPath,
		TrimFront:   2,
	})
	assert.NoError(t, err)

	_, err = os.Stat(oldPath)
	assert.True(t, os.IsNotExist(err))

	_, err = os.Stat(newPath)
	assert.NoError(t, err)
}

func TestRename_TrimBack(t *testing.T) {
	tmpDir := t.TempDir()
	oldPath := filepath.Join(tmpDir, "abcdef.txt")
	newPath := filepath.Join(tmpDir, "abcd.txt")

	err := os.WriteFile(oldPath, []byte("hello"), 0o644)
	assert.NoError(t, err)

	dum := &Command{Log: setupMockLogger(t)}
	err = dum.rename(context.Background(), Options{
		FilePattern: oldPath,
		TrimBack:    2,
	})
	assert.NoError(t, err)

	_, err = os.Stat(oldPath)
	assert.True(t, os.IsNotExist(err))

	_, err = os.Stat(newPath)
	assert.NoError(t, err)
}

func TestRename_Replace(t *testing.T) {
	tmpDir := t.TempDir()
	oldPath := filepath.Join(tmpDir, "_DSC1234.NEF")
	newPath := filepath.Join(tmpDir, "bar1234.NEF")

	err := os.WriteFile(oldPath, []byte("hello"), 0o644)
	assert.NoError(t, err)

	dum := &Command{Log: setupMockLogger(t)}
	err = dum.rename(context.Background(), Options{
		FilePattern: oldPath,
		Source:      "_DSC",
		Replace:     "bar",
	})
	assert.NoError(t, err)

	_, err = os.Stat(oldPath)
	assert.True(t, os.IsNotExist(err))

	_, err = os.Stat(newPath)
	assert.NoError(t, err)
}

func TestRename_Lowercase(t *testing.T) {
	tmpDir := t.TempDir()
	oldPath := filepath.Join(tmpDir, "FILE.JPG")

	err := os.WriteFile(oldPath, []byte("hello"), 0o644)
	assert.NoError(t, err)

	dum := &Command{Log: setupMockLogger(t)}
	err = dum.rename(context.Background(), Options{
		FilePattern: oldPath,
		Lowercase:   true,
	})
	assert.NoError(t, err)

	// On case-insensitive filesystems (macOS), os.Stat(oldPath) still succeeds
	// because "FILE.JPG" and "file.jpg" refer to the same inode. Verify by
	// reading the directory listing instead.
	entries, err := os.ReadDir(tmpDir)
	assert.NoError(t, err)
	assert.Len(t, entries, 1)
	assert.Equal(t, "file.jpg", entries[0].Name())
}

func TestRename_Iterate(t *testing.T) {
	tmpDir := t.TempDir()
	oldPath := filepath.Join(tmpDir, "photo.jpg")
	newPath := filepath.Join(tmpDir, "photo_1.jpg")

	err := os.WriteFile(oldPath, []byte("hello"), 0o644)
	assert.NoError(t, err)

	dum := &Command{Log: setupMockLogger(t)}
	err = dum.rename(context.Background(), Options{
		FilePattern: oldPath,
		Index:       0,
		Iterate:     1,
	})
	assert.NoError(t, err)

	_, err = os.Stat(oldPath)
	assert.True(t, os.IsNotExist(err))

	_, err = os.Stat(newPath)
	assert.NoError(t, err)
}

func TestRename_Combined(t *testing.T) {
	tmpDir := t.TempDir()
	oldPath := filepath.Join(tmpDir, "_DSC1234.NEF")
	newPath := filepath.Join(tmpDir, "bar1234.nef")

	err := os.WriteFile(oldPath, []byte("hello"), 0o644)
	assert.NoError(t, err)

	dum := &Command{Log: setupMockLogger(t)}
	err = dum.rename(context.Background(), Options{
		FilePattern: oldPath,
		Source:      "_DSC",
		Replace:     "bar",
		Lowercase:   true,
	})
	assert.NoError(t, err)

	_, err = os.Stat(oldPath)
	assert.True(t, os.IsNotExist(err))

	_, err = os.Stat(newPath)
	assert.NoError(t, err)
}

func TestRename_DryRun(t *testing.T) {
	tmpDir := t.TempDir()
	oldPath := filepath.Join(tmpDir, "old.txt")

	err := os.WriteFile(oldPath, []byte("hello"), 0o644)
	assert.NoError(t, err)

	dum := &Command{Log: setupMockLogger(t)}
	err = dum.rename(context.Background(), Options{
		FilePattern: oldPath,
		DryRun:      true,
	})
	assert.NoError(t, err)

	_, err = os.Stat(oldPath)
	assert.NoError(t, err, "file should still exist after dry run")
}

func TestRename_OSError(t *testing.T) {
	dum := &Command{Log: setupMockLogger(t)}
	err := dum.rename(context.Background(), Options{
		FilePattern: filepath.Join(t.TempDir(), "nonexistent.txt"),
	})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "error renaming file")
}

func TestNewRenameCommand_RunE_SourceWithoutReplace(t *testing.T) {
	logger := &mockLogger{}
	logger.Test(t)
	logger.On("SetLevel", mock.Anything).Return()
	logger.On("Debugf", mock.Anything, mock.Anything).Return()
	logger.On("Debugf", mock.Anything, mock.Anything, mock.Anything).Return()
	logger.On("Debug",
		mock.Anything, mock.Anything, mock.Anything, mock.Anything,
		mock.Anything, mock.Anything, mock.Anything, mock.Anything,
		mock.Anything, mock.Anything, mock.Anything, mock.Anything,
		mock.Anything, mock.Anything, mock.Anything, mock.Anything,
		mock.Anything, mock.Anything, mock.Anything,
	).Return()

	dum := &Command{Log: logger}
	cmd := NewRenameCommand("dum", dum)
	cmd.SetContext(context.Background())
	cmd.SetArgs([]string{"--source", "_DSC", "foo.txt"})

	err := cmd.Execute()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "cannot replace _DSC without replace")
}

func TestNewRenameCommand_RunE_ReplaceWithoutSource(t *testing.T) {
	logger := &mockLogger{}
	logger.Test(t)
	logger.On("SetLevel", mock.Anything).Return()
	logger.On("Debugf", mock.Anything, mock.Anything).Return()
	logger.On("Debugf", mock.Anything, mock.Anything, mock.Anything).Return()
	logger.On("Debug",
		mock.Anything, mock.Anything, mock.Anything, mock.Anything,
		mock.Anything, mock.Anything, mock.Anything, mock.Anything,
		mock.Anything, mock.Anything, mock.Anything, mock.Anything,
		mock.Anything, mock.Anything, mock.Anything, mock.Anything,
		mock.Anything, mock.Anything, mock.Anything,
	).Return()

	dum := &Command{Log: logger}
	cmd := NewRenameCommand("dum", dum)
	cmd.SetContext(context.Background())
	cmd.SetArgs([]string{"--replace", "bar", "foo.txt"})

	err := cmd.Execute()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "cannot replace for bar without source")
}
