#!/bin/bash

set -e -x

clone_dir="${clone_dir:-prod-conf}"

echo "Clone data source"
git clone https://github.com/cloudfoundry/bosh-io-releases                  $clone_dir/releases
git clone https://github.com/cloudfoundry/bosh-io-releases-index            $clone_dir/releases-index
git clone https://github.com/cloudfoundry/bosh-io-stemcells-legacy-index    $clone_dir/stemcells-legacy-index
git clone https://github.com/cloudfoundry/bosh-io-stemcells-core-index      $clone_dir/stemcells-core-index
git clone https://github.com/cloudfoundry/bosh-io-stemcells-cpi-index       $clone_dir/stemcells-cpi-index
git clone https://github.com/cloudfoundry/bosh-io-stemcells-windows-index   $clone_dir/stemcells-windows-index
git clone https://github.com/cloudfoundry/bosh-io-stemcells-softlayer-index $clone_dir/stemcells-softlayer-index
git clone https://github.com/cloudfoundry-incubator/stemcells-alicloud-index $clone_dir/stemcells-alicloud-index

(
  echo "In 1s"
  sleep 1

  echo "Warm up"
	time curl -I -q http://localhost:8080/releases
	time curl -I -q http://localhost:8080/stemcells
) &

disown
