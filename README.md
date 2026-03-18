# dum

A CLI tool for managing software installations and configurations on macOS/Linux laptops.

## Overview

Dum is a hand-rolled CLI tool similar to Ansible, written in Go. It manages software installations and configurations using a playbook-based system defined in YAML.

## Installation

```shell
make clean build
ls dist
```

## Usage

### Subcommands

| Command | Alias | Description |
|---------|-------|-------------|
| `install` | `i` | Runs plays and tasks for software installations |
| `list` | `ls` | Lists plays and tasks from the playbook |
| `rename` | `r` | Renames files with various transformations |
| `log` | `lg` | Logging utilities for shell scripts |

### Install Command

Installs software based on the playbook configuration:

```shell
dum install                    # Install all plays
dum install --group <name>    # Install specific play group
dum install --dry-run        # Preview without installing
dum install -f <config.yml>   # Use custom config file
```

### List Command

Lists available plays and tasks:

```shell
dum list                      # List all plays
dum list --group <name>       # List specific play group
dum list -f <config.yml>      # Use custom config file
```

### Rename Command

Renames files with various transformations:

```shell
dum rename -s <source> -r <replace> <files...>  # Replace text in filenames
dum rename -l <files...>                         # Convert to lowercase
dum rename -f <n> <files...>                   # Trim n chars from front
dum rename -b <n> <files...>                   # Trim n chars from back
dum rename -i <n> <files...>                   # Add numeric suffix starting at n
dum rename --dry-run <files...>                # Preview without renaming
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

By default, the playbook is loaded from `~/.config/installer.yml`. You can specify a custom file with the `-f` flag or set the `INSTALLER_CONFIG` environment variable.

## Task Types

- `brew` / `cellar` / `cask` - Homebrew installations
- `dir` - Create directories
- `git` - Git clone operations
- `link` - Create symbolic links
- `vscode` - VS Code extensions
- `mas` - Mac App Store apps
- `jetbrains` - JetBrains IDE plugins
