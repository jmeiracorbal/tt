#!/bin/bash
# tt installer
# Usage: curl -sSf https://raw.githubusercontent.com/jmeiracorbal/tt/main/install.sh | bash
#
# Environment overrides:
#   TT_VERSION=v1.0.0 bash install.sh
#   TT_INSTALL_DIR=/usr/local/bin bash install.sh
#   TT_DRY_RUN=true bash install.sh

set -e

REPO="jmeiracorbal/tt"
INSTALL_DIR="${TT_INSTALL_DIR:-$HOME/.local/bin}"
DRY_RUN="${TT_DRY_RUN:-false}"
TT_VERSION="${TT_VERSION:-}"

# ── helpers ────────────────────────────────────────────────────────────────────

info() { printf "\033[1;34m[tt]\033[0m %s\n" "$*"; }
ok()   { printf "\033[1;32m[tt]\033[0m %s\n" "$*"; }
err()  { printf "\033[1;31m[tt]\033[0m %s\n" "$*" >&2; exit 1; }
warn() { printf "\033[1;33m[tt]\033[0m %s\n" "$*"; }

dry() {
  if [ "$DRY_RUN" = "true" ]; then
    printf "\033[2m  (dry-run) %s\033[0m\n" "$*"
  else
    eval "$@"
  fi
}

# ── detect platform ────────────────────────────────────────────────────────────

detect_platform() {
  local os arch

  case "$(uname -s)" in
    Darwin) os="darwin" ;;
    Linux)  os="linux" ;;
    *)      err "Unsupported OS: $(uname -s)" ;;
  esac

  case "$(uname -m)" in
    arm64|aarch64) arch="arm64" ;;
    x86_64)        arch="amd64" ;;
    *)             err "Unsupported architecture: $(uname -m)" ;;
  esac

  echo "${os}-${arch}"
}

# ── fetch helpers ──────────────────────────────────────────────────────────────

fetch() {
  local url="$1" dest="$2"
  if command -v curl >/dev/null 2>&1; then
    curl -sSfL "$url" -o "$dest"
  elif command -v wget >/dev/null 2>&1; then
    wget -qO "$dest" "$url"
  else
    err "curl or wget required"
  fi
}

fetch_stdout() {
  local url="$1"
  if command -v curl >/dev/null 2>&1; then
    curl -sSfL "$url"
  else
    wget -qO- "$url"
  fi
}

probe_url() {
  local url="$1"
  if command -v curl >/dev/null 2>&1; then
    curl -sSfI "$url" >/dev/null 2>&1
  else
    wget -q --spider "$url" 2>/dev/null
  fi
}

# ── version resolution ─────────────────────────────────────────────────────────

fetch_latest_version() {
  local version
  version=$(fetch_stdout "https://api.github.com/repos/${REPO}/releases/latest" \
    | grep '"tag_name"' \
    | sed 's/.*"tag_name": *"\([^"]*\)".*/\1/')
  [ -z "$version" ] && err "Could not fetch latest release version"
  echo "$version"
}

check_version_compat() {
  local version="$1" platform="$2"
  local base_url="https://github.com/${REPO}/releases/download/${version}"

  info "Checking compatibility of pinned version ${version}..."
  probe_url "${base_url}/tt-${platform}.sha256" \
    || err "Release ${version} does not ship a binary for ${platform}. Unset TT_VERSION to use the latest."
  ok "Release ${version} is compatible."
}

# ── download and verify ────────────────────────────────────────────────────────

download_binary() {
  local version="$1" platform="$2"
  local base_url="https://github.com/${REPO}/releases/download/${version}"
  local binary_url="${base_url}/tt-${platform}"
  local checksum_url="${base_url}/tt-${platform}.sha256"
  local dest="${INSTALL_DIR}/tt"

  info "Downloading tt ${version} for ${platform}..."

  if [ "$DRY_RUN" = "true" ]; then
    dry "curl -sSfL \"${binary_url}\" -o \"${dest}\""
    dry "curl -sSfL \"${checksum_url}\" | shasum -a 256 -c"
    dry "chmod +x \"${dest}\""
    return
  fi

  mkdir -p "$INSTALL_DIR"

  local tmp checksum_file
  tmp=$(mktemp)
  checksum_file=$(mktemp)
  trap 'rm -f "$tmp" "$checksum_file"' EXIT

  fetch "$binary_url" "$tmp"       || err "Download failed: ${binary_url}"
  fetch_stdout "$checksum_url" > "$checksum_file" || err "Checksum download failed: ${checksum_url}"

  local expected actual
  expected=$(awk '{print $1}' "$checksum_file")
  actual=$(shasum -a 256 "$tmp" | awk '{print $1}')

  [ "$expected" = "$actual" ] || err "Checksum mismatch — aborting. Expected: ${expected}, got: ${actual}"

  mv "$tmp" "$dest"
  chmod +x "$dest"
  ok "Installed: ${dest}"
}

# ── PATH check ─────────────────────────────────────────────────────────────────

check_path() {
  if ! echo "$PATH" | tr ':' '\n' | grep -qx "$INSTALL_DIR"; then
    warn "${INSTALL_DIR} is not in your PATH."
    warn "Add this to your shell profile (~/.zshrc or ~/.bashrc):"
    warn "  export PATH=\"\$HOME/.local/bin:\$PATH\""
  fi
}

# ── main ───────────────────────────────────────────────────────────────────────

main() {
  [ "$DRY_RUN" = "true" ] && info "Dry-run mode — no changes will be made"

  local platform version
  platform=$(detect_platform)
  version="${TT_VERSION:-$(fetch_latest_version)}"

  info "Latest release: ${version}"

  if [ -n "${TT_VERSION}" ] && [ "$DRY_RUN" != "true" ]; then
    check_version_compat "$version" "$platform"
  fi

  download_binary "$version" "$platform"
  check_path

  ok "Done. Run 'tt --help' to get started."
}

main "$@"
