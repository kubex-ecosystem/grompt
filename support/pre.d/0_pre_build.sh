#!/usr/bin/env bash
# shellcheck disable=SC2065,SC2015

set -o nounset  # Treat unset variables as an error
set -o errexit  # Exit immediately if a command exits with a non-zero status
set -o pipefail # Prevent errors in a pipeline from being masked
set -o errtrace # If a command fails, the shell will exit immediately
set -o functrace # If a function fails, the shell will exit immediately
shopt -s inherit_errexit # Inherit the errexit option in functions

IFS=$'\n\t'

build_frontend() {
  local _ROOT_DIR="${_ROOT_DIR:-$(git rev-parse --show-toplevel)}"

  cd "${_ROOT_DIR}/frontend" || {
      echo "Failed to change directory to ${_ROOT_DIR}/frontend"
      exit 1
  }

  if command -v npm &>/dev/null; then
      echo "Installing frontend dependencies..."
      npm i --no-audit --no-fund --prefer-offline || {
          echo "Failed to install frontend dependencies."
          exit 1
      }

      npm run build || {
          echo "Failed to build frontend assets."
          exit 1
      }

      if [[ -d './build' ]]; then
          echo "Frontend assets built successfully."
      else
          echo "Build directory does not exist."
          exit 1
      fi

      if [[ -d "${_ROOT_DIR}/internal/services/server/" ]]; then
          echo "Removing old build directory..."
          rm -rf "${_ROOT_DIR}/internal/services/server/build"
      fi

      mv './build' "${_ROOT_DIR}/internal/services/server/build" || {
          echo "Failed to move build directory to server."
          exit 1
      }



      echo "Frontend build moved to server directory successfully."
  else
      echo "npm is not installed. Please install Node.js and npm to continue."
      exit 1
  fi
}

(build_frontend) || {
  echo "An error occurred during the pre-build process."
  exit 1
}
