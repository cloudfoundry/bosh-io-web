#!/bin/bash

set -e

echo "Configuring go"
export GOARCH=amd64
export GOOS=linux
export GO111MODULE=off
export GOTOOLDIR=$(go env GOROOT)/pkg/linux_amd64
export GOPATH="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && cd ../../../../.. && pwd )"

echo "Building bosh-hub"
go build -o bosh-hub github.com/bosh-io/web/main
