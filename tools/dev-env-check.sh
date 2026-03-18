#!/usr/bin/env bash

set -euo pipefail

# Ensure homebrew is available
if ! command -v brew &>/dev/null; then
  echo >&2 "Homebrew is required to install dependencies. Please install Homebrew and try again."
  exit 1
fi

if [[ ":$PATH:" != *":$GOPATH/bin:"* ]] && [[ ":$PATH:" != *"$GOPATH/bin"* ]]; then
  echo >&2 "\$GOPATH/bin is not in your \$PATH."
  echo >&2 "Add the following to your shell config: export PATH=\"\$PATH:\$GOPATH/bin\""
  exit 1
fi


REQUIRED_COMMANDS=(
  "direnv"
  "git-lfs"
  "go"
  "go-enum"
  "goimports"
  "golangci-lint"
  "goreleaser"
  "jq"
  "mockery"
  "pre-commit"
  "shellcheck"
  "trufflehog"
  "yamlfmt"
  "yq"
)

for cmd in "${REQUIRED_COMMANDS[@]}"; do
  if ! command -v "$cmd" &>/dev/null; then
    echo >&2 "$cmd is required to build and test. Please re-run scripts/dev-prep.sh and try again."
    exit 1
  fi
done

echo "All good!"
