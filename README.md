# tilt-tui

> `tt` on the command line. A terminal UI for [Tilt](https://github.com/tilt-dev/tilt).

[![Go Version](https://img.shields.io/badge/go-1.22+-00ADD8?style=flat-square&logo=go)](https://go.dev)
[![License](https://img.shields.io/badge/license-MIT-blue?style=flat-square)](LICENSE)
[![Go Report Card](https://goreportcard.com/badge/github.com/jmeiracorbal/tt)](https://goreportcard.com/report/github.com/jmeiracorbal/tt)
[![Release](https://img.shields.io/github/v/release/jmeiracorbal/tt?style=flat-square&include_prereleases)](https://github.com/jmeiracorbal/tt/releases)

Monitor your local services without leaving the terminal.

## Features

- Live resource status (building / ok / error)
- Per-resource log viewer with scroll
- Keyboard-driven navigation
- Connects to any running Tilt instance

## Requirements

### Tilt

tt requires a running Tilt instance. Install Tilt from [tilt.dev](https://docs.tilt.dev/install.html):

```bash
# macOS / Linux
curl -fsSL https://raw.githubusercontent.com/tilt-dev/tilt/master/scripts/install.sh | bash
```

```bash
# macOS via Homebrew
brew install tilt
```

For other platforms see the [official install guide](https://docs.tilt.dev/install.html).

### Go

Go 1.22 or later (only required if building from source).

## Install

```bash
curl -sSf https://raw.githubusercontent.com/jmeiracorbal/tt/main/install.sh | bash
```

Installs to `~/.local/bin/tt`. Override the destination:

```bash
TT_INSTALL_DIR=/usr/local/bin curl -sSf https://raw.githubusercontent.com/jmeiracorbal/tt/main/install.sh | bash
```

Pin a specific version:

```bash
TT_VERSION=v1.0.0 curl -sSf https://raw.githubusercontent.com/jmeiracorbal/tt/main/install.sh | bash
```

Or build from source:

```bash
git clone https://github.com/jmeiracorbal/tt
cd tt
go build -o tt .
```

## Usage

With Tilt running (default port 10350):

```bash
tt
```

Custom port:

```bash
tt --port 10351
```

## Keybindings

| Key | Action |
|-----|--------|
| `↑` / `k` | Previous resource |
| `↓` / `j` | Next resource |
| `PgUp` / `PgDn` | Scroll logs |
| `r` | Force refresh |
| `q` / `Ctrl+C` | Quit |

## License

MIT
