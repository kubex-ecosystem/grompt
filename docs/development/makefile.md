---
title: Comandos Makefile
---

# Comandos Makefile

Atalhos essenciais para desenvolvimento, build e testes. Rode a partir da raiz do repositório.

## Build

- `make build-dev linux amd64`: compila binário de desenvolvimento para `dist/` (usa o script de instalação embutido).
- `make clean`: remove artefatos de build em `dist/` e temporários.

## Testes

- `make test`: executa toda a suíte de testes Go.
- `go test ./... -v`: executa testes diretamente, com logs verbosos.
- `go test ./... -cover`: executa testes com relatório de cobertura.

### Flags úteis

- `GROMPT_E2E_GEMINI=1 GEMINI_API_KEY=... go test ./... -run TestGemini`: habilita testes de rede do Gemini (gated por env).
- `OPENAI_API_KEY=... go test ./... -run TestOpenAI`: executa testes focados em OpenAI.

## Docs

- `make pub-docs`: builda e publica a documentação (MkDocs) para GitHub Pages.

## Dicas

- Use `go mod tidy` antes de commits que adicionam dependências.
- Prefira `make test` localmente antes de abrir PRs.
