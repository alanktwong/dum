package external

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCodeImpl_InstallExtension(t *testing.T) {
	installFakeCodeCLI(t, "#!/bin/sh\nexit 0\n")
	code := NewCode()

	err := code.InstallExtension(context.Background(), "esbenp.prettier-vscode")
	assert.NoError(t, err)
}

func TestCodeImpl_ListExtensions(t *testing.T) {
	installFakeCodeCLI(t, "#!/bin/sh\necho esbenp.prettier-vscode\n")
	code := NewCode()

	output, err := code.ListExtensions(context.Background())
	assert.NoError(t, err)
	assert.NotEmpty(t, output)
}

func installFakeCodeCLI(t *testing.T, script string) {
	t.Helper()
	binDir := t.TempDir()
	path := filepath.Join(binDir, "code")
	assert.NoError(t, os.WriteFile(path, []byte(script), 0o755))
	t.Setenv("PATH", binDir+string(os.PathListSeparator)+os.Getenv("PATH"))
}
