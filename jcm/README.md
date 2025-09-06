# kbx - O Maestro da Squad de IAs do Kubex

`kbx` é a ferramenta de linha de comando que serve como ponto de entrada único e abstrato para o ecossistema de desenvolvimento assistido por IA do Kubex. Ele orquestra os diferentes agentes (`codex`, `g-agent`, `copilot`) para garantir que cada contribuição, de qualquer colaborador, esteja alinhada com os princípios e a qualidade definidos no Manifesto Kubex.

> **Nosso lema:** Um comando para começar. A filosofia do projeto para garantir a qualidade.

## 🚀 Início Rápido

### Pré-requisitos

- Go 1.25+
- Make

### Instalação

```bash
# Compila e instala o binário kbx no seu GOPATH/bin
make install
Configuração
O kbx é projetado para funcionar out-of-the-box em um repositório Kubex. Ele procura automaticamente por uma pasta .kubex/ contendo os arquivos de configuração e manifestos.

🛠️ Uso Básico
Verifique se o ambiente está pronto:

Bash

kbx doctor
Inicie uma nova tarefa com um objetivo claro:

Bash

kbx objective "Refatorar o parser de configuração para suportar a nova flag --strict"
Audite suas mudanças atuais contra os princípios do projeto:

Bash

kbx audit
Prepare seu commit, garantindo a conformidade com os padrões:

Bash

kbx commit