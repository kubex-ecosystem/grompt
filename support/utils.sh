#!/usr/bin/env bash
# lib/utils.sh – Utility functions

# set -o posix
set -o nounset  # Treat unset variables as an error
set -o errexit  # Exit immediately if a command exits with a non-zero status
set -o pipefail # Prevent errors in a pipeline from being masked
set -o errtrace # If a command fails, the shell will exit immediately
set -o functrace # If a function fails, the shell will exit immediately
shopt -s inherit_errexit # Inherit the errexit option in functions
IFS=$'\n\t'

# Color codes for logs
_SUCCESS="\033[0;32m"
_WARN="\033[0;33m"
_ERROR="\033[0;31m"
_INFO="\033[0;36m"
_NOTICE="\033[0;35m"
_FATAL="\033[0;41m"
_TRACE="\033[0;34m"
_NC="\033[0m"

log() {
  local _type=${1:-info}
  local _message=${2:-}
  local _debug=${3:-}
  _debug="${_debug:-${DEBUG:-${_DEBUG:-false}}}"

  case $_type in
    question|_QUESTION|-q|-Q)
      if [[ "${_debug:-false}" == "true" ]]; then
        printf '%b[QUESTION]%b ❓  %s: ' "${_NOTICE:-\033[0;35m}" "${_NC:-\033[0m}" "$_message"
      fi
      ;;
    notice|_NOTICE|-n|-N)
      if [[ "${_debug:-false}" == "true" ]]; then
        printf '%b[NOTICE]%b 📝  %s\n' "${_NOTICE:-\033[0;35m}" "${_NC:-\033[0m}" "$_message"
      fi
      ;;
    info|_INFO|-i|-I)
      if [[ "${_debug:-false}" == "true" ]]; then
        printf '%b[INFO]%b ℹ️  %s\n' "${_INFO:-\033[0;36m}" "${_NC:-\033[0m}" "$_message"
      fi
      ;;
    warn|_WARN|-w|-W)
      if [[ "${_debug:-false}" == "true" ]]; then
        printf '%b[WARN]%b ⚠️  %s\n' "${_WARN:-\033[0;33m}" "${_NC:-\033[0m}" "$_message"
      fi
      ;;
    error|_ERROR|-e|-E)
      printf '%b[ERROR]%b ❌  %s\n' "${_ERROR:-\033[0;31m}" "${_NC:-\033[0m}" "$_message" >&2
      ;;
    success|_SUCCESS|-s|-S)
      printf '%b[SUCCESS]%b ✅  %s\n' "${_SUCCESS:-\033[0;32m}" "${_NC:-\033[0m}" "$_message"
      ;;
    fatal|_FATAL|-f|-F)
      printf '%b[FATAL]%b 💀  %s\n' "${_FATAL:-\033[0;41m}" "${_NC:-\033[0m}" "Exiting due to fatal error: $_message" >&2
      clear_build_artifacts || true
      clear_script_cache || true

      # shellcheck disable=SC2317
      exit 1 || kill -9 $$
      ;;
    separator|_SEPARATOR|hr|-hr|-HR|line)
      # if [[ "${_debug:-false}" != "true" ]]; then
        local _columns=${COLUMNS:-$(tput cols || echo 80)}
        local _margin=$(( _columns - ( _columns / 2 ) ))
        _message="${_message// /¬}"
        _message="$(printf '%b%s%b %*s' "${_TRACE:-\033[0;34m}" "${_message:-}" "${_NC:-\033[0m}" "$((_columns - ( "${#_message}" + _margin )))" '')"
        _message="${_message// /\#}"
        _message="${_message//¬/ }"
        printf '%s\n' "${_message:-}" >&2
      # fi
      ;;
    *)
      log "info" "$_message" "${_debug:-false}" || true
      ;;
  esac

  return 0
}

clear_screen() {
  if [[ "${_QUIET:-false}" != "true" && "${_DEBUG:-false}" != "true" ]]; then
    printf "\033[H\033[2J"
  fi
}

get_current_shell() {
  local shell_proc
  shell_proc=$(cat /proc/$$/comm)
  case "${0##*/}" in
    ${shell_proc}*)
      local shebang
      shebang=$(head -1 "$0")
      printf '%s\n' "${shebang##*/}"
      ;;
    *)
      printf '%s\n' "$shell_proc"
      ;;
  esac
}

# Creates a temporary directory for cache
_TEMP_DIR="${_TEMP_DIR:-$(mktemp -d)}"
if [[ -d "${_TEMP_DIR:-}" ]]; then
    log info "Temporary directory created: ${_TEMP_DIR:-}"
else
    log error "Failed to create the temporary directory."
fi

clear_script_cache() {
  trap - EXIT HUP INT QUIT ABRT ALRM TERM
  if [[ ! -d "${_TEMP_DIR:-}" ]]; then
    return 0
  fi
  rm -rf "${_TEMP_DIR:-}" || true
  if [[ -d "${_TEMP_DIR:-}" ]] && sudo -v 2>/dev/null; then
    sudo rm -rf "${_TEMP_DIR:-}"
    if [[ -d "${_TEMP_DIR:-}" ]]; then
      printf '%b[ERROR]%b ❌  %s\n' "${_ERROR:-\033[0;31m}" "${_NC:-\033[0m}" "Failed to remove the temporary directory: ${_TEMP_DIR:-}"
    else
      printf '%b[SUCCESS]%b ✅  %s\n' "${_SUCCESS:-\033[0;32m}" "${_NC:-\033[0m}" "Temporary directory removed: ${_TEMP_DIR:-}"
    fi
  fi
  return 0
}


clear_build_artifacts() {
  clear_script_cache
  local build_dir="${_ROOT_DIR:-$(realpath '../')}/dist"
  if [[ -d "${build_dir}" ]]; then
    rm -rf "${build_dir}" || true
    if [[ -d "${build_dir}" ]]; then
      log error "Failed to remove build artifacts in ${build_dir}."
    else
      log success "Build artifacts removed from ${build_dir}."
    fi
  else
    log notice "No build artifacts found in ${build_dir}."
  fi
}

set_trap() {
  local current_shell=""
  current_shell=$(get_current_shell)
  case "${current_shell}" in
    *ksh|*zsh|*bash)
      declare -a FULL_SCRIPT_ARGS=("$@")
      if [[ "${FULL_SCRIPT_ARGS[*]}" == *--debug* ]]; then
          set -x
      fi
      if [[ "${current_shell}" == "bash" ]]; then
        set -o errexit
        set -o pipefail
        set -o errtrace
        set -o functrace
        shopt -s inherit_errexit
      fi
      trap 'clear_script_cache' EXIT HUP INT QUIT ABRT ALRM TERM
      ;;
  esac
}

export -f log
export -f clear_screen
export -f get_current_shell
export -f clear_script_cache
export -f clear_build_artifacts
export -f set_trap
