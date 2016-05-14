#!/bin/bash

# This script is called by the buildpack
# (https://github.com/shageman/buildpack-binary)

configPath=$1
assetsID=$2
privateToken=$3
debug=$4

if [ -z "$configPath" ]; then
  configPath=prod-conf/web.json
fi

if [ -z "$assetsID" ]; then
  assetsID=$(cat prod-conf/assets-id)
fi

if [ -z "$privateToken" ]; then
  privateToken=$(cat prod-conf/private-token)
fi

export PATH=/usr/local/bin:/usr/bin:/bin:/app/bin:$PATH

# Make bosh-blostore-s3 available
export PATH=$PWD:$PATH

# Make bosh_cli available
export PATH=$PWD/vendor/bundle/bin:$PATH

if [ -z "$debug" ]; then
  # Martini will cache compiled templates
  export MARTINI_ENV=production
fi

chmod +x ./bosh-hub

exec ./bosh-hub -configPath $configPath -assetsID "$assetsID" -privateToken "$privateToken" $debug
