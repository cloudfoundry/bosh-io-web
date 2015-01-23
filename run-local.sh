#!/bin/bash

set -e

go build -o bosh-hub github.com/cppforlife/bosh-hub/main

if [ ! `which bosh` ]; then
  echo 'Missing `bosh` executable on PATH'
  exit 1
fi

config=prod-conf/local.json

if [ ! -f $config ]; then
  config=conf/local.json

  if [ ! -f $config ]; then
    echo "Missing $config file"
	exit 1
  fi
fi

exec ./run.sh $config -debug
