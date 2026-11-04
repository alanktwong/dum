// Package tarball creates and restores verified tar archives.
package tarball

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"

	yamlv3 "gopkg.in/yaml.v3"
)

// Input describes the source paths used by list-based backup and restore.
type Input struct {
	Paths []string `yaml:"paths"`
}

// LoadInput reads and validates an input.yml file. Relative paths in the file
// resolve from the directory containing the input file.
func LoadInput(ctx context.Context, inputPath string) (*Input, error) {
	if ctx == nil {
		return nil, fmt.Errorf("load input: nil context")
	}
	if err := ctx.Err(); err != nil {
		return nil, fmt.Errorf("load input: %w", err)
	}

	absoluteInputPath, err := filepath.Abs(inputPath)
	if err != nil {
		return nil, fmt.Errorf("resolve input path %q: %w", inputPath, err)
	}

	data, err := os.ReadFile(absoluteInputPath)
	if err != nil {
		return nil, fmt.Errorf("read input file %q: %w", absoluteInputPath, err)
	}

	decoder := yamlv3.NewDecoder(bytes.NewReader(data))
	decoder.KnownFields(true)

	var input Input
	if err := decoder.Decode(&input); err != nil {
		return nil, fmt.Errorf("parse input file %q: %w", absoluteInputPath, err)
	}
	var extra any
	if err := decoder.Decode(&extra); !errors.Is(err, io.EOF) {
		if err != nil {
			return nil, fmt.Errorf("parse input file %q: %w", absoluteInputPath, err)
		}
		return nil, fmt.Errorf("parse input file %q: multiple YAML documents are not supported", absoluteInputPath)
	}

	if len(input.Paths) == 0 {
		return nil, fmt.Errorf("input file %q must contain at least one path", absoluteInputPath)
	}

	inputDir := filepath.Dir(absoluteInputPath)
	seen := make(map[string]struct{}, len(input.Paths))
	for index, path := range input.Paths {
		if path == "" {
			return nil, fmt.Errorf("input file %q has an empty path at index %d", absoluteInputPath, index)
		}

		resolvedPath := path
		if !filepath.IsAbs(resolvedPath) {
			resolvedPath = filepath.Join(inputDir, resolvedPath)
		}
		resolvedPath = filepath.Clean(resolvedPath)
		resolvedPath, err = filepath.Abs(resolvedPath)
		if err != nil {
			return nil, fmt.Errorf("resolve input path %q: %w", path, err)
		}
		if _, ok := seen[resolvedPath]; ok {
			return nil, fmt.Errorf("input file %q contains duplicate path %q", absoluteInputPath, resolvedPath)
		}
		seen[resolvedPath] = struct{}{}
		input.Paths[index] = resolvedPath
	}

	return &input, nil
}
