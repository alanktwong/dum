# Viper Configuration Integration Design

**Date**: 2026-08-13

## Goal

Use `github.com/spf13/viper` for Dum's environment-backed configuration resolution while keeping all Viper interaction inside `cmd/dum`. Preserve current configuration and logging behavior.

## Scope and boundaries

Viper is used only by `cmd/dum/cli` and command wiring. It owns resolution of:

- `INSTALLER_CONFIG`
- `XDG_CONFIG_HOME`
- `ZSH_LOG_LEVEL`

`internal/` packages must not import Viper or read environment variables for these settings.

`internal/yaml` remains responsible for parsing, typing, and validating `installer.yml`. `internal/factory` continues receiving `factory.InputOptions` with a resolved file path. `internal/logging` remains a logger abstraction and constructor, without environment lookup.

## Configuration resolution

`cli.GetDefaultConfig()` creates and configures an isolated Viper instance per call. Resolution precedence remains:

1. `INSTALLER_CONFIG`, when non-empty.
2. `XDG_CONFIG_HOME/dum/installer.yml`, when `XDG_CONFIG_HOME` is non-empty.
3. `~/.config/dum/installer.yml`.

The returned default preserves existing `~` behavior. `--file` remains a Cobra flag whose explicit value overrides the environment-derived default.

Viper must not parse or unmarshal installer YAML.

## Logging resolution

The same isolated Viper configuration path reads `ZSH_LOG_LEVEL`. `cmd/dum/cli` converts supported strings and numeric values to `charmbracelet/log.Level`:

- named levels: `debug`, `info`, `warn`, `error`, `fatal`
- numeric thresholds retain current behavior: `<=10` debug, `<=20` info, `<=30` warn, `<=40` error, `<=50` fatal
- empty, invalid, and values over 50 resolve to warn

`cmd/dum/root.go` resolves the level and passes it to `lg.NewLogger`. `cmd/dum/log` no longer calls `lg.EnvLevel()`. Remove `internal/logging.EnvLevel()` and its environment-parser tests.

## Data flow

```text
NewDum
  -> cmd/dum/cli resolver
  -> typed config path + log level
  -> Cobra defaults / logger options
  -> command execution
  -> factory.InputOptions.File
  -> internal/factory
  -> internal/yaml typed loader
```

No Viper object or runtime Viper state crosses into `internal/`. Resolver calls use isolated instances to prevent test leakage and command-order effects.

## Error handling

Environment values retain current fallback behavior; malformed or missing environment values do not fail command construction. Installer file and YAML errors remain owned by `internal/factory` and `internal/yaml`, preserving existing wrapped errors.

## Tests and verification

Update or add `cmd/dum` tests for:

- `INSTALLER_CONFIG` precedence
- `XDG_CONFIG_HOME` fallback
- home fallback
- empty environment values
- explicit `--file` override
- isolated resolver calls
- named and numeric log-level parsing
- invalid and over-50 log-level fallback
- logger initialization with resolved level

Remove obsolete `internal/logging.EnvLevel` tests. Verify with:

```sh
go test ./cmd/dum/...
go test ./...
make lint
make build
```

Command names, aliases, flags, help text, installer YAML behavior, logging output, and release version behavior remain unchanged except for implementation ownership of environment resolution.
