package tarball

import (
	"archive/tar"
	"context"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"syscall"
	"testing"
	"time"

	gen "github.com/alanktwong/dum/internal/logging/gen"
	"github.com/stretchr/testify/assert"
	mck "github.com/stretchr/testify/mock"
)

func TestRestore_ConflictsWithoutOverwrite(t *testing.T) {
	root := t.TempDir()
	artifacts := writeValidRestoreFixture(t, filepath.Join(root, "archive.tgz"))
	destination := filepath.Join(root, "destination")
	filePath := filepath.Join(destination, "Documents", "file.txt")
	if err := os.MkdirAll(filepath.Dir(filePath), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filePath, []byte("keep me"), 0o600); err != nil {
		t.Fatal(err)
	}

	_, err := NewService(nil).Restore(context.Background(), RestoreOptions{
		Archive:     artifacts.ArchivePath,
		Destination: destination,
	})
	assert.Error(t, err)
	content, readErr := os.ReadFile(filePath)
	if assert.NoError(t, readErr) {
		assert.Equal(t, []byte("keep me"), content)
	}
}

func TestRestore_OverwriteReplacesOnlyArchiveEntries(t *testing.T) {
	root := t.TempDir()
	artifacts := writeValidRestoreFixture(t, filepath.Join(root, "archive.tgz"))
	destination := filepath.Join(root, "destination")
	filePath := filepath.Join(destination, "Documents", "file.txt")
	keepPath := filepath.Join(destination, "keep.txt")
	if err := os.MkdirAll(filepath.Dir(filePath), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filePath, []byte("old"), 0o600); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(keepPath, []byte("unrelated"), 0o600); err != nil {
		t.Fatal(err)
	}

	_, err := NewService(nil).Restore(context.Background(), RestoreOptions{
		Archive:     artifacts.ArchivePath,
		Destination: destination,
		Overwrite:   true,
	})
	if !assert.NoError(t, err) {
		return
	}
	content, err := os.ReadFile(filePath)
	if !assert.NoError(t, err) {
		return
	}
	assert.Equal(t, []byte("payload"), content)
	content, err = os.ReadFile(keepPath)
	if assert.NoError(t, err) {
		assert.Equal(t, []byte("unrelated"), content)
	}
}

func TestRestore_OverwriteReplacesLeafSymlink(t *testing.T) {
	root := t.TempDir()
	archivePath := filepath.Join(root, "symlink.tgz")
	manifest := Manifest{SchemaVersion: 1, Entries: []ManifestEntry{
		{Path: "Documents", Type: EntryTypeDirectory},
		{Path: "Documents/settings", Type: EntryTypeFile},
		{Path: "Documents/settings-link", Type: EntryTypeSymlink, Target: "settings"},
	}}
	artifacts := writeTestArchiveFixture(t, archivePath, []testArchiveMember{
		{Header: tar.Header{Name: "Documents/", Typeflag: tar.TypeDir, Mode: 0o755}},
		{
			Header: tar.Header{Name: "Documents/settings", Typeflag: tar.TypeReg, Mode: 0o600, Size: 7},
			Data:   []byte("setting"),
		},
		{
			Header: tar.Header{
				Name:     "Documents/settings-link",
				Typeflag: tar.TypeSymlink,
				Linkname: "settings",
				Mode:     0o777,
			},
		},
	}, manifest)
	destination := filepath.Join(root, "destination")
	if err := os.MkdirAll(filepath.Join(destination, "Documents"), 0o755); err != nil {
		t.Fatal(err)
	}
	linkPath := filepath.Join(destination, "Documents", "settings-link")
	if err := os.Symlink("/etc/hosts", linkPath); err != nil {
		t.Fatal(err)
	}

	_, err := NewService(nil).Restore(context.Background(), RestoreOptions{
		Archive: artifacts.ArchivePath, Destination: destination, Overwrite: true,
	})
	if !assert.NoError(t, err) {
		return
	}
	linkTarget, err := os.Readlink(linkPath)
	if assert.NoError(t, err) {
		assert.Equal(t, "settings", linkTarget)
	}
}

func TestRestore_UnsafeSymlinkRequiresOptInAtApply(t *testing.T) {
	root := t.TempDir()
	archivePath := filepath.Join(root, "unsafe.tgz")
	manifest := Manifest{SchemaVersion: 1, Entries: []ManifestEntry{{
		Path: "link", Type: EntryTypeSymlink, Target: "/etc/passwd", Unsafe: true,
		Warning: "absolute symlink target; restore requires --allow-unsafe-links",
	}}}
	artifacts := writeTestArchiveFixture(t, archivePath, []testArchiveMember{{
		Header: tar.Header{Name: "link", Typeflag: tar.TypeSymlink, Linkname: "/etc/passwd", Mode: 0o777},
	}}, manifest)

	destination := filepath.Join(root, "destination")
	_, err := NewService(nil).Restore(context.Background(), RestoreOptions{
		Archive: artifacts.ArchivePath, Destination: destination,
	})
	assert.Error(t, err)
	assert.NoDirExists(t, destination)

	_, err = NewService(nil).Restore(context.Background(), RestoreOptions{
		Archive: artifacts.ArchivePath, Destination: destination, AllowUnsafeLinks: true,
	})
	if !assert.NoError(t, err) {
		return
	}
	linkTarget, err := os.Readlink(filepath.Join(destination, "link"))
	if assert.NoError(t, err) {
		assert.Equal(t, "/etc/passwd", linkTarget)
	}
}

func TestRestore_DirectoryTypeConflictsFailEvenWithOverwrite(t *testing.T) {
	t.Run("archive file collides with directory", func(t *testing.T) {
		root := t.TempDir()
		artifacts := writeValidRestoreFixture(t, filepath.Join(root, "archive.tgz"))
		destination := filepath.Join(root, "destination")
		conflictPath := filepath.Join(destination, "Documents", "file.txt")
		sentinel := filepath.Join(conflictPath, "sentinel")
		if err := os.MkdirAll(conflictPath, 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(sentinel, []byte("keep"), 0o600); err != nil {
			t.Fatal(err)
		}

		_, err := NewService(nil).Restore(context.Background(), RestoreOptions{
			Archive: artifacts.ArchivePath, Destination: destination, Overwrite: true,
		})
		assert.Error(t, err)
		content, readErr := os.ReadFile(sentinel)
		if assert.NoError(t, readErr) {
			assert.Equal(t, []byte("keep"), content)
		}
	})

	t.Run("archive directory collides with file", func(t *testing.T) {
		root := t.TempDir()
		archivePath := filepath.Join(root, "directory.tgz")
		manifest := Manifest{
			SchemaVersion: 1,
			Entries:       []ManifestEntry{{Path: "Documents", Type: EntryTypeDirectory}},
		}
		artifacts := writeTestArchiveFixture(t, archivePath, []testArchiveMember{{
			Header: tar.Header{Name: "Documents/", Typeflag: tar.TypeDir, Mode: 0o755},
		}}, manifest)
		destination := filepath.Join(root, "destination")
		if err := os.MkdirAll(destination, 0o755); err != nil {
			t.Fatal(err)
		}
		conflictPath := filepath.Join(destination, "Documents")
		if err := os.WriteFile(conflictPath, []byte("keep"), 0o600); err != nil {
			t.Fatal(err)
		}

		_, err := NewService(nil).Restore(context.Background(), RestoreOptions{
			Archive: artifacts.ArchivePath, Destination: destination, Overwrite: true,
		})
		assert.Error(t, err)
		content, readErr := os.ReadFile(conflictPath)
		if assert.NoError(t, readErr) {
			assert.Equal(t, []byte("keep"), content)
		}
	})
}

func TestRestore_RejectsPreExistingSymlinkAncestor(t *testing.T) {
	root := t.TempDir()
	artifacts := writeValidRestoreFixture(t, filepath.Join(root, "archive.tgz"))
	destination := filepath.Join(root, "destination")
	outside := filepath.Join(root, "outside")
	if err := os.MkdirAll(destination, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(outside, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.Symlink(outside, filepath.Join(destination, "Documents")); err != nil {
		t.Fatal(err)
	}

	_, err := NewService(nil).Restore(context.Background(), RestoreOptions{
		Archive: artifacts.ArchivePath, Destination: destination,
	})
	assert.Error(t, err)
	assert.NoFileExists(t, filepath.Join(outside, "file.txt"))
}

func TestRestore_CreatesImplicitParents(t *testing.T) {
	root := t.TempDir()
	archivePath := filepath.Join(root, "nested-file.tgz")
	manifest := Manifest{
		SchemaVersion: 1,
		Entries:       []ManifestEntry{{Path: "notes/todo.txt", Type: EntryTypeFile}},
	}
	artifacts := writeTestArchiveFixture(t, archivePath, []testArchiveMember{{
		Header: tar.Header{Name: "notes/todo.txt", Typeflag: tar.TypeReg, Mode: 0o600, Size: 4},
		Data:   []byte("todo"),
	}}, manifest)
	destination := filepath.Join(root, "destination")

	_, err := NewService(nil).Restore(context.Background(), RestoreOptions{
		Archive: artifacts.ArchivePath, Destination: destination,
	})
	if !assert.NoError(t, err) {
		return
	}
	content, err := os.ReadFile(filepath.Join(destination, "notes", "todo.txt"))
	if assert.NoError(t, err) {
		assert.Equal(t, []byte("todo"), content)
	}
}

func TestRestore_CreatesDestinationWithRestrictedPermissions(t *testing.T) {
	root := t.TempDir()
	artifacts := writeValidRestoreFixture(t, filepath.Join(root, "archive.tgz"))
	destination := filepath.Join(root, "destination")

	_, err := NewService(nil).Restore(context.Background(), RestoreOptions{
		Archive: artifacts.ArchivePath, Destination: destination,
	})
	if !assert.NoError(t, err) {
		return
	}
	info, err := os.Stat(destination)
	if !assert.NoError(t, err) {
		return
	}
	permissions := info.Mode().Perm()
	assert.Equal(t, os.FileMode(0o700), permissions&0o700)
	assert.Zero(t, permissions&^os.FileMode(0o750))
}

func TestRestore_DryRunDoesNotCreateDestination(t *testing.T) {
	root := t.TempDir()
	artifacts := writeValidRestoreFixture(t, filepath.Join(root, "archive.tgz"))
	destination := filepath.Join(root, "destination")

	result, err := NewService(nil).Restore(context.Background(), RestoreOptions{
		Archive: artifacts.ArchivePath, Destination: destination, DryRun: true,
	})
	if !assert.NoError(t, err) {
		return
	}
	assert.ElementsMatch(t, []string{"Documents", "Documents/file.txt"}, result.Paths)
	assert.NoDirExists(t, destination)
}

func TestRestore_PreservesModeAndModificationTimeAndUsesCurrentOwner(t *testing.T) {
	root := t.TempDir()
	archivePath := filepath.Join(root, "metadata.tgz")
	modified := time.Date(2023, time.March, 4, 5, 6, 7, 0, time.Local)
	manifest := Manifest{SchemaVersion: 1, Entries: []ManifestEntry{{Path: "file.txt", Type: EntryTypeFile}}}
	artifacts := writeTestArchiveFixture(t, archivePath, []testArchiveMember{{
		Header: tar.Header{Name: "file.txt", Typeflag: tar.TypeReg, Mode: 0o4755, ModTime: modified, Size: 1},
		Data:   []byte("x"),
	}}, manifest)
	destination := filepath.Join(root, "destination")

	service := NewService(nil)
	ownerChanges := 0
	service.changeOwner = func(_ *os.Root, _ string, _, _ int, _ bool) error {
		ownerChanges++
		return nil
	}
	_, err := service.Restore(context.Background(), RestoreOptions{
		Archive: artifacts.ArchivePath, Destination: destination,
	})
	if !assert.NoError(t, err) {
		return
	}
	filePath := filepath.Join(destination, "file.txt")
	info, err := os.Stat(filePath)
	if !assert.NoError(t, err) {
		return
	}
	assert.Equal(t, os.FileMode(0o755), info.Mode().Perm())
	assert.True(t, info.ModTime().Equal(modified))
	owner, ok := info.Sys().(*syscall.Stat_t)
	if assert.True(t, ok) {
		assert.Equal(t, uint32(os.Geteuid()), owner.Uid)
		assert.Equal(t, uint32(os.Getegid()), owner.Gid)
	}
	assert.Zero(t, ownerChanges)
}

func TestPermissionsFromTarModeKeepsOnlyOrdinaryPermissionBits(t *testing.T) {
	tests := []struct {
		name string
		mode int64
		want os.FileMode
	}{
		{name: "ordinary bits", mode: 0o640, want: 0o640},
		{name: "setuid stripped", mode: 0o4755, want: 0o755},
		{name: "high bits ignored", mode: 1<<50 | 0o604, want: 0o604},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, err := permissionsFromTarMode(test.mode)
			if !assert.NoError(t, err) {
				return
			}
			assert.Equal(t, test.want, got)
		})
	}

	_, err := permissionsFromTarMode(-1)
	assert.Error(t, err)
}

func TestRestore_PreserveOwnerOptIn(t *testing.T) {
	root := t.TempDir()
	archivePath := filepath.Join(root, "owner.tgz")
	uid := int64(os.Geteuid() + 1000)
	gid := int64(os.Getegid() + 1000)
	manifest := Manifest{SchemaVersion: 1, Entries: []ManifestEntry{{
		Path: "file.txt", Type: EntryTypeFile, UID: new(uid), GID: new(gid),
	}}}
	artifacts := writeTestArchiveFixture(t, archivePath, []testArchiveMember{{
		Header: tar.Header{Name: "file.txt", Typeflag: tar.TypeReg, Mode: 0o600, Size: 1, Uid: int(uid), Gid: int(gid)},
		Data:   []byte("x"),
	}}, manifest)
	service := NewService(nil)
	var gotUID, gotGID int
	service.changeOwner = func(_ *os.Root, _ string, uid, gid int, _ bool) error {
		gotUID, gotGID = uid, gid
		return nil
	}

	_, err := service.Restore(context.Background(), RestoreOptions{
		Archive: artifacts.ArchivePath, Destination: filepath.Join(root, "destination"), PreserveOwner: true,
	})
	if !assert.NoError(t, err) {
		return
	}
	assert.Equal(t, int(uid), gotUID)
	assert.Equal(t, int(gid), gotGID)
}

func TestRestore_PreserveOwnerChecksReusedDirectoryEvenWhenArchiveOwnerIsCurrent(t *testing.T) {
	root := t.TempDir()
	archivePath := filepath.Join(root, "current-owner.tgz")
	uid := int64(os.Geteuid())
	gid := int64(os.Getegid())
	manifest := Manifest{SchemaVersion: 1, Entries: []ManifestEntry{{
		Path: "Documents", Type: EntryTypeDirectory, UID: new(uid), GID: new(gid),
	}}}
	artifacts := writeTestArchiveFixture(t, archivePath, []testArchiveMember{{
		Header: tar.Header{Name: "Documents/", Typeflag: tar.TypeDir, Mode: 0o755, Uid: int(uid), Gid: int(gid)},
	}}, manifest)
	destination := filepath.Join(root, "destination")
	if err := os.MkdirAll(filepath.Join(destination, "Documents"), 0o700); err != nil {
		t.Fatal(err)
	}
	service := NewService(nil)
	ownerChanges := 0
	service.ownerOf = func(os.FileInfo) (int64, int64, bool) {
		return uid + 1, gid + 1, true
	}
	service.changeOwner = func(_ *os.Root, _ string, gotUID, gotGID int, _ bool) error {
		ownerChanges++
		assert.Equal(t, int(uid), gotUID)
		assert.Equal(t, int(gid), gotGID)
		return nil
	}

	_, err := service.Restore(context.Background(), RestoreOptions{
		Archive: artifacts.ArchivePath, Destination: destination, PreserveOwner: true,
	})
	if !assert.NoError(t, err) {
		return
	}
	assert.Equal(t, 1, ownerChanges)
}

func TestRestore_ClearedOwnerMetadataUsesCurrentOwner(t *testing.T) {
	root := t.TempDir()
	archivePath := filepath.Join(root, "owner-cleared.tgz")
	manifest := Manifest{SchemaVersion: 1, Entries: []ManifestEntry{{Path: "file.txt", Type: EntryTypeFile}}}
	artifacts := writeTestArchiveFixture(t, archivePath, []testArchiveMember{{
		Header: tar.Header{Name: "file.txt", Typeflag: tar.TypeReg, Mode: 0o600, Size: 1},
		Data:   []byte("x"),
	}}, manifest)
	logger := gen.NewMockLogger(t)
	logger.On("Warn", mck.Anything, mck.Anything).Return().Once()
	service := NewService(logger)
	destination := filepath.Join(root, "destination")

	_, err := service.Restore(context.Background(), RestoreOptions{
		Archive: artifacts.ArchivePath, Destination: destination, PreserveOwner: true,
	})
	if !assert.NoError(t, err) {
		return
	}
	info, err := os.Stat(filepath.Join(destination, "file.txt"))
	if !assert.NoError(t, err) {
		return
	}
	owner, ok := info.Sys().(*syscall.Stat_t)
	if assert.True(t, ok) {
		assert.Equal(t, uint32(os.Geteuid()), owner.Uid)
		assert.Equal(t, uint32(os.Getegid()), owner.Gid)
	}
}

func TestRestore_PreserveOwnerDeniedRetainsReusedDirectoryOwner(t *testing.T) {
	root := t.TempDir()
	archivePath := filepath.Join(root, "owner-denied.tgz")
	uid := int64(os.Geteuid() + 1000)
	gid := int64(os.Getegid() + 1000)
	manifest := Manifest{SchemaVersion: 1, Entries: []ManifestEntry{{
		Path: "Documents", Type: EntryTypeDirectory, UID: new(uid), GID: new(gid),
	}}}
	artifacts := writeTestArchiveFixture(t, archivePath, []testArchiveMember{{
		Header: tar.Header{Name: "Documents/", Typeflag: tar.TypeDir, Mode: 0o755, Uid: int(uid), Gid: int(gid)},
	}}, manifest)
	destination := filepath.Join(root, "destination")
	if err := os.MkdirAll(filepath.Join(destination, "Documents"), 0o700); err != nil {
		t.Fatal(err)
	}
	infoBefore, err := os.Stat(filepath.Join(destination, "Documents"))
	if err != nil {
		t.Fatal(err)
	}
	ownerBefore, ok := infoBefore.Sys().(*syscall.Stat_t)
	if !assert.True(t, ok) {
		return
	}

	logger := gen.NewMockLogger(t)
	logger.On("Warn", "archive ownership not restored", mck.Anything).Return().Once()
	service := NewService(logger)
	service.changeOwner = func(_ *os.Root, _ string, _, _ int, _ bool) error {
		return os.ErrPermission
	}
	_, err = service.Restore(context.Background(), RestoreOptions{
		Archive: artifacts.ArchivePath, Destination: destination, PreserveOwner: true,
	})
	if !assert.NoError(t, err) {
		return
	}
	infoAfter, err := os.Stat(filepath.Join(destination, "Documents"))
	if !assert.NoError(t, err) {
		return
	}
	ownerAfter, ok := infoAfter.Sys().(*syscall.Stat_t)
	if assert.True(t, ok) {
		assert.Equal(t, ownerBefore.Uid, ownerAfter.Uid)
		assert.Equal(t, ownerBefore.Gid, ownerAfter.Gid)
	}
}

func TestRestore_DoesNotApplyPAXExtendedAttributes(t *testing.T) {
	if runtime.GOOS != "darwin" {
		t.Skip("conditional xattr assertion is implemented for macOS")
	}
	xattr, err := exec.LookPath("xattr")
	if err != nil {
		t.Skip("xattr command is unavailable")
	}
	root := t.TempDir()
	probe := filepath.Join(root, "xattr-probe")
	if err := os.WriteFile(probe, []byte("probe"), 0o600); err != nil {
		t.Fatal(err)
	}
	const attribute = "com.example.dum-tarball"
	if output, err := exec.Command(xattr, "-w", attribute, "probe-value", probe).CombinedOutput(); err != nil {
		t.Skipf("filesystem does not support xattrs: %v (%s)", err, output)
	}
	if output, err := exec.Command(
		xattr,
		"-p",
		attribute,
		probe,
	).Output(); err != nil ||
		string(output) != "probe-value\n" {
		t.Skip("xattr probe could not be read back")
	}

	archivePath := filepath.Join(root, "xattr.tgz")
	manifest := Manifest{SchemaVersion: 1, Entries: []ManifestEntry{{Path: "file.txt", Type: EntryTypeFile}}}
	artifacts := writeTestArchiveFixture(t, archivePath, []testArchiveMember{{
		Header: tar.Header{
			Name: "file.txt", Typeflag: tar.TypeReg, Mode: 0o600, Size: 1, Format: tar.FormatPAX,
			PAXRecords: map[string]string{"SCHILY.xattr." + attribute: "archive-value"},
		},
		Data: []byte("x"),
	}}, manifest)
	destination := filepath.Join(root, "destination")
	if _, err := NewService(nil).Restore(context.Background(), RestoreOptions{
		Archive: artifacts.ArchivePath, Destination: destination,
	}); !assert.NoError(t, err) {
		return
	}

	restored := filepath.Join(destination, "file.txt")
	if _, err := os.Stat(restored); !assert.NoError(t, err) {
		return
	}
	if output, err := exec.Command(xattr, "-p", attribute, restored).CombinedOutput(); err == nil {
		t.Fatalf("restore unexpectedly applied extended attribute: %s", output)
	}
}
