package external

import (
	"context"
	"fmt"
	"os"
	"os/user"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDefaultExt_IsInstalled(t *testing.T) {
	u := NewExt()
	tests := []struct {
		name    string
		command string
		want    bool
	}{
		{
			name:    "bash is installed",
			command: "bash",
			want:    true,
		},
		{
			name:    "tar is installed",
			command: "tar",
			want:    true,
		},
		{
			name:    "nonexistent command",
			command: "nonexistentcommand12345",
			want:    false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := u.IsInstalled(tt.command)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestDefaultExt_IsOSX(t *testing.T) {
	u := NewExt()
	got := u.IsOSX()
	assert.Equal(t, runtime.GOOS == "darwin", got)
}

func TestDefaultExt_IsLinux(t *testing.T) {
	u := NewExt()
	got := u.IsLinux()
	assert.Equal(t, runtime.GOOS == "linux", got)
}

func TestDefaultExt_IsDir(t *testing.T) {
	u := NewExt()
	tests := []struct {
		name string
		path string
		want bool
	}{
		{
			name: "current directory is dir",
			path: ".",
			want: true,
		},
		{
			name: "nonexistent path",
			path: "/nonexistent/path/12345",
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := u.IsDir(tt.path)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestDefaultExt_IsSymlink(t *testing.T) {
	u := NewExt()
	tmpDir := t.TempDir()
	symlinkPath := filepath.Join(tmpDir, "symlink")
	err := os.Symlink("/tmp", symlinkPath)
	assert.NoError(t, err)

	tests := []struct {
		name string
		path string
		want bool
	}{
		{
			name: "symlink is symlink",
			path: symlinkPath,
			want: true,
		},
		{
			name: "regular file is not symlink",
			path: tmpDir,
			want: false,
		},
		{
			name: "nonexistent path",
			path: "/nonexistent/path/12345",
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := u.IsSymlink(tt.path)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestDefaultExt_ExpandUser(t *testing.T) {
	u := NewExt()
	tests := []struct {
		name    string
		path    string
		want    string
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name:    "tilde expands",
			path:    "~/test",
			want:    "",
			wantErr: assert.NoError,
		},
		{
			name:    "no tilde returns as is",
			path:    "/absolute/path",
			want:    "/absolute/path",
			wantErr: assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := u.ExpandUser(tt.path)
			if !tt.wantErr(t, err, fmt.Sprintf("ExpandUser(%v)", tt.path)) {
				return
			}
			if tt.path == "~/test" {
				assert.True(t, got != "")
			} else {
				assert.Equalf(t, tt.want, got, "ExpandUser(%v)", tt.path)
			}
		})
	}
}

func TestDefaultExt_RunCommand(t *testing.T) {
	u := NewExt()
	tests := []struct {
		name    string
		command string
		sudo    bool
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name:    "echo command succeeds",
			command: "echo hello",
			sudo:    false,
			wantErr: assert.NoError,
		},
		{
			name:    "false command fails",
			command: "false",
			sudo:    false,
			wantErr: assert.Error,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := u.RunCommand(context.Background(), tt.command, tt.sudo)
			tt.wantErr(t, err, fmt.Sprintf("RunCommand(%v)", tt.command))
		})
	}
}

func TestDefaultExt_CreateDirectory(t *testing.T) {
	u := NewExt()
	tmpDir := t.TempDir()
	testPath := filepath.Join(tmpDir, "testdir")

	tests := []struct {
		name    string
		path    string
		sudo    bool
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name:    "create directory succeeds",
			path:    testPath,
			sudo:    false,
			wantErr: assert.NoError,
		},
		{
			name:    "create existing directory succeeds",
			path:    tmpDir,
			sudo:    false,
			wantErr: assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := u.CreateDirectory(context.Background(), tt.path, tt.sudo)
			tt.wantErr(t, err, fmt.Sprintf("CreateDirectory(%v)", tt.path))
		})
	}
}

func TestDefaultExt_SoftLink(t *testing.T) {
	u := NewExt()
	tmpDir := t.TempDir()
	targetPath := filepath.Join(tmpDir, "target")

	err := os.WriteFile(targetPath, []byte("test"), 0o644)
	assert.NoError(t, err)

	tests := []struct {
		name    string
		root    string
		src     string
		target  string
		sudo    bool
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name:    "create symlink succeeds",
			root:    tmpDir,
			src:     targetPath,
			target:  "link",
			sudo:    false,
			wantErr: assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := u.SoftLink(context.Background(), tt.root, tt.src, tt.target, tt.sudo)
			tt.wantErr(t, err, fmt.Sprintf("SoftLink(%v, %v, %v)", tt.root, tt.src, tt.target))
		})
	}
}

func TestDefaultUtils_ToAbsolutePath(t *testing.T) {
	u, err := user.Current()
	assert.NoError(t, err)
	username := u.Username
	utils := NewExt()
	type args struct {
		path string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "should FAIl with empty",
			args: args{
				path: "",
			},
			want:    "",
			wantErr: assert.Error,
		},
		{
			name: "should succeed with ~",
			args: args{
				path: "~/test",
			},
			want:    fmt.Sprintf("/Users/%v/test", username),
			wantErr: assert.NoError,
		},
		{
			name: "should succeed with $HOME",
			args: args{
				path: "$HOME/test",
			},
			want:    fmt.Sprintf("/Users/%v/test", username),
			wantErr: assert.NoError,
		},
		{
			name: "should succeed with $GOPATH",
			args: args{
				path: "$GOPATH/src/test",
			},
			want:    fmt.Sprintf("/Users/%v/go/src/test", username),
			wantErr: assert.NoError,
		},
		{
			name: "should expand env vars",
			args: args{
				path: "$HOME/test",
			},
			want:    fmt.Sprintf("/Users/%v/test", username),
			wantErr: assert.NoError,
		},
		{
			name: "should fail with empty path",
			args: args{
				path: "",
			},
			want:    "",
			wantErr: assert.Error,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := utils.ToAbsolutePath(tt.args.path)
			if !tt.wantErr(t, err, fmt.Sprintf("ToAbsolutePath(%v)", tt.args.path)) {
				return
			}
			assert.Equalf(t, tt.want, got, "toAbsolutePath(%v)", tt.args.path)
		})
	}
}
