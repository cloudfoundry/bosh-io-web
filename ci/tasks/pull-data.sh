#!/bin/bash

set -eu

export CF_HOME=/tmp/bosh-io-web-pull-data-$$
./bin/configure-cf

app=bosh-hub-new

echo "Pulling..."
cf app "$app"|grep '#'|cut -f1 -d' '|cut -f2 -d'#'|xargs -P 100 -I{} cf ssh -i {} "$app" -T -c 'cd app && ./bin/git-pull.sh'
