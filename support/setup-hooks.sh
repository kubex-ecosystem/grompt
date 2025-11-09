#!/usr/bin/env bash
# shellcheck disable=SC2015,SC1091

# Script Metadata
__secure_logic_version="1.0.0"
__secure_logic_date="$( date +%Y-%m-%d )"
__secure_logic_author="Rafael Mori"
__secure_logic_use_type="exec"
__secure_logic_init_timestamp="$(date +%s)"
__secure_logic_elapsed_time=0

set -o errexit # Exit immediately if a command exits with a non-zero status
set -o nounset # Treat unset variables as an error when substituting
set -o pipefail # Return the exit status of the last command in the pipeline that failed
set -o errtrace # If a command fails, the shell will exit immediately
set -o functrace # If a function fails, the shell will exit immediately
shopt -s inherit_errexit # Inherit the errexit option in functions

# Get the root directory of the git project
_SCRIPT_DIR="$(git rev-parse --show-toplevel)"
cd "$_SCRIPT_DIR" || exit 1

_default_pre_commit_config() {
  echo "🚀 Configurando pre-commit hooks (defaults)..."

  # Create support/hooks/pre-commit-config.yaml if it doesn't exist
  if [[ ! -f support/hooks/.pre-commit-config.yaml || -z "$(cat support/hooks/.pre-commit-config.yaml)" ]]; then
    echo "🛠️  Creating support/hooks/.pre-commit-config.yaml..."
    printf '%s\n' '
# Pre-commit configuration file

repos:
  - repo: local
    hooks:
      - id: docs-sanitize
        name: docs-sanitize (deny secrets-looking strings in docs)
        entry: python3 support/hooks/.check_docs_secrets.py
        language: system
        files: \.(md|mdx|rst|txt|adoc|html)$
        fail_fast: false
        always_run: true

  # -------- Hygiene básica --------
  - repo: https://github.com/pre-commit/pre-commit-hooks
    rev: v6.0.0
    hooks:
      - id: trailing-whitespace
      - id: end-of-file-fixer
      - id: check-yaml
      - id: check-json
      - id: check-added-large-files
      - id: check-merge-conflicts
        args: ["--maxkb=1512"]
        fail_fast: false
        always_run: true
        verbose: false

  # -------- Formatação de código --------
  - repo: https://github.com/pre-commit/mirrors-eslint
    rev: "v8.45.0"
    hooks:
      - id: eslint
        additional_dependencies:
          - eslint@8.45.0
          - typescript-eslint/eslint-plugin@5.61.0
          - eslint-plugin-react@7.32.2
          - eslint-plugin-react-hooks@4.6.0
          - eslint-plugin-import@2.26.0
          - eslint-plugin-jsx-a11y@6.7.1
          - babel-eslint@10.1.0
        name: eslint (fix)
        entry: eslint
        language: system
        fail_fast: false
        always_run: true
        verbose: true
        files: \.(js|jsx|ts|tsx)$
        args:
          [
            "--config",
            "support/hooks/.eslintrc.json",
            "--ext",
            ".js,.jsx,.ts,.tsx,.mjs,.cjs,.json",
            "--fix"
          ]

  # -------- Documentação --------
  - repo: https://github.com/pre-commit/mirrors-markdownlint
    rev: v0.1.0
    hooks:
      - id: markdownlint
        name: markdownlint (fix)
        files: \.(md|mdx)$
        args: ["--config", "support/hooks/.markdownlint.yaml"]
        additional_dependencies:
          - markdownlint-cli@0.33.0
        entry: markdownlint-cli
        language: system
        fail_fast: true
        verbose: true
        always_run: true

  # -------- Segurança: semgrep --------
  - repo: https://github.com/returntocorp/semgrep
    rev: v1.35.0
    hooks:
      - id: semgrep
        name: semgrep (security and code quality)
        entry: semgrep
        language: system
        fail_fast: true
        pass_filenames: false
        always_run: true
        files: ^(.*\.go$|.*\.yaml$|.*\.yml$|.*\.json$|.*\.tf$|.*\.tfvars$|.*\.sh$|.*\.ps1$|.*\.ts$|.*\.tsx$|.*\.js$|.*\.jsx$)
        args:
          [
            "--exclude=vendor/",
            "--exclude=node_modules/",
            "--exclude=dist/",
            "--exclude=build/",
            "--exclude=coverage/",
            "--exclude=support/",
            "--exclude=tests/",
            "--timeout=120s",
            "--config=support/hooks/.semgrep.yaml"
          ]

  # -------- Segurança: detect-secrets --------
  - repo: https://github.com/Yelp/detect-secrets
    rev: v1.5.0
    hooks:
      - id: detect-secrets
        args: ["--baseline", ".secrets.baseline"]
        exclude: (^docs/.*$|^README\.md$|^internal/sockets/messagery/rabbitmq\.go$|^frontend/tsconfig\.json$|^docs/swagger\.json$|^frontend/src/locales/.*$)
        files: ^(.*\.go$|.*\.yaml$|.*\.yml$|.*\.json$|.*\.tf$|.*\.tfvars$|.*\.sh$|.*\.ps1$|.*\.ts$|.*\.tsx$|.*\.js$|.*\.jsx$)
        pragma: allowlist dsn := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable",pgConfig.Username, pgConfig.Password, pgConfig.Host, port, pgConfig.Name)
        fail_fast: false
        always_run: true
        verbose: true

  # -------- Segurança: gitleaks --------
  - repo: https://github.com/zricethezav/gitleaks
    rev: v8.28.0
    hooks:
      - id: gitleaks
        name: gitleaks
        entry: gitleaks
        language: system
        pass_filenames: false
        args:
          [
            "protect",
            "--staged",
            "--no-banner",
            "--config=support/hooks/.gitleaks.toml"
          ]
        files: ^(.*\.go$|.*\.yaml$|.*\.yml$|.*\.json$|.*\.tf$|.*\.tfvars$|.*\.sh$|.*\.ps1$|.*\.ts$|.*\.tsx$|.*\.js$|.*\.jsx$)
        verbose: true
        fail_fast: false
        always_run: true

' | tee "support/hooks/.pre-commit-config.yaml"
  # else
  #   cat support/pre-commit-config.yaml
  fi
}

_install_pre_commit_tools() {
  echo "🚀 Configurando pre-commit hooks (installs)..."
  # Create and activate a virtual environment for hooks
  if [[ ! -d .venv-hooks ]]; then
    python3 -m venv .venv-hooks
    echo '.venv-hooks' >> .gitignore
  fi
  if [[ ! -f .venv-hooks/bin/activate ]]; then
    echo "❌ Falha ao encontrar o ambiente virtual em .venv-hooks"
    exit 1
  fi

  # shellcheck source=/dev/null
  . .venv-hooks/bin/activate

  # Install requirements file from support/ if it exists
  if [[ -f support/hooks/requirements-hooks.txt ]]; then
    pip install -r support/hooks/requirements-hooks.txt
  fi

  pip install -U pip setuptools wheel
  pip install pre-commit detect-secrets

  # Install pre-commit hooks
  pre-commit install --config support/hooks/.pre-commit-config.yaml --install-hooks
}

_create_baseline() {
  # Create a baseline for detect-secrets if it doesn't exist
  if [[ ! -f .secrets.baseline ]]; then
    detect-secrets scan > .secrets.baseline
    git add .secrets.baseline
    git commit -m "chore(secrets): add baseline" || true
  fi
}

_main() {
  # First we check if pre-commit is already configured
  if git config --get core.hooksPath &>/dev/null; then
    echo "⚠️  Pre-commit hooks are already configured. Aborting..."
    return 0
  fi

  _default_pre_commit_config

  _install_pre_commit_tools

  _create_baseline

  echo "✅ Pre-commit hooks configured successfully!"
}

_main "$@"

__secure_logic_elapsed_time=$(( $(date +%s) - __secure_logic_init_timestamp ))
echo "⏱️  Script executed in $__secure_logic_elapsed_time seconds."
