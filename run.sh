#!/usr/bin/env bash

set -e

function shutdown() {
  docker compose -f docker/docker-compose.yml down --remove-orphans --volumes || true
}
trap shutdown EXIT

docker compose -f docker/docker-compose.yml up -d --build
docker compose -f docker/docker-compose.yml logs -f -t
