package tarball

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"path"
	"sort"
	"strings"

	yamlv3 "gopkg.in/yaml.v3"
)

const manifestSchemaVersion = 1

// EntryType identifies an entry stored in a tar archive.
type EntryType string

const (
	// EntryTypeFile identifies a regular file entry.
	EntryTypeFile EntryType = "file"
	// EntryTypeDirectory identifies a directory entry.
	EntryTypeDirectory EntryType = "directory"
	// EntryTypeSymlink identifies a symbolic link entry.
	EntryTypeSymlink EntryType = "symlink"
)

// Manifest is the versioned inventory stored beside each archive.
type Manifest struct {
	SchemaVersion int             `yaml:"schema_version"`
	Exemptions    []string        `yaml:"exemptions,omitempty"`
	UseGitignore  bool            `yaml:"use_gitignore,omitempty"`
	Entries       []ManifestEntry `yaml:"entries"`
}

// ManifestEntry describes one archive member.
type ManifestEntry struct {
	Path    string    `yaml:"path"`
	Type    EntryType `yaml:"type"`
	Target  string    `yaml:"target,omitempty"`
	Unsafe  bool      `yaml:"unsafe,omitempty"`
	Warning string    `yaml:"warning,omitempty"`
	UID     *int64    `yaml:"uid,omitempty"`
	GID     *int64    `yaml:"gid,omitempty"`
}

func marshalManifest(manifest Manifest) ([]byte, error) {
	if err := validateManifest(manifest); err != nil {
		return nil, fmt.Errorf("validate manifest: %w", err)
	}

	manifest.Entries = append([]ManifestEntry(nil), manifest.Entries...)
	sort.Slice(manifest.Entries, func(left, right int) bool {
		return manifest.Entries[left].Path < manifest.Entries[right].Path
	})

	data, err := yamlv3.Marshal(manifest)
	if err != nil {
		return nil, fmt.Errorf("marshal manifest: %w", err)
	}
	return data, nil
}

func parseManifest(data []byte) (Manifest, error) {
	decoder := yamlv3.NewDecoder(bytes.NewReader(data))
	decoder.KnownFields(true)

	var manifest Manifest
	if err := decoder.Decode(&manifest); err != nil {
		return Manifest{}, fmt.Errorf("decode manifest: %w", err)
	}
	var extra any
	if err := decoder.Decode(&extra); !errors.Is(err, io.EOF) {
		if err != nil {
			return Manifest{}, fmt.Errorf("decode manifest: %w", err)
		}
		return Manifest{}, fmt.Errorf("decode manifest: multiple YAML documents are not supported")
	}
	if err := validateManifest(manifest); err != nil {
		return Manifest{}, fmt.Errorf("validate manifest: %w", err)
	}
	return manifest, nil
}

func validateManifest(manifest Manifest) error {
	if manifest.SchemaVersion != manifestSchemaVersion {
		return fmt.Errorf("unsupported schema version %d", manifest.SchemaVersion)
	}

	seen := make(map[string]struct{}, len(manifest.Entries))
	for index, entry := range manifest.Entries {
		if err := validateMemberPath(entry.Path); err != nil {
			return fmt.Errorf("entry %d path %q: %w", index, entry.Path, err)
		}
		if _, ok := seen[entry.Path]; ok {
			return fmt.Errorf("duplicate archive path %q", entry.Path)
		}
		seen[entry.Path] = struct{}{}

		if (entry.UID == nil) != (entry.GID == nil) {
			return fmt.Errorf("entry %q must include both uid and gid or neither", entry.Path)
		}
		if entry.UID != nil && (*entry.UID < 0 || *entry.GID < 0) {
			return fmt.Errorf("entry %q has negative uid or gid", entry.Path)
		}

		switch entry.Type {
		case EntryTypeFile, EntryTypeDirectory:
			if entry.Target != "" || entry.Unsafe || entry.Warning != "" {
				return fmt.Errorf("non-symlink entry %q has symlink metadata", entry.Path)
			}
		case EntryTypeSymlink:
			if entry.Target == "" {
				return fmt.Errorf("symlink entry %q has an empty target", entry.Path)
			}
			unsafe, warning := classifySymlinkTarget(entry.Path, entry.Target)
			if unsafe != entry.Unsafe {
				return fmt.Errorf("symlink entry %q has incorrect unsafe marker", entry.Path)
			}
			if unsafe && (entry.Warning == "" || !strings.Contains(entry.Warning, warning)) {
				return fmt.Errorf("unsafe symlink entry %q has missing or incorrect warning", entry.Path)
			}
			if !unsafe && entry.Warning != "" {
				return fmt.Errorf("safe symlink entry %q has an unsafe warning", entry.Path)
			}
		default:
			return fmt.Errorf("entry %q has unsupported type %q", entry.Path, entry.Type)
		}
	}
	return nil
}

func validateMemberPath(member string) error {
	if member == "" || member == "." || member == ".." {
		return fmt.Errorf("path must be a non-empty relative member path")
	}
	if path.IsAbs(member) {
		return fmt.Errorf("absolute paths are not allowed")
	}
	if cleaned := path.Clean(member); cleaned != member {
		return fmt.Errorf("path is not normalized")
	}
	if strings.HasPrefix(member, "../") {
		return fmt.Errorf("parent traversal is not allowed")
	}
	return nil
}
