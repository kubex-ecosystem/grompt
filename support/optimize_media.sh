#!/usr/bin/env bash
# shellcheck disable=SC2015,SC1091

# Script Metadata
__opt_md_secure_logic_version="1.0.0"
__opt_md_secure_logic_date="$( date +%Y-%m-%d )"
__opt_md_secure_logic_author="Rafael Mori"
__opt_md_secure_logic_use_type="lib"
__opt_md_secure_logic_init_timestamp="$(date +%s)"
__opt_md_secure_logic_elapsed_time=0

# Check if verbose mode is enabled
if [[ "${MYNAME_VERBOSE:-false}" == "true" ]]; then
  set -x  # Enable debugging
fi

IFS=$'\n\t'

_ROOT_DIR="$(git rev-parse --show-toplevel 2>/dev/null)"
_ROOT_DIR="${_ROOT_DIR:-$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)}"

declare -a _opt_md_main_args=( "$@" )

__get_output_tty() {
  if [[ -t 1 ]]; then
    echo '/dev/tty'
  else
    echo '/dev/stderr'
  fi
}

__opt_md_secure_logic_sourced_name() {
  # prefer BASH_SOURCE[0], fallback to $0 for shells where BASH_SOURCE may be absent
  local _self="${BASH_SOURCE[0]:-$0}"
  _self="${_self//${_ROOT_DIR:-}/}"
  _self="${_self//.sh/}"
  _self="${_self//-/_}"
  _self="${_self//\//_}"
  _self="${_self//./_}"
  _self="${_self// /_}"
  echo "_was_sourced_${_self//__/_}"
  return 0
}

__opt_md_first(){
  if [ "$EUID" -eq 0 ] || [ "$UID" -eq 0 ]; then
    printf '%s\n' "Please do not run as root." >"$(__get_output_tty)"
    exit 1
  elif [ -n "${SUDO_USER:-}" ]; then
    printf '%s\n' "Please do not run as root, but with sudo privileges." >"$(__get_output_tty)"
    exit 1
  else
    # shellcheck disable=SC2155
    local _ws_name="$(__opt_md_secure_logic_sourced_name)"

    # detect if script was sourced: compare BASH_SOURCE[0] (or fallback) with $0
    if [ "${BASH_SOURCE[0]:-$0}" != "$0" ]; then
      # if test ${__opt_md_secure_logic_use_type:-} != "lib"; then
      #   printf '%s\n' "This script is not intended to be sourced." >"$(__get_output_tty)"
      #   printf '%s\n' "Please run it directly." >"$(__get_output_tty)"
      #   exit 1
      # fi
      # If the script is sourced, we set the variable to true
      # and export it to the environment without changing
      # the shell options.
      export "${_ws_name:-}"="true"
    else
      if test "${__opt_md_secure_logic_use_type:-}" != "exec"; then
        printf '%s\n' "This script is not intended to be executed directly." >"$(__get_output_tty)"
        printf '%s\n' "Please source it instead." >"$(__get_output_tty)"
        exit 1
      fi
      # If the script is executed directly, we set the variable to false
      # and export it to the environment. We also set the shell options
      # to ensure a safe execution.
      export "${_ws_name:-}"="false"
      set -o errexit # Exit immediately if a command exits with a non-zero status
      set -o nounset # Treat unset variables as an error when substituting
      set -o pipefail # Return the exit status of the last command in the pipeline that failed
      set -o errtrace # If a command fails, the shell will exit immediately
      set -o functrace # If a function fails, the shell will exit immediately
      # shopt -s inherit_errexit # Inherit the errexit option in functions

      if [[ "${_DEBUG:-}" == "true" ]]; then
        set -x
      fi
    fi
  fi
}

_QUIET=${_QUIET:-${QUIET:-false}}
_DEBUG=${_DEBUG:-${DEBUG:-false}}
_HIDE_ABOUT=${_HIDE_ABOUT:-${HIDE_ABOUT:-false}}

__opt_md_first "${_opt_md_main_args[@]}" >&2 || {
  echo "Error: This script must be run directly, not sourced." >&2
  exit 1
}

__opt_md_source_script_if_needed() {
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

_SCRIPT_DIR="${_ROOT_DIR:-$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)}/support"
__opt_md_source_script_if_needed "log" "${_SCRIPT_DIR}/log.sh" || exit 1
__opt_md_source_script_if_needed "__get_values_from_manifest" "${_SCRIPT_DIR}/apply_manifest.sh" || exit 1

__opt_md_main_functions() {
  if [[ $# -gt 0 ]]; then
    local func_name="${1:-}"
    local _full_args=( "${_main_args[@]}" )
    if declare -F "$func_name" >'/dev/null' 2>&1; then
      "$func_name" "${@:2}"
    else
      if test -x "${1:-}"; then
        _default "${_full_args[@]:1}"
      else
        log error "Function '${func_name}' not found."
        return 1
      fi
    fi
    return $?
  fi
}

optimize_media() {
  __get_values_from_manifest

  local _media_dir="${_ROOT_DIR}/media"
  local _backup_dir=""
  _backup_dir="${_ROOT_DIR}/bkp/media_optimization_$(date +%Y%m%d_%H%M%S)"
  local _image_quality=85  # Adjust quality as needed (0-100)
  local _max_width=1920    # Maximum width for resizing
  local _max_height=1080   # Maximum height for resizing

  if ! command -v mogrify &>/dev/null; then
    printf '%s\n' "ImageMagick is not installed. Please install it to optimize images." >"$(__get_output_tty)"
    return 1
  fi

  if [[ ! -d "$_media_dir" ]]; then
    printf '%s\n' "Media directory '$_media_dir' does not exist. Exiting." >"$(__get_output_tty)"
    return 1
  fi

  mkdir -p "$_backup_dir"

  find "$_media_dir" -type f \( -iname '*.jpg' -o -iname '*.jpeg' \) | while read -r img; do
    local img_rel_path="${img#"$_media_dir"/}"
    local backup_img_path="$_backup_dir/$img_rel_path"
    mkdir -p "$(dirname "$backup_img_path")"
    cp "$img" "$backup_img_path"
    mogrify -strip -interlace Plane -quality "$_image_quality" -resize "${_max_width}x${_max_height}>" "$img"
    printf '%s\n' "Optimized JPEG: $img (backup at $backup_img_path)" >"$(__get_output_tty)"
  done

  find "$_media_dir" -type f -iname '*.png' | while read -r img; do
    local img_rel_path="${img#"$_media_dir"/}"
    local backup_img_path="$_backup_dir/$img_rel_path"
    mkdir -p "$(dirname "$backup_img_path")"
    cp "$img" "$backup_img_path"
    mogrify -strip -quality "$_image_quality" -resize "${_max_width}x${_max_height}>" "$img"
    printf '%s\n' "Optimized PNG: $img (backup at $backup_img_path)" >"$(__get_output_tty)"
  done

  printf '%s\n' "Image optimization completed. Backups are stored in '$_backup_dir'." >"$(__get_output_tty)"
  return 0
}

__secure_logic_opt_md_main() {
  local _ws_name
  _ws_name="$(__opt_md_secure_logic_sourced_name)"
  local _ws_name_val
  _ws_name_val=$(eval "echo \${${_ws_name:-}}")
  if test "${_ws_name_val:-}" != "true"; then
    __opt_md_main_functions "${_main_args[@]}"
    return $?
  else
    # shellcheck disable=SC2207
    local _functions_to_export=( $(compgen -A function | grep -v -e '^__') )
    for func in "${_functions_to_export[@]}"; do
      # shellcheck disable=SC2163
      export -f "$func" || {
        echo "Error: Could not export function '$func'." >&2
        return 1
      }
    done
    return 0
  fi
}

__secure_logic_opt_md_main "${_main_args[@]}"

