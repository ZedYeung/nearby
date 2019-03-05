#!/bin/bash
docker-compose build
docker-compose push
kompose convert --stdout | kubectl apply -f -