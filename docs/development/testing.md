# Testing Guidelines

This project uses Go's built-in `testing` package with external test packages and centralized mocks.

## Structure
- Tests live in sibling folders named `tests_[package]` (e.g., `tests_types`, `tests_engine`).
- Test packages use `package <target>_test` to stay external.
- Reuse stable import aliases in tests: `eng`, `t`, `prov`, `tpl`, `cli`, `mod`, `m`.

## Mocks
- Use `mocks/` for shared test doubles.
- `mocks.APIConfigMock` implements `types.IAPIConfig`.
- `mocks.ConfigMock` implements `types.IConfig`.

## Commands
- `make test` – run all tests via the installer script.
- `go test ./... -cover` – run with coverage.

## Network
- Unit tests should avoid network/IO.
- Optional E2E calls (e.g., Gemini) must be gated by env flags (e.g., `GROMPT_E2E_GEMINI=1`, `GEMINI_API_KEY`).

## Conventions
- Prefer table-driven tests.
- Keep cases small, focused, and deterministic.
- Mocks go in `mocks/`, not inline.

