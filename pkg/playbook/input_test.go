package playbook

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateInput(t *testing.T) {
	input := createTestInput(t)
	assert.Equal(t, "", input.Task, "should initialize with empty task")
}
