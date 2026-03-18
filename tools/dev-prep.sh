#!/usr/bin/env bash
# set -x
set -euo pipefail

# determine directory of this script, resolving symlinks
SCRIPT_DIR="$(cd "$(dirname "$(readlink "$0" || echo "$0")")" && pwd)"

# ensure homebrew is available
if ! command -v brew &>/dev/null; then
  echo >&2 "Homebrew is required to install dependencies. Please install Homebrew and try again."
  exit 1
fi

# ensure the docker command is available
if ! command -v docker &>/dev/null; then
  echo >&2 "Docker is required to build and test. Please install Docker and try again."
  exit 1
fi

# install required dependencies
brew install \
  buf \
  direnv \
  git-lfs \
  jq \
  node@20 \
  pre-commit \
  shellcheck \
  temporal \
  trufflehog \
  yamlfmt \
  yq

(cd "$SCRIPT_DIR/.." && direnv allow)
(cd "$SCRIPT_DIR/.." && pre-commit install && pre-commit install --hook-type commit-msg)

echo ""
echo "Checking gradle wrapper..."
(cd "$SCRIPT_DIR" && ./check-gradle-wrapper.sh)

echo
echo " >>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>"
echo "Setup complete. You must update your ~/.bashrc, ~/.zshrc, or fish config to include the following line:"
echo
echo "  eval \"\$(direnv hook bash)\""
echo "# or "
echo "  eval \"\$(direnv hook zsh)\""
echo "# or "
echo "  direnv hook fish | source"
echo
echo "You need to add this line to your ~/.bashrc or ~/.zshrc."
# shellcheck disable=SC2016
echo 'export PATH="/opt/homebrew/opt/node@20/bin:$PATH'
echo "Or if you don't want to change your path, you can run:"
echo "brew link -f node@20"
echo
echo "then restart your terminal to apply the changes."
echo " <<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<"
echo
echo
