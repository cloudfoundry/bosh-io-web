#!/bin/bash

cd "$( dirname "${BASH_SOURCE[0]}" )/.."

fly --target be \
  set-pipeline \
  --pipeline bosh-io-web \
  --config ci/pipeline.yml \
  --load-vars-from <( lpass show --fixed-strings --notes "bosh-io/web ci/pipeline.yml secrets (VMware IT PCF - formerly known as PCF Photon)" )
