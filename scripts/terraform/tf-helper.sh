#!/usr/bin/env bash
#
##################################################################################
# Script Name: Terraform Modules Upgrader
#
# Author: Alex Torres (github.com/Excoriate), alex_torres@outlook.com
#
# Usage: ./script.sh --modules_dir=modules_dir_path [--max_depth_flag=depth]
#
# Description: This bash script searches for Terraform modules within a specified directory
#              up to a certain depth and runs 'terraform init -upgrade' in each module's directory.
#
# Parameters:
#    --modules_dir:      The path to the directory containing Terraform modules [mandatory].
#    --max_depth_flag:   The maximum depth for directory traversal. Default is 3 [optional].
#
# Examples:
#    Upgrade modules: ./script.sh --modules_dir=./modules --max_depth_flag=2
#
# Note: The script assumes that Terraform is installed and available on the system's PATH.
#
# For further details and support, contact the author.
#
##################################################################################

set -euo pipefail

# Constants
readonly DEFAULT_ACTION="upgrade" # it also supports "docs"

# Log a message
log() {
    local -r msg="${1}"
    echo "${msg}" >&2
}

# Obtain the Terraform modules' paths
get_terraform_module_paths() {
    local -r modules_dir="${1}"
    local -r max_depth_flag="${2:-3}" # Default max depth is 3 if not provided

    find "${modules_dir}" -mindepth 1 -maxdepth "${max_depth_flag}" -type d | while read -r dir; do
        if compgen -G "${dir}"/*.tf > /dev/null; then
            echo "${dir}"
        fi
    done
}

# Upgrade the Terraform modules
upgrade_terraform_modules() {
    local -r modules_dir="${1}"
    local -r max_depth_flag="${2:-3}" # Default max depth is 3 if not provided

    local modules_path
    modules_path=$(get_terraform_module_paths "${modules_dir}" "${max_depth_flag}")
    if [[ -z "${modules_path}" ]]; then
        log "No Terraform modules found in the directory: ${modules_dir}"
        return 0
    fi

    echo "${modules_path}" | while read -r module_path; do
        log "Upgrading Terraform module at: ${module_path}"
        (cd "${module_path}" && terraform init -upgrade)
    done

    log "Terraform modules upgraded successfully"
}

generate_terraform_docs() {
    local -r modules_dir="${1}"
    local -r max_depth_flag="${2:-3}" # Default max depth is 3 if not provided

    local modules_path
    modules_path=$(get_terraform_module_paths "${modules_dir}" "${max_depth_flag}")
    if [[ -z "${modules_path}" ]]; then
        log "No Terraform modules found in the directory: ${modules_dir}"
        return 0
    fi

    echo "${modules_path}" | while read -r module_path; do
        log "Generating Terraform docs for module at: ${module_path}"
        (cd "${module_path}" && terraform-docs md . > README.md)
    done

    log "Terraform docs generated successfully"
}

# Main entry point
main() {
    local modules_dir=""
    local max_depth_flag=""
    local action=""


    # Print the arguments received
    log "Arguments received: $*"

    while (( $# )); do
        case "$1" in
            --modules_dir=*)
                modules_dir="${1#*=}"
                shift
                ;;
            --max_depth_flag=*)
                max_depth_flag="${1#*=}"
                shift
                ;;
            --action=*)
                action="${1#*=}"
                shift
                ;;
            *)
                log "Error: Invalid argument."
                exit 1
        esac
    done

    if [[ -z "${modules_dir}" ]]; then
        log "Error: Missing mandatory argument '--modules_dir'."
        exit 1
    fi

    if [[ -z "${action}" ]]; then
        action="${DEFAULT_ACTION}"
    fi

    # Check the action and then run the function
    case "${action}" in
        upgrade)
            upgrade_terraform_modules "${modules_dir}" "${max_depth_flag}"
            ;;
        docs)
            generate_terraform_docs "${modules_dir}" "${max_depth_flag}"
            ;;
        *)
            log "Error: Invalid action."
            exit 1
    esac
}

main "$@"
