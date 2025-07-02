#!/usr/bin/env bash
# lib/utils.sh – Funções utilitárias

set -euo pipefail
set -o errtrace
set -o functrace
set -o posix
IFS=$'\n\t'

# Códigos de cor para logs
_SUCCESS="\033[0;32m"
_WARN="\033[0;33m"
_ERROR="\033[0;31m"
_INFO="\033[0;36m"
_NC="\033[0m"

log() {
  local type=${1:-info}
  local message=${2:-}
  local debug=${3:-${DEBUG:-false}}

  case $type in
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
    *)
      if [[ "$debug" == true ]]; then
        log "info" "$message" "$debug"
      fi
      ;;
  esac
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

# Cria um diretório temporário para cache
_TEMP_DIR="${_TEMP_DIR:-$(mktemp -d)}"
if [[ -d "${_TEMP_DIR}" ]]; then
    log info "Diretório temporário criado: ${_TEMP_DIR}"
else
    log error "Falha ao criar o diretório temporário."
fi

clear_script_cache() {
  trap - EXIT HUP INT QUIT ABRT ALRM TERM
  if [[ ! -d "${_TEMP_DIR}" ]]; then
    exit 0
  fi
  rm -rf "${_TEMP_DIR}" || true
  if [[ -d "${_TEMP_DIR}" ]] && sudo -v 2>/dev/null; then
    sudo rm -rf "${_TEMP_DIR}"
    if [[ -d "${_TEMP_DIR}" ]]; then
      printf '%b[_ERROR]%b ❌  %s\n' "$_ERROR" "$_NC" "Falha ao remover o diretório temporário: ${_TEMP_DIR}"
    else
      printf '%b[_SUCCESS]%b ✅  %s\n' "$_SUCCESS" "$_NC" "Diretório temporário removido: ${_TEMP_DIR}"
    fi
  fi
  exit 0
}

set_trap() {
  local current_shell=""
  current_shell=$(get_current_shell)
  case "${current_shell}" in
    *ksh|*zsh|*bash)
      declare -a FULL_SCRIPT_ARGS=("$@")
      if [[ "${FULL_SCRIPT_ARGS[*]}" =~ -d ]]; then
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
