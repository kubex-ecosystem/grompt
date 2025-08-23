# Description: Makefile for building and installing a Go application
# Author: Rafael Mori
# Copyright (c) 2025 Rafael Mori
# License: MIT License

# Define the application name and root directory
ROOT_DIR := $(dir $(abspath $(lastword $(MAKEFILE_LIST))))
TARGET_MANIFEST = $(ROOT_DIR)internal/module/info/manifest.json
APP_NAME := $(shell jq -r '.name' < $(TARGET_MANIFEST))

ifeq ($(APP_NAME),)
APP_NAME := $(shell  echo $(basename $(CURDIR)) | tr '[:upper:]' '[:lower:]')
endif
ORGANIZATION := $(shell jq -r '.organization' < $(TARGET_MANIFEST))
ifeq ($(ORGANIZATION),)
ORGANIZATION := $(shell git config --get user.name | tr '[:upper:]' '[:lower:]')
endif
ifeq ($(ORGANIZATION),)
ORGANIZATION := $(shell git config --get user.email | cut -d '@' -f 1 | tr '[:upper:]' '[:lower:]')
endif
ifeq ($(ORGANIZATION),)
ORGANIZATION := $(shell echo $(USER) | tr '[:upper:]' '[:lower:]')
endif
REPOSITORY := $(shell jq -r '.repository' < $(TARGET_MANIFEST))
ifeq ($(REPOSITORY),)
REPOSITORY := $(shell git config --get remote.origin.url)
endif
ifeq ($(REPOSITORY),)
REPOSITORY := $(shell git config --get remote.upstream.url)
endif
ifeq ($(REPOSITORY),)
REPOSITORY := $(printf 'https://github.com/%s/%s.git' $(ORGANIZATION) $(APP_NAME))
endif
DESCRIPTION := $(shell jq -r '.description' < $(TARGET_MANIFEST))
ifeq ($(DESCRIPTION),)
DESCRIPTION := $(shell git log -1 --pretty=%B | head -n 1)
endif
BINARY_NAME := $(shell jq -r '.bin' < $(TARGET_MANIFEST))
ifeq ($(BINARY_NAME),)
BINARY_NAME := $(ROOT_DIR)bin/$(APP_NAME)
else
BINARY_NAME := $(ROOT_DIR)bin/$(BINARY_NAME)
endif
CMD_DIR := $(ROOT_DIR)cmd

# Define the color codes
COLOR_GREEN := \033[32m
COLOR_YELLOW := \033[33m
COLOR_RED := \033[31m
COLOR_BLUE := \033[34m
COLOR_RESET := \033[0m

# Logging Functions
log = @printf "%b%s%b %s\n" "$(COLOR_BLUE)" "[LOG]" "$(COLOR_RESET)" "$(1)"
log_info = @printf "%b%s%b %s\n" "$(COLOR_BLUE)" "[INFO]" "$(COLOR_RESET)" "$(1)"
log_success = @printf "%b%s%b %s\n" "$(COLOR_GREEN)" "[SUCCESS]" "$(COLOR_RESET)" "$(1)"
log_warning = @printf "%b%s%b %s\n" "$(COLOR_YELLOW)" "[WARNING]" "$(COLOR_RESET)" "$(1)"
log_break = @printf "%b%s%b\n" "$(COLOR_BLUE)" "[INFO]" "$(COLOR_RESET)"
log_error = @printf "%b%s%b %s\n" "$(COLOR_RED)" "[ERROR]" "$(COLOR_RESET)" "$(1)"

ARGUMENTS := $(MAKECMDGOALS)
INSTALL_SCRIPT = $(ROOT_DIR)support/main.sh
CMD_STR := $(strip $(firstword $(ARGUMENTS)))
ARGS := $(filter-out $(strip $(CMD_STR)), $(ARGUMENTS))

# Default target: help
.DEFAULT_GOAL := help

# Build the binary using the install script.
build:
	@bash $(INSTALL_SCRIPT) build $(ARGS)
	$(shell exit 0)

build-dev:
	@bash $(INSTALL_SCRIPT) build-dev $(ARGS)
	$(shell exit 0)

# Install the binary and configure the environment.
install:
	@bash $(INSTALL_SCRIPT) install $(ARGS)
	$(shell exit 0)

# Uninstall the binary and clean up.
uninstall:
	@bash $(INSTALL_SCRIPT) uninstall $(ARGS)
	$(shell exit 0)

# Clean up build artifacts.
clean:
	@bash $(INSTALL_SCRIPT) clean $(ARGS)
	$(shell exit 0)

# Run tests.
test:
	@bash $(INSTALL_SCRIPT) test $(ARGS)
	$(shell exit 0)

## Run dynamic commands with arguments calling the install script.
%:
	@:
	$(call log_info, Running command: $(CMD_STR))
	$(call log_info, Args: $(ARGS))
	@bash $(INSTALL_SCRIPT) $(CMD_STR) $(ARGS)
	$(shell exit 0)

# Display help message
help:
	$(call log, $(APP_NAME) Build System Help)
	$(call log_break)
	$(call log, Available targets:)
	$(call log,   make build      - Build the binary using install script)
	$(call log,   make install    - Install the binary and configure environment)
	$(call log,   make docs       - Start API documentation server)
	$(call log,   make build-docs - Build documentation server binary)
	$(call log,   make clean      - Clean up build artifacts)
	$(call log,   make test       - Run tests)
	$(call log,   make help       - Display this help message)
	$(call log_break)
	$(call log, Documentation:)
	$(call log,   make docs       - Starts beautiful API documentation at http://localhost:8080/docs)
	$(call log,   ./start-docs.sh - Alternative way to start documentation server)
	$(call log_break)
	$(call log, Usage with arguments:)
	$(call log,   make install ARGS='--custom-arg value' - Pass custom arguments to the install script)
	$(call log_break)
	$(call log, Example:)
	$(call log,   make install ARGS='--prefix /usr/local')
	$(call log_break)
	$(call log, Description:)
	$(call log,   $(DESCRIPTION))
	$(call log_break)
	$(call log, For more information, visit:)
	$(call log,  $(REPOSITORY))
	$(call log_break)
	$(call log_success, End of help message)
	$(shell exit 0)
