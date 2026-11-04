package tarball

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoadInput_ResolvesPathsRelativeToInputFile(t *testing.T) {
	configDir := filepath.Join(t.TempDir(), "config")
	if err := os.MkdirAll(configDir, 0o755); err != nil {
		t.Fatal(err)
	}

	inputPath := filepath.Join(configDir, "input.yml")
	data := []byte("paths:\n  - Documents\n  - /mnt/data/Photos\n  - 'folder with spaces'\n")
	if err := os.WriteFile(inputPath, data, 0o600); err != nil {
		t.Fatal(err)
	}

	got, err := LoadInput(context.Background(), inputPath)
	if !assert.NoError(t, err) {
		return
	}

	assert.Equal(t, []string{
		filepath.Join(configDir, "Documents"),
		filepath.Clean("/mnt/data/Photos"),
		filepath.Join(configDir, "folder with spaces"),
	}, got.Paths)
}

func TestLoadInput_ResolvesRelativeInputFileFromCurrentDirectory(t *testing.T) {
	dir := t.TempDir()
	inputPath := filepath.Join(dir, "input.yml")
	if err := os.WriteFile(inputPath, []byte("paths:\n  - source\n"), 0o600); err != nil {
		t.Fatal(err)
	}

	t.Chdir(dir)

	got, err := LoadInput(context.Background(), "input.yml")
	if !assert.NoError(t, err) {
		return
	}
	assert.Equal(t, []string{filepath.Join(dir, "source")}, got.Paths)
}

func TestLoadInput_RejectsInvalidInput(t *testing.T) {
	tests := []struct {
		name string
		data string
	}{
		{name: "malformed yaml", data: "paths: [\n"},
		{name: "missing paths", data: "{}\n"},
		{name: "empty paths", data: "paths: []\n"},
		{name: "empty path", data: "paths:\n  - ''\n"},
		{name: "unknown field", data: "paths:\n  - source\nother: value\n"},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			inputPath := filepath.Join(t.TempDir(), "input.yml")
			if err := os.WriteFile(inputPath, []byte(test.data), 0o600); err != nil {
				t.Fatal(err)
			}

			_, err := LoadInput(context.Background(), inputPath)
			assert.Error(t, err)
		})
	}
}

func TestLoadInput_RejectsDuplicateResolvedPaths(t *testing.T) {
	inputPath := filepath.Join(t.TempDir(), "input.yml")
	data := []byte("paths:\n  - source\n  - ./source\n")
	if err := os.WriteFile(inputPath, data, 0o600); err != nil {
		t.Fatal(err)
	}

	_, err := LoadInput(context.Background(), inputPath)
	assert.Error(t, err)
}
