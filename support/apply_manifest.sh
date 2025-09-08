#!/usr/bin/env bash

# set -o posix
set -o nounset  # Treat unset variables as an error
set -o errexit  # Exit immediately if a command exits with a non-zero status
set -o pipefail # Prevent errors in a pipeline from being masked
set -o errtrace # If a command fails, the shell will exit immediately
set -o functrace # If a function fails, the shell will exit immediately
shopt -s inherit_errexit # Inherit the errexit option in functions
IFS=$'\n\t'

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

_MANIFEST_SUBPATH=${_MANIFEST_SUBPATH:-'internal/module/info/manifest.json'}

__get_values_from_manifest() {
  # # Define the root directory (assuming this script is in lib/ under the root)
  _ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"

  # shellcheck disable=SC2005
  _APP_NAME="$(jq -r '.bin' "$_ROOT_DIR/$_MANIFEST_SUBPATH" 2>/dev/null || echo "$(basename "${_ROOT_DIR:-}")")"
  _DESCRIPTION="$(jq -r '.description' "$_ROOT_DIR/$_MANIFEST_SUBPATH" 2>/dev/null || echo "No description provided.")"
  _OWNER="$(jq -r '.organization' "$_ROOT_DIR/$_MANIFEST_SUBPATH" 2>/dev/null || echo "kubex-ecosystem")"
  _OWNER="${_OWNER,,}"  # Converts to lowercase
  _BINARY_NAME="${_APP_NAME}"
  _PROJECT_NAME="$(jq -r '.name' "$_ROOT_DIR/$_MANIFEST_SUBPATH" 2>/dev/null || echo "$_APP_NAME")"
  _AUTHOR="$(jq -r '.author' "$_ROOT_DIR/$_MANIFEST_SUBPATH" 2>/dev/null || echo "Rafa Mori")"
  _VERSION=$(jq -r '.version' "$_ROOT_DIR/$_MANIFEST_SUBPATH" 2>/dev/null || echo "v0.0.0")
  _LICENSE="$(jq -r '.license' "$_ROOT_DIR/$_MANIFEST_SUBPATH" 2>/dev/null || echo "MIT")"
  _REPOSITORY="$(jq -r '.repository' "$_ROOT_DIR/$_MANIFEST_SUBPATH" 2>/dev/null || echo "kubex-ecosystem/${_APP_NAME}")"
  _PRIVATE_REPOSITORY="$(jq -r '.private' "$_ROOT_DIR/$_MANIFEST_SUBPATH" 2>/dev/null || echo "false")"
  _VERSION_GO=$(grep '^go ' "$_ROOT_DIR/go.mod" | awk '{print $2}')
  _PLATFORMS_SUPPORTED="$(jq -r '.platforms[]' "$_ROOT_DIR/$_MANIFEST_SUBPATH" 2>/dev/null || echo "linux, macOS, windows")"
  _PLATFORMS_SUPPORTED="$(printf '%s ' "${_PLATFORMS_SUPPORTED[*]//
/, }")" # Converts to comma-separated list
  _PLATFORMS_SUPPORTED="${_PLATFORMS_SUPPORTED,,}"  # Converts to lowercase

  return 0
}

__replace_project_name() {
  local _old_bin_name="gobe"
  local _new_bin_name="${_BINARY_NAME}"

  if [[ ! -d "$_ROOT_DIR/bkp" ]]; then
    mkdir -p "$_ROOT_DIR/bkp"
  fi

  # Backup the original files before making changes
  tar --exclude='bkp' --exclude='*.tar.gz' --exclude='go.sum' -czf "$_ROOT_DIR/bkp/$(date +%Y%m%d_%H%M%S)_goforge_backup.tar.gz" -C "$_ROOT_DIR" . || {
    log fatal "Could not create backup. Please check if the directory exists and is writable." true
    return 1
  }

  local _files_to_remove=(
    "$_ROOT_DIR/README.md"
    "$_ROOT_DIR/CHANGELOG.md"
    "$_ROOT_DIR/docs/README.md"
    "$_ROOT_DIR/docs/assets/*"
    "$_ROOT_DIR/go.sum"
  )
  for _file in "${_files_to_remove[@]}"; do
    if [[ -f "$_file" ]]; then
      rm -f "$_file" || {
        log error "Could not remove $_file. Please check if the file exists and is writable." true
        continue
      }
      log info "Removed $_file"
    else
      log warn "File $_file does not exist, skipping."
    fi
  done

  local _files_to_rename=(
    "$_ROOT_DIR/go${_old_bin_name}.go"
    "$_ROOT_DIR/"**"/${_old_bin_name}.go"
  )
  for _file in "${_files_to_rename[@]}"; do
    if [[ -f "$_file" ]]; then
      local _new_file="${_file//${_old_bin_name}/$_BINARY_NAME}"
      mv "$_file" "$_new_file" || {
        log error "Could not rename $_file to $_new_file. Please check if the file exists and is writable." true
        continue
      }
      log info "Renamed $_file to $_new_file"
    else
      log warn "File $_file does not exist, skipping."
    fi
  done

  local _files_to_update=(
    "$_ROOT_DIR/go.mod"
    "$_ROOT_DIR/"**/*.go
    "$_ROOT_DIR/"**/*.md
    "$_ROOT_DIR/"*/*.go
    "$_ROOT_DIR/"*.md
  )
  for _file in "${_files_to_update[@]}"; do
    if [[ -f "$_file" ]]; then
      sed -i "s/$_old_bin_name/$_new_bin_name/g" "$_file" || {
        log error "Could not update $_file. Please check if the file exists and is writable." true
        continue
      }
      log info "Updated $_file"
    else
      log warn "File $_file does not exist, skipping."
    fi
  done

  cd "$_ROOT_DIR" || {
    log error "Could not change directory to $_ROOT_DIR. Please check if the directory exists." true
    return 1
  }

  go mod tidy || {
    log error "Could not run 'go mod tidy'. Please check if Go is installed and configured correctly." true
    return 1
  }

  return 0
}

change_project_name() {
  __replace_project_name || return 1
  return 0
}

apply_manifest() {
  __get_values_from_manifest || return 1
  return 0
}

export -f apply_manifest
export -f change_project_name
