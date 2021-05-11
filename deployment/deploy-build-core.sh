#!/usr/bin/env bash

DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" &> /dev/null && pwd )"

nomad job run "$DIR"/rabbitmq.hcl
nomad job run "$DIR"/registry.hcl