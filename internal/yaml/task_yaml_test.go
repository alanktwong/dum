package yaml

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTaskYAML_Validate(t *testing.T) {
	tests := []struct {
		name    string
		task    TaskYAML
		wantErr bool
		errMsg  string
	}{
		{
			name:    "valid dir task",
			task:    TaskYAML{ID: "t1", Type: "dir"},
			wantErr: false,
		},
		{
			name:    "empty id",
			task:    TaskYAML{Type: "dir"},
			wantErr: true,
			errMsg:  "task id cannot be empty",
		},
		{
			name:    "empty type",
			task:    TaskYAML{ID: "t1"},
			wantErr: true,
			errMsg:  "task type cannot be empty",
		},
		{
			name:    "invalid type",
			task:    TaskYAML{ID: "t1", Type: "unknown"},
			wantErr: true,
			errMsg:  "invalid task type",
		},
		{
			name:    "bash with command",
			task:    TaskYAML{ID: "t1", Type: "bash", Command: "echo hello"},
			wantErr: false,
		},
		{
			name:    "bash with script",
			task:    TaskYAML{ID: "t1", Type: "bash", Script: "echo hello"},
			wantErr: false,
		},
		{
			name:    "bash without command or script",
			task:    TaskYAML{ID: "t1", Type: "bash"},
			wantErr: true,
			errMsg:  "bash task must have either command or script",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.task.Validate()
			if tt.wantErr {
				assert.Error(t, err)
				if tt.errMsg != "" {
					assert.Contains(t, err.Error(), tt.errMsg)
				}
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
