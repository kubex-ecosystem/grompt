# Makefile for Grompt Documentation

.PHONY: help docs-install docs-serve docs-build docs-clean docs-deploy

# Variables
DOCS_DIR = support/docs
VENV_DIR = $(DOCS_DIR)/.venv
PYTHON = $(VENV_DIR)/bin/python
PIP = $(VENV_DIR)/bin/pip

help: ## Show this help message
	@echo "Grompt Documentation Commands:"
	@echo ""
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-20s\033[0m %s\n", $$1, $$2}'

docs-install: ## Install documentation dependencies
	@echo "🔧 Installing documentation dependencies..."
	cd $(DOCS_DIR) && python3 -m venv .venv
	$(PIP) install --upgrade pip
	$(PIP) install mkdocs mkdocs-material
	$(PIP) install mkdocs-git-revision-date-localized-plugin
	$(PIP) install pymdown-extensions
	@echo "✅ Dependencies installed!"

docs-serve: ## Serve documentation locally
	@echo "🚀 Starting documentation server..."
	@echo "📖 Documentation will be available at: http://127.0.0.1:8000"
	cd $(DOCS_DIR) && $(PYTHON) -m mkdocs serve

docs-build: ## Build documentation
	@echo "🔨 Building documentation..."
	cd $(DOCS_DIR) && $(PYTHON) -m mkdocs build
	@echo "✅ Documentation built in $(DOCS_DIR)/site/"

docs-clean: ## Clean documentation build
	@echo "🧹 Cleaning documentation build..."
	cd $(DOCS_DIR) && rm -rf site/
	@echo "✅ Documentation build cleaned!"

docs-deploy: ## Deploy documentation to GitHub Pages
	@echo "🚀 Deploying documentation to GitHub Pages..."
	cd $(DOCS_DIR) && $(PYTHON) -m mkdocs gh-deploy --force
	@echo "✅ Documentation deployed!"

docs-check: ## Check documentation for issues
	@echo "🔍 Checking documentation..."
	cd $(DOCS_DIR) && $(PYTHON) -m mkdocs build --strict
	@echo "✅ Documentation check completed!"

docs-new-page: ## Create a new documentation page (usage: make docs-new-page PAGE=path/to/page.md)
	@if [ -z "$(PAGE)" ]; then \
		echo "❌ Usage: make docs-new-page PAGE=path/to/page.md"; \
		exit 1; \
	fi
	@echo "📝 Creating new page: $(PAGE)"
	@mkdir -p $(DOCS_DIR)/docs/$(dir $(PAGE))
	@echo "# $(shell basename $(PAGE) .md | tr '-' ' ' | sed 's/\b\w/\u&/g')" > $(DOCS_DIR)/docs/$(PAGE)
	@echo "" >> $(DOCS_DIR)/docs/$(PAGE)
	@echo "Esta página está em desenvolvimento." >> $(DOCS_DIR)/docs/$(PAGE)
	@echo "✅ Page created: $(DOCS_DIR)/docs/$(PAGE)"

docs-setup: docs-install ## Setup complete documentation environment
	@echo "🏗️  Setting up documentation environment..."
	@echo "✅ Documentation environment ready!"
	@echo ""
	@echo "Next steps:"
	@echo "  1. Run 'make docs-serve' to start local server"
	@echo "  2. Edit files in $(DOCS_DIR)/docs/"
	@echo "  3. Run 'make docs-build' to build for production"

docs-all: docs-clean docs-build ## Clean and build documentation

# Development helpers
docs-watch: ## Watch for changes and rebuild (requires entr)
	@echo "👀 Watching for changes..."
	find $(DOCS_DIR)/docs -name "*.md" | entr -r make docs-build

docs-lint: ## Lint markdown files (requires markdownlint)
	@echo "🔍 Linting markdown files..."
	@if command -v markdownlint >/dev/null 2>&1; then \
		markdownlint $(DOCS_DIR)/docs/**/*.md; \
		echo "✅ Markdown linting completed!"; \
	else \
		echo "⚠️  markdownlint not found. Install with: npm install -g markdownlint-cli"; \
	fi

docs-spellcheck: ## Spell check documentation (requires aspell)
	@echo "📝 Spell checking documentation..."
	@if command -v aspell >/dev/null 2>&1; then \
		find $(DOCS_DIR)/docs -name "*.md" -exec aspell check {} \;; \
		echo "✅ Spell check completed!"; \
	else \
		echo "⚠️  aspell not found. Install with your package manager"; \
	fi

# Quick commands
docs: docs-serve ## Alias for docs-serve
build-docs: docs-build ## Alias for docs-build
clean-docs: docs-clean ## Alias for docs-clean

# Info commands
docs-info: ## Show documentation info
	@echo "📊 Documentation Information:"
	@echo "  Documentation directory: $(DOCS_DIR)"
	@echo "  Virtual environment: $(VENV_DIR)"
	@echo "  Python executable: $(PYTHON)"
	@echo ""
	@if [ -d "$(VENV_DIR)" ]; then \
		echo "  Status: ✅ Environment ready"; \
		echo "  MkDocs version: $$($(PYTHON) -c 'import mkdocs; print(mkdocs.__version__)' 2>/dev/null || echo 'Not installed')"; \
	else \
		echo "  Status: ❌ Environment not set up"; \
		echo "  Run 'make docs-install' to set up"; \
	fi
	@echo ""
	@if [ -d "$(DOCS_DIR)/site" ]; then \
		echo "  Built documentation: ✅ Available in $(DOCS_DIR)/site/"; \
	else \
		echo "  Built documentation: ❌ Not built yet"; \
	fi

docs-urls: ## Show important documentation URLs
	@echo "🔗 Documentation URLs:"
	@echo "  Local development: http://127.0.0.1:8000"
	@echo "  GitHub Pages: https://rafa-mori.github.io/grompt/"
	@echo "  GitHub Repository: https://github.com/rafa-mori/grompt"
	@echo "  Issues: https://github.com/rafa-mori/grompt/issues"
