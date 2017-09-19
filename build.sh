#!/bin/bash

set -e

echo "Configuring go"
export GOARCH=amd64
export GOOS=linux
export GOTOOLDIR=$(go env GOROOT)/pkg/linux_amd64

echo "Building bosh-hub"
go build -o bosh-hub github.com/bosh-io/web/main

echo "Building docs"
./build-docs.sh
