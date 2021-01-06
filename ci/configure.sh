#!/bin/bash

cd "$( dirname "${BASH_SOURCE[0]}" )/.."

fly --target be \
  set-pipeline \
  --pipeline bosh-io-web \
  --config ci/pipeline.yml
