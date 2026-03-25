package plays

import (
	pg "awong/dotfiles/pkg/plays/gen"
	ty "awong/dotfiles/pkg/types"
	"testing"

	"github.com/stretchr/testify/assert"
)

func newTestPlayInput(t *testing.T, playbookID string) *PlayInput {
	pbi := pg.NewMockPlayBookInfo(t)
	pbi.On("GetID").Return(playbookID)
	return &PlayInput{
		Play:     "test-play",
		PlayBook: pbi,
	}
}

func TestNewPlayResult_AllSuccess(t *testing.T) {
	input := newTestPlayInput(t, "test-book")
	tasks := []*ty.TaskResult{
		{Task: "t1", Success: true},
		{Task: "t2", Success: true},
	}
	result := NewPlayResult(input, tasks)
	assert.NotNil(t, result)
	assert.True(t, result.Success)
	assert.Equal(t, "test-book", result.PlayBook)
	assert.Equal(t, "test-play", result.Play)
}

func TestNewPlayResult_SomeFailed(t *testing.T) {
	input := newTestPlayInput(t, "test-book")
	tasks := []*ty.TaskResult{
		{Task: "t1", Success: true},
		{Task: "t2", Success: false},
	}
	result := NewPlayResult(input, tasks)
	assert.NotNil(t, result)
	assert.False(t, result.Success)
}

func TestNewPlayResult_Empty(t *testing.T) {
	input := newTestPlayInput(t, "test-book")
	tasks := []*ty.TaskResult{}
	result := NewPlayResult(input, tasks)
	assert.NotNil(t, result)
	assert.True(t, result.Success)
}

func TestCalculateSuccess_AllPlays(t *testing.T) {
	plays := []*PlayResult{
		{Play: "p1", Success: true},
		{Play: "p2", Success: true},
	}
	assert.True(t, CalculateSuccess(plays))
}

func TestCalculateSuccess_SomeFailed(t *testing.T) {
	plays := []*PlayResult{
		{Play: "p1", Success: true},
		{Play: "p2", Success: false},
	}
	assert.False(t, CalculateSuccess(plays))
}

func TestCalculateSuccess_Empty(t *testing.T) {
	plays := []*PlayResult{}
	assert.True(t, CalculateSuccess(plays))
}
