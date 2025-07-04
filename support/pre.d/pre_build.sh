#!/usr/bin/env bash
# shellcheck disable=SC2065,SC2015

set -o nounset  # Treat unset variables as an error
set -o errexit  # Exit immediately if a command exits with a non-zero status
set -o pipefail # Prevent errors in a pipeline from being masked
set -o errtrace # If a command fails, the shell will exit immediately
set -o functrace # If a function fails, the shell will exit immediately
shopt -s inherit_errexit # Inherit the errexit option in functions


build_frontend() {
  _ROOT_DIR="$(git rev-parse --show-toplevel)"


  if test -d "${_ROOT_DIR}/frontend"; then
    cd "${_ROOT_DIR}/frontend" || exit 1

    npm i -f --silent || {
      echo "❌ Falha ao instalar dependências do frontend. Verifique o log acima."
      exit 1
    }

    npm run build --silent || {
      echo "❌ Falha ao construir o frontend. Verifique o log acima."
      exit 1
    }
  else
    echo "❌ O diretório 'frontend' não foi encontrado em '${_ROOT_DIR}/frontend'."
    exit 1
  fi

  _CP_EXEC=$(cp -r "${_ROOT_DIR}/frontend/build" "${_ROOT_DIR}/internal/services/server" 2>/dev/null || true)

  if [[ ! -d "${_ROOT_DIR}/internal/services/server/build" ]]; then
    echo "❌ O diretório de build não foi encontrado em '${_ROOT_DIR}/internal/services/server/build'."
    exit 1
  else
    echo "✅ Build copiado com sucesso para '${_ROOT_DIR}/internal/services/server/build'."
  fi
}

build_frontend "$@" 2>/dev/null || {
  echo "❌ Ocorreu um erro durante o processo de build do frontend. Verifique o log acima."
  exit 1
}