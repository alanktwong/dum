package playbook

import (
	pl "alanktwong/dum/internal/plays"
	ty "alanktwong/dum/internal/types"
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

func TestPlayBook_GetPlays_PlayBookDisabled_PlayEnabled(t *testing.T) {
	attr := &ty.Attributes{ID: "test-pb", Enabled: false}
	plays := []*pl.Play{
		{Attributes: ty.Attributes{ID: "play-1", Enabled: true}},
	}
	pb, err := NewPlayBook(attr, plays, nil)
	assert.NoError(t, err)

	result := pb.GetPlays(true)
	assert.Equal(t, 1, result.Len())
	assert.True(t, result.Has("play-1"))
}

func TestPlayBook_GetAttributes(t *testing.T) {
	attr := &ty.Attributes{ID: "test-pb", Description: "test desc", Enabled: true, Sudo: false}
	pb, err := NewPlayBook(attr, nil, nil)
	assert.NoError(t, err)

	got := pb.GetAttributes()
	assert.Equal(t, *attr, got)
}

func TestPlayBook_GetID(t *testing.T) {
	attr := &ty.Attributes{ID: "test-pb"}
	pb, err := NewPlayBook(attr, nil, nil)
	assert.NoError(t, err)

	got := pb.GetID()
	assert.Equal(t, "test-pb", got)
}

func TestPlayBook_GetJetBrainsApps(t *testing.T) {
	attr := &ty.Attributes{ID: "test-pb"}
	apps := map[string]string{"IntelliJ": "/Applications/IntelliJ"}
	pb, err := NewPlayBook(attr, nil, apps)
	assert.NoError(t, err)

	got := pb.GetJetBrainsApps()
	assert.Equal(t, apps, got)
}
