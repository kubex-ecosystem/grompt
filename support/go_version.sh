#!/usr/bin/env bash

set -euo pipefail

# ====== Go Version Management ======
# Ensures exact Go version match with go.mod
# Integrates with GoSetup for installation

get_required_go_version() {
  local go_mod_path="${1:-go.mod}"

  if [[ ! -f "${go_mod_path}" ]]; then
    echo "1.21" # fallback
    return
  fi

  # Extract go version from go.mod
  awk '/^go / {print $2; exit}' "${go_mod_path}"
}

get_current_go_version() {
  if ! command -v go >/dev/null 2>&1; then
    echo "not-installed"
    return
  fi

  go version | awk '{print $3}' | sed 's/go//'
}

check_go_version_compatibility() {
  local required_version current_version

  required_version="${1:-$(get_required_go_version "go.mod")}"
  current_version="${2:-$(get_current_go_version)}"

  if [[ "${current_version}" == "not-installed" ]]; then
    log error "Go is not installed"
    return 1
  fi

  if [[ "${current_version}" != "${required_version}" ]]; then
    log warn "Go version mismatch:"
    log warn "  Required: ${required_version} (from go.mod)"
    log warn "  Current:  ${current_version}"
    log warn "  Use GoSetup to install: gosetup --version ${required_version}"
    return 1
  fi

  log info "Go version OK: ${current_version}"
  return 0
}

auto_install_go_with_gosetup() {
  local required_version go_setup_url

  required_version="${1:-$(get_required_go_version "go.mod")}"
  go_setup_url='https://raw.githubusercontent.com/kubex-ecosystem/gosetup/main/go.sh'

  log info "Installing Go ${required_version} using GoSetup..."

  local go_installation_output
  if [[ -t 0 ]]; then
    # Interactive mode
    go_installation_output="$(bash -c "$(curl -sSfL "${go_setup_url}")" -s --version "${required_version}" 2>&1)"
  else
    # Non-interactive mode
    go_installation_output="$(export NON_INTERACTIVE=true; bash -c "$(curl -sSfL "${go_setup_url}")" -s --version "${required_version}" 2>&1)"
  fi

  # shellcheck disable=SC2181
  if [[ $? -eq 0 ]]; then
    log success "Go ${required_version} installed successfully via GoSetup"
    log info "GoSetup output: ${go_installation_output}"
    return 0
  else
    log error "Failed to install Go ${required_version} via GoSetup"
    log error "Output: ${go_installation_output}"
    return 1
  fi
}

# Export functions
export -f get_required_go_version
export -f get_current_go_version
export -f check_go_version_compatibility
export -f auto_install_go_with_gosetup
