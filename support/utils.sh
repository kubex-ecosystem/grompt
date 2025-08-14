#!/usr/bin/env bash
# lib/utils.sh – Utility functions

set -euo pipefail
set -o errtrace
set -o functrace
set -o posix
IFS=$'\n\t'

# Color codes for logs
_SUCCESS="\033[0;32m"
_WARN="\033[0;33m"
_ERROR="\033[0;31m"
_INFO="\033[0;36m"
_NOTICE="\033[0;35m"
_NC="\033[0m"

log() {
  local type=${1:-info}
  local message=${2:-}
  local debug=${3:-${DEBUG:-false}}

  case $type in
    question|_QUESTION|-q|-Q)
      if [[ "$debug" == true ]]; then
        printf '%b[_QUESTION]%b ❓  %s: ' "$_NOTICE" "$_NC" "$message"
      fi
      ;;
    notice|_NOTICE|-n|-N)
      if [[ "$debug" == true ]]; then
        printf '%b[_NOTICE]%b 📝  %s\n' "$_NOTICE" "$_NC" "$message"
      fi
      ;;
    info|_INFO|-i|-I)
      if [[ "$debug" == true ]]; then
        printf '%b[_INFO]%b ℹ️  %s\n' "$_INFO" "$_NC" "$message"
      fi
      ;;
    warn|_WARN|-w|-W)
      if [[ "$debug" == true ]]; then
        printf '%b[_WARN]%b ⚠️  %s\n' "$_WARN" "$_NC" "$message"
      fi
      ;;
    error|_ERROR|-e|-E)
      printf '%b[_ERROR]%b ❌  %s\n' "$_ERROR" "$_NC" "$message"
      ;;
    success|_SUCCESS|-s|-S)
      printf '%b[_SUCCESS]%b ✅  %s\n' "$_SUCCESS" "$_NC" "$message"
      ;;
    fatal|_FATAL|-f|-F)
      printf '%b[_FATAL]%b 💀  %s\n' "$_FATAL" "$_NC" "$message"
      if [[ "$debug" == true ]]; then
        printf '%b[_FATAL]%b 💀  %s\n' "$_FATAL" "$_NC" "Exiting due to fatal error."
      fi
      clear_build_artifacts
      exit 1
      ;;
    *)
      if [[ "$debug" == true ]]; then
        log "info" "$message" "$debug"
      fi
      ;;
  esac
  return 0
}

clear_screen() {
  printf "\033[H\033[2J"
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
if [[ -d "${_TEMP_DIR}" ]]; then
    log info "Temporary directory created: ${_TEMP_DIR}"
else
    log error "Failed to create the temporary directory."
fi

clear_script_cache() {
  trap - EXIT HUP INT QUIT ABRT ALRM TERM
  if [[ ! -d "${_TEMP_DIR}" ]]; then
    return 0
  fi
  rm -rf "${_TEMP_DIR}" || true
  if [[ -d "${_TEMP_DIR}" ]] && sudo -v 2>/dev/null; then
    sudo rm -rf "${_TEMP_DIR}"
    if [[ -d "${_TEMP_DIR}" ]]; then
      printf '%b[_ERROR]%b ❌  %s\n' "$_ERROR" "$_NC" "Failed to remove the temporary directory: ${_TEMP_DIR}"
    else
      printf '%b[_SUCCESS]%b ✅  %s\n' "$_SUCCESS" "$_NC" "Temporary directory removed: ${_TEMP_DIR}"
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
