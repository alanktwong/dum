package plays

import (
	i "awong/dotfiles/pkg/tasks/installer"
	"testing"

	"github.com/stretchr/testify/assert"
)

func newPlayExecutorWithMocks(t *testing.T) *PlayExecutor {
	return &PlayExecutor{
		Log: i.NewMockLogger(t),
		Ext: &i.MockExt{},
	}
}

func TestNewPlayExecutor_NoArgs(t *testing.T) {
	e := NewPlayExecutor()
	assert.NotNil(t, e)
	assert.NotNil(t, e.Log)
	assert.NotNil(t, e.Ext)
}

func TestPlayExecutor_Run(t *testing.T) {
	e := newPlayExecutorWithMocks(t)
	assert.NotNil(t, e)
}
