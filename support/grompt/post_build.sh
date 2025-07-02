#!/usr/bin/env bash


_ROOT_DIR="$(git rev-parse --show-toplevel)"

_CP_EXEC=$(cp -r "${_ROOT_DIR}/frontend/build" "${_ROOT_DIR}/internal/services/server" 2>/dev/null || true)

if [[ ! -d "${_ROOT_DIR}/internal/services/server/build" ]]; then
  echo "❌ O diretório de build não foi encontrado em '${_ROOT_DIR}/internal/services/server/build'."
  exit 1
else
  echo "✅ Build copiado com sucesso para '${_ROOT_DIR}/internal/services/server/build'."
fi

