#!/usr/bin/env bash

# Development dependency setup and validation for the repository dispatcher.

setup_install_tools() {
    setup_check_git_version || return
    hk validate || return
    setup_check_no_global_hk_hooks || return
    git lfs install --force
    hk install --mise
}

setup_check_git_version() {
    local version major minor
    version="$(git --version | awk '{print $3}')"
    if [[ ! "$version" =~ ^([0-9]+)\.([0-9]+)(\.[0-9]+([.-][[:alnum:]]+)*)?$ ]]; then
        echo >&2 "Could not parse Git version '$version'; Git 2.54 or newer is required."
        return 1
    fi

    major=$((10#${BASH_REMATCH[1]}))
    minor=$((10#${BASH_REMATCH[2]}))
    if (( major < 2 || (major == 2 && minor < 54) )); then
        echo >&2 "Git 2.54 or newer is required (found $version)."
        return 1
    fi

    return 0
}

setup_check_no_global_hk_hooks() {
    local status
    if git config --global --get-regexp '^hook\.hk-' >/dev/null 2>&1; then
        echo >&2 "Global hk hooks are configured; this repository uses repo-local hk. Choose one installation scope. Run 'hk uninstall --global' only if you intend to switch to repo-local hooks."
        return 1
    else
        status=$?
        if (( status != 1 )); then
            echo >&2 "Could not inspect global Git hook configuration; refusing repo-local installation."
            return 1
        fi
    fi

    return 0
}

setup_check_tools() {
    local -a required_commands=(
        "git"
        "git-lfs"
        "go"
        "go-enum"
        "go-test-coverage"
        "goimports"
        "golangci-lint"
        "goreleaser"
        "hk"
        "jq"
        "mockery"
        "node"
        "npm"
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

    setup_check_git_version || return

    echo "All good!"
}
