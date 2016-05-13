#!/bin/bash

# This script is called by the buildpack
# (https://github.com/shageman/buildpack-binary)

config=$1
debug=$2

if [ -z "$config" ]; then
  config=prod-conf/web.json
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

# Generate assets-id only if it has not been already generated
if [ ! -f "./public/assets-id" ]; then
  echo -n "dev" > ./public/assets-id
fi

chmod +x ./bosh-hub

exec ./bosh-hub -configPath $config $debug
