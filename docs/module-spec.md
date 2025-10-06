# Module Specification

```go
// Module is an interface that defines the methods a module must implement.
type Module interface {
  Alias() string
  ShortDescription() string
  LongDescription() string
  Usage() string
  Examples() []string
  Active() bool
  Module() string
  Execute() error
  Command() *cobra.Command
  SetParentCmdName(rtCmd string)
  concatenateExamples() string
}
```

```go
// Simple Generic Module Concrete Case
type ModuleConcrete struct {
 parentCmdName string
 printBanner   bool
}

/*
  Abaixo estão os métodos que representam wrappers
  de dados básicos do módulo, enquanto CLI e Lib.

  Esses métodos são obrigatórios para qualquer módulo.
 */
func (m *ModuleConcrete) Alias() string {
 return "module"
}
func (m *ModuleConcrete) ShortDescription() string {
 return "ModuleConcrete short description."
}
func (m *ModuleConcrete) LongDescription() string {
 return `ModuleConcrete long description.`
}
func (m *ModuleConcrete) Usage() string {
 return "module [command] [args]"
}
func (m *ModuleConcrete) Examples() []string {
 return []string{
  "module doctor <subcommand>",
  "module platform <subcommand>",
 }
}
func (m *ModuleConcrete) Active() bool {
 return true
}
func (m *ModuleConcrete) Module() string {
 return "module"
}
func (m *ModuleConcrete) Execute() error {
 return m.Command().Execute()
}

func (m *ModuleConcrete) Command() *cobra.Command {
/*
  Todo módulo, pode opcionalmente usar o `Annotations`
  como uma camada de metadados para o comando.
  Ele contém dois argumentos:
    - Descriptions: []string{}
      Ele irá lidar com os Short e Long descriptions
      contemplando a abstração que há trazendo estilização
      e formatação.

    - hideBanner: bool
      Se verdadeiro, oculta a exibição do banner do módulo e
      torna o comportamento do Cobra em relação ao Short e Long
      descriptions nativo novamente.
 */

 rootCmd := &cobra.Command{
  Use:   m.Module(),
  Alias: m.Alias(),
  Short: m.ShortDescription(),
  Long:  m.LongDescription(),
  // Aqui há um tratamento quebrando as linhas do
  // array de exemplos na exibição. Permitindo
  // que seja veículado vários exemplos para o mesmo comando
  // de forma organizada.
  Example: m.concatenateExamples(),

  Annotations: m.GetDescriptions(
    []string{
      m.ShortDescription(),
      m.LongDescription(),
    }, m.printBanner,
  ),
  RunE: func(cmd *cobra.Command, args []string) error {
   return m.Execute()
  },
 }

// Aqui você pode adicionar subcomandos ao comando raiz
// Ex: rootCmd.AddCommand(subCmd1, subCmd2)

// Após inserir os comandos, é necessário acionar o
// registro do template e cores para o Cobra.
// Isso faz com que o CLI tenha uma aparência
// mais amigável e estilizada, além de padronizada
// para todos os módulos do ecossistema.
 setUsageDefinition(rtCmd)

 // Esse loop existe para garantir que os padrões de
 // estilo também sejam aplicados aos subcomandos.
 for _, c := range rtCmd.Commands() {
  setUsageDefinition(c)
  if !strings.Contains(strings.Join(os.Args, " "), c.Use) {
   if c.Short == "" {
    c.Short = c.Annotations["description"]
   }
  }
 }

// Retorna o comando completo.
 return rootCmd
}
```
/// jcm.ltk ///
/// .codex/config.toml ///
[settings]
model_reasoning_summary = "concise"
user_intent_summary     = "detailed"
project_summary         = "detailed"
default_max_tokens_out  = 700
fail_fast               = true
temperature             = 0.1

[manifests]
# cada repo terá estes princípios leves sempre no contexto
paths = [
  "docs/brand/**/kubex_brand_visual_spec*.md",
  "docs/architecture/**/Resilience-by-Design*.md",
  "docs/governance/**/KUBEX_MANIFESTO*.md"
]
max_tokens = 300
update_strategy = "incremental"
priority = "high"

[context]
scopes = [
  { name = "objective",   max_tokens = 120, summarization = "bullet_goals" },
  { name = "constraints", max_tokens = 120, summarization = "hard_rules"   },
  { name = "facts",       max_tokens = 200, summarization = "code_facts"   },
  { name = "delta",       max_tokens = 180, summarization = "git_diff"     },
  { name = "plan",        max_tokens = 160, summarization = "step_plan"    },
  { name = "exec",        max_tokens = 120, summarization = "checklist"    }
]
eviction_order = ["exec", "plan", "delta"]

[output]
style = "succinct_actionable"
postamble = "✅ Valide: build/test passam e objetivo foi atendido."
# foco em produto: não deixar “relatórios de tarefa” no repo
generate_task_summaries = false
allow_repo_writing_of_reports = false

# ajuste a confiança conforme o repo
[projects]
# exemplo: substitua <REPO_NAME> ou remova se preferir deteção automática
"/<ABS_PATH>/<REPO_NAME>" = { trust_level = "trusted" }

[guards]
dangerous_patterns    = ["rm -rf", "dd if=", "mkfs", "drop table", "helm uninstall --all"]
require_confirmation  = true
auto_dry_run          = true

[behaviors]
goal_driven_context   = true
prefer_recent_diffs   = true
refuse_when_unsafe    = true

[summarizers]
bullet_goals = """
- 3–5 bullets com verbo imperativo + artefato esperado.
"""
hard_rules = """
- Restrições duras (build tags, compatibilidade, lint), sem opinião estética.
"""
code_facts = """
- Fatos verificáveis do código (assinaturas, tipos, caminhos), sem interpretação.
"""
git_diff = """
- Mudanças por arquivo; +funções, -removidas, contratos alterados; destaque breaking.
"""
step_plan = """
- 3–6 passos atômicos, cada um com checagem objetiva.
"""
checklist = """
- Checklist [ ] executável.
"""
/// .kubex/config.toml ///
Este arquivo seria a fonte da verdade para o internal/config/loader.go
Mantendo a estrutura que definimos anteriormente.
[settings]
model_reasoning_summary = "concise"

[manifests]
paths = [
".kubex/docs/brand/kubex_brand_visual_spec.v0.0.1.md",
".kubex/docs/architecture/Resilience-by-Design.md",
".kubex/docs/governance/KUBEX_MANIFESTO.md"
]

... resto da configuração ...
/// .kubex/docs/governance/AI_SQUAD_ROLES.md ///
Papéis da AI Squad — Ecossistema Kubex
Mantendo o documento que define os papéis como parte da configuração base.

// (e assim por diante para os outros manifestos) //
/// .vscode/keybindings.json ///
[
  { "key": "ctrl+alt+o", "command": "codex.ide.runObjective" },
  { "key": "ctrl+alt+p", "command": "codex.ide.planFromChanges" },
  { "key": "ctrl+alt+d", "command": "codex.ide.applyDiffSafely" },
  { "key": "ctrl+alt+t", "command": "codex.ide.toggleTrace" }
]
/// .vscode/settings.json ///
{
  "codex.ide.manifests.paths": [
    "docs/brand/**/kubex_brand_visual_spec*.md",
    "docs/architecture/**/Resilience-by-Design*.md",
    "docs/governance/**/KUBEX_MANIFESTO*.md"
  ],
  "codex.ide.manifests.maxTokens": 300,
  "codex.ide.manifests.updateStrategy": "incremental",
  "codex.ide.output.style": "succinct_actionable",
  "codex.ide.output.postamble": "✅ Valide: build/test passam e objetivo foi atendido.",
  "codex.ide.context.goalDriven": true,
  "codex.ide.context.evictionOrder": ["exec", "plan", "delta"],
  "codex.ide.scopes": [
    { "name": "objective",   "maxTokens": 120, "summarization": "bullet_goals" },
    { "name": "constraints", "maxTokens": 120, "summarization": "hard_rules"   },
    { "name": "facts",       "maxTokens": 200, "summarization": "code_facts"   },
    { "name": "delta",       "maxTokens": 180, "summarization": "git_diff"     },
    { "name": "plan",        "maxTokens": 160, "summarization": "step_plan"    },
    { "name": "exec",        "maxTokens": 120, "summarization": "checklist"    }
  ],
  "codex.ide.summarizers": {
    "bullet_goals": "- 3–5 bullets com verbo imperativo + artefato esperado.",
    "hard_rules": "- Restrições duras (build tags, compatibilidade, lint), sem opinião estética.",
    "code_facts": "- Fatos verificáveis do código (assinaturas, tipos, caminhos), sem interpretação.",
    "git_diff": "- Mudanças por arquivo; +funções, -removidas, contratos alterados; destaque breaking.",
    "step_plan": "- 3–6 passos atômicos, cada um com checagem objetiva.",
    "checklist": "- Checklist [ ] executável."
  },
  "codex.ide.projects": {
    "/<ABS_PATH>/<REPO_NAME>": {
      "trustLevel": "trusted",
      "languageBias": ["go","ts","tsx","bash","make","md","rust","wasm"],
      "rules": [
        "Sem relatórios de tarefa no repo; apenas artefatos de produto.",
        "Tests e linters como fonte de verdade.",
        "Padrões Kubex e Brand visual consistentes."
      ]
    }
  },
  "codex.ide.guards.dangerousPatterns": ["rm -rf","dd if=","mkfs","drop table","helm uninstall --all"],
  "codex.ide.guards.requireConfirmation": true,
  "codex.ide.guards.autoDryRun": true,

  "editor.formatOnSave": true,
  "files.trimFinalNewlines": true,

  "github.copilot.inlineSuggest.enable": true,
  "github.copilot.editor.enableAutoCompletions": true,
  "github.copilot.enable": { "*": true, "markdown": false, "plaintext": false, "yaml": false },
  "github.copilot.advanced": { "excludedFolders": ["docs/brand","docs/architecture","docs/governance",".codex"] }
}
/// .vscode/tasks.json ///
{
  "version": "2.0.0",
  "tasks": [
    { "label": "codex: objective (CLI)", "type": "shell", "command": "codex-cli objective --trace=summary" },
    { "label": "codex: plan from git diff", "type": "shell", "command": "git diff | codex-cli plan --stdin --fail-fast" },
    { "label": "codex: apply (dry-run)", "type": "shell", "command": "codex-cli apply --dry-run" },
    { "label": "codex: apply (confirm)", "type": "shell", "command": "codex-cli apply --yes" }
  ]
}
/// README.md ///
# kbx - O Maestro da Squad de IAs do Kubex

`kbx` é a ferramenta de linha de comando que serve como ponto de entrada único e abstrato para o ecossistema de desenvolvimento assistido por IA do Kubex. Ele orquestra os diferentes agentes (`codex`, `g-agent`, `copilot`) para garantir que cada contribuição, de qualquer colaborador, esteja alinhada com os princípios e a qualidade definidos no Manifesto Kubex.

> **Nosso lema:** Um comando para começar. A filosofia do projeto para garantir a qualidade.

## 🚀 Início Rápido

### Pré-requisitos

- Go 1.25+
- Make

### Instalação

```bash
# Compila e instala o binário kbx no seu GOPATH/bin
make install
Configuração
O kbx é projetado para funcionar out-of-the-box em um repositório Kubex. Ele procura automaticamente por uma pasta .kubex/ contendo os arquivos de configuração e manifestos.

🛠️ Uso Básico
Verifique se o ambiente está pronto:

Bash

kbx doctor
Inicie uma nova tarefa com um objetivo claro:

Bash

kbx objective "Refatorar o parser de configuração para suportar a nova flag --strict"
Audite suas mudanças atuais contra os princípios do projeto:

Bash

kbx audit
Prepare seu commit, garantindo a conformidade com os padrões:

Bash

kbx commit
/// docs/architecture/Resilience-by-Design.md ///
# Resilience by Design — Kortex/Kubex (Outline)

> Descrever a arquitetura resiliente e o “offline-first com melhoria progressiva”.
> Documento curto e objetivo, focado em produto.

## 1. Princípios
- *Offline-first*, reconexão automática, *event broadcasting* interno
- UI com indicadores de conectividade e frescor
- Fallback seguro sem MCP, upgrades com WebSockets quando presente

## 2. Padrão de Estados
- `AppContext` global (Next.js/TS strict)
- Filas de eventos, handlers idempotentes
- Retentativas com *backoff* e circuit breakers

## 3. Integração MCP
- Conexão WebSocket com reconexão
- Logs e métricas em tempo real
- Estratégia de degradação graciosa

## 4. Build/Deploy
- Static export (`output: 'export'`) para GitHub Pages
- Imagens `unoptimized: true`, trailingSlash
- CI com lint/test/build e checagens de qualidade

## 5. Aceite
- Build estável
- Sem *runtime errors* offline
- Realtime upgrades funcionando quando MCP disponível
/// docs/brand/kubex_brand_visual_spec.v0.0.1.md ///
# Kubex Brand Visual Spec
A visual identity guide for all Kubex modules and projects. This document ensures consistent branding across repositories, banners, and promotional assets.

## 🎨 Color Palette
- **Base / Background**: Deep dark blue-black (`#0a0f14` → `#10151b`), subtle grid/dots
- **Primary**: Ciano `#00f0ff`, Verde `#00e676`
- **Secondary**: Roxo `#7c4dff`, Lilás `#8e24aa`
- **Contrast/Text**: `#e0f7fa`

## 🔤 Typography
- **Title**: futuristic sans (Orbitron/Exo2/Bank Gothic), uppercase + neon glow
- **Subtitle**: monospace (IBM Plex Mono/Source Code Pro), sutil

## 🔷 Elements
Hexágonos com outline neon, linhas de conexão finas, grid/circuitos, ícones minimalistas (Docker, Helm, chat, charts). Radial glow sutil atrás do título.

## 📐 Composition
1) Título central (neon)  
2) Subtítulo curto abaixo  
3) Ícones periféricos em hexágonos conectados  
4) Espaço para “respirar” (sem poluição visual)

## 🧩 Modular Abstraction
- Mesmo DNA visual; cor primária distinta por módulo (ex.: GHbex→Ciano, Kortex→Teal, Logz→Green, GoBE→Purple)
- Ícones relevantes por domínio
- Composição consistente

## ✅ Checklist de uso
- [ ] Background dark + textura leve
- [ ] Título central neon
- [ ] Subtítulo monospace
- [ ] 2–3 ícones periféricos
- [ ] Linhas finas conectando
- [ ] Cor primária do módulo
- [ ] Grid sutil ao fundo
/// docs/governance/KUBEX_MANIFESTO.md ///
# Kubex — Manifesto de Co-Autoria (Agentes Cofundadores)

> **Princípio**: agentes não “entregam tarefas”, **entregam produto**. Nada de arquivos de relato de tarefa no repo.

## Diretrizes
1. **Objetivo verificável**: cada change precisa de resultado observável (build/test/endpoint/UX).
2. **Sem task-reports**: bloquear `TASK_SUMMARY*.md`, `agent_report*.md` e afins.
3. **Documentação útil**: README/guia/ADRs curtos; nada de diários de execução.
4. **Consistência Kubex**: seguir Brand Visual + Craftsmanship Standards (Go/TS).
5. **Qualidade**: testes, lint, security checks.

## “Done”
- ✅ Build ok, sem warnings críticos
- ✅ Testes cobrindo paths felizes/erro
- ✅ Docs de uso atualizados quando aplicável
- ✅ Nenhum relatório de tarefa inserido
# Module Specification

```go
// Module is an interface that defines the methods a module must implement.
type Module interface {
  Alias() string
  ShortDescription() string
  LongDescription() string
  Usage() string
  Examples() []string
  Active() bool
  Module() string
  Execute() error
  Command() *cobra.Command
  SetParentCmdName(rtCmd string)
  concatenateExamples() string
}
```

```go
// Simple Generic Module Concrete Case
type ModuleConcrete struct {
 parentCmdName string
 printBanner   bool
}

/*
  Abaixo estão os métodos que representam wrappers
  de dados básicos do módulo, enquanto CLI e Lib.

  Esses métodos são obrigatórios para qualquer módulo.
 */
func (m *ModuleConcrete) Alias() string {
 return "module"
}
func (m *ModuleConcrete) ShortDescription() string {
 return "ModuleConcrete short description."
}
func (m *ModuleConcrete) LongDescription() string {
 return `ModuleConcrete long description.`
}
func (m *ModuleConcrete) Usage() string {
 return "module [command] [args]"
}
func (m *ModuleConcrete) Examples() []string {
 return []string{
  "module doctor <subcommand>",
  "module platform <subcommand>",
 }
}
func (m *ModuleConcrete) Active() bool {
 return true
}
func (m *ModuleConcrete) Module() string {
 return "module"
}
func (m *ModuleConcrete) Execute() error {
 return m.Command().Execute()
}

func (m *ModuleConcrete) Command() *cobra.Command {
/*
  Todo módulo, pode opcionalmente usar o `Annotations`
  como uma camada de metadados para o comando.
  Ele contém dois argumentos:
    - Descriptions: []string{}
      Ele irá lidar com os Short e Long descriptions
      contemplando a abstração que há trazendo estilização
      e formatação.

    - hideBanner: bool
      Se verdadeiro, oculta a exibição do banner do módulo e
      torna o comportamento do Cobra em relação ao Short e Long
      descriptions nativo novamente.
 */

 rootCmd := &cobra.Command{
  Use:   m.Module(),
  Alias: m.Alias(),
  Short: m.ShortDescription(),
  Long:  m.LongDescription(),
  // Aqui há um tratamento quebrando as linhas do
  // array de exemplos na exibição. Permitindo
  // que seja veículado vários exemplos para o mesmo comando
  // de forma organizada.
  Example: m.concatenateExamples(),

  Annotations: m.GetDescriptions(
    []string{
      m.ShortDescription(),
      m.LongDescription(),
    }, m.printBanner,
  ),
  RunE: func(cmd *cobra.Command, args []string) error {
   return m.Execute()
  },
 }

// Aqui você pode adicionar subcomandos ao comando raiz
// Ex: rootCmd.AddCommand(subCmd1, subCmd2)

// Após inserir os comandos, é necessário acionar o
// registro do template e cores para o Cobra.
// Isso faz com que o CLI tenha uma aparência
// mais amigável e estilizada, além de padronizada
// para todos os módulos do ecossistema.
 setUsageDefinition(rtCmd)

 // Esse loop existe para garantir que os padrões de
 // estilo também sejam aplicados aos subcomandos.
 for _, c := range rtCmd.Commands() {
  setUsageDefinition(c)
  if !strings.Contains(strings.Join(os.Args, " "), c.Use) {
   if c.Short == "" {
    c.Short = c.Annotations["description"]
   }
  }
 }

// Retorna o comando completo.
 return rootCmd
}
```
