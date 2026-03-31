package yaml

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPlayBookYAML_Validate(t *testing.T) {
	tests := []struct {
		name    string
		pb      PlayBookYAML
		wantErr bool
	}{
		{
			name: "valid playbook",
			pb: PlayBookYAML{
				ID:    "test",
				Plays: []PlayYAML{{ID: "play-1", Tasks: []TaskYAML{{ID: "t1", Type: "dir"}}}},
			},
			wantErr: false,
		},
		{
			name:    "empty id",
			pb:      PlayBookYAML{},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.pb.Validate()
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
