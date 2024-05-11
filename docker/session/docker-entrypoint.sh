#!/usr/bin/env bash

set -e

function shutdown() {
  touch /tmp/.shutdown
}
trap shutdown SIGTERM

_="${SESSION_UUID:?SESSION_UUID env var missing}"
_="${BASE_FOLDER_PATH:?BASE_FOLDER_PATH env var missing}"
_="${BUILD_CMD:?BUILD_CMD env var missing}"
_="${RUN_CMD:?RUN_CMD env var missing}"

LOCAL_IP=$(ip addr show eth0 | grep inet | head -n 1 | xargs | cut -d ' ' -f 2 | cut -d '/' -f 1)
BROADCAST_IP=$(ip addr show eth0 | grep inet | head -n 1 | xargs | cut -d ' ' -f 4)

export SESSION_UUID
export LOCAL_IP
export BROADCAST_IP

screen -a -A -S session -d -m bash -c '/loop.sh'

bash

# while true; do

#   if test -e "/tmp/.shutdown"; then
#     exit 0
#   fi

#   gotty \
#     --address 0.0.0.0 \
#     --port "${GOTTY_PORT:?PORT env var missing}" \
#     --path "${GOTTY_PATH:?PATH env var missing}" \
#     --ws-origin '.*' \
#     screen -a -A -x session

# done
