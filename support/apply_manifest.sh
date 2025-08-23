#!/usr/bin/env bash

set -euo pipefail
set -o errtrace
set -o functrace
set -o posix

IFS=$'\n\t'

__get_values_from_manifest() {
  # Define the root directory (assuming this script is in lib/ under the root)
  _ROOT_DIR="$(cd "$(dirname "${0}")/.." && pwd)" || return 1

  # shellcheck disable=SC2005
  _APP_NAME="$(jq -r '.bin' "$_ROOT_DIR/internal/module/info/manifest.json" 2>/dev/null || echo "$(basename "${_ROOT_DIR}")")" || return 1
  _DESCRIPTION="$(jq -r '.description' "$_ROOT_DIR/internal/module/info/manifest.json" 2>/dev/null || echo "No description provided.")" || return 1
  _OWNER="$(jq -r '.organization' "$_ROOT_DIR/internal/module/info/manifest.json" 2>/dev/null || echo "rafa-mori")" || return 1
  _OWNER="${_OWNER,,}" || return 1
  _BINARY_NAME="${_APP_NAME}" || return 1
  _PROJECT_NAME="$(jq -r '.name' "$_ROOT_DIR/internal/module/info/manifest.json" 2>/dev/null || echo "$_APP_NAME")" || return 1
  _AUTHOR="$(jq -r '.author' "$_ROOT_DIR/internal/module/info/manifest.json" 2>/dev/null || echo "Rafa Mori")" || return 1
  _VERSION=$(jq -r '.version' "$_ROOT_DIR/internal/module/info/manifest.json" 2>/dev/null || echo "v0.0.0") || return 1
  _LICENSE="$(jq -r '.license' "$_ROOT_DIR/internal/module/info/manifest.json" 2>/dev/null || echo "MIT")" || return 1
  _REPOSITORY="$(jq -r '.repository' "$_ROOT_DIR/internal/module/info/manifest.json" 2>/dev/null || echo "rafa-mori/${_APP_NAME}")" || return 1
  _PRIVATE_REPOSITORY="$(jq -r '.private' "$_ROOT_DIR/internal/module/info/manifest.json" 2>/dev/null || echo "false")" || return 1

  return 0
}

__replace_project_name() {
  local _old_bin_name="goforge"
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
    "$_ROOT_DIR/goforge.go"
    "$_ROOT_DIR/"**/goforge.go
  )
  for _file in "${_files_to_rename[@]}"; do
    if [[ -f "$_file" ]]; then
      local _new_file="${_file//goforge/$_BINARY_NAME}"
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

apply_manifest() {
  __get_values_from_manifest || return 1
  __replace_project_name || return 1
  return 0
}

export -f apply_manifest
