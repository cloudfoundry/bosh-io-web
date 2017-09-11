#!/bin/bash

set -e -x

echo "Fetching"

cd prod-conf/releases
git fetch origin

cd ../releases-index
git fetch origin

cd ../stemcells-legacy-index
git fetch origin
