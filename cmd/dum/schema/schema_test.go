package schema

import (
	"context"
	"io"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewSchemaCommand_Stdout(t *testing.T) {
	dum := &Command{}
	cmd := NewSchemaCommand("dum", dum)
	cmd.SetContext(context.Background())

	oldStdout := os.Stdout
	defer func() { os.Stdout = oldStdout }()

	pr, pw, _ := os.Pipe()
	os.Stdout = pw

	err := cmd.RunE(cmd, []string{})

	pw.Close()
	out, _ := io.ReadAll(pr)

	os.Stdout = oldStdout

	assert.NoError(t, err)
	assert.NotZero(t, len(out))
	assert.Equal(t, schemaData, out)
}

func TestNewSchemaCommand_File(t *testing.T) {
	dum := &Command{}
	cmd := NewSchemaCommand("dum", dum)
	cmd.SetContext(context.Background())

	tmpDir := t.TempDir()
	outputPath := filepath.Join(tmpDir, "test.schema.json")

	err := cmd.Flags().Set("output", outputPath)
	assert.NoError(t, err)

	err = cmd.RunE(cmd, []string{})

	assert.NoError(t, err)

	content, err := os.ReadFile(outputPath)
	assert.NoError(t, err)
	assert.Equal(t, schemaData, content)
}
