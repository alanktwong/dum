package tarball

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestManifest_MarshalAndParse(t *testing.T) {
	uid := int64(501)
	gid := int64(20)
	manifest := Manifest{
		SchemaVersion: 1,
		Entries: []ManifestEntry{
			{Path: "Documents/config-link", Type: EntryTypeSymlink, Target: "settings"},
			{Path: "Documents/project/readme.txt", Type: EntryTypeFile, UID: &uid, GID: &gid},
			{Path: "Documents", Type: EntryTypeDirectory},
		},
	}

	data, err := marshalManifest(manifest)
	if !assert.NoError(t, err) {
		return
	}
	assert.Contains(t, string(data), "schema_version: 1")
	assert.Contains(t, string(data), "use_gitignore: false")
	assert.Less(
		t,
		strings.Index(string(data), "path: Documents\n"),
		strings.Index(string(data), "path: Documents/config-link"),
	)

	got, err := parseManifest(data)
	if !assert.NoError(t, err) {
		return
	}
	assert.Equal(t, Manifest{
		SchemaVersion: 1,
		Entries: []ManifestEntry{
			{Path: "Documents", Type: EntryTypeDirectory},
			{Path: "Documents/config-link", Type: EntryTypeSymlink, Target: "settings"},
			{Path: "Documents/project/readme.txt", Type: EntryTypeFile, UID: &uid, GID: &gid},
		},
	}, got)
}

func TestValidateManifest_RejectsInvalidEntries(t *testing.T) {
	uid := int64(501)
	tests := []struct {
		name    string
		entries []ManifestEntry
	}{
		{
			name: "absolute path",
			entries: []ManifestEntry{
				{Path: "/etc/passwd", Type: EntryTypeFile},
			},
		},
		{
			name: "parent traversal",
			entries: []ManifestEntry{
				{Path: "../escape", Type: EntryTypeFile},
			},
		},
		{
			name: "non-normalized path",
			entries: []ManifestEntry{
				{Path: "Documents/./file", Type: EntryTypeFile},
			},
		},
		{
			name: "unsupported type",
			entries: []ManifestEntry{
				{Path: "Documents/file", Type: EntryType("device")},
			},
		},
		{
			name: "incomplete ownership pair",
			entries: []ManifestEntry{
				{Path: "Documents/file", Type: EntryTypeFile, UID: &uid},
			},
		},
		{
			name: "duplicate member path",
			entries: []ManifestEntry{
				{Path: "Documents/file", Type: EntryTypeFile},
				{Path: "Documents/file", Type: EntryTypeFile},
			},
		},
		{
			name: "symlink missing target",
			entries: []ManifestEntry{
				{Path: "Documents/link", Type: EntryTypeSymlink},
			},
		},
		{
			name: "unsafe symlink missing warning",
			entries: []ManifestEntry{
				{Path: "Documents/link", Type: EntryTypeSymlink, Target: "/etc/passwd", Unsafe: true},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			assert.Error(t, validateManifest(Manifest{SchemaVersion: 1, Entries: test.entries}))
		})
	}
}

func TestValidateManifest_RejectsUnsupportedSchemaVersion(t *testing.T) {
	manifest := Manifest{SchemaVersion: 2, Entries: []ManifestEntry{{Path: "file", Type: EntryTypeFile}}}
	assert.Error(t, validateManifest(manifest))
}

func TestParseManifest_RejectsUnknownFields(t *testing.T) {
	data := []byte("schema_version: 1\nentries: []\nextra: true\n")
	_, err := parseManifest(data)
	assert.Error(t, err)
}

func TestValidateManifest_RejectsFalselySafeSymlink(t *testing.T) {
	manifest := Manifest{
		SchemaVersion: 1,
		Entries: []ManifestEntry{{
			Path:   "Documents/link",
			Type:   EntryTypeSymlink,
			Target: "../../outside",
		}},
	}
	assert.Error(t, validateManifest(manifest))
}
