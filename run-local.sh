#!/bin/bash

set -e

config=$1

if [ -z "$config" ]; then
  config=conf/local.json
fi

go build -o bosh-hub github.com/bosh-io/web/main

if [ ! -f $config ]; then
  echo "Missing $config file"
  exit 1
fi

exec ./run.sh $config -debug
