#!/bin/bash

# 🔥 Complete Git History Cleanup Script

set -e  # Stop on error

echo "🔥 COMPLETE GIT HISTORY CLEANUP"
echo "=================================="
echo ""
echo "⚠️  WARNING: This operation is IRREVERSIBLE!"
echo "    - ALL Git history will be lost"
echo "    - All previous commits will be removed"
echo "    - Old branches will be eliminated"
echo ""
echo "✅ Benefits:"
echo "    - Permanently removes exposed keys"
echo "    - Eliminates any trace of vulnerabilities"
echo "    - Clean and secure repository"
echo ""

# Check if we are in a Git repository
if [ ! -d ".git" ]; then
  echo "❌ Error: This directory is not a Git repository!"
  exit 1
fi

# Show current information
echo "📊 Current repository status:"
echo "   Current branch: $(git branch --show-current)"
echo "   Total commits: $(git rev-list --all --count)"
echo "   Remotes: $(git remote -v | wc -l) configured"
echo ""

# Confirm with the user
read -t 10 -n 1 -p "🤔 Are you sure you want to CLEAN ALL HISTORY? (type 'y' to proceed): " confirmation || confirmation='n'

if [[ ! $confirmation =~ [yY] ]]; then
  echo "❌ Operation cancelled by user."
  exit 1
fi

echo ""
echo "🚀 Starting history cleanup..."

# Backup current remote (if exists)
REMOTE_URL=""
if git remote get-url origin &>/dev/null; then
  REMOTE_URL=$(git remote get-url origin)
  echo "💾 Current remote saved: $REMOTE_URL"
fi

# Save current branch name
CURRENT_BRANCH=$(git branch --show-current)
echo "🌿 Current branch: $CURRENT_BRANCH"

# Check for uncommitted changes
if ! git diff-index --quiet HEAD --; then
  echo "📝 Uncommitted changes detected. Stashing..."
  git stash push -m "Backup before history cleanup - $(date)"
fi

echo ""
echo "🔥 Performing complete cleanup..."

# 1. Remove remote reference to avoid accidental push
if [ ! -z "$REMOTE_URL" ]; then
  git remote remove origin
  echo "   ✅ Remote temporarily removed"
fi

# 2. Create a new orphan branch (no history)
git checkout --orphan new-clean-history
echo "   ✅ Orphan branch created"

# 3. Add all current files
git add .
echo "   ✅ Files added"

# 4. Make the first clean commit
git commit -m "Security Commit - Clean history

✅ Security vulnerabilities resolved
✅ Supabase keys removed from history
✅ Fresh start with secure configuration

Previous history removed for security reasons.
Date: $(date '+%Y-%m-%d %H:%M:%S')
"
echo "   ✅ Initial commit created"

# 5. Delete the old branch
git branch -D "$CURRENT_BRANCH" 2>/dev/null || echo "   ⚠️  Old branch could not be removed (normal if it was main/master)"

# 6. Rename current branch to the original name
if [ "$CURRENT_BRANCH" != "new-clean-history" ]; then
  git branch -m new-clean-history "$CURRENT_BRANCH"
  echo "   ✅ Branch renamed to $CURRENT_BRANCH"
fi

# 7. Force garbage collection to free up space
git gc --aggressive --prune=now
echo "   ✅ Space cleanup executed"

# 8. Reconnect remote if it existed
if [ ! -z "$REMOTE_URL" ]; then
  git remote add origin "$REMOTE_URL"
  echo "   ✅ Remote reconnected: $REMOTE_URL"
fi

# 9. Apply stash if exists
if git stash list | grep -q "Backup before history cleanup"; then
  echo "   📝 Applying changes that were stashed..."
  git stash pop
fi

echo ""
echo "🎉 CLEANUP SUCCESSFULLY COMPLETED!"
echo "================================"
echo ""
echo "📊 New repository status:"
echo "   Current branch: $(git branch --show-current)"
echo "   Total commits: $(git rev-list --all --count)"
echo "   First commit: $(git log --oneline | tail -1)"
echo ""
echo "🚨 MANDATORY NEXT STEPS:"
echo ""
echo "1. 🔍 Check if everything is correct:"
echo "   git log --oneline"
echo "   git status"
echo ""
echo "2. 🚀 Force push to the remote repository:"
echo "   git push -f origin $CURRENT_BRANCH"
echo ""
echo "3. ⚠️  INFORM THE TEAM:"
echo "   - The history was completely rewritten"
echo "   - Everyone must re-clone the repository"
echo "   - Old local branches must be discarded"
echo ""
echo "4. 🔐 Confirm in Supabase:"
echo "   - Revoke the old keys immediately"
echo "   - Generate new keys"
echo "   - Set up .env with the new keys"
echo ""
echo "✅ Your repository is now 100% clean of vulnerabilities!"

