# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Overview

**Grompt** is a modern prompt engineering tool that transforms unstructured ideas into clean, effective prompts for AI models. Built with Go (backend) and React 19 (frontend), it runs as a single binary with zero dependencies.

The project is part of the **Kubex ecosystem** and follows standardized architectural patterns defined in AGENTS.md.

## Build and Development Commands

### Building

```bash
# Build production binary (outputs to ./dist/grompt)
make build

# Build for development (with debug symbols)
make build-dev

# Build for specific platform
make build linux amd64
make build darwin arm64
make build windows amd64

# Install to system (builds + installs to PATH)
make install

# Clean build artifacts
make clean
```

### Testing

```bash
# Run all tests
make test

# Run specific test files
go test ./tests/tests_engine/...
go test ./tests/tests_module/...

# Run tests with coverage
go test -cover ./...
```

### Frontend Development

```bash
cd frontend

# Install dependencies
npm run setup

# Development server (Vite)
npm run dev

# Build production frontend
npm run build

# Preview production build
npm run preview
```

### Running the Application

```bash
# Start web server (default: http://localhost:8080)
./grompt start

# Start on custom port
./grompt start -p 5000

# CLI: Ask a question directly
./grompt ask "What is Go?" --provider gemini --model gemini-2.0-flash

# CLI: Generate prompt from ideas
./grompt generate --idea "REST API" --idea "authentication" --purpose code

# CLI: Generate AI squad recommendations
./grompt squad "Build a payment microservice with Stripe"

# Gateway mode (enterprise features)
./grompt gateway start
```

## Architecture Overview

### Project Structure (Kubex Standard)

```
grompt/
├── cmd/                        # Application entrypoints
│   ├── main.go                 # Main CLI entrypoint
│   └── cli/                    # CLI command definitions (ask, generate, squad, etc.)
├── internal/                   # Core private logic
│   ├── module/                 # Module interface & wrapper (RegX pattern)
│   │   ├── module.go           # Implements universal Kubex module interface
│   │   ├── wrpr.go             # RegX() wrapper for module registration
│   │   ├── logger/             # Universal logger wrapper (logz)
│   │   ├── info/               # Application metadata & banners
│   │   └── usage.go            # Custom CLI styling
│   ├── engine/                 # Core prompt engineering logic
│   ├── providers/              # AI provider integrations (OpenAI, Claude, etc.)
│   ├── gateway/                # Enterprise gateway server
│   │   ├── server.go           # HTTP server with CORS & middleware
│   │   ├── transport/          # HTTP routing
│   │   ├── middleware/         # Production middleware (rate limiting, etc.)
│   │   └── registry/           # Provider registry
│   ├── services/               # Business logic services
│   │   ├── agents/             # AI agent generation
│   │   ├── squad/              # AI squad recommendations
│   │   └── server/             # Web server
│   └── types/                  # Type definitions & interfaces
├── factory/                    # Public API constructors
│   ├── providers/              # Provider factory functions
│   ├── engine/                 # Engine factory functions
│   └── templates/              # Template management
├── frontend/                   # React 19 + Vite + TypeScript
├── tests/                      # All test files
├── support/                    # Build hooks & scripts
│   ├── main.sh                 # Main build orchestrator
│   ├── pre.d/                  # Pre-build hooks
│   └── pos.d/                  # Post-build hooks
├── grompt.go                   # Public library interface
├── Makefile                    # Build system (delegates to support/main.sh)
└── go.mod                      # Go 1.25.1+
```

### Key Architectural Patterns

#### 1. Module Interface (Kubex Standard)

All Kubex modules implement a universal interface defined in `internal/module/module.go`:

```go
type Grompt struct {
    parentCmdName string
    PrintBanner   bool
}

// Required methods:
func (m *Grompt) Alias() string
func (m *Grompt) ShortDescription() string
func (m *Grompt) LongDescription() string
func (m *Grompt) Usage() string
func (m *Grompt) Examples() []string
func (m *Grompt) Active() bool
func (m *Grompt) Module() string
func (m *Grompt) Execute() error
func (m *Grompt) Command() *cobra.Command
```

#### 2. RegX Wrapper Pattern

- **Location**: `internal/module/wrpr.go`
- **Purpose**: Provides `func RegX() *Grompt` - the canonical way to access the module
- **Critical**: Do NOT modify `wrpr.go`; customizations go in `module.go`

#### 3. Internal-First Architecture

- All core logic resides in `internal/`
- External access is provided via:
  - **Interfaces** exported through `factory/` or `api/`
  - **Public API** in root-level `grompt.go`

Example from `grompt.go`:

```go
// Public interface for library use
type PromptEngine interface {
    ProcessPrompt(template string, vars map[string]interface{}) (*Result, error)
    GetProviders() []Provider
    BatchProcess(prompts []string, vars map[string]interface{}) ([]Result, error)
}

// Factory constructor
func NewPromptEngine(config Config) PromptEngine {
    return engine.NewEngine(config)
}
```

#### 4. CLI Command Organization

- `cmd/main.go`: Minimal entrypoint - calls `module.RegX().Command().Execute()`
- `cmd/cli/`: Contains command implementations
  - `ai.go`: ask, generate, chat commands
  - `service.go`: start command (web server)
  - `squad.go`: squad command
  - `gateway.go`: gateway commands
  - `daemon.go`: daemon commands

#### 5. Configuration & Environment

- Config loading happens **per-command** in `cmd/cli/`, NOT in `main.go`
- All envvars have **fallback defaults** for resilience
- Manifest: `internal/module/info/manifest.json` (metadata, version, etc.)

#### 6. Logger (Universal Kubex Pattern)

Always use the centralized logger:

```go
import gl "github.com/kubex-ecosystem/grompt/internal/module/logger"

gl.Log("info", "Starting server...")
gl.Log("debug", "Processing request")
gl.Log("error", "Failed to connect")
gl.Log("fatal", "Critical error")
```

#### 7. Provider Architecture

- **Internal**: `internal/providers/` - concrete implementations
- **Factory**: `factory/providers/` - constructor functions
- **Interface**: `internal/types/provider.go` - Provider interface
- **Supported**: OpenAI, Claude (Anthropic), Gemini, DeepSeek, Ollama, ChatGPT

Providers are initialized in `internal/engine/engine.go:initializeProviders()`

#### 8. Gateway (Enterprise Features)

The gateway (`internal/gateway/`) provides:

- Multi-provider routing with registry pattern
- Production middleware (rate limiting, auth, metrics)
- CORS support for web clients
- BYOK (Bring Your Own Key) via HTTP headers

Key headers: `X-API-Key`, `X-OPENAI-Key`, `X-CLAUDE-Key`, `X-GEMINI-Key`

### Frontend Architecture

**Stack**: React 19, Vite, TypeScript, TailwindCSS

- **Entry**: `frontend/src/main.tsx`
- **AI SDKs**: `@anthropic-ai/sdk`, `openai`, `@google/genai`
- **Features**: Dark/light themes, responsive UI, real-time prompt engineering

Build output is embedded in the Go binary during `make build`.

## Development Guidelines

### Go Standards (from .github/copilot-instructions.md)

1. **Module Management**: Keep `go.mod` clean; avoid indirect dependencies
2. **Project Structure**: Follow `cmd/`, `internal/`, `factory/`, `tests/` layout
3. **Package Comments**: Every package MUST have a single-line `// Package <name> ...` comment
4. **Testing**: Table-driven tests with `testing` package; use `testify` for complex assertions
5. **Mocking**: Mock via interfaces, NEVER via globals
6. **Naming**: CamelCase for exported, camelCase for internal; avoid package name stuttering
7. **Error Handling**: Handle errors explicitly; return early; avoid nesting
8. **Context**: Always use `context.Context` for cancellation/timeouts; pass explicitly
9. **Documentation**: Exported functions MUST have godoc comments starting with the name
10. **Composition**: Favor composition over inheritance; accept interfaces, return structs

### Testing

- **Unit tests**: `tests/tests_*/` directories
- **Integration tests**: `tests/*.sh` scripts
- **Example**: `tests/analyzer_core_test.go` (table-driven test pattern)

### Custom Build Hooks

Support modular build automation via:

- `support/pre.d/`: Scripts executed **before** build
- `support/pos.d/`: Scripts executed **after** build

Scripts run in lexicographic order (use `01-`, `02-` prefixes). Each runs in a subshell.

Example:

```bash
# support/pre.d/10-setup-env.sh
#!/usr/bin/env bash
echo "🔧 [pre.d] Setting up environment..."
export GROMPT_ENV="dev"
```

## API Endpoints

When running `grompt start` or `grompt gateway start`:

```
GET  /api/config      # Available providers & configuration
GET  /api/health      # Server health status
GET  /api/models      # Available models per provider

POST /api/unified     # Unified endpoint (all providers)
POST /api/openai      # OpenAI-specific
POST /api/gemini      # Gemini-specific
POST /api/claude      # Claude-specific
POST /api/deepseek    # DeepSeek-specific
POST /api/ollama      # Ollama-specific
```

### BYOK Support

External API keys can be provided per request:

```bash
curl -X POST http://localhost:8080/api/unified \
  -H "Content-Type: application/json" \
  -H "X-API-Key: sk-your-key" \
  -d '{"provider":"openai","prompt":"Hello","max_tokens":100}'
```

## Environment Variables

```bash
# Server
PORT=8080                         # Server port
DEBUG=true                        # Enable debug mode

# AI Providers (all optional)
OPENAI_API_KEY=sk-...             # OpenAI
CLAUDE_API_KEY=sk-ant-...         # Anthropic Claude
GEMINI_API_KEY=...                # Google Gemini
DEEPSEEK_API_KEY=...              # DeepSeek
CHATGPT_API_KEY=...               # ChatGPT
OLLAMA_ENDPOINT=http://localhost:11434  # Ollama
```

## Library Usage (External Consumption)

Grompt can be used as a library:

```go
import "github.com/kubex-ecosystem/grompt"

// Create prompt engine
config := grompt.DefaultConfig("")
engine := grompt.NewPromptEngine(config)

// Process prompt
result, err := engine.ProcessPrompt("Explain {{topic}}", map[string]interface{}{
    "topic": "quantum computing",
})

// Get available providers
providers := engine.GetProviders()
```

## Important Files

- `internal/module/info/manifest.json`: Version, metadata, platform targets
- `internal/module/info/application.go`: Application banners & branding
- `internal/module/usage.go`: Custom CLI design (colors, layout)
- `grompt.go`: Public library interface
- `Makefile`: Build orchestration (delegates to `support/main.sh`)

## Common Tasks

### Adding a New CLI Command

1. Create command in `cmd/cli/<name>.go`
2. Implement using Cobra pattern
3. Register in `internal/module/module.go:Command()` via `rtCmd.AddCommand()`
4. Add examples to module's `Examples()` method

### Adding a New AI Provider

1. Implement `internal/types/provider.go` interface
2. Create constructor in `factory/providers/`
3. Register in `internal/engine/engine.go:initializeProviders()`
4. Add configuration to `internal/types/config.go`

### Updating Version

1. Edit `internal/module/info/manifest.json` → `"version"`
2. Rebuild: `make build`

## References

- Main README: `README.md` (English)
- Portuguese README: `docs/README.pt-BR.md`
- Kubex Standards: `AGENTS.md`
- Go Copilot Instructions: `.github/copilot-instructions.md`
- Contributing: `docs/CONTRIBUTING.md`

## Notes

- Binary size: ~15MB (includes embedded React frontend)
- Memory: ~20MB idle, ~50MB under load
- Platform targets: Linux, macOS, Windows (amd64, arm64, 386)
- Go version: 1.25.1+
- React version: 19+
