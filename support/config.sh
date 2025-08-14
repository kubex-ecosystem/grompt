#!/usr/bin/env bash
# shellcheck disable=SC2005

set -euo pipefail
set -o errtrace
set -o functrace
set -o posix

IFS=$'\n\t'

# Define environment variables for the current platform and architecture
# Converts to lowercase for compatibility
_CURRENT_PLATFORM="$(uname -s | tr '[:upper:]' '[:lower:]')"
_CURRENT_ARCH="$(uname -m | tr '[:upper:]' '[:lower:]')"

# Define the root directory (assuming this script is in lib/ under the root)
_ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
_APP_NAME="$(jq -r '.bin' "$_ROOT_DIR/info/manifest.json" 2>/dev/null || echo "$(basename "${_ROOT_DIR}")")"
_DESCRIPTION="$(jq -r '.description' "$_ROOT_DIR/info/manifest.json" 2>/dev/null || echo "No description provided.")"
_OWNER="$(jq -r '.organization' "$_ROOT_DIR/info/manifest.json" 2>/dev/null || echo "rafa-mori")"
_OWNER="${_OWNER,,}"  # Converts to lowercase
_BINARY_NAME="${_APP_NAME}"
_PROJECT_NAME="$(jq -r '.name' "$_ROOT_DIR/info/manifest.json" 2>/dev/null || echo "$_APP_NAME")"
_AUTHOR="$(jq -r '.author' "$_ROOT_DIR/info/manifest.json" 2>/dev/null || echo "Rafa Mori")"
_VERSION=$(jq -r '.version' "$_ROOT_DIR/info/manifest.json" 2>/dev/null || echo "v0.0.0")
_LICENSE="$(jq -r '.license' "$_ROOT_DIR/info/manifest.json" 2>/dev/null || echo "MIT")"
_REPOSITORY="$(jq -r '.repository' "$_ROOT_DIR/info/manifest.json" 2>/dev/null || echo "rafa-mori/${_APP_NAME}")"
_PRIVATE_REPOSITORY="$(jq -r '.private' "$_ROOT_DIR/info/manifest.json" 2>/dev/null || echo "false")"
_VERSION_GO=$(grep '^go ' "$_ROOT_DIR/go.mod" | awk '{print $2}')
_PLATFORMS_SUPPORTED="$(jq -r '.platforms[]' "$_ROOT_DIR/info/manifest.json" 2>/dev/null || echo "Linux, MacOS, Windows")"
_FORCE="${FORCE:-${_FORCE:-n}}"
_DEBUG="${DEBUG:-${_DEBUG:-n}}"
_DRY_RUN="${DRY_RUN:-${_DRY_RUN:-n}}"
_NON_INTERACTIVE="${NON_INTERACTIVE:-${_NON_INTERACTIVE:-n}}"

_ABOUT="  Name: ${_PROJECT_NAME} (${_APP_NAME})
  Version: ${_VERSION}
  License: ${_LICENSE}
  Supported OS: ${_PLATFORMS_SUPPORTED}
  Description: ${_DESCRIPTION}
  Author: ${_AUTHOR}
  Organization: https://github.com/${_OWNER}
  Repository: ${_REPOSITORY}
  Notes:
  - The binary is compiled with Go ${_VERSION_GO}
  - To report issues, visit: ${_REPOSITORY}/issues
###########################################################################################"

_BANNER="###########################################################################################

               ██   ██ ██     ██ ██████   ████████ ██     ██
              ░██  ██ ░██    ░██░█░░░░██ ░██░░░░░ ░░██   ██
              ░██ ██  ░██    ░██░█   ░██ ░██       ░░██ ██
              ░████   ░██    ░██░██████  ░███████   ░░███
              ░██░██  ░██    ░██░█░░░░ ██░██░░░░     ██░██
              ░██░░██ ░██    ░██░█    ░██░██        ██ ░░██
              ░██ ░░██░░███████ ░███████ ░████████ ██   ░░██
              ░░   ░░  ░░░░░░░  ░░░░░░░  ░░░░░░░░ ░░     ░░"

# Paths for the build
_CMD_PATH="$_ROOT_DIR/cmd"
_BUILD_PATH="$(dirname "$_CMD_PATH")"
_BINARY="$_BUILD_PATH/$_APP_NAME"
_LOCAL_BIN="${HOME:-"~"}/.local/bin"
_GLOBAL_BIN="/usr/local/bin"

show_about() {
    printf '%s\n\n' "${_ABOUT:-}"
}

show_banner() {
    printf '%s\n\n' "${_BANNER:-}"
}

show_headers() {
    show_banner || return 1
    show_about || return 1
}

export -f show_about
export -f show_banner
export -f show_headers
