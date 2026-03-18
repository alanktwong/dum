# Installer for setting up a Mac OSX laptop

The installer CLI is a hand-rolled (by me in Go) version of Ansible.

Why roll your own?

* I originally wrote a library of bash scripts similar in design to Omakub.
* That was a lot of code with repetitive lines, limited configurability, logging, etc.
* So translated them to Go.
* Design is similar to Ansible with similar concepts: tasks, plays, playbooks, etc.
* Immutable configuration is not a design goal like Terraform or NixOS.

Key concepts:

* Task is a specific task for installing something on a computer. E.g. `BrewTask` runs `brew install ...`
* Play is a group of tasks
* Playbook is a grouping of plays.

Why not use Ansible?

* Don't want to use a huge library that I only use a small subset of. I.e. Ansible is too big and complex for the managing the setup of my own laptop.
* Wanted to write my own tasks coordinated by a CLI using some form of configuration.

macOS gatekeeper looks for an Apple code signing cert within the download of an app.

```bash
sudo xattr -dr com.apple.com /path/to/app
```

Above is the terminal command to override quarantine by macOS gatekeeper so that it recognizes an app as legitimate.

## Building

Assumes that:

* Using Linux or MacOS
* Installed Go 1.24.4+
* Installed curl
* Installed zip/unzip
* Installed curl
* Installed tar
* Install repo in normal $GOPATH/src/dotfiles


```shell
make help
```

for the normal help

```shell
make clean build
ls dist
```

You should see an `installer-M.N.X` executor with semantic versioning.

## Running

To see help

```shell
./dum-M.N.X install --help
```

## Configuring

See installer.yml

## Legacy Doc

### Apps

* Datavault
* Dropbox
* Kaleidoscope
* Sourcetree
* Zoom
* Onedrive
* MS 365
* MS Teams
* Grand Perspective
* Firefox
* Adobe Acrobat Reader
* MySQL
* PostgresSQL
* Anaconda
* Express VPN
* Adobe Creative Cloud
* All In One Messenger
* Azure Data Studio
* GraphiQL
* [AppCleaner](https://freemacsoft.net/appcleaner/)

### installer script

#### [zsh](https://www.zsh.org/)

#### [bash](https://www.gnu.org/software/bash/)

Note Apple has not updated bash since they switched to zsh. To use the latest & greatest, you have to install a newer version so that it overrides the ancient one from Apple on the PATH.

#### [homebrew](https://brew.sh/)

Following installed via home brew

* [ack](https://beyondgrep.com/) - a grep-like source code search tool.
* [autoconf](https://www.gnu.org/software/autoconf) - Automatic configure script builder
* [bat](https://github.com/sharkdp/bat) - cat clone with syntax highlighting, line numbers and git integration
* [btop](https://github.com//btop) resource monitor TUI
* [buf](https://buf.build/product/cli) - all in one protobuf toolchain
* [cmake](https://cmake.org/) de-facto standard for building C++ code
* [cmatrix](https://github.com/abishekvashok/cmatrix) "Matrix" like effect in your terminal
* [ctop](https://ctop.sh/) - concise commandline monitoring for containers
* [direnv](https://direnv.net/) augments existing shells with a new feature that can load and unload environment variables depending on the current directory.
* [dive](https://github.com/wagoodman/dive) - tool for exploring each layer of Docker image
* [duf](https://github.com/muesli/duf) improved `df`
* [dust](https://github.com/bootandy/dust) = `du` + rust ... an improvement on `du`
* [eza](https://eza.rocks/) - modernization of `ls`
* [fastfetch](https://github.com/fastfetch-cli/fastfetch) is a neofetch-like tool for fetching system information and displaying it prettily.
* [fzf](https://junegunn.github.io/fzf/) - Command-line fuzzy finder written in Go
* [gettext](https://www.gnu.org/software/gettext/manual/gettext.html) - ?
* [grpcurl](https://github.com/fullstorydev/grpcurl) - Like cURL, but for gRPC: Command-line tool for interacting with gRPC servers
* [gtrash](https://github.com/umlx5h/gtrash) provides an alternative to `rm` that puts files in the trash, searches files in the trash, restores files from the trash and empties the trash.
* [gum](https://github.com/charmbracelet/gum) provides highly configurable, ready-to-use utilities to help you write useful shell scripts such as: choose, confirm, file, filter, format, input, join, pager, spin, log, etc.
* [jq](https://jqlang.github.io/jq/) - command-line JSON processor
* [jsonnet](https://jsonnet.org/) DSL to build JSON
* [lua](https://www.lua.org/) - Powerful, lightweight programming language
* [openssl](https://openssl.org/) - Cryptography and SSL/TLS Toolkit
* [pandoc](https://pandoc.org/) convert files from one markup format into another
* [pre-commit](https://pre-commit.com/) A framework for managing and maintaining multi-language pre-commit hooks.
* [protobuf](https://protobuf.dev/) - language-neutral, platform-neutral extensible mechanisms for serializing structured data.
* [ripgrep](https://github.com/BurntSushi/ripgrep) is a line-oriented search tool that recursively searches the current directory for a regex pattern. By default, ripgrep will respect gitignore rules and automatically skip hidden files/directories and binary files.
* [scc](https://github.com/boyter/scc) counts blank lines, comment lines, and physical lines of source code in many programming languages.
* [sevenzip](https://7-zip.org/) file archiver with a high compression ratio.
* [shellcheck](https://www.shellcheck.net/) finds bugs in shell scripts
* [tree](https://github.com/MrRaindrop/tree-cli#readme) - Displays directories as trees (with optional color/HTML output)
* [trufflehog](https://trufflesecurity.com/) secrets detection engine
* [xz](https://tukaani.org/xz/) - General-purpose data compression with high compression ratio
* [yq](https://mikefarah.gitbook.io/yq) - command line YAML processor
* [zoxide](https://github.com/ajeetdsouza/zoxide)  smarter cd command, inspired by z and autojump. It remembers which directories you use most frequently, so you can "jump" to them in just a few keystrokes.

Cellars:

* [boost](https://github.com/wg/wrk) - Collection of portable C++ source libraries
* [icu4c](https://icu.unicode.org/home) - C/C++ and Java libraries for Unicode and globalization
* [libyaml](https://icu.unicode.org/home) - YAML Parser
* [oniguruma](https://github.com/kkos/oniguruma/) - Regular expressions library
* [readline](https://tiswww.case.edu/php/chet/readline/rltop.html) - Library for command-line editing
* [sqlite](https://sqlite.org/) - A file based SQL engine
* [zlib](https://tiswww.case.edu/php/chet/readline/rltop.html) - General-purpose lossless data-compression library
* [zsh-autosuggestion](https://github.com/zsh-users/zsh-autosuggestions) - Fish-like autosuggestions for zsh
* [zsh-syntax-highlighting](https://github.com/zsh-users/zsh-syntax-highlighting) - Fish shell like syntax highlighting for Zsh.
* pkg-config - Manage compile and link flags for libraries

#### TBD

* [exiftool](https://exiftool.org/) - ExifTool is a platform-independent Perl library plus a [command-line application](https://exiftool.org/exiftool_pod.html) for reading, writing and editing meta information in a [wide variety of media files](https://exiftool.org/#supported).
* [ffmpeg](https://ffmpeg.org/) - A complete, cross-platform solution to record, convert and stream audio and video. FFmpeg is an amazing tool.
* [fish?](https://fishshell.com/) - fish is a smart and user-friendly command line shell for Linux, macOS, and the rest of the family.
* [gdbm](https://kubernetes.io/) - GNU database manager
* [httpie](https://httpie.io/) - A simple, easy-to-use, tool for making http calls. Much more intuitive than curl or wget.
* [hyperline](https://github.com/sharkdp/hyperfine) - command line benchmarking tool
* [procs](https://github.com/dalance/procs) - modern replacements for `ps`
* [sd](https://github.com/chmln/sd) - Intuitive find & replace CLI (`sed` alternative)
* [trash-cli](https://github.com/andreafrancia/trash-cli) trashes files recording the original path, deletion date, and permissions.
* [wrk](https://github.com/wg/wrk) - modern HTTP benchmarking tool

#### Apps

* [adobe-creative-cloud](https://www.adobe.com/creativecloud.html)
* [grandperspective](https://grandperspectiv.sourceforge.net/) is a small utility application for macOS that graphically shows the disk usage within a file system
* [rar](https://www.rarlab.com/) - WinRAR is a powerful archive manager.

#### Dev Tools

* [awscli](https://aws.amazon.com/cli/) - CLI for AWS
* [Azure CLI](https://learn.microsoft.com/en-us/cli/azure/)
* [google-cloud-sdk](https://cloud.google.com/sdk)
* [graphiql](https://www.gnu.org/software/gettext/manual/gettext.html) - Light, Electron-based Wrapper around GraphiQL
* [insomnia](https://insomnia.rest/) - Postman alternative
* [kaleidoscope](https://kaleidoscope.app/) - merge/diff tool
* [VS Code](https://code.visualstudio.com/)

#### IDEs

* [DataGrip](https://www.jetbrains.com/datagrip/)
* [GoLand](https://www.jetbrains.com/go/)
* [IntelliJ IDEA](https://www.jetbrains.com/idea/)
* [PyCharm](https://www.jetbrains.com/pycharm/)
* [WebStorm](https://www.jetbrains.com/webstorm/)

### [git](https://git-scm.com/)

* [gh: github cli](https://cli.github.com/)
* [glab: gitlab cli](https://gitlab.com/gitlab-org/cli#installation)

### [vim](https://www.vim.org/)

### [neovim](https://neovim.io/)

I use [lazyvim](https://www.lazyvim.org/)

### [tmux](https://github.com/tmux/tmux/wiki)

* [tpm](https://github.com/tmux-plugins/tpm) - Tmux Plugin Manager
* [tmux-resurrrect](https://github.com/tmux-plugins/tmux-resurrect) - Persists tmux environment across system restarts.

### [NerdFonts](https://www.nerdfonts.com/)

See also [powerline fonts](https://github.com/powerline/fonts) - Patched fonts for Powerline users

* font-fira-code
* font-fira-code-nerd-font
* font-hack-nerd-font
* font-jetbrains-mono
* font-jetbrains-mono-nerd-font
* font-menlo-for-powerline
* font-meslo-for-powerlevel10k
* font-meslo-for-powerline
* font-meslo-lg
* font-meslo-lg-dz
* font-meslo-lg-nerd-font

### [starship](https://starship.rs/)

A modern prompt

### Terminal Emulators

* [iterm2](https://iterm2.com/) and [iterm2 color schemes](https://iterm2colorschemes.com/) - install with git clone

#### [wezterm](https://wezfurlong.org/wezterm/index.html)

#### [kitty](https://sw.kovidgoyal.net/kitty/)

#### [ghostty](https://ghostty.org/)

### [lazygit](https://github.com/jesseduffield/lazygit) TUI

### [sdkman](https://sdkman.io/)

Manages Java, its JVM languages, its build tools, and other things in the system

### [Python](https://www.python.org/)

[Hitchhiker's Guide to Python](https://docs.python-guide.org/)

Manage with [pyenv](https://github.com/pyenv/pyenv) with [pyenv-virtualenv](https://github.com/pyenv/pyenv-virtualenv) to manage virtualenv

[pipx](https://pipx.pypa.io/stable/) - install/run python in an isolate environment

#### Package Management

`pip` is the OOTB package installer for Python. See [its home page](https://pip.pypa.io/en/stable/)
`pipenv` is a package manager for Python similar to yarn. See [its home page](https://pipenv.pypa.io/en/latest/)
[Conda Package Management](https://docs.conda.io/en/latest/) is the standard package management for Python data science.

#### Virtual Environments

This is intend to isolate a bundle of packages per project
rather than per Python system install.

Beware of polluting the space of your system Python's package space.

<https://realpython.com/python-virtual-environments-a-primer/>

`virtualenv` ... its [home page](https://virtualenv.pypa.io/en/latest/)

### [node](https://nodejs.org/en)

manage with [nvm](https://github.com/nvm-sh/nvm)

### [Go](https://go.dev/)

Its build task maanagers:

* [golangci-lint](https://golangci-lint.run/)
* [go-releaser](https://goreleaser.com/)
* [go-task](https://taskfile.dev/)
* [mage](https://npf.io/2018/09/mage/)

Its version managers

* [goenv](https://github.com/go-nv/goenv)
* [gvm](https://github.com/moovweb/gvm)

### [Ruby](https://www.ruby-lang.org/en/)

Manage with [rbenv](https://github.com/rbenv/rbenv)

* [rbenv-default-gems](https://github.com/sstephenson/rbenv-default-gems) - Auto-installs gems for Ruby Auto-installs
* [rbenv-gemset](https://github.com/jf/rbenv-gemset) - Adds basic gemset support to rbenv

### [lf](https://github.com/gokcehan/lf) - terminal file manager TUI

### [mani](https://manicli.com/)

A CLI tool that helps you manage multiple repositories. It's useful when you are working with microservices, multi-project systems, many libraries or just a bunch of repositories and want a central place for pulling all repositories and running commands over them.

Alternatives:

* [gardens](https://garden-rs.gitlab.io/index.html) is all about making it easy to remix and reuse libraries maintained in separate Git repositories.
* [gita](https://github.com/nosarthur/gita) also manages multiple git repos
* [gr](https://github.com/mixu/gr)
* [meta](https://github.com/mateodelnorte/meta)
* [mu-repo](https://github.com/fabioz/mu-repo)
* [myrepos](https://myrepos.branchable.com/)
* [repo](https://source.android.com/setup/develop/repo)
* [vcstool](https://github.com/dirk-thomas/vcstool)

### [kubernetes](https://kubernetes.io/)

* [colima](https://github.com/abiosoft/colima) CLI to manage container runtimes f
* [docker](https://www.docker.com/) & [docker-compose](https://docs.docker.com/compose/)
* [helm](https://helm.sh/) - Kubernetes package manager
* [k9s](https://k9scli.io/) TUI for Kubernetes
  * [kubectl](https://kubernetes.io/docs/reference/kubectl/) - kubernetes CLI - i.e. `kubectl`
* [minikube](https://minikube.sigs.k8s.io/docs/)
* [Rancher Desktop](https://rancherdesktop.io/)

## Old or not installed via Homebrew

### Languages

* [perl](https://www.perl.org/) - programming language
* [Rust](https://www.rust-lang.org/) - using [rustup](https://rustup.rs/)

### Todo

Apps: Adobe Lightroom,  Slack,

### To reconsider or wait

* Adobe Photoshop CC 2019
* Anaconda Navigator
* Astah Community (no longer supported)
* Cornerstone
* FileZilla
* ForkLift
* FreeMind
* Hermes JMS
* Malwarebytes Anti-Malware
* Microsoft Office?
* Omni Desk Sweeperman
* Quicken 2016
* Quicksilver
* QuickTime Player
* SecureID
* Sequel Pro (why when you have Data Grip?)
* Skype
* Skype for Businss
* Skype Meetings App
* SOAP UI 5.4.0
* Spotify
* [stats](https://github.com/exelban/stats) - system monitor for Mac
* Tunnelblick
* Virtual Box
* Visual VM
* We Chat
* We Clean

### Apple Software

* iMovie
* iPhoto