#!/usr/bin/env bash

# set -o posix
set -o nounset  # Treat unset variables as an error
set -o errexit  # Exit immediately if a command exits with a non-zero status
set -o pipefail # Prevent errors in a pipeline from being masked
set -o errtrace # If a command fails, the shell will exit immediately
set -o functrace # If a function fails, the shell will exit immediately
shopt -s inherit_errexit # Inherit the errexit option in functions
IFS=$'\n\t'

get_release_url() {
    local os="${_PLATFORM%%-*}"
    local format
    if [[ "$os" == "windows" ]]; then
      format="zip"
    else
      format="tar.gz"
    fi
    local arch="${_PLATFORM##*-}"
    local release_url="${_REPOSITORY}/releases/download/${_VERSION}/${_APP_NAME}_${_VERSION}_${os}_${arch}.${format}"
    echo "${release_url}"
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
      "x86_64"|"amd64") _arch="amd64" ;;
      "armv6") _arch="armv6l" ;;
      "armv8"|"aarch64"|"arm64") _arch="arm64" ;;
      *386*) _arch="386" ;;
    esac
    platform="linux-${_arch}"
    ;;
  *Darwin*)
    _os="darwin"
    case "${_arch}" in
      "x86_64"|"amd64") _arch="amd64" ;;
      "armv8"|"aarch64"|"arm64") _arch="arm64" ;;
    esac
    platform="darwin-${_arch}"
    ;;
  MINGW*|MSYS*|CYGWIN*|Win*)
    _os="windows"
    case "${_arch}" in
      "x86_64"|"amd64") _arch="amd64" ;;
      "armv8"|"aarch64"|"arm64") _arch="arm64" ;;
    esac
    platform="windows-${_arch}"
    ;;
  *)
    log error "Unsupported OS: ${_os} with architecture: ${_arch}"
    log error "Please report this issue to the project maintainers."
    return 1
    ;;
  esac

  export _PLATFORM_WITH_ARCH="${platform//-/_}"
  export _PLATFORM="${_os}"
  export _ARCH="${_arch}"

  return 0
}

_get_os_arr_from_args() {
  local _platform="${1:-"$(uname -s | tr '[:upper:]' '[:lower:]')"}"

  case "${_platform}" in
    all|ALL|a|A|-a|-A)
      echo "windows darwin linux"
      return 0
    ;;

    win|WIN|windows|WINDOWS|w|W|-w|-W)
      echo "windows"
      return 0
    ;;

    linux|LINUX|l|L|-l|-L)
      echo "linux"
      return 0
    ;;

    darwin|DARWIN|macOS|MACOS|m|M|-m|-M)
      echo "darwin"
      return 0
    ;;

    *)
      local _argArr=( "${_platform}" )
      for arg in "${_argArr[@]}"; do
        echo "${arg}"
      done
    ;;
  esac
}

_get_arch_arr_from_args() {
  local _platform="${1:-"$(uname -s | tr '[:upper:]' '[:lower:]')"}"
  local _arch="${2:-"$(uname -m | tr '[:upper:]' '[:lower:]')"}"

  case "${_platform:-"$(uname -s | tr '[:upper:]' '[:lower:]')"}" in
    darwin|DARWIN|macOS|MACOS|m|M|-m|-M)
      case "${_arch:-"$(uname -m | tr '[:upper:]' '[:lower:]')"}" in
        all|ALL|a|A|-a|-A)
          echo "amd64 arm64"
          return 0
          ;;
        armv8|arm64|ARM64|aarch64|AARCH64)
          echo "arm64"
          return 0
          ;;
        amd64|AMD64|x86_64|X86_64|x64|X64)
          echo "amd64"
          return 0
          ;;
        *)
          log fatal "Invalid architecture: '$_arch'. Valid options: amd64, arm64."
          return 1
          ;;
      esac
    ;;

    linux|LINUX|l|L|-l|-L)
      case "${_arch:-"$(uname -m | tr '[:upper:]' '[:lower:]')"}" in
        all|ALL|a|A|-a|-A)
          echo "amd64 arm64 armv6l 386"
          return 0
          ;;
        armv8|arm64|ARM64|aarch64|AARCH64)
          echo "arm64"
          ;;
        amd64|AMD64|x86_64|X86_64|x64|X64)
          echo "amd64"
          ;;
        386|I386)
          echo "386"
          ;;
        armv6l|ARMV6L)
          echo "armv6l"
          ;;
        *)
          log fatal "Invalid architecture: '$_arch'. Valid options: amd64, arm64, 386, armv6l."
          return 1
          ;;
      esac
    ;;

    windows|WINDOWS|w|W|-w|-W)
      case "${_arch:-"$(uname -m | tr '[:upper:]' '[:lower:]')"}" in
        all|ALL|a|A|-a|-A)
          echo "amd64 arm64"
          return 0
          ;;
        amd64|AMD64|x86_64|X86_64|x64|X64)
          echo "amd64"
          return 0
          ;;
        *)
          log fatal "Invalid architecture: '${_arch:-}'. Valid options: amd64, 386." true
          return 1
          ;;
      esac
    ;;

    *)
      log fatal "${_arch:-} is invalid for ${_platform:-}." true
    ;;
  esac

  return 1
}

_get_os_from_args() {
  local _platform="${1:-"$(uname -s | tr '[:upper:]' '[:lower:]')"}"

  case "${_platform:-"$(uname -s | tr '[:upper:]' '[:lower:]')"}" in
    all|ALL|a|A|-a|-A)
      echo "all"
    ;;

    win|WIN|windows|WINDOWS|w|W|-w|-W)
      echo "windows"
    ;;

    linux|LINUX|l|L|-l|-L)
      echo "linux"
    ;;

    darwin|DARWIN|macOS|MACOS|m|M|-m|-M)
      echo "darwin"
    ;;

    *)
      log fatal "Invalid platform: '${_platform:-}'. Valid options: windows, linux, darwin, all."
    ;;

  esac
}

_get_arch_from_args() {
  local _platform="${1:-$(uname -s | tr '[:upper:]' '[:lower:]')}"
  local _arch="${2:-$(uname -m | tr '[:upper:]' '[:lower:]')}"

  # Normalize common arch names
  case "${_arch}" in
    # First we handle with different names
    x86_64|X86_64) echo "amd64" ;;
    armv8|aarch64|AARCH64) echo "arm64" ;;
    i386|I386) echo "386" ;;
    ARMV6L) echo "armv6l" ;;

    # Then we handle with common names, including all option
    all|ALL|a|A|-a|-A) echo "all" ;;

    amd64|arm64|386|armv6l) echo "${_arch}" ;;
    *) uname -m | tr '[:upper:]' '[:lower:]' ;; # "${_arch}" ;;
  esac
}

export -f _get_os_arr_from_args
export -f _get_arch_arr_from_args
export -f _get_os_from_args
export -f _get_arch_from_args
export -f get_release_url
export -f what_platform
