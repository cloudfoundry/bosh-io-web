#!/bin/bash

set -e -x

echo "Fetching"

cd prod-conf/releases
git pull

cd ../releases-index
git pull

cd ../stemcells-legacy-index
git pull

pkill -SIGHUP bosh-hub

sleep 1

echo "Warm up"
time curl -I -q http://localhost:8080
