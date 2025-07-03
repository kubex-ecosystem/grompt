#!/bin/bash

# 🚀 Script Rápido de Limpeza do Histórico Git
# Versão simplificada para execução direta

echo "🔥 LIMPEZA RÁPIDA DO HISTÓRICO GIT"
echo "================================"
echo ""
echo "⚠️  Esta operação remove TODO o histórico do Git!"
echo "✅ Mantém apenas os arquivos atuais"
echo ""

# Verificar se estamos em um repositório Git
if [ ! -d ".git" ]; then
    echo "❌ Erro: Não está em um repositório Git!"
    exit 1
fi

# Mostrar status atual
echo "📊 Status atual:"
echo "   Branch: $(git branch --show-current)"
echo "   Commits: $(git rev-list --all --count)"
echo ""

# Confirmar
echo "🤔 Para prosseguir, digite exatamente: LIMPAR HISTORICO"
read -p "Confirmação: " confirm

if [ "$confirm" != "LIMPAR HISTORICO" ]; then
    echo "❌ Operação cancelada."
    exit 1
fi

echo ""
echo "🚀 Executando limpeza..."

# Salvar informações
BRANCH=$(git branch --show-current)
REMOTE=$(git remote get-url origin 2>/dev/null || echo "")

# Fazer stash se necessário
if ! git diff-index --quiet HEAD --; then
    git stash push -m "Backup antes da limpeza"
    echo "   💾 Mudanças salvas em stash"
fi

# Remover remote temporariamente
if [ ! -z "$REMOTE" ]; then
    git remote remove origin
fi

# Criar branch órfão
git checkout --orphan temp-clean-branch
echo "   🌿 Branch órfão criado"

# Adicionar todos os arquivos
git add .
echo "   📁 Arquivos adicionados"

# Commit inicial
git commit -m "🎉 Clean history - Security fix

✅ Removed exposed Supabase keys from history
✅ Implemented secure environment configuration  
✅ Fresh start for security compliance

Previous commits removed for security reasons.
Cleaned on: $(date)"

echo "   💾 Commit inicial criado"

# Deletar branch antigo
git branch -D "$BRANCH" 2>/dev/null || true

# Renomear branch
git branch -m temp-clean-branch "$BRANCH"
echo "   🔄 Branch renomeado"

# Reconectar remote
if [ ! -z "$REMOTE" ]; then
    git remote add origin "$REMOTE"
    echo "   🔗 Remote reconectado"
fi

# Aplicar stash se existe
if git stash list | grep -q "Backup antes da limpeza"; then
    git stash pop
    echo "   📝 Mudanças restauradas"
fi

# Limpeza final
git gc --aggressive --prune=now
echo "   🧹 Limpeza de cache concluída"

echo ""
echo "🎉 LIMPEZA CONCLUÍDA!"
echo "==================="
echo ""
echo "📊 Novo status:"
echo "   Branch: $(git branch --show-current)"  
echo "   Commits: $(git rev-list --all --count)"
echo ""
echo "🚀 PRÓXIMO PASSO - Execute o comando:"
echo ""
echo "   git push -f origin $BRANCH"
echo ""
echo "⚠️  IMPORTANTE: Avise a equipe para re-clonar o repositório!"
echo "✅ O histórico foi completamente limpo de vulnerabilidades."
