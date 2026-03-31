# YAML Schema Validation Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Add schema validation to the YAML installer config with typed Go structs and JSON Schema generation for IDE autocomplete support.

**Architecture:** Create a new `internal/yaml/` package with typed structs that mirror the YAML config structure. Replace `map[string]any` unmarshaling with typed structs. Generate JSON Schema at build time and embed it in the binary. Add a `dum schema` subcommand for schema extraction.

**Tech Stack:** Go 1.26, yaml.v3, quicktype or go-jsonschema for schema generation, //go:embed for binary embedding

---

## File Structure

The new `internal/yaml/` package will contain:
- `internal/yaml/config.go` - Typed struct definitions for PlayBookYAML, PlayYAML, TaskYAML
- `internal/yaml/config_test.go` - Tests for struct validation
- `internal/yaml/loader.go` - Functions to load and validate YAML with the typed structs

Existing files to modify:
- `internal/factory/factory.go` - Update to use new yaml package
- `internal/factory/factory_test.go` - May need updates
- `internal/cmd/dum.go` - Add schema subcommand registration
- `internal/cmd/schema.go` - New file for schema subcommand
- `Makefile` - Add schema generation to generate target
- `cfg/goreleaser.yaml` - Add schema to archives

Generated files:
- `cfg/installer.schema.json` - Generated JSON Schema (gitignored)

---

## Task 1: Create YAML Config Package with Typed Structs

**Files:**
- Create: `internal/yaml/config.go`
- Test: `internal/yaml/config_test.go`

- [ ] **Step 1: Create config.go with typed structs**

```go
package yaml

import (
    "fmt"

    tyg "awong/dotfiles/internal/types/gen"
)

type PlayBookYAML struct {
    ID          string                `yaml:"id"`
    Description string                `yaml:"description"`
    Enabled     bool                  `yaml:"enabled"`
    Sudo        bool                  `yaml:"sudo"`
    JetBrains   []map[string]string   `yaml:"jetbrains"`
    Plays       []PlayYAML            `yaml:"plays"`
}

type PlayYAML struct {
    ID          string    `yaml:"id"`
    Description string    `yaml:"description"`
    Enabled     bool      `yaml:"enabled"`
    Tasks       []TaskYAML `yaml:"tasks"`
}

type TaskYAML struct {
    Type        string   `yaml:"type"`
    ID          string   `yaml:"id"`
    Description string   `yaml:"description"`
    Enabled     bool     `yaml:"enabled"`
    // Task-specific fields
    Command     string   `yaml:"command,omitempty"`
    Script      string   `yaml:"script,omitempty"`
    Root        string   `yaml:"root,omitempty"`
    Target      string   `yaml:"target,omitempty"`
    Name        string   `yaml:"name,omitempty"`
    Tap         string   `yaml:"tap,omitempty"`
    Apps        []string `yaml:"apps,omitempty"`
}

func (p *PlayBookYAML) Validate() error {
    if p.ID == "" {
        return fmt.Errorf("playbook id cannot be empty")
    }
    for _, play := range p.Plays {
        if err := play.Validate(); err != nil {
            return fmt.Errorf("play %s: %w", play.ID, err)
        }
    }
    return nil
}

func (p *PlayYAML) Validate() error {
    if p.ID == "" {
        return fmt.Errorf("play id cannot be empty")
    }
    for _, task := range p.Tasks {
        if err := task.Validate(); err != nil {
            return fmt.Errorf("task %s: %w", task.ID, err)
        }
    }
    return nil
}

func (t *TaskYAML) Validate() error {
    if t.ID == "" {
        return fmt.Errorf("task id cannot be empty")
    }
    if t.Type == "" {
        return fmt.Errorf("task type cannot be empty")
    }
    if !tyg.TaskType(t.Type).IsValid() {
        return fmt.Errorf("invalid task type: %s", t.Type)
    }
    if t.Type == string(tyg.TaskTypeBash) && t.Command == "" && t.Script == "" {
        return fmt.Errorf("bash task must have either command or script")
    }
    return nil
}
```

- [ ] **Step 2: Create config_test.go with validation tests**

```go
package yaml

import (
    "testing"

    "github.com/stretchr/testify/assert"
)

func TestPlayBookYAML_Validate(t *testing.T) {
    tests := []struct {
        name    string
        pb      PlayBookYAML
        wantErr bool
    }{
        {
            name: "valid playbook",
            pb: PlayBookYAML{
                ID:   "test",
                Plays: []PlayYAML{{ID: "play-1", Tasks: []TaskYAML{{ID: "t1", Type: "dir"}}}},
            },
            wantErr: false,
        },
        {
            name: "empty id",
            pb:   PlayBookYAML{},
            wantErr: true,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            err := tt.pb.Validate()
            if tt.wantErr {
                assert.Error(t, err)
            } else {
                assert.NoError(t, err)
            }
        })
    }
}

func TestTaskYAML_Validate(t *testing.T) {
    tests := []struct {
        name    string
        task    TaskYAML
        wantErr bool
        errMsg  string
    }{
        {
            name:    "valid dir task",
            task:    TaskYAML{ID: "t1", Type: "dir"},
            wantErr: false,
        },
        {
            name:    "empty id",
            task:    TaskYAML{Type: "dir"},
            wantErr: true,
            errMsg:  "task id cannot be empty",
        },
        {
            name:    "empty type",
            task:    TaskYAML{ID: "t1"},
            wantErr: true,
            errMsg:  "task type cannot be empty",
        },
        {
            name:    "invalid type",
            task:    TaskYAML{ID: "t1", Type: "unknown"},
            wantErr: true,
            errMsg:  "invalid task type",
        },
        {
            name:    "bash with command",
            task:    TaskYAML{ID: "t1", Type: "bash", Command: "echo hello"},
            wantErr: false,
        },
        {
            name:    "bash with script",
            task:    TaskYAML{ID: "t1", Type: "bash", Script: "echo hello"},
            wantErr: false,
        },
        {
            name:    "bash without command or script",
            task:    TaskYAML{ID: "t1", Type: "bash"},
            wantErr: true,
            errMsg:  "bash task must have either command or script",
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            err := tt.task.Validate()
            if tt.wantErr {
                assert.Error(t, err)
                if tt.errMsg != "" {
                    assert.Contains(t, err.Error(), tt.errMsg)
                }
            } else {
                assert.NoError(t, err)
            }
        })
    }
}
```

- [ ] **Step 3: Run tests to verify they pass**

Run: `go test -v ./internal/yaml/...`
Expected: PASS (all tests should pass)

- [ ] **Step 4: Commit**

```bash
git add internal/yaml/config.go internal/yaml/config_test.go
git commit -m "feat: add typed YAML config structs with validation"
```

---

## Task 2: Create YAML Loader with Unmarshal Functions

**Files:**
- Create: `internal/yaml/loader.go`
- Test: `internal/yaml/loader_test.go`

- [ ] **Step 1: Create loader.go with unmarshal functions**

```go
package yaml

import (
    "fmt"
    "os"

    "gopkg.in/yaml.v3"
)

type Config struct {
    PlayBook PlayBookYAML `yaml:"playbook"`
}

func Load(filePath string) (*Config, error) {
    data, err := os.ReadFile(filePath)
    if err != nil {
        return nil, fmt.Errorf("failed to read file: %w", err)
    }

    var config Config
    if err := yaml.Unmarshal(data, &config); err != nil {
        return nil, fmt.Errorf("failed to unmarshal YAML: %w", err)
    }

    if err := config.PlayBook.Validate(); err != nil {
        return nil, fmt.Errorf("validation failed: %w", err)
    }

    return &config, nil
}
```

- [ ] **Step 2: Create loader_test.go**

```go
package yaml

import (
    "os"
    "path/filepath"
    "testing"

    "github.com/stretchr/testify/assert"
)

func TestLoad_Success(t *testing.T) {
    config, err := Load("../../factory/testdata/test_installer.yml")
    assert.NoError(t, err)
    assert.NotNil(t, config)
    assert.Equal(t, "my-test-playbook", config.PlayBook.ID)
    assert.Len(t, config.PlayBook.Plays, 2)
}

func TestLoad_FileNotFound(t *testing.T) {
    _, err := Load("/nonexistent/path/installer.yml")
    assert.Error(t, err)
    assert.Contains(t, err.Error(), "failed to read file")
}

func TestLoad_InvalidYaml(t *testing.T) {
    tmpDir := t.TempDir()
    badFile := filepath.Join(tmpDir, "bad.yml")
    err := os.WriteFile(badFile, []byte("{{invalid yaml: [}"), 0644)
    assert.NoError(t, err)

    _, err = Load(badFile)
    assert.Error(t, err)
    assert.Contains(t, err.Error(), "failed to unmarshal YAML")
}

func TestLoad_ValidationError(t *testing.T) {
    tmpDir := t.TempDir()
    badFile := filepath.Join(tmpDir, "invalid.yml")
    err := os.WriteFile(badFile, []byte("playbook:\n  plays:\n    - id: \"\"\n"), 0644)
    assert.NoError(t, err)

    _, err = Load(badFile)
    assert.Error(t, err)
    assert.Contains(t, err.Error(), "validation failed")
}
```

- [ ] **Step 3: Run tests to verify they pass**

Run: `go test -v ./internal/yaml/...`
Expected: PASS

- [ ] **Step 4: Commit**

```bash
git add internal/yaml/loader.go internal/yaml/loader_test.go
git commit -m "feat: add YAML loader with validation"
```

---

## Task 3: Integrate New YAML Package with Factory

**Files:**
- Modify: `internal/factory/factory.go`
- Test: Existing tests should still pass

- [ ] **Step 1: Add import and helper function to factory.go**

Add to imports:
```go
yml "awong/dotfiles/internal/yaml"
```

Add new function:
```go
func (f *Factory) loadFromTypedYAML(file string) (*yml.Config, error) {
    return yml.Load(file)
}
```

- [ ] **Step 2: Run existing tests to ensure they pass**

Run: `go test -v ./internal/factory/...`
Expected: PASS (existing tests should still work)

- [ ] **Step 3: Commit**

```bash
git add internal/factory/factory.go
git commit -m "feat: integrate typed YAML config with factory"
```

---

## Task 4: Generate JSON Schema

**Files:**
- Modify: `internal/yaml/config.go` (add go:generate directive)
- Create: `cfg/installer.schema.json` (generated, gitignored)

- [ ] **Step 1: Add go:generate directive to config.go**

Add at top of config.go:
```go
//go:generate go run -tags gen_schema ./gen_schema.go

// ConfigForSchema is used by the schema generator to get the structure
type ConfigForSchema struct {
    PlayBook PlayBookYAML `yaml:"playbook"`
}
```

- [ ] **Step 2: Create schema generation tool**

Create: `internal/yaml/gen_schema/main.go`

```go
//go:build gen_schema

package main

import (
    "encoding/json"
    "fmt"
    "os"

    "gopkg.in/yaml.v3"
)

type PlayBookYAML struct {
    ID          string                `yaml:"id"`
    Description string                `yaml:"description"`
    Enabled     bool                  `yaml:"enabled"`
    Sudo        bool                  `yaml:"sudo"`
    JetBrains   []map[string]string   `yaml:"jetbrains"`
    Plays       []PlayYAML            `yaml:"plays"`
}

type PlayYAML struct {
    ID          string    `yaml:"id"`
    Description string    `yaml:"description"`
    Enabled     bool      `yaml:"enabled"`
    Tasks       []TaskYAML `yaml:"tasks"`
}

type TaskYAML struct {
    Type        string   `yaml:"type"`
    ID          string   `yaml:"id"`
    Description string   `yaml:"description"`
    Enabled     bool     `yaml:"enabled"`
    Command     string   `yaml:"command,omitempty"`
    Script      string   `yaml:"script,omitempty"`
    Root        string   `yaml:"root,omitempty"`
    Target      string   `yaml:"target,omitempty"`
    Name        string   `yaml:"name,omitempty"`
    Tap         string   `yaml:"tap,omitempty"`
    Apps        []string `yaml:"apps,omitempty"`
}

type Config struct {
    PlayBook PlayBookYAML `yaml:"playbook"`
}

func main() {
    // Use yaml.Tags to generate JSON schema
    // This is a simplified version - for production, use quicktype or go-jsonschema
    
    // Read the struct and generate schema manually for now
    schema := map[string]interface{}{
        "$schema":     "http://json-schema.org/draft-07/schema#",
        "title":       "Dum Installer Config",
        "type":        "object",
        "properties": map[string]interface{}{
            "playbook": map[string]interface{}{
                "type": "object",
                "properties": map[string]interface{}{
                    "id":          map[string]interface{}{"type": "string"},
                    "description": map[string]interface{}{"type": "string"},
                    "enabled":     map[string]interface{}{"type": "boolean"},
                    "sudo":        map[string]interface{}{"type": "boolean"},
                    "jetbrains": map[string]interface{}{
                        "type": "array",
                        "items": map[string]interface{}{
                            "type": "object",
                            "additionalProperties": map[string]interface{}{"type": "string"},
                        },
                    },
                    "plays": map[string]interface{}{
                        "type": "array",
                        "items": map[string]interface{}{
                            "type": "object",
                            "properties": map[string]interface{}{
                                "id":          map[string]interface{}{"type": "string"},
                                "description": map[string]interface{}{"type": "string"},
                                "enabled":     map[string]interface{}{"type": "boolean"},
                                "tasks": map[string]interface{}{
                                    "type": "array",
                                    "items": map[string]interface{}{
                                        "type": "object",
                                        "properties": map[string]interface{}{
                                            "type":        map[string]interface{}{"type": "string", "enum": []string{"bash", "dir", "link", "git", "brew", "cask", "cellar", "function", "mas", "vscode", "jetbrains"}},
                                            "id":          map[string]interface{}{"type": "string"},
                                            "description": map[string]interface{}{"type": "string"},
                                            "enabled":     map[string]interface{}{"type": "boolean"},
                                            "command":     map[string]interface{}{"type": "string"},
                                            "script":      map[string]interface{}{"type": "string"},
                                            "root":        map[string]interface{}{"type": "string"},
                                            "target":      map[string]interface{}{"type": "string"},
                                            "name":        map[string]interface{}{"type": "string"},
                                            "tap":         map[string]interface{}{"type": "string"},
                                            "apps": map[string]interface{}{
                                                "type":  "array",
                                                "items": map[string]interface{}{"type": "string"},
                                            },
                                        },
                                        "required": []string{"type", "id"},
                                    },
                                },
                            },
                            "required": []string{"id", "tasks"},
                        },
                    },
                },
                "required": []string{"id", "plays"},
            },
        },
        "required": []string{"playbook"},
    }

    data, err := json.MarshalIndent(schema, "", "  ")
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error marshaling schema: %v\n", err)
        os.Exit(1)
    }

    outputPath := "cfg/installer.schema.json"
    if len(os.Args) > 1 {
        outputPath = os.Args[1]
    }

    err = os.WriteFile(outputPath, data, 0644)
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error writing schema: %v\n", err)
        os.Exit(1)
    }

    fmt.Printf("Generated JSON schema at %s\n", outputPath)
}
```

- [ ] **Step 3: Run schema generation**

Run: `go generate ./internal/yaml/...`
Expected: Generates `cfg/installer.schema.json`

- [ ] **Step 4: Verify schema file was created**

Run: `ls -la cfg/installer.schema.json`
Expected: File exists

- [ ] **Step 5: Add to .gitignore**

Ensure `cfg/installer.schema.json` is gitignored (check .gitignore)

- [ ] **Step 6: Commit**

```bash
git add internal/yaml/gen_schema/ cfg/installer.schema.json .gitignore 2>/dev/null || true
git commit -m "feat: add JSON schema generation for installer config"
```

---

## Task 5: Embed Schema in Binary

**Files:**
- Create: `internal/cmd/schema_embed.go` (or add to existing file)
- Modify: `internal/cmd/dum.go` (register command)

- [ ] **Step 1: Create schema_embed.go**

```go
package cmd

import (
    "embed"
    "os"

    "github.com/spf13/cobra"
)

var schemaFS embed.FS

var schemaCmd = &cobra.Command{
    Use:   "schema",
    Short: "Output or extract the JSON schema for installer.yml",
    Long: `Outputs the JSON schema for installer.yml configuration file.

This can be used for IDE autocomplete and validation in editors like
VS Code and JetBrains IDEs.

Examples:
  # Print schema to stdout
  dum schema

  # Save schema to a file
  dum schema --output installer.schema.json

  # Save to a specific location
  dum schema -o ~/.config/dum/installer.schema.json
`,
    RunE: func(cmd *cobra.Command, args []string) error {
        outputPath, err := cmd.Flags().GetString("output")
        if err != nil {
            return err
        }

        schemaData, err := schemaFS.ReadFile("installer.schema.json")
        if err != nil {
            return err
        }

        if outputPath == "" {
            // Print to stdout
            os.Stdout.Write(schemaData)
            return nil
        }

        // Write to file
        return os.WriteFile(outputPath, schemaData, 0644)
    },
}

func init() {
    schemaCmd.Flags().StringP("output", "o", "", "Output file path (default: stdout)")
}
```

- [ ] **Step 2: Add embed directive**

Add to `internal/cmd/dum.go` or create a new file:
```go
//go:embed installer.schema.json
var schemaFS embed.FS
```

- [ ] **Step 3: Add command to dum.go**

Add to `NewDum()`:
```go
schemaCmd := NewSchemaCommand(rootUse, dum)
rootCmd.AddCommand(schemaCmd)
```

- [ ] **Step 4: Create schema command constructor**

Create: `internal/cmd/schema.go`

```go
package cmd

func NewSchemaCommand(rootUse string, dum *Dum) *cobra.Command {
    use := "schema"
    cmd := &cobra.Command{
        Use:   use,
        Short: "Output or extract the JSON schema for installer.yml",
        Long: `Outputs the JSON schema for installer.yml configuration file.

This can be used for IDE autocomplete and validation in editors like
VS Code and JetBrains IDEs.

Examples:
  # Print schema to stdout
  dum schema

  # Save schema to a file
  dum schema --output installer.schema.json
`,
        RunE: func(cmd *cobra.Command, args []string) error {
            outputPath, err := cmd.Flags().GetString("output")
            if err != nil {
                return err
            }

            schemaData, err := schemaFS.ReadFile("installer.schema.json")
            if err != nil {
                return err
            }

            if outputPath == "" {
                // Print to stdout
                os.Stdout.Write(schemaData)
                return nil
            }

            // Write to file
            return os.WriteFile(outputPath, schemaData, 0644)
        },
    }

    cmd.Flags().StringP("output", "o", "", "Output file path (default: stdout)")
    return cmd
}
```

- [ ] **Step 5: Run build to verify it compiles**

Run: `go build -v ./cmd/dum/...`
Expected: SUCCESS

- [ ] **Step 6: Test the schema command**

Run: `./dist/dum-* schema`
Expected: Outputs JSON schema

- [ ] **Step 7: Commit**

```bash
git add internal/cmd/schema.go
git commit -m "feat: add dum schema command to extract JSON schema"
```

---

## Task 6: Update Build and Release Process

**Files:**
- Modify: `Makefile`
- Modify: `cfg/goreleaser.yaml`

- [ ] **Step 1: Add schema generation to Makefile generate target**

The `generate` target already runs `go generate ./...`, so schema should be generated automatically.

- [ ] **Step 2: Add schema to goreleaser archive**

Modify `cfg/goreleaser.yaml` to include the schema in the archive:
```yaml
archives:
  - formats: ["tar.gz"]
    # ... existing config
    files:
      - cfg/installer.schema.json
```

- [ ] **Step 3: Commit**

```bash
git add Makefile cfg/goreleaser.yaml
git commit -m "build: include JSON schema in release archives"
```

---

## Task 7: Documentation

**Files:**
- Modify: `README.md` (add editor setup instructions)

- [ ] **Step 1: Add editor integration docs to README**

```markdown
## Editor Integration

For IDE autocomplete and validation when editing `installer.yml`, you can use the JSON schema:

### VS Code

Add to your `settings.json`:

```json
{
  "yaml.schemas": {
    "installer.schema.json": "installer.yml"
  }
}
```

Extract the schema:
```bash
dum schema --output installer.schema.json
```

### JetBrains IDEs

Place the schema in your IDE config directory:
- macOS: `~/Library/Application Support/JetBrains/<IDE>/schemas`
- Linux: `~/.config/JetBrains/<IDE>/schemas`

The IDE should automatically detect and use the schema for `installer.yml` files.

### Other Editors

Extract the schema and configure your editor:
```bash
dum schema --output installer.schema.json
```

Most YAML-aware editors support JSON Schema validation.
```

- [ ] **Step 2: Commit**

```bash
git commit -m "docs: add editor integration instructions for YAML schema"
```

---

## Execution Commands Summary

After all tasks, run these to verify:

```bash
# Build
make build

# Run tests
go test -v ./internal/yaml/... ./internal/factory/...

# Test schema command
./dist/dum-* schema --output /tmp/test.schema.json
cat /tmp/test.schema.json | head -20

# Run linter
make lint

# Run full check
make check
```

---

## Plan Complete

**Two execution options:**

**1. Subagent-Driven (recommended)** - I dispatch a fresh subagent per task, review between tasks, fast iteration

**2. Inline Execution** - Execute tasks in this session using executing-plans, batch execution with checkpoints

Which approach?
