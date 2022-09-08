#!/usr/bin/env bash

set -e

cd "${BASE_FOLDER_PATH:?BASE_FOLDER_PATH env var missing}"

echo ""
echo "HOSTNAME=${HOSTNAME}"
echo "LOCAL_IP=${LOCAL_IP}"
echo "BROADCAST_IP=${BROADCAST_IP}"
echo ""

echo -e "\nRunning..."

${RUN_CMD:?RUN_CMD env var missing}
