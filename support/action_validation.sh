#!/usr/bin/env bash

set -euo pipefail
set -o errtrace
set -o functrace
set -o nounset
set -o posix
IFS=$'\n\t'

_SOURCED=false

ensure_marker() {
  local MARKER_PATH="${1:-}"
  if [[ -z "$MARKER_PATH" ]]; then
    printf '%s\n' "❌ Marker path is not set. Cannot proceed."
    return 1
  fi
  if [[ ! -f "$MARKER_PATH" ]]; then
    touch "$MARKER_PATH" || {
      printf '%s\n' "❌ Failed to create marker file: $MARKER_PATH"
      return 1
    }
    chmod 644 "$MARKER_PATH"
    printf '%d\n' "COUNT=0" > "$MARKER_PATH"
  fi
  return 0
}

get_git_tag(){
  get_tag_local_git(){
    git describe --tags --abbrev=0 2>/dev/null || echo "v0.0.0"
  }
  get_tag_env_git(){
    echo "${GITHUB_REF:-"$(get_tag_local_git)"}"
  }
  get_tag_remote_git(){
    git ls-remote --tags origin | grep -o 'refs/tags/v[0-9]\+\.[0-9]\+\.[0-9]\+' | sort -V | tail -n1 || echo "v0.0.0"
  }
  local TAG="${1:-}"
  if [[ -z "$TAG" ]]; then
    TAG="$(get_tag_env_git)"
  fi
  if [[ "$TAG" == "v0.0.0" ]]; then
    TAG="$(get_tag_local_git)"
  fi
  if [[ "$TAG" == "v0.0.0" ]]; then
    TAG="$(get_tag_remote_git)"
  fi
  if [[ "$TAG" == "v0.0.0" ]]; then
    printf '%s\n' "❌ No valid tag found. Please ensure you have tags in your repository."
    exit 1
  fi
  printf '%s\n' "$TAG"
}

sanitize_version(){
  local VERSION="${1:-"${REF:-"${GITHUB_REF:-"refs/tags/v0.0.0"}"}"}"

  VERSION="${VERSION##*tags/}"
  VERSION="${VERSION##*v}"
  VERSION="${VERSION%%-*}" # Remove any pre-release suffix
  VERSION="${VERSION%%+*}" # Remove any build metadata
  VERSION="${VERSION//[^0-9.]/}" # Remove any non-numeric characters

  printf '%s\n' "$VERSION"
}

get_version_hash(){
  # Generate a hash for the version
  printf '%s\n' "$(sanitize_version "$(get_git_tag "${1:-}")")" | sha256sum | awk '{print $1}'
}

parse_vars(){
  # Marker directory by default
  local MARKER_DIR="docs/vsctl"

  # Check if the marker directory is set and exists
  if [[ ! -d "${MARKER_DIR:-}" ]]; then
    mkdir -p "${MARKER_DIR:-}" && chmod 755 "${MARKER_DIR:-}" || {
      printf '%s\n' "❌ Failed to create marker directory: ${MARKER_DIR:-}"
      return 1
    }
  fi

  # Repo vars
  local REPO=""
  local ACTOR=""
  local EVENT=""
  local IS_FORK=""
  local REF=""
  local CONTRIBUTORS=""
  local VERSION=""

  # Marker vars
  local MARKER_NAME=""
  local MARKER_PATH=""
  local CUR_TIMESTAMP=""

  # Control variables
  local WILL_PROCEED=false
  local COUNT=0

  # Extract version from the first argument or fallback to GITHUB_REF
  VERSION="$(sanitize_version "${1:-"${VERSION:-"${GITHUB_REF:-"refs/tags/v0.0.0"}"}"}")"

  # Create a associative array to hold the arguments
  declare -A _ARGS_LIST=(
    ["version"]="${VERSION}"
    ["version_hash"]="$(get_version_hash "${_ARGS_LIST["version"]}")"
    ["repo"]="${2:-"${REPO:-${GITHUB_REPOSITORY:-$(git config --get remote.origin.url | sed 's|.*://||; s|\.git$||')}}"}"
    ["actor"]="${3:-"${ACTOR:-$(git config --get user.name)}"}"
    ["event"]="${4:-"${EVENT:-${GITHUB_EVENT_NAME:-"unknown"}}"}"
    ["is_fork"]="${5:-"${IS_FORK:-"$(git config --get remote.origin.url | grep -q 'fork' && echo "true" || echo "false")"}"}"
    ["ref"]="${6:-"${REF:-${GITHUB_REF:-"refs/tags/v0.0.0"}}"}"
  )
  
  # Get MARKER's directory and absolute path for current repo based on the version hash and version
  MARKER_NAME=".kubex_publish_marker_${_ARGS_LIST["version"]}_${_ARGS_LIST["version_hash"]}"
  MARKER_PATH="${MARKER_DIR}/${MARKER_NAME}"

  # Internal function to validate context with shared scope
  validate_context(){
    # Check if actor is in contributors
    # shellcheck disable=SC2076
    if [[ ! " $CONTRIBUTORS " =~ " $ACTOR " ]]; then
      printf '%s\n' "❌ Actor '$ACTOR' is not an authorized contributor."
      return 1
    fi
    # Check event type and repository fork status
    if [[ -z "$EVENT" ]]; then
      printf '%s\n' "❌ EVENT is not set. Cannot proceed."
      return 1
    fi
    # Check if the event is a push or workflow_dispatch
    if [[ "$EVENT" != "push" && "$EVENT" != "workflow_dispatch" ]]; then
      printf '%s\n' "❌ Only 'push' or manual dispatch allowed. Got '$EVENT'."
      return 1
    fi
    # Check if the repository is a fork
    if [[ "$IS_FORK" == "true" ]]; then
      printf '%s\n' "❌ Workflow cannot run from a fork."
      return 1
    fi
    # Check if the ref who triggered the workflow
    if [[ -z "$REF" ]]; then
      printf '%s\n' "❌ REF is not set. Cannot proceed."
      return 1
    fi
    # Check if the ref is a tag and matches the semver format
    if [[ ! "$REF" =~ ^refs/tags/v[0-9]+\.[0-9]+\.[0-9]+$ ]]; then
      printf '%s\n' "❌ Tag does not match semver format: $REF"
      return 1
    fi
    # Check if the version is set
    if [[ -z "$VERSION" ]]; then
      printf '%s\n' "❌ Version is not set. Cannot proceed."
      return 1
    fi
    
    return 0
  }

  validate_marker() {
    # Extract repository contributors username list
    CONTRIBUTORS="$(gh api repos/"${_ARGS_LIST["repo"]}"/contributors --jq '.[].login' 2>/dev/null || true)"
    # Remove trailing space
    CONTRIBUTORS="$(echo "${CONTRIBUTORS:-}" | sed 's/[[:space:]]*$//')"

    # Initialize control variables
    WILL_PROCEED=false
    COUNT=0

    if [ -n "${MARKER_PATH}" ]; then
      # Check for existing marker file
      if test -f "${MARKER_PATH}"; then
        # Read the count from the marker file
        COUNT=$(cat "${MARKER_PATH}" | grep '^COUNT=' | awk -F'=' '{print $2}')

        # Increment the count
        COUNT=$((COUNT + 1))
        printf '%s\n' "COUNT=$COUNT" > "${MARKER_PATH}"

        # Output the current count
        printf '%s\n' "📊 Current publish marker count for ${VERSION}: $COUNT"
        printf '%s\n' "🗓️ Last execution: $(date -d "@${CUR_TIMESTAMP}" +"%Y-%m-%d %H:%M:%S")"

        # If count exceeds 2, abort the publication
        # and inform the user
        if [ $COUNT -ge 2 ]; then
          WILL_PROCEED=false
          printf '%s\n' "❌ Publish marker count for ${VERSION} has reached $COUNT executions. Aborting deployment."
          printf '%s\n' "Please check the marker file: ${MARKER_PATH}"
          printf '%s\n' "If you want to reset the count, please remove the marker file manually."
        else
          WILL_PROCEED=true
        fi
      else
        WILL_PROCEED=true
      fi
    else 
      WILL_PROCEED=true
    fi

    if test "$WILL_PROCEED" = true; then
      # Log the deployment proceeding
      printf '%s\n' "✅ Proceeding with deployment for version ${VERSION}."

      # Ensure the marker file exists
      # and create/set the initial count if it doesn't exist
      ensure_marker "${MARKER_PATH}" || {
        printf '%s\n' "❌ Failed to ensure marker file: ${MARKER_PATH}"
        return 1
      }

      # After ensuring the marker file, this condition NEEDS to be checked again and pass
      if test -f "${MARKER_PATH}"; then
        # Read the count from the marker file
        COUNT=$(grep '^COUNT=' "${MARKER_PATH}" | awk -F'=' '{print $2}')

        # Read the timestamp from the marker file and check if it is older than 24 hours
        TIMESTAMP=$(grep '^TIMESTAMP=' "${MARKER_PATH}" | awk -F'=' '{print $2}')
        CURRENT_TIMESTAMP=$(date +%s)
        MARKER_TIMESTAMP=$(date -d "$TIMESTAMP" +%s)
        TIME_DIFF=$((CURRENT_TIMESTAMP - MARKER_TIMESTAMP))

        if [ $TIME_DIFF -gt 86400 ]; then
          # If the marker is older than 24 hours, reset the count
          printf '%s\n' "🕒 Marker file is older than 24 hours. Removing marker file, allowing runners again over this tag."
          rm -f "${MARKER_PATH}"
          COUNT=0
        else
          # Increment the count
          COUNT=$((COUNT + 1))
        fi
      else 
        printf '%s\n' "❌ Error ensuring marker file: ${MARKER_PATH}"
        return 1
      fi

      # Update the marker file with the new count and timestamp
      CUR_TIMESTAMP=$(date +%s)
      printf '%s\n' "COUNT=${COUNT}" > "${MARKER_PATH}"
      printf '%s\n' "TIMESTAMP=$(date -d "@${CUR_TIMESTAMP}" +"%Y-%m-%d %H:%M:%S")" >> "${MARKER_PATH}"

      # Add the marker file to git
      if [ $COUNT -gt 0 ]; then
          git add "${MARKER_PATH}" || true
      else
          git rm "${MARKER_PATH}" || true
      fi

      git commit -m "add publish marker for version ${VERSION}, count: ${COUNT}" || true
      git push origin HEAD:main
    else
      printf '%s\n' "❌ Deployment will not proceed due to marker validation failure."
      return 1
    fi

    return 0
  }

  # Validate the context variables in the shared scope
  validate_context || {
    printf '%s\n' "❌ Context validation failed. Cannot proceed."
    return 1
  }

  validate_marker || {
    printf '%s\n' "❌ Marker validation failed. Cannot proceed."
    return 1
  }
}

# 1: steps.extract_info.outputs.version
# 2: steps.extract_info.outputs.will_proceed
# 3: steps.extract_info.outputs.marker
# 4: steps.extract_info.outputs.count
main () {
  # Store all arguments in an array
  # This allows us to pass them to the function
  # without losing the original arguments, and this sanitizes them
  # to avoid issues with special characters or spaces
  local args=()

  # Check if the script is run with arguments in enlaced mode
  # If no arguments are provided, print usage and return an error
  if [[ $# -eq 0 ]]; then
    printf '%s\n' "❌ No arguments provided."
    printf '%s\n' "Usage: $0 <version> <will_proceed> <marker> <count>"
    return 1
  else 
    args=("$@")
  fi

  # Shift all arguments to the left, leaving an empty array
  shift $(( $# )) 

  declare -a _cmd=(
    validate_marker 
    "${args[@]}"
  )

  # If the script is run directly, execute the command
  # Otherwise, export the function for later use
  if [[ "${BASH_SOURCE[0]}" == "${0}" ]]; then
    # Sanitize input arguments
    if [[ $# -eq 0 ]]; then
      printf '%s\n' "Usage: $0 <version> <will_proceed> <marker> <count>"
      return 1
    fi

    "${_cmd[@]}" || {
      printf '%s\n' "❌ Error: Validation failed."
      return 1
    }
  else
    # If the script is sourced, we don't want to exit
    _SOURCED=true

    # If sourced, define the function for later use
    export -f validate_marker
  fi
}

main "${@:-}"

# If the script is sourced, we don't want to exit
# shellcheck disable=SC2319
test $_SOURCED = false && exit $?

# End of script