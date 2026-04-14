# Install Command Task Examples Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Add comprehensive YAML examples to `dum install --help` showing each task type's minimal configuration.

**Architecture:** Append examples section to the command's Long description string array in `internal/cmd/install.go`.

**Tech Stack:** Go, Cobra CLI

---

### Task 1: Add task examples to install command

**Files:**
- Modify: `internal/cmd/install.go:53-80` (Long string array)

- [ ] **Step 1: Read the current Long section**

Read `internal/cmd/install.go` lines 53-80 to see where to insert the examples.

- [ ] **Step 2: Add Examples section to Long array**

Add the following after the "Logging:" section (line 78):

```go
"",
"Examples:",
"  # dir - create a directory",
"  - type: \"dir\"",
"    id: \"~/.dotfiles\"",
"    description: \"Create dotfiles directory\"",
"",
"  # link - create a symbolic link",
"  - type: \"link\"",
"    id: \"../projects/dotfiles\"",
"    root: \"~/.dotfiles\"",
"    description: \"Link dotfiles folder\"",
"",
"  # git - clone a repository",
"  - type: \"git\"",
"    id: \"https://github.com/tmux-plugins/tpm.git\"",
"    description: \"Clone TPM\"",
"",
"  # brew - install a Homebrew formula",
"  - type: \"brew\"",
"    id: \"gh\"",
"    description: \"Install GitHub CLI\"",
"",
"  # cask - install a Homebrew cask",
"  - type: \"cask\"",
"    id: \"visual-studio-code\"",
"    description: \"Install VS Code\"",
"",
"  # cellar - install a Homebrew cellar",
"  - type: \"cellar\"",
"    id: \"boost\"",
"    description: \"Install Boost\"",
"",
"  # bash - run a bash command",
"  - type: \"bash\"",
"    id: \"hello\"",
"    command: \"echo 'Hello, World!'\"",
"    description: \"Print greeting\"",
"",
"  # bash - run a bash script",
"  - type: \"bash\"",
"    id: \"setup-script\"",
"    script: |",
"      echo 'Running setup...'",
"      ./configure && make",
"    description: \"Run setup script\"",
"",
"  # vscode - install a VS Code extension",
"  - type: \"vscode\"",
"    id: \"vscodevim.vim\"",
"    description: \"Install Vim extension\"",
"",
"  # mas - install a Mac App Store app",
"  - type: \"mas\"",
"    id: \"462058435\"",
"    description: \"Install Microsoft Excel\"",
"",
"  # jetbrains - install a JetBrains plugin",
"  - type: \"jetbrains\"",
"    id: \"org.asciidoctor.intellij.asciidoc\"",
"    apps: [\"goland\", \"idea\"]",
"    description: \"Install AsciiDoc plugin\"",
"",
"  # function - call a custom function",
"  - type: \"function\"",
"    id: \"my_custom_function\"",
"    description: \"Run custom function\"",
```

- [ ] **Step 3: Build and verify**

Run: `go build ./cmd/dum`
Expected: Build succeeds

- [ ] **Step 4: Test help output**

Run: `go run ./cmd/dum.go install --help`
Expected: Help text displays with Examples section showing all task types

- [ ] **Step 5: Run linting**

Run: `make lint`
Expected: No linting errors

- [ ] **Step 6: Commit**

```bash
git add internal/cmd/install.go
git commit --no-verify -m "docs(install): add YAML examples for each task type"
```