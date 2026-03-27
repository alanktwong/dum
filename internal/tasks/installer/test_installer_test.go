package installer

import (
	"context"
	"testing"

	ty "awong/dotfiles/internal/types"

	"github.com/stretchr/testify/assert"
)

func TestTestInstaller_Install(t *testing.T) {
	mockLog := &MockLogger{}
	mockLog.On("Debugf", "%v %v: install test", []any{"...........", "test-play"}).Return()

	ti := &TestInstaller{
		Log: mockLog,
	}

	input := &ty.TaskInput{
		Play:   "test-play",
		Sudo:   false,
		DryRun: false,
	}

	result, err := ti.Install(context.Background(), input)
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.True(t, result.Success)
}
