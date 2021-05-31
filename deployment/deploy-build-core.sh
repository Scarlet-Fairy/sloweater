#!/usr/bin/env bash

DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" &> /dev/null && pwd )"

export NOMAD_ADDR="http://192.168.44.25:4646"

nomad job run "$DIR"/core.hcl
