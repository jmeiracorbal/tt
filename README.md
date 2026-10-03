# tt

[![Go Version](https://img.shields.io/badge/go-1.21+-00ADD8?style=flat-square&logo=go)](https://go.dev)
[![License](https://img.shields.io/badge/license-MIT-blue?style=flat-square)](LICENSE)
[![Go Report Card](https://goreportcard.com/badge/github.com/jmeiracorbal/tt?style=flat-square)](https://goreportcard.com/report/github.com/jmeiracorbal/tt)
[![Release](https://img.shields.io/github/v/release/jmeiracorbal/tt?style=flat-square)](https://github.com/jmeiracorbal/tt/releases)

Terminal UI for [Tilt](https://tilt.dev) — monitor your local services without leaving the terminal.

![tt screenshot](docs/screenshot.png)

## Features

- Live resource status (building / ok / error)
- Per-resource log viewer with scroll
- Keyboard-driven navigation
- Connects to any running Tilt instance

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

## Requirements

- Go 1.21+
- [Tilt](https://tilt.dev) running locally

## Example project

An example project with FastAPI + Docker Compose is available under [`../example/`](../example/). Run it with:

```bash
cd example
tilt up
```

Then in another terminal:

```bash
tt
```

## License

MIT
