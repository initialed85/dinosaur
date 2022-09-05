#!/usr/bin/env bash

set -e -x

function shutdown() {
  touch /tmp/.shutdown
}
trap shutdown SIGTERM

_="${BASE_FOLDER_PATH:?BASE_FOLDER_PATH env var missing}"
_="${BUILD_CMD:?BUILD_CMD env var missing}"
_="${RUN_CMD:?RUN_CMD env var missing}"

while true; do

  if test -e "/tmp/.shutdown"; then
    exit 0
  fi

  gotty \
    --address 0.0.0.0 \
    --port "${GOTTY_PORT:?PORT env var missing}" \
    --path "${GOTTY_PATH:?PATH env var missing}" \
    --ws-origin '.*' \
    /loop.sh

done
