#!/usr/bin/env bash
#
# Script to build binary from Golang source and move binary to the correct location

# Constants
SCRIPT_NAME=$(basename "$0")

# Global variables.
declare OS_NAME
declare ARCH_NAME

# Log a message
log() {
  local MESSAGE="$1"
  echo "${MESSAGE}"
}

# Display usage information
usage() {
  echo "Usage: ${SCRIPT_NAME} --binary <binary-name> --path <go-src-path>"
}

# Remove binary file if it exists
remove_old_binary() {
  local BINARY_FULL_PATH="$1"
  local rm_cmd

  if [[ -f ${BINARY_FULL_PATH} ]]; then
    log "Removing old binary..."
    rm_cmd=$(rm "${BINARY_FULL_PATH}")

    if [[ "$rm_cmd" -ne "0" ]]; then
      log "Failed to remove old binary at ${BINARY_FULL_PATH}"
      exit 1
    fi
  fi
}

# Build binary
# shellcheck disable=SC2034
build_binary() {
  local BINARY_FULL_PATH="$1"
  local GO_SRC="$2"
  local build_linux
  local build_m1_mac
  local build_non_m1_mac

  log "Building binary in path ${BINARY_FULL_PATH}..."
  build_linux=$(GOOS=linux GOARCH=amd64 go build -o "${BINARY_FULL_PATH}" "${GO_SRC}")
  build_m1_mac=$(GOOS=darwin GOARCH=arm64 go build -o "${BINARY_FULL_PATH}" "${GO_SRC}")
  build_non_m1_mac=$(GOOS=darwin GOARCH=amd64 go build -o "${BINARY_FULL_PATH}" "${GO_SRC}")

  #Check OS and ARCH
  if [[ "$OS_NAME" == "linux" ]] && [[ "$ARCH_NAME" == "amd64" ]]; then
      build_linux=$(GOOS=linux GOARCH=amd64 go build -o "${BINARY_FULL_PATH}" "${GO_SRC}")
  elif [[ "$OS_NAME" == "darwin" ]] && [[ "$ARCH_NAME" == "arm64" ]]; then
      build_m1_mac=$(GOOS=darwin GOARCH=arm64 go build -o "${BINARY_FULL_PATH}" "${GO_SRC}")
 elif [[ "$OS_NAME" == "mac" ]] && [[ "$ARCH_NAME" == "arm64" ]]; then
      build_m1_mac=$(GOOS=darwin GOARCH=arm64 go build -o "${BINARY_FULL_PATH}" "${GO_SRC}")
 elif [[ "$OS_NAME" == "mac" ]] && [[ "$ARCH_NAME" == "amd64" ]]; then
      build_non_m1_mac=$(GOOS=darwin GOARCH=amd64 go build -o "${BINARY_FULL_PATH}" "${GO_SRC}")
  elif [[ "$OS_NAME" == "darwin" ]] && [[ "$ARCH_NAME" == "amd64" ]]; then
      build_non_m1_mac=$(GOOS=darwin GOARCH=amd64 go build -o "${BINARY_FULL_PATH}" "${GO_SRC}")
  elif [[ "$OS_NAME" == "UNKNOWN" ]] || [[ "$ARCH_NAME" == "UNKNOWN" ]]; then
      log "Failed to determine OS or ARCH"
      exit 1
  fi

  log "Binary built successfully"
}

# Add binary to .gitignore at the root of the git repository
add_to_gitignore_if_not_exist() {
  local BINARY_NAME="$1"
  local CURRENT_DIR
  CURRENT_DIR="$(pwd)"
  local GIT_ROOT_DIR
  local GITIGNORE_FILE

  # Find the root .git repository
  GIT_ROOT_DIR=$(git rev-parse --show-toplevel 2>/dev/null)
  if [[ -z "${GIT_ROOT_DIR}" ]]; then
    log "Failed to locate the root directory of a git repository."
    return 1
  fi

  GITIGNORE_FILE="${GIT_ROOT_DIR}/.gitignore"

  if [[ ! -f "${GITIGNORE_FILE}" ]]; then
    log "Failed to find .gitignore file, creating one at the git repository root."
    touch "${GITIGNORE_FILE}"
  fi

  if ! grep -q "^${BINARY_NAME}$" "${GITIGNORE_FILE}"; then
    log "Adding binary to .gitignore file..."
    echo "${BINARY_NAME}" >> "${GITIGNORE_FILE}"
  else
    log "Binary already exists in .gitignore file."
  fi

  cd "${CURRENT_DIR}" || exit # Restore original working directory if needed
}

set_current_os() {
  local uname_output
  local os_name_found
  uname_output=$(uname -s)

  case "${uname_output}" in
    Linux*)     os_name_found=linux;;
    Darwin*)    os_name_found=mac;;
    *)          os_name_found="UNKNOWN:${uname_output}"
  esac

  log "Current OS: ${os_name_found}"
  export OS_NAME="${os_name_found}"
}

set_current_arch() {
  local uname_output
  local arch_name_found
  uname_output=$(uname -m)

  case "${uname_output}" in
    x86_64*)    arch_name_found=amd64;;
    arm64*)     arch_name_found=arm64;;
    *)          arch_name_found="UNKNOWN:${uname_output}"
  esac

  log "Current Arch: ${arch_name_found}"
  export ARCH_NAME="${arch_name_found}"
}

# Main Function
main() {
  local BINARY_NAME=""
  local GO_SRC=""

  set_current_os
  set_current_arch

  # Parse arguments
  while (( "$#" )); do
    case "$1" in
      --binary)
        BINARY_NAME="$2"
        shift 2
        ;;
      --path)
        GO_SRC="$2"
        shift 2
        ;;
      *)
        usage
        exit 1
        ;;
    esac
  done

  # Check if BINARY_NAME and GO_SRC are provided
  if [[ -z ${BINARY_NAME} ]] || [[ -z ${GO_SRC} ]]; then
    usage
    exit 1
  fi

  # Check if GO_SRC file exist.
  if [[ ! -f ${GO_SRC} ]]; then
    log "Failed to find Go source file: ${GO_SRC}"
    exit 1
  fi

  remove_old_binary "${BINARY_NAME}"
  build_binary "${BINARY_NAME}" "${GO_SRC}"
  add_to_gitignore_if_not_exist "${BINARY_NAME}"
}

main "$@"
