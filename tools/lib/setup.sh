#!/usr/bin/env bash

# Development dependency setup and validation for the repository dispatcher.

setup_install_tools() {
    pre-commit install
    pre-commit install --hook-type commit-msg
    git lfs install --force
}

setup_check_tools() {
    local -a required_commands=(
        "git-lfs"
        "go"
        "go-enum"
        "go-test-coverage"
        "goimports"
        "golangci-lint"
        "goreleaser"
        "jq"
        "mockery"
        "node"
        "npm"
        "pre-commit"
        "python3"
        "shellcheck"
        "svu"
        "trufflehog"
        "yamlfmt"
        "yq"
    )
    local command_name
    for command_name in "${required_commands[@]}"; do
        if ! command -v "$command_name" &>/dev/null; then
            echo >&2 "$command_name is required to build and test. Please run 'mise install --locked' and try again."
            return 1
        fi
    done

    echo "All good!"
}
