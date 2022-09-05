#!/usr/bin/env bash

set -e

cd "${BASE_FOLDER_PATH:?BASE_FOLDER_PATH env var missing}"

echo -e "Building..."

${BUILD_CMD:?BUILD_CMD env var missing}
