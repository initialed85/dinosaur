#!/usr/bin/env bash

set -e

while true; do

  if test -e "/tmp/.shutdown"; then
    echo "Shutting down..."
    exit 0
  fi

  if ! find cmd/ -type f | grep -E '\w+'; then
    echo -e "Waiting for source files..."

    while ! find cmd/ -type f | grep -E '\w+'; do
      sleep 1
    done

    echo ""
  fi

  if ! /watch.sh; then
    echo "Watch failed..."
    sleep 1
  fi

done
