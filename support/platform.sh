#!/usr/bin/env bash

set -euo pipefail
set -o errtrace
set -o functrace
set -o posix
IFS=$'\n\t'

get_release_url() {
    local os="${_PLATFORM%%-*}"
    local format
    if [[ "$os" == "windows" ]]; then
      format="zip"
    else
      format="tar.gz"
    fi
    echo "'https://github.com/${_OWNER}/${_PROJECT_NAME}/releases/download/${_VERSION}/${_PROJECT_NAME}_.${format}'"
}

what_platform() {
  local _os
  _os="$(uname -s)"
  local _arch
  _arch="$(uname -m)"
  local platform=""

  case "${_os}" in
  *Linux*|*Nix*)
    _os="linux"
    case "${_arch}" in
      "x86_64") _arch="amd64" ;;
      "armv6") _arch="armv6l" ;;
      "armv8"|"aarch64") _arch="arm64" ;;
      *386*) _arch="386" ;;
    esac
    platform="linux-${_arch}"
    ;;
  *Darwin*)
    _os="darwin"
    case "${_arch}" in
      "x86_64") _arch="amd64" ;;
      "arm64") _arch="arm64" ;;
    esac
    platform="darwin-${_arch}"
    ;;
  MINGW*|MSYS*|CYGWIN*|Win*)
    _os="windows"
    case "${_arch}" in
      "x86_64") _arch="amd64" ;;
      "arm64") _arch="arm64" ;;
    esac
    platform="windows-${_arch}"
    ;;
  *)
    log error "Plataforma não suportada: ${_os} ${_arch}"
    log error "Informe este problema aos mantenedores do projeto."
    return 1
    ;;
  esac

  export _PLATFORM_WITH_ARCH="${platform//-/_}"
  export _PLATFORM="${_os}"
  export _ARCH="${_arch}"

  return 0
}

_get_os_arr_from_args() {
  local _PLATFORM_ARG=$1
  if [[ "${_PLATFORM_ARG}" == "all" ]]; then
    echo "windows darwin linux"
  else
    echo "${_PLATFORM_ARG}"
  fi
}

_get_arch_arr_from_args() {
  local _ARCH_ARG=$1
  if [[ "${_ARCH_ARG}" == "all" ]]; then
    echo "amd64 386 arm64"
  else
    echo "${_ARCH_ARG}"
  fi
}

_get_os_from_args() {
  local arg=$1
  case "$arg" in
    all|ALL|a|A|-a|-A) echo "all" ;;
    win|WIN|windows|WINDOWS|w|W|-w|-W) echo "windows" ;;
    linux|LINUX|l|L|-l|-L) echo "linux" ;;
    darwin|DARWIN|macOS|MACOS|m|M|-m|-M) echo "darwin" ;;
    *)
      log error "Plataforma inválida: '${arg}'. Opções válidas: windows, linux, darwin, all."
      exit 1
      ;;
  esac
}

_get_arch_from_args() {
  local arg=$1
  case "$arg" in
    all|ALL|a|A|-a|-A) echo "all" ;;
    amd64|AMD64|x86_64|X86_64|x64|X64) echo "amd64" ;;
    arm64|ARM64|aarch64|AARCH64) echo "arm64" ;;
    386|i386|I386) echo "386" ;;
    *)
      log error "Arquitetura inválida: '${arg}'. Opções válidas: amd64, arm64, 386."
      exit 1
      ;;
  esac
}

export -f _get_os_arr_from_args
export -f _get_arch_arr_from_args
export -f _get_os_from_args
export -f _get_arch_from_args
export -f get_release_url
export -f what_platform

what_platform "${@}"
