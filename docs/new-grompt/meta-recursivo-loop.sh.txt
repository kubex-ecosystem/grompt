#!/usr/bin/env bash
# shellcheck disable=SC2065,SC2046,SC2317

# LookAtni + Grompt Meta-Recursive Virtuous Refactor Cycle
# "The most virtuous development cycle that's running!"

# Set project root path (Grompt)
_project_root_path="${PROJECT_ROOT_PATH:-${_project_root_path:-$(printf '%s\n' "/srv/apps/LIFE/KUBEX/grompt")}}"

# Keep without brakets, it is intentional, to prevent empty path and another issues
cd "$_project_root_path" || exit 1

# Set Internal Field Separator
IFS=$'\n\t'

# If true, adds a small delay after success messages for better readability
_lazy_exec="${_lazy_exec:-${LAZY_EXEC:-false}}"

# Quiet mode (no info or warning messages)
_quiet=${_quiet:-${QUIET:-false}}
_hide_about=${_hide_about:-${_HIDE_ABOUT:-false}}

# Debug mode
_debug=${_debug:-${DEBUG:-false}}

# Color codes for logs
_SUCCESS="\033[0;32m"
_WARN="\033[0;33m"
_ERROR="\033[0;31m"
_INFO="\033[0;36m"
_NOTICE="\033[0;35m"
_FATAL="\033[0;41m"
_TRACE="\033[0;34m"
_NC="\033[0m"

_get_stdout_alt() {
  local _dev_null=""
  _dev_null="/dev/null"
  test -f "${_dev_null:-/tmp/null}" || {
      touch "${_dev_null:-}" || {
        _dev_null="&2"
      }
      test -f "${_dev_null:-}" || {
          _dev_null="&2"
      }
  }
  _dev_null=">${_dev_null:-}"

  printf '%s\n' "${_dev_null:-}"
}

_provision() {
  printf '%s\n' "Provisioning environment..."
  _temp_combined_prompt="${_temp_combined_prompt:-$(mktemp -t combined_prompt.XXXXXX || echo "")}"
  if [[ -n "${_temp_combined_prompt:-}" && -f "${_temp_combined_prompt:-}" ]]; then
    rm -f "${_temp_combined_prompt:-}" || true
  fi
  return 0
}

log() {
  local _type=${1:-info}
  local _message=${2:-}
  local _debug_arg=${3:-}
  _debug_arg="${_debug_arg:-${_debug:-${DEBUG:-${_DEBUG:-false}}}}"

  case $_type in
    question|_QUESTION|-q|-Q)
      if [[ "${_quiet:-false}" == "true" && "${_debug_arg:-false}" == "true" ]]; then
        printf '%b[QUESTION]%b ❓  %s: ' "${_NOTICE:-\033[0;35m}" "${_NC:-\033[0m}" "$_message"
      fi
      ;;
    notice|_NOTICE|-n|-N)
      if [[ "${_quiet:-false}" == "true" && "${_debug_arg:-false}" == "true" ]]; then
        printf '%b[NOTICE]%b %s\n' "${_NOTICE:-\033[0;35m}" "${_NC:-\033[0m}" "$_message"
      fi
      ;;
    info|_INFO|-i|-I)
      if [[ "${_quiet:-false}" == "true" && "${_debug_arg:-false}" == "true" ]]; then
        printf '%b[INFO]%b  %s\n' "${_INFO:-\033[0;36m}" "${_NC:-\033[0m}" "$_message"
      fi
      ;;
    warn|_WARN|-w|-W)
      if [[ "${_quiet:-false}" == "true" && "${_debug_arg:-false}" == "true" ]]; then
        printf '%b[WARN]%b %s\n' "${_WARN:-\033[0;33m}" "${_NC:-\033[0m}" "$_message"
      fi
      ;;
    error|_ERROR|-e|-E)
      printf '%b[ERROR]%b %s\n' "${_ERROR:-\033[0;31m}" "${_NC:-\033[0m}" "$_message" >&2
      ;;
    success|_SUCCESS|-s|-S)
       printf '%b[SUCCESS]%b %s\n' "${_SUCCESS:-\033[0;32m}" "${_NC:-\033[0m}" "$_message"
       if [[ "${_lazy_exec:-false}" == "true" ]]; then
         sleep 2
       fi
      ;;
    fatal|_FATAL|-f|-F)
      printf '%b[FATAL]%b %s\n' "${_FATAL:-\033[0;41m}" "${_NC:-\033[0m}" "Exiting due to fatal error: $_message" >&2
      # clear_build_artifacts || true
      test $(declare -f _clear_globals >$(_get_stdout_alt)) && _clear_globals
      exit 1 || kill -9 $$
      ;;
    separator|_SEPARATOR|hr|-hr|-HR|line)
      # if [[ "${_debug_arg:-false}" != "true" ]]; then
        local _columns=${COLUMNS:-$(tput cols || echo 80)}
        local _margin=$(( _columns - ( _columns / 2 ) ))
        _message="${_message// /¬}"
        _message="$(printf '%b%s%b %*s' "${_TRACE:-\033[0;34m}" "${_message:-}" "${_NC:-\033[0m}" "$((_columns - ( "${#_message}" + _margin )))" '')"
        _message="${_message// /\#}"
        _message="${_message//¬/ }"
        printf '%s\n' "${_message:-}" >&2
      # fi
      ;;
    *)
      log "info" "$_message" "${_debug_arg:-false}" || true
      ;;
  esac

  return 0
}

clear_screen() {
  if [[ "${_quiet:-false}" != "true" && "${_debug:-false}" != "true" ]]; then
    printf "\033[H\033[2J"
  fi
}

clear_build_artifacts() {
  test $(declare -f _clear_globals >$(_get_stdout_alt)) && _clear_globals

  local build_dir="${_ROOT_DIR:-$(realpath '../')}/dist"
  if [[ -d "${build_dir}" ]]; then
    rm -rf "${build_dir}" || true
    if [[ -d "${build_dir}" ]]; then
      log error "Failed to remove build artifacts in ${build_dir}."
    else
      log success "Build artifacts removed from ${build_dir}."
    fi
  else
    log notice "No build artifacts found in ${build_dir}."
  fi
}

get_current_shell() {
  local shell_proc
  shell_proc=$(cat /proc/$$/comm)
  case "${0##*/}" in
    ${shell_proc}*)
      local shebang
      shebang=$(head -1 "$0")
      printf '%s\n' "${shebang##*/}"
      ;;
    *)
      printf '%s\n' "$shell_proc"
      ;;
  esac
}

_ensure_globals(){
  # I added _ to all VARS to avoid colliding with any system ones,
  # which by convention don't use this prefix, plus only use uppercase,
  # so all were also converted to lowercase.
  # Lets start everything now with all basic functions declared...

  _example_artifact="${_example_artifact:-}" || true
  _artifact_content="${_artifact_content:-}" || true

  _refactored_content="${_refactored_content:-}" || true
  _refactored_output="${_refactored_output:-}" || true

  _grompt_generated_prompt="${_grompt_generated_prompt:-}" || true
  _grompt_ask="${_grompt_ask:-}" || true

  _combined_prompt="${_combined_prompt:-}" || true
  _prompt_content="${_prompt_content:-}" || true
  _example_prompt="${_example_prompt:-}" || true

  _duration="${_duration:-}" || true
  _start_time="${_start_time:-}" || true
  _end_time="${_end_time:-}" || true

  _start_time=$(date +%s)                                             # SCRIPT START TIME VAR INITIALIZATION
  # readonly _start_time                                                # SCRIPT START TIME LOCK
  _end_time=$(date +%s)                                               # SCRIPT END TIME VAR INITIALIZATION

  _temp_combined_prompt="${_temp_combined_prompt:-}"                                            # SCRIPT TEMPORARY PROMPT VAR INITIALIZATION

  # LookAtni Settings
  _artifact_content="${_artifact_content:-}"                                                # LOOKATNI ARTIFACT CONTENT VAR INITIALIZATION

  # LLM Settings
  # Markers for replace
  _llm_provider="gemini"                                              # LLM PROVIDER
  _llm_model="gemini-2.0-flash"                                       # LLM MODEL
  _llm_api_key='AIzaSyBTcVVqCp7VxLlBApsnpxxeFsK2TigIBaU'              # GEMINI API KEY [OPTIONAL, CAN BE SET AS ENV VAR]
  _max_tokens=32000                                                   # MAX TOKENS (GEMINI-2.0-FLASH SUPPORTS UP TO 32000)

  # Workspace Path
  _workspace_path="/srv/apps/LIFE/KUBEX"                              # WORKSPACE PATH [OPTIONAL]

  # LookAtni Paths
  _lookatni_path="${_workspace_path:-$(pwd)}/lookatni-file-markers"   # LOOKATNI ROOT PATH
  _lookatni_bin="/home/user/.local/bin/lookatniCli"                   # LOOKATNI BINARY

  # Grompt Paths
  _grompt_path="${_workspace_path:-$(pwd)}/grompt"                    # GROMPT ROOT PATH
  _grompt_bin="${_grompt_path}/dist/grompt_linux_amd64"               # GROMPT BINARY

  # Examples Path
  _example_parent="${_grompt_path}/docs/examples/virt-cycles"         # TARGET TEST PARENT FOLDER

  # Example Project Paths
  _example_project="${_example_parent:-$(pwd)}/test-project"            # 1: TEST TARGET
  _example_artifact="${_example_parent}/test-project-artifact.md"       # 2: LOOKATNI ARTIFACT
  _example_prompt="${_example_parent}/improvement-prompt.md"            # 3: GROMPT GENERATED SCREENING PROMPT
  _example_result="${_example_parent}/test-project-refactored"          # 4: PROJECT REFACTORED, ALREADY RE-EXPANDED WITH LOOKATNI
  _example_res_artifact="${_example_parent}/test-project-refactored.md" # 5: PROJECT ARTIFACT REFACTORED
}

_clear_globals() {
  unset _example_artifact || true
  unset _artifact_content || true

  unset _refactored_content || true
  unset _refactored_output || true

  unset _grompt_generated_prompt || true
  unset _grompt_ask || true

  unset _combined_prompt || true
  unset _prompt_content || true
  unset _example_prompt || true

  unset _example_res_artifact || true

  unset _duration || true
  unset _end_time || true

  # unset _start_time || true
}

# Creates a temporary directory for cache
_provision || log fatal "Provisioning failed!"

# Ensure global variables are initialized
_ensure_globals || log fatal "Failed to ensure global variables!"

# Set a trap to clear script cache on exit
set_trap() {
  local current_shell=""
  current_shell=$(get_current_shell)

  # Check the shell type and set options accordingly
  case "${current_shell}" in
    *ksh|*zsh|*bash)

      declare -a FULL_SCRIPT_ARGS=("$@")
      if [[ "${FULL_SCRIPT_ARGS[*]}" == *--debug* ]]; then
          set -x
      fi

      # Enable strict mode
      if [[ "${current_shell}" == "bash" ]]; then
        set -o nounset  # Treat unset variables as an error
        set -o errexit  # Exit immediately if a command exits with a non-zero status
        set -o pipefail # Prevent errors in a pipeline from being masked
        set -o errtrace # If a command fails, the shell will exit immediately
        set -o functrace # If a function fails, the shell will exit immediately
        shopt -s inherit_errexit # Inherit the errexit option in functions
      fi

      # Set a trap to clear the script cache on exit
      trap '_clear_globals' EXIT HUP INT QUIT ABRT ALRM TERM
      ;;
  esac
}

# Clear global variables and temporary files
cleanup() {
  # Clean the script trap to avoid "gremlins"
  trap - EXIT HUP INT QUIT ABRT ALRM TERM

  # Don't make the PC or anyone waste time... lol
  if [[ -n "${_temp_combined_prompt:-}" && -f "${_temp_combined_prompt:-}" ]]; then
    rm -f "${_temp_combined_prompt:-}" || true
  fi

  # Clear globals if the function exists
  if test $(declare -f _clear_globals >$(_get_stdout_alt)); then
    _clear_globals || true
    unset -f _clear_globals || true
  fi

  return 0
}

set_trap "$@"

first_step() {
  log info "STEP 1: Generating project artifact..." true

  "${_lookatni_bin}" generate "${_example_parent:-}" "${_example_artifact:-}"

  # Check if artifact was generated
  test -f "${_example_artifact:-}" || {
    log fatal "Artifact generation failed!" true
  }

  log success "Artifact generated: ${_example_artifact:-}" && sleep 2

  return 0
}

second_step() {
  log info "STEP 2: Generating professional prompt..."

  _grompt_generated_prompt=$("${_grompt_bin:-}" generate \
      --provider "${_llm_provider:-gemini}" \
      --apikey "${_llm_api_key:-${GEMINI_API_KEY:-}}" \
      --model "${_llm_model:-gemini-2.0-flash}" \
      --ideas 'Analyze this Go project and identify code improvements following Go best practices' \
      --ideas 'Focus on: error handling, naming conventions, idiomatic code, performance' \
      --ideas 'Maintain LookAtni file structure (markers //<ASCII[28]>/ filename /<ASCII[28]>//)' \
      --ideas 'The placeholder <ASCII[28]> represents ASCII character 28 (File Separator - ) and must be PRINTED in the presented result.' \
      --ideas 'Return the complete refactored code with explanations of changes, without title or footer, but explanations in comments within the code itself.' \
      --ideas 'The final file and the files within the code must maintain the LookAtni markers exactly as in the original artifact.' \
      --max-tokens "${_max_tokens:-10240}" \
      --purpose 'code' || echo '')

  test -n "${_grompt_generated_prompt:-}" || {
      log fatal "Prompt generation failed!" true
  }

  printf '%b\n' "${_grompt_generated_prompt:-}" > "${_example_prompt:-}"
  test -f "${_example_prompt:-}" || {
    log fatal "Failed to save generated prompt!" true
  }

  return 0
}

third_step() {
  log info "STEP 3: Executing AI refactoring..." true

  # Create combined prompt
  _artifact_content="$(cat "${_example_artifact:-}" --show-nonprinting)"
  test -n "${_artifact_content}" || {
    log fatal "Failed to read artifact!"
  }

  if [[ " ${#_artifact_content} " -lt $(( _max_tokens / 2 - 1000 )) ]]; then
      _combined_prompt="$(printf '%b\nTARGET CONTENT:\n%s\n%s\n%s\n' "$(cat "${_example_prompt:-}" --show-nonprinting)" '```plaintext' "${_artifact_content:-}" '```')"
  else
      _combined_prompt="$(printf '%s\n' "$(cat "${_example_prompt:-}" --show-nonprinting)")"
  fi

  if [[ ! -f "${_temp_combined_prompt:-}" && -n "${_temp_combined_prompt:-}" ]]; then
    touch "${_temp_combined_prompt:-}"
  elif [[ ! -f "${_temp_combined_prompt:-}" && -z "${_temp_combined_prompt:-}" ]]; then
    _temp_combined_prompt=$(mktemp -t "combined_prompt.XXXXXX" || echo "")
    if [[ -z "${_temp_combined_prompt:-}" || ! -f "${_temp_combined_prompt:-}" ]]; then
      log fatal "Failed to create temporary prompt file!"
    fi
    export _temp_combined_prompt
  fi

  # Now save to the real temporary file so it doesn't bother if something goes wrong.. hehe
  printf '%b\n' "${_combined_prompt:-}" > "${_temp_combined_prompt:-}"

  test -f "${_temp_combined_prompt:-}" || {
      log fatal "Failed to generate combined prompt!"
  }

  log success "Combined prompt generated: ${_temp_combined_prompt:-}" && sleep 2

  _grompt_task_extract="$(cat "${_temp_combined_prompt:-}" --show-nonprinting)"

  # Execute with Gemini (using temporary file to work around input limit)
  _grompt_ask="$("${_grompt_bin:-}" ask \
      --prompt "${_grompt_task_extract:-}" \
      --provider "${_llm_provider:-gemini}" \
      --apikey "${_llm_api_key:-${GEMINI_API_KEY:-}}" \
      --max-tokens "${_max_tokens:-8000}" && true || echo '')" # Just for a cool ensure with this true...

  # Check if it's filled
  test -n "${_grompt_ask:-}" || {
      log fatal "Refactoring failed!"
  }

  # Fill the file
  printf '%b\n' "${_grompt_ask:-}" > "${_example_res_artifact:-}"

  # Check if the file was filled
  test -f "${_example_res_artifact:-}"|| {
      log fatal "Refactoring failed!"
  }

  # Everything ok? Print it!
  log success "Refactoring completed: ${_example_res_artifact:-}" && sleep 2

  return 0
}

fourth_step() {
  # I'm adding printf to ensure the invisible character will be printed without
  # any kind of expansion, etc...
  # The cat is good to keep with --show-nonprinting to ensure the invisible character
  _exemplo="$(cat "${_example_res_artifact:-}" --show-nonprinting)"
  _exemplo="${_exemplo#*\`\`\`go}"
  _exemplo="${_exemplo//\/\/<ASCII\[28\]>\//$(printf "//\x1C/")}"
  _exemplo="${_exemplo//\/<ASCII\[28\]>\/\//$(printf "/\x1C//")}"

  printf '%b\n' "${_exemplo:-}" > "${_example_res_artifact:-}" || {
    log fatal "Failed to process markers!"
  }

  test -f "${_example_res_artifact:-}" || {
    log fatal "Refactored artifact not found!"
  }

  # For this to pass we need the following:
  # 1: Remove the last line of the generated $_example_res_artifact file.
  # 2: Replace the ␜ character with the real one
  local _validate_lookatni=""
  _validate_lookatni=$("${_lookatni_bin:-}" validate "${_example_res_artifact:-}" || echo '')

  # shellcheck disable=SC2143
  if [[ -z "${_validate_lookatni:-}" || $(echo "${_validate_lookatni}" | grep 'Total markers: 0' -q) ]] ; then
    log fatal "Validation failed!"
  fi

  log info "STEP 4: Extracting refactored project..."

  "${_lookatni_bin:-}" extract "${_example_res_artifact:-}" "${_example_result:-}" --overwrite --create-dirs

  test -d "${_example_result:-}" || {
    log fatal "Extraction failed!"
  }

  log success "Refactored project extracted: ${_example_result:-}"

  return 0
}

print_summary() {
  local _key=""

  if [[ -n "${_key}" && ${#_key} -gt 8 ]]; then
    _key=${_llm_api_key:-}
    _key="${_key:0:4}****${_key: -4}"
  fi

  log hr
  log success "META-RECURSIVITY COMPLETE!"
  log success "Workspace: ${_workspace_path:-}"
  log success "Original project: ${_example_project:-}"
  log success "Artifact: ${_example_artifact:-}"
  log success "Result: ${_example_result:-}"
  log success "Refactored: ${_example_res_artifact:-}"
  log hr
  log success "Provider: ${_llm_provider:-gemini}"
  log success "Model: ${_llm_model:-gemini-2.0-flash}"
  log success "Max Tokens: ${_max_tokens:-10240}"
  log success "API Key: ${_key:-[None]}"
  log success "Prompt: ${_example_prompt:-}"
  log hr
  log success "End of Virtuous Cycle! 🚀"
  _end_time=$(date +%s)
  _duration=$(( _end_time - _start_time ))
  if [ $_duration -gt 60 ]; then
    log notice "Time elapsed: $(( _duration / 60 )) minutes and $(( _duration % 60 )) seconds" true
  else
    log notice "Time elapsed: ${_duration} seconds" true
  fi
  log hr

  return 0
}

main() {
  log hr
  log info "STARTING LOOKATNI + GROMPT META-RECURSIVITY!" true
  log hr

  first_step || log fatal "Step 1 failed!"
  second_step || log fatal "Step 2 failed!"
  third_step || log fatal "Step 3 failed!"
  fourth_step || log fatal "Step 4 failed!"
  print_summary || log fatal "Summary printing failed!"
}

main || log fatal "Script execution failed!"
