#!/bin/bash

set -eu

cp -rp docroot src/github.com/bosh-io/web/templates/docs

cd src/github.com/bosh-io/web

# Application names
old=bosh-io-old
new=bosh-io-new
curr=bosh-io

mkdir -p prod-conf
echo "$WEB_CONFIG" | jq --arg private_key "${PRIVATE_KEY}" '.Repos.ReleaseTarballLinker.PrivateKey = $private_key' > prod-conf/web.json

export CF_HOME=/tmp/bosh-io-web-push-$$
./bin/configure-cf

./bin/build.sh

echo "Pushing to cf"
if cf app $new; then echo "$new must not exist"; exit 1; fi
if cf app $old; then echo "$old must not exist"; exit 1; fi

echo "Pushing new version"
cf push $new -i 10 -k 2G -m 1G -b binary_buildpack -c './run.sh'
rm prod-conf/web.json

echo "Testing new version"
./bin/test-server "https://$new.de.a9sapp.eu"

echo "Mapping routes to new version"
cf unmap-route $new de.a9sapp.eu -n $new
cf map-route $new de.a9sapp.eu -n $curr
cf map-route $new bosh.cloudfoundry.org
cf map-route $new www.bosh.io
cf map-route $new bosh.io

echo "Swapping version: current->old"

cf rename $curr $old

echo "Unmapping routes from old version"
cf map-route $old de.a9sapp.eu -n $old
cf unmap-route $old de.a9sapp.eu -n $curr
cf unmap-route $old bosh.cloudfoundry.org
cf unmap-route $old www.bosh.io
cf unmap-route $old bosh.io

echo "Swapping version: current->old, new->current"
cf rename $new $curr

echo "Testing current version"
./bin/test-server "https://bosh.io"

echo "Deleting old version"
cf delete -f $old

rm -fr ./prod-conf
rm -fr "$CF_HOME"
