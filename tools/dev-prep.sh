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

# install required dependencies
brew install \
  direnv \
  git-lfs \
  go \
  golangci-lint \
  goreleaser \
  jq \
  mockery \
  pre-commit \
  shellcheck \
  trufflehog \
  yamlfmt \
  yq

# install go tools
go_install_tools=(
  "github.com/abice/go-enum@latest"
  "golang.org/x/tools/cmd/goimports@latest"
)

for tool in "${go_install_tools[@]}"; do
  go install "$tool"
done


(direnv allow || true)
# PROJECT_ROOT="$(cd "$SCRIPT_DIR/.." && pwd)"
# (cd "$PROJECT_ROOT" && pre-commit install && pre-commit install --hook-type commit-msg || true)

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
echo "then restart your terminal to apply the changes."
echo " <<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<"
echo
