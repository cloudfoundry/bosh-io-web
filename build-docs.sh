#!/bin/bash

set -e

if [ ! -d ./docs-bosh ]; then
  echo 'Automatically cloning docs-bosh...'
  git clone --recurse-submodules https://github.com/cloudfoundry/docs-bosh.git
fi

docker run --rm -it \
  -v "${PWD}/docs-bosh:/docs" \
  -v "${PWD}/templates/docs:/site" \
  squidfunk/mkdocs-material:2.7.2 \
  build --site-dir=/site
