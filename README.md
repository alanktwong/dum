# dum

[![Release](https://img.shields.io/github/v/release/alanktwong/dum)](https://github.com/alanktwong/dum/releases)
[![License](https://img.shields.io/github/license/alanktwong/dum)](LICENSE)
[![Build status](https://github.com/alanktwong/dum/actions/workflows/ci.yml/badge.svg)](https://github.com/alanktwong/dum/actions/workflows/ci.yml)
[![Conventional Commits](https://img.shields.io/badge/Conventional%20Commits-1.0.0-yellow.svg)](https://www.conventionalcommits.org)
[![GoDoc](https://pkg.go.dev/badge/github.com/alanktwong/dum.svg)](https://pkg.go.dev/github.com/alanktwong/dum)

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
mise trust && mise install --locked
mise run clean && mise run build
ls dist
```

## Usage

### Subcommands

| Command   | Alias | Description                                     |
|-----------|-------|-------------------------------------------------|
| `install` | `i`   | Runs plays and tasks for software installations |
| `rename`  | `r`   | Renames files with various transformations      |
| `log`     | `lg`  | Logging utilities for shell scripts             |
| `schema`  | `s`   | Prints or extracts the embedded JSON schema for `dum.yml` |

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

### CLI surfaces

Every `main.go` binary in the repository:

| Path             | Purpose                                                                                                                     |
|------------------|-----------------------------------------------------------------------------------------------------------------------------|
| `cmd/dum`        | The `dum` binary: `install` (`i`), `rename` (`r`), `log` (`lg`) with levels `debug d`, `info i`, `success s`, `warn w`, `error e`, and `schema` (`s`) for the embedded JSON schema |
| `cmd/gen-schema` | Build-time tool run via `go:generate`; regenerates `cmd/dum/installer.schema.json`. Not user-facing                          |
| `cmd/linters`    | Standalone `go/analysis` linter (own go.mod) flagging single-character import aliases; wired into golangci-lint as a plugin  |

### Build automation

The repository's shell recipes live in `tools/build.sh` and `tools/lib/*.sh`. Use mise tasks for
normal development; they load the locked tools and environment before calling the dispatcher.
`tools/build.sh` remains available for direct dispatcher access and must be run from the repository
root.

```shell
mise tasks ls                 # List supported tasks
mise run help                 # Same task list
mise run build                # Run normal build task
mise run test                 # Run tests with coverage
./tools/build.sh help         # List direct dispatcher subcommands
./tools/build.sh check-lfs-hook  # Run direct hook check
```

The former standalone tool scripts are dispatcher subcommands. Use
`./tools/build.sh check-lfs-hook`, `./tools/build.sh coverage-check`, or
`./tools/build.sh cpd` when those direct checks are needed.

### Developer setup

Install mise once using a supported method from the
[mise installation guide](https://mise.jdx.dev/installing-mise.html); Homebrew is optional and is
not a project dev-tool dependency. In the repository, trust the config and install locked tools:

```shell
mise trust && mise install --locked
```

For project tools and `.env` to activate automatically on directory changes—including for Git
hooks—add the appropriate activation command to your shell startup file:

```shell
eval "$(mise activate bash)"  # bash
eval "$(mise activate zsh)"   # zsh
mise activate fish | source   # fish
```

Remove any old `direnv hook` line from your shell startup file. The repository no longer has
`.envrc`, so a stale direnv hook silently stops setting `__PROJECT_DIRENV_INSTALLED` and loading
Git LFS/pre-commit setup and optional `.env`. Explicit `mise run ...` tasks work without shell
activation.

Run repository hook and Git LFS setup once after installing tools:

```shell
mise run install
```

### Build Commands

```shell
mise run build           # Local build
mise run test            # Run tests with coverage
mise run lint            # Run golangci-lint
mise run vet             # Go vet
mise run fmt             # Go fmt
mise run fix             # Run go fix to update code to new language features
mise run coverage        # Generate the HTML coverage report
mise run check-coverage  # Check coverage meets minimum threshold
mise run all             # Run build and quality checks
mise run clean           # Clean generated files
mise run tag             # Compute next semver from conventional commits and create git tag
```

### Versioning and Releases

Contributions follow [Conventional Commits](https://www.conventionalcommits.org); commitlint
validates local commit messages via pre-commit. Node is managed by mise.

To cut a release: run `mise run tag` ([svu](https://github.com/caarlos0/svu) computes the next semver
from merged conventional commits), then `git push origin <tag>`. Pushing the tag triggers
GoReleaser, which publishes per-platform archives and a SHA256 checksums file to the GitHub
Release.

### Enums

We are using [go-enum](https://github.com/abice/go-enum) to generate enums from Go struct field tags.
The generated enums, by convention, are suffixed as `_enum.go` and are gitignored.

### Mockery

We are using [mockery](https://vektra.github.io/mockery/v3.0/) to generate
mocks of Go interfaces for testing.
The mockery configurations live in `cfg/mockery.*.yml` and the generated
mocks are gitignored since they should follow the naming convention `*_mocks.go`.

### JSON Schema

The JSON schema for `dum.yml` is generated by `cmd/gen-schema` via a
`go:generate` directive. Run `mise run generate` to regenerate
`cmd/dum/installer.schema.json`, which is embedded into the binary; commit the
updated file when task types or config fields change.

### Coverage Threshold

Default minimum coverage is **80%**. Override via environment variable:

```shell
COVERAGE_THRESHOLD=90 mise run check-coverage
```

### Linting Principles

This project enforces strict import alias standards and enables comprehensive linting:

1. **No single-letter import aliases** - Aliases like `l`, `i`, `t` are banned. Use descriptive 2+ character aliases.
2. **Standardized aliases** - All internal and vendor packages must use defined aliases:
   - Internal: `cd`, `ext`, `fy`, `lg`, `pb`, `pl`, `plg`, `tk`, `ti`, `ty`, `tyg`
   - Vendor: `clog`, `omv3`, `ca`, `tt`, `asrt`, `mck`, `yamlv3`
3. **Comprehensive linting** - 34 linters enabled via golangci-lint covering correctness, security, performance, and style.

Run `mise run lint` to lint. Use `golangci-lint run ./... --fix` to auto-fix import alias issues.
