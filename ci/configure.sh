#!/bin/bash

cd "$( dirname "${BASH_SOURCE[0]}" )/.."

fly --target production \
  set-pipeline \
  --pipeline web \
  --config ci/pipeline.yml \
  --load-vars-from <( lpass show --fixed-strings --notes "bosh-io/web ci/pipeline.yml secrets" )
