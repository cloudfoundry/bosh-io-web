#!/bin/bash

set -e -x

echo "Clone data source"
git clone https://github.com/bosh-io/releases                  prod-conf/releases
git clone https://github.com/bosh-io/releases-index            prod-conf/releases-index
git clone https://github.com/bosh-io/stemcells-legacy-index    prod-conf/stemcells-legacy-index
git clone https://github.com/bosh-io/stemcells-core-index      prod-conf/stemcells-core-index
git clone https://github.com/bosh-io/stemcells-cpi-index       prod-conf/stemcells-cpi-index
git clone https://github.com/bosh-io/stemcells-windows-index   prod-conf/stemcells-windows-index
git clone https://github.com/bosh-io/stemcells-softlayer-index prod-conf/stemcells-softlayer-index

(
  echo "In 1s"
  sleep 1

  echo "Warm up"
	time curl -I -q http://localhost:8080/releases
	time curl -I -q http://localhost:8080/stemcells
) &

disown
