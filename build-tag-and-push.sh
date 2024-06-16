#!/bin/bash

set -e

docker build --progress plain --platform=linux/amd64 -t kube-registry:5000/dinosaur-frontend:latest -f ./docker/frontend/Dockerfile .
docker build --progress plain --platform=linux/amd64 -t kube-registry:5000/dinosaur-backend:latest -f ./docker/backend/Dockerfile .
docker build --progress plain --platform=linux/amd64 -t kube-registry:5000/dinosaur-session:latest -f ./docker/session/Dockerfile ./docker/session

docker image push kube-registry:5000/dinosaur-frontend:latest
docker image push kube-registry:5000/dinosaur-backend:latest
docker image push kube-registry:5000/dinosaur-session:latest
