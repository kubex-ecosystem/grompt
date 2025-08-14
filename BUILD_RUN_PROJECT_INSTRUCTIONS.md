# Instruções de como compilar e executar o projeto

## Absolute Real Path

O diretório REAL e absoluto desse projeto é o:

```bash
cd '/srv/apps/LIFE/PROJECTS/grompt'
```

## Build

Para compilar o projeto, execute o seguinte comando na raiz do diretório do projeto:

```bash
FORCE=y make build-dev
```

Esse comando acima irá gerar um arquivo binário chamado lookatni no diretório `./dist`. OBSERVE BEM A DIFERENÇA DA SAÍDA ENTRE OS PROJETOS. Esse projeto está no dist por vários motivos e já há planos para normalização dos módulos que estamos criando, mas isso virá em ooooutra task.. rsrs

O `FORCE=y` é necessário porque o fluxo todo de distribuição foi projetado pra atender tanto esteiras automatizadas quando users, então se já houver binário com o mesmo nome no diretório de destino, ele, num ambiente interativo como estamos irá perguntar se deseja substituir ou não, o que te impediria de concluir o build em função de só haver interação sua no shell para enviar comandos e não exatamente como seria a de um user.

Então, o próximo passo, naturalmente...

## Executar

Para executar o projeto, utilize o seguinte comando:

```bash
 ./dist/grompt -h # Ou alguma outra flag que prefira na hora...
```
