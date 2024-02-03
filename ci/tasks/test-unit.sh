#!/bin/bash

set -e

cd web

go run github.com/onsi/ginkgo/v2/ginkgo \
  -r \
  --keep-going \
  --race \
  --randomize-all \
  --randomize-suites \
  .
