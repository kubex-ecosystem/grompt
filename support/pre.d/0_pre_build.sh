#!/usr/bin/env bash
# shellcheck disable=SC2065,SC2015

set -o nounset  # Treat unset variables as an error
set -o errexit  # Exit immediately if a command exits with a non-zero status
set -o pipefail # Prevent errors in a pipeline from being masked
set -o errtrace # If a command fails, the shell will exit immediately
set -o functrace # If a function fails, the shell will exit immediately
shopt -s inherit_errexit # Inherit the errexit option in functions

IFS=$'\n\t'

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

_SCRIPT_DIR="$(cd "$(dirname "${0}")" && pwd)"
__source_script_if_needed "get_current_shell" "${_SCRIPT_DIR:-}/utils.sh" || {
  echo "Error: Could not source utils.sh. Please ensure it exists." >&2
  exit 1
}

build_frontend() {
  local _ROOT_DIR="${_ROOT_DIR:-$(git rev-parse --show-toplevel)}"

  if [[ ! -d "${_ROOT_DIR}/frontend" ]]; then
      echo "Frontend directory does not exist."
      exit 1
  fi

  cd "${_ROOT_DIR}/frontend" || {
    log fatal "Failed to change directory to ${_ROOT_DIR}/frontend" || echo "Failed to change directory to ${_ROOT_DIR}/frontend" >&2
    exit 1
  }

  if command -v npm &>/dev/null; then
      log info "Building frontend..." true

      _frontend_install_output="$(npm i --no-fund --no-audit --loglevel=error --legacy-peer-deps || {
          echo "Failed to install frontend dependencies."
      })"

      _frontend_build_output="$(npm run build --silent > /dev/null 2>&1 || {
          echo "Failed to build frontend assets."
      })"

      if [[ "${_frontend_build_output:-}" == "Failed to build frontend assets." ]] || [[ -n "${_frontend_build_output:-}" && "${_QUIET:-false}" != "true" ]]; then
          log error "${_frontend_build_output}" true
          log fatal "Frontend build failed." true
          exit 1
      fi

      if [[ -d './dist' ]]; then
          log success "Frontend assets built successfully." true
      else
          log fatal "Build directory does not exist." true
          exit 1
      fi

      if [[ -d "${_ROOT_DIR}/internal/grompt/embedded/" ]]; then
          log notice "Removing old build directory..."
          rm -rf "${_ROOT_DIR}/internal/grompt/embedded/guiweb" || true
      fi

      mv './dist' "${_ROOT_DIR}/internal/grompt/embedded/guiweb" || {
          log fatal "Failed to move build directory to server." true
          exit 1
      }

      log success "Frontend build moved to server directory successfully." true
  else
      log fatal "npm is not installed. Please install Node.js and npm to continue." true
      exit 1
  fi
}

(build_frontend) || {
  log fatal "An error occurred during the pre-build process." | tee /dev/stderr
  exit 1
}
