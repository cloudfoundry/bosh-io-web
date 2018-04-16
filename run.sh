#!/bin/bash

set -e

# This script is called by the buildpack
# (https://github.com/shageman/buildpack-binary)

configPath=$1
debug=$2

if [ -z "$configPath" ]; then
  configPath=prod-conf/web.json
  ./git-init-clone.sh
fi

export PATH=/usr/local/bin:/usr/bin:/bin:/app/bin:$PATH

if [ -z "$debug" ]; then
  # Martini will cache compiled templates
  export MARTINI_ENV=production
fi

chmod +x ./bosh-hub

exec ./bosh-hub -configPath $configPath $debug
