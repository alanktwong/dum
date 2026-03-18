#!/usr/bin/env bash
set -eou pipefail

# determine directory of this script, resolving symlinks
SCRIPT_DIR="$(cd "$(dirname "$(readlink "$0" || echo "$0")")" && pwd)"

set -euo pipefail

BASE_DIR="$(cd "$SCRIPT_DIR/.." && pwd)"

# Ensure that git-lfs is installed into the repo properly
grep -q '^git lfs pre-push' "${BASE_DIR}/.git/hooks/pre-push" || {
  echo >&2 "git-lfs is not installed properly. Please run 'git lfs install' and try again."
  exit 1
}

exit 0
