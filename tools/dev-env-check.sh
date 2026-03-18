#!/usr/bin/env bash

set -euo pipefail

# Ensure homebrew is available
if ! command -v brew &>/dev/null; then
  echo >&2 "Homebrew is required to install dependencies. Please install Homebrew and try again."
  exit 1
fi

# Ensure the docker command is available
if ! command -v docker &>/dev/null; then
  echo >&2 "Docker is required to build and test. Please install Docker and try again."
  exit 1
fi

REQUIRED_COMMANDS=(
  "buf"
  "direnv"
  "git-lfs"
  "jq"
  "pre-commit"
  "shellcheck"
  "temporal"
  "trufflehog"
  "yq"
)

for cmd in "${REQUIRED_COMMANDS[@]}"; do
  if ! command -v "$cmd" &>/dev/null; then
    echo >&2 "$cmd is required to build and test. Please re-run scripts/dev-prep.sh and try again."
    exit 1
  fi
done

if [[ "$__MAP_DIRENV_INSTALLED" != "1" ]]; then
  echo >&2 "You haven't installed/setup direnv properly, or allowed it."
  echo >&2 "Follow the install instructions at https://direnv.net/docs/hook.html"
  exit 1
fi

echo "All good!"
