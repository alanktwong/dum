// Package factory provides the Factory for constructing Input.
package factory

import (
	ext "awong/dotfiles/pkg/external"
	l "awong/dotfiles/pkg/logging"
	pb "awong/dotfiles/pkg/playbook"
	pl "awong/dotfiles/pkg/plays"
	ty "awong/dotfiles/pkg/types"
	tyg "awong/dotfiles/pkg/types/gen"
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

// Factory provides the Factory for constructing Input from YML files.
type Factory struct {
	Log         l.Logger
	Utils       ext.Ext
	PlayFactory *pl.PlayFactory
}

// NewFactory returns a new Factory for constructing Input.
func NewFactory() *Factory {
	return &Factory{
		Log:         l.Log(),
		Utils:       ext.NewExt(),
		PlayFactory: pl.NewPlayFactory(),
	}
}

// InputOptions provides the parameterized struct for constructing a Input.
type InputOptions struct {
	File   string
	Group  string
	DryRun bool
}

// Provide constructs a Input given Input Options.
func (f *Factory) Provide(options InputOptions) (*pb.Input, error) {
	yml, err := f.getYaml(options.File)
	if err != nil {
		return nil, err
	}
	playBook, err := f.ProvidePlayBook(yml)
	if err != nil {
		return nil, fmt.Errorf("failed to provide playbook: %w", err)
	}
	return &pb.Input{
		DryRun:   options.DryRun,
		Play:     options.Group,
		PlayBook: playBook,
		Sudo:     playBook.Sudo,
	}, nil
}

func (f *Factory) getYaml(file string) (map[string]any, error) {
	f.Log.Debug("Loading playbook from file", "file", file)
	absPath, err := f.Utils.ToAbsolutePath(file)
	if err != nil {
		return nil, fmt.Errorf("failed to get absolute path for file %s: %w", file, err)
	}
	byteArray, err := os.ReadFile(absPath)
	if err != nil {
		return nil, fmt.Errorf("failed to get read file %s: %w", absPath, err)
	}
	var yamlMap map[string]any
	err = yaml.Unmarshal(byteArray, &yamlMap)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal playbook from file %s: %w", file, err)
	}
	return yamlMap, nil
}

// ProvidePlayBook constructs a PlayBook given a map that expresses the YML.
func (f *Factory) ProvidePlayBook(yml map[string]any) (*pb.PlayBook, error) {
	pbData, ok := yml["playbook"].(map[string]any)
	if !ok {
		return nil, fmt.Errorf("playbook key not found or invalid format in YAML")
	}

	plays, err := f.PlayFactory.ProvidePlays(pbData)
	if err != nil {
		return nil, fmt.Errorf("failed to provide plays: %w", err)
	}
	id := f.Utils.GetString(pbData, "id", "")
	description := f.Utils.GetString(pbData, "description", id)
	apps := f.provideJetBrainsApps(pbData)
	attributes, err := ty.NewAttributes(
		id,
		description,
		f.Utils.GetBool(pbData, "enabled", true),
		f.Utils.GetBool(pbData, "sudo", false),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create attributes: %w", err)
	}
	playBook, err := pb.NewPlayBook(attributes, plays, apps)
	if err != nil {
		return nil, fmt.Errorf("failed to create playbook: %w", err)
	}
	return playBook, nil
}

func (f *Factory) provideJetBrainsApps(yml map[string]any) map[string]string {
	apps := make(map[string]string)
	appkeys := tyg.JetBrainsTypeNames()

	if jetbrainsArray, ok := yml["jetbrains"].([]any); ok {
		for _, each := range appkeys {
			for _, it := range jetbrainsArray {
				if m, ok2 := it.(map[string]any); ok2 {
					version := f.Utils.GetString(m, each, "")
					if version != "" {
						apps[each] = version
					}
				}
			}
		}
	}
	return apps
}
