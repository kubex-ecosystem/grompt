---
title: Publicando Docs
---

# Publicando a Documentação

Este projeto usa MkDocs (Material) com a configuração em `support/docs/mkdocs.yml` e o conteúdo em `docs/`.

## Fluxo rápido

1. Edite/adicione páginas em `docs/`.
2. Atualize a navegação em `support/docs/mkdocs.yml` (se necessário).
3. Rode `make pub-docs` para buildar e publicar no GitHub Pages.

O site é servido em: `https://kubex-ecosystem.github.io/grompt/`.

## Estrutura

- `support/docs/mkdocs.yml`: tema, plugins, navegação e integrações.
- `docs/`: páginas Markdown, ativos e seções.

## Pré-visualização local

Se desejar testar localmente:

```bash
pip install -r support/docs/requirements.txt
mkdocs serve -f support/docs/mkdocs.yml
# Acesse http://127.0.0.1:8000
```

## Boas práticas

- Mantenha títulos e slugs em PT-BR consistentes.
- Prefira exemplos executáveis e curtos; use trechos com `bash`, `http`, `json`.
- Adicione páginas novas à `nav` para ficarem descobertas no site.
- Quebre documentos longos em seções claras; use admonitions quando fizer sentido.

## Links de edição

Habilitamos `edit_uri` para exibir “Editar esta página” apontando para `edit/main/docs/` no GitHub. Prefira commits pequenos e descritivos.

## Link checking

O config possui o plugin `linkcheck` comentado. Em ambientes com rede habilitada, você pode ativar para validar os links externos:

```yaml
plugins:
  - linkcheck:
      fail_on_dead_links: true
      timeout: 10
      retry_count: 2
```

Se encontrar links off-line, prefira substituí-los por fontes estáveis.
