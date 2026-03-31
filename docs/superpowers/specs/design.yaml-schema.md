# YAML Schema Validation Design

**Date**: 2026-03-27

## Overview

Add schema validation to the YAML installer config (`installer.yml`) to provide:
1. Clear, typed Go errors at runtime
2. IDE autocomplete and inline validation for config authors

## Current State

The factory currently unmarshals YAML into `map[string]any` and manually type-asserts values.
This provides no compile-time or early validation.

## Proposed Solution

### 1. Define Typed Go Structs

Create a new package `internal/yaml` with structs that mirror the YAML structure:

```go
// PlayBookYAML represents the playbook config structure
type PlayBookYAML struct {
    ID          string     `yaml:"id"`
    Description string     `yaml:"description"`
    Enabled     bool       `yaml:"enabled"`
    Sudo        bool       `yaml:"sudo"`
    JetBrains   []map[string]string `yaml:"jetbrains"`
    Plays       []PlayYAML `yaml:"plays"`
}

type PlayYAML struct {
    ID          string    `yaml:"id"`
    Description string    `yaml:"description"`
    Enabled     bool      `yaml:"enabled"`
    Tasks       []TaskYAML `yaml:"tasks"`
}

type TaskYAML struct {
    Type        string `yaml:"type"`
    ID          string `yaml:"id"`
    Description string `yaml:"description"`
    Enabled     bool   `yaml:"enabled"`
    // Task-specific fields
    Command     string   `yaml:"command,omitempty"`   // bash
    Script      string   `yaml:"script,omitempty"`     // bash
    Root        string   `yaml:"root,omitempty"`       // link, git
    Apps        []string `yaml:"apps,omitempty"`      // jetbrains
}
```

**Fields by Task Type:**
| Task Type   | Fields |
|-------------|--------|
| `bash`      | `command`, `script` |
| `dir`       | (none besides common) |
| `link`      | `root` |
| `git`       | `root` |
| `brew`      | (none besides common) |
| `cask`      | (none besides common) |
| `cellar`    | (none besides common) |
| `function`  | (none besides common) |
| `mas`       | (none besides common) |
| `vscode`    | (none besides common) |
| `jetbrains` | `apps` |

### 2. Use yaml.v3 with Custom Validation

- Replace `map[string]any` with typed structs
- Add `Validate()` methods for business rules:
  - Task type must be valid (use existing `tyg.ParseTaskType`)
  - Either `command` or `script` required for bash task
  - Task ID cannot be empty

### 3. Generate and Distribute JSON Schema

Use [quicktype](https://github.com/glideapps/quicktype) or
[go-jsonschema](https://github.com/alecthomas/go-jsonschema) to generate `installer.schema.json`
from the Go structs.

**Schema distribution:**
1. Generate at build time via `go generate` to `cfg/installer.schema.json`
2. Embed in binary using `//go:embed` so it's available at runtime
3. Include in distro archive so users can extract it for IDE integration

**Runtime extraction:** Add `dum schema` command to extract embedded schema:
```bash
dum schema --output installer.schema.json
```

**Editor integration:** Users point their IDE to the extracted schema:
- **VS Code**: Add to `settings.json`:
  ```json
  "yaml.schemas": {
      "installer.schema.json": "installer.yml"
  }
  ```
- **JetBrains**: Place schema in config directory, IDE auto-detects

## Implementation Steps

1. Create `internal/yaml/` package with struct definitions
2. Add validation methods to structs
3. Update `factory.go` to use new config structs
4. Add `go generate` directive to create JSON schema to `cfg/installer.schema.json`
5. Embed schema in binary using `//go:embed`
6. Add `dum schema` subcommand to extract schema to file
7. Update goreleaser to include schema in distro archive
8. Document editor setup in README

## Trade-offs

| Approach | Pros | Cons |
|----------|------|------|
| Go structs only | Simple, typed errors | No IDE autocomplete |
| Go + JSON Schema | IDE support + typed errors | Extra generation step |
| Schema-only (yamale) | Quick to implement | Separate schema from code |

**Recommendation:** Go structs + JSON Schema (this design). Single source of truth in Go,
IDE support via generated schema.
