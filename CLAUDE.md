# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

Grompt is a modern prompt engineering tool built with Go backend and React frontend. It transforms messy, unstructured thoughts into clean, effective prompts for AI models. The project runs as a single binary with zero dependencies and supports multiple AI providers (OpenAI, Claude, DeepSeek, Ollama, Gemini).

## Key Commands

### Building and Development
- `make build` - Build production binary with compression
- `make build-dev` - Build development binary without compression
- `make install` - Interactive install (download pre-compiled or build locally)
- `make clean` - Clean build artifacts
- `make test` - Run Go tests
- `make run` - Run the application

### Frontend Development
- `cd frontend && npm run setup` - Install frontend dependencies
- `cd frontend && npm run dev` - Start Vite development server
- `cd frontend && npm run build` - Build frontend for production

### Documentation
- `make build-docs` - Build documentation using MkDocs
- `make serve-docs` - Serve documentation at http://localhost:8081/docs
- `make pub-docs` - Publish documentation to GitHub Pages

### Internationalization (i18n)
- `make i18n.used` - Extract used i18n keys from frontend
- `make i18n.avail` - List available i18n keys
- `make i18n.diff` - Compare used vs available keys
- `make i18n.check` - Validate i18n completeness

## Architecture Overview

### Backend (Go)
- **Entry point**: `cmd/main.go` - Simple main that calls module system
- **Module system**: `internal/module/` - CLI framework built on Cobra
- **Gateway**: `internal/gateway/` - AI provider gateway with middleware, health checks, circuit breakers
- **Services**: `internal/services/` - GitHub integration, notifications, orchestration
- **Providers**: `internal/providers/` - AI provider implementations (OpenAI, etc.)
- **Config**: `config/` directory contains YAML configurations

### Frontend (React)
- **Framework**: React 19 + Vite + TypeScript
- **Styling**: TailwindCSS + PostCSS
- **Build**: Vite bundler with embedded output in Go binary
- **Location**: `frontend/` directory

### Key Architecture Patterns

1. **Gateway Pattern**: The `internal/gateway/` provides a unified interface to multiple AI providers with production middleware (rate limiting, circuit breakers, health checks)

2. **Module System**: CLI commands are organized through a module pattern in `internal/module/` using Cobra framework

3. **Embedded Frontend**: React frontend is built and embedded into the Go binary via embed directives

4. **Middleware Stack**: Production-ready middleware for health monitoring, rate limiting, and circuit breaking

5. **Provider Registry**: Pluggable AI provider system with consistent interfaces

## Configuration

- Main config: `config/config.yml` (with examples in `config/`)
- Environment variables for API keys (OPENAI_API_KEY, CLAUDE_API_KEY, etc.)
- Docker support with multiple Dockerfiles (dev, production, Koyeb)

## Build System

The project uses a sophisticated build system:
- **Makefile**: High-level interface calling `support/main.sh`
- **support/main.sh**: Main build script with extensive validation and cross-platform support
- **Custom hooks**: `support/pre.d/` and `support/post.d/` for custom build steps
- **Cross-platform**: Supports Windows, Linux, macOS with multiple architectures

## Development Notes

- Go version: 1.25.1+
- Frontend: React 19 + Next.js 15 patterns
- The project follows a monolithic architecture but with clean separation of concerns
- Extensive use of interfaces for testability and modularity
- GitHub integration for webhooks and automated workflows

## Testing

- `make test` runs the full Go test suite
- Frontend tests use standard React testing patterns
- Integration tests in `tests/` directory

## Deployment

- Single binary deployment with embedded frontend
- Docker containers supported
- Koyeb deployment configuration included
- Health check endpoints available at `/api/health`