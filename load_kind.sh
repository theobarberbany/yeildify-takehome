#!/usr/bin/env bash
set -euo pipefail

# source utils
. scripts/util.sh --source-only

echo "Tag:" $1

echo docker build -f $(get_root_path)/docker/Dockerfile.reverse . -t tb15/reverse:${1}
docker build -f $(get_root_path)/docker/Dockerfile.reverse . -t tb15/reverse:${1}

echo docker build -f $(get_root_path)/docker/Dockerfile.random . -t tb15/random:${1}
docker build -f $(get_root_path)/docker/Dockerfile.random . -t tb15/random:${1}


kind load docker-image tb15/reverse:$1
kind load docker-image tb15/random:$1

# deploy ?
