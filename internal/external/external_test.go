package external

import (
	"fmt"
	"os/user"
	"testing"

	"github.com/stretchr/testify/assert"
)

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
