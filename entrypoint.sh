#!/bin/sh
set -o errexit
set -o nounset
set -o pipefail

# Print all environment variables for debugging (optional, remove in production)
# echo "Debug: Received Environment Variables:"
# env | sort

docker_run_command=(
    docker run
    -e "INPUT_HOST=${INPUT_HOST:-}"
    -e "INPUT_PROTOCOL=${INPUT_PROTOCOL:-}"
    -e "INPUT_PORT=${INPUT_PORT:-}"
    -e "INPUT_USER=${INPUT_USER:-}"
    -e "INPUT_PASSWORD=${INPUT_PASSWORD:-}"
    -e "INPUT_PKEY=${INPUT_PKEY:-}"
    -e "INPUT_PASSPHRASE=${INPUT_PASSPHRASE:-}"
    -e "INPUT_GITHUB_REPO=${INPUT_GITHUB_REPO:-}"
    -e "INPUT_STARTUP_SCRIPT=${INPUT_STARTUP_SCRIPT:-}"
    -e "INPUT_API_PORT=${INPUT_API_PORT:-}"
    -e "INPUT_ENV=${INPUT_ENV:-}"

    # Remove the interactive flags if running in a CI environment
    -i -t
    
    # Your Golang SSH client image name
    do-nginx
)

# Execute the docker run command
echo "Executing: ${docker_run_command[*]}"
"${docker_run_command[@]}"