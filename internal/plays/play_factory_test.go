package plays

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewPlayFactory_NoArgs(t *testing.T) {
	f := NewPlayFactory()
	assert.NotNil(t, f)
	assert.NotNil(t, f.Log)
	assert.NotNil(t, f.Utils)
	assert.NotNil(t, f.TaskFactory)
}

func TestPlayFactory_ProvidePlay_Success(t *testing.T) {
	f := NewPlayFactory()
	yml := map[string]any{
		"id":          "test-play",
		"description": "a test play",
		"enabled":     true,
		"tasks": []any{
			map[string]any{
				"id":     "test-link",
				"type":   "link",
				"root":   "/tmp/root",
				"target": "/tmp/target",
			},
		},
	}
	play, err := f.ProvidePlay(yml)
	assert.NoError(t, err)
	assert.NotNil(t, play)
	assert.Equal(t, "test-play", play.ID)
	assert.Equal(t, "a test play", play.Description)
	assert.True(t, play.Enabled)
	assert.Equal(t, 1, len(play.Tasks))
}

func TestPlayFactory_ProvidePlay_NoTasks(t *testing.T) {
	f := NewPlayFactory()
	yml := map[string]any{
		"id":          "test-play",
		"description": "a test play",
	}
	play, err := f.ProvidePlay(yml)
	assert.Error(t, err)
	assert.Nil(t, play)
	assert.Contains(t, err.Error(), "has no tasks")
}

func TestPlayFactory_ProvidePlay_Disabled(t *testing.T) {
	f := NewPlayFactory()
	yml := map[string]any{
		"id":      "test-play",
		"enabled": false,
		"tasks": []any{
			map[string]any{
				"id":     "test-link",
				"type":   "link",
				"root":   "/tmp/root",
				"target": "/tmp/target",
			},
		},
	}
	play, err := f.ProvidePlay(yml)
	assert.NoError(t, err)
	assert.NotNil(t, play)
	assert.False(t, play.IsEnabled())
}

func TestPlayFactory_ProvidePlays_Empty(t *testing.T) {
	f := NewPlayFactory()
	yml := map[string]any{}
	plays, err := f.ProvidePlays(yml)
	assert.NoError(t, err)
	assert.Empty(t, plays)
}

func TestPlayFactory_ProvidePlays_Multiple(t *testing.T) {
	f := NewPlayFactory()
	yml := map[string]any{
		"plays": []any{
			map[string]any{
				"id": "play-1",
				"tasks": []any{
					map[string]any{
						"id":     "link-1",
						"type":   "link",
						"root":   "/tmp/root1",
						"target": "/tmp/target1",
					},
				},
			},
			map[string]any{
				"id": "play-2",
				"tasks": []any{
					map[string]any{
						"id":     "link-2",
						"type":   "link",
						"root":   "/tmp/root2",
						"target": "/tmp/target2",
					},
				},
			},
		},
	}
	plays, err := f.ProvidePlays(yml)
	assert.NoError(t, err)
	assert.Equal(t, 2, len(plays))
	assert.Equal(t, "play-1", plays[0].ID)
	assert.Equal(t, "play-2", plays[1].ID)
}
