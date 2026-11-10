package tarball

import (
	"context"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestRestore_RoundTripsBackupArchives(t *testing.T) {
	root := t.TempDir()
	home := filepath.Join(root, "home")
	inside := filepath.Join(home, "Documents")
	external := filepath.Join(root, "external", "Photos")
	writeTestFile(t, filepath.Join(inside, "settings.json"), []byte("{}"), 0o600, time.Time{})
	if err := os.Symlink("settings.json", filepath.Join(inside, "settings-link")); err != nil {
		t.Fatal(err)
	}
	writeTestFile(t, filepath.Join(external, "photo.bin"), []byte("photo"), 0o644, time.Time{})

	service := NewService(nil)
	service.now = func() time.Time { return time.Date(2026, 11, 3, 10, 11, 12, 0, time.Local) }
	backups, err := service.Backup(context.Background(), BackupOptions{
		Paths: []string{inside, external}, Target: filepath.Join(root, "archives"), Home: home,
	})
	if !assert.NoError(t, err) {
		return
	}

	destination := filepath.Join(root, "restore")
	for _, artifact := range backups.Artifacts {
		if _, err := service.Restore(context.Background(), RestoreOptions{
			Archive: artifact.ArchivePath, Destination: destination,
		}); !assert.NoError(t, err) {
			return
		}
	}

	content, err := os.ReadFile(filepath.Join(destination, "Documents", "settings.json"))
	if !assert.NoError(t, err) {
		return
	}
	assert.Equal(t, []byte("{}"), content)
	assert.FileExists(t, filepath.Join(destination, "Photos", "photo.bin"))
	linkPath := filepath.Join(destination, "Documents", "settings-link")
	linkTarget, err := os.Readlink(linkPath)
	if assert.NoError(t, err) {
		assert.Equal(t, "settings.json", linkTarget)
	}
}
