<!-- ---
title: Shellscript Craftsmanship Standards
version: 0.1.0
owner: kubex
audience: dev
languages: [pt-BR]
sources: ["support/string_utils.sh"]
assumptions: ["Linux moderno (bash >= 4.4)", "scripts seguem padrão de header/footer provido no repositório"]
--- -->

# Shellscript Craftsmanship Standards

Este documento descreve as normas mínimas de qualidade para scripts shell do monorepo Kubex e fornece o modelo de header e footer padrão (proteção de execução / exportação). Siga o header para definir se o script será usado como biblioteca (sourced) ou executável. Use o footer padrão para expor funções quando o script for importado. Testes devem forçar modo exec via cópia temporária e validar comportamento de sourcing sem modificar o arquivo original.

Conteúdo principal

- Objetivo
  - Padronizar headers e footers de scripts para garantir determinismo e proteção da lógica interna.
  - Facilitar reuse: permitir exportar funções quando for intencional e proteger execução direta quando necessário.

- Requisitos
  - Ambiente: Linux moderno com bash compatível (bash >= 4.4).
  - Ferramentas usuais: sed, awk, rev (ou perl fallback), git, tput.

- Convenções principais
  - Metadata no header: __secure_logic_* (versão, author, use_type, timestamps).
  - use_type: escolha explícita entre "lib" (sourced-only) e "exec" (executável-only). Definido no próprio arquivo para garantir proteção determinística.
  - Mensagens de erro → escreva para stderr.
  - Sem efeitos colaterais de módulo ao declarar funções (nenhuma execução ao importar salvo via export intencional no footer).

Modelo de header (copiar/colar e ajustar metadados)

```bash
# filepath: /projects/kubex/grompt/support/string_utils.sh (exemplo)
#!/usr/bin/env bash
# shellcheck disable=SC2015,SC1091

# Script Metadata
__secure_logic_version="1.0.0"
__secure_logic_date="$( date +%Y-%m-%d )"
__secure_logic_author="Rafael Mori"
# __secure_logic_use_type="exec"  # "exec" or "lib"
__secure_logic_use_type="lib"
__secure_logic_init_timestamp="$(date +%s)"
__secure_logic_elapsed_time=0

# Check if verbose mode is enabled
if [[ "${MYNAME_VERBOSE:-false}" == "true" ]]; then
  set -x  # Enable debugging
fi

IFS=$'\n\t'

_ROOT_DIR="$(git rev-parse --show-toplevel 2>/dev/null)"
_ROOT_DIR="${_ROOT_DIR:-$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)}"

declare -a _main_args=( "$@" )
```

- Notas sobre o header
  - Mantenha o comentário que indica a opção __secure_logic_use_type (para documentar intenção).
  - Não altere __secure_logic_use_type em runtime para escapar da proteção; o propósito é que o arquivo contenha a decisão final.
  - Use _ROOT_DIR via git rev-parse quando disponível, com fallback para caminho relativo.

Modelo de footer (padrão para controlar export / execução)

```bash
# filepath: /projects/kubex/grompt/support/string_utils.sh (exemplo)
__secure_logic_main() {
  local _ws_name
  _ws_name="$(__secure_logic_sourced_name)"
  local _ws_name_val
  _ws_name_val=$(eval "echo \${${_ws_name:-}}")
  if test "${_ws_name_val:-}" != "true"; then
    __main_functions "${_main_args[@]}"
    return $?
  else
    # shellcheck disable=SC2207
    local _functions_to_export=( $(compgen -A function | grep -v -e '^__') )
    for func in "${_functions_to_export[@]}"; do
      # shellcheck disable=SC2163
      export -f "$func" || {
        echo "Error: Could not export function '$func'." >&2
        return 1
      }
    done
    return 0
  fi
}

__secure_logic_main "${_main_args[@]}"
```

- Comportamento esperado do footer
  - Quando o script é executado diretamente (modo exec), __main_functions é chamado com os args originais.
  - Quando o script é sourceado e a variável de controle interna indica sourced=true, o footer exporta (export -f) as funções públicas (não iniciadas por __).
  - O nome da variável de controle é gerada por __secure_logic_sourced_name e depende do caminho do script. Não sobreponha esse mecanismo.

How to run / Repro

- Testar execução (modo exec) sem modificar arquivo original:
  - Crie cópia temporária substituindo apenas a linha de use_type por exec; execute a cópia.
  - Exemplo rápido:
    - awk '...sub(...,"__secure_logic_use_type=\"exec\"")...' support/foo.sh > /tmp/foo_exec.sh
    - chmod +x /tmp/foo_exec.sh
    - /tmp/foo_exec.sh toLowerCase "HeLLo"
- Testar sourcing (modo lib) sem alterar arquivo:
  - export __secure_logic_use_type="lib"
  - source support/foo.sh
  - chame funções diretamente: toLowerCase "HeLLo"

Riscos & Mitigações

- Risco: Override de __secure_logic_use_type externamente (tentativa de burlar proteção).
  - Mitigação: Decisão do modo deve residir no arquivo (comentário e valor); documente políticas de revisão de MR para não alterar esse bloco.
- Risco: export -f não disponível em shells diferentes do bash.
  - Mitigação: Exigir bash na shebang e documentar requisito (Linux moderno, bash >= 4.4).
- Risco: efeitos colaterais ao sourcear (execuções não intencionais).
  - Mitigação: evitar comandos top-level fora de funções; usar __first para prevenir execução indevida.

Práticas recomendadas rápidas

- Sempre documente __secure_logic_use_type no header (comentário e valor).
- Evite dependências de estado global; prefira receber parâmetros em funções.
- Mensagens de erro para stderr, dados para stdout.
- Escreva testes pequenos que validem ambos os modos: exec (cópia temporária) e source (export e source).
- Use shellcheck nas funções (não desabilite regras sem justificativa clara).

Próximos passos (máx. 3)

1. Aplicar header/footer a scripts novos e migrar scripts críticos para este padrão.
2. Adicionar job CI que executa testes básicos (exec + source) para cada script em support/.
3. Documentar no CONTRIBUTING.md a convenção (exemplo de header/footer + checklist de revisão).

Changelog

- 0.1.0 — Documento inicial com header/footer padrão e instruções de teste.
