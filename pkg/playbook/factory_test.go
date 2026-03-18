package playbook

import (
	e "awong/dotfiles/pkg/enums"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFactory_Provide(t *testing.T) {
	options := InputOptions{
		File:   "../../pkg/testdata/test_installer.yml",
		Group:  "test",
		DryRun: true,
	}
	f := NewFactory()
	want := NewInput(options.DryRun, options.Group, &PlayBook{
		Attributes: Attributes{
			ID:          "my-test-playbook",
			Description: "My test playbook",
			Enabled:     true,
			Sudo:        false,
		},
		JetBrainsApps: map[string]string{
			e.JetBrainsTypeDatagrip.String():  "DataGrip2025.1",
			e.JetBrainsTypeGoland.String():    "GoLand2025.1",
			e.JetBrainsTypeIdea.String():      "IntelliJIdea2025.1",
			e.JetBrainsTypePycharm.String():   "PyCharm2025.1",
			e.JetBrainsTypeRustrover.String(): "RustRover2025.1",
			e.JetBrainsTypeWebstorm.String():  "WebStorm2025.1",
		},
		Plays: []*Play{
			{
				Attributes: Attributes{
					ID:          "play-1",
					Description: "Play 1",
					Enabled:     true,
					Sudo:        false,
				},
				Tasks: []Task{
					createTestDirTask(t,
						Attributes{
							ID:          "~/.dotfiles",
							Description: "~/.dotfiles",
							Enabled:     true,
							Sudo:        false,
						},
					),
					createTestLinkTask(t,
						Attributes{
							ID:          "../projects/dotfiles",
							Description: "../projects/dotfiles",
							Enabled:     true,
							Sudo:        false,
						},
						"~/.dotfiles",
						"",
					),
					createTestBrewTask(t,
						Attributes{
							ID:          "scc",
							Description: "scc",
							Enabled:     true,
							Sudo:        false,
						},
					),
					createTestFunctionTask(t,
						Attributes{
							ID:          "install_test",
							Description: "install_test",
							Enabled:     true,
							Sudo:        false,
						},
					),
				},
			},
			{
				Attributes: Attributes{
					ID:          "play-2",
					Description: "Play 2",
					Enabled:     false,
					Sudo:        false,
				},
				Tasks: []Task{
					createTestGitTask(t,
						Attributes{
							ID:          "https://github.com/assertj/assertj.git",
							Description: "https://github.com/assertj/assertj.git",
							Enabled:     false,
							Sudo:        false,
						},
						"~/projects",
						"",
					),
					createTestMasTask(t,
						Attributes{
							ID:          "462058435",
							Description: "MS Excel",
							Enabled:     true,
							Sudo:        false,
						},
					),
				},
			},
		},
	})

	// when
	got, err := f.Provide(options)
	// then
	assert.NoError(t, err)
	assert.Equal(t, want.Task, got.Task, "Provide(%v): Task", options)
	assert.Equal(t, want.DryRun, got.DryRun, "Provide(%v): DryRun", options)
	gotPb := got.PlayBook
	assert.Equal(t, want.PlayBook.ID, gotPb.ID, "Provide(%v): PlayBook.ID", options)
	assert.Equal(t, want.PlayBook.Description, gotPb.Description, "Provide(%v): PlayBook.Description", options)
	assert.Equal(t, want.PlayBook.Enabled, gotPb.Enabled, "Provide(%v): PlayBook.Enabled", options)
	assert.Equal(t, want.PlayBook.Sudo, gotPb.Sudo, "Provide(%v): PlayBook.Sudo", options)
	assert.Equal(t, want.PlayBook.JetBrainsApps, gotPb.JetBrainsApps, "Provide(%v): PlayBook.JetBrainsApp", options)
}

func TestFactory_ProvidePlayBook(t *testing.T) {
	type args struct {
		yml map[string]interface{}
	}
	f := NewFactory()
	test := []struct {
		name    string
		args    args
		want    *PlayBook
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "ProvidePlayBook should provide playbook",
			args: args{
				yml: map[string]interface{}{
					"playbook": map[string]interface{}{
						"id":          "default",
						"description": "Default Playbook",
						"enabled":     true,
						"jetbrains": []interface{}{
							map[string]interface{}{
								e.JetBrainsTypeIdea.String(): "2023.4",
							},
						},
						"plays": []interface{}{
							map[string]interface{}{
								"id":          "play-1",
								"description": "Play 1",
								"enabled":     true,
								"tasks": []interface{}{
									map[string]interface{}{
										"type": "dir",
										"id":   "~/.dotfiles",
									},
									map[string]interface{}{
										"type": "function",
										"id":   "install_test",
									},
								},
							},
							map[string]interface{}{
								"id":          "play-2",
								"description": "Play 2",
								"enabled":     false,
								"tasks": []interface{}{
									map[string]interface{}{
										"type": "brew",
										"id":   "scc",
									},
									map[string]interface{}{
										"type": "git",
										"id":   "https://github.com/assertj/assertj.git",
									},
								},
							},
						},
					},
				},
			},
			want: &PlayBook{
				Attributes: Attributes{
					ID:          "default",
					Description: "Default Playbook",
					Enabled:     true,
					Sudo:        false,
				},
				JetBrainsApps: map[string]string{
					e.JetBrainsTypeIdea.String(): "2023.4",
				},
				Plays: []*Play{
					{
						Attributes: Attributes{
							ID:          "play-1",
							Description: "Play 1",
							Enabled:     true,
							Sudo:        false,
						},
						Tasks: []Task{
							createTestDirTask(t,
								Attributes{
									ID:          "~/.dotfiles",
									Description: "~/.dotfiles",
									Enabled:     true,
									Sudo:        false,
								},
							),
							createTestFunctionTask(t,
								Attributes{
									ID:          "install_test",
									Description: "install_test",
									Enabled:     true,
									Sudo:        false,
								},
							),
						},
					},
					{
						Attributes: Attributes{
							ID:          "play-2",
							Description: "Play 2",
							Enabled:     false,
							Sudo:        false,
						},
						Tasks: []Task{
							createTestBrewTask(t,
								Attributes{
									ID:          "scc",
									Description: "scc",
									Enabled:     true,
									Sudo:        false,
								},
							),
							createTestGitTask(t,
								Attributes{
									ID:          "https://github.com/assertj/assertj.git",
									Description: "https://github.com/assertj/assertj.git",
									Enabled:     true,
									Sudo:        false,
								},
								"~/projects",
								"",
							),
						},
					},
				},
			},
			wantErr: assert.NoError,
		},
	}
	for _, tt := range test {
		t.Run(tt.name, func(t *testing.T) {
			got, err := f.ProvidePlayBook(tt.args.yml)
			if !tt.wantErr(t, err, "providePlayBook(%v)", tt.args.yml) {
				return
			}
			assert.Equalf(t, tt.want, got, "providePlayBook(%v)", tt.args.yml)
		})
	}
}

func TestFactory_providePlays(t *testing.T) {
	type args struct {
		yml map[string]interface{}
	}
	f := NewFactory()
	test := []struct {
		name    string
		args    args
		want    []*Play
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "ProvidePlays should provide plays",
			args: args{
				yml: map[string]interface{}{
					"plays": []interface{}{
						map[string]interface{}{
							"id":          "play1",
							"description": "Play 1",
							"enabled":     true,
							"tasks": []interface{}{
								map[string]interface{}{
									"type": "dir",
									"id":   "~/.dotfiles",
								},
							},
						},
						map[string]interface{}{
							"id":          "play2",
							"description": "Play 2",
							"enabled":     false,
							"tasks": []interface{}{
								map[string]interface{}{
									"type": "dir",
									"id":   "~/.config",
								},
							},
						},
					},
				},
			},
			want: []*Play{
				{
					Attributes: Attributes{
						ID:          "play1",
						Description: "Play 1",
						Enabled:     true,
						Sudo:        false,
					},
					Tasks: []Task{
						createTestDirTask(t,
							Attributes{
								ID:          "~/.dotfiles",
								Description: "~/.dotfiles",
								Enabled:     true,
								Sudo:        false,
							},
						),
					},
				},
				{
					Attributes: Attributes{
						ID:          "play2",
						Description: "Play 2",
						Enabled:     false,
						Sudo:        false,
					},
					Tasks: []Task{
						createTestDirTask(t,
							Attributes{
								ID:          "~/.config",
								Description: "~/.config",
								Enabled:     true,
								Sudo:        false,
							},
						),
					},
				},
			},
			wantErr: assert.NoError,
		},
		{
			name: "ProvidePlays should NOT provide play when it has no tasks",
			args: args{
				yml: map[string]interface{}{
					"plays": []interface{}{
						map[string]interface{}{
							"id":          "play1",
							"description": "Play 1",
							"enabled":     true,
						},
					},
				},
			},
			want:    nil,
			wantErr: assert.Error,
		},
	}
	for _, tt := range test {
		t.Run(tt.name, func(t *testing.T) {
			got, err := f.providePlays(tt.args.yml)
			if !tt.wantErr(t, err, "providePlays(%v)", tt.args.yml) {
				return
			}
			assert.Equalf(t, tt.want, got, "providePlays(%v)", tt.args.yml)
		})
	}
}

func TestFactory_provideTask(t *testing.T) {
	type args struct {
		yml map[string]interface{}
	}
	f := NewFactory()
	test := []struct {
		name    string
		args    args
		want    Task
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "should NOT provide invalid task",
			args: args{
				yml: map[string]interface{}{
					"type": "invalid",
					"id":   "invalid-task",
				},
			},
			want:    nil,
			wantErr: assert.Error,
		},
		{
			name: "should provide dir task",
			args: args{
				yml: map[string]interface{}{
					"type": "dir",
					"id":   "~/.dotfiles",
				},
			},
			want: createTestDirTask(t,
				Attributes{
					ID:          "~/.dotfiles",
					Description: "~/.dotfiles",
					Enabled:     true,
					Sudo:        false,
				},
			),
			wantErr: assert.NoError,
		},
		{
			name: "should provide brew task",
			args: args{
				yml: map[string]interface{}{
					"type": "brew",
					"id":   "scc",
				},
			},
			want: createTestBrewTask(t,
				Attributes{
					ID:          "scc",
					Description: "scc",
					Enabled:     true,
					Sudo:        false,
				}),
			wantErr: assert.NoError,
		},
		{
			name: "should provide cask task",
			args: args{
				yml: map[string]interface{}{
					"type":    "cask",
					"id":      "grandperspective",
					"enabled": false,
				},
			},
			want: createTestBrewCaskTask(t,
				Attributes{
					ID:          "grandperspective",
					Description: "grandperspective",
					Enabled:     false,
					Sudo:        false,
				}),
			wantErr: assert.NoError,
		},
		{
			name: "should provide cellar task",
			args: args{
				yml: map[string]interface{}{
					"type":    "cellar",
					"id":      "sqlite",
					"enabled": false,
				},
			},
			want: createTestBrewCellarTask(t,
				Attributes{
					ID:          "sqlite",
					Description: "sqlite",
					Enabled:     false,
					Sudo:        false,
				},
			),
			wantErr: assert.NoError,
		},
		{
			name: "should provide link task",
			args: args{
				yml: map[string]interface{}{
					"type": "link",
					"id":   "~/.dotfiles",
				},
			},
			want: createTestLinkTask(t,
				Attributes{
					ID:          "~/.dotfiles",
					Description: "~/.dotfiles",
					Enabled:     true,
					Sudo:        false,
				},
				"",
				"",
			),
			wantErr: assert.NoError,
		},
		{
			name: "should provide function task",
			args: args{
				yml: map[string]interface{}{
					"type": "function",
					"id":   "install_test",
				},
			},
			want: createTestFunctionTask(t,
				Attributes{
					ID:          "install_test",
					Description: "install_test",
					Enabled:     true,
					Sudo:        false,
				},
			),
			wantErr: assert.NoError,
		},
		{
			name: "should provide git task",
			args: args{
				yml: map[string]interface{}{
					"type": "git",
					"id":   "https://github.com/assertj/assertj.git",
				},
			},
			want: createTestGitTask(t,
				Attributes{
					ID:          "https://github.com/assertj/assertj.git",
					Description: "https://github.com/assertj/assertj.git",
					Enabled:     true,
					Sudo:        false,
				},
				"~/projects",
				"",
			),
			wantErr: assert.NoError,
		},
		{
			name: "should provide jetbrains task",
			args: args{
				yml: map[string]interface{}{
					"type":        "jetbrains",
					"id":          "org.asciidoctor.intellij.asciidoc",
					"description": "asciidoc",
					"apps": []interface{}{
						e.JetBrainsTypeIdea.String(),
					},
				},
			},
			want: createTestJetBrainsTask(t,
				Attributes{
					ID:          "org.asciidoctor.intellij.asciidoc",
					Description: "asciidoc",
					Enabled:     true,
					Sudo:        false,
				},
				[]string{e.JetBrainsTypeIdea.String()},
			),
			wantErr: assert.NoError,
		},
		{
			name: "should provide mas task",
			args: args{
				yml: map[string]interface{}{
					"type":        "mas",
					"id":          "462058435",
					"description": "MS Excel",
				},
			},
			want: createTestMasTask(t,
				Attributes{
					ID:          "462058435",
					Description: "MS Excel",
					Enabled:     true,
					Sudo:        false,
				},
			),
			wantErr: assert.NoError,
		},
		{
			name: "should provide vscode task",
			args: args{
				yml: map[string]interface{}{
					"type":        "vscode",
					"id":          "asciidoctor.asciidoctor-vscode",
					"description": "Asciidoctor",
				},
			},
			want: createTestVsCodeTask(t,
				Attributes{
					ID:          "asciidoctor.asciidoctor-vscode",
					Description: "Asciidoctor",
					Enabled:     true,
					Sudo:        false,
				},
			),
			wantErr: assert.NoError,
		},
	}
	for _, tt := range test {
		t.Run(tt.name, func(t *testing.T) {
			got, err := f.provideTask(tt.args.yml)
			if !tt.wantErr(t, err, "provideTask(%v)", tt.args.yml) {
				return
			}
			assert.Equalf(t, tt.want, got, "provideTask(%v)", tt.args.yml)
		})
	}
}

func TestFactory_provideJetBrainsApps(t *testing.T) {
	type args struct {
		yml map[string]interface{}
	}
	f := NewFactory()
	test := []struct {
		name    string
		args    args
		want    map[string]string
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "provideJetBrainsApps should return per configuration",
			args: args{
				yml: map[string]interface{}{
					"jetbrains": []interface{}{
						map[string]interface{}{
							e.JetBrainsTypePhpstorm.String(): "PhpStorm2023.5",
						},
						map[string]interface{}{
							e.JetBrainsTypePycharm.String(): "PyCharm2023.6",
						},
						map[string]interface{}{
							e.JetBrainsTypeClion.String(): "CLion2023.1",
						},
						map[string]interface{}{
							e.JetBrainsTypeRubymine.String(): "RubyMine2023.8",
						},
						map[string]interface{}{
							e.JetBrainsTypeRustrover.String(): "RustRover2023.9",
						},
						map[string]interface{}{
							e.JetBrainsTypeWebstorm.String(): "WebStorm2024.1",
						},
						map[string]interface{}{
							"excel": "Excel 2024.2",
						},
						map[string]interface{}{
							e.JetBrainsTypeDatagrip.String(): "DataGrip2023.2",
						},
						map[string]interface{}{
							e.JetBrainsTypeGoland.String(): "GoLand2023.3",
						},
						map[string]interface{}{
							e.JetBrainsTypeIdea.String(): "IntelliJ2023.4",
						},
						map[string]interface{}{
							e.JetBrainsTypeRider.String(): "Rider2023.7",
						},
					},
				},
			},
			want: map[string]string{
				e.JetBrainsTypeClion.String():     "CLion2023.1",
				e.JetBrainsTypeDatagrip.String():  "DataGrip2023.2",
				e.JetBrainsTypeGoland.String():    "GoLand2023.3",
				e.JetBrainsTypeIdea.String():      "IntelliJ2023.4",
				e.JetBrainsTypePhpstorm.String():  "PhpStorm2023.5",
				e.JetBrainsTypePycharm.String():   "PyCharm2023.6",
				e.JetBrainsTypeRider.String():     "Rider2023.7",
				e.JetBrainsTypeRubymine.String():  "RubyMine2023.8",
				e.JetBrainsTypeRustrover.String(): "RustRover2023.9",
				e.JetBrainsTypeWebstorm.String():  "WebStorm2024.1",
			},
			wantErr: assert.NoError,
		},
		{
			name: "provideJetBrainsApps should return nothing per configuration",
			args: args{
				yml: map[string]interface{}{
					"jetbrains": []interface{}{
						map[string]interface{}{
							"excel": "Excel 2024.2",
						},
					},
				},
			},
			want:    map[string]string{},
			wantErr: assert.NoError,
		},
		{
			name: "provideJetBrainsApps should return nothing per incorrect map",
			args: args{
				yml: map[string]interface{}{
					"jetbrains": map[string]interface{}{
						e.JetBrainsTypeGoland.String(): "GoLand2023.3",
					},
				},
			},
			want:    map[string]string{},
			wantErr: assert.NoError,
		},
	}

	for _, tt := range test {
		t.Run(tt.name, func(t *testing.T) {
			got := f.provideJetBrainsApps(tt.args.yml)
			assert.Equalf(t, tt.want, got, "provideJetBrainsApps(%v)", tt.args.yml)
		})
	}
}

func TestFactory_getString(t *testing.T) {
	type args struct {
		yml map[string]interface{}
		key string
		def string
	}
	f := NewFactory()
	test := []struct {
		name string
		args args
		want string
	}{
		{
			name: "getString should return default value if key not found",
			args: args{
				yml: map[string]interface{}{"key1": "value1", "key2": 42},
				key: "key3",
				def: "default",
			},
			want: "default",
		},
		{
			name: "getString should return default value if value is not a string",
			args: args{
				yml: map[string]interface{}{"key1": "value1", "key2": 42},
				key: "key2",
				def: "default",
			},
			want: "default",
		},
		{
			name: "getString should return value if key found",
			args: args{
				yml: map[string]interface{}{"key1": "value1", "key2": 42},
				key: "key1",
				def: "default",
			},
			want: "value1",
		},
	}
	for _, tt := range test {
		t.Run(tt.name, func(t *testing.T) {
			got := f.getString(tt.args.yml, tt.args.key, tt.args.def)
			assert.Equalf(t, tt.want, got, "getString(%v, %v, %v)", tt.args.yml, tt.args.key, tt.args.def)
		})
	}
}

func TestFactory_getStrings(t *testing.T) {
	type args struct {
		yml map[string]interface{}
		key string
		def []string
	}
	f := NewFactory()
	test := []struct {
		name string
		args args
		want []string
	}{
		{
			name: "getStrings should return default value if key not found",
			args: args{
				yml: map[string]interface{}{
					"key1": []interface{}{
						"foo",
						"bar",
					},
					"key2": 42,
				},
				key: "key3",
				def: []string{"value1", "value2"},
			},
			want: []string{"value1", "value2"},
		},
		{
			name: "getStrings should return default value if value is not a string",
			args: args{
				yml: map[string]interface{}{
					"key1": []interface{}{
						"foo",
						"bar",
					},
					"key2": 42,
				},
				key: "key2",
				def: []string{"value1", "value2"},
			},
			want: []string{"value1", "value2"},
		},
		{
			name: "getStrings should return value if key found",
			args: args{
				yml: map[string]interface{}{
					"key1": []interface{}{
						"foo",
						"bar",
					},
					"key2": 42,
				},
				key: "key1",
				def: []string{"foo", "bar"},
			},
			want: []string{"foo", "bar"},
		},
	}
	for _, tt := range test {
		t.Run(tt.name, func(t *testing.T) {
			got := f.getStrings(tt.args.yml, tt.args.key, tt.args.def)
			assert.Equalf(t, tt.want, got, "getString(%v, %v, %v)", tt.args.yml, tt.args.key, tt.args.def)
		})
	}
}

func TestFactory_getBool(t *testing.T) {
	type args struct {
		yml map[string]interface{}
		key string
		def bool
	}
	f := NewFactory()
	test := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "getBool should return default value if key not found",
			args: args{
				yml: map[string]interface{}{"key1": true, "key2": 42},
				key: "key3",
				def: false,
			},
			want: false,
		},
		{
			name: "getBool should return default value if value is not a bool",
			args: args{
				yml: map[string]interface{}{"key1": true, "key2": 42},
				key: "key2",
				def: false,
			},
			want: false,
		},
		{
			name: "getBool should return value if key found",
			args: args{
				yml: map[string]interface{}{"key1": true, "key2": 42},
				key: "key1",
				def: false,
			},
			want: true,
		},
	}
	for _, tt := range test {
		t.Run(tt.name, func(t *testing.T) {
			got := f.getBool(tt.args.yml, tt.args.key, tt.args.def)
			assert.Equalf(t, tt.want, got, "getBool(%v, %v, %v)", tt.args.yml, tt.args.key, tt.args.def)
		})
	}
}
