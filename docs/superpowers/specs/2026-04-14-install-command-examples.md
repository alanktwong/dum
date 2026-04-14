# Design: Add Task Examples to Install Command

**Date**: 2026-04-14
**Status**: Approved

## Overview

Add comprehensive YAML examples to the `dum install` command's help text to demonstrate how to configure each task type in `installer.yml`.

## Problem

The install command lists task types in its Long description but provides no examples of the YAML structure. Users must guess the correct format for each task type.

## Solution

Add an `Examples:` section to the command's Long description showing minimal YAML for all 11 task types.

## Implementation

Modify `internal/cmd/install.go` - append examples after the "Logging:" section in the `Long` string array.

## Examples Format

Each task type shown with required fields only:
- `type`: Task type identifier
- `id`: Unique task identifier
- `description`: Human-readable description

Optional fields demonstrated where applicable (e.g., `root`, `apps`, `command`, `script`).

## Task Types Covered

| Type | Key Fields | Description |
|------|-----------|-------------|
| dir | id | Create directory |
| link | id, root | Create symbolic link |
| git | id | Clone repository |
| brew | id | Install Homebrew formula |
| cask | id | Install Homebrew cask |
| cellar | id | Install Homebrew cellar |
| bash | id, command | Run bash command |
| bash | id, script | Run bash script |
| vscode | id | Install VS Code extension |
| mas | id | Install Mac App Store app |
| jetbrains | id, apps | Install JetBrains plugin |
| function | id | Call custom function |

## Acceptance Criteria

1. `dum install --help` shows examples for all task types
2. Each example uses minimal/required fields only
3. Examples are valid YAML syntax
4. Help text remains readable and not overly long