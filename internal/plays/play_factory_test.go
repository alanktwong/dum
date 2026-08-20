package plays

import (
	"testing"

	yml "github.com/alanktwong/dum/internal/yaml"

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
	playYAML := yml.PlayYAML{
		ID:          "test-play",
		Description: "a test play",
		Enabled:     true,
		Tasks: []yml.TaskYAML{
			{
				ID:     "test-link",
				Type:   "link",
				Root:   "/tmp/root",
				Target: "/tmp/target",
			},
		},
	}
	play, err := f.providePlay(playYAML)
	assert.NoError(t, err)
	assert.NotNil(t, play)
	assert.Equal(t, "test-play", play.ID)
	assert.Equal(t, "a test play", play.Description)
	assert.True(t, play.Enabled)
	assert.Equal(t, 1, len(play.Tasks))
}

func TestPlayFactory_ProvidePlay_NoTasks(t *testing.T) {
	f := NewPlayFactory()
	playYAML := yml.PlayYAML{
		ID:          "test-play",
		Description: "a test play",
		Tasks:       []yml.TaskYAML{},
	}
	play, err := f.providePlay(playYAML)
	assert.Error(t, err)
	assert.Nil(t, play)
	assert.Contains(t, err.Error(), "has no tasks")
}

func TestPlayFactory_ProvidePlay_Disabled(t *testing.T) {
	f := NewPlayFactory()
	playYAML := yml.PlayYAML{
		ID:      "test-play",
		Enabled: false,
		Tasks: []yml.TaskYAML{
			{
				ID:     "test-link",
				Type:   "link",
				Root:   "/tmp/root",
				Target: "/tmp/target",
			},
		},
	}
	play, err := f.providePlay(playYAML)
	assert.NoError(t, err)
	assert.NotNil(t, play)
	assert.False(t, play.IsEnabled())
}

func TestPlayFactory_ProvidePlays_Empty(t *testing.T) {
	f := NewPlayFactory()
	playsYAML := []yml.PlayYAML{}
	plays, err := f.ProvidePlays(playsYAML)
	assert.NoError(t, err)
	assert.Empty(t, plays)
}

func TestPlayFactory_ProvidePlays_Multiple(t *testing.T) {
	f := NewPlayFactory()
	playsYAML := []yml.PlayYAML{
		{
			ID: "play-1",
			Tasks: []yml.TaskYAML{
				{
					ID:     "link-1",
					Type:   "link",
					Root:   "/tmp/root1",
					Target: "/tmp/target1",
				},
			},
		},
		{
			ID: "play-2",
			Tasks: []yml.TaskYAML{
				{
					ID:     "link-2",
					Type:   "link",
					Root:   "/tmp/root2",
					Target: "/tmp/target2",
				},
			},
		},
	}
	plays, err := f.ProvidePlays(playsYAML)
	assert.NoError(t, err)
	assert.Equal(t, 2, len(plays))
	assert.Equal(t, "play-1", plays[0].ID)
	assert.Equal(t, "play-2", plays[1].ID)
}
