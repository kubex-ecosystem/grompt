#!/usr/bin/env bash
set -euo pipefail

# 1) Delete workflow runs (with filters)
# usage:
#   delete_workflow_runs "<name_or_file.yml|all>" <limit> <days> [--dry-run]
# ex: delete_workflow_runs "build.yml" 100 30 --dry-run
delete_workflow_runs() {
  local wf="${1:-all}" limit="${2:-200}" days="${3:-0}" dry="${4:-}"
  local wfFlag=(); [[ "$wf" != "all" ]] && wfFlag=(--workflow "$wf")
  local jqFilter='.[]
    | select(.status!="in_progress" and .status!="queued")
    | .databaseId'
  # filter by age (createdAt < now - days)
  if [[ "$days" -gt 0 ]]; then
    jqFilter=".[]
      | select(.status!=\"in_progress\" and .status!=\"queued\")
      | select((.createdAt | fromdateiso8601) < (now - (${days}*86400)))
      | .databaseId"
  fi
  mapfile -t ids < <(gh run list "${wfFlag[@]}" --limit "$limit" \
    --json databaseId,status,createdAt --jq "$jqFilter")
  [[ "${#ids[@]}" -eq 0 ]] && { echo "Nothing to delete."; return 0; }
  echo "Runs to delete (${#ids[@]}): ${ids[*]}"
  [[ "$dry" == "--dry-run" ]] && return 0
  printf '%s\n' "${ids[@]}" | xargs -r -n1 gh run delete
}

# 2) Delete releases + tags (remote and local)
delete_releases_and_tags() {
  # releases
  mapfile -t rels < <(gh release list --limit 1000 | awk '{print $1}')
  if [[ "${#rels[@]}" -gt 0 ]]; then
    printf '%s\n' "${rels[@]}" | xargs -r -n1 gh release delete -y
  fi
  # remote and local tags
  mapfile -t tags < <(git tag -l)
  if [[ "${#tags[@]}" -gt 0 ]]; then
    printf '%s\n' "${tags[@]}" | xargs -r -n1 git push --delete origin || true
    printf '%s\n' "${tags[@]}" | xargs -r -n1 git tag -d
  fi
}

# 3) Delete Actions artifacts (frees up space)
#   requires gh CLI >= v2.60
purge_artifacts() {
  mapfile -t ids < <(gh api repos/{owner}/{repo}/actions/artifacts --paginate \
    --jq '.artifacts[].id')
  [[ "${#ids[@]}" -eq 0 ]] && { echo "No artifacts."; return 0; }
  printf '%s\n' "${ids[@]}" | xargs -r -n1 -I{} gh api \
    -X DELETE repos/{owner}/{repo}/actions/artifacts/{}
}

# 4) Delete Actions caches (needs official extension)
#   gh extension install actions/gh-actions-cache
purge_caches() {
  if ! gh actions-cache --help &>/dev/null; then
    echo "Installing gh-actions-cache extension..."
    gh extension install actions/gh-actions-cache
  fi
  gh actions-cache list | awk 'NR>1 {print $1}' | xargs -r -n1 \
    gh actions-cache delete --confirm
}

# 5) Secrets scan (gitleaks + trufflehog)
scan_secrets() {
  command -v gitleaks >/dev/null || {
    curl -sSL https://raw.githubusercontent.com/gitleaks/gitleaks/master/install.sh | bash
  }
  gitleaks detect --no-git --redact -v || true
  command -v trufflehog >/dev/null || pipx install trufflehog || true
  trufflehog filesystem . || true
}

# 6) Your original helper (kept, with minor adjustments)
# usage: delete_workflow_history "<workflow>" <limit>
delete_workflow_history() {
  local wf="${1:-none}" qty="${2:-20}"
  gh run list --workflow "$wf" --json databaseId --limit "$qty" \
    | jq -e '.[].databaseId' | xargs -r -n1 gh run delete
}

# 7) (Optional) Full history reset
#    -> use the dedicated script “clean_repo.sh” (already ready and commented)  ←
#       It creates an orphan branch, makes the first commit and gives forced push instructions.
#       Read before using (destructive operation).
