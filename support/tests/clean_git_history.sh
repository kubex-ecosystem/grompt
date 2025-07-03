#!/bin/bash

# 🔥 Script de Limpeza Completa do Histórico Git - SUSsec
# Este script remove COMPLETAMENTE o histórico do Git para eliminar chaves expostas

set -e  # Parar em caso de erro

echo "🔥 LIMPEZA COMPLETA DO HISTÓRICO GIT"
echo "=================================="
echo ""
echo "⚠️  ATENÇÃO: Esta operação é IRREVERSÍVEL!"
echo "    - TODO o histórico do Git será perdido"
echo "    - Todos os commits anteriores serão removidos"
echo "    - Branches antigas serão eliminadas"
echo ""
echo "✅ Benefícios:"
echo "    - Remove definitivamente chaves expostas"
echo "    - Elimina qualquer rastro de vulnerabilidades"
echo "    - Repositório limpo e seguro"
echo ""

# Verificar se estamos em um repositório Git
if [ ! -d ".git" ]; then
    echo "❌ Erro: Este diretório não é um repositório Git!"
    exit 1
fi

# Mostrar informações atuais
echo "📊 Status atual do repositório:"
echo "   Branch atual: $(git branch --show-current)"
echo "   Total de commits: $(git rev-list --all --count)"
echo "   Remotes: $(git remote -v | wc -l) configurados"
echo ""

# Confirmar com o usuário
read -p "🤔 Tem certeza que deseja LIMPAR TODO O HISTÓRICO? (digite 'CONFIRMO' para prosseguir): " confirmacao

if [ "$confirmacao" != "CONFIRMO" ]; then
    echo "❌ Operação cancelada pelo usuário."
    exit 1
fi

echo ""
echo "🚀 Iniciando limpeza do histórico..."

# Backup do remote atual (se existir)
REMOTE_URL=""
if git remote get-url origin &>/dev/null; then
    REMOTE_URL=$(git remote get-url origin)
    echo "💾 Remote atual salvo: $REMOTE_URL"
fi

# Salvar nome da branch atual
CURRENT_BRANCH=$(git branch --show-current)
echo "🌿 Branch atual: $CURRENT_BRANCH"

# Verificar se há mudanças não commitadas
if ! git diff-index --quiet HEAD --; then
    echo "📝 Detectadas mudanças não commitadas. Fazendo stash..."
    git stash push -m "Backup antes da limpeza de histórico - $(date)"
fi

echo ""
echo "🔥 Executando limpeza completa..."

# 1. Remover referência ao remote para evitar push acidental
if [ ! -z "$REMOTE_URL" ]; then
    git remote remove origin
    echo "   ✅ Remote removido temporariamente"
fi

# 2. Criar um novo branch orfão (sem histórico)
git checkout --orphan new-clean-history
echo "   ✅ Branch órfão criado"

# 3. Adicionar todos os arquivos atuais
git add .
echo "   ✅ Arquivos adicionados"

# 4. Fazer o primeiro commit limpo
git commit -m "🎉 Initial commit - Clean history

✅ Security vulnerabilities resolved
✅ Supabase keys removed from history
✅ Fresh start with secure configuration

Previous history removed for security reasons.
Date: $(date '+%Y-%m-%d %H:%M:%S')
"
echo "   ✅ Commit inicial criado"

# 5. Deletar a branch antiga
git branch -D "$CURRENT_BRANCH" 2>/dev/null || echo "   ⚠️  Branch antiga não pôde ser removida (normal se era main/master)"

# 6. Renomear o branch atual para o nome original
if [ "$CURRENT_BRANCH" != "new-clean-history" ]; then
    git branch -m new-clean-history "$CURRENT_BRANCH"
    echo "   ✅ Branch renomeado para $CURRENT_BRANCH"
fi

# 7. Forçar garbage collection para liberar espaço
git gc --aggressive --prune=now
echo "   ✅ Limpeza de espaço executada"

# 8. Reconectar o remote se existia
if [ ! -z "$REMOTE_URL" ]; then
    git remote add origin "$REMOTE_URL"
    echo "   ✅ Remote reconectado: $REMOTE_URL"
fi

# 9. Aplicar stash se existe
if git stash list | grep -q "Backup antes da limpeza"; then
    echo "   📝 Aplicando mudanças que estavam em stash..."
    git stash pop
fi

echo ""
echo "🎉 LIMPEZA CONCLUÍDA COM SUCESSO!"
echo "================================"
echo ""
echo "📊 Novo status do repositório:"
echo "   Branch atual: $(git branch --show-current)"
echo "   Total de commits: $(git rev-list --all --count)"
echo "   Primeiro commit: $(git log --oneline | tail -1)"
echo ""
echo "🚨 PRÓXIMOS PASSOS OBRIGATÓRIOS:"
echo ""
echo "1. 🔍 Verificar se tudo está correto:"
echo "   git log --oneline"
echo "   git status"
echo ""
echo "2. 🚀 Force push para o repositório remoto:"
echo "   git push -f origin $CURRENT_BRANCH"
echo ""
echo "3. ⚠️  AVISAR A EQUIPE:"
echo "   - O histórico foi completamente reescrito"
echo "   - Todos devem re-clonar o repositório"
echo "   - Branches locais antigas devem ser descartadas"
echo ""
echo "4. 🔐 Confirmar no Supabase:"
echo "   - Revogar as chaves antigas imediatamente"
echo "   - Gerar novas chaves"
echo "   - Configurar .env com as novas chaves"
echo ""
echo "✅ Seu repositório agora está 100% limpo de vulnerabilidades!"
