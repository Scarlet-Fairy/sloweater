#!/usr/bin/env bash

DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" &> /dev/null && pwd )"

nomad agent \
  -dev-connect \
  -log-level INFO \
  -config "$DIR"/nomad-config-dev.hcl