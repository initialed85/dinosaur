#!/usr/bin/env bash

set -e

function shutdown() {
  docker compose -f docker/docker-compose.yml down --remove-orphans --volumes || true
}
trap shutdown EXIT

# not strictly necessary because the SessionManager does it as well, just here so the developer can watch it happen (rather than
# wonder why the backend isn't handling requests yet
docker build -t dinosaur-session -f docker/session/Dockerfile ./docker/session

docker compose -f docker/docker-compose.yml up -d --build
docker compose -f docker/docker-compose.yml logs -f -t
