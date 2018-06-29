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

export CF_HOME=/tmp/bosh-io-web-push-$$
./bin/configure-cf

./bin/build.sh

echo "Pushing to cf"
if cf app $new; then echo "$new must not exist"; exit 1; fi
if cf app $old; then echo "$old must not exist"; exit 1; fi

echo "Pushing new version"
cf push $new -i 5 -k 2G -b https://github.com/shageman/buildpack-binary
rm prod-conf/web.json

echo "Testing new version"
./bin/test-server "https://$new.cfapps.io"

echo "failing because you're supposed to be testing"
exit 1

echo "Mapping routes to new version"
cf unmap-route $new cfapps.io -n $new
cf map-route $new cfapps.io -n $curr
cf map-route $new cfapps.io -n bosh
cf map-route $new cloudfoundry.org -n bosh
cf map-route $new bosh.io -n www
cf map-route $new bosh.io

echo "Swapping version: current->old"

cf rename $curr $old

echo "Unmapping routes from old version"
cf map-route $old cfapps.io -n $old
cf unmap-route $old cfapps.io -n $curr
cf unmap-route $old cfapps.io -n bosh
cf unmap-route $old cloudfoundry.org -n bosh
cf unmap-route $old bosh.io -n www
cf unmap-route $old bosh.io

echo "Swapping version: current->old, new->current"
cf rename $new $curr

echo "Testing current version"
./bin/test-server "https://bosh.io"

echo "Deleting old version"
cf delete $old

rm -fr ./prod-conf
rm -fr "$CF_HOME"
