package factory

import (
	"awong/dotfiles/internal/tasks"
	"os"
	"path/filepath"
	"testing"

	yml "awong/dotfiles/internal/yaml"

	"github.com/stretchr/testify/assert"
)

func TestNewFactory(t *testing.T) {
	f := NewFactory()
	assert.NotNil(t, f)
	assert.NotNil(t, f.Log)
	assert.NotNil(t, f.Utils)
	assert.NotNil(t, f.PlayFactory)
}

func TestFactory_Provide_HappyPath(t *testing.T) {
	f := NewFactory()
	input, err := f.Provide(InputOptions{
		File:   "testdata/test_installer.yml",
		Group:  "play-1",
		DryRun: true,
	})
	assert.NoError(t, err)
	assert.NotNil(t, input)
	assert.True(t, input.DryRun)
	assert.Equal(t, "play-1", input.Play)
	assert.False(t, input.Sudo)

	pb := input.PlayBook
	assert.NotNil(t, pb)
	assert.Equal(t, "my-test-playbook", pb.ID)
	assert.Equal(t, "My test playbook", pb.Description)
	assert.True(t, pb.Enabled)
	assert.False(t, pb.Sudo)
	assert.Len(t, pb.Plays, 2)
	assert.Len(t, pb.JetBrainsApps, 6)
}

func TestFactory_Provide_FileNotFound(t *testing.T) {
	f := NewFactory()
	_, err := f.Provide(InputOptions{
		File: "/nonexistent/path/installer.yml",
	})
	assert.Error(t, err)
}

func TestFactory_Provide_InvalidYaml(t *testing.T) {
	f := NewFactory()

	tmpDir := t.TempDir()
	badFile := filepath.Join(tmpDir, "bad.yml")
	err := os.WriteFile(badFile, []byte("{{invalid yaml: [}"), 0o644)
	assert.NoError(t, err)

	_, err = f.Provide(InputOptions{File: badFile})
	assert.Error(t, err)
}

func TestFactory_ProvidePlayBook_Success(t *testing.T) {
	f := NewFactory()
	pbYAML := &yml.PlayBookYAML{
		ID:          "test-pb",
		Description: "a test playbook",
		Enabled:     true,
		Plays: []yml.PlayYAML{
			{
				ID: "play-1",
				Tasks: []yml.TaskYAML{
					{
						ID:   "dir-1",
						Type: "dir",
					},
				},
			},
		},
	}
	result, err := f.ProvidePlayBook(pbYAML)
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "test-pb", result.ID)
	assert.Equal(t, "a test playbook", result.Description)
	assert.True(t, result.Enabled)
	assert.Len(t, result.Plays, 1)
}

func TestFactory_ProvidePlayBook_MissingID(t *testing.T) {
	f := NewFactory()
	pbYAML := &yml.PlayBookYAML{
		Plays: []yml.PlayYAML{
			{
				ID: "play-1",
				Tasks: []yml.TaskYAML{
					{
						ID:   "dir-1",
						Type: "dir",
					},
				},
			},
		},
	}
	_, err := f.ProvidePlayBook(pbYAML)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "attribute ID cannot be empty")
}

func TestFactory_ProvidePlayBook_NoPlays(t *testing.T) {
	f := NewFactory()
	pbYAML := &yml.PlayBookYAML{
		ID:          "no-plays",
		Description: "empty playbook",
	}
	result, err := f.ProvidePlayBook(pbYAML)
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "no-plays", result.ID)
	assert.Empty(t, result.Plays)
}

func TestFactory_ProvidePlayBook_MultiplePlays(t *testing.T) {
	f := NewFactory()
	pbYAML := &yml.PlayBookYAML{
		ID: "multi",
		Plays: []yml.PlayYAML{
			{
				ID: "play-a",
				Tasks: []yml.TaskYAML{
					{ID: "t1", Type: "dir"},
				},
			},
			{
				ID: "play-b",
				Tasks: []yml.TaskYAML{
					{ID: "t2", Type: "dir"},
				},
			},
		},
	}
	result, err := f.ProvidePlayBook(pbYAML)
	assert.NoError(t, err)
	assert.Len(t, result.Plays, 2)
	assert.Equal(t, "play-a", result.Plays[0].ID)
	assert.Equal(t, "play-b", result.Plays[1].ID)
}

func TestFactory_ProvidePlayBook_DisabledPlay(t *testing.T) {
	f := NewFactory()
	pbYAML := &yml.PlayBookYAML{
		ID: "with-disabled",
		Plays: []yml.PlayYAML{
			{
				ID:      "disabled-play",
				Enabled: false,
				Tasks: []yml.TaskYAML{
					{ID: "t1", Type: "dir"},
				},
			},
		},
	}
	result, err := f.ProvidePlayBook(pbYAML)
	assert.NoError(t, err)
	assert.Len(t, result.Plays, 1)
	assert.False(t, result.Plays[0].Enabled)
}

func TestFactory_ProvidePlayBook_Sudo(t *testing.T) {
	f := NewFactory()
	pbYAML := &yml.PlayBookYAML{
		ID:   "sudo-pb",
		Sudo: true,
		Plays: []yml.PlayYAML{
			{
				ID: "play-1",
				Tasks: []yml.TaskYAML{
					{ID: "t1", Type: "dir"},
				},
			},
		},
	}
	result, err := f.ProvidePlayBook(pbYAML)
	assert.NoError(t, err)
	assert.True(t, result.Sudo)
}

func TestFactory_ProvidePlayBook_JetBrainsApps(t *testing.T) {
	f := NewFactory()
	pbYAML := &yml.PlayBookYAML{
		ID: "jetbrains-pb",
		JetBrains: []yml.JetBrainsYAML{
			{"goland": "GoLand2025.1"},
			{"idea": "IntelliJIdea2025.1"},
		},
		Plays: []yml.PlayYAML{
			{
				ID: "play-1",
				Tasks: []yml.TaskYAML{
					{ID: "t1", Type: "dir"},
				},
			},
		},
	}
	result, err := f.ProvidePlayBook(pbYAML)
	assert.NoError(t, err)
	assert.Equal(t, "GoLand2025.1", result.JetBrainsApps["goland"])
	assert.Equal(t, "IntelliJIdea2025.1", result.JetBrainsApps["idea"])
	assert.Len(t, result.JetBrainsApps, 2)
}

func TestFactory_ProvidePlayBook_EmptyID(t *testing.T) {
	f := NewFactory()
	pbYAML := &yml.PlayBookYAML{
		Plays: []yml.PlayYAML{
			{
				ID: "play-1",
				Tasks: []yml.TaskYAML{
					{ID: "t1", Type: "dir"},
				},
			},
		},
	}
	_, err := f.ProvidePlayBook(pbYAML)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "attribute ID cannot be empty")
}

func TestFactory_loadFromTypedYAML_Success(t *testing.T) {
	f := NewFactory()
	cfg, err := f.loadFromTypedYAML("testdata/test_installer.yml")
	assert.NoError(t, err)
	assert.NotNil(t, cfg)
	assert.Equal(t, "my-test-playbook", cfg.PlayBook.ID)
}

func TestFactory_loadFromTypedYAML_FileNotFound(t *testing.T) {
	f := NewFactory()
	_, err := f.loadFromTypedYAML("/nonexistent/path/installer.yml")
	assert.Error(t, err)
}

func TestFactory_loadFromTypedYAML_InvalidYaml(t *testing.T) {
	f := NewFactory()

	tmpDir := t.TempDir()
	badFile := filepath.Join(tmpDir, "bad.yml")
	err := os.WriteFile(badFile, []byte("{{invalid: [}"), 0o644)
	assert.NoError(t, err)

	_, err = f.loadFromTypedYAML(badFile)
	assert.Error(t, err)
}

func TestFactory_Provide_SudoFlag(t *testing.T) {
	f := NewFactory()
	input, err := f.Provide(InputOptions{
		File:  "testdata/test_installer.yml",
		Group: "play-1",
	})
	assert.NoError(t, err)
	assert.NotNil(t, input)

	assert.False(t, input.Sudo)
	assert.Equal(t, input.PlayBook.Sudo, input.Sudo)
}

func TestFactory_Provide_DryRunFalse(t *testing.T) {
	f := NewFactory()
	input, err := f.Provide(InputOptions{
		File:   "testdata/test_installer.yml",
		Group:  "play-1",
		DryRun: false,
	})
	assert.NoError(t, err)
	assert.False(t, input.DryRun)
}

func TestFactory_ProvidePlayBook_EmptyPlayBook(t *testing.T) {
	f := NewFactory()
	pbYAML := &yml.PlayBookYAML{}
	_, err := f.ProvidePlayBook(pbYAML)
	assert.Error(t, err)
}

func TestFactory_ProvidePlayBook_InvalidTaskType(t *testing.T) {
	f := NewFactory()
	pbYAML := &yml.PlayBookYAML{
		ID: "bad-task",
		Plays: []yml.PlayYAML{
			{
				ID: "play-1",
				Tasks: []yml.TaskYAML{
					{ID: "t1", Type: "unknown-type"},
				},
			},
		},
	}
	_, err := f.ProvidePlayBook(pbYAML)
	assert.Error(t, err)
}

func TestFactory_Provide_HappyPath_Plays(t *testing.T) {
	f := NewFactory()
	input, err := f.Provide(InputOptions{
		File: "testdata/test_installer.yml",
	})
	assert.NoError(t, err)
	assert.NotNil(t, input)

	plays := input.PlayBook.Plays
	assert.Len(t, plays, 2)

	assert.Equal(t, "play-1", plays[0].ID)
	assert.Equal(t, "Play 1", plays[0].Description)
	assert.True(t, plays[0].Enabled)

	assert.Equal(t, "play-2", plays[1].ID)
	assert.Equal(t, "Play 2", plays[1].Description)
	assert.False(t, plays[1].Enabled)
}

func TestFactory_Provide_HappyPath_TaskTypes(t *testing.T) {
	f := NewFactory()
	input, err := f.Provide(InputOptions{
		File: "testdata/test_installer.yml",
	})
	assert.NoError(t, err)

	play1 := input.PlayBook.Plays[0]
	assert.Len(t, play1.Tasks, 12)

	_, ok := play1.Tasks[0].(*tasks.DirTask)
	assert.True(t, ok, "expected DirTask")
	_, ok = play1.Tasks[1].(*tasks.LinkTask)
	assert.True(t, ok, "expected LinkTask")
	_, ok = play1.Tasks[2].(*tasks.BashTask)
	assert.True(t, ok, "expected BashTask")
	_, ok = play1.Tasks[4].(*tasks.BrewTask)
	assert.True(t, ok, "expected BrewTask")
	_, ok = play1.Tasks[5].(*tasks.BrewCaskTask)
	assert.True(t, ok, "expected BrewCaskTask")
	_, ok = play1.Tasks[6].(*tasks.BrewCellarTask)
	assert.True(t, ok, "expected BrewCellarTask")
	_, ok = play1.Tasks[7].(*tasks.FunctionTask)
	assert.True(t, ok, "expected FunctionTask")
	_, ok = play1.Tasks[8].(*tasks.GitTask)
	assert.True(t, ok, "expected GitTask")
	_, ok = play1.Tasks[9].(*tasks.MasTask)
	assert.True(t, ok, "expected MasTask")
	_, ok = play1.Tasks[10].(*tasks.VsCodePluginTask)
	assert.True(t, ok, "expected VsCodePluginTask")
	_, ok = play1.Tasks[11].(*tasks.JetBrainsPluginTask)
	assert.True(t, ok, "expected JetBrainsPluginTask")
}

func TestFactory_Provide_HappyPath_JetBrainsApps(t *testing.T) {
	f := NewFactory()
	input, err := f.Provide(InputOptions{
		File: "testdata/test_installer.yml",
	})
	assert.NoError(t, err)

	apps := input.PlayBook.JetBrainsApps
	assert.Equal(t, "DataGrip2025.1", apps["datagrip"])
	assert.Equal(t, "GoLand2025.1", apps["goland"])
	assert.Equal(t, "IntelliJIdea2025.1", apps["idea"])
	assert.Equal(t, "PyCharm2025.1", apps["pycharm"])
	assert.Equal(t, "RustRover2025.1", apps["rustrover"])
	assert.Equal(t, "WebStorm2025.1", apps["webstorm"])
	assert.Len(t, apps, 6)
}
