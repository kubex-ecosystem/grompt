
# i18nwatcher (kit mínimo)

**O que faz:**

* Observa `.ts/.tsx/.js/.jsx` do seu projeto.
* Para cada chamada i18n (ex.: `t("…")` / `useTranslation`), extrai **contexto** (componente, JSX, linha, etc.).
* Gera **keys determinísticas** (`component.element.slug`) e grava **drafts** no `i18n/i18n.vault.json`.

## Rodando local

```bash
go run ./cmd/hello-parser -- /caminho/do/seu/projeto
```

* Edite qualquer arquivo fonte → o vault é atualizado automaticamente.
* Arquivo do vault: `<seu-projeto>/i18n/i18n.vault.json`.

## Próximos passos (planejados)

* **Hardcoded scanner** (detectar textos de UI fora de `t()`).
* **CLI TUI** para aprovar/renomear keys.
* **Rewriter AST**:

  * Fase 1: substituir por `T("Texto")` (pass-through).
  * Fase 2: substituir por `t("key")` quando `approved`.
* **Stats** e cobertura real.
* **Adapter runtime** sem lock-in (i18next/Lingui ou seu próprio).
