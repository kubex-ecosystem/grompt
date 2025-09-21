# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

Grompt is a modern prompt engineering tool built with Go (backend) and React 19 (frontend) that transforms unstructured thoughts into effective AI prompts. It runs as a single binary with zero dependencies and supports multiple AI providers (OpenAI, Claude, Gemini, DeepSeek, Ollama).

## Commands

### Build and Development
```bash
# Build using Makefile (recommended)
make build              # Build production binary to ./dist/grompt
make build-dev          # Build development version
make install            # Build and install binary to system PATH
make clean              # Clean build artifacts
make test               # Run tests

# Frontend development (requires Node.js)
cd frontend
npm run setup           # Install dependencies
npm run dev             # Start development server
npm run build           # Build for production
npm run build:static    # Build static production version
```

### Testing
```bash
make test               # Run Go tests
cd frontend && npm test # Run frontend tests (if available)
```

### Documentation
```bash
make build-docs         # Build documentation
make serve-docs         # Start documentation server at localhost:8080/docs
make pub-docs           # Publish documentation
```

### CLI Usage
```bash
# Direct AI queries
./grompt ask "What is the capital of France?" --provider openai --model gpt-4

# Generate structured prompts from ideas
./grompt generate --idea "REST API" --idea "User auth" --provider gemini

# Generate AI agent recommendations
./grompt squad "Need a backend for payments with Stripe integration"
```

## Architecture

### Backend (Go)
- **Entry Point**: `cmd/main.go` - CLI application entry
- **Library Interface**: `grompt.go` - Public API for library usage
- **Core Modules**:
  - `internal/engine/` - Prompt processing engine
  - `internal/module/` - CLI command modules and registration
  - `internal/services/` - Business logic services
  - `internal/types/` - Type definitions and interfaces
  - `internal/providers/` - AI provider implementations
- **Factories**: `factory/` - Provider, config, and template factories
- **Configuration**: Environment variables and config files

### Frontend (React 19)
- **Framework**: Vite + React 19 + TypeScript
- **Styling**: Tailwind CSS
- **Components**: `frontend/src/components/` - Reusable UI components
- **Services**: `frontend/src/services/` - API communication
- **Core**: `frontend/src/core/` - Core application logic

### Key Patterns
- **Interface-based design**: Heavy use of Go interfaces for modularity
- **Factory pattern**: For creating providers, configs, and templates
- **Module registration**: CLI commands register themselves through `internal/module`
- **Type safety**: Strong typing throughout with custom interfaces
- **Security flags**: Atomic bitflag system for security controls
- **Job states**: Atomic state management for background tasks

## API Providers

Support for multiple AI providers with unified interface:
- OpenAI (GPT-4, GPT-3.5-turbo)
- Anthropic Claude (Claude 3.5 Sonnet, Claude 3 Haiku)
- Google Gemini (Gemini 1.5 Pro, Gemini 2.0 Flash)
- DeepSeek (DeepSeek Chat, DeepSeek Coder)
- Ollama (Local models)
- Demo mode (No API key required)

Configuration via environment variables:
```bash
export OPENAI_API_KEY=sk-...
export CLAUDE_API_KEY=sk-ant-...
export GEMINI_API_KEY=...
export DEEPSEEK_API_KEY=...
export OLLAMA_ENDPOINT=http://localhost:11434
```

## Development Guidelines

Follow the patterns established in `.github/copilot-instructions.md`:
- Use Go modules for dependency management
- Organize code using idiomatic Go structure
- Write table-driven tests with standard `testing` package
- Use interfaces for mocking, avoid globals
- Handle errors explicitly, return early
- Use `context.Context` for cancellation and timeouts
- Document exported functions with godoc comments
- Favor composition over inheritance

## Build System

The project uses a sophisticated Makefile that delegates to `support/main.sh` for cross-platform builds. The build system supports:
- Custom build hooks in `support/pre.d/` (pre-build) and `support/pos.d/` (post-build)
- Cross-platform compilation for multiple architectures
- Metadata extraction from `internal/module/info/manifest.json`
- Docker deployment and documentation generation

## File Structure
```
grompt/
├── cmd/                    # CLI entry points
├── internal/               # Private application code
│   ├── engine/            # Prompt processing engine
│   ├── module/            # CLI command modules
│   ├── services/          # Business logic
│   └── types/             # Type definitions
├── factory/               # Provider and config factories
├── frontend/              # React frontend application
├── support/               # Build scripts and hooks
├── tests/                 # Test files
├── docs/                  # Documentation
├── grompt.go             # Library interface
└── Makefile              # Build system
```

## Testing

The project follows Go testing conventions:
- Tests are located alongside source files (`*_test.go`)
- Use table-driven tests for comprehensive coverage
- Mock dependencies via interfaces
- Run tests with `make test`

## Configuration

The application supports configuration through:
- Environment variables for API keys and settings
- Command-line flags for runtime options
- Configuration files (format determined by the config factory)
- Default values defined in `types.Config`