// Package factory provides the Factory for constructing Input.
package factory

import (
	ext "awong/dotfiles/internal/external"
	l "awong/dotfiles/internal/logging"
	pb "awong/dotfiles/internal/playbook"
	pl "awong/dotfiles/internal/plays"
	ty "awong/dotfiles/internal/types"
	tyg "awong/dotfiles/internal/types/gen"
	yml "awong/dotfiles/internal/yaml"
	"fmt"
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
	cfg, err := f.loadFromTypedYAML(options.File)
	if err != nil {
		return nil, err
	}
	playBook, err := f.ProvidePlayBook(&cfg.PlayBook)
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

// ProvidePlayBook constructs a PlayBook given a PlayBookYAML.
func (f *Factory) ProvidePlayBook(pbYAML *yml.PlayBookYAML) (*pb.PlayBook, error) {
	plays, err := f.PlayFactory.ProvidePlays(pbYAML.Plays)
	if err != nil {
		return nil, fmt.Errorf("failed to provide plays: %w", err)
	}
	apps := f.provideJetBrainsAppsFromYAML(pbYAML.JetBrains)
	attributes, err := ty.NewAttributes(
		pbYAML.ID,
		pbYAML.Description,
		pbYAML.Enabled,
		pbYAML.Sudo,
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

func (f *Factory) provideJetBrainsAppsFromYAML(jetbrains []yml.JetBrainsYAML) map[string]string {
	apps := make(map[string]string)
	appkeys := tyg.JetBrainsTypeNames()

	for _, each := range appkeys {
		for _, m := range jetbrains {
			version := m[each]
			if version != "" {
				apps[each] = version
			}
		}
	}
	return apps
}

func (f *Factory) loadFromTypedYAML(file string) (*yml.Config, error) {
	cfg, err := yml.Load(file)
	if err != nil {
		return nil, fmt.Errorf("failed to load typed YAML from %s: %w", file, err)
	}
	return cfg, nil
}
