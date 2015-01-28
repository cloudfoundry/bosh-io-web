#!/bin/bash

docs=docs-bosh/
conf=prod-conf/

if [ -d $docs ]; then
  echo "Repo $docs already exists"
  exit 1
fi

git clone -b master https://github.com/cppforlife/docs-bosh $docs

if [ -d $conf ]; then
  echo "Repo $conf already exists"
  exit 1
fi

git clone -b master git@github.com:cppforlife/bosh-hub-prod-conf.git $conf
