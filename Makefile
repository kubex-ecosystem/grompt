# Description: Makefile for building and installing a Go application
# Author: Rafael Mori
# Copyright (c) 2025 Rafael Mori
# License: MIT License

# This Makefile is used to build and install a Go application.
# It provides commands for building the binary, installing it, cleaning up build artifacts,
# and running tests. It also includes a help command to display usage information.
# The Makefile uses color codes for logging messages and provides a consistent interface
# for interacting with the application.

# Define the application name and root directory
APP_NAME := $(shell echo $(basename $(CURDIR)) | tr '[:upper:]' '[:lower:]')
ROOT_DIR := $(dir $(abspath $(lastword $(MAKEFILE_LIST))))
BINARY_NAME := $(ROOT_DIR)$(APP_NAME)
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
INSTALL_SCRIPT=$(ROOT_DIR)support/install.sh
CMD_STR := $(strip $(firstword $(ARGUMENTS)))
ARGS := $(filter-out $(strip $(CMD_STR)), $(ARGUMENTS))

# Build the binary using the install script.
build:
	$(call log_info, Building $(APP_NAME) binary)
	$(call log_info, Args: $(ARGS))
	@bash $(INSTALL_SCRIPT) build $(ARGS)
	$(shell exit 0)

# Install the binary and configure the environment.
install:
	$(call log_info, Installing $(APP_NAME) binary)
	$(call log_info, Args: $(ARGS))
	@bash $(INSTALL_SCRIPT) install $(ARGS)
	$(shell exit 0)

# Clean up build artifacts.
clean:
	$(call log_info, Cleaning up build artifacts)
	$(call log_info, Args: $(ARGS))
	@bash $(INSTALL_SCRIPT) clean $(ARGS)
	$(shell exit 0)

# Run tests.
test:
	$(call log_info, Running tests)
	$(call log_info, Args: $(ARGS))
	@bash $(INSTALL_SCRIPT) test $(ARGS)
	$(shell exit 0)

## Run dynamic commands with arguments calling the install script.
%:
	@:
	$(call log_info, Running command: $(CMD_STR))
	$(call log_info, Args: $(ARGS))
	@bash $(INSTALL_SCRIPT) $(CMD_STR) $(ARGS)
	$(shell exit 0)

# Display help message.
help:
	$(call log, $(APP_NAME) Makefile)
	$(call log_break)
	$(call log, Usage:)
	$(call log,   make [target] [ARGS='--custom-arg value'])
	$(call log_break)
	$(call log, Available targets:)
	$(call log,   make build      - Build the binary using install script)
	$(call log,   make install    - Install the binary and configure environment)
	$(call log,   make clean      - Clean up build artifacts)
	$(call log,   make test       - Run tests)
	$(call log,   make help       - Display this help message)
	$(call log_break)
	$(call log, Usage with arguments:)
	$(call log,   make install ARGS='--custom-arg value' - Pass custom arguments to the install script)
	$(call log_break)
	$(call log, Example:)
	$(call log,   make install ARGS='--prefix /usr/local')
	$(call log_break)
	$(call log, $(APP_NAME) is a tool for managing Kubernetes resources)
	$(call log_break)
	$(call log, For more information, visit:)
	$(call log, 'https://github.com/rafa-mori/'$(APP_NAME))
	$(call log_break)
	$(call log_success, End of help message)
	$(shell exit 0)


