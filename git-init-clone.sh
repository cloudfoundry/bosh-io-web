#!/bin/bash

set -e -x

echo "Clone data source"
git clone https://github.com/bosh-io/releases               prod-conf/releases
git clone https://github.com/bosh-io/releases-index         prod-conf/releases-index
git clone https://github.com/bosh-io/stemcells-legacy-index prod-conf/stemcells-legacy-index

(
  echo "In 1s"
  sleep 1

  echo "Warm up"
	time curl -I -q http://localhost:8080
) &

disown
