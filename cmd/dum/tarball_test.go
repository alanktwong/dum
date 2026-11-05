package main

import (
	"bytes"
	"context"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"

	tb "github.com/alanktwong/dum/internal/tarball"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
)

func TestTarballCommand_BackupExplicitPathsBypassInputAndReportArtifacts(t *testing.T) {
	root := t.TempDir()
	home := filepath.Join(root, "home")
	source := filepath.Join(home, "Documents")
	if err := os.MkdirAll(source, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(source, "file.txt"), []byte("data"), 0o600); err != nil {
		t.Fatal(err)
	}
	t.Setenv("HOME", home)
	target := filepath.Join(root, "archives")
	missingInput := filepath.Join(root, "does-not-exist.yml")

	command := newTestTarballCommand()
	var output bytes.Buffer
	command.SetOut(&output)
	command.SetContext(context.Background())
	command.SetArgs([]string{"backup", "--target", target, "--input", missingInput, source})

	err := command.Execute()
	if !assert.NoError(t, err) {
		return
	}
	assert.Contains(t, output.String(), filepath.Join(target, "Documents-"))
	archives, err := filepath.Glob(filepath.Join(target, "Documents-*.tgz"))
	if !assert.NoError(t, err) || !assert.Len(t, archives, 1) {
		return
	}
	assert.FileExists(t, archives[0]+".sha256")
	assert.FileExists(t, archives[0]+".manifest.yaml")
}

func TestTarballCommand_BackupListModeRequiresInput(t *testing.T) {
	command := newTestTarballCommand()
	var output bytes.Buffer
	var stderr bytes.Buffer
	command.SetOut(&output)
	command.SetErr(&stderr)
	command.SetContext(context.Background())
	command.SetArgs([]string{"backup"})

	err := command.Execute()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "--input")
}

func TestTarballCommand_TargetDefaultPrecedence(t *testing.T) {
	home := t.TempDir()
	dataHome := filepath.Join(t.TempDir(), "xdg-data")
	t.Setenv("HOME", home)
	t.Setenv("XDG_DATA_HOME", dataHome)
	t.Setenv("TGZ_TARGET", filepath.Join(t.TempDir(), "env-target"))

	command := newTestTarballCommand()
	backup := findCommand(command.Commands(), "backup")
	if !assert.NotNil(t, backup) {
		return
	}
	target, err := backup.Flags().GetString(tarballTargetFlag)
	if assert.NoError(t, err) {
		assert.Equal(t, os.Getenv("TGZ_TARGET"), target)
	}

	t.Setenv("TGZ_TARGET", "")
	command = newTestTarballCommand()
	backup = findCommand(command.Commands(), "backup")
	if !assert.NotNil(t, backup) {
		return
	}
	target, err = backup.Flags().GetString(tarballTargetFlag)
	if assert.NoError(t, err) {
		assert.Equal(t, filepath.Join(dataHome, "tarballs"), target)
	}

	t.Setenv("XDG_DATA_HOME", "")
	command = newTestTarballCommand()
	backup = findCommand(command.Commands(), "backup")
	if !assert.NotNil(t, backup) {
		return
	}
	target, err = backup.Flags().GetString(tarballTargetFlag)
	if assert.NoError(t, err) {
		assert.Equal(t, filepath.Join(home, ".local", "share", "tarballs"), target)
	}
}

func TestTarballCommand_RestoreExplicitArchiveBypassesInput(t *testing.T) {
	root := t.TempDir()
	home := filepath.Join(root, "home")
	source := filepath.Join(home, "Documents")
	if err := os.MkdirAll(source, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(source, "file.txt"), []byte("data"), 0o600); err != nil {
		t.Fatal(err)
	}
	target := filepath.Join(root, "archives")
	service := tb.NewService(nil)
	backup, err := service.Backup(context.Background(), tb.BackupOptions{
		Paths: []string{source}, Target: target, Home: home,
	})
	if !assert.NoError(t, err) {
		return
	}
	if !assert.Len(t, backup.Artifacts, 1) {
		return
	}

	destination := filepath.Join(root, "destination")
	command := newTestTarballCommand()
	var output bytes.Buffer
	var stderr bytes.Buffer
	command.SetOut(&output)
	command.SetErr(&stderr)
	command.SetContext(context.Background())
	command.SetArgs(
		[]string{
			"restore",
			backup.Artifacts[0].ArchivePath,
			"--input",
			filepath.Join(root, "missing.yml"),
			"--dest",
			destination,
		},
	)
	err = command.Execute()
	if !assert.NoError(t, err) {
		return
	}
	content, err := os.ReadFile(filepath.Join(destination, "Documents", "file.txt"))
	if assert.NoError(t, err) {
		assert.Equal(t, []byte("data"), content)
	}
	assert.Contains(t, output.String(), backup.Artifacts[0].ArchivePath)
}

func TestTarballCommand_DoesNotReportSuccessAfterPartialRestoreError(t *testing.T) {
	deps := &TarballCommand{
		Service: tb.NewService(nil),
		restoreArchive: func(context.Context, tb.RestoreOptions) (tb.RestoreResult, error) {
			return tb.RestoreResult{
				ArchivePath: "archive.tgz",
				Destination: "/restore",
				Paths:       []string{"file.txt"},
			}, errors.New("apply failed after writing a file")
		},
	}
	command := NewTarballCommand("dum", deps)
	var output bytes.Buffer
	var stderr bytes.Buffer
	command.SetOut(&output)
	command.SetErr(&stderr)
	command.SetContext(context.Background())
	command.SetArgs([]string{"restore", "archive.tgz"})

	err := command.Execute()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "apply failed")
	assert.Contains(t, stderr.String(), "apply failed")
	assert.NotContains(t, output.String(), "restored:")
}

func TestTarballCommand_RestoreResolvesRelativeArchiveAndDefaultsToHome(t *testing.T) {
	root := t.TempDir()
	home := filepath.Join(root, "home")
	t.Setenv("HOME", home)
	source := filepath.Join(root, "external", "Books")
	if err := os.MkdirAll(source, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(source, "readme.txt"), []byte("book"), 0o600); err != nil {
		t.Fatal(err)
	}
	target := filepath.Join(root, "archives")
	backup, err := tb.NewService(nil).Backup(context.Background(), tb.BackupOptions{
		Paths: []string{source}, Target: target, Home: home,
	})
	if !assert.NoError(t, err) || !assert.Len(t, backup.Artifacts, 1) {
		return
	}
	t.Chdir(target)

	command := newTestTarballCommand()
	var output bytes.Buffer
	command.SetOut(&output)
	command.SetContext(context.Background())
	command.SetArgs([]string{
		"restore", filepath.Base(backup.Artifacts[0].ArchivePath),
		"--input", filepath.Join(root, "missing-input.yml"),
	})
	if !assert.NoError(t, command.Execute()) {
		return
	}
	content, err := os.ReadFile(filepath.Join(home, "Books", "readme.txt"))
	if assert.NoError(t, err) {
		assert.Equal(t, []byte("book"), content)
	}
}

func TestTarballCommand_DryRunReportsPlansWithoutWriting(t *testing.T) {
	t.Run("backup", func(t *testing.T) {
		root := t.TempDir()
		home := filepath.Join(root, "home")
		source := filepath.Join(home, "Documents")
		if err := os.MkdirAll(source, 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(filepath.Join(source, "notes.txt"), []byte("notes"), 0o600); err != nil {
			t.Fatal(err)
		}
		t.Setenv("HOME", home)
		envTarget := filepath.Join(root, "environment-target")
		t.Setenv("TGZ_TARGET", envTarget)
		flagTarget := filepath.Join(root, "flag-target")
		command := newTestTarballCommand()
		var output bytes.Buffer
		command.SetOut(&output)
		command.SetContext(context.Background())
		command.SetArgs([]string{"backup", "--target", flagTarget, "--dry-run", source})
		if !assert.NoError(t, command.Execute()) {
			return
		}
		assert.Contains(t, output.String(), filepath.Join(flagTarget, "Documents-"))
		assert.Contains(t, output.String(), ".tgz.sha256")
		assert.Contains(t, output.String(), ".tgz.manifest.yaml")
		assert.NoDirExists(t, flagTarget)
		assert.NoDirExists(t, envTarget)
	})

	t.Run("restore", func(t *testing.T) {
		root := t.TempDir()
		home := filepath.Join(root, "home")
		source := filepath.Join(root, "external", "Books")
		if err := os.MkdirAll(source, 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(filepath.Join(source, "readme.txt"), []byte("book"), 0o600); err != nil {
			t.Fatal(err)
		}
		backup, err := tb.NewService(nil).Backup(context.Background(), tb.BackupOptions{
			Paths: []string{source}, Target: filepath.Join(root, "archives"), Home: home,
		})
		if !assert.NoError(t, err) || !assert.Len(t, backup.Artifacts, 1) {
			return
		}
		destination := filepath.Join(root, "destination")
		command := newTestTarballCommand()
		var output bytes.Buffer
		command.SetOut(&output)
		command.SetContext(context.Background())
		command.SetArgs([]string{"restore", backup.Artifacts[0].ArchivePath, "--dest", destination, "--dry-run"})
		if !assert.NoError(t, command.Execute()) {
			return
		}
		assert.Contains(t, output.String(), "would restore:")
		assert.Contains(t, output.String(), filepath.Join(destination, "Books"))
		assert.Contains(t, output.String(), "Books/readme.txt")
		assert.NoDirExists(t, destination)
	})
}

func TestTarballCommand_RestoreListUsesInputYAML(t *testing.T) {
	root := t.TempDir()
	inputDir := filepath.Join(root, "config")
	source := filepath.Join(inputDir, "Documents")
	if err := os.MkdirAll(source, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(source, "file.txt"), []byte("data"), 0o600); err != nil {
		t.Fatal(err)
	}
	inputPath := filepath.Join(inputDir, "input.yml")
	if err := os.WriteFile(inputPath, []byte("paths:\n  - Documents\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	target := filepath.Join(root, "archives")
	service := tb.NewService(nil)
	if _, err := service.Backup(context.Background(), tb.BackupOptions{
		Paths: []string{source}, Target: target, Home: filepath.Join(root, "home"),
	}); !assert.NoError(t, err) {
		return
	}

	destination := filepath.Join(root, "destination")
	command := newTestTarballCommand()
	var output bytes.Buffer
	command.SetOut(&output)
	command.SetContext(context.Background())
	command.SetArgs([]string{"restore", "--input", inputPath, "--target", target, "--dest", destination})
	if !assert.NoError(t, command.Execute()) {
		return
	}
	content, err := os.ReadFile(filepath.Join(destination, "Documents", "file.txt"))
	if assert.NoError(t, err) {
		assert.Equal(t, []byte("data"), content)
	}
	assert.Contains(t, output.String(), "restored:")
}

func TestTarballCommand_ListRestoreContinuesAfterMissingMatch(t *testing.T) {
	root := t.TempDir()
	inputDir := filepath.Join(root, "config")
	source := filepath.Join(inputDir, "Documents")
	if err := os.MkdirAll(source, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(source, "file.txt"), []byte("data"), 0o600); err != nil {
		t.Fatal(err)
	}
	inputPath := filepath.Join(inputDir, "input.yml")
	if err := os.WriteFile(inputPath, []byte("paths:\n  - Missing\n  - Documents\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	target := filepath.Join(root, "archives")
	if _, err := tb.NewService(nil).Backup(context.Background(), tb.BackupOptions{
		Paths: []string{source}, Target: target, Home: filepath.Join(root, "home"),
	}); !assert.NoError(t, err) {
		return
	}

	destination := filepath.Join(root, "destination")
	command := newTestTarballCommand()
	var output bytes.Buffer
	var stderr bytes.Buffer
	command.SetOut(&output)
	command.SetErr(&stderr)
	command.SetContext(context.Background())
	command.SetArgs([]string{"restore", "--input", inputPath, "--target", target, "--dest", destination})
	err := command.Execute()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "Missing")
	assert.Contains(t, stderr.String(), "Missing")
	content, readErr := os.ReadFile(filepath.Join(destination, "Documents", "file.txt"))
	if assert.NoError(t, readErr) {
		assert.Equal(t, []byte("data"), content)
	}
}

func TestTarballCommand_ListRestoreContinuesAfterInvalidArchive(t *testing.T) {
	root := t.TempDir()
	inputDir := filepath.Join(root, "config")
	documents := filepath.Join(inputDir, "Documents")
	photos := filepath.Join(inputDir, "Photos")
	for _, directory := range []string{documents, photos} {
		if err := os.MkdirAll(directory, 0o755); err != nil {
			t.Fatal(err)
		}
	}
	if err := os.WriteFile(filepath.Join(documents, "notes.txt"), []byte("notes"), 0o600); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(photos, "image.bin"), []byte("photo"), 0o600); err != nil {
		t.Fatal(err)
	}
	inputPath := filepath.Join(inputDir, "input.yml")
	if err := os.WriteFile(inputPath, []byte("paths:\n  - Documents\n  - Photos\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	target := filepath.Join(root, "archives")
	backup, err := tb.NewService(nil).Backup(context.Background(), tb.BackupOptions{
		Paths: []string{documents, photos}, Target: target, Home: filepath.Join(root, "home"),
	})
	if !assert.NoError(t, err) || !assert.Len(t, backup.Artifacts, 2) {
		return
	}
	if err := os.WriteFile(backup.Artifacts[0].ChecksumPath, []byte("invalid checksum\n"), 0o600); err != nil {
		t.Fatal(err)
	}

	destination := filepath.Join(root, "destination")
	command := newTestTarballCommand()
	var output bytes.Buffer
	var stderr bytes.Buffer
	command.SetOut(&output)
	command.SetErr(&stderr)
	command.SetContext(context.Background())
	command.SetArgs([]string{"restore", "--input", inputPath, "--target", target, "--dest", destination})
	err = command.Execute()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "Documents")
	assert.Contains(t, stderr.String(), "checksum")
	assert.Contains(t, output.String(), "Photos-")
	assert.NoFileExists(t, filepath.Join(destination, "Documents", "notes.txt"))
	photo, readErr := os.ReadFile(filepath.Join(destination, "Photos", "image.bin"))
	if assert.NoError(t, readErr) {
		assert.Equal(t, []byte("photo"), photo)
	}
}

func TestTarballCommand_RestoreMultipleMatchesUsesPicker(t *testing.T) {
	root := t.TempDir()
	inputDir := filepath.Join(root, "config")
	source := filepath.Join(inputDir, "Documents")
	if err := os.MkdirAll(source, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(source, "file.txt"), []byte("data"), 0o600); err != nil {
		t.Fatal(err)
	}
	inputPath := filepath.Join(inputDir, "input.yml")
	if err := os.WriteFile(inputPath, []byte("paths:\n  - Documents\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	target := filepath.Join(root, "archives")
	backup, err := tb.NewService(nil).Backup(context.Background(), tb.BackupOptions{
		Paths: []string{source}, Target: target, Home: filepath.Join(root, "home"),
	})
	if !assert.NoError(t, err) {
		return
	}
	first := backup.Artifacts[0]
	secondArchive := filepath.Join(target, "Documents-2099-01-01-000000.tgz")
	for _, pair := range [][2]string{
		{first.ArchivePath, secondArchive},
		{first.ChecksumPath, secondArchive + ".sha256"},
		{first.ManifestPath, secondArchive + ".manifest.yaml"},
	} {
		data, err := os.ReadFile(pair[0])
		if err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(pair[1], data, 0o600); err != nil {
			t.Fatal(err)
		}
	}

	deps := &TarballCommand{Service: tb.NewService(nil)}
	var picked []string
	deps.pickArchive = func(_ context.Context, candidates []string) (string, error) {
		picked = append([]string(nil), candidates...)
		return secondArchive, nil
	}
	command := NewTarballCommand("dum", deps)
	var output bytes.Buffer
	command.SetOut(&output)
	command.SetContext(context.Background())
	command.SetArgs(
		[]string{"restore", "--input", inputPath, "--target", target, "--dest", filepath.Join(root, "destination")},
	)
	if !assert.NoError(t, command.Execute()) {
		return
	}
	assert.Len(t, picked, 2)
	assert.Contains(t, output.String(), secondArchive)
}

func TestTarballCommand_MultipleMatchesWithoutPickerListsCandidates(t *testing.T) {
	root := t.TempDir()
	inputDir := filepath.Join(root, "config")
	source := filepath.Join(inputDir, "Documents")
	if err := os.MkdirAll(source, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(source, "file.txt"), []byte("data"), 0o600); err != nil {
		t.Fatal(err)
	}
	inputPath := filepath.Join(inputDir, "input.yml")
	if err := os.WriteFile(inputPath, []byte("paths:\n  - Documents\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	target := filepath.Join(root, "archives")
	backup, err := tb.NewService(nil).Backup(context.Background(), tb.BackupOptions{
		Paths: []string{source}, Target: target, Home: filepath.Join(root, "home"),
	})
	if !assert.NoError(t, err) {
		return
	}
	secondArchive := filepath.Join(target, "Documents-2099-01-01-000000.tgz")
	for _, pair := range [][2]string{
		{backup.Artifacts[0].ArchivePath, secondArchive},
		{backup.Artifacts[0].ChecksumPath, secondArchive + ".sha256"},
		{backup.Artifacts[0].ManifestPath, secondArchive + ".manifest.yaml"},
	} {
		data, err := os.ReadFile(pair[0])
		if err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(pair[1], data, 0o600); err != nil {
			t.Fatal(err)
		}
	}

	noTools := t.TempDir()
	t.Setenv("PATH", noTools)
	command := newTestTarballCommand()
	var output bytes.Buffer
	var stderr bytes.Buffer
	command.SetOut(&output)
	command.SetErr(&stderr)
	command.SetContext(context.Background())
	command.SetArgs(
		[]string{"restore", "--input", inputPath, "--target", target, "--dest", filepath.Join(root, "destination")},
	)
	err = command.Execute()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "multiple tarballs")
	assert.Contains(t, err.Error(), secondArchive)
	assert.Contains(t, err.Error(), backup.Artifacts[0].ArchivePath)
}

func newTestTarballCommand() *cobra.Command {
	return NewTarballCommand("dum", &TarballCommand{Service: tb.NewService(nil)})
}

func TestTarballCommand_PickerOutputMustBeCandidate(t *testing.T) {
	selected := filepath.Join(t.TempDir(), "not-a-candidate.tgz")
	_, err := validateArchiveSelection(selected, []string{"candidate.tgz"})
	assert.Error(t, err)
	assert.True(t, strings.Contains(err.Error(), "not a candidate"))

	valid, err := validateArchiveSelection("candidate.tgz\n", []string{"candidate.tgz"})
	if assert.NoError(t, err) {
		assert.Equal(t, "candidate.tgz", valid)
	}
}

func TestTarballCommand_HelpExplainsInputYAMLShape(t *testing.T) {
	for _, name := range []string{"backup", "restore"} {
		command := newTestTarballCommand()
		sub := findCommand(command.Commands(), name)
		if !assert.NotNil(t, sub, "missing %s command", name) {
			return
		}
		var output bytes.Buffer
		sub.SetOut(&output)
		sub.SetErr(&output)
		assert.NoError(t, sub.Help())

		help := output.String()
		assert.Contains(t, help, "INPUT.yml", name)
		assert.Contains(t, help, "  paths:", name)
		assert.Contains(t, help, "- Documents", name)
		assert.Contains(t, help, "relative to the input.yml directory", name)
		assert.Contains(t, help, "duplicate paths are rejected", name)
	}
}

func TestTarballCommand_RegistersRestoreSafetyFlags(t *testing.T) {
	command := newTestTarballCommand()
	restore := findCommand(command.Commands(), "restore")
	if !assert.NotNil(t, restore) {
		return
	}

	for _, name := range []string{
		tarballInputFlag,
		tarballTargetFlag,
		tarballDestinationFlag,
		tarballOverwriteFlag,
		tarballAllowUnsafeLinksFlag,
		tarballPreserveOwnerFlag,
		tarballDryRunFlag,
	} {
		assert.NotNil(t, restore.Flags().Lookup(name), "missing --%s", name)
	}
}
