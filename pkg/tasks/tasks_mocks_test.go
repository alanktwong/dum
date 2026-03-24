package tasks

import (
	"context"

	ext "awong/dotfiles/pkg/external"
	i "awong/dotfiles/pkg/tasks/installer"

	mock "github.com/stretchr/testify/mock"
)

func NewTestLogger(t interface {
	mock.TestingT
	Cleanup(func())
}) *i.MockLogger {
	l := i.NewMockLogger(t)
	l.On("Debugf", mock.Anything, mock.Anything).Return().Maybe()
	l.On("Infof", mock.Anything, mock.Anything).Return().Maybe()
	l.On("Debug", mock.Anything, mock.Anything).Return().Maybe()
	l.On("Errorf", mock.Anything, mock.Anything).Return().Maybe()
	l.On("Printlnf", mock.Anything, mock.Anything).Return(nil).Maybe()
	l.On("Warnf", mock.Anything, mock.Anything).Return().Maybe()
	return l
}

type MockGit struct {
	mock.Mock
}

func (m *MockGit) Clone(ctx context.Context, url, name, path string, sudo bool) error {
	ret := m.Called(ctx, url, name, path, sudo)
	//nolint:wrapcheck
	return ret.Error(0)
}

func (m *MockGit) AlreadyExists(targetPath string) bool {
	ret := m.Called(targetPath)
	return ret.Bool(0)
}

type MockMas struct {
	mock.Mock
}

func (m *MockMas) Install(ctx context.Context, app string) error {
	ret := m.Called(ctx, app)
	//nolint:wrapcheck
	return ret.Error(0)
}

func (m *MockMas) List(ctx context.Context) (string, error) {
	ret := m.Called(ctx)
	return ret.String(0), ret.Error(1)
}

type MockCode struct {
	mock.Mock
}

func (m *MockCode) InstallExtension(ctx context.Context, extension string) error {
	ret := m.Called(ctx, extension)
	//nolint:wrapcheck
	return ret.Error(0)
}

func (m *MockCode) ListExtensions(ctx context.Context) (string, error) {
	ret := m.Called(ctx)
	return ret.String(0), ret.Error(1)
}

type MockJetBrainsApp struct {
	mock.Mock
}

func (m *MockJetBrainsApp) Install(ctx context.Context, app, plugin string) error {
	ret := m.Called(ctx, app, plugin)
	//nolint:wrapcheck
	return ret.Error(0)
}

func (m *MockJetBrainsApp) IsInstalled(ideName, plugin string) bool {
	ret := m.Called(ideName, plugin)
	return ret.Bool(0)
}

type MockExt struct {
	mock.Mock
}

func (m *MockExt) IsInstalled(command string) bool {
	ret := m.Called(command)
	return ret.Bool(0)
}

func (m *MockExt) IsOSX() bool {
	ret := m.Called()
	return ret.Bool(0)
}

func (m *MockExt) IsLinux() bool {
	ret := m.Called()
	return ret.Bool(0)
}

func (m *MockExt) IsUserInFileGroup(filePath string) (bool, error) {
	ret := m.Called(filePath)
	return ret.Bool(0), ret.Error(1)
}

func (m *MockExt) CreateDirectory(ctx context.Context, path string, sudo bool) error {
	ret := m.Called(ctx, path, sudo)
	//nolint:wrapcheck
	return ret.Error(0)
}

func (m *MockExt) SoftLink(ctx context.Context, rootPath, src, target string, sudo bool) error {
	ret := m.Called(ctx, rootPath, src, target, sudo)
	//nolint:wrapcheck
	return ret.Error(0)
}

func (m *MockExt) ExpandUser(path string) (string, error) {
	ret := m.Called(path)
	return ret.String(0), ret.Error(1)
}

func (m *MockExt) ToAbsolutePath(path string) (string, error) {
	ret := m.Called(path)
	return ret.String(0), ret.Error(1)
}

func (m *MockExt) IsDir(path string) bool {
	ret := m.Called(path)
	return ret.Bool(0)
}

func (m *MockExt) IsSymlink(path string) bool {
	ret := m.Called(path)
	return ret.Bool(0)
}

func (m *MockExt) RunCommand(ctx context.Context, command string, sudo bool) error {
	ret := m.Called(ctx, command, sudo)
	//nolint:wrapcheck
	return ret.Error(0)
}

func (m *MockExt) GetString(data map[string]interface{}, key string, defaultValue string) string {
	ret := m.Called(data, key, defaultValue)
	return ret.String(0)
}

func (m *MockExt) GetStrings(data map[string]interface{}, key string, defaultValue []string) []string {
	ret := m.Called(data, key, defaultValue)
	return ret.Get(0).([]string)
}

func (m *MockExt) GetBool(data map[string]interface{}, key string, defaultValue bool) bool {
	ret := m.Called(data, key, defaultValue)
	return ret.Bool(0)
}

type MockBrew struct {
	mock.Mock
}

func (m *MockBrew) Install(ctx context.Context, formula string) error {
	ret := m.Called(ctx, formula)
	//nolint:wrapcheck
	return ret.Error(0)
}

func (m *MockBrew) InstallCask(ctx context.Context, formula string) error {
	ret := m.Called(ctx, formula)
	//nolint:wrapcheck
	return ret.Error(0)
}

func (m *MockBrew) Tap(ctx context.Context, formula string) error {
	ret := m.Called(ctx, formula)
	//nolint:wrapcheck
	return ret.Error(0)
}

func (m *MockBrew) Prefix(ctx context.Context) (string, error) {
	ret := m.Called(ctx)
	return ret.String(0), ret.Error(1)
}

func (m *MockBrew) InPath(ctx context.Context, prefix, id string) bool {
	ret := m.Called(ctx, prefix, id)
	return ret.Bool(0)
}

var _ ext.Git = (*MockGit)(nil)
var _ ext.Mas = (*MockMas)(nil)
var _ ext.Code = (*MockCode)(nil)
var _ ext.JetBrainsApp = (*MockJetBrainsApp)(nil)
var _ ext.Ext = (*MockExt)(nil)
var _ ext.Brew = (*MockBrew)(nil)
