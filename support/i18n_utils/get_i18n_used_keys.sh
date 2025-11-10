#!/usr/bin/env bash

# set -o posix
set -o nounset  # Treat unset variables as an error
set -o errexit  # Exit immediately if a command exits with a non-zero status
set -o pipefail # Prevent errors in a pipeline from being masked
set -o errtrace # If a command fails, the shell will exit immediately
set -o functrace # If a function fails, the shell will exit immediately
# shopt -s inherit_errexit # Inherit the errexit option in functions
IFS=$'\n\t'

_ROOT_DIR="$(git rev-parse --show-toplevel 2>/dev/null)"
_ROOT_DIR="${_ROOT_DIR:-$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)}"
_STRING_UTILS_SCRIPT="$(realpath "$_ROOT_DIR/support/string_utils.sh")"

# Load string utilities
if [[ -f "$_STRING_UTILS_SCRIPT" ]]; then
  # shellcheck source=/dev/null
  source "$_STRING_UTILS_SCRIPT"
else
  echo "Error: Could not find string utilities script at '$_STRING_UTILS_SCRIPT'. Please ensure it exists." >&2
  exit 1
fi

_scan_for_i18n_usage() {
  if ! command -v rg &> /dev/null; then
    echo "The 'rg' (ripgrep) command is required but not installed. Please install it and try again."
    exit 1
  fi

  local _what_part_is_scanning="${1:-FRONTEND}"
  local _what_part_lowercase=""
  _what_part_lowercase="$(toLowerCase "$_what_part_is_scanning")"


  local _final_report_message=""

  if [[ ! -d "$_ROOT_DIR/$_what_part_lowercase" ]]; then
    echo "Error: $_what_part_lowercase directory '$_ROOT_DIR/$_what_part_lowercase' does not exist. Skipping $_what_part_lowercase i18n scan."
  else
    # Scan for i18n usage in frontend
    rg -no --pcre2 "t\(\s*['\"\$(]([A-Za-z][\w-]+)\.([A-Za-z0-9_.-]+)['\")]\s*(?:,|\))" "$_ROOT_DIR/$_what_part_lowercase" \
    | awk -F: '{print $3}' \
    | sed -E "s/^t\(['\"\`]//; s/['\"\`].*$//" \
    | sort -u > "$_ROOT_DIR/REPORT_${_what_part_is_scanning}-i18n_used_keys.txt"

      # Remove duplicates and sort the final report
    sort -u "$_ROOT_DIR/REPORT_${_what_part_is_scanning}-i18n_used_keys.txt" -o "$_ROOT_DIR/REPORT_${_what_part_is_scanning}-i18n_used_keys.txt"

    _final_report_message="${_what_part_lowercase} i18n used keys report generated at: $(_ROOT_DIR)/REPORT_${_what_part_is_scanning}-i18n_used_keys.txt"
  fi

  if [[ -z "$_final_report_message" ]]; then
    echo "No i18n usage found in $_what_part_lowercase."
    return 1
  fi

  echo "i18n used keys reports generated at: "
  echo "$_final_report_message"
}

_scan_for_i18n_usage "$@"


