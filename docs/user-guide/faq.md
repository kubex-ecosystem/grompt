---
title: Perguntas Frequentes (FAQ)
---

# Perguntas Frequentes

## Qual provedor devo usar?

- **OpenAI (gpt-4/4-turbo)**: melhor equilíbrio para geral/código.
- **Claude (sonnet/opus)**: raciocínio e análise fortes.
- **DeepSeek (coder)**: foco em programação e custo.
- **Ollama (local)**: privacidade e offline; requer download de modelos.

## Como configurar as chaves com segurança?

- Exporte via ambiente (`export OPENAI_API_KEY=...`) ou use `~/.gromptrc`.
- Garanta permissões `chmod 600 ~/.gromptrc`.
- Nunca versione `.env`/`.gromptrc`.

## Os testes fazem chamadas reais?

Não por padrão. Testes de rede são opcionais e controlados por env (ex.: `GROMPT_E2E_GEMINI=1`). Veja “Desenvolvimento > Testes”.

## O binário não inicia na porta 8080

- Defina `PORT` livre (ex.: `PORT=8081 ./grompt`).
- Verifique se outro processo não ocupa a porta (`lsof -i :8080`).

## Erro: chave de API ausente

- Verifique `OPENAI_API_KEY`/`CLAUDE_API_KEY`/`GEMINI_API_KEY`/etc.
- Para Ollama, confirme `OLLAMA_ENDPOINT` e `ollama serve` rodando.

## Como publicar as docs?

Use `make pub-docs`. Para pré-visualizar, rode `mkdocs serve -f support/docs/mkdocs.yml`.

## Existe suporte em Windows/macOS/Linux?

Sim. Baixe o binário correspondente nas releases. Para Ollama no Windows, use o instalador do site do projeto.

## Onde acho exemplos práticos?

Veja “Guia do Usuário > Exemplos de Uso” e “CLI > Comandos”.
