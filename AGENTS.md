# Kubex AGENTS.md (Universal)

Este documento serve como guia universal para configuração dos **Agents** em todos os repositórios Kubex. Ele unifica padrões já estabelecidos em cada projeto, evitando divergências e garantindo consistência para automações, bots e contribuidores humanos.

## Contexto

Os agentes (Claude, Codex, Copilot, etc.) leem este documento como referência de **estrutura, padrões e práticas vigentes** no ecossistema Kubex. O objetivo é reduzir erros de interpretação e facilitar a manutenção de múltiplos repositórios.

## Project Standards (Current Reality)

- **Manifest.json**: por padrão em `internal/module/info/manifest.json`. Se movido, é necessário atualizar manualmente as referências.
- **Wrappers de módulo**: `internal/module/module.go` é a estrutura principal; `internal/module/wrpr.go` contém wrappers auxiliares.
- **cmd/**: `main.go` atua como entrypoint principal. `cmd/cli/` guarda os entrypoints para wrappers de comandos CLI.
- **Envvars**: sempre com fallback para valores padrão, garantindo resiliência.
- **Carregamento de configs**: não ocorre no `main.go`; cada comando em `cmd/cli/` carrega apenas o necessário.
- **READMEs**: todos os projetos têm `README.md` em inglês e opcional `docs/README.pt-BR.md`, com link no README em inglês.
- **Makefile e support/**: genéricos, reagem às chaves definidas no `manifest.json`.
- **Mocks**: sempre centralizados em módulos específicos, nunca hardcoded em lógicas reais.
- **Logger universal (logz)**: todos os projetos usam o wrapper padrão em `internal/module/logger`. Recomenda-se importar com alias para consistência:

  ```go
  import (
      gl "github.com/kubex-ecosystem/analyzer/internal/module/logger"
  )
  ```

- **Wrapper RegX**: todo módulo possui `internal/module/wrpr.go`, que contém exclusivamente o wrapper `func RegX() *[NOME_DO_MODULO]`. Isso padroniza o acesso ao módulo e mantém o `module.go` livre para customizações específicas sem risco de alterar o wrapper global.

- **Internal como núcleo**: toda lógica central do projeto, ou qualquer parte que trate de algo sensível, crítico ou de complexidade média‑alta/alta, deve ficar em `internal/`.
  - Se essa lógica precisar ser exposta para outro módulo ou uso externo, deve ser feita via **interfaces**, com construtores que permitam a instanciação dessas interfaces.
  - A exportação deve ocorrer preferencialmente em `api/` ou `factory/`. Somente em casos específicos (como CLI ou domínios não ligados ao `internal/`) a exportação poderá estar fora desses packages.

- **Interface universal dos módulos**: o arquivo `internal/module/module.go` implementa a interface comum de **todos os módulos Kubex**. Cada módulo segue esse padrão, com métodos como `Alias()`, `ShortDescription()`, `LongDescription()`, `Usage()`, `Examples()`, `Active()`, `Module()`, `Execute()` e `Command()`. Essa estrutura garante consistência entre módulos e integração fluida com Cobra para CLI.

  Exemplo simplificado:

  ```go
  type Analyzer struct {
      parentCmdName string
      hideBanner   bool
  }

  func (m *Analyzer) Alias() string           { return "" }
  func (m *Analyzer) ShortDescription() string { return "AI tools help in the editor, but they stop antes do PR, lacking governance." }
  func (m *Analyzer) LongDescription() string  { return `Analyzer: An AI-powered tool...` }
  func (m *Analyzer) Usage() string            { return "Analyzer [command] [args]" }
  func (m *Analyzer) Active() bool             { return true }
  func (m *Analyzer) Module() string           { return "Analyzer" }
  func (m *Analyzer) Execute() error           { return m.Command().Execute() }
  // ... demais métodos garantindo padronização
  ```

- **Banners**: todos os banners estão em `internal/module/info/application.go`, junto com o método auxiliar que cuida da lógica de impressão. Eles não ficam em `module/` porque também são usados pelos comandos em `cmd/cli/`. Se estivessem em `main.go`, criariam dependência cíclica, por isso foram isolados nesse arquivo específico.

- **CLI design customizado**: todo design customizado do CLI (cores, layout e estilos de exibição) está centralizado em `internal/module/usage.go`. Esse arquivo define a aparência dos comandos e mensagens, mantendo o padrão visual consistente entre os módulos.
