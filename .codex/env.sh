#!/usr/bin/env bash
# envAI - Kubex Adaptive Context Bootstrap
# v0.1.0  |  https://kubex.world/envai

set -e

_HOME_DIR="${HOME:-"/home/user"}"

# VERSION

_ENVAI_VERSION='0.1.0'

# MAIN GUIDANCE

_MAINS_GUIDANCE_SUFFIX='kubex'
_MAINS_GUIDANCE_CONTENT="../**/.$_MAINS_GUIDANCE_SUFFIX/*.md"

_version(){
  echo "${_ENVAI_VERSION:-unknown}"
}

__show_org_files() {
  local _org_files
  _org_files=$(find . -type f -name "*.$_MAINS_GUIDANCE_SUFFIX" -print0 | xargs -0 ls -lt | awk '{print $9}')

  if [[ -z "${_org_files}" ]]; then
    echo "No '$_MAINS_GUIDANCE_SUFFIX' files found."
    return
  fi

  echo "Recently edited '$_MAINS_GUIDANCE_SUFFIX' files:"
  echo

  local _file
  for _file in ${_org_files}; do
    echo "- ${_file}"
  done
}

__show_org_files_tree() {
  tree -L 4 "${PROJECT_ROOT_DIR:-}" -I 'temp|tasks|support|scripts|backups|node_modules|.venv|.venv-hooks|lab|go.work|images|.git|.github|.gitignore|.gitattributes|*.log|*.tmp|*.bak|*.swp|*.swo'
}

__load_available_tools(){
  [ -d "${HOME:-$_HOME_DIR}/.go/bin" ] && export PATH="$PATH:${HOME:-$_HOME_DIR}/.go/bin"
  [ -d "${HOME:-$_HOME_DIR}/go/bin" ] && export PATH="$PATH:${HOME:-$_HOME_DIR}/go/bin"

  # GOPATH opcional
  [ -z "${GOPATH:-}" ] && [ -d "${HOME:-$_HOME_DIR}/go" ] && export GOPATH="${HOME:-$_HOME_DIR}/go" && export PATH="$PATH:$GOPATH/bin"

  # NVM
  export NVM_DIR="${NVM_DIR:-${HOME:-$_HOME_DIR}/.nvm}"
  if [ -s "$NVM_DIR/nvm.sh" ]; then
    # shellcheck disable=SC1090,SC1091
    . "$NVM_DIR/nvm.sh" >/dev/null 2>&1
    # Se tiver "default", ativa
    nvm_alias="$NVM_DIR/alias/default"
    if [ -f "$nvm_alias" ]; then
      ver="$(cat "$nvm_alias" 2>/dev/null || true)"
      if [ -n "$ver" ]; then
        nvm use "$ver" >/dev/null 2>&1 || true
        export NODE_PATH_SET=1
      fi
    fi
  fi

  # fallback: ensure the default Node version bin dir is on PATH even when nvm use fails
  if [ -z "${NODE_PATH_SET:-}" ]; then
    if [ -d "$NVM_DIR/versions/node" ]; then
      # pick default alias if available, otherwise latest version directory
      node_ver=""
      if [ -f "$NVM_DIR/alias/default" ]; then
        node_ver="$(cat "$NVM_DIR/alias/default" 2>/dev/null || true)"
      fi
      if [ -n "$node_ver" ] && [ -d "$NVM_DIR/versions/node/$node_ver/bin" ]; then
        export PATH="$PATH:$NVM_DIR/versions/node/$node_ver/bin"
      else
        latest_node_dir="$(find "$NVM_DIR/versions/node" -maxdepth 1 -mindepth 1 -type d -printf '%f\n' 2>/dev/null | sort -V | tail -n1)"
        if [ -n "$latest_node_dir" ] && [ -d "$NVM_DIR/versions/node/$latest_node_dir/bin" ]; then
          export PATH="$PATH:$NVM_DIR/versions/node/$latest_node_dir/bin"
        fi
      fi
    fi
    export NODE_PATH_SET=1
  fi

  # Corepack (habilita yarn/pnpm se existir)
  if command -v corepack >/dev/null 2>&1; then
    corepack enable >/dev/null 2>&1 || true
  fi
  echo "[agent_env] PATH=$PATH" 1>&2
  command -v go   >/dev/null 2>&1 && go version 1>&2 || echo "[agent_env] go not found" 1>&2
  command -v node >/dev/null 2>&1 && node -v   1>&2 || echo "[agent_env] node not found" 1>&2
  command -v npm  >/dev/null 2>&1 && npm -v    1>&2 || echo "[agent_env] npm not found" 1>&2
  command -v pnpm >/dev/null 2>&1 && pnpm -v   1>&2 || echo "[agent_env] pnpm not found" 1>&2

  return 0
}

if [[ "${BASH_SOURCE[0]:-}" == "${0}" ]]; then
  echo "envAI version: $(_version)"
  echo "Loading available tools..."
  __load_available_tools
  __show_org_files
elif [[ "${ENVAI_AUTOLOAD:-1}" == "1" ]]; then
  # Quando o script é "sourced", executa o bootstrap por padrão
  echo "envAI version: $(_version)"
  echo "Loading available tools..."
  __load_available_tools
fi


# End of file
