#!/usr/bin/env bash

# Indent output by (indent) levels
#
# Example:
#   echo "example" | indent 2
#   cat file.txt | indent
indent() {
  local indentSize=2
  local indent=1
  if [ -n "${1:-}" ]; then
    indent="$1"
  fi
  pr --omit-header --indent=$((indent * indentSize))
}

# Check if dependencies are available.
#
# If any dependencies are missing, an error message will be printed to
# stderr and the program will exit, running traps on EXIT beforehand.
#
# Example:
#   check_dependencies git bash node
check_dependencies() {
  local missing=false
  for command in "$@"; do
    if ! command -v "$command" &> /dev/null; then
      echo "$0: script requires '$command', but it is not in your PATH" >&2
      missing=true
    fi
  done

  if [ $missing = true ]; then
    exit 1
  fi
}

# Emit a message to stderr and exit.
#
# This prints the arguments to stderr before exiting.
error() {
  echo "$@" >&2
  exit 1
}

# Run a command, with tracing.
#
# This prints a command to stderr for tracing, in a format similar to
# the bash xtrace option (i.e. set -x, set -o xtrace), then runs it using
# eval if the first argument (dry run) is false.
#
# If dry run is true, the command and arguments will be printed, but
# not executed.
#
# Example:
#   xtrace $DRY_RUN rm -rf /
run_trace() {
  local args=("$@")
  local dry_run=${args[0]}
  args=("${args[@]:1}")

  echo "+ ${args[*]}" >&2
  if [ "$dry_run" = false ]; then
    eval "${args[@]}"
  fi
}
