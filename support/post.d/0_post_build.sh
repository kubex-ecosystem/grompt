#!/usr/bin/env bash
# shellcheck disable=SC2065,SC2015

set -o nounset  # Treat unset variables as an error
set -o errexit  # Exit immediately if a command exits with a non-zero status
set -o pipefail # Prevent errors in a pipeline from being masked
set -o errtrace # If a command fails, the shell will exit immediately
set -o functrace # If a function fails, the shell will exit immediately
shopt -s inherit_errexit # Inherit the errexit option in functions

## Example of post build script
#_ROOT_DIR="$(git rev-parse --show-toplevel)"

#echo "Root directory: ${_ROOT_DIR}" > /dev/tty