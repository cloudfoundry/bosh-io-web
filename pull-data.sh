#!/bin/bash

set -e -x

echo "Fetching"
cf app bosh-hub|grep '#'|cut -f1 -d' '|cut -f2 -d'#'|xargs -P 100 -I{} cf ssh -i {} bosh-hub -T -c 'cd app && ./git-fetch.sh'

echo "Pulling"
cf app bosh-hub|grep '#'|cut -f1 -d' '|cut -f2 -d'#'|xargs -P 100 -I{} cf ssh -i {} bosh-hub -T -c 'cd app && ./git-pull.sh'
