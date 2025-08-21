# Golang Craftsmanship Standards

Use Go Modules for dependency management. Keep `go.mod` and `go.sum` clean and minimal. Avoid indirect dependencies when possible.

Organize projects using idiomatic structure: `cmd/`, `cmd/cli/`, `internal/`, `internal/types`, `internal/interfaces`, `api/`, `support/`, `support/instructions`, `tests/`.  

Place the main CLI entrypoint in `cmd/main.go` and the library entrypoint in the root withe package at the same name of project.

Every package must contain a comment describing its purpose before the package declaration in one line. Use `// Package <name> ...` format.

Write **table-driven tests** with the standard `testing` package. For complex assertions, use `testify`. Coverage should be high on business logic, especially for error paths.

Mock dependencies via interfaces — never via `globals` or side effects. Benchmark performance-sensitive functions. Keep tests fast and deterministic.

Naming: `CamelCase` for exported, `camelCase` for internal. Avoid stutter in package names (e.g., `user.User` is wrong).

Functions must be small and cohesive. Return early. Nesting is a code smell. Handle errors explicitly. Don’t ignore them — even temporarily.

Favor **composition over inheritance**. Accept interfaces, return concrete structs. Document behavior at interface boundaries.

Always use `context.Context` for cancellation, timeouts, and tracing. Pass it explicitly — do not store it in structs.

Exported types, functions, and packages MUST include **godoc-compatible comments**. Start with the function/type name. Include usage examples when applicable.

README must be clear, technical and up to date. Include build instructions, feature summary, and example usage. If possible, add architecture diagrams and CLI reference.

Be consistent. Be happy. Be fast. Be safe.
