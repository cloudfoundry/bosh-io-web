#!/bin/bash

set -e -x

echo "Fetching"

set +x
PKEY=/tmp/pkey
echo "${GIT_CLONE_KEY}" > "${PKEY}"
echo >> "${PKEY}"
chmod 0400 "${PKEY}"
set -x

export GIT_SSH_COMMAND="ssh -i ${PKEY} -o StrictHostKeyChecking=no"

for dir in \
  prod-conf/releases \
  prod-conf/releases-index \
  prod-conf/stemcells-legacy-index \
  prod-conf/stemcells-core-index \
  prod-conf/stemcells-cpi-index \
  prod-conf/stemcells-windows-index \
  prod-conf/stemcells-softlayer-index \
  prod-conf/stemcells-alicloud-index
do
  pushd "$dir"
  git pull
  popd > /dev/null
done

pkill -SIGHUP bosh-hub

sleep 1

echo "Warm up"
time curl -I -q http://localhost:8080/releases
time curl -I -q http://localhost:8080/stemcells
