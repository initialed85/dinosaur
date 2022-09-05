#!/usr/bin/env bash

set -e

cd "${BASE_FOLDER_PATH:?BASE_FOLDER_PATH env var missing}"

echo -e "\nRunning..."

${RUN_CMD:?RUN_CMD env var missing}
