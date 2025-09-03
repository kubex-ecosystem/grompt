#!/usr/bin/env bash

# LookAtni + Grompt Meta-Recursivo Refactor Loop
# "O ciclo mais virtuoso de desenvolvimento que ta rolando!" 🚀

set -o nounset  # Treat unset variables as an error
set -o errexit  # Exit immediately if a command exits with a non-zero status
set -o pipefail # Prevent errors in a pipeline from being masked
set -o errtrace # If a command fails, the shell will exit immediately
set -o functrace # If a function fails, the shell will exit immediately
shopt -s inherit_errexit # Inherit the errexit option in functions

IFS=$'\n\t'

cd /srv/apps/LIFE/KUBEX/grompt || exit 1

_temp_combined_prompt=""
_start_time=$(date +%s)
_end_time=$(date +%s)

# Coloquei _ em todas as VARS pra não colidir com nenhuma do sistema mesmo,
# que por convenção não usa esse prefixo, além de só usar uppercase, por
# isso todas também foram convertidas para lowercase.

# LookAtni Settings
_init_pattern='///'
_end_pattern='///'
_artifact_content=""

# LLM Settings
# Markers for replace
_max_tokens=16000
_gemini_api_key='AIzaSyAGVRdfCOiW5HZdp09Bbtf4cwqn0mfLUv8'

# Paths
_workspace_path="/srv/apps/LIFE/KUBEX"

_lookatni_path="${_workspace_path}/lookatni-file-markers"
_lookatni_bin="${_lookatni_path}/dist/lookatni"
_original_target_path="${_lookatni_path}/test-project"

_grompt_path="${_workspace_path}/grompt"
_grompt_bin="${_grompt_path}/dist/grompt_linux_amd64"
_output_target_path="${_grompt_path}/docs/prompt"
_refactored_output="${_output_target_path}/refactored/test-project"

_artifact_file="${_output_target_path}/test-project-artifact.md"
_prompt_file="${_output_target_path}/improvement-prompt.md"

_refactored_file="${_output_target_path}/refactored-project.md"
_refactored_output="${_output_target_path}/refactored/test-project"

cleanup() {
  # Não fazer o PC nem ngm perder tempo... rsrs
  if [[ -n "${_temp_combined_prompt:-}" && -f "${_temp_combined_prompt:-}" ]]; then
    rm -f "${_temp_combined_prompt:-}" || true
  fi

  # Limpa a trap do script pra evitar "gremlins"
  trap - EXIT

  return 0
}

trap cleanup EXIT

provision() {
  printf '%s\n' "🔧 Provisionando ambiente..."
  _temp_combined_prompt="$(mktemp -t combined_prompt.XXXXXX || echo "")"
  test -f "${_temp_combined_prompt:-}" || {
      printf '%s\n' "❌ Falha ao gerar prompt combinado!" >&2
      return 1
  }
  return 0
}

first_step() {
  printf '%s\n' "📦 PASSO 1: Gerando artefato do projeto..." >&2
  "${_lookatni_bin}" generate "${_original_target_path:-}" "${_artifact_file:-}"

  # Verificar se o artefato foi gerado
  test -f "$_artifact_file" && printf '%s\n' "✅ Artefato gerado: ${_artifact_file:-}" >&2 || return 1
  sleep 2

  return 0
}

second_step() {
  printf '%s\n' "🧠 PASSO 2: Gerando prompt profissional..." >&2

  _grompt_generated_prompt=$("${_grompt_bin:-}" generate \
      --provider gemini \
      --apikey "${_gemini_api_key:-${GEMINI_API_KEY:-}}" \
      --model 'gemini-2.0-flash' \
      --ideas 'Analise este projeto Go e identifique melhorias de código seguindo Go best practices' \
      --ideas 'Foque em: error handling, naming conventions, código idiomático, performance' \
      --ideas 'Mantenha a estrutura de arquivos LookAtni (markers //<ASCII[28]>/ filename /<ASCII[28]>//)' \
      --ideas 'O placeholder <ASCII[28]> representa o caractere ASCII 28 (File Separator - ) e deve ser IMPRESSO no resultado apresentado.' \
      --ideas 'Retorne o código refatorado completo com explicações das mudanças, sem título ou rodapé, mas explicações em comentários no próprio código.' \
      --max-tokens 10000 \
      --purpose 'code')

  test -n "${_grompt_generated_prompt:-}" || {
      printf '%s\n' "❌ Prompt falhou!" >&2
      return 1
  }

  printf '%s\n' "${_grompt_generated_prompt:-}" > "${_prompt_file:-}"
  test -f "${_prompt_file:-}" || {
    printf '%s\n' "❌ Falha ao salvar prompt!" >&2
    return 1
  }

  return 0
}

third_step() {
  printf '%s\n' "🤖 PASSO 3: Executando refatoração com IA..." >&2

  # # Criar prompt combinado
  _artifact_content="$(cat "${_artifact_file:-}" --show-nonprinting)"
  test -n "${_artifact_content}" || {
    printf '%s\n' "❌ Falha ao ler artefato!" >&2
    return 1
  }

  if [[ " ${#_artifact_content} " -lt $(( _max_tokens / 2 - 1000 )) ]]; then
      _combined_prompt="$(printf '%s\nTARGET CONTENT:\n%s\n%s\n%s\n' "$(cat "${_prompt_file:-}" --show-nonprinting)" '```plaintext' "${_artifact_content:-}" '```')"
  else
      _combined_prompt="$(printf '%s\n' "$(cat "${_prompt_file:-}" --show-nonprinting)")"
  fi

  # Agora salvo no temporário de verdade pra não incomodar caso dê algo errado.. hehe
  printf '%s\n' "${_combined_prompt:-}" > "${_temp_combined_prompt:-}"

  test -f "${_temp_combined_prompt:-}" || {
      printf '%s\n' "❌ Falha ao gerar prompt combinado!" >&2
      return 1
  }

  printf '%s\n' "✅ Prompt combinado gerado: ${_temp_combined_prompt:-}" >&2
  sleep 2

  # # Executar com Gemini (usando arquivo temporário para contornar limite de input)
  _grompt_ask="$("${_grompt_bin:-}" ask \
      --prompt "$(cat "${_temp_combined_prompt:-}" --show-nonprinting)" \
      --provider 'gemini' \
      --apikey "${_gemini_api_key:-${GEMINI_API_KEY:-}}" \
      --max-tokens 8000 && true || echo '')" # Só pra um ensure pica esse true...

  # Checa se tá preenchido
  test -n "${_grompt_ask:-}" || {
      printf '%s\n' "❌ Refatoração falhou!"
      return 1
  }

  # Preenche o arquivo
  printf '%s\n' "${_grompt_ask:-}" > "${_refactored_file:-}"

  # Checa se o arquivo foi preenchido
  test -f "${_refactored_file:-}"|| {
      printf '%s\n' "❌ Refatoração falhou!"
      return 1
  }

  # Tudo ok? Printa!
  printf '%s\n' "✅ Refatoração concluída: ${_refactored_file:-}"

  return 0
}

fourth_step() {
  # To inserindo o printf pra garantir que o caracter invisível será impresso sem
  # nenhuma espécie de expansão, etc....
  # O cat é bom manter com o --show-nonprinting pra garantir que o caracter invisível
  _exemplo="$(cat "${_refactored_file:-}" --show-nonprinting)"
  _exemplo="${_exemplo#*\`\`\`go}"
  _exemplo="${_exemplo//\/\/<ASCII\[28\]>\//$(printf "//\x1C/")}"
  _exemplo="${_exemplo//\/<ASCII\[28\]>\/\//$(printf "/\x1C//")}"

  printf '%b\n' "${_exemplo:-}" > "${_refactored_file:-}" || {
    printf '%s\n' "❌ Falha ao processar marcadores!"
    return 1
  }

  test -f "${_refactored_file:-}" || {
    printf '%s\n' "❌ Artefato refatorado não encontrado!"
    return 1
  }

  # Pra isso aqui passar falta o seguinte:
  # 1: Remover a última linha do arquivo $_refactored_file gerado.
  # 2: Substituir o caracter ␜ pelo  real

  if ! "${_lookatni_bin:-}" validate "${_refactored_file:-}"; then
    printf '%s\n' "❌ Validação falhou!"
    return 1
  fi

  printf '%s\n' "📁 PASSO 4: Extraindo projeto refatorado..."

  "${_lookatni_bin:-}" extract "${_refactored_file:-}" "${_refactored_output:-}" --overwrite --create-dirs
  test -d "${_refactored_output:-}" || {
    printf '%s\n' "❌ Extração falhou!"
    return 1
  }

  printf '%s\n' "✅ Projeto refatorado extraído: ${_refactored_output}"
  return 0
}

print_summary() {
  printf '\n%s\n' "🎉 META-RECURSIVIDADE COMPLETA!" >&2
  printf '%s\n' "================================" >&2
  printf '%s\n' "📂 Projeto original: ${_original_target_path}" >&2
  printf '%s\n' "📄 Artefato: ${_artifact_file}" >&2
  printf '%s\n' "🧠 Prompt: ${_prompt_file}" >&2
  printf '%s\n' "🤖 Refatorado: ${_refactored_file}" >&2
  printf '%s\n' "📁 Projeto final: ${_refactored_output}" >&2
  printf '\n%s\n' "🔥 BOOM! SÓ ALEGRIA E EVOLUÇÃO NO CICLO MAIS VIRTUOSO! 🚀" >&2
  printf '%s\n' "================================" >&2

  _end_time=$(date +%s)
  _duration=$(( _end_time - _start_time ))
  if [ "$_duration" -gt 60 ]; then
    printf '%s\n' "⏱️ Duração total: $(( _duration / 60 )) minutos e $(( _duration % 60 )) segundos" >&2
  else
    printf '%s\n' "⏱️ Duração total: ${_duration} segundos" >&2
  fi
}

main() {
  printf '\n%s\n' "=================================================" >&2
  printf '%s\n' "🚀 INICIANDO META-RECURSIVIDADE DO LOOKATNI + GROMPT!" >&2
  printf '%s\n' "=================================================" >&2

  provision || {
    printf '%s\n' "❌ Falha no provisionamento!" >&2
    exit 1
  }

  first_step || {
    printf '%s\n' "❌ Falha no passo 1!" >&2
    exit 1
  }
  second_step || {
    printf '%s\n' "❌ Falha no passo 2!" >&2
    exit 1
  }
  third_step || {
    printf '%s\n' "❌ Falha no passo 3!" >&2
    exit 1
  }
  fourth_step || {
    printf '%s\n' "❌ Falha no passo 4!" >&2
    exit 1
  }
  cleanup || {
    printf '%s\n' "❌ Falha na limpeza!" >&2
    exit 1
  }


}

main || {
  printf '%s\n' "❌ Falha na execução do script!" >&2
  exit 1
}

print_summary
