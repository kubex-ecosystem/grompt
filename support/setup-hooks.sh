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

  # Create support/pre-commit-config.yaml if it doesn't exist
  if [[ ! -f support/.pre-commit-config.yaml || -z "$(cat support/.pre-commit-config.yaml)" ]]; then
    echo "🛠️  Creating support/.pre-commit-config.yaml..."
    printf '%s\n' '
# This is the default pre-commit configuration file.
# It includes basic hygiene checks, security scans
# Pre-commit configuration file
# Documentation: https://pre-commit.com/
repos:
  # -------- Hygiene básica --------
  - repo: https://github.com/pre-commit/pre-commit-hooks
    rev: v6.0.0
    hooks:
      - id: trailing-whitespace
      - id: end-of-file-fixer
      - id: check-yaml
      - id: check-json
      - id: check-added-large-files
        args: ["--maxkb=1512"]

  # -------- Segurança: detect-secrets --------
  - repo: https://github.com/Yelp/detect-secrets
    rev: v1.5.0
    hooks:
      - id: detect-secrets
        args: ["--baseline", ".secrets.baseline"]
        exclude: '"'"'(^docs/.*$|^README\.md$|^internal/sockets/messagery/rabbitmq\.go$|^frontend/tsconfig\.json$|^docs/swagger\.json$|^frontend/src/locales/.*$)'"'"'
        files: '"'"'^(.*\.go$|.*\.yaml$|.*\.yml$|.*\.json$|.*\.tf$|.*\.tfvars$|.*\.sh$|.*\.ps1$)'"'"'
        pragma: allowlist dsn := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable", pgConfig.Username, pgConfig.Password, pgConfig.Host, port, pgConfig.Name)

  # -------- Segurança: semgrep --------
  # - repo: https://github.com/returntocorp/semgrep
  #   rev: v1.35.0
  #   hooks:
  #     - id: semgrep
  #       args:
  #         [
  #           "--config=auto",
  #           "--exclude=vendor/",
  #           "--exclude=docs/",
  #           "--exclude=support/",
  #         ]
  #       files: '"'"'^(.*\.go$|.*\.yaml$|.*\.yml$|.*\.json$|.*\.tf$|.*\.tfvars$|.*\.sh$|.*\.ps1$)'"'"'

  # -------- Segurança: bandit --------
  # - repo: https://github.com/PyCQA/bandit
  #   rev: v1.7.0
  #   hooks:
  #     - id: bandit
  #       args: ["-r", "."]

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
            "--config=support/.gitleaks.toml",
          ]

  # -------- Go tools --------
  - repo: https://github.com/dnephin/pre-commit-golang
    rev: v0.5.1
    hooks:
      - id: go-fmt
      - id: go-vet
      - id: go-mod-tidy
        # - id: golangci-lint
        # args: ["run", "--config=./support/.golangci.yaml"]
' | tee "support/.pre-commit-config.yaml"
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
  if [[ -f support/requirements-hooks.txt ]]; then
    pip install -r support/requirements-hooks.txt
  fi

  pip install -U pip setuptools wheel
  pip install pre-commit detect-secrets

  # Install pre-commit hooks
  pre-commit install --config support/.pre-commit-config.yaml --install-hooks
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
