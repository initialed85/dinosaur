#!/bin/bash

set -e

find . -type f | entr -n -r -a -c -s 'bash -c "/build.sh && /run.sh"'
