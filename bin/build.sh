#!/bin/bash

set -e

echo "Configuring go"
export GOARCH=amd64
export GOOS=linux

echo "Building bosh-hub"
go build -o bosh-hub github.com/bosh-io/web/main
