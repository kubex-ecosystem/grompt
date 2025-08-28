# Solução de Problemas

Guia para resolver problemas comuns do Grompt.

## 🔧 Problemas de Instalação

### Erro: "Permission denied"

```bash
chmod +x grompt
```

### Erro: "Command not found"

```bash
# Verificar se está no PATH
echo $PATH
which grompt

# Mover para diretório no PATH
sudo mv grompt /usr/local/bin/
```

### Erro: "Port already in use"

```bash
# Verificar qual processo usa a porta
lsof -i :8080

# Usar porta diferente
grompt --port 8081
```

## 🌐 Problemas de Conectividade

### API key não funciona

```bash
# Testar conectividade
grompt ask "teste" --provider openai --dry-run

# Verificar variáveis de ambiente
env | grep -E "(OPENAI|CLAUDE|GEMINI)"
```

### Timeout de conexão

```bash
# Verificar conectividade
curl -I https://api.openai.com

# Configurar proxy se necessário
export HTTP_PROXY="http://proxy:8080"
```

## 📝 Problemas da Interface

### Página não carrega

1. Verificar se o servidor está rodando
2. Checar logs de erro
3. Limpar cache do navegador

### Interface lenta

1. Verificar uso de CPU/memória
2. Reduzir número de ideias simultâneas
3. Usar provedor mais rápido

## 🚀 Problemas de Performance

### Alto uso de memória

```bash
# Verificar uso de recursos
top | grep grompt
```

### Respostas lentas

1. Usar modelos menores
2. Reduzir max_tokens
3. Verificar latência de rede

---

Em desenvolvimento: mais soluções serão adicionadas.
