#!/bin/bash

export GO111MODULE=off

go build -o bosh-hub github.com/bosh-io/web/main

exec ./run.sh conf/local.json -debug
