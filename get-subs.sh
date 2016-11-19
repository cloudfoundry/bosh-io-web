#!/bin/bash

docs=docs-bosh/

if [ -d $docs ]; then
  echo "Repo $docs already exists"
  exit 1
fi

git clone -b master https://github.com/cloudfoundry/docs-bosh $docs
