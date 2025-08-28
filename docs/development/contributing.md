# Contribuindo

Obrigado por seu interesse em contribuir com o Grompt! Este guia explica como você pode ajudar a melhorar o projeto.

## 🎯 Formas de Contribuir

### 📝 Reportar Bugs

Encontrou um bug? Ajude-nos a corrigi-lo!

1. **Verifique se já existe** uma issue sobre o problema
2. **Crie uma nova issue** com informações detalhadas:
   - Descrição clara do problema
   - Passos para reproduzir
   - Comportamento esperado vs atual
   - Screenshots ou logs, se aplicável
   - Informações do ambiente (SO, versão do Go, etc.)

**Template de Bug Report:**

```markdown
## Descrição
[Descrição clara e concisa do bug]

## Passos para Reproduzir
1. Execute o comando '...'
2. Acesse a página '...'
3. Clique em '...'
4. Veja o erro

## Comportamento Esperado
[O que deveria acontecer]

## Comportamento Atual
[O que está acontecendo]

## Ambiente
- OS: [e.g. Ubuntu 22.04]
- Grompt Version: [e.g. v1.2.3]
- Go Version: [e.g. 1.21.5]
- Browser: [e.g. Chrome 120]

## Logs/Screenshots
[Adicione logs ou screenshots]
```

### ✨ Solicitar Funcionalidades

Tem uma ideia para melhorar o Grompt?

1. **Verifique se já foi solicitada** nas issues
2. **Descreva claramente** a funcionalidade desejada
3. **Explique o caso de uso** e benefícios
4. **Considere a implementação** se possível

**Template de Feature Request:**

```markdown
## Funcionalidade
[Descrição clara da funcionalidade]

## Motivação
[Por que esta funcionalidade é útil]

## Solução Proposta
[Como você imagina que deveria funcionar]

## Alternativas Consideradas
[Outras abordagens que você considerou]

## Contexto Adicional
[Screenshots, mockups, links, etc.]
```

### 💻 Contribuir com Código

#### Configuração do Ambiente

1. **Fork** o repositório
2. **Clone** seu fork:

   ```bash
   git clone https://github.com/seu-usuario/grompt.git
   cd grompt
   ```

3. **Configure o ambiente:**

   ```bash
   # Instalar dependências
   make deps

   # Verificar se tudo está funcionando
   make test
   make build
   ```

#### Fluxo de Desenvolvimento

1. **Crie uma branch** para sua funcionalidade:

   ```bash
   git checkout -b feature/nome-da-funcionalidade
   ```

2. **Faça suas alterações** seguindo os padrões do projeto

3. **Execute os testes:**

   ```bash
   make test
   make lint
   ```

4. **Commit suas alterações:**

   ```bash
   git add .
   git commit -m "feat: adiciona nova funcionalidade X"
   ```

5. **Push para seu fork:**

   ```bash
   git push origin feature/nome-da-funcionalidade
   ```

6. **Abra um Pull Request**

## 📋 Padrões de Código

### Go (Backend)

- **Siga os padrões do Go:** use `gofmt`, `golint`, `go vet`
- **Documentação:** todas as funções exportadas devem ter comentários
- **Testes:** mantenha cobertura de testes alta
- **Nomenclatura:** use convenções idiomáticas do Go

**Exemplo:**

```go
// Package engine provides prompt engineering capabilities
package engine

// GeneratePrompt creates a structured prompt from raw ideas
func GeneratePrompt(ideas []string, purpose Purpose) (*Prompt, error) {
    if len(ideas) == 0 {
        return nil, errors.New("at least one idea is required")
    }

    // Implementation...
    return prompt, nil
}
```

### TypeScript/React (Frontend)

- **TypeScript:** use tipagem estrita
- **Componentes funcionais:** prefira hooks
- **Styled-components:** para estilização
- **ESLint/Prettier:** para formatação consistente

**Exemplo:**

```typescript
interface PromptFormProps {
  onSubmit: (ideas: string[]) => void;
  isLoading?: boolean;
}

export const PromptForm: React.FC<PromptFormProps> = ({
  onSubmit,
  isLoading = false
}) => {
  const [ideas, setIdeas] = useState<string[]>([]);

  const handleSubmit = useCallback((e: FormEvent) => {
    e.preventDefault();
    onSubmit(ideas);
  }, [ideas, onSubmit]);

  return (
    <StyledForm onSubmit={handleSubmit}>
      {/* Component implementation */}
    </StyledForm>
  );
};
```

### Commits

Use **Conventional Commits:**

- `feat:` nova funcionalidade
- `fix:` correção de bug
- `docs:` documentação
- `style:` formatação (sem mudança de código)
- `refactor:` refatoração
- `test:` adição/correção de testes
- `chore:` tarefas de manutenção

**Exemplos:**

```bash
feat: adiciona suporte para provedor Gemini
fix: corrige erro de timeout em requests longos
docs: atualiza guia de instalação
test: adiciona testes para módulo de templates
```

## 🔍 Process de Review

### Checklist do Pull Request

- [ ] **Código** segue os padrões estabelecidos
- [ ] **Testes** passam e cobertura não diminui
- [ ] **Documentação** está atualizada
- [ ] **Commit messages** seguem conventional commits
- [ ] **Breaking changes** estão documentadas
- [ ] **Performance** não foi degradada

### O que Esperamos

1. **Código limpo e legível**
2. **Testes abrangentes**
3. **Documentação atualizada**
4. **Compatibilidade mantida**
5. **Performance considerada**

### Processo de Review

1. **Automated checks** devem passar
2. **Manual review** por pelo menos um mantenedor
3. **Feedback** pode ser solicitado
4. **Merge** após aprovação

## 🏗️ Arquitetura do Projeto

### Estrutura de Diretórios

```plaintext
grompt/
├── cmd/                    # CLI entrypoints
│   ├── main.go            # Servidor principal
│   └── cli/               # Comandos CLI
├── internal/              # Código interno
│   ├── engine/            # Core prompt engineering
│   ├── providers/         # AI providers
│   ├── services/          # Business logic
│   └── types/             # Type definitions
├── frontend/              # React application
│   ├── src/
│   │   ├── components/    # React components
│   │   ├── hooks/         # Custom hooks
│   │   └── utils/         # Utilities
├── docs/                  # Documentação
├── tests/                 # Testes
└── support/               # Scripts de build
```

### Componentes Principais

1. **Engine:** Core prompt engineering logic
2. **Providers:** Integrações com APIs de IA
3. **Services:** Lógica de negócio e orquestração
4. **Frontend:** Interface React
5. **CLI:** Interface de linha de comando

## 🧪 Testes

### Executar Testes

```bash
# Todos os testes
make test

# Testes do backend
make test-backend

# Testes do frontend
make test-frontend

# Testes de integração
make test-integration

# Coverage
make coverage
```

### Escrevendo Testes

#### Go (Backend Tests)

```go
func TestGeneratePrompt(t *testing.T) {
    tests := []struct {
        name     string
        ideas    []string
        purpose  Purpose
        expected *Prompt
        wantErr  bool
    }{
        {
            name:    "valid ideas",
            ideas:   []string{"create API", "use Node.js"},
            purpose: PurposeCode,
            expected: &Prompt{
                Content: "Create a Node.js API...",
                Purpose: PurposeCode,
            },
            wantErr: false,
        },
        // More test cases...
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            result, err := GeneratePrompt(tt.ideas, tt.purpose)

            if tt.wantErr {
                assert.Error(t, err)
                return
            }

            assert.NoError(t, err)
            assert.Equal(t, tt.expected, result)
        })
    }
}
```

#### React (Frontend)

```typescript
describe('PromptForm', () => {
  it('should submit ideas when form is submitted', () => {
    const onSubmit = jest.fn();
    render(<PromptForm onSubmit={onSubmit} />);

    const input = screen.getByLabelText('Add idea');
    const submitButton = screen.getByRole('button', { name: 'Generate' });

    fireEvent.change(input, { target: { value: 'test idea' } });
    fireEvent.click(screen.getByRole('button', { name: 'Add' }));
    fireEvent.click(submitButton);

    expect(onSubmit).toHaveBeenCalledWith(['test idea']);
  });
});
```

## 📚 Documentação

### Onde Documentar

- **README:** visão geral e quick start
- **Docs:** documentação detalhada (MkDocs)
- **Código:** comentários inline
- **API:** documentação de endpoints

### Atualizando Documentação

```bash
# Servir docs localmente
cd support/docs
mkdocs serve

# Build da documentação
mkdocs build

# Deploy (apenas mantenedores)
mkdocs gh-deploy
```

## 🎯 Primeiras Contribuições

### Good First Issues

Procure por issues marcadas com:

- `good-first-issue`: ideal para iniciantes
- `help-wanted`: ajuda é bem-vinda
- `documentation`: melhorias na documentação
- `bug`: correção de bugs simples

### Áreas que Precisam de Ajuda

1. **Documentação:** sempre pode ser melhorada
2. **Testes:** aumentar cobertura
3. **Performance:** otimizações
4. **UI/UX:** melhorias na interface
5. **Provedores de IA:** novos integrações

## 💬 Comunicação

### Onde Buscar Ajuda

- **GitHub Issues:** problemas específicos
- **GitHub Discussions:** discussões gerais
- **Código:** comentários em PRs

### Diretrizes de Comunicação

- **Seja respeitoso** e profissional
- **Use português ou inglês** conforme preferir
- **Seja específico** em perguntas e descrições
- **Compartilhe contexto** relevante

## 🏆 Reconhecimento

### Contributors

Todos os contribuidores são reconhecidos:

- **README:** lista de contributors
- **CHANGELOG:** créditos por release
- **Commits:** histórico permanente

### Como ser Reconhecido

1. **Contribua regularmente**
2. **Ajude outros contribuidores**
3. **Mantenha qualidade alta**
4. **Participe de discussões**

## 📄 Licença

Ao contribuir, você concorda que suas contribuições serão licenciadas sob a **MIT License**.

---

## 📞 Contato

- **GitHub Issues:** [github.com/rafa-mori/grompt/issues](https://github.com/rafa-mori/grompt/issues)
- **GitHub Discussions:** [github.com/rafa-mori/grompt/discussions](https://github.com/rafa-mori/grompt/discussions)
- **Email:** [através do GitHub](https://github.com/rafa-mori)

---

***Obrigado por contribuir! 🙏***

*Cada contribuição, por menor que seja, faz a diferença na comunidade open source.*
