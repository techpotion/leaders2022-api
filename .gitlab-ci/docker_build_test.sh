#!/usr/bin/env bash

set +x
set -e

RED="\e[31m"
GREEN="\033[32m"
YELLOW="\033[1;33m"
BLUE="\033[34m"
CLEAR="\e[0m"

TAG="${CI_PROJECT_NAME}:test-latest"

echo -e "------------------------------------------------------------------------------------------"
echo -e " "
echo -e "[ " "${BLUE}--- Building the image from Dockerfile ---${CLEAR}" " ]"
echo -e " "

docker build \
    --no-cache --pull -t ${TAG} .

echo -e " "
echo -e "------------------------------------------------------------------------------------------"
