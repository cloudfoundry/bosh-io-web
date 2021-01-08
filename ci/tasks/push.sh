#!/bin/bash

set -eu

cp -rp docroot src/github.com/bosh-io/web/templates/docs

cd src/github.com/bosh-io/web

# Application names
old=bosh-hub-old
new=bosh-hub-new
curr=bosh-hub

mkdir -p prod-conf
echo "$WEB_CONFIG" > prod-conf/web.json
cat prod-conf/web.json
exit 0

export CF_HOME=/tmp/bosh-io-web-push-$$
./bin/configure-cf

./bin/build.sh

echo "Pushing to cf"
if cf app $new; then echo "$new must not exist"; exit 1; fi
if cf app $old; then echo "$old must not exist"; exit 1; fi

echo "Pushing new version"
cf push $new -i 10 -k 2G -m 1G -b https://github.com/shageman/buildpack-binary
rm prod-conf/web.json

echo "Testing new version"
./bin/test-server "https://$new.sc2-04-pcf1-apps.oc.vmware.com"

echo "Mapping routes to new version"
cf unmap-route $new sc2-04-pcf1-apps.oc.vmware.com -n $new
cf map-route $new sc2-04-pcf1-apps.oc.vmware.com -n $curr
cf map-route $new sc2-04-pcf1-apps.oc.vmware.com -n bosh
cf map-route $new cloudfoundry.org -n bosh
cf map-route $new bosh.io -n test               #TODO: remove after migration
cf map-route $new bosh.io -n www
cf map-route $new bosh.io

echo "Swapping version: current->old"

cf rename $curr $old

echo "Unmapping routes from old version"
cf map-route $old sc2-04-pcf1-apps.oc.vmware.com -n $old
cf unmap-route $old sc2-04-pcf1-apps.oc.vmware.com -n $curr
cf unmap-route $old sc2-04-pcf1-apps.oc.vmware.com -n bosh
cf unmap-route $old cloudfoundry.org -n bosh
cf unmap-route $old bosh.io -n test   #TODO: remove after migration
cf unmap-route $old bosh.io -n www
cf unmap-route $old bosh.io

echo "Swapping version: current->old, new->current"
cf rename $new $curr

echo "Testing current version"
./bin/test-server "https://test.bosh.io"  #TODO: change to bosh.io after migration

echo "Deleting old version"
cf delete -f $old

rm -fr ./prod-conf
rm -fr "$CF_HOME"
