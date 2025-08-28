# Documentação do Grompt

Esta é a documentação oficial do Grompt, construída com [MkDocs](https://www.mkdocs.org/) e [Material Theme](https://squidfunk.github.io/mkdocs-material/).

## 🚀 Início Rápido

### Pré-requisitos

- Python 3.8+
- pip

### Configuração

```bash
# Instalar dependências
make docs-install

# Servir localmente
make docs-serve

# Abrir no navegador: http://127.0.0.1:8000
```

## 📁 Estrutura

```plaintext
support/docs/
├── mkdocs.yml          # Configuração do MkDocs
├── Makefile            # Comandos de automação
├── requirements.txt    # Dependências Python
└── docs/               # Conteúdo da documentação
    ├── index.md        # Página inicial
    ├── getting-started/
    │   ├── installation.md
    │   └── quickstart.md
    ├── user-guide/
    │   ├── cli-commands.md
    │   ├── configuration.md
    │   ├── usage-examples.md
    │   └── api-reference.md
    ├── development/
    │   ├── contributing.md
    │   ├── architecture.md
    │   └── troubleshooting.md
    ├── community/
    │   ├── code-of-conduct.md
    │   ├── support.md
    │   └── authors.md
    └── project-info/
        ├── license.md
        ├── notice.md
        ├── security.md
        └── changelog.md
```

## 🛠️ Comandos Disponíveis

### Desenvolvimento

```bash
make docs-serve        # Servir localmente (http://127.0.0.1:8000)
make docs-build        # Build para produção
make docs-clean        # Limpar build
make docs-check        # Verificar por erros
```

### Deployment

```bash
make docs-deploy       # Deploy para GitHub Pages
```

### Utilidades

```bash
make docs-info         # Informações sobre o ambiente
make docs-urls         # URLs importantes
make docs-lint         # Verificar markdown (requer markdownlint)
make docs-new-page PAGE=caminho/pagina.md  # Criar nova página
```

### Ajuda

```bash
make help             # Mostrar todos os comandos
```

## 📝 Escrevendo Documentação

### Criando Nova Página

1. **Via Makefile:**

   ```bash
   make docs-new-page PAGE=user-guide/nova-pagina.md
   ```

2. **Manualmente:**

   ```bash
   touch docs/user-guide/nova-pagina.md
   ```

3. **Adicionar ao mkdocs.yml:**

   ```yaml
   nav:
     - User Guide:
         - Nova Página: user-guide/nova-pagina.md
   ```

### Padrões de Escrita

#### Cabeçalhos

```markdown
# Título Principal (H1)
## Seção (H2)
### Subseção (H3)
#### Detalhe (H4)
```

#### Blocos de Código

````markdown
```bash
# Comando bash
grompt --help
```

```javascript
// Código JavaScript
const grompt = new GromptClient();
```
````

#### Admonições

```markdown
!!! info "Informação"
    Informação útil para o usuário.

!!! warning "Atenção"
    Algo que requer atenção especial.

!!! danger "Perigo"
    Algo que pode causar problemas.

!!! success "Sucesso"
    Operação realizada com sucesso.
```

#### Links

```markdown
[Texto do link](caminho/para/pagina.md)
[Link externo](https://exemplo.com)
[Link para seção](#seção)
```

#### Imagens

```markdown
![Alt text](../assets/imagem.png)
![Alt text](https://url-externa.com/imagem.png)
```

#### Tabelas

```markdown
| Coluna 1 | Coluna 2 | Coluna 3 |
|----------|----------|----------|
| Valor 1  | Valor 2  | Valor 3  |
| Valor 4  | Valor 5  | Valor 6  |
```

### Extensões Disponíveis

O MkDocs está configurado com várias extensões:

- **Admonitions:** Blocos de destaque
- **Code Highlighting:** Syntax highlighting
- **Tables:** Suporte a tabelas
- **Footnotes:** Notas de rodapé
- **Abbreviations:** Abreviações
- **Math:** Fórmulas matemáticas
- **Mermaid:** Diagramas

## 🎨 Temas e Estilização

### Configuração do Material Theme

```yaml
theme:
  name: material
  palette:
    primary: "indigo"
    accent: "indigo"
  features:
    - navigation.tabs
    - navigation.top
    - search.highlight
```

### Assets

Coloque arquivos estáticos em:

```plaintext
docs/assets/
├── images/
├── css/
└── js/
```

## 🔄 Workflow de Contribuição

### 1. Desenvolvimento Local

```bash
# 1. Fazer fork do repositório
# 2. Clonar localmente
git clone https://github.com/seu-usuario/grompt.git
cd grompt/support/docs

# 3. Configurar ambiente
make docs-install

# 4. Fazer alterações
# 5. Testar localmente
make docs-serve

# 6. Verificar build
make docs-check
```

### 2. Submissão

```bash
# 1. Commit das alterações
git add .
git commit -m "docs: adiciona nova página sobre X"

# 2. Push para fork
git push origin feature/docs-nova-pagina

# 3. Abrir Pull Request
```

### 3. Deploy Automático

O deploy é automático via GitHub Actions quando:

- **Push para main:** Deploy para GitHub Pages
- **Pull Request:** Build de verificação

## 📋 Checklist para Novas Páginas

- [ ] **Conteúdo escrito** em português claro e objetivo
- [ ] **Navegação atualizada** no mkdocs.yml
- [ ] **Links internos** verificados
- [ ] **Código testado** (se aplicável)
- [ ] **Imagens otimizadas** (se aplicável)
- [ ] **SEO considerado** (títulos, meta descrições)
- [ ] **Build local** executado sem erros
- [ ] **Review de português** realizado

## 🔍 Solução de Problemas

### Erro: "Module not found"

```bash
# Reinstalar dependências
make docs-clean
make docs-install
```

### Erro: "Port already in use"

```bash
# Usar porta diferente
cd support/docs
python -m mkdocs serve --dev-addr=127.0.0.1:8001
```

### Erro de Build

```bash
# Verificar logs detalhados
make docs-check
```

### Problemas de Encoding

Certifique-se que arquivos estão em UTF-8:

```bash
file -bi docs/arquivo.md
# Deve mostrar: text/plain; charset=utf-8
```

## 📊 Métricas e Analytics

A documentação inclui:

- **Google Analytics:** Configurado no mkdocs.yml
- **Search:** Busca integrada
- **Git revision dates:** Datas de última modificação

## 🔗 Links Úteis

- **MkDocs:** <https://www.mkdocs.org/>
- **Material Theme:** <https://squidfunk.github.io/mkdocs-material/>
- **Markdown Guide:** <https://www.markdownguide.org/>
- **GitHub Pages:** <https://pages.github.com/>

## 📞 Suporte

Para dúvidas sobre a documentação:

- **Issues:** [github.com/rafa-mori/grompt/issues](https://github.com/rafa-mori/grompt/issues)
- **Discussions:** [github.com/rafa-mori/grompt/discussions](https://github.com/rafa-mori/grompt/discussions)

---

***Obrigado por contribuir com a documentação! 📚***
