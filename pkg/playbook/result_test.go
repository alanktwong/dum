package playbook

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewResult_Success(t *testing.T) {
	pb := &PlayBook{}
	pb.ID = "test-pb"
	input := &Input{PlayBook: pb}
	result := NewResult(input, true)
	assert.NotNil(t, result)
	assert.Equal(t, "test-pb", result.PlayBook)
	assert.True(t, result.Success)
}

func TestNewResult_Failure(t *testing.T) {
	pb := &PlayBook{}
	pb.ID = "test-pb"
	input := &Input{PlayBook: pb}
	result := NewResult(input, false)
	assert.False(t, result.Success)
}
