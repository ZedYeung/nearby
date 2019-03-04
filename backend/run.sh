#!/bin/bash
# -a is equivalent to allexport
set -a
source .env
rm -rf docker-compose.yml
envsubst < docker-compose.yml.template > docker-compose.yml
docker-compose build
docker-compose up