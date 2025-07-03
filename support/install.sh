#!/usr/bin/env bash
# shellcheck disable=SC2065,SC2015

# Script Metadata
__secure_logic_version="1.0.0"
__secure_logic_date="$( date +%Y-%m-%d )"
__secure_logic_author="Rafael Mori"
__secure_logic_use_type="exec"
__secure_logic_init_timestamp="$(date +%s)"
__secure_logic_elapsed_time=0

# Check if verbose mode is enabled
if [[ "${MYNAME_VERBOSE:-false}" == "true" ]]; then
  set -x  # Enable debugging
fi

IFS=$'\n\t'

__secure_logic_sourced_name() {
  local _self="${BASH_SOURCE-}"
  _self="${_self//${_kbx_root:-$()}/}"
  _self="${_self//\.sh/}"
  _self="${_self//\-/_}"
  _self="${_self//\//_}"
  echo "_was_sourced_${_self//__/_}"
  return 0
}

__first(){
  if [ "$EUID" -eq 0 ] || [ "$UID" -eq 0 ]; then
    echo "Please do not run as root." 1>&2 > /dev/tty
    exit 1
  elif [ -n "${SUDO_USER:-}" ]; then
    echo "Please do not run as root, but with sudo privileges." 1>&2 > /dev/tty
    exit 1
  else
    # shellcheck disable=SC2155
    local _ws_name="$(__secure_logic_sourced_name)"

    if test "${BASH_SOURCE-}" != "${0}"; then
      if test $__secure_logic_use_type != "lib"; then
        echo "This script is not intended to be sourced." 1>&2 > /dev/tty
        echo "Please run it directly." 1>&2 > /dev/tty
        exit 1
      fi
      # If the script is sourced, we set the variable to true
      # and export it to the environment without changing
      # the shell options.
      export "${_ws_name}"="true"
    else
      if test $__secure_logic_use_type != "exec"; then
        echo "This script is not intended to be executed directly." 1>&2 > /dev/tty
        echo "Please source it instead." 1>&2 > /dev/tty
        exit 1
      fi
      # If the script is executed directly, we set the variable to false
      # and export it to the environment. We also set the shell options
      # to ensure a safe execution.
      export "${_ws_name}"="false"
      set -o errexit # Exit immediately if a command exits with a non-zero status
      set -o nounset # Treat unset variables as an error when substituting
      set -o pipefail # Return the exit status of the last command in the pipeline that failed
      set -o errtrace # If a command fails, the shell will exit immediately
      set -o functrace # If a function fails, the shell will exit immediately
      shopt -s inherit_errexit # Inherit the errexit option in functions
    fi
  fi
}

_DEBUG=${DEBUG:-false}
_HIDE_ABOUT=${HIDE_ABOUT:-false}

__first "$@" >/dev/tty || exit 1

# Carrega os arquivos de biblioteca
_SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
#shellcheck source=/dev/null
test -z "${_BANNER:-}" && source "${_SCRIPT_DIR}/config.sh" || true
#shellcheck source=/dev/null
test -z "$(declare -f log)" >/dev/null && source "${_SCRIPT_DIR}/utils.sh" || true
#shellcheck source=/dev/null
test -z "$(declare -f what_platform)" >/dev/null && source "${_SCRIPT_DIR}/platform.sh" || true
#shellcheck source=/dev/null
test -z "$(declare -f check_dependencies)" >/dev/null && source "${_SCRIPT_DIR}/validate.sh" || true
#shellcheck source=/dev/null
test -z "$(declare -f detect_shell_rc)" >/dev/null && source "${_SCRIPT_DIR}/install_funcs.sh" || true
#shellcheck source=/dev/null
test -z "$(declare -f build_binary)" >/dev/null && source "${_SCRIPT_DIR}/build.sh" || true
#shellcheck source=/dev/null
test -z "$(declare -f show_banner)" >/dev/null && source "${_SCRIPT_DIR}/info.sh" || true

# Inicializa os traps
set_trap "$@"

clear_screen

__main() {
  if ! what_platform; then
    log error "Plataforma não suportada: ${_PLATFORM}"
    exit 1
  fi

  if [[ "${_DEBUG}" != true ]]; then
    show_headers
    if [[ -z "${_HIDE_ABOUT}" ]]; then
      show_about
    fi
  else
    log info "Modo debug ativado; banner será ignorado..."
    if [[ -z "${_HIDE_ABOUT}" ]]; then
      show_about
    fi
  fi

  _ARGS=( "$@" )
  local default_label='Auto detect'
  local arrArgs=( "${_ARGS[@]:0:$#}" )
  local PLATFORM_ARG
  PLATFORM_ARG=$(_get_os_from_args "${arrArgs[1]:-${_PLATFORM}}")
  local ARCH_ARG
  ARCH_ARG=$(_get_arch_arr_from_args "${arrArgs[2]:-${_ARCH}}")

  log info "Comando: ${arrArgs[0]:-}" true
  log info "Plataforma: ${PLATFORM_ARG:-$default_label}" true
  log info "Arquitetura: ${ARCH_ARG:-$default_label}" true
  log info "Args: ${_ARGS[*]:-}" true

  case "${arrArgs[0]:-}" in
    build|BUILD|-b|-B)
      # validate_versions
      log info "Executando comando de build..."
      build_binary "${PLATFORM_ARG}" "${ARCH_ARG}" || exit 1
      ;;
    install|INSTALL|-i|-I)
      log info "Executando comando de instalação..."
      read -r -p "Deseja baixar o binário pré-compilado? [y/N] (Caso contrário, fará build local): " choice </dev/tty
      log info "Escolha do usuário: ${choice}"
      if [[ "$choice" == "y" || "$choice" == "Y" ]]; then
          log info "Baixando binário pré-compilado..."
          install_from_release
      else
          log info "Realizando build local..."
          validate_versions
          build_binary "${PLATFORM_ARG}" "${ARCH_ARG}" || exit 1
          install_binary
      fi
      summary
      ;;
    clear|clean|CLEAN|-c|-C)
      log info "Executando comando de limpeza..."
      clean_artifacts
      log success "Clean executado com sucesso."
      ;;
    *)
      log error "Comando inválido: ${arrArgs[0]:-}"
      echo "Uso: $0 {build|install|clean}"
      ;;
  esac
}

# Função para limpar artefatos de build
clean_artifacts() {
    log info "Limpando artefatos de build..."
    local platforms=("windows" "darwin" "linux")
    local archs=("amd64" "386" "arm64")
    for platform in "${platforms[@]}"; do
        for arch in "${archs[@]}"; do
            local output_name
            output_name=$(printf '%s_%s_%s' "${_BINARY}" "${platform}" "${arch}")
            if [[ "${platform}" != "windows" ]]; then
                local compress_name="${output_name}.tar.gz"
            else
                output_name="${output_name}.exe"
                local compress_name="${_BINARY}_${platform}_${arch}.zip"
            fi
            rm -f "${output_name}" || true
            rm -f "${compress_name}" || true
        done
    done
    log success "Artefatos de build removidos."
}

__secure_logic_main() {
  local _ws_name
  _ws_name="$(__secure_logic_sourced_name)"
  local _ws_name_val
  _ws_name_val=$(eval "echo \${$_ws_name}")
  if test "${_ws_name_val}" != "true"; then
    __main "$@"
    return $?
  else
    # If the script is sourced, we export the functions
    log error "This script is not intended to be sourced."
    log error "Please run it directly."
    return 1
  fi
}

# echo "MAKE ARGS: ${ARGS[*]:-}"
log info "Starting installation script..."
__secure_logic_main "$@"

__secure_logic_elapsed_time="$(($(date +%s) - __secure_logic_init_timestamp))"

if [[ "${MYNAME_VERBOSE:-false}" == "true" || "${_DEBUG:-false}" == "true" ]]; then
  log info "Script executed in ${__secure_logic_elapsed_time} seconds."
fi

# End of script logic