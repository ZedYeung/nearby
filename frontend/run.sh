#!/bin/bash
set -a
source .env.development
echo "build..."
docker build -t ${DOCKER_REGISTRY}/${DOCKER_USER}/${APP}_frontend .
echo "run..."
docker run -it --rm -e "BACKEND=${BACKEND}" -p ${PORT}:80 ${DOCKER_REGISTRY}/${DOCKER_USER}/${APP}_frontend