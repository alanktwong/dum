# dum

A CLI tool for managing software installations and configurations on macOS/Linux laptops.

## Overview

Dum is a hand-rolled CLI tool similar to Ansible, written in Go. It manages software installations and configurations using a playbook-based system defined in YAML.

## Installation

Requires [Go](https://go.dev) 1.26 or newer:

```shell
go install github.com/alanktwong/dum/cmd/dum@latest
```

Or build from source:

```shell
git clone https://github.com/alanktwong/dum
cd dum
make clean build
ls dist
```

## Usage

### Subcommands

| Command   | Alias | Description                                     |
|-----------|-------|-------------------------------------------------|
| `install` | `i`   | Runs plays and tasks for software installations |
| `rename`  | `r`   | Renames files with various transformations      |
| `log`     | `lg`  | Logging utilities for shell scripts             |

### Install Command

Installs software based on the playbook configuration:

```shell
dum install                    # Install all plays
dum install --group <name>    # Install specific play group
dum install --dryrun       # Preview without installing (includes disabled plays and tasks)
dum install -f <config.yml>   # Use custom config file
```

Use `--dryrun` to preview what would run without making changes. Dry run output includes disabled plays and tasks, showing the full playbook manifest.

### Rename Command

Renames files with various transformations:

```shell
dum rename -s <source> -r <replace> <files...>  # Replace text in filenames
dum rename -l <files...>                         # Convert to lowercase
dum rename -f <n> <files...>                   # Trim n chars from front
dum rename -b <n> <files...>                   # Trim n chars from back
dum rename -i <n> <files...>                   # Add numeric suffix starting at n
dum rename --dryrun <files...>                # Preview without renaming
```

### Log Command

Logging utilities for shell scripts with subcommands:

```shell
dum log debug <message>    # Debug level log
dum log info <message>     # Info level log
dum log success <message>  # Success level log
dum log warn <message>     # Warn level log
dum log error <message>    # Error level log
```

## Configuration

By default, the playbook is resolved in this order (first match wins):

1. `DUM_CONFIG` environment variable
2. `./dum.yml` in the current directory
3. `$XDG_CONFIG_HOME/dum/dum.yml` (defaults to `~/.config/dum/dum.yml`)
4. Legacy fallback: `$XDG_CONFIG_HOME/dum/installer.yml`

You can also specify a custom file with the `-f` flag.

## Task Types

- `brew` / `cellar` / `cask` - Homebrew installations
- `dir` - Create directories
- `git` - Git clone operations
- `link` - Create symbolic links
- `vscode` - VS Code extensions
- `mas` - Mac App Store apps
- `jetbrains` - JetBrains IDE plugins

## Editor Integration

For IDE autocomplete and validation when editing `dum.yml`, you can use the JSON schema:

### VS Code

Add to your `settings.json`:

```json
{
  "yaml.schemas": {
    "installer.schema.json": "dum.yml"
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

The IDE should automatically detect and use the schema for `dum.yml` files.

### Other Editors

Extract the schema and configure your editor:
```bash
dum schema --output installer.schema.json
```

Most YAML-aware editors support JSON Schema validation.

## Similar Tools

Dum sits in the space between Ansible (full configuration management) and shell-based bootstrap scripts (one-shot setup). Here are related projects:

| Tool | Language | Approach | Platform |
|------|----------|----------|----------|
| [Preflight](https://github.com/felixgeelhaar/preflight) | Go | Declarative YAML, compiler model (`plan → apply`), multi-platform package management | macOS, Linux |
| [Omakub](https://github.com/basecamp/omakub) | Shell | Opinionated Ubuntu setup via `wget \| bash`, curated tool selections | Ubuntu |
| [Omarchy](https://omarchy.org) | Shell | Ambitious take on Linux setup (companion to Omakub) | Linux |
| [dotfiles-cli](https://github.com/wsoule/dotfiles-cli) | Go | Cross-platform package management (Homebrew/apt/pacman/yum), GNU Stow integration, JSON config | macOS, Linux |
| [Chezmoi](https://www.chezmoi.io/) | Go | Declarative dotfile manager with templating, encryption, and host-specific configs | macOS, Linux, Windows |
| [Ansible](https://github.com/ansible/ansible) | Python | Full configuration management with YAML playbooks, idempotent tasks, and a large module ecosystem | macOS, Linux |

**What makes Dum different:** Dum combines Go performance with an Ansible-like playbook/task model specifically tailored for personal workstation setup. It offers JSON schema validation for config files, a `--dryrun` preview mode, and a growing set of opinionated task types (brew, git, vscode, jetbrains, mas) without the overhead of a full configuration management tool.

## Development

### Build automation

The repository's shell automation has one entry point: `tools/build.sh`. It must be run from the repository root. The Makefile targets are the supported wrappers for normal development; they pass the output directory, version, and coverage settings to the dispatcher while preserving the commands below.

```shell
make help                 # List the Makefile wrappers
./tools/build.sh help     # List dispatcher subcommands
./tools/build.sh make_out # Create the output directory
./tools/build.sh build    # Invoke the build recipe directly
./tools/build.sh test     # Invoke the test recipe directly
```

The former standalone tool scripts are now dispatcher subcommands. Use `./tools/build.sh check-lfs-hook`, `./tools/build.sh coverage-check`, or `./tools/build.sh cpd` when those direct checks are needed. Run the dispatcher from the repository root; Makefile wrappers handle that automatically.

### Build Commands

```shell
make build           # Local build
make test            # Run tests with coverage
make lint            # Run golangci-lint
make vet             # Go vet
make fmt             # Go fmt
make fix             # Runs go fix to update code to new language features
make coverage        # View coverage report in browser and check threshold
make check-coverage  # Check coverage meets minimum threshold
make all             # Runs build and quality checks
make clean           # Clean generated files
```

### Enums

We are using [go-enum](https://github.com/abice/go-enum) to generate enums from Go struct field tags.
The generated enums, by convention, are suffixed as `_enum.go` and are gitignored.

### Mockery

We are using [mockery](https://vektra.github.io/mockery/v3.0/) to generate
mocks of Go interfaces for testing.
The mockery configurations live in `cfg/mockery.*.yml` and the generated
mocks are gitignored since they should follow the naming convention `*_mocks.go`.

### Coverage Threshold

Default minimum coverage is **80%**. Override via environment variable:

```shell
COVERAGE_THRESHOLD=90 make check-coverage
```

### Linting Principles

This project enforces strict import alias standards and enables comprehensive linting:

1. **No single-letter import aliases** - Aliases like `l`, `i`, `t` are banned. Use descriptive 2+ character aliases.
2. **Standardized aliases** - All internal and vendor packages must use defined aliases:
   - Internal: `cd`, `ext`, `fy`, `lg`, `pb`, `pl`, `plg`, `tk`, `ti`, `ty`, `tyg`
   - Vendor: `clog`, `omv3`, `ca`, `tt`, `asrt`, `mck`, `yamlv3`
3. **Comprehensive linting** - 34 linters enabled via golangci-lint covering correctness, security, performance, and style.

Run `make lint` to lint. Use `golangci-lint run ./... --fix` to auto-fix import alias issues.
