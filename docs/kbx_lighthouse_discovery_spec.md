//  / .ops/kbx_discovery.sh /  //
#!/usr/bin/env bash
# KBX Lighthouse discovery — guia, não juiz.
# Descobre e normaliza metadados de um repo Kubex-compatível.
# Saídas: JSON (default), pretty JSON, ou exports ENV.
# Níveis de adoção: L0 (heurístico), L1 (breadcrumb/.kbx ou hints), L2 (lighthouse/manifest), L3 (full: manifest + docs + workflows).

set -euo pipefail

FORMAT="json"      # json | pretty | env
STRICT=0            # --strict => exit!=0 se nada encontrado
PRINT_PATH=0        # --path-only
PRINT_LEVEL=0       # --level-only
QUIET=0             # --quiet

usage() {
  cat <<USAGE
KBX Discovery
Usage: $0 [--format json|pretty|env] [--path-only] [--level-only] [--strict] [--quiet]
USAGE
}

log(){ [[ $QUIET -eq 1 ]] || echo "[$1] $2" >&2; }
exists(){ command -v "$1" >/dev/null 2>&1; }

# Caminhos preferenciais (ordem)
_kbx_manifest_path(){
  local p
  for p in \
    internal/module/info/manifest.json \
    .kubex/manifest.toml \
    manifest.json \
    .kbx; do
    [[ -f "$p" ]] && { echo "$p"; return 0; }
  done
  return 1
}

# Heurísticas úteis
_has(){ [[ -e "$1" ]]; }
_hasdir(){ [[ -d "$1" ]]; }
_first_hint_file(){
  # Procura hints // kbx: em .go
  grep -RIl '^[[:space:]]*//[[:space:]]*kbx:' -- '*.go' 2>/dev/null | head -n1 || true
}
_module_name_from_go(){
  if [[ -f go.mod ]]; then awk '/^module /{print $2; exit}' go.mod; fi
}
_basename(){ basename "$(pwd)"; }

# Normalizadores (cada um imprime JSON canônico)
_from_json_manifest(){
  local path="$1"
  jq -n --argfile m "$path" '($m) as $x | {
    name: ($x.name // ""),
    version: ($x.version // ""),
    description: ($x.description // ""),
    bin: ($x.bin // ($x.name // "")),
    repo: ($x.repo // ""),
    license: ($x.license // ""),
    capabilities: ($x.capabilities // []),
    build: {
      supports: {
        os:   ($x.build.supports.os   // ["linux"]),
        arch: ($x.build.supports.arch // ["amd64"]) },
      upxDefault: ($x.build.upxDefault // true)
    },
    entrypoints: ($x.entrypoints // {})
  }'
}

_from_toml_manifest(){
  local path="$1"
  python3 - "$path" <<'PY'
import sys, json
try:
    import tomllib  # py3.11+
except Exception:
    print('{}')
    sys.exit(0)

p=sys.argv[1]
with open(p,'rb') as f:
    t=tomllib.load(f)
# mapear campos comuns
x={
  'name': t.get('name',''),
  'version': t.get('version',''),
  'description': t.get('description',''),
  'bin': t.get('bin', t.get('name','')),
  'repo': t.get('repo',''),
  'license': t.get('license',''),
  'capabilities': t.get('capabilities',[]) or [],
  'build': {
    'supports': {
      'os':   (t.get('build',{}).get('supports',{}).get('os') or ['linux']),
      'arch': (t.get('build',{}).get('supports',{}).get('arch') or ['amd64'])
    },
    'upxDefault': t.get('build',{}).get('upxDefault', True)
  },
  'entrypoints': t.get('entrypoints',{}) or {}
}
print(json.dumps(x))
PY
}

_from_kbx_breadcrumb(){
  local path="$1"
  python3 - "$path" <<'PY'
import sys, json
s=open(sys.argv[1],'r',encoding='utf-8',errors='replace').read().strip()
# formato: key=val;key=val; cap pode ser lista separada por vírgula
pairs=[p for p in s.split(';') if p.strip()]
d={}
for p in pairs:
    if '=' in p:
        k,v=p.split('=',1)
        k=k.strip(); v=v.strip().strip('"\'')
        d[k]=v
name=d.get('name','')
bin_=d.get('bin', name)
cap = d.get('cap','')
os  = d.get('os','linux').split(',') if d.get('os') else ['linux']
arch= d.get('arch','amd64').split(',') if d.get('arch') else ['amd64']
obj={
  'name': name,
  'version': d.get('version',''),
  'description': d.get('description',''),
  'bin': bin_,
  'repo': d.get('repo',''),
  'license': d.get('license',''),
  'capabilities': [c.strip() for c in cap.split(',') if c.strip()] if cap else [],
  'build': {
    'supports': {'os': os, 'arch': arch},
    'upxDefault': d.get('upxDefault','true').lower()!='false'
  },
  'entrypoints': {}
}
print(json.dumps(obj))
PY
}

_from_hints(){
  local file="$1"
  python3 - "$file" <<'PY'
import sys, re, json
pat=re.compile(r"^[ \t]*//[ \t]*kbx:(.*)$")
kv=re.compile(r"(\w+)\s*=\s*(?:\"([^\"]*)\"|'([^']*)'|([^\s]+))")
content=open(sys.argv[1],'r',encoding='utf-8',errors='replace').read().splitlines()
d={}
for line in content:
    m=pat.match(line)
    if not m: continue
    rest=m.group(1)
    for k in kv.finditer(rest):
        key=k.group(1)
        val=k.group(2) or k.group(3) or k.group(4) or ''
        d[key]=val
name=d.get('name','') or ''
obj={
  'name': name,
  'version': d.get('version',''),
  'description': d.get('description',''),
  'bin': d.get('bin', name),
  'repo': d.get('repo',''),
  'license': d.get('license',''),
  'capabilities': [c.strip() for c in d.get('cap','').split(',') if c.strip()],
  'build': {
    'supports': {
      'os':  [x.strip() for x in (d.get('os','linux') or 'linux').split(',') if x.strip()],
      'arch':[x.strip() for x in (d.get('arch','amd64') or 'amd64').split(',') if x.strip()]
    },
    'upxDefault': (d.get('upxDefault','true').lower()!='false')
  },
  'entrypoints': {
    'buildDev': d.get('entry.buildDev',''),
    'runHelp': d.get('entry.runHelp','')
  }
}
print(json.dumps(obj))
PY
}

_from_heuristic(){
  local mod name bin
  mod=$(_module_name_from_go || true)
  if [[ -n "${mod:-}" ]]; then name="${mod##*/}"; else name="$(_basename)"; fi
  bin="$name"
  local cap=(); [[ -f main.go || -d cmd ]] && cap+=("cli")
  [[ -f package.json ]] && cap+=("frontend")
  # imprimir JSON mínimo
  python3 - <<PY
import json
print(json.dumps({
  'name': '${name}',
  'version': '',
  'description': '',
  'bin': '${bin}',
  'repo': '',
  'license': '',
  'capabilities': ${cap[@]+['$(printf "%s','" "${cap[@]}" | sed "s/,'$//")']},
  'build': {'supports': {'os':['linux'], 'arch':['amd64']}, 'upxDefault': True},
  'entrypoints': {}
}))
PY
}

# Argumentos
while [[ $# -gt 0 ]]; do
  case "$1" in
    --format) FORMAT="${2:-json}"; shift 2;;
    --path-only) PRINT_PATH=1; shift;;
    --level-only) PRINT_LEVEL=1; shift;;
    --strict) STRICT=1; shift;;
    --quiet) QUIET=1; shift;;
    -h|--help) usage; exit 0;;
    *) log WARN "arg ignorado: $1"; shift;;
  esac
done

# Descoberta
mp=$(_kbx_manifest_path || true)
json_path=""; toml_path=""; kbx_path=""; hints_file=""; source=""; level="L0"

if [[ -n "${mp:-}" ]]; then
  case "$mp" in
    *.json) json_path="$mp"; source="lighthouse-json";;
    *.toml) toml_path="$mp"; source="lighthouse-toml";;
    .kbx)   kbx_path="$mp"; source="breadcrumb";;
  esac
else
  hints_file="$(_first_hint_file || true)"
  if [[ -n "${hints_file:-}" ]]; then source="hints"; else source="heuristic"; fi
fi

# Nível
if [[ -n "$json_path" || -n "$toml_path" ]]; then
  # L2 por padrão; L3 se docs + workflows
  if _hasdir docs/architecture && _hasdir .github/workflows; then level="L3"; else level="L2"; fi
elif [[ -n "$kbx_path" || -n "$hints_file" ]]; then
  level="L1"
else
  level="L0"
fi

# Saídas rápidas
if [[ $PRINT_PATH -eq 1 ]]; then
  echo "${mp:-}"
  exit 0
fi
if [[ $PRINT_LEVEL -eq 1 ]]; then
  echo "$level"
  exit 0
fi

# Monta JSON canônico
data="{}"
case "$source" in
  lighthouse-json)
    data=$(_from_json_manifest "$json_path") ;;
  lighthouse-toml)
    data=$(_from_toml_manifest "$toml_path") ;;
  breadcrumb)
    data=$(_from_kbx_breadcrumb "$kbx_path") ;;
  hints)
    data=$(_from_hints "$hints_file") ;;
  heuristic)
    data=$(_from_heuristic) ;;
esac

# Enriquecer com metadados
if exists jq; then
  data=$(jq -n --arg src "$source" --arg lvl "$level" \
             --arg jp "$json_path" --arg tp "$toml_path" --arg kp "$kbx_path" --arg hf "$hints_file" \
             --arg has_docs "$( _hasdir docs/architecture && echo true || echo false )" \
             --arg has_wf   "$( _hasdir .github/workflows && echo true || echo false )" \
             --arg repo_dir "$(pwd)" \
             --argjson base "$data" '
    $base + {
      _meta: {
        source: $src,
        level: $lvl,
        paths: { json:$jp, toml:$tp, kbx:$kp, hints:$hf },
        hasDocs: ($has_docs=="true"),
        hasWorkflows: ($has_wf=="true"),
        repoDir: $repo_dir
      }
    }')
fi

# Saída no formato pedido
case "$FORMAT" in
  json)   echo "$data" ;;
  pretty) echo "$data" | (exists jq && jq . || cat) ;;
  env)
    # Exporta KBX_* para shells/CI
    if exists jq; then
      eval $(echo "$data" | jq -r '
        "export KBX_NAME=\"" + (.name // "") + "\"\n" +
        "export KBX_VERSION=\"" + (.version // "") + "\"\n" +
        "export KBX_BIN=\"" + (.bin // "") + "\"\n" +
        "export KBX_OS=\"" + ((.build.supports.os // ["linux"]) | join(",")) + "\"\n" +
        "export KBX_ARCH=\"" + ((.build.supports.arch // ["amd64"]) | join(",")) + "\"\n" +
        "export KBX_LEVEL=\"" + (._meta.level // "L0") + "\"\n" +
        "export KBX_SOURCE=\"" + (._meta.source // "") + "\"\n"
      ')
    fi
    env | grep '^KBX_' || true
    ;;
  *) usage; exit 2;;
fi

# Strict mode: falha se heurístico e pediram estrito
if [[ $STRICT -eq 1 && "$level" == "L0" ]]; then
  log ERROR "nenhuma pista KBX encontrada (L0)."; exit 3
fi

exit 0


//  / docs/architecture/LIGHTHOUSE.md /  //
# KBX Lighthouse — guia, não juiz

> **Ideia:** o *Lighthouse* é um **farol**. Orienta e compõe quando presente; não bloqueia quando ausente. Repositórios pequenos continuam pequenos.

## Níveis de adoção (opt‑in, com degradação elegante)
- **L0 — Heurístico:** sem arquivo KBX. Ferramentas deduzem por layout (`go.mod`, `cmd/`, `main.go`, `package.json`).
- **L1 — Breadcrumb:** existe `.kbx` (1 linha) **ou** *hints* em comentário (`// kbx: ...`).
- **L2 — Lighthouse:** manifest dedicado (ex.: `internal/module/info/manifest.json` ou `.kubex/manifest.toml`).
- **L3 — Full:** L2 **+** docs (`docs/architecture/`) **+** workflows (`.github/workflows/`).

> Nada obriga subir de nível. Um *smart‑contract* pode viver feliz em L1 para sempre.

## Ordem de descoberta (fallback)
1. `internal/module/info/manifest.json` (JSON — preferido)
2. `.kubex/manifest.toml` (TOML)
3. `manifest.json` (fallback)
4. `.kbx` (breadcrumb de 1 linha)
5. *Hints* de comentário `// kbx:` em arquivos `.go`
6. **Heurística** (L0) — nome do módulo por `go.mod`, *bin* = nome do repo, `cli` se `cmd/` ou `main.go`, `frontend` se `package.json`

## Formatos suportados
### Lighthouse JSON (leve)
```json
{
  "kind": "kbx.lighthouse",
  "name": "kbxctl",
  "version": "0.3.0",
  "bin": "kbxctl",
  "description": "Kubex Control CLI",
  "capabilities": ["cli"],
  "build": {"supports": {"os": ["linux","darwin"], "arch": ["amd64","arm64"]}, "upxDefault": true},
  "entrypoints": {"buildDev": "make build-dev", "runHelp": "./dist/kbxctl --help"}
}
```

### Breadcrumb `.kbx` (ultramini)
```
kbx>=0.3;name=kbxctl;bin=kbxctl;profile=light;cap=cli;os=linux,darwin;arch=amd64,arm64
```

### Hints em comentário (para repositórios microscópicos)
```go
// kbx: compat=">=0.3" name="chaincode-demo" bin="ccdemo" profile="chaincode" cap="cli"
// kbx: entry.buildDev="go build ./..." entry.runHelp="./ccdemo --help"
package main
```

## Perfil *smart‑contract* (3 arquivos)
```
go.mod
main.go        // com hints // kbx: ...
chaincode.go
```
> Opcional: `.kbx` ou `Makefile` mini para DX local.

## Comportamento das ferramentas (opt‑in)
- **kbxctl/GoBE/agents** usam o Lighthouse **se existir**; caso contrário, heurística.
- Geradores de doc (README/BUILD_ARTIFACTS) só atuam quando há farol/breadcrumb/hints.
- CI pode validar o manifest **apenas** quando presente.

## CLI utilitário: `.ops/kbx_discovery.sh`
- **JSON (default):** `bash .ops/kbx_discovery.sh`
- **Bonito:** `bash .ops/kbx_discovery.sh --format pretty`
- **ENV:** `eval "$(bash .ops/kbx_discovery.sh --format env)"` (exporta `KBX_*`)
- **Caminho do manifest:** `bash .ops/kbx_discovery.sh --path-only`
- **Nível:** `bash .ops/kbx_discovery.sh --level-only`
- **Estrito:** `bash .ops/kbx_discovery.sh --strict` (falha se L0)

## Por que é assim?
- **Farol, não dogma:** convive bem com repositórios mínimos (ex.: chaincode em 3 arquivos).
- **Compatibilidade progressiva:** qualquer projeto externo pode ser “KBX‑compatível” só deixando `.kbx` ou `// kbx:`.
- **Sem atrito:** teu fluxo de build/Make continua igual; o farol apenas melhora *descoberta e composição* no ecossistema.



//  / .ops/kbx_apply_manifest.sh /  //
#!/usr/bin/env bash
# KBX: aplicar/ler manifest com degradação elegante.
# Usa .ops/kbx_discovery.sh se existir; caso contrário, faz fallback por jq.
set -euo pipefail
shopt -s nullglob

kbx::log(){ printf '[%s] %s
' "$1" "$2" >&2; }
kbx::repo_root(){ git rev-parse --show-toplevel 2>/dev/null || pwd; }

_ROOT_DIR=${_ROOT_DIR:-"$(kbx::repo_root)"}

kbx::load_env(){
  if [[ -x "${_ROOT_DIR}/.ops/kbx_discovery.sh" ]]; then
    # Preferir o farol
    eval "$(${_ROOT_DIR}/.ops/kbx_discovery.sh --format env)"
  else
    # Fallback: caminhos comuns
    local M="";
    for p in \
      "${_ROOT_DIR}/internal/module/info/manifest.json" \
      "${_ROOT_DIR}/manifest.json"; do
      [[ -f "$p" ]] && { M="$p"; break; }
    done
    [[ -n "$M" ]] || { kbx::log WARN "manifest não encontrado"; return 0; }
    export KBX_NAME=$(jq -r '.name // ""' "$M")
    export KBX_VERSION=$(jq -r '.version // ""' "$M")
    export KBX_BIN=$(jq -r '.bin // .name // ""' "$M")
    export KBX_OS=$(jq -r '(.build.supports.os // .platforms // ["linux"]) | join(",")' "$M")
    export KBX_ARCH=$(jq -r '(.build.supports.arch // ["amd64"]) | join(",")' "$M")
  fi
}

kbx::sed_i(){ # sed -i compat (Linux/macOS)
  if sed --version >/dev/null 2>&1; then sed -i "$@"; else sed -i '' "$@"; fi
}

apply_manifest(){
  kbx::load_env
  # Vars públicas usadas pelo fluxo atual
  export _ROOT_DIR
  export _APP_NAME="${_APP_NAME:-${KBX_BIN:-${KBX_NAME:-$(basename "${_ROOT_DIR}")}}}"
  export _BINARY_NAME="${_BINARY_NAME:-${_APP_NAME}}"
  export _PROJECT_NAME="${_PROJECT_NAME:-${KBX_NAME:-${_APP_NAME}}}"
  export _VERSION="${_VERSION:-${KBX_VERSION:-v0.0.0}}"
  export _DESCRIPTION="${_DESCRIPTION:-No description provided.}"
  export _OWNER="${_OWNER:-rafa-mori}"; _OWNER="${_OWNER,,}"
  export _AUTHOR="${_AUTHOR:-Rafa Mori}"
  export _LICENSE="${_LICENSE:-MIT}"
  # repo/repository compat
  local M SUB="internal/module/info/manifest.json"; M="${_ROOT_DIR}/${SUB}"
  export _REPOSITORY="${_REPOSITORY:-$(jq -r '.repo // .repository // empty' "$M" 2>/dev/null || true)}"
  export _PRIVATE_REPOSITORY="${_PRIVATE_REPOSITORY:-$(jq -r '.private // false' "$M" 2>/dev/null || echo false)}"
  # Go version resiliente
  if [[ -f "${_ROOT_DIR}/go.mod" ]]; then
    export _VERSION_GO=$(awk '/^go /{print $2; exit}' "${_ROOT_DIR}/go.mod" || true)
  else
    export _VERSION_GO=""
  fi
  # OS/ARCH (já vem de KBX_* se discovery rodou)
  export _PLATFORMS_SUPPORTED
  if [[ -n "${KBX_OS:-}" ]]; then
    _PLATFORMS_SUPPORTED="${KBX_OS,,}"
  else
    _PLATFORMS_SUPPORTED=$(jq -r '(.platforms // .build.supports.os // ["linux"]) | join(", ")' "$M" 2>/dev/null || echo "linux, darwin")
    _PLATFORMS_SUPPORTED="${_PLATFORMS_SUPPORTED,,}"
  fi
}

# renomeador seguro: usa find e evita .git/dist/node_modules/vendor
change_project_name(){
  local _old_bin_name="${1:-gobe}"; shift || true
  local _new_bin_name="${1:-${_BINARY_NAME:-${_APP_NAME}}}"; shift || true
  local DRY=${DRY_RUN:-1}

  mkdir -p "${_ROOT_DIR}/bkp"
  tar --exclude='bkp' --exclude='*.tar.gz' --exclude='.git' --exclude='dist' --exclude='node_modules' \
      -czf "${_ROOT_DIR}/bkp/$(date +%Y%m%d_%H%M%S)_backup.tar.gz" -C "${_ROOT_DIR}" .

  # Renomear arquivos que contém o nome antigo
  mapfile -t files_to_mv < <(find "${_ROOT_DIR}" \
    -path "${_ROOT_DIR}/.git" -prune -o \
    -path "${_ROOT_DIR}/dist" -prune -o \
    -path "${_ROOT_DIR}/node_modules" -prune -o \
    -type f -name "*${_old_bin_name}*" -print)
  for f in "${files_to_mv[@]}"; do
    nf="${f//${_old_bin_name}/${_new_bin_name}}"
    if [[ "$DRY" -eq 1 ]]; then kbx::log INFO "mv (dry) $f -> $nf"; else mv "$f" "$nf"; fi
  done

  # Substituir conteúdo em arquivos de interesse
  mapfile -t files_to_sed < <(find "${_ROOT_DIR}" \
    -path "${_ROOT_DIR}/.git" -prune -o \
    -path "${_ROOT_DIR}/dist" -prune -o \
    -path "${_ROOT_DIR}/node_modules" -prune -o \
    -type f \( -name "*.go" -o -name "*.md" -o -name "go.mod" \) -print)
  for f in "${files_to_sed[@]}"; do
    if [[ "$DRY" -eq 1 ]]; then kbx::log INFO "sed (dry) $f"; else kbx::sed_i "s/${_old_bin_name}/${_new_bin_name}/g" "$f"; fi
  done

  # Ajustar go mod
  if [[ -f "${_ROOT_DIR}/go.mod" && "$DRY" -ne 1 ]]; then (cd "${_ROOT_DIR}" && go mod tidy || true); fi
}

export -f apply_manifest change_project_name


//  / docs/architecture/LIGHTHOUSE.md (appendix: quick usage) /  //
## Appendix — Quick usage with scripts
```bash
# carregar variáveis a partir do farol
source .ops/kbx_apply_manifest.sh
apply_manifest

# renomear projeto (dry-run por padrão)
DRY_RUN=1 change_project_name gobe "$__APP_NAME"
# efetivar
DRY_RUN=0 change_project_name gobe "$__APP_NAME"
```

