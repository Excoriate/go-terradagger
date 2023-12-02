#!/usr/bin/env bash
#
##################################################################################
# Script Name: Docker Command Manager
#
# Author: Alex Torres (github.com/Excoriate), alex_torres@outlook.com
#
# Usage: ./script.sh --image=image-name [--rebuild=rebuild-flag] [--action=action] [--dockerfile=docker-file-name]
#
# Description: This bash script provides functionality to build and/or run Docker images.
#              It's designed to be used as a simple CLI tool to manage Docker workflows.
#
# Parameters:
#    --image:        The name of the Docker image [mandatory].
#    --rebuild:      If true, it forces the rebuild of the Docker image if it already exists. Default is "false" [optional].
#    --action:       Can be one of "build", "run", or "all". Default is "all" indicating to build and run the image [optional].
#    --dockerfile :  Dockerfile to be used, default is "Dockerfile" [optional].
#
# Examples:
#    Build and Run: ./script.sh --image=my-app --rebuild=true --action=all --dockerfile=myDockerfile
#    Build:         ./script.sh --image=my-app --rebuild=true --action=build
#    Run:           ./script.sh --image=my-app --rebuild=false --action=run
#
# Note: Currently, the script supports Docker images that are compatible with linux/amd64 and linux/arm64.
#
# For further details and support, contact the author.
#
##################################################################################
set -euo pipefail

# Constants
readonly DEFAULT_ACTION="all"
readonly DEFAULT_DOCKER_FILE_NAME="Dockerfile"

# Log a message
log() {
    local -r msg="${1}"
    echo "${msg}"
}

# Build Docker image
build_image() {
    local -r image_name="${1}"
    local -r rebuild_flag="${2}"
    local -r docker_file_name="${3}"

    if [[ ! -f "${docker_file_name}" ]]; then
        log "Docker file ${docker_file_name} does not exist. Exiting..."
        exit 1
    fi

    if [[ $(docker images -q "${image_name}" 2> /dev/null) != "" && "${rebuild_flag}" = "false" ]]; then
        log "Docker image ${image_name} already exists. Skipping docker build..."
    else
        log "Building docker image ${image_name} from Dockerfile ${docker_file_name}..."
        docker build -t "${image_name}" -f "${docker_file_name}" .
    fi
}

# Run Docker image
run_image() {
    local -r image_name="${1}"
    log "Running docker image ${image_name} ..."
    docker run -it --rm -v "$(pwd)":/app -w /app "${image_name}"
}

# Main entry point
main() {
    # parse parameters
    local image_name=""
    local rebuild_flag="false"
    local action="${DEFAULT_ACTION}"
    local docker_file_name="${DEFAULT_DOCKER_FILE_NAME}"

    # Print the arguments received
    log "Arguments received: $*"

    while (( "$#" )); do
        case "$1" in
            --image=*)
                image_name="${1#*=}"
                shift
                ;;
            --rebuild=*)
                rebuild_flag="${1#*=}"
                shift
                ;;
            --action=*)
                action="${1#*=}"
                shift
                ;;
            --dockerfile=*)
                docker_file_name="${1#*=}"
                shift
                ;;
            *)
                echo "Error: Invalid argument."
                exit 1
        esac
    done

    if [[ -z "${image_name}" ]]; then
        echo "Error: The image name is mandatory."
        exit 1
    fi

    log "Running with image-name='${image_name}', rebuild-flag='${rebuild_flag}', action='${action}', docker-file-name='${docker_file_name}'"

    case "${action}" in
        build)
          build_image "${image_name}" "${rebuild_flag}" "${docker_file_name}"
          ;;
        run)
          run_image "${image_name}"
          ;;
        all)
          build_image "${image_name}" "${rebuild_flag}" "${docker_file_name}"
          run_image "${image_name}"
          ;;
        *)
          echo "Error: Invalid action '${action}'. The supported actions are: 'build', 'run', or 'all'"
          exit 1
    esac
}

main "$@"
