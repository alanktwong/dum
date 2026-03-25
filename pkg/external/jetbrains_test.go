package external

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestJetBrainsImpl_IsInstalled(t *testing.T) {
	jetbrains := NewJetBrains()

	pluginPath := filepath.Join(os.Getenv("HOME"), "Library", "Application Support", "JetBrains", "IntelliJIDEA", "plugins", "test-plugin")
	os.MkdirAll(filepath.Dir(pluginPath), 0755)
	defer os.RemoveAll(filepath.Dir(pluginPath))

	result := jetbrains.IsInstalled("IntelliJIDEA", "test-plugin")
	assert.True(t, result)
}
