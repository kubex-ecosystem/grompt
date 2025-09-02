#!/usr/bin/env bash
# shellcheck disable=SC2005

# set -o posix
set -o nounset  # Treat unset variables as an error
set -o errexit  # Exit immediately if a command exits with a non-zero status
set -o pipefail # Prevent errors in a pipeline from being masked
set -o errtrace # If a command fails, the shell will exit immediately
set -o functrace # If a function fails, the shell will exit immediately
shopt -s inherit_errexit # Inherit the errexit option in functions
IFS=$'\n\t'

# Define the relative path to the manifest file
_MANIFEST_SUBPATH=${_MANIFEST_SUBPATH:-'internal/module/info/manifest.json'}

# Define environment variables for the current platform and architecture
# Converts to lowercase for compatibility
_CURRENT_PLATFORM="$(uname -s | tr '[:upper:]' '[:lower:]')"
_CURRENT_ARCH="$(uname -m | tr '[:upper:]' '[:lower:]')"

# Define variables to hold manifest values
_ROOT_DIR="${_ROOT_DIR:-}"
_APP_NAME="${_APP_NAME:-}"
_DESCRIPTION="${_DESCRIPTION:-}"
_OWNER="${_OWNER:-}"
_BINARY_NAME="${_BINARY_NAME:-}"
_PROJECT_NAME="${_PROJECT_NAME:-}"
_AUTHOR="${_AUTHOR:-}"
_VERSION="${_VERSION:-}"
_LICENSE="${_LICENSE:-}"
_REPOSITORY="${_REPOSITORY:-}"
_PRIVATE_REPOSITORY="${_PRIVATE_REPOSITORY:-}"
_VERSION_GO="${_VERSION_GO:-}"
_PLATFORMS_SUPPORTED="${_PLATFORMS_SUPPORTED:-}"

__source_script_if_needed() {
  local _check_declare="${1:-}"
  local _script_path="${2:-}"
  # shellcheck disable=SC2065
  if test -z "$(declare -f "${_check_declare:-}")" >/dev/null; then
    # shellcheck source=/dev/null
    source "${_script_path:-}" || {
      echo "Error: Could not source ${_script_path:-}. Please ensure it exists." >&2
      return 1
    }
  fi
  return 0
}

# Quiet, force, debug, hide about, dry run
_QUIET="${QUIET:-${_QUIET:-false}}"
_FORCE="${FORCE:-${_FORCE:-false}}"
_DEBUG="${DEBUG:-${_DEBUG:-false}}"
_HIDE_ABOUT="${HIDE_ABOUT:-${_HIDE_ABOUT:-false}}"
_DRY_RUN="${DRY_RUN:-${_DRY_RUN:-false}}"
_NON_INTERACTIVE="${NON_INTERACTIVE:-${_NON_INTERACTIVE:-n}}"

# Paths for the build
_CMD_PATH="${_ROOT_DIR:-}/cmd"
_BUILD_PATH="$(dirname "${_CMD_PATH:-}")"
_BINARY="${_BUILD_PATH:-}/${_APP_NAME:-}"
_LOCAL_BIN="${HOME:-"~"}/.local/bin"
_GLOBAL_BIN="/usr/local/bin"

_SCRIPT_DIR="$(cd "$(dirname "${0:-${BASH_SOURCE[0]}}")" && pwd)"
__source_script_if_needed "apply_manifest" "${_SCRIPT_DIR:-}/apply_manifest.sh" || exit 1
__source_script_if_needed "get_current_shell" "${_SCRIPT_DIR:-}/utils.sh" || exit 1


if [[ -z "${_ROOT_DIR:-}" || -z "${_APP_NAME:-}" || -z "${_DESCRIPTION:-}" || -z "${_OWNER:-}" || -z "${_OWNER:-}" || -z "${_BINARY_NAME:-}" || -z "${_PROJECT_NAME:-}" || -z "${_AUTHOR:-}" || -z "${_VERSION:-}" || -z "${_LICENSE:-}" || -z "${_REPOSITORY:-}" || -z "${_PRIVATE_REPOSITORY:-}" || -z "${_VERSION_GO:-}" ]]; then
  apply_manifest "$@" || return 1
fi

show_about() {
  local _build_target=""
  _build_target="${BUILD_TARGET:-${_BUILD_TARGET:-}}"
  _build_target="${_build_target:-${_BUILD_TARGET:-}}"

  local _about=""
  local _about_origin=""
  local _about_repo=""

  local _platform="${_PLATFORM:-}"
  local _arch="${_ARCH:-}"

  _about_repo="  Repository: ${_REPOSITORY:-}
  Version: ${_VERSION:-}
  Description: ${_DESCRIPTION:-}
  Supported OS: ${_PLATFORMS_SUPPORTED:-}
  Notes:
  - The binary is compiled with Go ${_VERSION_GO:-}
  - To report issues, visit: ${_REPOSITORY:-}/issues"

  _about_origin="  Author: ${_AUTHOR:-}
  License: ${_LICENSE:-}
  Organization: https://github.com/${_OWNER:-}"

  if [[ "${_QUIET:-false}" == "true" ]]; then
    # _about_origin=""
    _about_repo=""
  fi

  _about="  Name: ${_PROJECT_NAME:-} (${_APP_NAME:-})
${_about_origin:-}
${_about_repo:-}"

  if [[ "${_HIDE_ABOUT:-false}" == "true" ]]; then
    _about=""
  fi

  _about=$(printf '%s\n\n' "${_about:-}")

  log hr " "

  printf '%s\n' "${_about:-}" >&2 || true

  log hr " "
}

show_banner() {
  if [[ "${_QUIET:-false}" != "true" && "${_HIDE_BANNER:-false}" != "true" ]]; then
    printf '%s\n' "#####################################################

               ██   ██ ██     ██ ██████   ████████ ██     ██
              ░██  ██ ░██    ░██░█░░░░██ ░██░░░░░ ░░██   ██
              ░██ ██  ░██    ░██░█   ░██ ░██       ░░██ ██
              ░████   ░██    ░██░██████  ░███████   ░░███
              ░██░██  ░██    ░██░█░░░░ ██░██░░░░     ██░██
              ░██░░██ ░██    ░██░█    ░██░██        ██ ░░██
              ░██ ░░██░░███████ ░███████ ░████████ ██   ░░██
              ░░   ░░  ░░░░░░░  ░░░░░░░  ░░░░░░░░ ░░     ░░" >&2
  fi
}

show_headers() {
  show_banner || return 1
  show_about || return 1
}

show_summary() {
  local install_dir="$_BINARY"
  local _cmd_executed=
  check_path "$install_dir"
}

export -f show_about
export -f show_banner
export -f show_headers
export -f show_summary
