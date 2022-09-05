#!/usr/bin/env bash

set -e

while true; do

  if test -e "/tmp/.shutdown"; then
    exit 0
  fi

  if ! find . -type f | grep -E '\w+'; then
    echo -e "Waiting for source files..."

    while ! find . -type f | grep -E '\w+'; do
      sleep 1
    done

    echo ""
  fi

  if ! find . -type f | entr -n -r -a -c -s 'bash -c "/build.sh && /run.sh"'; then
    sleep 1
  fi

  break

done
