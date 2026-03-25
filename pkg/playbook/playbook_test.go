package playbook

import (
	pl "awong/dotfiles/pkg/plays"
	ty "awong/dotfiles/pkg/types"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewPlayBook_NilAttributes(t *testing.T) {
	pb, err := NewPlayBook(nil, nil, nil)
	assert.Error(t, err)
	assert.Nil(t, pb)
	assert.Contains(t, err.Error(), "attributes cannot be nil")
}

func TestNewPlayBook_EmptyID(t *testing.T) {
	attr := &ty.Attributes{ID: "", Description: "test"}
	pb, err := NewPlayBook(attr, nil, nil)
	assert.Error(t, err)
	assert.Nil(t, pb)
	assert.Contains(t, err.Error(), "playbook ID cannot be empty")
}

func TestNewPlayBook_Success(t *testing.T) {
	attr := &ty.Attributes{ID: "test-pb", Description: "desc"}
	pb, err := NewPlayBook(attr, nil, nil)
	assert.NoError(t, err)
	assert.NotNil(t, pb)
	assert.Equal(t, "test-pb", pb.ID)
	assert.Equal(t, "desc", pb.Description)
}

func TestNewPlayBook_Success_Enabled(t *testing.T) {
	attr := &ty.Attributes{ID: "test-pb", Enabled: true}
	pb, err := NewPlayBook(attr, nil, nil)
	assert.NoError(t, err)
	assert.True(t, pb.IsEnabled())
}

func TestPlayBook_GetPlays_All(t *testing.T) {
	attr := &ty.Attributes{ID: "test-pb", Enabled: true}
	plays := []*pl.Play{
		{Attributes: ty.Attributes{ID: "play-1", Enabled: true}},
		{Attributes: ty.Attributes{ID: "play-2", Enabled: false}},
	}
	pb, err := NewPlayBook(attr, plays, nil)
	assert.NoError(t, err)

	result := pb.GetPlays(false)
	assert.Equal(t, 2, result.Len())
	assert.True(t, result.Has("play-1"))
	assert.True(t, result.Has("play-2"))
}

func TestPlayBook_GetPlays_Active(t *testing.T) {
	attr := &ty.Attributes{ID: "test-pb", Enabled: true}
	plays := []*pl.Play{
		{Attributes: ty.Attributes{ID: "play-1", Enabled: true}},
		{Attributes: ty.Attributes{ID: "play-2", Enabled: false}},
	}
	pb, err := NewPlayBook(attr, plays, nil)
	assert.NoError(t, err)

	result := pb.GetPlays(true)
	assert.Equal(t, 1, result.Len())
	assert.True(t, result.Has("play-1"))
	assert.False(t, result.Has("play-2"))
}

func TestPlayBook_GetPlays_DisabledPlayBook(t *testing.T) {
	attr := &ty.Attributes{ID: "test-pb", Enabled: false}
	plays := []*pl.Play{
		{Attributes: ty.Attributes{ID: "play-1", Enabled: true}},
	}
	pb, err := NewPlayBook(attr, plays, nil)
	assert.NoError(t, err)

	result := pb.GetPlays(true)
	assert.Equal(t, 0, result.Len())
}
