package yaml

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoad_Success(t *testing.T) {
	config, err := Load("../../internal/factory/testdata/test_installer.yml")
	assert.NoError(t, err)
	assert.NotNil(t, config)
	assert.Equal(t, "my-test-playbook", config.PlayBook.ID)
	assert.Len(t, config.PlayBook.Plays, 2)
}

func TestLoad_FileNotFound(t *testing.T) {
	_, err := Load("/nonexistent/path/installer.yml")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to read file")
}

func TestLoad_InvalidYaml(t *testing.T) {
	tmpDir := t.TempDir()
	badFile := filepath.Join(tmpDir, "bad.yml")
	err := os.WriteFile(badFile, []byte("{{invalid yaml: [}"), 0644)
	assert.NoError(t, err)

	_, err = Load(badFile)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to unmarshal YAML")
}

func TestLoad_ValidationError(t *testing.T) {
	tmpDir := t.TempDir()
	badFile := filepath.Join(tmpDir, "invalid.yml")
	err := os.WriteFile(badFile, []byte("playbook:\n  plays:\n    - id: \"\"\n"), 0644)
	assert.NoError(t, err)

	_, err = Load(badFile)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "validation failed")
}
