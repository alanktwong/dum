package external

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

type testExt struct {
	isDirResult    bool
	createDirError error
}

func (t *testExt) IsInstalled(command string) bool                 { return true }
func (t *testExt) IsOSX() bool                                     { return true }
func (t *testExt) IsLinux() bool                                   { return false }
func (t *testExt) IsUserInFileGroup(filePath string) (bool, error) { return false, nil }
func (t *testExt) CreateDirectory(ctx context.Context, path string, sudo bool) error {
	return t.createDirError
}
func (t *testExt) SoftLink(ctx context.Context, rootPath, src, target string, sudo bool) error {
	return nil
}
func (t *testExt) ExpandUser(path string) (string, error)                            { return path, nil }
func (t *testExt) ToAbsolutePath(path string) (string, error)                        { return path, nil }
func (t *testExt) IsDir(path string) bool                                            { return t.isDirResult }
func (t *testExt) IsSymlink(path string) bool                                        { return false }
func (t *testExt) RunCommand(ctx context.Context, command string, sudo bool) error   { return nil }
func (t *testExt) GetString(data map[string]any, key string, def string) string      { return def }
func (t *testExt) GetStrings(data map[string]any, key string, def []string) []string { return def }
func (t *testExt) GetBool(data map[string]any, key string, def bool) bool            { return def }

func TestBrewImpl_Install_Error(t *testing.T) {
	brew := &BrewImpl{Utils: &testExt{}}

	err := brew.Install(context.Background(), "nonexistent-formula-12345")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to brew install")
}

func TestBrewImpl_InstallCask_Error(t *testing.T) {
	brew := &BrewImpl{Utils: &testExt{}}

	err := brew.InstallCask(context.Background(), "nonexistent-cask-12345")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to brew install --cask")
}

func TestBrewImpl_Tap_Error(t *testing.T) {
	brew := &BrewImpl{Utils: &testExt{}}

	err := brew.Tap(context.Background(), "nonexistent-tap-12345")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to brew tap")
}

func TestBrewImpl_Prefix(t *testing.T) {
	brew := &BrewImpl{Utils: &testExt{}}

	prefix, err := brew.Prefix(context.Background())
	assert.NoError(t, err)
	assert.NotEmpty(t, prefix)
}

func TestBrewImpl_InPath_Found(t *testing.T) {
	mockUtils := &testExt{isDirResult: true}
	brew := &BrewImpl{Utils: mockUtils}

	found := brew.InPath(context.Background(), "bin", "git")
	assert.True(t, found)
}

func TestBrewImpl_InPath_NotFound(t *testing.T) {
	mockUtils := &testExt{isDirResult: false}
	brew := &BrewImpl{Utils: mockUtils}

	found := brew.InPath(context.Background(), "bin", "nonexistent")
	assert.False(t, found)
}

func TestBrewImpl_InPath_PrefixError(t *testing.T) {
	brew := &BrewImpl{Utils: &testExt{}}

	found := brew.InPath(context.Background(), "bin", "git")
	assert.False(t, found)
}
