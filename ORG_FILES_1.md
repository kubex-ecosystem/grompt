## Prompt Final (com regras de alteração de arquivos)

> Você é um engenheiro sênior atuando no projeto Grompt.
>
> Seu objetivo é implementar uma feature para **geração automática de squads de agentes AI** (AGENTS.md), baseada na análise de requisitos textuais fornecidos pelo usuário.
>
> **Siga rigorosamente as seguintes diretrizes de estrutura e alteração de arquivos do projeto:**
>
> ---
>
> ### Estrutura do Projeto (respeitar e manter):
>
> ```
> .
> ├── LICENSE
> ├── NOTICE.md
> ├── README.md
> ├── SECURITY.md
> ├── Makefile
> ├── support/
> ├── cmd/
> │   ├── cli/
> │   ├── main.go   # ponto de entrada do CLI, NÃO MODIFICAR
> │   ├── wrpr.go   # todos os novos comandos CLI devem ser referenciados aqui
> │   └── usage.go
> ├── docs/
> ├── frontend/
> ├── internal/
> │   ├── services/
> │   └── types/
> └── utils/
> ```
>
> #### Regras e restrições:
>
> * **Arquivos de licença, segurança, NOTICE, Makefile e suporte devem ser mantidos como estão.**
> * **A estrutura das pastas 'docs', 'support', 'frontend', 'internal', e 'cmd' deve ser mantida.**
> * **Toda implementação de comandos CLI deve ser na pasta 'cmd/cli', com referência obrigatória em 'cmd/wrpr.go'.**
> * **O arquivo 'cmd/main.go' NÃO PODE SER MODIFICADO.**
> * **O código de frontend deve ser React, otimizado para produção, mantido na pasta 'frontend', e nunca perder funcionalidade existente.**
> * **A documentação (docs) deve ser aprimorada apenas quando necessário, sempre em inglês e Markdown, sem alterar a estrutura.**
> * **Toda lógica Go deve estar em 'internal'; manter ou aprimorar existente, criar novos módulos mantendo a estrutura.**
> * **O AGENTS.md deve ser gerado na raiz do projeto, no padrão markdown.**
> * **Arquivos de configuração como `go.mod`, `package.json`, ou outros arquivos essenciais à operação do projeto DEVEM ser alterados conforme necessário para o funcionamento da nova feature.**
> * **Arquivos não listados acima só devem ser alterados se e somente se for estritamente necessário para o funcionamento da feature, e sempre mantendo o restante da funcionalidade do projeto intacta.**
>
> ---
>
> ### Sobre a nova feature:
>
> 1. Implemente um comando CLI (na pasta 'cmd/cli') chamado `squad` ou similar, que recebe requisitos em texto livre e gera um AGENTS.md, conforme padrões de squads multi-agent (roles, skills, stack, limitações, exemplos de prompt).
> 2. O comando deve ser registrado em `cmd/wrpr.go`.
> 3. Use funções/modularização em 'internal' para parsing, construção do AGENTS.md, e validação de input.
> 4. A geração do AGENTS.md deve ser automatizada e idempotente, respeitando inputs e restrições (preferências de stack, granularidade dos roles, exemplos concretos).
> 5. Mantenha tudo o que já existe e só adicione/otimize conforme necessário, SEM sobrescrever lógicas não relacionadas.
> 6. O código precisa ser legível, extensível, e alinhado com o restante do projeto.
> 7. **NÃO criar nem alterar arquivos ou lógicas fora do escopo acima, exceto se realmente necessário para o funcionamento da feature.**
>
> ---
>
> ### Exemplo de uso:
>
> ```
> grompt squad "Quero um microserviço de backend para pagamentos, com autenticação, integração com Stripe, testes automatizados, e deploy em Docker. Prefiro Go ou Python, sem Java."
> ```
>
> Isso deve gerar um AGENTS.md como:
>
> ```
> # Agents
>
> ## Backend Developer (Go/Python)
> - Role: Implement backend services
> - Skills: Go, Python, REST, Stripe API
> - Restrictions: No Java
> - Prompt Example: ...
>
> ## DevOps Engineer
> - Role: Setup Docker and CI/CD pipelines
> - Skills: Docker, CI/CD, Cloud
> - Prompt Example: ...
>
> ## QA Engineer
> - Role: Write and execute automated tests
> - Skills: Automated Testing, Go test, PyTest
> - Prompt Example: ...
> ```
>
> ---
>
> Se precisar de contexto, exemplos ou padrões para AGENTS.md, peça explicitamente antes de tentar implementar.
> Implemente primeiro a versão mínima funcional, depois incremente se solicitado.
>
> **Siga todas as restrições de estrutura. Se não for explicitamente permitido, NÃO altere.**
