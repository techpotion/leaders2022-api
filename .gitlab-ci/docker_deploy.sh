#!/usr/bin/env bash

set +x
set -e

RED="\e[31m"
GREEN="\033[32m"
YELLOW="\033[1;33m"
BLUE="\033[34m"
CLEAR="\e[0m"

TAG="${CI_PROJECT_NAME}:latest"

echo -e "------------------------------------------------------------------------------------------"
echo -e " "
echo -e "[ " "${BLUE}--- Deploying ---${CLEAR}" " ]"
echo -e " "

CONTAINER_NAME=${APP}_${PORT}

echo -e " "
echo -e "[ " "${BLUE}Stopping and removing existing "${CONTAINER_NAME}" container...${CLEAR}" " ]"
echo -e " "

docker stop ${CONTAINER_NAME} > /dev/null || true;
docker rm -f ${CONTAINER_NAME} > /dev/null || true;

echo -e " "
echo -e "[ " "${BLUE}Starting new ${CI_PROJECT_NAME} container on port ${PORT}...${CLEAR}" " ]"
echo -e " "

docker run \
    -d \
    --name "${CONTAINER_NAME}" \
    --restart always \
    -p ${PORT}:${PORT} \
    -e ENV=prod \
    -e PORT=${PORT} \
    -e HTTP_SERVER_PORT=${PORT} \
    -e DB_POSTGRESQL_URI=${DB_POSTGRESQL_URI} \
    -e MODEL_MICROSERVICE_URI=${MODEL_MICROSERVICE_URI} \
    -e MODEL_MICROSERVICE_PREDICT_ENDPOINT=${MODEL_MICROSERVICE_PREDICT_ENDPOINT} \
    -e PLOT_MICROSERVICE_URI=${PLOT_MICROSERVICE_URI} \
    -e PLOT_MICROSERVICE_EFFICIENCY_ENDPOINT=${PLOT_MICROSERVICE_EFFICIENCY_ENDPOINT} \
    ${TAG} > /dev/null;

echo -e " "
echo -e "[ " "${GREEN}Successfully deployed${CLEAR}" " ]"
echo -e " "

echo -e " "
echo -e "------------------------------------------------------------------------------------------"
