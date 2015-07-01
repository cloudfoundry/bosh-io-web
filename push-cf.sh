#!/bin/bash

set -e

# Application names
old=bosh-hub-old
new=bosh-hub-new
curr=bosh-hub

if [ ! -f prod-conf/web.json ]; then
  echo 'Missing web.json config'
  exit 1
fi

./build.sh

echo "Generating new assets ID"
echo `date|md5` > ./public/assets-id

echo "Pushing to cf"
if cf app $new; then echo "$new must not exist"; exit 1; fi
if cf app $old; then echo "$old must not exist"; exit 1; fi

echo "Pushing new version"
cf push $new -i 2 -k 2G -b https://github.com/ddollar/heroku-buildpack-multi.git

read -p "Map routes to new version (y/n)? " CONT
if [ "$CONT" != "y" ]; then
  echo "Exiting. $new is running. Delete before redeploying."
  exit 1
fi

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

read -p "Delete old version (y/n) ? " CONT
if [ "$CONT" != "y" ]; then
  echo "Exiting. $old is running. Delete before redeploying."
  exit 1
fi

echo "Deleting old version"
cf delete $old
