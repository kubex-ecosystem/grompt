# Repository Guidelines

## Project Structure & Module Organization
Grompt couples a Go backend with a React/Vite UI. `cmd/` hosts CLI entrypoints (`cli/` for commands, `main.go` for the binary). Runtime logic is grouped inside `internal/` with gateway transports, provider clients, shared services, metrics, and module utilities; cross-cutting definitions stay in `types/`. The bundled frontend lives in `frontend/` with TypeScript sources, Tailwind config, and static assets. `support/` and `scripts/` provide automation used by the Make targets, while `docs/` captures architecture notes and contributor references. Tests and scenario harnesses reside in `tests/`, with additional fixtures in `factory/` and configuration samples in `config/`.

## Build, Test, and Development Commands
`make build` compiles a release binary into `dist/`, and `make build-dev` adds developer extras. `make run` launches the Go gateway plus embedded UI; use `go run cmd/main.go gateway serve -p 3000` for rapid backend iteration. `make test` executes the backend unit suite and shell integration checks, while `make validate` layers on pre-flight diagnostics. Frontend work happens with `cd frontend && npm run dev`; ship-ready assets come from `npm run build` or `npm run preview`.

## Coding Style & Naming Conventions
Apply `gofmt` (tabs, idiomatic imports) and keep packages aligned to a single responsibility (e.g., `internal/providers/anthropic`). Exported Go symbols use `CamelCase`; locals stay `camelCase`; constants favor `ALL_CAPS`. Run `go vet ./...` or `golangci-lint run` before submitting. Frontend components are PascalCase under `frontend/src/components`, hooks begin with `use`, and Tailwind utility classes drive styling—avoid custom CSS unless necessary.

## Testing Guidelines
Unit coverage follows Go conventions with `*_test.go` companions. Use `go test ./...` for complete checks or focus with package-specific paths (`go test ./internal/gateway/...`). Shell-based end-to-end scripts live in `tests/` and expect a freshly built binary; run them through `make test` or individually like `bash tests/test_gateway.sh`.

## Commit & Pull Request Guidelines
Commits use conventional prefixes (`fix:`, `feat:`, `chore:`) as seen in history and stay scoped to one logical change. Pull requests should summarise motivation, list impacted modules, and link issues. Include screenshots for UI updates. Always confirm `make test` (and relevant frontend builds) succeed before requesting review.

## Configuration & Security Notes
Copy `.env.example` to `.env` or `.env.local` when integrating providers; keep keys out of version control. Docker-driven workflows rely on `docker-compose.dev.yml` and share port settings with `support/config.sh`, so adjust carefully.
