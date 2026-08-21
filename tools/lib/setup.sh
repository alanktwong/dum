#!/usr/bin/env bash

# Development dependency setup and validation for the repository dispatcher.

setup_install_tools() {
    # Ensure Homebrew is available before installing dependencies.
    if ! command -v brew &>/dev/null; then
        echo >&2 "Homebrew is required to install dependencies. Please install Homebrew and try again."
        return 1
    fi

    # Install required dependencies.
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
        svu \
        trufflehog \
        yamlfmt \
        yq

    # Install Go tools.
    local -a go_install_tools=(
        "github.com/abice/go-enum@latest"
        "golang.org/x/tools/cmd/goimports@latest"
        "github.com/vladopajic/go-test-coverage/v2@latest"
    )
    local tool
    for tool in "${go_install_tools[@]}"; do
        go install "$tool"
    done

    (direnv allow || true)

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
}

setup_check_tools() {
    # Ensure Homebrew is available.
    if ! command -v brew &>/dev/null; then
        echo >&2 "Homebrew is required to install dependencies. Please install Homebrew and try again."
        return 1
    fi

    if [[ ":$PATH:" != *":$GOPATH/bin:"* ]] && [[ ":$PATH:" != *"$GOPATH/bin"* ]]; then
        echo >&2 "\$GOPATH/bin is not in your \$PATH."
        echo >&2 "Add the following to your shell config: export PATH=\"\$PATH:\$GOPATH/bin\""
        return 1
    fi

    local -a required_commands=(
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
        "svu"
        "trufflehog"
        "yamlfmt"
        "yq"
    )
    local command_name
    for command_name in "${required_commands[@]}"; do
        if ! command -v "$command_name" &>/dev/null; then
            echo >&2 "$command_name is required to build and test. Please run './tools/build.sh install' and try again."
            return 1
        fi
    done

    echo "All good!"
}

setup_ci_tools() {
    go install github.com/abice/go-enum@latest
    go install github.com/vektra/mockery/v3@latest
    go install golang.org/x/tools/cmd/goimports@latest
}
