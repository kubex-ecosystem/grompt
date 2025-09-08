#!/usr/bin/env bash
# lib/validate.sh – Validação da versão do Go e dependências

# Source go version management functions
# shellcheck disable=SC1091
source "$(dirname "${BASH_SOURCE[0]}")/go_version.sh"

validate_versions() {
    local go_setup_url='https://raw.githubusercontent.com/kubex-ecosystem/gosetup/main/go.sh'
    local current_version required_version

    # Use modular functions for version checking
    current_version="$(get_current_go_version)"
    required_version="$(get_required_go_version)"

    if [[ "${current_version}" == "not-installed" ]]; then
        log error "Go is not installed or not found in PATH."
        return 1
    fi

    if [[ -z "$required_version" ]]; then
        log error "Could not determine the target Go version from go.mod."
        return 1
    fi

    # Check version compatibility
    if [[ "$current_version" != "$required_version" ]]; then
        log warn "Go version mismatch: current=${current_version}, required=${required_version}"

        local go_installation_output
        if [[ -t 0 ]]; then
            go_installation_output="$(bash -c "$(curl -sSfL "${go_setup_url}")" -s --version "$required_version" 2>&1)"
        else
            go_installation_output="$(export NON_INTERACTIVE=true; bash -c "$(curl -sSfL "${go_setup_url}")" -s --version "$required_version" 2>&1)"
        fi

        # shellcheck disable=SC2181
        if [[ $? -ne 0 ]]; then
            log error "Failed to install Go version ${required_version}. Output: ${go_installation_output}"
            return 1
        fi
    fi

    # Validate other dependencies from manifest
    local dependencies manifest_file
    manifest_file="${_ROOT_DIR:-$(git rev-parse --show-toplevel)}/${_MANIFEST_SUBPATH:-/internal/module/info/manifest.json}"

    if [[ -f "${manifest_file}" ]]; then
        mapfile -t dependencies < <(jq -r '.dependencies[]?' "${manifest_file}")
        check_dependencies "${dependencies[@]}" || return 1
    fi
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
          if [[ "${_FORCE:-}" =~ ^[Yy]$ ]]; then
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
