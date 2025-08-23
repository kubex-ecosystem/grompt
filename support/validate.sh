#!/usr/bin/env bash
# lib/validate.sh – Validação da versão do Go e dependências

validate_versions() {
    local _GO_SETUP='https://raw.githubusercontent.com/rafa-mori/gosetup/main/go.sh'
    local go_version
    go_version=$(go version | awk '{print $3}' | tr -d 'go' || echo "")
    if [[ -z "$go_version" ]]; then
        log error "Go is not installed or not found in PATH."
        return 1
    fi
    local version_target=""
    version_target="$(grep '^go ' go.mod | awk '{print $2}')"
    if [[ -z "$version_target" ]]; then
        log error "Could not determine the target Go version from go.mod."
        return 1
    fi
    if [[ "$go_version" != "$version_target" ]]; then
      local _go_installation_output
      if [[ -t 0 ]]; then
        _go_installation_output="$(bash -c "$(curl -sSfL "${_GO_SETUP}")" -s --version "$version_target" >/dev/tty)"
      else
        _go_installation_output="$(export NON_INTERACTIVE=true; bash -c "$(curl -sSfL "${_GO_SETUP}")" -s --version "$version_target")"
      fi
      if [[ $? -ne 0 ]]; then
          log error "Failed to install Go version ${version_target}. Output: ${_go_installation_output}"
          return 1
      fi
    fi
    local _DEPENDENCIES=( $(cat "${_ROOT_DIR:-$(git rev-parse --show-toplevel)}/internal/module/info/manifest.json" | jq -r '.dependencies[]?') )
    check_dependencies "${_DEPENDENCIES[@]}" || return 1
    return 0
}

check_dependencies() {
  for dep in "$@"; do
    if ! command -v "$dep" > /dev/null; then
      if ! dpkg -l --selected-only "$dep" | grep "$dep" -q >/dev/null; then
        log error "$dep is not installed." true
        if [[ -z "${_NON_INTERACTIVE:-}" ]]; then
          log warn "$dep is required for this script to run." true
          local answer=""
          if [[ -z "${_FORCE:-}" ]]; then
            log question "Would you like to install it now? (y/n)" true
            read -r -n 1 -t 10 answer || answer="n"
          elif [[ "${_FORCE:-n}" == [Yy] ]]; then
            log warn "Force mode is enabled. Installing $dep without confirmation."
            answer="y"
          fi
          if [[ $answer =~ ^[Yy]$ ]]; then
            sudo apt-get install -y "$dep" || {
              log error "Failed to install $dep. Please install it manually."
              return 1
            }
            log info "$dep has been installed successfully."
          fi
        else
          log warn "$dep is required for this script to run. Installing..." true
          if [[ $_FORCE =~ ^[Yy]$ ]]; then
            log warn "Force mode is enabled. Installing $dep without confirmation."
            sudo apt-get install -y "$dep" || {
            log error "Failed to install $dep. Please install it manually."
              return 1
            }
            log info "$dep has been installed successfully."
          else
            log error "Failed to install $dep. Please install it manually before running this script."
            return 1
          fi
        fi
      fi
    fi
  done
}

export -f validate_versions
export -f check_dependencies
